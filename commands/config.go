package commands

import (
	"bonvoy/consul"
	"bonvoy/envoy"
	"encoding/json"
	"github.com/spf13/cobra"
)

type Config struct {
	Command *cobra.Command
}

func (r *Registry) Config() *Config {
	cmd := &cobra.Command{
		Use: "config",
		Short: "Envoy configuration commands",
	}
	cmd.AddCommand(r.BuildConfigDumpCommand())
	return &Config{
		Command: cmd,
	}
}

func (r *Registry) BuildConfigDumpCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "dump",
		Short: "Output the Envoy config",
		Long:  `Output the entire Envoy configuration JSON blob`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			controller := ConfigDumpController{
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

type ConfigDumpController struct {
	ServiceName string
	Consul consul.Client
}

func (c *ConfigDumpController) Run() (ConfigDumpResponse, error) {
	resp := ConfigDumpResponse{
		ServiceName: c.ServiceName,
	}
	e, err := envoy.NewFromServiceName(c.ServiceName)
	if err != nil { return resp, err }

	resp.Envoy = &e

	result, cErr := e.Config().Dump()
	if cErr != nil { return resp, cErr }

	var r map[string]interface{}
	err = json.Unmarshal([]byte(result), &r)
	if err != nil { return resp, err }

	resp.Output = r
	return resp, nil
}

type ConfigDumpResponse struct {
	ServiceName string `json:"service"`
	Envoy *envoy.Instance `json:"envoy"`
	Output map[string]interface{} `json:"output,json"`
}

func (r ConfigDumpResponse) String() string {
	v, _ := json.MarshalIndent(r.Output, "", "  ")
	return string(v)
}