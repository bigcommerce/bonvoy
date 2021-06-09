package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "bonvoy",
	Short: "Envoy and Consul Connect interaction",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Init() {
	BuildServerInfoCommand(rootCmd)
	BuildSetLogLevelCommand(rootCmd)
	BuildListenersCommand(rootCmd)
	BuildExpiredCertificatesController(rootCmd)
	BuildVersionCommand(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

//type Runner interface {
//	Init([]string) error
//	Run() error
//	Name() string
//}
//
//func All() []Runner {
//	return []Runner{
//		BuildListeners(),
//		BuildVersion(),
//		BuildExpiredCertificatesCommand(),
//		BuildSetLogLevelCommand(),
//		BuildServerInfoCommand(),
//	}
//}