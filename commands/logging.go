package commands

import (
	"bonvoy/envoy"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	logLevel string
)

type SetLogLevelController struct {
	ServiceName string
}

func BuildSetLogLevelCommand(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use: "log level",
		Short: "Set Envoy log level",
		Long:  `Set the Envoy sidecar log level`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			controller := SetLogLevelController{
				ServiceName: args[1],
			}
			return controller.Run()
		},
	}
	cmd.Flags().StringVarP(&logLevel, "level", "l", "", "Desired log level (debug/info/warning/error")
	rootCmd.AddCommand(cmd)
}

func (c *SetLogLevelController) Run() error {
	if logLevel == "" {
		logLevel = c.SelectLogLevel()
	}
	e, err := envoy.NewFromServiceName(c.ServiceName)
	if err != nil {
		return err
	}
	e.Logging().SetLevel(logLevel)
	return nil
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