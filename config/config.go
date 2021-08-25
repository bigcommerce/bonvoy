package config

import "github.com/spf13/viper"

func Load() {
	viper.AutomaticEnv()
	viper.SetDefault("ENVOY_HOST", "http://127.0.0.2:19001")
	viper.SetDefault("CONSUL_API_HOST", "http://0.0.0.0:8500")
	viper.SetDefault("NOMAD_ADDR", "https://0.0.0.0:4646")
}