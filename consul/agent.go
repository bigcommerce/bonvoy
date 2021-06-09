package consul

import (
	"encoding/json"
	"net/http"
)

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

	r, err := http.Get(c.address + "/v1/agent/connect/ca/leaf/" + svc)
	if err != nil { return response, nil }

	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&response)
	if err != nil { return response, nil }

	return response, nil
}