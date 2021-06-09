package envoy

import (
	"encoding/json"
	"strings"
)

type Certificates struct {
	i *Instance
	endpoints CertificatesEndpoints
}

type CertificatesEndpoints struct {
	list string
}

func (i *Instance) Certificates() *Certificates {
	return &Certificates{
		i: i,
		endpoints: CertificatesEndpoints{
			list: i.Address + "/certs",
		},
	}
}

type CertsResponse struct {
	Certificates []CertificateConfig `json:"certificates"`
}
type CertificateConfig struct {
	CaCertificates []Certificate   `json:"ca_cert"`
	CertificateChain []Certificate `json:"cert_chain"`
}
type SubjectAltName struct {
	Uri string `json:"uri"`
}
type Certificate struct {
	Path string                      `json:"path"`
	SerialNumber string              `json:"serial_number"`
	SubjectAltNames []SubjectAltName `json:"subject_alt_names"`
	DaysUntilExpiration string       `json:"days_until_expiration"`
	ValidFrom string                 `json:"valid_from"`
	ExpirationTime string            `json:"expiration_time"`
}
type CertificateChain struct {
	Path string                      `json:"path"`
	SerialNumber string              `json:"serial_number"`
	SubjectAltNames []SubjectAltName `json:"subject_alt_names"`
}

func (c *Certificates) Get() (CertsResponse, error) {
	rawJson, err := c.i.nsenter.Curl("-s", c.endpoints.list)

	jsonData := []byte(strings.Trim(rawJson, " "))

	var response CertsResponse
	err = json.Unmarshal(jsonData, &response)
	if err != nil {
		return CertsResponse{}, err
	}
	return response, nil
}