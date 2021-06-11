package commands

import (
	"bonvoy/consul"
	"bonvoy/envoy"
	"fmt"
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
			return controller.Run()
		},
	}
	cmd.Flags().BoolP("restart", "r", false, "If passed, will restart all sidecars that have expired certificates")
	return cmd
}

type ConfigDumpController struct {
	ServiceName string
	Consul consul.Client
}

func (c *ConfigDumpController) Run() error {
	e, err := envoy.NewFromServiceName(c.ServiceName)
	if err != nil { return err }

	result, cErr := e.Config().Dump()
	if cErr != nil { return cErr }

	fmt.Println(result)
	return nil
}