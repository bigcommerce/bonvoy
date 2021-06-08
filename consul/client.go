package consul

import (
	"github.com/spf13/viper"
	"net/http"
)

type Client struct {
	cli *http.Client
	address string
}

func NewClient() Client {
	return Client{
		address: viper.GetString("CONSUL_API_HOST"),
	}
}