package podman

import (
	"github.com/containers/podman/v3/libpod/define"
	"github.com/containers/podman/v3/pkg/bindings/containers"
)

func (ctx *Context) Inspect(name string) (*define.InspectContainerData, error) {
	return containers.Inspect(ctx, name, nil)
}
