package commands

import (
	"bonvoy/consul"
	"bonvoy/envoy"
	"fmt"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"os"
)

type Certificates struct {
	Command *cobra.Command
}

func (r *Registry) Certificates() *Certificates {
	cmd := &cobra.Command{
		Use: "certificates",
		Short: "Certificates-related commands",
	}
	cmd.AddCommand(r.BuildCertificatesListCommand())
	cmd.AddCommand(r.BuildCertificatesExpiredCommand())
	return &Certificates{
		Command: cmd,
	}
}

func (r *Registry) BuildCertificatesListCommand() *cobra.Command {
	return &cobra.Command{
		Use: "list",
		Short: "Show all certificates",
		Long:  `Display all certificates registered on the Envoy sidecar`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			controller := CertificatesListController{
				ServiceName: args[0],
				Consul: consul.NewClient(),
			}
			return controller.Run()
		},
	}
}

type CertificatesListController struct {
	ServiceName string
	Consul consul.Client
}

func (c *CertificatesListController) Run() error {
	e, err := envoy.NewFromServiceName(c.ServiceName)
	if err != nil { return err }

	certs, cErr := e.Certificates().Get()
	if cErr != nil { return cErr }

	_, _ = color.New(color.FgGreen).Add(color.Bold).Println(c.ServiceName + " Envoy (PID " + cast.ToString(e.Pid) + ")")
	color.Green("-----------------------------------------------------------")

	for _, r := range certs.Certificates {
		// Certs
		color.Green("CertificateChain")
		color.Green("--------------------------------")
		err = c.DisplayCertificateList(r.CertificateChain)
		if err != nil { return err }

		// CA certs
		fmt.Println("")
		color.Green("CA Certificates")
		color.Green("-------------------------------")
		err = c.DisplayCertificateList(r.CaCertificates)
		if err != nil { return err }
	}
	return nil
}

func (c * CertificatesListController) DisplayCertificateList(certs []envoy.Certificate) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"SAN",
		"Serial #",
		"Path",
		"Valid From",
		"Expiration Time",
		"Days Until Expiration",
	})
	table.SetBorder(false)
	table.SetTablePadding("\t")
	readCas := map[string]bool{}
	var car [][]string
	for _, ca := range certs {
		if readCas[ca.SerialNumber] {
			continue // ignore duplicates
		}

		car = append(car, []string{
			ca.SubjectAltNames[0].Uri,
			ca.SerialNumber,
			ca.Path,
			ca.ValidFrom,
			ca.ExpirationTime,
			cast.ToString(ca.DaysUntilExpiration),
		})
	}
	table.AppendBulk(car)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
	return nil
}

// certificates expired

func (r *Registry) BuildCertificatesExpiredCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "expired",
		Short: "Show all expired certificates",
		Long:  `Display all expired sidecar certificates as compared to the local Consul agent`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			controller := CertificatesExpiredController{
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

type CertificatesExpiredController struct {
	ServiceName string
	Consul consul.Client
}

func (c *CertificatesExpiredController) Run(restart bool) error {
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
		color.Green(e.ServiceName)
		fmt.Println("  Envoy Process ID:", e.Pid)
		fmt.Printf("  Envoy Certificate Expiry: %s (%d days)\n", e.EnvoyExpiration, e.EnvoyDaysUntilExpiration)
		fmt.Println("  Consul Agent Certificate Expiry: ", e.ConsulExpiration)

		if restart == true {
			color.Green("    Restarting " + e.ServiceName + " Envoy...")
			err := e.Envoy.Restart()
			if err != nil { return err }

			color.Green("    ...done.")
		}
	}

	return nil
}