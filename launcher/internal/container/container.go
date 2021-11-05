package container

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/containerd/containerd/log"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/pterodactyl/wings/system"
)

type ImagePullProgressDetails struct {
	Current int `json:"current"`
	Total   int `json:"total"`
}

type ImagePullStatus struct {
	Id             string                    `json:"id"`
	Status         string                    `json:"status"`
	Progress       string                    `json:"progress"`
	ProgressDetail *ImagePullProgressDetails `json:"progressDetail"`
}

type ContainerEvent struct {
	Type ContainerEventType
	Data interface{}
}

type LaunchableContainer struct {
	onContainerEventFuncs []OnContainerEventFuc
	image                 string
	containerBody         *container.ContainerCreateCreatedBody
	containerName         string
	containerEnv          []string
	containerBinds        []string
	state                 ContainerState
	stream                *types.HijackedResponse
	// Tracks the environment state.
	st *system.AtomicString
}

type ContainerState int

const (
	ContainerOfflineState ContainerState = iota
	ContainerPullingState
	ContainerStartingState
	ContainerRunningState
	ContainerStoppingState
)

type ContainerEventType int

const (
	ImagePullStatusChanged ContainerEventType = iota
	ContainerStateChanged
	ConsoleOutput
)

type OnContainerEventFuc func(event ContainerEvent)

// A custom console writer that allows us to keep a function blocked until the
// given stream is properly closed. This does nothing special, only exists to
// make a noop io.Writer.
type noopWriter struct{}

var _ io.Writer = noopWriter{}

// Implement the required Write function to satisfy the io.Writer interface.
func (nw noopWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

var ctx = context.Background()
var dockerClient *client.Client = nil

func New() error {
	var err error
	dockerClient, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.L.Error("Unable to connect to docker daemon")
		dockerClient = nil
		return err
	}
	return nil
}

func IsConnectedToDocker() bool {
	return dockerClient != nil
}

func (c *LaunchableContainer) OnContainerEvent(callback OnContainerEventFuc) {
	c.onContainerEventFuncs = append(c.onContainerEventFuncs, callback)
}

func newContainer(image string, containerId string, containerEnv []string, containerBinds []string) (error, *LaunchableContainer) {

	c := LaunchableContainer{
		image:          image,
		containerName:  containerId,
		containerEnv:   containerEnv,
		containerBinds: containerBinds,
	}

	return nil, &c
}

func (c *LaunchableContainer) Launch() error {

	if err := c.Stop(); err != nil {
		return err
	}

	if err := c.Remove(); err != nil {
		return err
	}

	if err := c.Pull(); err != nil {
		return err
	}

	if err := c.Start(); err != nil {
		return err
	}

	if err := c.Attach(); err != nil {
		return err
	}

	return nil
}

func (c *LaunchableContainer) Stop() error {
	if err := dockerClient.ContainerStop(ctx, c.containerName, nil); err != nil && !client.IsErrNotFound(err) {
		log.L.Error("Could not stop old container: " + err.Error())
		return err
	}
	c.setState(ContainerOfflineState)
	return nil
}

func (c *LaunchableContainer) Remove() error {
	if err := dockerClient.ContainerRemove(ctx, c.containerName, types.ContainerRemoveOptions{Force: true}); err != nil && !client.IsErrNotFound(err) {
		log.L.Error("Could not remove old container: " + err.Error())
		return err
	}
	c.setState(ContainerOfflineState)
	return nil
}

func (c *LaunchableContainer) Pull() error {
	c.setState(ContainerPullingState)
	defer c.setState(ContainerOfflineState)

	reader, err := dockerClient.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		s := ImagePullStatus{}

		//fmt.Println(scanner.Text())
		if err := json.Unmarshal(scanner.Bytes(), &s); err == nil {
			c.handleContainerEvent(ContainerEvent{
				Type: ImagePullStatusChanged,
				Data: s,
			})
		}
	}
	return nil
}

func (c *LaunchableContainer) Start() error {
	c.setState(ContainerStartingState)
	cont, err := dockerClient.ContainerCreate(
		ctx,
		&container.Config{
			Image: image,
			Env:   c.containerEnv,
		},
		&container.HostConfig{
			Binds: c.containerBinds,
		}, nil, nil, c.containerName)

	if err != nil {
		c.setState(ContainerOfflineState)
		return err
	}

	c.containerBody = &cont

	err = dockerClient.ContainerStart(ctx, cont.ID, types.ContainerStartOptions{})
	if err != nil {
		c.setState(ContainerOfflineState)
		return err
	}

	c.setState(ContainerRunningState)
	return nil
}

func (c *LaunchableContainer) Attach() error {

	if err := c.followOutput(); err != nil {
		return err
	}

	opts := types.ContainerAttachOptions{
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Stream: true,
	}

	// Set the stream again with the container.
	if st, err := dockerClient.ContainerAttach(ctx, c.containerBody.ID, opts); err != nil {
		return err
	} else {
		c.stream = &st
	}

	go func() {
		defer c.stream.Close()
		defer func() {
			c.setState(ContainerOfflineState)
			c.Remove()
			c.stream = nil
		}()

		// Block the completion of this routine until the container is no longer running. This allows
		// the pollResources function to run until it needs to be stopped. Because the container
		// can be polled for resource usage, even when stopped, we need to have this logic present
		// in order to cancel the context and therefore stop the routine that is spawned.
		//
		// For now, DO NOT use client#ContainerWait from the Docker package. There is a nasty
		// bug causing containers to hang on deletion and cause servers to lock up on the system.
		//
		// This weird code isn't intuitive, but it keeps the function from ending until the container
		// is stopped and therefore the stream reader ends up closed.
		// @see https://github.com/moby/moby/issues/41827
		w := new(noopWriter)
		if _, err := io.Copy(w, c.stream.Reader); err != nil {
			c.log("error", "could not copy from environment stream to noop writer: "+err.Error())
		}
	}()

	return nil
}

func (c *LaunchableContainer) followOutput() error {

	opts := types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Follow:     true,
		Since:      time.Now().Format(time.RFC3339),
	}

	reader, err := dockerClient.ContainerLogs(context.Background(), c.containerBody.ID, opts)
	if err != nil {
		return err
	}

	go c.scanOutput(reader)

	return nil
}

func (c *LaunchableContainer) scanOutput(reader io.ReadCloser) {
	defer reader.Close()

	if err := system.ScanReader(reader, func(line string) {
		c.handleContainerEvent(ContainerEvent{
			Type: ConsoleOutput,
			Data: line,
		})
	}); err != nil && err != io.EOF {
		c.log("error", "error processing scanner line in console output: "+err.Error())
		return
	}

	if c.state == ContainerStoppingState || c.state == ContainerOfflineState {
		return
	}

	_ = reader.Close()

	go c.followOutput()
}

func (c *LaunchableContainer) handleContainerEvent(e ContainerEvent) {
	for _, callbackFunc := range c.onContainerEventFuncs {
		callbackFunc(e)
	}
}

func (c *LaunchableContainer) setState(s ContainerState) {
	if s < ContainerOfflineState || s > ContainerStoppingState {
		panic(errors.New(fmt.Sprintf("invalid container state received: %d", s)))
	}

	// Emit the event to any listeners that are currently registered.
	if c.state != s {
		// If the state changed make sure we update the internal tracking to note that.
		c.state = s
		c.handleContainerEvent(ContainerEvent{
			Type: ContainerStateChanged,
			Data: s,
		})
	}
}

func (c *LaunchableContainer) log(l string, m string) {
	log.L.Info(m)
}
