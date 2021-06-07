package nsenter

import "github.com/Devatoria/go-nsenter"

func BuildConfig(pid int) nsenter.Config {
	return nsenter.Config{
		Net: true,
		IPC: true,
		Target: pid,
	}
}