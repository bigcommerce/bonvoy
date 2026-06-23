package docker

import (
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	t.Run("returns a Client from env", func(t *testing.T) {
		_ = os.Setenv("DOCKER_API_VERSION", "1.40")
		client := NewClient()
		if client.cli.ClientVersion() != "1.40" {
			t.Error("Invalid client version returned: " + client.cli.ClientVersion() + " != 1.40")
		}
	})
}
