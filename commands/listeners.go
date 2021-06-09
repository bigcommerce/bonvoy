package commands

import (
	"bonvoy/envoy"
	"flag"
	"fmt"
	"github.com/spf13/cobra"
)

type ListenersController struct {
	ServiceName string
}

func BuildListenersCommand(rootCmd *cobra.Command) {
	rootCmd.AddCommand(&cobra.Command{
		Use: "listeners",
		Short: "Show Envoy listeners",
		Long:  `Display all registered Envoy sidecar listeners`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			controller := ListenersController{
				ServiceName: args[1],
			}
			return controller.Run()
		},
	})
}

func (s *ListenersController) Run() error {
	e, err := envoy.NewFromServiceName(s.ServiceName)
	if err != nil {
		return err
	}
	listeners := e.Listeners().Get()
	fmt.Println("LISTENERS:")
	fmt.Println("----------------------------------------------------------------------")
	for _, listener := range listeners {
		fmt.Printf("%s: %s\n", listener.Name, listener.TargetAddress)
	}
	return nil
}
type ListenersCommand struct {
	fs *flag.FlagSet
	name string
}