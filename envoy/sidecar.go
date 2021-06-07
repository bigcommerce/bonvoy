package envoy

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
	"strings"
)

func GetHost() string {
	return viper.GetString("ENVOY_HOST")
}

func GetPid(docker *client.Client, name string) int {
	filter := filters.NewArgs()
	filter.Add("name", "connect-proxy-" + name)

	containers, err := docker.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
		Filters: filter,
	})
	if err != nil {
		panic(err)
	}
	var containerNames []string
	for _, container := range containers {
		containerNames = append(containerNames, container.Names...)
	}

	desiredName := ""
	if len(containerNames) > 1 {
		desiredName = SelectContainer(containerNames)
	} else {
		desiredName = containerNames[0]
	}

	fmt.Printf("Entering %s\n\n", strings.TrimLeft(desiredName, "/"))

	var desiredId = ""
	for _, container := range containers {
		if container.Names[0] == desiredName {
			desiredId = container.ID
		}
	}
	container, err := docker.ContainerInspect(context.Background(), desiredId)
	if err != nil {
		fmt.Printf("Failed to inspect container %v\n", err)
		panic(err)
	}
	return container.State.Pid
}

func SelectContainer(names []string) string {
	prompt := promptui.Select{
		Label: "Please Select Sidecar Container",
		Items: names,
	}
	_, desiredName, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		panic(err)
	}
	return desiredName
}