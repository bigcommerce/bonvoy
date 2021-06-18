package consul

import (
	"bonvoy/config"
	"github.com/stretchr/testify/require"
	mock "gopkg.in/h2non/gentleman-mock.v2"
	"testing"
)

func TestClient_GetConnectLeafCaCertificate(t *testing.T) {
	defer mock.Disable()
	config.Load()

	mock.New(GetDefaultAddress()).
		Get("/v1/agent/connect/ca/leaf/auth-grpc").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		JSON(map[string]interface{}{
		"SerialNumber": "01:7e",
		"CertPEM": "test-cert",
		"PrivateKeyPEM": "test-pk",
		"Service": "auth-grpc",
		"ServiceURI": "spiffe://89749198-098b-d4bc-802b-04d28ed8af0a.consul/ns/default/dc/youngeducated/svc/auth-grpc",
		"ValidAfter": "2021-06-07T17:48:13Z",
		"ValidBefore": "2021-06-10T17:48:13Z",
		"CreateIndex": 271994,
		"ModifyIndex": 271994,
	})

	client := NewClient()
	client.client.Use(mock.Plugin)
	res, err := client.GetConnectLeafCaCertificate("auth-grpc")
	if err != nil {
		t.Error(err)
	}
	if res.Service != "auth-grpc" {
		t.Error("Failed to return proper Service name")
	}

	expectedCert := ConnectLeafCaCertificate{
		SerialNumber: "01:7e",
		CertPEM: "test-cert",
		PrivateKeyPEM: "test-pk",
		Service: "auth-grpc",
		ServiceURI: "spiffe://89749198-098b-d4bc-802b-04d28ed8af0a.consul/ns/default/dc/youngeducated/svc/auth-grpc",
		ValidAfter: "2021-06-07T17:48:13Z",
		ValidBefore: "2021-06-10T17:48:13Z",
		CreateIndex: 271994,
		ModifyIndex: 271994,
	}
	require.Equal(t, expectedCert, res)
}