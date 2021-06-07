package consul

import (
	"encoding/json"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

var consulClient = &http.Client{Timeout: 2 * time.Second}
type AgentConnectLeafCaCertificate struct {
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

func GetHost() string {
	return viper.GetString("CONSUL_API_HOST")
}

func AgentLeafCaCertificate(svc string) AgentConnectLeafCaCertificate {
	r, err := consulClient.Get(GetHost() + "v1/agent/connect/ca/leaf/" + svc)
	if err != nil {
		panic(err)
	}
	var response AgentConnectLeafCaCertificate
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&response)
	return response
}
