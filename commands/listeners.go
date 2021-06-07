package commands

import (
	"bonvoy/docker"
	"bonvoy/envoy"
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
	//gc.fs.StringVar(&gc.name, "service", "", "name of the service whose sidecar to enter")
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

	var cli = docker.NewClient()
	var pid = envoy.GetPid(cli, name)
	config := nsenter.BuildConfig(pid)
	stdout, stderr, err := config.Execute("curl", "0.0.0.0:19001/listeners")
	if err != nil {
		fmt.Println(stderr)
		panic(err)
	}
	fmt.Println(stdout)

	return nil
}
