package commands

import (
	"bonvoy/consul"
	"bonvoy/envoy"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"strings"
)

type Clusters struct {
	Command *cobra.Command
}

func (r *Registry) Clusters() *Clusters {
	cmd := &cobra.Command{
		Use: "clusters",
		Short: "Clusters-related commands",
	}
	cmd.AddCommand(r.BuildListClustersCommand())
	return &Clusters{
		Command: cmd,
	}
}

func (r *Registry) BuildListClustersCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "list",
		Short: "List clusters statistics",
		Long:  `Display all clusters statistics for a given service`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cluster, err := cmd.Flags().GetString("cluster")
			if err != nil { return err }
			if cluster == "" && len(args) >= 2 { // allow `clusters list [service] [cluster]` syntax
				cluster = args[1]
			}

			controller := ListClustersController{
				ServiceName: args[0],
				Cluster: cluster,
				Consul: consul.NewClient(),
			}
			o, cErr := controller.Run()
			if cErr != nil { return cErr }

			return r.Output(o)
		},
	}
	cmd.Flags().String("cluster", "", "Filter to a specific cluster")
	return cmd
}

type ListClustersController struct {
	ServiceName string
	Cluster string
	Consul consul.Client
}

func (c *ListClustersController) Run() (ListClustersResponse, error) {
	resp := ListClustersResponse{
		ServiceName: c.ServiceName,
	}
	e, err := envoy.NewFromServiceName(c.ServiceName)
	if err != nil { return resp, err }

	resp.Envoy = &e

	stats, gErr := e.Clusters().GetStatistics(c.Cluster)
	if gErr != nil { return resp, gErr }

	resp.ClusterStatistics = stats
	return resp, nil
}

type ListClustersResponse struct {
	ServiceName string `json:"service"`
	Envoy *envoy.Instance `json:"envoy"`
	ClusterStatistics map[string]envoy.ClusterStatistics `json:"clusters"`
}

func (r ListClustersResponse) String() string {
	o := ""
	for _, stats := range r.ClusterStatistics {
		o += Ok("")
		o += Ok(color.New(color.FgGreen).Add(color.Bold).Sprint(stats.Host))
		o += Ok("------------------------------------------------------------------------------------------------------------")
		d := [][]string{
			{
				"Outlier: Success Rate", stats.Outlier.SuccessRateAverage,
				"Outlier: Success Rate Ejection Threshold", stats.Outlier.SuccessRateEjectionThreshold,
			},
			{
				"Outlier: Local Origin - Success Rate", stats.Outlier.LocalOriginSuccessRateAverage,
				"Outlier: Local Origin - Success Rate Ejection Threshold", stats.Outlier.LocalOriginSuccessRateEjectionThreshold,
			},
			{
				"Default Priority - Max Connections", cast.ToString(stats.DefaultPriority.MaxConnections),
				"Default Priority - Max Retries", cast.ToString(stats.DefaultPriority.MaxRetries),
			},
			{
				"Default Priority - Max Pending Requests", cast.ToString(stats.DefaultPriority.MaxPendingRequests),
				"Default Priority - Max Requests", cast.ToString(stats.DefaultPriority.MaxRequests),
			},
			{
				"High Priority - Max Connections", cast.ToString(stats.HighPriority.MaxConnections),
				"High Priority - Max Retries", cast.ToString(stats.HighPriority.MaxRetries),
			},
			{
				"High Priority - Max Pending Requests", cast.ToString(stats.HighPriority.MaxPendingRequests),
				"High Priority - Max Requests", cast.ToString(stats.HighPriority.MaxRequests),
			},
		}
		tableString := strings.Builder{}
		table := tablewriter.NewWriter(&tableString)
		table.SetBorder(false)
		table.SetTablePadding("\t")
		table.AppendBulk(d)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.Render()
		o += tableString.String()

		if len(stats.Instances) > 0 {
			o += Ok("")
			o += Ok("---------------------")
			o += Ok("- Cluster Instances -")
			o += Ok("---------------------")

			tableString = strings.Builder{}
			table = tablewriter.NewWriter(&tableString)
			table.SetHeader([]string{
				"Host",
				"Cx Active",
				"Cx Failed",
				"Cx Total",
				"Req Active",
				"Req Timeout",
				"Req Success",
				"Req Error",
				"Req Total",
				"Success Rate",
				"Local Success Rate",
				"Health Flags",
				"Region",
				"Zone",
				"SubZone",
				"Canary",
			})

			d = [][]string{}
			for _, i := range stats.Instances {
				d = append(d, []string{
					i.Hostname,
					cast.ToString(i.Connections.Active),
					cast.ToString(i.Connections.Failed),
					cast.ToString(i.Connections.Total),
					cast.ToString(i.Requests.Active),
					cast.ToString(i.Requests.Timeout),
					cast.ToString(i.Requests.Success),
					cast.ToString(i.Requests.Error),
					cast.ToString(i.Requests.Total),
					i.SuccessRate,
					i.LocalOriginSuccessRate,
					i.HealthFlags,
					i.Region,
					i.Zone,
					i.SubZone,
					cast.ToString(i.Canary),
				})
			}
			table.SetBorder(false)
			table.SetTablePadding("\t")
			table.AppendBulk(d)
			table.SetAlignment(tablewriter.ALIGN_RIGHT)
			table.Render()
			o += tableString.String()
		}
	}
	return o
}
