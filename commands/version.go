package commands

import (
	"bonvoy/version"
	"github.com/spf13/cobra"
)

type Version struct {
	Command *cobra.Command
}

type GetVersionResponse struct {
	Version string `json:"version"`
}

func (pi GetVersionResponse) String() string {
	return pi.Version
}

func (r *Registry) Version() *Version {
	return &Version{
		Command: &cobra.Command{
			Use:   "version",
			Short: "Display Bonvoy version",
			Long:  `Display the Bonvoy CLI version`,
			RunE: func(cmd *cobra.Command, args []string) error {
				pi := GetVersionResponse{
					Version: version.Version,
				}
				return r.Output(pi)
			},
		},
	}
}