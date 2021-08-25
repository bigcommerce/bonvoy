package envoy

import (
	"bonvoy/docker"
	"bonvoy/nomad"
	"bonvoy/nsenter"
	"github.com/docker/docker/api/types"
	"github.com/spf13/viper"
	"os/exec"
	"strconv"
	"strings"
)

type Instance struct {
	Address string `json:"-"`
	Pid int `json:"pid"`
	NomadAllocationID string `json:"nomad_allocation_id"`
	Nomad nomad.Client `json:"-"`
	docker docker.Client
	nsenter nsenter.Client
}

func GetDefaultHost() string {
	return viper.GetString("ENVOY_HOST")
}

func NewFromServiceName(name string) (Instance, error) {
	dci := docker.NewClient()
	container, err := dci.GetSidecarContainer(name)
	if err != nil { return Instance{}, err }

	return NewFromSidecarContainer(container)
}

func NewFromSidecarContainer(container types.ContainerJSON) (Instance, error) {
	dci := docker.NewClient()

	return Instance{
		Address: GetDefaultHost(),
		Pid: container.State.Pid,
		NomadAllocationID: container.Config.Labels["com.hashicorp.nomad.alloc_id"],
		docker: dci,
		Nomad: nomad.NewClient(),
		nsenter: nsenter.NewClient(container.State.Pid),
	}, nil
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

// Restarts the Envoy instance
func (i *Instance) Restart() error {
	err := i.Nomad.Allocations().Restart(i.NomadAllocationID)
	if err != nil { return err }

	return nil
}