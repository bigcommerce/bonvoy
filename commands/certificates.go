package commands

import (
	"bonvoy/consul"
	"bonvoy/envoy"
	"flag"
	"fmt"
	"strings"
)

type ExpiredCertificatesCommand struct {
	fs *flag.FlagSet
	consul consul.Client
	name string
}

func BuildExpiredCertificatesCommand() *ExpiredCertificatesCommand {
	gc := &ExpiredCertificatesCommand{
		fs: flag.NewFlagSet("certs-expired", flag.ContinueOnError),
		consul: consul.NewClient(),
	}
	gc.fs.Arg(0)
	return gc
}

func (g *ExpiredCertificatesCommand) Name() string {
	return g.fs.Name()
}

func (g *ExpiredCertificatesCommand) Init(args []string) error {
	return g.fs.Parse(args)
}

func (g *ExpiredCertificatesCommand) Run() error {
	var name = g.name
	if name == "" {
		name = g.fs.Arg(0)
	}

	e, err := envoy.NewFromServiceName(name)
	if err != nil {
		return err
	}
	data := e.Certificates().Get()

	for _, certs := range data.Certificates {
		for _, c := range certs.CertificateChain {
			a := strings.Split(c.SubjectAltNames[0].Uri, "/")
			svc := a[len(a)-1]
			leaf := g.consul.Agent().GetConnectLeafCaCertificate(svc)

			if c.ExpirationTime != leaf.ValidBefore {
				fmt.Println(svc)
				fmt.Println("  Envoy Process ID:", e.Pid)
				fmt.Printf("  Envoy Certificate Expiry: %s (%s days)\n", c.ExpirationTime, c.DaysUntilExpiration)
				fmt.Println("  Consul Agent Certificate Expiry: ", leaf.ValidBefore)
			}
		}
	}

	return nil
}
