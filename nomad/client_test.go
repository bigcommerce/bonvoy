package nomad

import "testing"

func TestNewClient(t *testing.T) {
	t.Run("returns a Client from env", func(t *testing.T) {
		client := NewClient()
		if client.address != GetDefaultAddress() {
			t.Error("Invalid address set")
		}
	})
}
