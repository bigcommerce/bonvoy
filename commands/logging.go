package commands

import (
	"bonvoy/envoy"
	"flag"
	"fmt"
	"github.com/manifoldco/promptui"
)

type SetLogLevelCommand struct {
	fs *flag.FlagSet
	name string
	level string
}

// ListenersCommand
func BuildSetLogLevelCommand() *SetLogLevelCommand {
	gc := &SetLogLevelCommand{
		fs: flag.NewFlagSet("log-level", flag.ContinueOnError),
	}
	gc.fs.StringVar(&gc.name, "service", "", "name of the service sidecar to set logging for")
	gc.fs.StringVar(&gc.level, "level", "", "Value to set logging to. Must be one of debug/info/warning/error")
	return gc
}

func (g *SetLogLevelCommand) Name() string {
	return g.fs.Name()
}

func (g *SetLogLevelCommand) Init(args []string) error {
	return g.fs.Parse(args)
}

func (g *SetLogLevelCommand) Run() error {
	var name = g.name
	if name == "" {
		name = g.fs.Arg(0)
	}
	var level = g.level
	if level == "" {
		level = g.fs.Arg(1)
		if level == "" {
			level = g.SelectLogLevel()
		}
	}

	e, err := envoy.NewFromServiceName(name)
	if err != nil {
		return err
	}
	e.Logging().SetLevel(level)
	return nil
}

func (g *SetLogLevelCommand) SelectLogLevel() string {
	prompt := promptui.Select{
		Label: "Please Select a Log Level",
		Items: []string{"debug", "info", "warning", "error"},
	}
	_, desired, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		panic(err)
	}
	return desired
}