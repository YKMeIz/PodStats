package main

import (
	"encoding/json"
	"fmt"
	"github.com/YKMeIz/PodStats/internal/osstat"
	"github.com/YKMeIz/PodStats/internal/podman"
	"github.com/YKMeIz/PodStats/internal/ws"
	"github.com/containers/podman/v3/pkg/domain/entities"
	"github.com/containers/podman/v3/utils"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
)

type context struct {
	*ws.Hub
	*podman.Context
	eventReports []eventReport
}

func floatToPercentString(f float64) string {
	strippedFloat, err := utils.RemoveScientificNotationFromFloat(f)
	if err != nil || strippedFloat == 0 {
		// If things go bazinga, return a safe value
		return "--"
	}
	return fmt.Sprintf("%.2f", strippedFloat) + "%"
}

func main() {
	ctx := context{
		Hub:     ws.Run(),
		Context: podman.Connect(),
	}

	log.Println("websocket hub started; podman connected")

	ctx.StatsReportHandlerFunc = ctx.stats
	ctx.EventHandlerFunc = ctx.event

	log.Println("start serve podman")
	go ctx.Serve()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ws" {
			ws.ServeWs(ctx.Hub, w, r)
			return
		}

		if _, err := os.Stat("dist" + r.URL.Path); os.IsNotExist(err) {
			http.ServeFile(w, r, "dist/index.html")
			return
		}
		http.ServeFile(w, r, "dist"+r.URL.Path)
	})

	log.Fatal(http.ListenAndServe(":9023", mux))
}

func (ctx *context) stats(report entities.ContainerStatsReport) {
	sort.Slice(ctx.eventReports, func(i, j int) bool {
		return ctx.eventReports[i].Time > ctx.eventReports[j].Time
	})
	if len(ctx.eventReports) > 10 {
		ctx.eventReports = ctx.eventReports[:10]
	}
	if b, err := json.Marshal(message{
		Type:    "event-reports",
		Content: ctx.eventReports,
	}); err == nil {
		ctx.Broadcast(b)
	}

	if report.Error != nil {
		log.Println(report.Error)
		return
	}

	if sysStat, err := osstat.GetOSStat(); err != nil {
		log.Println(err)
	} else {
		status := "Fully Online"
		var service float64

		if containers, err := ctx.List(); err != nil {
			status = "Unknown"
		} else {
			if len(report.Stats) < len(containers) {
				status = "Partially Online "
			} else if len(report.Stats) == 0 {
				status = "Fully Offline"
			}
			service = float64(len(report.Stats)) / float64(len(containers)) * 100
		}

		b, err := json.Marshal(message{
			Type: "system-report",
			Content: systemReport{
				Cpu:     sysStat.CpuUsage,
				Memory:  sysStat.MemoryUsage,
				Service: service,
				Status:  status,
			},
		})

		if err == nil {
			ctx.Broadcast(b)
		}
	}

	var (
		statusReports    []statusReport
		containerReports []containerReport
	)

	for _, v := range report.Stats {
		sr := statusReport{
			Name: v.Name,
		}
		info, err := ctx.Inspect(v.Name)
		if err == nil {
			sr.Description = formatName(info.ImageName)
			sr.Status = info.State.Status
		}

		statusReports = append(statusReports, sr)

		containerReports = append(containerReports, containerReport{
			Name:         v.Name,
			Cpu:          formatFloat(v.AvgCPU * float64(len(report.Stats))),
			Memory:       formatFloat(v.MemPerc * float64(len(report.Stats))),
			NetworkIO:    ByteSize(v.NetInput) + "/" + ByteSize(v.NetOutput),
			BlockIO:      ByteSize(v.BlockInput) + "/" + ByteSize(v.BlockOutput),
			Created:      info.Created.String(),
			StartedAt:    info.State.StartedAt.String(),
			RestartCount: info.RestartCount,
		})
	}

	if b, err := json.Marshal(message{
		Type:    "status-reports",
		Content: statusReports,
	}); err == nil {
		ctx.Broadcast(b)
	}

	if b, err := json.Marshal(message{
		Type:    "container-reports",
		Content: containerReports,
	}); err == nil {
		ctx.Broadcast(b)
	}
}

func (ctx *context) event(event entities.Event) {
	if event.Type == "container" {
		name := formatName(event.Actor.Attributes["image"])

		if len(ctx.eventReports) > 0 {
			item := ctx.eventReports[len(ctx.eventReports)-1]
			if item.ID == event.Actor.ID && item.Time == event.Time {
				if event.Action == "restart" {
					item.Action = "restarted"
				} else if item.Action != "restarted" && event.Action == "start" {
					item.Action = "started"
				} else if item.Action != "restarted" && event.Action == "stop" {
					item.Action = "stopped"
				} else {
					return
				}
				ctx.eventReports[len(ctx.eventReports)-1] = item
			}
		}

		action := event.Action
		if action == "start" || action == "restart" {
			action += "ed"
		} else if action == "stop" {
			action = "stopped"
		} else {
			return
		}

		ctx.eventReports = append(ctx.eventReports, eventReport{
			ID:     event.Actor.ID,
			Name:   name,
			Action: action,
			Time:   event.Time,
		})

	}
}

func formatName(n string) string {
	res := strings.Split(n, "/")
	return strings.Split(res[len(res)-1], ":")[0]
}

func formatFloat(f float64) float64 {
	strippedFloat, err := utils.RemoveScientificNotationFromFloat(f)
	if err != nil {
		return -1
	}
	return strippedFloat
}
