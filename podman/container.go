package podman

import (
	"context"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/bindings/images"
	"github.com/containers/podman/v5/pkg/specgen"
	spec "github.com/opencontainers/runtime-spec/specs-go"
	"log"
	"os"
	"strings"
)

//RunContainer

func RunContainer(ctx context.Context, app App, appName string, projectName string) {
	containerSpec := specgen.SpecGenerator{}
	containerSpec.Name = projectName + "-" + appName
	if app.PrivilegedMode == "enabled" {
		truePrivileged := true
		containerSpec.Privileged = &truePrivileged
		mountInfo := spec.Mount{}
		mountInfo.Source = "/var/run/docker.sock"
		mountInfo.Destination = "/var/run/docker.sock"
		mountInfo.Type = "bind"
		containerSpec.Mounts = append(containerSpec.Mounts, mountInfo)
	}
	if len(app.Command) > 0 {
		if len(app.Command) > 1 {
			containerSpec.Entrypoint = []string{app.Command[0]}
			containerSpec.Command = app.Command[1:]
		}
		if len(app.Command) == 1 {
			containerSpec.Entrypoint = []string{app.Command[0]}
		}
	}
	if app.Image != "" {
		containerSpec.Image = app.Image
	} else {
		log.Fatal("Image for app " + appName + " is not provided")
	}

	if app.Tag != "" {
		containerSpec.Image = containerSpec.Image + ":" + app.Tag
	} else {
		containerSpec.Image = containerSpec.Image + ":latest"
	}

	containerSpec.Pod = projectName + "-pod"
	var labels = make(map[string]string)
	for _, value := range app.Labels {
		splittedLabel := strings.Split(value, "=")
		labels[splittedLabel[0]] = splittedLabel[1]
	}
	containerSpec.Labels = labels
	containerSpec.WorkDir = "/app"
	if app.MountCode != "disabled" {
		containerSpec.Mounts = make([]spec.Mount, 0)
		mountInfo := spec.Mount{}
		if app.MountContainerDir != "" {
			mountInfo.Destination = app.MountContainerDir
		} else {
			mountInfo.Destination = "/app"
		}
		curDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		mountInfo.Source = curDir
		mountInfo.Type = "bind"
		containerSpec.WorkDir = mountInfo.Destination
		containerSpec.Mounts = append(containerSpec.Mounts, mountInfo)
	}

	containerSpec.Env = app.Env
	pull, err := images.Pull(ctx, containerSpec.Image, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(pull)
	report, err := containers.CreateWithSpec(ctx, &containerSpec, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(report)

	err = containers.Start(ctx, containerSpec.Name, nil)
	if err != nil {
		log.Println(err)
	}

}

func RunBox(ctx context.Context, box Box, boxName string, projectName string) {
	boxSpec := specgen.SpecGenerator{}
	boxSpec.Name = projectName + "-" + boxName
	if box.PrivilegedMode == "enabled" {
		truePrivileged := true
		boxSpec.Privileged = &truePrivileged
		mountInfo := spec.Mount{}
		mountInfo.Source = "/var/run/docker.sock"
		mountInfo.Destination = "/var/run/docker.sock"
		mountInfo.Type = "bind"
		boxSpec.Mounts = append(boxSpec.Mounts, mountInfo)
	}
	var labels = make(map[string]string)
	for _, value := range box.Labels {
		splittedLabel := strings.Split(value, "=")
		labels[splittedLabel[0]] = splittedLabel[1]
	}
	boxSpec.Labels = labels
	boxSpec.Pod = projectName + "-pod"
	if len(box.Command) > 0 {
		if len(box.Command) > 1 {
			boxSpec.Entrypoint = []string{box.Command[0]}
			boxSpec.Command = box.Command[1:]
		}
		if len(box.Command) == 1 {
			boxSpec.Entrypoint = []string{box.Command[0]}
		}
	}
	if box.Image != "" {
		boxSpec.Image = box.Image
	} else {
		log.Fatal("Image for box " + boxName + " is not provided")
	}
	if box.Tag != "" {
		boxSpec.Image = boxSpec.Image + ":" + box.Tag
	} else {
		boxSpec.Image = boxSpec.Image + ":latest"
	}
	pull, err := images.Pull(ctx, boxSpec.Image, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(pull)
	boxSpec.Env = box.Env
	report, err := containers.CreateWithSpec(ctx, &boxSpec, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(report)

	err = containers.Start(ctx, boxSpec.Name, nil)
	if err != nil {
		log.Println(err)
	}

}
