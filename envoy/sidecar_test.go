package envoy

import (
	"github.com/spf13/viper"
	"testing"
)

func TestGetHost(t *testing.T) {
	t.Run("Test GetHost == ENVOY_HOST", func(t *testing.T) {
		testAddr := "http://0.0.0.0:19000"
		viper.Set("ENVOY_HOST", testAddr)
		if GetHost() != testAddr {
			t.Error("ENVOY_HOST does not match")
		}
	})
}