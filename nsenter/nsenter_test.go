package nsenter

import "testing"

func TestBuildConfig(t *testing.T) {
	t.Run("Ensure --net", func (t *testing.T) {
		config := BuildConfig(1)
		if config.Net != true {
			t.Error("--net is not set")
		}
	})

	t.Run("Ensure --ipc", func (t *testing.T) {
		config := BuildConfig(1)
		if config.IPC != true {
			t.Error("--ipc is not set")
		}
	})

	t.Run("Ensure -t", func (t *testing.T) {
		config := BuildConfig(1)
		if config.Target != 1 {
			t.Error("target is not set")
		}
	})
}
