package envoy

import (
	"bonvoy/consul"
	"encoding/json"
	"strconv"
	"strings"
)

type Certificates struct {
	i *Instance
	endpoints CertificatesEndpoints
	consul consul.Client
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
		consul: consul.NewClient(),
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

type ExpiredCertificate struct {
	ServiceName string
	Pid int
	Envoy *Instance `json:"-"`
	EnvoyExpiration string
	EnvoyDaysUntilExpiration int
	ConsulExpiration string
}

func (c *Certificates) FindExpired() ([]ExpiredCertificate, error) {
	var expiredCerts []ExpiredCertificate

	data, err := c.Get()
	if err != nil { return expiredCerts, err }

	readSerials := map[string]bool{}

	for _, certs := range data.Certificates {
		for _, cert := range certs.CertificateChain {
			if readSerials[cert.SerialNumber] {
				continue // ignore duplicates
			}
			readSerials[cert.SerialNumber] = true
			a := strings.Split(cert.SubjectAltNames[0].Uri, "/")
			svc := a[len(a)-1]
			leaf, lErr := c.consul.Agent().GetConnectLeafCaCertificate(svc)
			if lErr != nil { return expiredCerts, lErr }

			if cert.ExpirationTime != leaf.ValidBefore {
				daysUntilExpiration, err := strconv.Atoi(cert.DaysUntilExpiration)
				if err != nil { daysUntilExpiration = -1 }

				expiredCerts = append(expiredCerts, ExpiredCertificate{
					ServiceName: svc,
					Envoy: c.i,
					Pid: c.i.Pid,
					EnvoyExpiration: cert.ExpirationTime,
					EnvoyDaysUntilExpiration: daysUntilExpiration,
					ConsulExpiration: leaf.ValidBefore,
				})
			}
		}
	}

	return expiredCerts, nil
}