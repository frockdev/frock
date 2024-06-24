package podman

import (
	"context"
	"github.com/containers/podman/v5/pkg/domain/entities/types"
	"log"
)
import "github.com/containers/podman/v5/pkg/bindings/pods"

func CreatePod(ctx context.Context, projectName string) string {
	podName := projectName + "-pod"
	exists, err := pods.Exists(ctx, podName, nil)
	if err != nil {
		log.Fatal(err)
	}
	if exists {
		remove, err := pods.Remove(ctx, podName, nil)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(remove)
	}
	podSpec := types.PodSpec{}
	podSpec.PodSpecGen.Name = podName
	report, err := pods.CreatePodFromSpec(ctx, &podSpec)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(report)
	return podName
}

func RemovePod() {
	connection := Connect()
	config, err := GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	removeOptions := pods.RemoveOptions{}
	force := true
	removeOptions.Force = &force
	report, err := pods.Remove(connection, config.ProjectName+"-pod", &removeOptions)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(report)
}
