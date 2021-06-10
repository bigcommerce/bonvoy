package envoy

import (
	"bonvoy/docker"
	"bonvoy/nsenter"
	"github.com/spf13/viper"
	"os/exec"
	"strconv"
	"strings"
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

func AllSidecars() ([]Instance, error) {
	var sidecars []Instance
	dci := docker.NewClient()

	pids, err := GetAllProcessIds()
	if err != nil { return sidecars, err }

	defaultHost := GetDefaultHost()

	for _, pid := range pids {
		sidecars = append(sidecars, Instance{
			Address: defaultHost,
			Pid:     pid,
			docker:  dci,
			nsenter: nsenter.NewClient(pid),
		})
	}
	return sidecars, nil
}

func GetAllProcessIds() ([]int, error) {
	var pids []int

	out, err := exec.Command("ps", "-C", "envoy", "-o", "pid", "--no-headers").Output()
	if err != nil { return pids, err }

	data := strings.Split(string(out), "\n")
	for _, o := range data {
		pidTrim := strings.TrimSpace(o)
		pid, pidErr := strconv.Atoi(pidTrim)
		if pidErr != nil { continue }
		if pid < 1 { continue }

		pids = append(pids, pid)
	}

	return pids, nil
}