package podman

import (
	"context"
	"encoding/json"
	"github.com/containers/podman/v3/pkg/bindings"
	"github.com/containers/podman/v3/pkg/bindings/containers"
	"github.com/containers/podman/v3/pkg/domain/entities"
	"log"
	"net/http"
	"reflect"
	"time"
)

func (ctx *Context) stats() {
	go func() {
		stream := true
		interval := 5
		var (
			ids []string
			err error
		)

		for {
			ids, err = ctx.List()
			if err != nil {
				log.Println(err)
				time.Sleep(time.Second)
				continue
			}

			cancel := make(chan bool)

			ctx.statsReport, err = modifiedStats(ctx.Context, ids, cancel, &containers.StatsOptions{
				Stream:   &stream,
				Interval: &interval,
			})
			if err != nil {
				log.Println(err)
				time.Sleep(time.Second)
				continue
			}

			for {
				time.Sleep(5 * time.Second)
				newIDs, err := ctx.List()
				if err != nil {
					log.Println(err)
					continue
				}
				if !reflect.DeepEqual(ids, newIDs) {
					ids = newIDs
					cancel <- true
					break
				}
			}
		}
	}()
}

// Original Stats() does not respect interval setting.
// See https://github.com/containers/podman/blob/v3.4.4/pkg/bindings/containers/containers.go#L221.
func modifiedStats(ctx context.Context, containerIDs []string, cancel chan bool, options *containers.StatsOptions) (chan entities.ContainerStatsReport, error) {
	if options == nil {
		options = new(containers.StatsOptions)
	}
	_ = options
	conn, err := bindings.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	params, err := options.ToParams()
	if err != nil {
		return nil, err
	}
	for _, c := range containerIDs {
		params.Add("containers", c)
	}

	response, err := conn.DoRequest(nil, http.MethodGet, "/containers/stats", params, nil)
	if err != nil {
		return nil, err
	}
	if !response.IsSuccess() {
		return nil, response.Process(nil)
	}

	statsChan := make(chan entities.ContainerStatsReport)

	go func() {
		defer close(statsChan)
		defer response.Body.Close()

		dec := json.NewDecoder(response.Body)
		doStream := true
		if options.Changed("Stream") {
			doStream = options.GetStream()
		}

	streamLabel: // label to flatten the scope
		select {
		case <-response.Request.Context().Done():
			return // lost connection - maybe the server quit
		case <-cancel:
			return // stop this goroutine
		default:
			// fall through and do some work
		}

		var report entities.ContainerStatsReport
		if err := dec.Decode(&report); err != nil {
			report = entities.ContainerStatsReport{Error: err}
		}
		statsChan <- report

		if report.Error != nil || !doStream {
			return
		}

		time.Sleep(time.Duration(options.GetInterval()) * time.Second)

		goto streamLabel
	}()

	return statsChan, nil
}
