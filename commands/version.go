package commands

import (
	"bonvoy/version"
	"flag"
	"fmt"
)

type VersionCommand struct {
	fs *flag.FlagSet
	name string
}

func BuildVersion() *VersionCommand {
	gc := &VersionCommand{
		fs: flag.NewFlagSet("listeners", flag.ContinueOnError),
	}
	gc.fs.Arg(0)
	//gc.fs.StringVar(&gc.name, "service", "", "name of the service whose sidecar to enter")
	return gc
}

func (g *VersionCommand) Name() string {
	return g.fs.Name()
}

func (g *VersionCommand) Init(args []string) error {
	return g.fs.Parse(args)
}

func (g *VersionCommand) Run() error {
	fmt.Println(version.Version)
	return nil
}
