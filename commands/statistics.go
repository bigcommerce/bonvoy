package commands

import (
	"bonvoy/consul"
	"bonvoy/envoy"
	"fmt"
	"github.com/spf13/cobra"
)

type Statistics struct {
	Command *cobra.Command
}

func (r *Registry) Statistics() *Statistics {
	cmd := &cobra.Command{
		Use: "stats",
		Short: "Envoy stats commands",
	}
	cmd.AddCommand(r.BuildStatisticsDumpCommand())
	return &Statistics{
		Command: cmd,
	}
}

func (r *Registry) BuildStatisticsDumpCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "dump",
		Short: "Output all statistics",
		Long:  `Output all statistics for the Envoy process`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			controller := StatisticsDumpController{
				ServiceName: args[0],
			}
			return controller.Run()
		},
	}
	cmd.Flags().BoolP("restart", "r", false, "If passed, will restart all sidecars that have expired certificates")
	return cmd
}

type StatisticsDumpController struct {
	ServiceName string
	Consul consul.Client
}

func (c *StatisticsDumpController) Run() error {
	e, err := envoy.NewFromServiceName(c.ServiceName)
	if err != nil { return err }

	result, cErr := e.Statistics().Dump()
	if cErr != nil { return cErr }

	fmt.Println(result)
	return nil
}