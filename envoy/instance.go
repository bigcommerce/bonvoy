package envoy

import (
	"bonvoy/docker"
	"bonvoy/nsenter"
	"github.com/spf13/viper"
)

type Instance struct {
	Address string
	Pid int
	docker docker.Client
	nsenter nsenter.Client
}

func GetDefaultHost() string {
	return viper.GetString("ENVOY_HOST")
}

func NewFromServiceName(name string) Instance {
	dci := docker.NewClient()
	pid := dci.GetEnvoyPid(name)
	nse := nsenter.NewClient(pid)

	return Instance{
		Address: GetDefaultHost(),
		Pid: pid,
		docker: dci,
		nsenter: nse,
	}
}