package config

import "github.com/spf13/viper"

func Load() {
	viper.AutomaticEnv()
	viper.SetDefault("ENVOY_HOST", "http://0.0.0.0:19001/")
	viper.SetDefault("CONSUL_API_HOST", "http://0.0.0.0:8500/")
}