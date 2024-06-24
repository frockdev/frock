package podman

import (
	"errors"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strings"
)

type Project struct {
	ProjectName string         `json:"projectName" yaml:"projectName"`
	MountPrefix string         `json:"mountPrefix,omitempty" yaml:"mountPrefix"`
	Apps        map[string]App `json:"apps,omitempty" yaml:"apps"`
	Boxes       map[string]Box `json:"boxes,omitempty" yaml:"boxes"`
	Commands    []Command      `json:"commands,omitempty" yaml:"commands"`
}

type App struct {
	Build             Build             `json:"build,omitempty" yaml:"build"`
	MountCode         string            `json:"mountCode,omitempty" jsonschema:"enum=enabled,enum=disabled" yaml:"mountCode"`
	MountContainerDir string            `json:"mountContainerDir,omitempty" yaml:"mountContainerDir"`
	Image             string            `json:"image,omitempty" yaml:"image"`
	Tag               string            `json:"tag,omitempty" yaml:"tag"`
	Command           []string          `json:"command,omitempty" yaml:"command"`
	Env               map[string]string `json:"env,omitempty" yaml:"env"`
	Labels            map[string]string `json:"labels,omitempty" yaml:"labels"`
	State             string            `json:"state,omitempty" jsonschema:"enum=disabled,enum=enabled" yaml:"state"`
	PrivilegedMode    string            `json:"privilegedMode,omitempty" jsonschema:"enum=enabled,enum=disabled" yaml:"privilegedMode"`
}

type Build struct {
	State      string            `json:"state,omitempty" jsonschema:"enum=disabled,enum=enabled" yaml:"state"`
	Context    string            `json:"context,omitempty" yaml:"context"`
	Dockerfile string            `json:"dockerfile,omitempty" yaml:"dockerfile"`
	Args       map[string]string `json:"args,omitempty" yaml:"args"`
	Target     string            `json:"target,omitempty" yaml:"target"`
}

type Box struct {
	State          string            `json:"state,omitempty" jsonschema:"enum=disabled,enum=enabled" yaml:"disabled"`
	Image          string            `json:"image,omitempty" yaml:"image"`
	Tag            string            `json:"tag,omitempty" yaml:"tag"`
	Env            map[string]string `json:"env,omitempty" yaml:"env"`
	Command        []string          `json:"command,omitempty" yaml:"command"`
	PrivilegedMode string            `json:"privilegedMode,omitempty" jsonschema:"enum=enabled,enum=disabled" yaml:"privilegedMode"`
	Labels         map[string]string `json:"labels,omitempty" yaml:"labels"`
}

//type DevelopmentPackage struct {
//	SSHLink             string `json:"sshLink" yaml:"sshLink"`
//	Branch              string `json:"branch" yaml:"branch"`
//	ComposerPackageName string `json:"composerPackageName" yaml:"composerPackageName"`
//}

type Command struct {
	WorkingDir  string   `json:"workingDir,omitempty" yaml:"workingDir"`
	AppName     string   `json:"appName,omitempty" yaml:"appName"`
	Signature   string   `json:"signature" yaml:"signature"`
	Type        string   `json:"type" jsonschema:"enum=local,enum=container" yaml:"type"`
	Description string   `json:"description,omitempty" yaml:"description"`
	Command     []string `json:"command" yaml:"command"`
	//Env         map[string]string `json:"env,omitempty" yaml:"env"`
	//FromFile    string            `json:"fromFile,omitempty" yaml:"fromFile"`
}

func checkFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	//return !os.IsNotExist(err)
	return !errors.Is(err, os.ErrNotExist)
}

// GetConfig reads the configuration file and returns the Project struct
// configuration file name is frock.yaml
// after that it will read frock.override.yaml and override the values
func GetConfig() (Project, error) {
	// first of all read frock.yaml into a structure Project
	filename := "frock.yaml"
	f, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
		return Project{}, err
	}
	defer f.Close()
	var project Project
	err = yaml.NewDecoder(f).Decode(&project)
	if err != nil {
		log.Fatal(err)
		return Project{}, err
	}
	if !checkFileExists("frock.override.yaml") {
		return project, nil
	}
	// read frock.override.yaml and override the values
	f, err = os.OpenFile("frock.override.yaml", os.O_RDONLY, 0644)
	if err != nil {
		return project, nil
	}
	defer f.Close()
	var override Project
	err = yaml.NewDecoder(f).Decode(&override)
	if err != nil {
		log.Fatal(err)
		return Project{}, err
	}
	//err = copier.CopyWithOption(&project, &override, copier.Option{IgnoreEmpty: false, DeepCopy: true})
	//if err != nil {
	//	return Project{}, err
	//}
	//err = copier.CopyWithOption(&project, &override, copier.Option{IgnoreEmpty: false, DeepCopy: true})
	project = mergeOverriden(project, override)
	project = normalizeToLowers(project)
	return project, nil
}

