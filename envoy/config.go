package envoy

type Config struct {
	i *Instance
	endpoints ConfigEndpoints
}

type ConfigEndpoints struct {
	dump string
}
func (i *Instance) Config() *Config {
	return &Config{
		i: i,
		endpoints: ConfigEndpoints{
			dump: i.Address + "/config_dump",
		},
	}
}

func (l *Config) Dump() (string, error) {
	result, err := l.i.nsenter.Curl(l.endpoints.dump)
	if err != nil { return "", err }

	return result, nil
}