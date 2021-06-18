package envoy

type Logging struct {
	i *Instance
	endpoints LoggingEndpoints
}

type LoggingEndpoints struct {
	logging string
}
func (i *Instance) Logging() *Logging {
	return &Logging{
		i: i,
		endpoints: LoggingEndpoints{
			logging: i.Address + "/logging",
		},
	}
}

func (l *Logging) SetLevel(level string) (string, error) {
	result, err := l.i.nsenter.Curl("-X", "POST", l.endpoints.logging + "?level="+level)
	if err != nil { return "", err }

	return result, nil
}