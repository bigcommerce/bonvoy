package consul

import (
	"bonvoy/test"
	"net/http"
	"testing"
	"github.com/stretchr/testify/require"
)

func AgentLeafCaCertificateMock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{
  "SerialNumber": "01:7e",
  "CertPEM": "test-cert",
  "PrivateKeyPEM": "test-pk",
  "Service": "auth-grpc",
  "ServiceURI": "spiffe://89749198-098b-d4bc-802b-04d28ed8af0a.consul/ns/default/dc/youngeducated/svc/auth-grpc",
  "ValidAfter": "2021-06-07T17:48:13Z",
  "ValidBefore": "2021-06-10T17:48:13Z",
  "CreateIndex": 271994,
  "ModifyIndex": 271994
}`))
}

func TestClient_GetConnectLeafCaCertificate(t *testing.T) {
	srv := test.ServerMock("/v1/agent/connect/ca/leaf/auth-grpc", AgentLeafCaCertificateMock)
	defer srv.Close()
	client := Client{address: srv.URL}

	cert := ConnectLeafCaCertificate{
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
	res, err := client.GetConnectLeafCaCertificate("auth-grpc")
	if err != nil {
		t.Error(err)
	}
	if res.Service != "auth-grpc" {
		t.Error("Failed to return proper Service name")
	}
	require.Equal(t, cert, res)
}