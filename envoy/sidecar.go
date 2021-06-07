package envoy

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/manifoldco/promptui"
	"strings"
)

func GetPid(cli *client.Client, name string) int {
	filter := filters.NewArgs()
	filter.Add("name", "connect-proxy-" + name)

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
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

	prompt := promptui.Select{
		Label: "Please Select Sidecar Container",
		Items: containerNames,
	}
	_, desiredName, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		panic(err)
	}

	fmt.Printf("Entering %s\n\n", strings.TrimLeft(desiredName, "/"))

	var desiredId = ""
	for _, container := range containers {
		if container.Names[0] == desiredName {
			desiredId = container.ID
		}
	}
	container, err := cli.ContainerInspect(context.Background(), desiredId)
	if err != nil {
		fmt.Printf("Failed to inspect container %v\n", err)
		panic(err)
	}
	return container.State.Pid
}
