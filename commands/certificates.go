package commands

import (
	"bonvoy/consul"
	"bonvoy/envoy"
	"fmt"
	"github.com/spf13/cobra"
)

type Certificates struct {
	Command *cobra.Command
}

func (r *Registry) Certificates() *Certificates {
	cmd := &cobra.Command{
		Use: "certificates",
		Short: "Certificates-related commands",
	}
	cmd.AddCommand(r.BuildExpiredCertificatesCommand())
	return &Certificates{
		Command: cmd,
	}
}

func (r *Registry) BuildExpiredCertificatesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "expired",
		Short: "Show all expired certificates",
		Long:  `Display all expired sidecar certificates as compared to the local Consul agent`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			controller := ExpiredCertificatesController{
				ServiceName: args[0],
				Consul: consul.NewClient(),
			}
			restart, err := cmd.Flags().GetBool("restart")
			if err != nil { return err }

			return controller.Run(restart)
		},
	}
	cmd.Flags().BoolP("restart", "r", false, "If passed, will restart all sidecars that have expired certificates")
	return cmd
}

type ExpiredCertificatesController struct {
	ServiceName string
	Consul consul.Client
}

func (c *ExpiredCertificatesController) Run(restart bool) error {
	var expiredCerts []envoy.ExpiredCertificate
	var sidecars []envoy.Instance

	if c.ServiceName != "all" {
		e, err := envoy.NewFromServiceName(c.ServiceName)
		if err != nil { return err }

		sidecars = append(sidecars, e)
	} else {
		sideResp, err := envoy.AllSidecars()
		if err != nil {
			return err
		}

		sidecars = append(sidecars, sideResp...)
	}

	for _, e := range sidecars {
		resp, lErr := e.Certificates().FindExpired()
		if lErr != nil { return lErr }

		expiredCerts = append(expiredCerts, resp...)
	}

	for _, e := range expiredCerts {
		fmt.Println(e.ServiceName)
		fmt.Println("  Envoy Process ID:", e.Pid)
		fmt.Printf("  Envoy Certificate Expiry: %s (%d days)\n", e.EnvoyExpiration, e.EnvoyDaysUntilExpiration)
		fmt.Println("  Consul Agent Certificate Expiry: ", e.ConsulExpiration)

		if restart == true {
			err := e.Envoy.Restart()
			if err != nil { return err }

			fmt.Println("...restarting...Done.")
		}
	}

	return nil
}