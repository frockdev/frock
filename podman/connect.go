package podman

import (
	"context"
	"log"
	"os"
	"os/user"
	"runtime"
	"strings"
	"time"
)

import (
	"github.com/containers/podman/v5/pkg/bindings"
)

func Connect() context.Context {
	//if operating system is macos, then connection is ssh://core@127.0.0.1:53523/runInContainer/user/501/podman/podman.sock or ssh://root@127.0.0.1:53523/runInContainer/podman/podman.sock
	//if operating system is linux, then connection is unix:///runInContainer/podman/podman.sock
	//if operating system is windows, then connection is npipe:////./pipe/podman
	var connectionUri string
	if runtime.GOOS == "windows" {
		connectionUri = "npipe:////./pipe/podman"
	} else {
		if runtime.GOOS == "darwin" {
			connectionUri = getDarwinUnixSocket()
		} else {
			connectionUri = "unix:///runInContainer/podman/podman.sock"
		}
	}
	conn, err := bindings.NewConnection(context.Background(), connectionUri)
	if err != nil {
		log.Fatal(err)
		//os.Exit(1)
	}
	return conn

}

func getDarwinUnixSocket() string {

	//lets read file ~/podman-socket
	//if file exists, then return the content
	//if file does not exist - write to terminal that file does not exist
	// and user should get unix path from podman's settings and put into this file
	usr, _ := user.Current()
	dir := usr.HomeDir

	file, err := os.OpenFile(dir+"/podman-socket", os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal("File does not exist. Please create file ~/podman-socket and put unix path from podman settings into it. You can runInContainer 'podman machine inspect' and get unix path from .ConnectionInfo.PodmanSocket.Path")
	}
	defer file.Close()
	buf := make([]byte, 1024)
	n, err := file.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Trim(string(buf[:n]), "\n")
}

type MachineInfo struct {
	ConfigDir struct {
		Path string `json:"Path"`
	} `json:"ConfigDir"`
	ConnectionInfo struct {
		PodmanSocket struct {
			Path string `json:"Path"`
		} `json:"PodmanSocket"`
		PodmanPipe interface{} `json:"PodmanPipe"`
	} `json:"ConnectionInfo"`
	Created   time.Time `json:"Created"`
	LastUp    time.Time `json:"LastUp"`
	Name      string    `json:"Name"`
	Resources struct {
		CPUs     int           `json:"CPUs"`
		DiskSize int           `json:"DiskSize"`
		Memory   int           `json:"Memory"`
		USBs     []interface{} `json:"USBs"`
	} `json:"Resources"`
	SSHConfig struct {
		IdentityPath   string `json:"IdentityPath"`
		Port           int    `json:"Port"`
		RemoteUsername string `json:"RemoteUsername"`
	} `json:"SSHConfig"`
	State              string `json:"State"`
	UserModeNetworking bool   `json:"UserModeNetworking"`
	Rootful            bool   `json:"Rootful"`
}
