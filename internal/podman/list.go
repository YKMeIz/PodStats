package podman

import (
	"github.com/containers/podman/v3/pkg/bindings/containers"
)

func (ctx *Context) List() ([]string, error) {
	all := true
	options := containers.ListOptions{
		All:     &all,
		Filters: ctx.filters,
	}

	containerList, err := containers.List(ctx.Context, &options)
	if err != nil {
		return []string{}, err
	}

	var res []string

	for _, v := range containerList {
		res = append(res, v.ID)
	}

	return res, nil
}
