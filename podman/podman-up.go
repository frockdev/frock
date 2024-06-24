package podman

import (
	"log"
)

func RunAppsBoxesInPod() {
	config, err := GetConfig()
	if err != nil {
		log.Fatal(err)
		return
	}
	podmanConnection := Connect()
	for _, app := range config.Apps {
		if app.Build.State == "enabled" {
			BuildImage(podmanConnection, app)
		}
	}
	_ = CreatePod(podmanConnection, config.ProjectName)

	for boxName, box := range config.Boxes {
		if box.State == "enabled" || box.State == "" {
			RunBox(podmanConnection, box, boxName, config.ProjectName)
		}
	}

	for appName, app := range config.Apps {
		if app.State == "enabled" || app.State == "" {
			RunContainer(podmanConnection, app, appName, config.ProjectName)
		}
	}
}
