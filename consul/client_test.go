package consul

import (
	"bonvoy/config"
	"testing"
)

func TestClient_NewClient(t *testing.T) {
	config.Load()

	t.Run("returns a Client from env", func(t *testing.T) {
		client := NewClient()
		if client.address != GetDefaultAddress() {
			t.Error("Invalid client address returned")
		}
	})
}
