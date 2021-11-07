package container

import (
	"fmt"
	"runtime"

	"github.com/keinenclue/sasm-docker/launcher/internal/config"
)

const image string = "ghcr.io/keinenclue/sasm-docker"

// NewSasmContainer creates a new launchable sasm container
func NewSasmContainer() (*LaunchableContainer, error) {
	containerEnv := []string{}
	containerBinds := []string{}

	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		containerEnv = append(containerEnv, "DISPLAY=host.docker.internal:0")
	} else if runtime.GOOS == "linux" {
		xauthority := "/tmp/.docker.xauth"
		containerEnv = append(containerEnv, "DISPLAY=:0")
		containerEnv = append(containerEnv, fmt.Sprintf("XAUTHORITY=%s", xauthority))
		containerBinds = append(containerBinds, fmt.Sprintf(`%[1]v:%[1]v`, xauthority))
		containerBinds = append(containerBinds, "/tmp/.X11-unix:/tmp/.X11-unix")
	}
	containerBinds = append(containerBinds, config.Get("dataPath").(string)+"/.config:/root/.config")
	containerBinds = append(containerBinds, config.Get("dataPath").(string)+"/Projects:/usr/share/sasm/Projects")

	return newContainer(image, "sasm_docker_container", containerEnv, containerBinds)
}
