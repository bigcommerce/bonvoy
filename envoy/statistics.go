package envoy

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