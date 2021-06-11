package docker

import (
	"context"
	"fmt"
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
	filter := filters.NewArgs()
	filter.Add("name", "connect-proxy-" + name)

	containers, err := c.cli.ContainerList(context.Background(), types.ContainerListOptions{
		All: false,
		Filters: filter,
	})
	if err != nil { return 0, err }

	var containerNames []string
	for _, container := range containers {
		containerNames = append(containerNames, container.Names...)
	}

	desiredName := ""
	if len(containerNames) > 1 {
		desiredName, err = c.SelectDesiredContainer(containerNames)
		if err != nil { return 0, err }

	} else if len(containerNames) == 1 {
		desiredName = containerNames[0]
	} else {
		return 0, fmt.Errorf("No sidecar found for name: " + name)
	}

	var desiredId = ""
	for _, container := range containers {
		if container.Names[0] == desiredName {
			desiredId = container.ID
		}
	}
	container, err := c.cli.ContainerInspect(context.Background(), desiredId)
	if err != nil {
		return 0, err
	}
	return container.State.Pid, nil
}

func (c *Client) SelectDesiredContainer(names []string) (string, error) {
	prompt := promptui.Select{
		Label: "Please Select Sidecar Container",
		Items: names,
	}
	_, desiredName, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return desiredName, nil
}