package nsenter

import "testing"

func TestNewClient(t *testing.T) {
	config := NewClient(1)

	t.Run("Ensure --net", func (t *testing.T) {
		if config.config.Net != true {
			t.Error("--net is not set")
		}
	})

	t.Run("Ensure --ipc", func (t *testing.T) {
		if config.config.IPC != true {
			t.Error("--ipc is not set")
		}
	})

	t.Run("Ensure -t", func (t *testing.T) {
		if config.config.Target != 1 {
			t.Error("target is not set")
		}
	})
}
