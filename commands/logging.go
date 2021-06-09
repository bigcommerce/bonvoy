package commands

import (
	"bonvoy/envoy"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type Logging struct {
	Command *cobra.Command
}

func (r *Registry) Logging() *Logging {
	cmd := &cobra.Command{
		Use: "log",
		Short: "Logging related commands",
	}
	cmd.AddCommand(r.BuildSetLogLevelCommand())
	return &Logging{
		Command: cmd,
	}
}

// log level

var (
	desiredLogLevel string
)

func (r *Registry) BuildSetLogLevelCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "level",
		Short: "Set Envoy log level",
		Long:  `Set the Envoy sidecar log level`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			controller := SetLogLevelController{
				ServiceName: args[0],
			}
			return controller.Run()
		},
	}
	cmd.Flags().StringVarP(&desiredLogLevel, "level", "l", "", "Desired log level (debug/info/warning/error")
	return cmd
}

type SetLogLevelController struct {
	ServiceName string
}

func (c *SetLogLevelController) Run() error {
	if desiredLogLevel == "" {
		desiredLogLevel = c.SelectLogLevel()
	}
	e, err := envoy.NewFromServiceName(c.ServiceName)
	if err != nil { return err }

	return e.Logging().SetLevel(desiredLogLevel)
}

func (c *SetLogLevelController) SelectLogLevel() string {
	prompt := promptui.Select{
		Label: "Please Select a Log Level",
		Items: []string{"debug", "info", "warning", "error"},
	}
	_, desired, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "info"
	}
	return desired
}