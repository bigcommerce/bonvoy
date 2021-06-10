package commands

import (
	"bonvoy/consul"
	"bonvoy/envoy"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
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
	return &cobra.Command{
		Use: "expired",
		Short: "Show all expired certificates",
		Long:  `Display all expired sidecar certificates as compared to the local Consul agent`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			controller := ExpiredCertificatesController{
				ServiceName: args[0],
				Consul: consul.NewClient(),
			}
			return controller.Run()
		},
	}
}

type ExpiredCertificatesController struct {
	ServiceName string
	Consul consul.Client
}

func (c *ExpiredCertificatesController) Run() error {
	if c.ServiceName != "all" {
		e, err := envoy.NewFromServiceName(c.ServiceName)
		if err != nil { return err }

		return c.LookupForSidecar(e)
	} else {
		sidecars, err := envoy.AllSidecars()
		if err != nil { return err }

		for _, e := range sidecars {
			err = c.LookupForSidecar(e)
			if err != nil { return err }
		}
	}

	return nil
}

func (c *ExpiredCertificatesController) LookupForSidecar(e envoy.Instance) error {
	data, err := e.Certificates().Get()
	if err != nil { return err }

	readSerials := map[string]bool{}

	for _, certs := range data.Certificates {
		for _, cert := range certs.CertificateChain {
			if readSerials[cert.SerialNumber] {
				continue // ignore duplicates
			}
			readSerials[cert.SerialNumber] = true
			a := strings.Split(cert.SubjectAltNames[0].Uri, "/")
			svc := a[len(a)-1]
			leaf, lErr := c.Consul.Agent().GetConnectLeafCaCertificate(svc)
			if lErr != nil { return lErr }

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