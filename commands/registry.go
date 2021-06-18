package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

const outputJson = "json"
const outputText = "text"

type Registry struct {}

func NewRegistry() *Registry {
	return &Registry{}
}

var rootCmd = &cobra.Command{
	Use:   "bonvoy",
	Short: "Envoy and Consul Connect interaction",
}

func (r *Registry) Init() {
	r.RegisterGlobalFlags()
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
	rootCmd.AddCommand(r.Clusters().Command)
	rootCmd.AddCommand(r.Certificates().Command)
	rootCmd.AddCommand(r.Config().Command)
	rootCmd.AddCommand(r.Statistics().Command)
	rootCmd.AddCommand(r.Version().Command)
}

func (r *Registry) RegisterGlobalFlags() {
	rootCmd.PersistentFlags().StringP("output", "o", "", "Optional output format flag. Accepts: 'json'")
}

func (r *Registry) IsJsonOutput() bool {
	return r.GetOutputFormat() == outputJson
}

// Get the desired output format
func (r *Registry) GetOutputFormat() string {
	format, err := rootCmd.PersistentFlags().GetString("output")
	if err != nil || format != outputJson { return outputText }

	return outputJson
}