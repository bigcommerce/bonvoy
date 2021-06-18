package envoy

import (
	"sort"
	"strings"
)

type Statistics struct {
	i *Instance
	endpoints StatisticsEndpoints
}

type StatisticsEndpoints struct {
	dump string
}
func (i *Instance) Statistics() *Statistics {
	return &Statistics{
		i: i,
		endpoints: StatisticsEndpoints{
			dump: i.Address + "/stats",
		},
	}
}

func (l *Statistics) Dump() (string, error) {
	result, err := l.i.nsenter.Curl(l.endpoints.dump)
	if err != nil { return "", err }

	return result, nil
}

type Statistic struct {
	Name string
	Value string
}

func (l *Statistics) List() ([]Statistic, error) {
	var stats []Statistic

	result, err := l.i.nsenter.Curl(l.endpoints.dump)
	if err != nil { return stats, err }

	ss := strings.Split(result, "\n")
	sort.Strings(ss)
	for _, s := range ss {
		stat := strings.Split(s, ": ")
		if len(stat) < 2 { continue }

		stats = append(stats, Statistic{
			Name: stat[0],
			Value: stat[1],
		})
	}

	return stats, nil
}