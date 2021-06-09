package commands

import (
	"bonvoy/envoy"
	"fmt"
	"github.com/spf13/cobra"
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

	fmt.Println("LISTENERS:")
	fmt.Println("----------------------------------------------------------------------")
	for _, listener := range listeners {
		fmt.Printf("%s: %s\n", listener.Name, listener.TargetAddress)
	}
	return nil
}