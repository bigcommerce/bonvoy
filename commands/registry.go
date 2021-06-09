package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

type Registry struct {}

func NewRegistry() *Registry {
	return &Registry{}
}

var rootCmd = &cobra.Command{
	Use:   "bonvoy",
	Short: "Envoy and Consul Connect interaction",
}

func (r *Registry) Init() {
	r.RegisterCommands()
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (r *Registry) RegisterCommands() {
	// Root level subcommands/commands
	rootCmd.AddCommand(r.Server().Command)
	rootCmd.AddCommand(r.Logging().Command)
	rootCmd.AddCommand(r.Listeners().Command)
	rootCmd.AddCommand(r.Certificates().Command)
	rootCmd.AddCommand(r.Version().Command)
}