package nomad

import (
	"crypto/tls"
	"github.com/spf13/viper"
	"gopkg.in/h2non/gentleman.v2"
	gtls "gopkg.in/h2non/gentleman.v2/plugins/tls"
)

type Client struct {
	client *gentleman.Client
	address string
}

func GetDefaultAddress() string {
	return viper.GetString("NOMAD_ADDR")
}

func NewClient() Client {
	cli := gentleman.New()
	cli.Use(gtls.Config(&tls.Config{InsecureSkipVerify: true}))
	cli.URL(GetDefaultAddress())
	return Client{
		address: GetDefaultAddress(),
		client: cli,
	}
}