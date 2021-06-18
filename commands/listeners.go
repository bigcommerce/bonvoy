package commands

import (
	"bonvoy/envoy"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"strings"
)

type Listeners struct {
	Command *cobra.Command
}

func (r *Registry) Listeners() *Listeners {
	cmd := &cobra.Command{
		Use:   "listeners",
		Short: "listeners-related commands",
	}

	cmd.AddCommand(r.BuildListListenersCommand())
	return &Listeners{
		Command: cmd,
	}
}

func (r *Registry) BuildListListenersCommand() *cobra.Command {
	return &cobra.Command{
		Use: "list",
		Short: "Show Envoy listeners",
		Long:  `Display all registered Envoy sidecar listeners`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			controller := ListListenersController{
				ServiceName: args[0],
			}
			o, err := controller.Run()
			if err != nil { return err }

			return r.Output(o)
		},
	}
}

type ListListenersController struct {
	ServiceName string
}

type ListListenersResponse struct {
	ServiceName string `json:"service"`
	Envoy *envoy.Instance `json:"envoy"`
	Listeners []envoy.Listener `json:"listeners"`
}

func (s *ListListenersController) Run() (ListListenersResponse, error) {
	resp := ListListenersResponse{
		ServiceName: s.ServiceName,
	}
	e, err := envoy.NewFromServiceName(s.ServiceName)
	if err != nil {
		return resp, err
	}

	resp.Envoy = &e

	listeners, err := e.Listeners().Get()
	if err != nil {
		return resp, err
	}

	resp.Listeners = listeners
	return resp, nil
}

func (r ListListenersResponse) String() string {
	o := ""
	o += Ok("Listeners for "+r.ServiceName+" Envoy (PID "+cast.ToString(r.Envoy.Pid)+")")
	o += Ok("----------------------------------------------------------------------")

	tableString := strings.Builder{}
	table := tablewriter.NewWriter(&tableString)
	table.SetHeader([]string{
		"Name",
		"Address",
	})
	table.SetBorder(false)
	table.SetTablePadding("\t")
	var d [][]string
	for _, i := range r.Listeners {
		d = append(d, []string{
			i.Name,
			i.TargetAddress,
		})
	}
	table.AppendBulk(d)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
	o += tableString.String()
	return o
}