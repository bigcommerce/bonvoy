package envoy

import (
	"os"
	"testing"
)

func TestGetHost(t *testing.T) {
	t.Run("Test GetHost == ENVOY_HOST", func(t *testing.T) {
		eh := os.Getenv("ENVOY_HOST")
		if GetHost() != eh {
			t.Error("ENVOY_HOST does not match")
		}
	})
}