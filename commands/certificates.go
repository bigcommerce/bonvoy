package commands

import (
	"bonvoy/consul"
	"bonvoy/docker"
	"bonvoy/envoy"
	envoyApi "bonvoy/envoy/api"
	"bonvoy/nsenter"
	"flag"
	"fmt"
	"strings"
)

type ExpiredCertificatesCommand struct {
	fs *flag.FlagSet
	name string
}

func BuildExpiredCertificatesCommand() *ExpiredCertificatesCommand {
	gc := &ExpiredCertificatesCommand{
		fs: flag.NewFlagSet("certs-expired", flag.ContinueOnError),
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
	var cli = docker.NewClient()
	var pid = envoy.GetPid(cli, name)
	config := nsenter.BuildConfig(pid)
	data := envoyApi.GetCertificates(config)

	for _, certs := range data.Certificates {
		for _, c := range certs.CertificateChain {
			a := strings.Split(c.SubjectAltNames[0].Uri, "/")
			svc := a[len(a)-1]
			leaf := consul.AgentLeafCaCertificate(svc)

			if c.ExpirationTime != leaf.ValidBefore {
				fmt.Println(svc)
				fmt.Println("  Envoy Process ID:", pid)
				fmt.Printf("  Envoy Certificate Expiry: %s (%s days)\n", c.ExpirationTime, c.DaysUntilExpiration)
				fmt.Println("  Consul Agent Certificate Expiry: ", leaf.ValidBefore)
			}
		}
	}

	return nil
}
