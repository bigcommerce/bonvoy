package commands

import (
	"bonvoy/envoy"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"strings"
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

/***********************************************************************************************************************
 * server info [service]
 **********************************************************************************************************************/
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
			o, err := controller.Run()
			if err != nil { return err }

			return r.Output(o)
		},
	}
}

type ServerInfoController struct {
	ServiceName string
}

func (s *ServerInfoController) Run() (ServerInfoResponse, error) {
	resp := ServerInfoResponse{
		ServiceName: s.ServiceName,
	}
	e, err := envoy.NewFromServiceName(s.ServiceName)
	if err != nil { return resp, err }

	resp.Envoy = &e

	response, err  := e.Server().Info()
	if err != nil { return resp, err }

	resp.Server = response
	return resp, nil
}

type ServerInfoResponse struct {
	ServiceName string `json:"service"`
	Envoy *envoy.Instance `json:"envoy"`
	Server envoy.ServerInfoJson `json:"server"`
}

func (s ServerInfoResponse) String() string {
	data := &s.Server
	o := ""

	o += Ok("----------------------")
	o += Ok("- Server Information -")
	o += Ok("----------------------")
	d := [][]string{
		{"Service", s.ServiceName},
		{"Envoy Pid", fmt.Sprintf("%d", s.Envoy.Pid)},
		{"Version", data.Version},
		{"Hot Restart Version", data.HotRestartVersion},
		{"State", data.State},
		{"Uptime", data.UptimeCurrentEpoch},
	}
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.AppendBulk(d)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
	o += tableString.String()

	o += Info("")
	o += Ok("--------------------")
	o += Ok("- Node Information -")
	o += Ok("--------------------")
	d = [][]string{
		{"Node ID", data.Node.ID},
		{"Node Cluster", data.Node.Cluster},
		{"User Agent", data.Node.UserAgentName},
		{"Envoy Version", data.Node.Metadata.EnvoyVersion},
		{"Namespace", data.Node.Metadata.Namespace},
	}
	tableString = &strings.Builder{}
	table = tablewriter.NewWriter(tableString)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.AppendBulk(d)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
	o += tableString.String()

	o += Info("")
	o += Ok("------------------------")
	o += Ok("- Command Line Options -")
	o += Ok("------------------------")
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

	tableString = &strings.Builder{}
	table = tablewriter.NewWriter(tableString)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.AppendBulk(d)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
	o += tableString.String()
	return o
}

/***********************************************************************************************************************
 * server memory [service]
 **********************************************************************************************************************/
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
			o, err := controller.Run()
			if err != nil { return err }

			return r.Output(o)
		},
	}
}

type ServerMemoryResponse struct {
	ServiceName string `json:"service"`
	Envoy *envoy.Instance `json:"envoy"`
	Memory envoy.ServerMemoryJson `json:"memory"`
}
func (s *ServerMemoryController) Run() (ServerMemoryResponse, error) {
	resp := ServerMemoryResponse{
		ServiceName: s.ServiceName,
	}
	e, err := envoy.NewFromServiceName(s.ServiceName)
	if err != nil {
		return resp, err
	}

	resp.Envoy = &e

	data, err := e.Server().Memory()
	if err != nil {
		return resp, err
	}
	resp.Memory = data
	return resp, nil
}

func (r ServerMemoryResponse) String() string {
	data := &r.Memory
	o := ""

	o += Ok("----------------------")
	o += Ok("- Server Memory Info -")
	o += Ok("----------------------")
	d := [][]string{
		{"Service", r.ServiceName},
		{"Envoy Pid", fmt.Sprintf("%d", r.Envoy.Pid)},
		{"Allocated", fmt.Sprintf("%d", data.Allocated)},
		{"Heap Size", fmt.Sprintf("%d", data.HeapSize)},
		{"Page Heap (Unmapped)", fmt.Sprintf("%d", data.PageHeapUnmapped)},
		{"Page Heap (Free)", fmt.Sprintf("%d", data.PageHeapFree)},
		{"Total Physical Bytes", fmt.Sprintf("%d", data.TotalPhysicalBytes)},
		{"Total Thread Cache", fmt.Sprintf("%d", data.TotalThreadCache)},
	}
	tableString := strings.Builder{}
	table := tablewriter.NewWriter(&tableString)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.AppendBulk(d)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
	o += tableString.String()
	return o
}

/***********************************************************************************************************************
 * server restart [service]
 **********************************************************************************************************************/
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
			Ok("Restarting " + controller.ServiceName + " Envoy...")
			o, err := controller.Run()
			if err != nil { return err }

			return r.Output(o)
		},
	}
}
type ServerRestartResponse struct {
	ServiceName string `json:"service"`
	Envoy *envoy.Instance `json:"envoy"`
	Ok bool `json:"ok"`
}

func (r ServerRestartResponse) String() string {
	o := ""
	o += Ok(r.ServiceName + " Envoy restarted.")
	return o
}

func (s *ServerRestartController) Run() (ServerRestartResponse, error) {
	resp := ServerRestartResponse{
		ServiceName: s.ServiceName,
		Ok: true,
	}
	e, err := envoy.NewFromServiceName(s.ServiceName)
	if err != nil { resp.Ok = false; return resp, err }
	resp.Envoy = &e

	rErr := e.Restart()
	if rErr != nil { resp.Ok = false; return resp, rErr }

	return resp, nil
}
