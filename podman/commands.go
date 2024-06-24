package podman

import (
	"bufio"
	"context"
	"github.com/containers/podman/v5/pkg/api/handlers"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"io"
	"log"
	"os"
	"os/exec"
)

func RunCommand(commandName string) {
	config, err := GetConfig()
	if err != nil {
		log.Fatal(err)
		return
	}
	podmanConnection := Connect()
	found := false
	for _, command := range config.Commands {
		if command.Signature == commandName {
			found = true
			if command.Type == "container" {
				var containerName string
				var containerWorkingDir string
				if command.AppName == "" {
					for appName, app := range config.Apps {
						containerName = appName
						if app.MountContainerDir != "" {
							containerWorkingDir = app.MountContainerDir
						} else {
							containerWorkingDir = "/app"
						}
						break
					}
				}
				var workingDir string
				if command.WorkingDir == "" {
					workingDir = containerWorkingDir
				} else {
					workingDir = command.WorkingDir
				}
				runInContainer(podmanConnection, command, containerName, workingDir, config.ProjectName)
				return

			} else if command.Type == "local" {
				runLocal(command)
			}
		}
	}
	if !found {
		log.Fatal("Command not found")
	}
}

func runLocal(command Command) {
	cmd := exec.Command(command.Command[0], command.Command[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run() // add error checking
	if err != nil {
		log.Fatal(err)
	}
}

//func copy(src, dst string) (int64, error) {
//	sourceFileStat, err := os.Stat(src)
//	if err != nil {
//		return 0, err
//	}
//
//	if !sourceFileStat.Mode().IsRegular() {
//		return 0, fmt.Errorf("%s is not a regular file", src)
//	}
//
//	source, err := os.Open(src)
//	if err != nil {
//		return 0, err
//	}
//	defer source.Close()
//
//	destination, err := os.Create(dst)
//	if err != nil {
//		return 0, err
//	}
//	defer destination.Close()
//	nBytes, err := io.Copy(destination, source)
//	return nBytes, err
//}

func runInContainer(ctx context.Context, command Command, containerName string, workingDir string, projectName string) {
	execOptions := handlers.ExecCreateConfig{}
	//wd, _ := os.Getwd()
	//if command.FromFile != "" {
	//	_, err := os.OpenFile(command.FromFile, os.O_RDONLY, 0644)
	//	if os.IsNotExist(err) {
	//		log.Fatal("File not found")
	//		return
	//	}
	//	if err == nil {
	//		_, err = copy(command.FromFile, wd+"/runnable")
	//	} else {
	//		log.Fatal(err)
	//		return
	//	}
	//	execOptions.Cmd = []string{}
	//	execOptions.Cmd = append(execOptions.Cmd, "bash")
	//	execOptions.Cmd = append(execOptions.Cmd, "-c")
	//	execOptions.Cmd = append(execOptions.Cmd, "chmod +x "+filepath.Base(command.FromFile)+" && "+"./runnable")
	//	log.Println(execOptions.Cmd)
	//} else {
	execOptions.Cmd = command.Command
	//}
	execOptions.WorkingDir = workingDir
	execOptions.AttachStdin = true
	execOptions.AttachStdout = true
	execOptions.AttachStderr = true
	execOptions.Tty = true
	//execOptions.Privileged = true

	//var envOptions []string
	//for key, value := range command.Env {
	//	execOptions.Env = append(envOptions, key+"="+value)
	//}
	//execOptions.Env = envOptions

	session, err := containers.ExecCreate(ctx, projectName+"-"+containerName, &execOptions)
	if err != nil {
		log.Fatal(err)
	}
	attachOptions := containers.ExecStartAndAttachOptions{}
	attachTrue := true
	attachOptions.AttachOutput = &attachTrue
	attachOptions.AttachError = &attachTrue
	attachOptions.AttachInput = &attachTrue

	attachOptions.InputStream = bufio.NewReader(os.Stdin)

	var writer io.Writer = os.Stdout
	attachOptions.OutputStream = &writer

	var errorWriter io.Writer = os.Stderr
	attachOptions.ErrorStream = &errorWriter

	err = containers.ExecStartAndAttach(ctx, session, &attachOptions)
	if err != nil {
		log.Fatal(err)
	}
	//if command.FromFile != "" {
	//	//delete file
	//	err = os.Remove(wd + "/runnable")
	//}
}
