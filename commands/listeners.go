package commands

import (
	"bonvoy/docker"
	"bonvoy/envoy"
	envoyApi "bonvoy/envoy/api"
	"bonvoy/nsenter"
	"flag"
	"fmt"
)

type ListenersCommand struct {
	fs *flag.FlagSet
	name string
}

// ListenersCommand
func BuildListeners() *ListenersCommand {
	gc := &ListenersCommand{
		fs: flag.NewFlagSet("listeners", flag.ContinueOnError),
	}
	gc.fs.Arg(0)
	gc.fs.StringVar(&gc.name, "service", "", "name of the service whose sidecar to enter")
	return gc
}

func (g *ListenersCommand) Name() string {
	return g.fs.Name()
}

func (g *ListenersCommand) Init(args []string) error {
	return g.fs.Parse(args)
}

func (g *ListenersCommand) Run() error {
	var name = g.name
	if name == "" {
		name = g.fs.Arg(0)
	}

	cli := docker.NewClient()
	pid := envoy.GetPid(cli, name)
	config := nsenter.BuildConfig(pid)
	listeners := envoyApi.GetListeners(config)

	fmt.Println("LISTENERS:")
	fmt.Println("----------------------------------------------------------------------")
	for _, listener := range listeners {
		fmt.Printf("%s:\t\t%s\n", listener.Name, listener.TargetAddress)
	}
	return nil
}
