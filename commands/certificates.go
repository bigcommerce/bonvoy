package commands

import (
	"bonvoy/consul"
	"bonvoy/envoy"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cast"
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
	cmd.AddCommand(r.BuildCertificatesListCommand())
	cmd.AddCommand(r.BuildCertificatesExpiredCommand())
	return &Certificates{
		Command: cmd,
	}
}

/***********************************************************************************************************************
 * certificates list [service]
 **********************************************************************************************************************/

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
			o, cErr := controller.Run()
			if cErr != nil { return cErr }

			return r.Output(o)
		},
	}
}

type CertificatesListController struct {
	ServiceName string
	Consul consul.Client
}

func (c *CertificatesListController) Run() (ListCertificatesResponse, error) {
	resp := ListCertificatesResponse{
		ServiceName: c.ServiceName,
	}
	e, err := envoy.NewFromServiceName(c.ServiceName)
	if err != nil { return resp, err }

	resp.Envoy = &e

	certs, cErr := e.Certificates().Get()
	if cErr != nil { return resp, cErr }

	var certificateChains []envoy.Certificate
	var caCertificates []envoy.Certificate

	for _, r := range certs.Certificates {
		certificateChains = append(certificateChains, r.CertificateChain...)
		caCertificates = append(caCertificates, r.CaCertificates...)
	}
	resp.CertificateChains = certificateChains
	resp.CaCertificates = caCertificates
	return resp, nil
}

type ListCertificatesResponse struct {
	ServiceName string `json:"service"`
	Envoy *envoy.Instance `json:"envoy"`
	CertificateChains []envoy.Certificate`json:"certificate_chains"`
	CaCertificates []envoy.Certificate`json:"ca_certificates"`
}

func (r ListCertificatesResponse) String() string {
	o := ""
	o += Ok("----------------------------------------------------------------------------")
	o += Ok(r.ServiceName + " Envoy (PID " + cast.ToString(r.Envoy.Pid) + ")")
	o += Ok("----------------------------------------------------------------------------")
	o += Info("Certificate Chains:")
	o += r.DisplayCertificateList(r.CertificateChains)
	o += Info("")
	o += Info("CA Certificates:")
	o += r.DisplayCertificateList(r.CaCertificates)
	return o
}

func (r *ListCertificatesResponse) DisplayCertificateList(certs []envoy.Certificate) string {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
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
	return tableString.String()
}

/***********************************************************************************************************************
 * certificates expired [service]
 **********************************************************************************************************************/
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

			o, cErr := controller.Run(restart)
			if cErr != nil { return cErr }

			return r.Output(o)
		},
	}
	cmd.Flags().BoolP("restart", "r", false, "If passed, will restart all sidecars that have expired certificates")
	return cmd
}

type CertificatesExpiredController struct {
	ServiceName string
	Consul consul.Client
}

type FindExpiredCertificatesResponse struct {
	ExpiredCertificates []envoy.ExpiredCertificate `json:"expiredCertificates"`
	restart             bool
}

func (r FindExpiredCertificatesResponse) String() string {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{
		"Service",
		"PID",
		"Envoy Expiry",
		"Days Left",
		"Consul Leaf Expiry",
		"Restarted",
	})
	restartStr := "NO"
	if r.restart { restartStr = "YES" }

	var d [][]string
	for _, e := range r.ExpiredCertificates {
		d = append(d, []string{
			e.ServiceName,
			fmt.Sprintf("%d", e.Pid),
			e.EnvoyExpiration,
			fmt.Sprintf("%d", e.EnvoyDaysUntilExpiration),
			e.ConsulExpiration,
			restartStr,
		})
	}

	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.AppendBulk(d)
	table.SetAlignment(tablewriter.ALIGN_RIGHT)
	table.Render()
	return tableString.String()
}

func (c *CertificatesExpiredController) Run(restart bool) (FindExpiredCertificatesResponse, error) {
	var expiredCerts []envoy.ExpiredCertificate
	var sidecars []envoy.Instance
	runResp := FindExpiredCertificatesResponse{
		restart: restart,
	}

	if c.ServiceName != "all" {
		e, err := envoy.NewFromServiceName(c.ServiceName)
		if err != nil { return runResp, err }

		sidecars = append(sidecars, e)
	} else {
		sideResp, err := envoy.AllSidecars()
		if err != nil {
			return runResp, err
		}

		sidecars = append(sidecars, sideResp...)
	}

	for _, e := range sidecars {
		resp, lErr := e.Certificates().FindExpired()
		if lErr != nil { return runResp, lErr }

		expiredCerts = append(expiredCerts, resp...)
	}

	runResp.ExpiredCertificates = expiredCerts

	for _, e := range expiredCerts {
		if restart == true {
			err := e.Envoy.Restart()
			if err != nil {
				return runResp, err
			}
		}
	}

	return runResp, nil
}