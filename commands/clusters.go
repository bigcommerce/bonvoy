package commands

import (
	"bonvoy/consul"
	"bonvoy/docker"
	"bonvoy/envoy"
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"net"
	"sort"
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
	cmd.AddCommand(r.BuildCompareClustersCommand())
	return &Clusters{
		Command: cmd,
	}
}

/***********************************************************************************************************************
 * clusters list [service]
 **********************************************************************************************************************/
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

/***********************************************************************************************************************
 * clusters compare
 **********************************************************************************************************************/
func (r *Registry) BuildCompareClustersCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use: "compare",
		Short: "Compare Envoy clusters with Consul",
		Long:  `Compare Envoy clusters with registered sidecars in Consul`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			controller := CompareClustersController{
				ServiceName: args[0],
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
type CompareClustersController struct {
	ServiceName string
	Cluster string
	Consul consul.Client
}

func (c *CompareClustersController) Run() (CompareClustersResponse, error) {
	resp := CompareClustersResponse{
		ServiceName: c.ServiceName,
	}

	var dockerClient = docker.NewClient()
	container, err := dockerClient.GetSidecarContainer(c.ServiceName)
	if err != nil { return resp, err }

	envoyInstance, eErr := envoy.NewFromSidecarContainer(container)
	if eErr != nil { return resp, err }

	// get all clusters for this container
	var clusterStats, cErr = envoyInstance.Clusters().GetStatistics("")
	if cErr != nil { return resp, cErr }

	var combinedClusters []CombinedCluster
	for _, cluster := range clusterStats {
		if cluster.Host == "local_app" || cluster.Host == "self_admin" || cluster.Host == "local_agent" {
			continue
		}
		clusterAddresses := make(map[string]struct{})
		for _, i := range cluster.Instances {
			clusterAddresses[i.Hostname] = struct{}{}
		}

		var consulName = cluster.GetConsulName()
		if consulName == "bcapp" {
			continue // this does not work with bcapp yet
		}

		consulSidecarServices, err := c.Consul.GetSidecarServices(consulName)
		if err != nil { return resp, err }

		consulAddresses := make(map[string]struct{})
		for _, ic := range consulSidecarServices {
			addr := ic.ServiceAddress + ":" + fmt.Sprintf("%d", ic.ServicePort)
			consulAddresses[addr] = struct{}{}
		}

		matchingAddrs := make([]string, 0)
		consulOnlyAddrs := make([]string, 0)
		envoyOnlyAddrs := make([]string, 0)

		for clusterAddr, _ := range clusterAddresses {
			_, isPresent := consulAddresses[clusterAddr]
			if isPresent {
				matchingAddrs = append(matchingAddrs, clusterAddr)
			} else {
				envoyOnlyAddrs = append(envoyOnlyAddrs, clusterAddr)
			}
		}

		for consulAddr, _ := range consulAddresses {
			_, found := c.inAddresses(matchingAddrs, consulAddr)
			if found {
				continue // we've already matched this address
			}
			_, isPresent := clusterAddresses[consulAddr]
			if isPresent {
				matchingAddrs = append(matchingAddrs, consulAddr)
			} else {
				consulOnlyAddrs = append(consulOnlyAddrs, consulAddr)
			}
		}

		combinedClusters = append(combinedClusters, CombinedCluster{
			Name: cluster.GetConsulName(),
			Host: cluster.Host,
			BothAddresses: matchingAddrs, //c.sortIPs(matchingAddrs),
			ConsulOnlyAddresses: consulOnlyAddrs, //c.sortIPs(consulOnlyAddrs),
			EnvoyOnlyAddresses: envoyOnlyAddrs, //c.sortIPs(envoyOnlyAddrs),
		})
	}

	resp.Clusters = combinedClusters

	return resp, nil
}

func (c *CompareClustersController) inAddresses(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}


func (c *CompareClustersController) sortIPs(ips []string) []string {
	realIPs := make([]net.IP, 0, len(ips))

	for _, ip := range ips {
		realIPs = append(realIPs, net.ParseIP(ip))
	}

	sort.Slice(realIPs, func(i, j int) bool {
		return bytes.Compare(realIPs[i], realIPs[j]) < 0
	})

	var sortedIPs []string
	for _, ip := range realIPs {
		sortedIPs = append(sortedIPs, fmt.Sprintf("%s\n", ip))
	}
	return sortedIPs
}

type CombinedCluster struct {
	Name string `json:"Name"`
	Host string `json:"Host"`
	BothAddresses []string `json:"BothAddresses"`
	ConsulOnlyAddresses []string `json:"ConsulOnlyAddresses"`
	EnvoyOnlyAddresses []string `json:"EnvoyOnlyAddresses"`
}

type CompareClustersResponse struct {
	ServiceName string `json:"ServiceName"`
	Clusters []CombinedCluster `json:"Clusters"`
}
func (r CompareClustersResponse) String() string {
	o := ""
	for _, cluster := range r.Clusters {
		if len(cluster.EnvoyOnlyAddresses) == 0 && len(cluster.ConsulOnlyAddresses) == 0 && len(cluster.BothAddresses) == 0 {
			continue // if no addresses, we can assume this service doesnt exist yet and we can skip it
		}

		o += Ok("")
		name := color.New(color.FgGreen).Add(color.Bold).Sprint(cluster.Name)
		o += Ok(name + " (" + cluster.Host + ")")
		o += Ok("------------------------------------------------------------------------------------------------------------")

		tableString := strings.Builder{}
		table := tablewriter.NewWriter(&tableString)
		table.SetHeader([]string{
			"Matching Addresses",
			"Envoy Only Addresses",
			"Consul Only Addresses",
		})

		var both string
		var eo string
		var co string
		for _, i := range cluster.BothAddresses {
			both += i + "\n"
		}
		for _, i := range cluster.EnvoyOnlyAddresses {
			eo += i + "\n"
		}
		for _, i := range cluster.ConsulOnlyAddresses {
			co += i + "\n"
		}

		d := [][]string{{both, eo, co}}
		table.SetBorder(false)
		table.SetTablePadding("\t")
		table.AppendBulk(d)
		table.SetAlignment(tablewriter.ALIGN_RIGHT)
		table.Render()
		o += tableString.String()
	}
	return o
}