package api

import (
	"bonvoy/envoy"
	"fmt"
	"github.com/Devatoria/go-nsenter"
	"strings"
)

type Listener struct {
	Name string
	TargetAddress string
}

func GetListeners(config nsenter.Config) []Listener {
	stdout, stderr, err := config.Execute("curl", envoy.GetHost() + "/listeners")
	if err != nil {
		fmt.Println(stderr)
		panic(err)
	}
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