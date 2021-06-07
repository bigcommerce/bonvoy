package api

import (
	"bonvoy/envoy"
	"encoding/json"
	"fmt"
	"github.com/Devatoria/go-nsenter"
	"strings"
)

type CertsResponse struct {
	Certificates []EnvoyCertificates `json:"certificates"`
}
type EnvoyCertificates struct {
	CaCertificates []Certificate `json:"ca_cert"`
	CertificateChain []Certificate `json:"cert_chain"`
}

type SubjectAltName struct {
	Uri string `json:"uri"`
}
type Certificate struct {
	Path string `json:"path"`
	SerialNumber string `json:"serial_number"`
	SubjectAltNames []SubjectAltName `json:"subject_alt_names"`
	DaysUntilExpiration string `json:"days_until_expiration"`
	ValidFrom string `json:"valid_from"`
	ExpirationTime string `json:"expiration_time"`
}
type CertificateChain struct {
	Path string `json:"path"`
	SerialNumber string `json:"serial_number"`
	SubjectAltNames []SubjectAltName `json:"subject_alt_names"`
}

func GetCertificates(config nsenter.Config) CertsResponse {
	rawJson, stderr, err := config.Execute("curl", "-s", envoy.GetHost() + "certs")
	jsonData := []byte(strings.Trim(rawJson, " "))

	var response CertsResponse
	err = json.Unmarshal(jsonData, &response)
	if err != nil {
		fmt.Println(stderr)
		panic(err)
	}
	return response
}