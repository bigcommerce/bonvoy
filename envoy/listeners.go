package envoy

import (
	"strings"
)

type Listeners struct {
	i *Instance
	endpoints ListenersEndpoints
}

type ListenersEndpoints struct {
	list string
}

func (i *Instance) Listeners() *Listeners {
	return &Listeners{
		i: i,
		endpoints: ListenersEndpoints{
			list: i.Address + "/listeners",
		},
	}
}

type Listener struct {
	Name string
	TargetAddress string
}

func (l *Listeners) Get() []Listener {
	stdout := l.i.nsenter.Curl(l.endpoints.list)
	data := strings.Split(stdout, "\n")
	var listeners []Listener
	for _, str := range data {
		maps := strings.Split(str, "::")
		if len(maps) == 2 {
			listeners = append(listeners, Listener{
				Name:          maps[0],
				TargetAddress: maps[1],
			})
		}
	}
	return listeners
}