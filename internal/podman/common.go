package podman

import (
	"context"
	"github.com/containers/podman/v3/pkg/domain/entities"
)

const (
	socketRoot = "/run/podman/podman.sock"
)

type Context struct {
	context.Context

	filters map[string][]string

	statsReport            chan entities.ContainerStatsReport
	event                  chan entities.Event
	StatsReportHandlerFunc func(report entities.ContainerStatsReport)
	EventHandlerFunc       func(event entities.Event)
}

func (ctx *Context) SetFilter(key string, value []string) {
	ctx.filters[key] = value
}

func (ctx *Context) SetFilterPodName(podName ...string) {
	ctx.SetFilter("pod", podName)
}

func (ctx *Context) Serve() {
	if ctx.StatsReportHandlerFunc != nil {
		ctx.stats()
	}
	if ctx.EventHandlerFunc != nil {
		ctx.events()
	}

	for {
		select {
		case report := <-ctx.statsReport:
			ctx.StatsReportHandlerFunc(report)
		case event := <-ctx.event:
			ctx.EventHandlerFunc(event)
		}
	}
}
