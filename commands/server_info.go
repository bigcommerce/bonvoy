package commands

import (
	"bonvoy/envoy"
	"flag"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
)

type ServerInfoCommand struct {
	fs *flag.FlagSet
	name string
}

// ListenersCommand
func BuildServerInfoCommand() *ServerInfoCommand {
	gc := &ServerInfoCommand{
		fs: flag.NewFlagSet("server-info", flag.ContinueOnError),
	}
	gc.fs.StringVar(&gc.name, "service", "", "name of the service sidecar to see server information for")
	return gc
}

func (g *ServerInfoCommand) Name() string {
	return g.fs.Name()
}

func (g *ServerInfoCommand) Init(args []string) error {
	return g.fs.Parse(args)
}

func (g *ServerInfoCommand) Run() error {
	var name = g.name
	if name == "" {
		name = g.fs.Arg(0)
	}
	e, err := envoy.NewFromServiceName(name)
	if err != nil {
		return err
	}
	response := e.Server().Info()
	g.DisplayOutput(response)
	return nil
}

func (g *ServerInfoCommand) DisplayOutput(data envoy.ServerInfoJson) {
	fmt.Println("----------------------")
	fmt.Println("- Server Information -")
	fmt.Println("----------------------")
	d := [][]string{
		{"Version", data.Version},
		{"Hot Restart Version", data.HotRestartVersion},
		{"State", data.State},
		{"Uptime", data.UptimeCurrentEpoch},
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.AppendBulk(d)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()

	fmt.Println("")
	fmt.Println("--------------------")
	fmt.Println("- Node Information -")
	fmt.Println("--------------------")
	d = [][]string{
		{"Node ID", data.Node.ID},
		{"Node Cluster", data.Node.Cluster},
		{"User Agent", data.Node.UserAgentName},
		{"Envoy Version", data.Node.Metadata.EnvoyVersion},
		{"Namespace", data.Node.Metadata.Namespace},
	}
	table = tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.AppendBulk(d)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()

	fmt.Println("")
	fmt.Println("------------------------")
	fmt.Println("- Command Line Options -")
	fmt.Println("------------------------")
	d = [][]string{
		{"Concurrency", fmt.Sprintf("%d", data.ServerCommandLineOptions.Concurrency)},
		{"Mode", data.ServerCommandLineOptions.Mode},
		{"Log Level", data.ServerCommandLineOptions.LogLevel},
		{"Component Log Level", data.ServerCommandLineOptions.ComponentLogLevel},
		{"Log Format", data.ServerCommandLineOptions.LogFormat},
		{"Drain Strategy", data.ServerCommandLineOptions.DrainStrategy},
		{"Drain Time", data.ServerCommandLineOptions.DrainTime},
		{"Config Path", data.ServerCommandLineOptions.ConfigPath},
		{"Parent Shutdown Time", data.ServerCommandLineOptions.ParentShutdownTime},
	}
	table = tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.AppendBulk(d)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
}