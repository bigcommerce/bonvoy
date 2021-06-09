package commands

import (
	"bonvoy/consul"
	"bonvoy/envoy"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

type ExpiredCertificatesController struct {
	ServiceName string
	Consul consul.Client
}

func BuildExpiredCertificatesController(rootCmd *cobra.Command) {
	rootCmd.AddCommand(&cobra.Command{
		Use: "certificates expired",
		Short: "Show all expired certificates",
		Long:  `Display all expired sidecar certificates as compared to the local Consul agent`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			controller := ExpiredCertificatesController{
				ServiceName: args[1],
				Consul: consul.NewClient(),
			}
			return controller.Run()
		},
	})
}

func (c *ExpiredCertificatesController) Run() error {
	e, err := envoy.NewFromServiceName(c.ServiceName)
	if err != nil {
		return err
	}
	data := e.Certificates().Get()

	for _, certs := range data.Certificates {
		for _, cert := range certs.CertificateChain {
			a := strings.Split(cert.SubjectAltNames[0].Uri, "/")
			svc := a[len(a)-1]
			leaf := c.Consul.Agent().GetConnectLeafCaCertificate(svc)

			if cert.ExpirationTime != leaf.ValidBefore {
				fmt.Println(svc)
				fmt.Println("  Envoy Process ID:", e.Pid)
				fmt.Printf("  Envoy Certificate Expiry: %s (%s days)\n", cert.ExpirationTime, cert.DaysUntilExpiration)
				fmt.Println("  Consul Agent Certificate Expiry: ", leaf.ValidBefore)
			}
		}
	}
	return nil
}