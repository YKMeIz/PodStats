package podman

import (
	"context"
	"github.com/containers/podman/v3/pkg/bindings"
	"log"
	"os"
)

func Connect() *Context {
	// Connect to Podman socket
	ctx, err := bindings.NewConnection(context.Background(), getSock())
	if err != nil {
		log.Fatalln(err)
	}
	return &Context{Context: ctx, filters: make(map[string][]string)}
}

func getSock() string {
	if _, err := os.Stat("/run/podman/podman.sock"); err == nil {
		log.Println("detect podman.sock under root user")
		return "unix:" + socketRoot
	}

	socketUser := "unix:" + os.Getenv("XDG_RUNTIME_DIR") + "/podman/podman.sock"
	if _, err := os.Stat(socketUser); err == nil {
		log.Println("detect podman.sock under " + socketUser + " user")
		return socketUser
	}

	log.Println("podman.sock is not activated")
	return ""
}
