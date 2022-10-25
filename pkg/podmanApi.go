package pkg

import (
	"context"
	"fmt"
	"github.com/containers/podman/v4/pkg/bindings"
	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/containers/podman/v4/pkg/bindings/images"
	"os"
)

type PodmanAPI struct {
	connection context.Context
}

func NewPodmanAPI() *PodmanAPI {
	conn, err := bindings.NewConnection(context.Background(), "unix://run/podman/podman.sock")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return &PodmanAPI{
		connection: conn,
	}
}
func (p *PodmanAPI) AddContainer() {
	_, err := images.Pull(p.connection, "quay.io/libpod/alpine_nginx", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	s := NewContainerBuilder("docker.io/nginx").
		GenName().
		Mount(MountTypeBind, "/tmp", "/data").
		Label("generator", "podder").
		Label("reason", "scaleup").
		MapPort(65102, 80).
		BuildSpecification()

	createResponse, err := containers.CreateWithSpec(p.connection, s, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Container created.")
	if err := containers.Start(p.connection, createResponse.ID, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Container started.")
}
