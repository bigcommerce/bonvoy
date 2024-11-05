package docker

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	dockerClient "github.com/docker/docker/client"
	"github.com/manifoldco/promptui"
)

type Client struct {
	cli *dockerClient.Client
}

func NewClient() Client {
	cli, err := dockerClient.NewClientWithOpts(dockerClient.FromEnv)
	if err != nil {
		panic(err)
	}
	return Client{
		cli: cli,
	}
}

func (c *Client) GetEnvoyPid(name string) (int, error) {
	container, err := c.GetSidecarContainer(name)
	if err != nil {
		return 0, err
	}

	return container.State.Pid, nil
}

func (c *Client) GetSidecarContainer(serviceName string) (types.ContainerJSON, error) {
	filter := filters.NewArgs()
	filter.Add("name", "connect-proxy-"+serviceName)

	containers, err := c.cli.ContainerList(context.Background(), types.ContainerListOptions{
		All:     false,
		Filters: filter,
	})
	if err != nil {
		return types.ContainerJSON{}, err
	}

	var containerNames []string
	for _, container := range containers {
		containerNames = append(containerNames, container.Names...)
	}

	desiredName := ""
	if len(containerNames) > 1 {
		desiredName, err = c.SelectDesiredContainer(containerNames)
		if err != nil {
			return types.ContainerJSON{}, err
		}

	} else if len(containerNames) == 1 {
		desiredName = containerNames[0]
	} else {
		return types.ContainerJSON{}, fmt.Errorf("No sidecar found for name: " + serviceName)
	}

	var desiredId = ""
	for _, container := range containers {
		if container.Names[0] == desiredName {
			desiredId = container.ID
		}
	}
	container, err := c.cli.ContainerInspect(context.Background(), desiredId)
	if err != nil {
		return types.ContainerJSON{}, err
	}
	return container, nil
}

var defaultContainerPageSize = 500

func (c *Client) SelectDesiredContainer(names []string) (string, error) {
	size, err := strconv.Atoi(os.Getenv("BONVOY_MAX_CONTAINER_PAGE_SIZE"))
	if err != nil {
		size = defaultContainerPageSize
	}
	prompt := promptui.Select{
		Label: "Please Select Sidecar Container",
		Items: names,
		Size:  size,
	}
	_, desiredName, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return desiredName, nil
}
