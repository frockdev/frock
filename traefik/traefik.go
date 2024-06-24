package traefik

import (
	"frock2/podman"
	nettypes "github.com/containers/common/libnetwork/types"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/bindings/images"
	"github.com/containers/podman/v5/pkg/specgen"
	spec "github.com/opencontainers/runtime-spec/specs-go"
	"log"
)

// this function will create a traefik container
// with port 80 mapped to host
func UpTraefik() {
	connection := podman.Connect()

	// check if traefik container exists
	// if it does, remove it
	// write a code:
	exists, err := containers.Exists(connection, "traefik-main", nil)
	if err != nil {
		log.Fatal(err)
	}
	if !exists {
		// create a traefik container
		// write a code:
		pull, err := images.Pull(connection, "traefik:v3.0", nil)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(pull)
		containerSpec := specgen.SpecGenerator{}
		containerSpec.Name = "traefik-main"
		containerSpec.Image = "traefik:v3.0"
		containerSpec.Command = []string{"traefik", "--api.insecure=true", "--providers.docker=true", "--entrypoints.web.address=:80"}
		privileged := true
		containerSpec.Privileged = &privileged
		containerSpec.PortMappings = make([]nettypes.PortMapping, 0)
		portMapping := nettypes.PortMapping{}
		portMapping.ContainerPort = 80
		portMapping.HostPort = 80
		portMapping.Protocol = "tcp"

		portMapping2 := nettypes.PortMapping{}
		portMapping2.ContainerPort = 8080
		portMapping2.HostPort = 8080
		portMapping2.Protocol = "tcp"

		containerSpec.PortMappings = append(containerSpec.PortMappings, portMapping)
		containerSpec.PortMappings = append(containerSpec.PortMappings, portMapping2)

		containerSpec.Mounts = make([]spec.Mount, 0)
		mountInfo := spec.Mount{}
		mountInfo.Source = "/var/run/docker.sock"
		mountInfo.Destination = "/var/run/docker.sock"
		mountInfo.Type = "bind"
		containerSpec.Mounts = append(containerSpec.Mounts, mountInfo)

		containerSpec.Labels = map[string]string{}

		report, err := containers.CreateWithSpec(connection, &containerSpec, nil)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(report)
	}
	err = containers.Start(connection, "traefik-main", nil)
	if err != nil {
		log.Fatal(err)
	}

}

func DownTraefik() {
	connection := podman.Connect()
	exists, err := containers.Exists(connection, "traefik-main", nil)
	if err != nil {
		log.Fatal(err)
	}
	_ = containers.Stop(connection, "traefik-main", nil)
	if exists {
		remove, err := containers.Remove(connection, "traefik-main", nil)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(remove)
	}
}
