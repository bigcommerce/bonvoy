package nomad

import (
	"bonvoy/config"
	"gopkg.in/h2non/gentleman-mock.v2"
	"testing"
)

func TestAllocations_RestartAllocation(t *testing.T) {
	defer mock.Disable()
	config.Load()

	mock.New(GetDefaultAddress()).
		Post("/v1/client/allocation/asdf-1234/restart").
		Reply(200).
		JSON(map[string]string{})

	client := NewClient()
	client.client.Use(mock.Plugin)
	err := client.RestartAllocation("asdf-1234")
	if err != nil {
		t.Error(err)
	}
}