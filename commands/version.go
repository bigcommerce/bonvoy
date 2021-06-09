package commands

import (
	"bonvoy/version"
	"fmt"
	"github.com/spf13/cobra"
)

func BuildVersionCommand(rootCmd *cobra.Command) {
	rootCmd.AddCommand(&cobra.Command{
		Use: "version",
		Short: "Display Bonvoy version",
		Long:  `Display the Bonvoy CLI version`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(version.Version)
			return nil
		},
	})
}