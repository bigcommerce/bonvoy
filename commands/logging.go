package commands

import (
	"bonvoy/envoy"
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

func GetAvailableLogLevels() []string {
	return []string{"debug", "info", "warning", "error"}
}

// log level

func (r *Registry) BuildSetLogLevelCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "level",
		Short: "Set Envoy log level",
		Long:  `Set the Envoy sidecar log level`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			level, fErr := cmd.Flags().GetString("level")
			if fErr != nil { return fErr }

			controller := SetLogLevelController{
				ServiceName: args[0],
				DesiredLogLevel: level,
			}

			o, err := controller.Run()
			if err != nil { return err }

			return r.Output(o)
		},
	}
	cmd.Flags().StringP("level", "l", "", "Desired log level (debug/info/warning/error")
	return cmd
}

type SetLogLevelController struct {
	ServiceName string
	DesiredLogLevel string
}

func (c *SetLogLevelController) Run() (SetLogLevelResponse, error) {
	resp := SetLogLevelResponse{
		ServiceName: c.ServiceName,
	}
	if c.DesiredLogLevel == "" {
		d, sErr := c.SelectLogLevel()
		if sErr != nil { return resp, sErr }
		c.DesiredLogLevel = d
	}
	resp.Level = c.DesiredLogLevel

	e, err := envoy.NewFromServiceName(c.ServiceName)
	if err != nil { return resp, err }

	resp.Envoy = &e

	result, sErr := e.Logging().SetLevel(c.DesiredLogLevel)
	if sErr != nil { return resp, err }

	resp.Output = result
	return resp, nil
}

type SetLogLevelResponse struct {
	ServiceName string `json:"service"`
	Envoy *envoy.Instance `json:"envoy"`
	Level string `json:"level"`
	Output string `json:"-"`
}

func (r SetLogLevelResponse) String() string {
	o := Ok("Set Log Level for", r.ServiceName)
	o += Ok("------------------------------------------------")
	o += Info(r.Output)
	return o
}

func (c *SetLogLevelController) SelectLogLevel() (string, error) {
	prompt := promptui.Select{
		Label: "Please Select a Log Level",
		Items: GetAvailableLogLevels(),
	}
	_, desired, err := prompt.Run()
	if err != nil { return "info", err }
	return desired, nil
}