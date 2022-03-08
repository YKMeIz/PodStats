package podman

import (
	"github.com/containers/podman/v3/pkg/bindings/system"
	"github.com/containers/podman/v3/pkg/domain/entities"
	"log"
)

func (ctx *Context) events() {
	stream := true
	since := "0"
	ctx.event = make(chan entities.Event)

	go func() {
		err := system.Events(ctx, ctx.event, make(chan bool), &system.EventsOptions{
			//Filters: ctx.filters,
			Since:  &since,
			Stream: &stream,
		})

		if err != nil {
			log.Fatalln(err)
		}
	}()
}
