package consul

import (
	"bonvoy/config"
	"fmt"
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

func TestClient_GetSidecarServices(t *testing.T) {
	defer mock.Disable()
	config.Load()

	mock.New(GetDefaultAddress()).
		Get("/v1/catalog/service/orders-grpc-sidecar-proxy").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		JSON([]byte(`[
	{
		"ID": "0844560c-ebb3-5b80-b618-0fb1e5068e97",
		"Node": "nomad-client-xxxx",
		"Address": "1.2.3.4",
		"Datacenter": "int-us-central1",
		"TaggedAddresses": {
			"lan": "1.2.3.4",
			"lan_ipv4": "1.2.3.4",
			"wan": "1.2.3.4",
			"wan_ipv4": "1.2.3.4"
		},
		"NodeMeta": {
			"cluster": "nomad",
			"consul-network-segment": "",
			"machine_type": "custom-12-53248",
			"zone": "us-central1-b"
		},
		"ServiceKind": "connect-proxy",
		"ServiceID": "_nomad-task-93ea0de1-3bd8-4aad-2758-beb04cdfa852-group-orders-rpc-orders-grpc-orders9778-sidecar-proxy",
		"ServiceName": "orders-grpc-sidecar-proxy",
		"ServiceTags": [
			"sidecar-for:orders-rpc"
		],
		"ServiceAddress": "1.2.3.4",
		"ServiceTaggedAddresses": {
			"lan_ipv4": {
				"Address": "1.2.3.4",
				"Port": 20305
			},
			"wan_ipv4": {
				"Address": "1.2.3.4",
				"Port": 20305
			}
		},
		"ServiceWeights": {
			"Passing": 1,
			"Warning": 1
		},
		"ServiceMeta": {
			"envoy-stats-port": "23031",
			"external-source": "nomad"
		},
		"ServicePort": 20305,
		"ServiceSocketPath": "",
		"ServiceEnableTagOverride": false,
		"ServiceProxy": {
			"DestinationServiceName": "orders-grpc",
			"DestinationServiceID": "_nomad-task-93ea0de1-3bd8-4aad-2758-beb04cdfa852-group-orders-rpc-orders-grpc-orders9778",
			"LocalServiceAddress": "127.0.0.1",
			"LocalServicePort": 9778,
			"Mode": "",
			"Config": {
				"bind_address": "0.0.0.0",
				"bind_port": 20305,
				"envoy_stats_bind_addr": "0.0.0.0:1239",
				"local_connect_timeout_ms": 5000,
				"local_request_timeout_ms": 0,
				"protocol": "grpc"
			},
			"Upstreams": [{
			  "DestinationType": "service",
			  "DestinationName": "bcapp",
			  "Datacenter": "",
			  "LocalBindPort": 8000,
			  "Config": {
				"protocol": "http"
			  },
			  "MeshGateway": {}
			}],
			"MeshGateway": {},
			"Expose": {}
		},
		"ServiceConnect": {},
		"CreateIndex": 333400859,
		"ModifyIndex": 333400859
	}]`))

	client := NewClient()
	client.client.Use(mock.Plugin)
	res, err := client.GetSidecarServices("orders-grpc")
	if err != nil {
		t.Error(err)
		return
	}
	if len(res) == 0 {
		t.Error("Empty services returned - should have returned 1")
		return
	}
	fmt.Printf("Slice: %v\n", res)

	var expectedServices []SidecarService
	var expectedUpstreams []ServiceUpstream
	expectedUpstreams = append(expectedUpstreams, ServiceUpstream{
		DestinationName: "bcapp",
		DestinationType: "service",
		Datacenter: "",
		LocalBindPort: 8000,
	})
	expectedServices = append(expectedServices, SidecarService{
		ServiceID: "_nomad-task-93ea0de1-3bd8-4aad-2758-beb04cdfa852-group-orders-rpc-orders-grpc-orders9778-sidecar-proxy",
		Node: "nomad-client-xxxx",
		Address: "1.2.3.4",
		Datacenter: "int-us-central1",
		ServiceAddress: "1.2.3.4",
		ServicePort: 20305,
		ServiceProxy: ServiceProxy{
			DestinationServiceName: "orders-grpc",
			LocalServiceAddress: "127.0.0.1",
			LocalServicePort:9778,
			Upstreams: expectedUpstreams,
		},
	})
	require.Equal(t, expectedServices, res)
}