package podman

import (
	"context"
	"github.com/containers/image/v5/types"
	"log"
)
import "github.com/containers/podman/v5/pkg/bindings/images"

func BuildImage(ctx context.Context, app App) {
	var dockerfile string
	if app.Build.Dockerfile == "" {
		dockerfile = "Dockerfile"
	} else {
		dockerfile = app.Build.Dockerfile
	}
	containerFiles := []string{}
	containerFiles = append(containerFiles, dockerfile)
	buildOptions := images.BuildOptions{}
	if app.Build.Context == "" {
		buildOptions.ContextDirectory = "."
	} else {
		buildOptions.ContextDirectory = app.Build.Context
	}
	buildOptions.Target = app.Build.Target
	buildOptions.SkipUnusedStages = types.OptionalBoolTrue
	buildOptions.AdditionalTags = make([]string, 0)
	var tag string
	if app.Image != "" {
		tag = app.Image
	} else {
		log.Fatal("If you want to build image, you need to use 'custom' field in app section of frock.yaml file.")
	}
	buildOptions.NoCache = false
	if app.Tag != "" {
		tag = tag + ":" + app.Tag
	} else {
		tag = tag + ":local"
	}

	buildOptions.AdditionalTags = append(buildOptions.AdditionalTags, tag)
	buildOptions.Args = app.Build.Args
	buildReport, err := images.Build(ctx, containerFiles, buildOptions)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(buildReport)
}
