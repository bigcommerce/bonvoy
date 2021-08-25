package docker

import "testing"

func TestNewClient(t *testing.T) {
	t.Run("returns a Client from env", func(t *testing.T) {
		client := NewClient()
		if client.cli.ClientVersion() != "1.39" {
			t.Error("Invalid client version returned: " + client.cli.ClientVersion() + " != 1.39")
		}
	})
}
