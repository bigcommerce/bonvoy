package envoy

import (
	"github.com/spf13/viper"
	"testing"
)

func TestGetDefaultHost(t *testing.T) {
	t.Run("Test GetHost == ENVOY_HOST", func(t *testing.T) {
		testAddr := "http://127.0.0.2:19000"
		viper.Set("ENVOY_HOST", testAddr)
		if GetDefaultHost() != testAddr {
			t.Error("ENVOY_HOST does not match")
		}
	})
}