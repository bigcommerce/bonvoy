package commands

import (
	"bonvoy/consul"
	"bonvoy/envoy"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"sort"
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
			o, err := controller.Run()
			if err != nil { return err }

			return r.Output(o)
		},
	}
	cmd.Flags().BoolP("restart", "r", false, "If passed, will restart all sidecars that have expired certificates")
	return cmd
}

type StatisticsDumpController struct {
	ServiceName string
	Consul consul.Client
}

func (c *StatisticsDumpController) Run() (StatisticsDumpResponse, error) {
	resp := StatisticsDumpResponse{
		ServiceName: c.ServiceName,
		Statistics: make(map[string]string),
	}
	e, err := envoy.NewFromServiceName(c.ServiceName)
	if err != nil { return resp, err }

	resp.Envoy = &e

	stats, cErr := e.Statistics().List()
	if cErr != nil { return resp, cErr }

	for _, stat := range stats {
		resp.Statistics[stat.Name] = stat.Value
	}
	return resp, nil
}

type StatisticsDumpResponse struct {
	ServiceName string `json:"service"`
	Envoy *envoy.Instance `json:"envoy"`
	Statistics map[string]string `json:"statistics"`
}

func (r StatisticsDumpResponse) String() string {
	o := ""
	o += Ok("Statistics for "+r.ServiceName+" Envoy (PID "+cast.ToString(r.Envoy.Pid)+")")
	o += Ok("----------------------------------------------------------------------")

	keys := make([]string, 0, len(r.Statistics))
	for k := range r.Statistics {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		o += Info(k + ": "+r.Statistics[k])
	}
	return o
}