func normalizeToLowers(project Project) Project {
	project.ProjectName = strings.ToLower(project.ProjectName)
	for appName, appFields := range project.Apps {
		if strings.ToLower(appName) != appName {
			project.Apps[strings.ToLower(appName)] = appFields
			delete(project.Apps, appName)
		}
	}
	for boxName, boxFields := range project.Boxes {
		if strings.ToLower(boxName) != boxName {
			project.Boxes[strings.ToLower(boxName)] = boxFields
			delete(project.Boxes, boxName)
		}
	}
	return project
}

func mergeOverriden(main Project, overriden Project) Project {
	project := Project{}
	project.ProjectName = main.ProjectName
	project.Apps = make(map[string]App)
	project.Boxes = make(map[string]Box)
	//project.DevPackages = make(map[string]DevelopmentPackage)
	project.Commands = make([]Command, 0)
	for appName, appFields := range main.Apps {
		var finalAppFields App
		finalAppFields = appFields
		if overridesApp, ok := overriden.Apps[appName]; ok {
			if overridesApp.Build.State != "" {
				finalAppFields.Build.State = overridesApp.Build.State
			}
			if overridesApp.Build.Context != "" {
				finalAppFields.Build.Context = overridesApp.Build.Context
			}
			if overridesApp.Build.Dockerfile != "" {
				finalAppFields.Build.Dockerfile = overridesApp.Build.Dockerfile
			}
			if len(overridesApp.Build.Args) > 0 {
				for key, value := range overridesApp.Build.Args {
					finalAppFields.Build.Args[key] = value
				}
			}
			if overridesApp.Build.Target != "" {
				finalAppFields.Build.Target = overridesApp.Build.Target
			}
			if overridesApp.MountCode != "" {
				finalAppFields.MountCode = overridesApp.MountCode
			}
			if overridesApp.MountContainerDir != "" {
				finalAppFields.MountContainerDir = overridesApp.MountContainerDir
			}
			if overridesApp.Image != "" {
				finalAppFields.Image = overridesApp.Image
			}
			if overridesApp.Tag != "" {
				finalAppFields.Tag = overridesApp.Tag
			}
			if len(overridesApp.Command) > 0 {
				finalAppFields.Command = overridesApp.Command
			}
			if len(overridesApp.Env) > 0 {
				for key, value := range overridesApp.Env {
					finalAppFields.Env[key] = value
				}
			}
			if len(overridesApp.Labels) > 0 {
				for key, label := range overridesApp.Labels {
					finalAppFields.Labels[key] = label
				}
			}
		}
		project.Apps[appName] = finalAppFields
	}
	for boxName, boxFields := range main.Boxes {
		var finalBoxFields Box
		finalBoxFields = boxFields
		if overridesBox, ok := overriden.Boxes[boxName]; ok {
			if overridesBox.State != "" {
				finalBoxFields.State = overridesBox.State
			}
			if overridesBox.Image != "" {
				finalBoxFields.Image = overridesBox.Image
			}
			if overridesBox.Tag != "" {
				finalBoxFields.Tag = overridesBox.Tag
			}
			if len(overridesBox.Env) > 0 {
				for key, value := range overridesBox.Env {
					finalBoxFields.Env[key] = value
				}
			}
		}
		project.Boxes[boxName] = finalBoxFields
	}
	//for devPackageName, devPackageFields := range main.DevPackages {
	//	var finalDevPackageFields DevelopmentPackage
	//	finalDevPackageFields = devPackageFields
	//	if overridesDevPackage, ok := overriden.DevPackages[devPackageName]; ok {
	//		if overridesDevPackage.SSHLink != "" {
	//			finalDevPackageFields.SSHLink = overridesDevPackage.SSHLink
	//		}
	//		if overridesDevPackage.Branch != "" {
	//			finalDevPackageFields.Branch = overridesDevPackage.Branch
	//		}
	//		if overridesDevPackage.ComposerPackageName != "" {
	//			finalDevPackageFields.ComposerPackageName = overridesDevPackage.ComposerPackageName
	//		}
	//	}
	//	project.DevPackages[devPackageName] = finalDevPackageFields
	//}
	for _, commandFields := range main.Commands {
		var finalCommandFields Command
		finalCommandFields = commandFields
		if len(overriden.Commands) > 0 {
			for _, command := range overriden.Commands {
				if command.Signature == commandFields.Signature {
					if command.Type != "" {
						finalCommandFields.Type = command.Type
					}
					if command.Description != "" {
						finalCommandFields.Description = command.Description
					}
					if len(command.Command) > 0 {
						finalCommandFields.Command = command.Command
					}
				}
			}
		}
		project.Commands = append(project.Commands, finalCommandFields)
	}
	return project
}
