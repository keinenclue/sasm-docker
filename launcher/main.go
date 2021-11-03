package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/user"
	"runtime"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

const image string = "ghcr.io/keinenclue/sasm-docker"

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println("Unable to create docker client")
		panic(err)
	}

	fmt.Println("Removing old container...")
	cli.ContainerRemove(ctx, "sasm_docker_container", types.ContainerRemoveOptions{})
	fmt.Println("Pulling image...")
	reader, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

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

	u, _ := user.Current()
	containerBinds = append(containerBinds, fmt.Sprintf("%s/sasm-data:/root", u.HomeDir))

	// TODO: Add storage bind!

	fmt.Printf("Env:\n%s\n\n", containerEnv)
	fmt.Printf("Volumes:\n%s\n\n", containerBinds)

	cont, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: image,
			Env:   containerEnv,
		},
		&container.HostConfig{
			Binds: containerBinds,
		}, nil, nil, "sasm_docker_container")

	if err != nil {
		panic(err)
	}

	cli.ContainerStart(ctx, cont.ID, types.ContainerStartOptions{})
	fmt.Printf("Container %s has been started.\n", cont.ID)

	statusCh, errCh := cli.ContainerWait(ctx, cont.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, cont.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
}
