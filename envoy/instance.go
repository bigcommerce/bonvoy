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

func NewFromServiceName(name string) (Instance, error) {
	dci := docker.NewClient()
	pid, err := dci.GetEnvoyPid(name)
	if err != nil {
		return Instance{}, err
	} else {
		return Instance{
			Address: GetDefaultHost(),
			Pid: pid,
			docker: dci,
			nsenter: nsenter.NewClient(pid),
		}, nil
	}
}