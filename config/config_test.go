package config

import (
	"github.com/spf13/viper"
	"testing"
)

func TestConfig_Load(t *testing.T) {
	Load()
	eh := viper.GetString("ENVOY_HOST")
	if eh != "http://0.0.0.0:19001" {
		t.Error("Default ENVOY_HOST is incorrect:", eh)
	}
	eh = viper.GetString("CONSUL_API_HOST")
	if eh != "http://0.0.0.0:8500" {
		t.Error("Default CONSUL_API_HOST is incorrect:", eh)
	}
	eh = viper.GetString("NOMAD_ADDR")
	if eh != "https://0.0.0.0:4646" {
		t.Error("Default NOMAD_ADDR is incorrect:", eh)
	}
}