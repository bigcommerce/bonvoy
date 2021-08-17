package consul

type Agent struct {
	c *Client
}

func (c *Client) Agent() *Agent {
	return &Agent{c}
}

type ConnectLeafCaCertificate struct {
	SerialNumber string `json:"SerialNumber"`
	CertPEM string `json:"CertPEM"`
	PrivateKeyPEM string `json:"PrivateKeyPEM"`
	Service string `json:"Service"`
	ServiceURI string `json:"ServiceURI"`
	ValidAfter string `json:"ValidAfter"`
	ValidBefore string `json:"ValidBefore"`
	CreateIndex int `json:"CreateIndex"`
	ModifyIndex int `json:"ModifyIndex"`
}

func (a *Agent) GetConnectLeafCaCertificate(svc string) (ConnectLeafCaCertificate, error) {
	return a.c.GetConnectLeafCaCertificate(svc)
}

func (c *Client) GetConnectLeafCaCertificate(svc string) (ConnectLeafCaCertificate, error) {
	var response ConnectLeafCaCertificate

	req := c.client.Request()
	req.Path("/v1/agent/connect/ca/leaf/" + svc)
	resp, err := req.Send()
	if err != nil { return response, nil }

	jErr := resp.JSON(&response)
	if jErr != nil { return response, jErr }

	return response, nil
}

type SidecarService struct {
	ServiceID string `json:"ServiceID"`
	Node string `json:"Node"`
	Address string `json:"Address"`
	Datacenter string `json:"Datacenter"`
	ServiceAddress string `json:"ServiceAddress"`
	ServicePort int `json:"ServicePort"`
	ServiceProxy ServiceProxy `json:"ServiceProxy"`
}

type ServiceProxy struct {
	DestinationServiceName string `json:"DestinationServiceName"`
	LocalServiceAddress string `json:"LocalServiceAddress"`
	LocalServicePort int `json:"LocalServicePort"`
	Upstreams []ServiceUpstream `json:"Upstreams,omitempty"`
}

type ServiceUpstream struct {
	DestinationType string `json:"DestinationType"`
	DestinationName string `json:"DestinationName"`
	Datacenter string `json:"Datacenter"`
	LocalBindPort int `json:"LocalBindPort"`
}

func (c *Client) GetSidecarServices(svc string) ([]SidecarService, error) {
	var services []SidecarService

	req := c.client.Request()
	req.Path("/v1/catalog/service/" + svc + "-sidecar-proxy")
	resp, err := req.Send()
	if err != nil { return services, nil }

	jErr := resp.JSON(&services)
	if jErr != nil { return services, jErr }

	return services, nil
}