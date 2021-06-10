package commands

import (
	"bonvoy/envoy"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"os"
)

type Listeners struct {
	Command *cobra.Command
}

func (r *Registry) Listeners() *Listeners {
	return &Listeners{
		Command: &cobra.Command{
			Use: "listeners",
			Short: "Show Envoy listeners",
			Long:  `Display all registered Envoy sidecar listeners`,
			Args: cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				controller := ListenersController{
					ServiceName: args[0],
				}
				return controller.Run()
			},
		},
	}
}

type ListenersController struct {
	ServiceName string
}

func (s *ListenersController) Run() error {
	e, err := envoy.NewFromServiceName(s.ServiceName)
	if err != nil { return err }

	listeners, err := e.Listeners().Get()
	if err != nil { return err }

	color.Green("Listeners for "+s.ServiceName+" Envoy (PID "+cast.ToString(e.Pid)+")")
	color.Green("----------------------------------------------------------------------")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Name",
		"Address",
	})
	table.SetBorder(false)
	table.SetTablePadding("\t")
	var d [][]string
	for _, i := range listeners {
		d = append(d, []string{
			i.Name,
			i.TargetAddress,
		})
	}
	table.AppendBulk(d)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
	return nil
}