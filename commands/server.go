package commands

import (
	"bonvoy/envoy"
	"fmt"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

type Server struct {
	Command *cobra.Command
}

func (r *Registry) Server() *Server {
	cmd := &cobra.Command{
		Use: "server",
		Short: "Envoy server commands",
	}
	cmd.AddCommand(r.BuildServerInfoCommand())
	cmd.AddCommand(r.BuildServerMemoryCommand())
	cmd.AddCommand(r.BuildServerRestartCommand())
	return &Server{
		Command: cmd,
	}
}

// server info

func (r *Registry) BuildServerInfoCommand() *cobra.Command {
	return &cobra.Command{
		Use: "info",
		Short: "Display envoy server information",
		Long:  `Display server information about the envoy sidecar`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			controller := ServerInfoController{
				ServiceName: args[0],
			}
			return controller.Run()
		},
	}
}

type ServerInfoController struct {
	ServiceName string
}

func (s *ServerInfoController) Run() error {
	e, err := envoy.NewFromServiceName(s.ServiceName)
	if err != nil { return err }

	response, err  := e.Server().Info()
	if err != nil { return err }

	s.DisplayOutput(response)
	return nil
}

func (s *ServerInfoController) DisplayOutput(data envoy.ServerInfoJson) {
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

// server memory

type ServerMemoryController struct {
	ServiceName string
}
func (r *Registry) BuildServerMemoryCommand() *cobra.Command {
	return &cobra.Command{
		Use: "memory",
		Short: "Display envoy memory information",
		Long:  `Display memory usage information about the envoy sidecar`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			controller := ServerMemoryController{
				ServiceName: args[0],
			}
			return controller.Run()
		},
	}
}

func (s *ServerMemoryController) Run() error {
	e, err := envoy.NewFromServiceName(s.ServiceName)
	if err != nil { return err }

	data, err  := e.Server().Memory()
	if err != nil { return err }

	fmt.Println("----------------------")
	fmt.Println("- Server Memory Info -")
	fmt.Println("----------------------")
	d := [][]string{
		{"Allocated", data.Allocated},
		{"Heap Size", data.HeapSize},
		{"Page Heap (Unmapped)", data.PageHeapUnmapped},
		{"Page Heap (Free)", data.PageHeapFree},
		{"Total Physical Bytes", data.TotalPhysicalBytes},
		{"Total Thread Cache", data.TotalThreadCache},
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.AppendBulk(d)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
	return nil
}

// server restart

type ServerRestartController struct {
	ServiceName string
}
func (r *Registry) BuildServerRestartCommand() *cobra.Command {
	return &cobra.Command{
		Use: "restart",
		Short: "Restart an Envoy sidecar",
		Long:  `Restarts an Envoy sidecar process`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			controller := ServerRestartController{
				ServiceName: args[0],
			}
			return controller.Run()
		},
	}
}

func (s *ServerRestartController) Run() error {
	e, err := envoy.NewFromServiceName(s.ServiceName)
	if err != nil { return err }

	err = e.Restart()
	if err != nil { return err }

	color.Green(s.ServiceName + " Envoy restarted.")
	// TODO: Would be nice to monitor and let the user know when the service comes back alive.
	return nil
}
