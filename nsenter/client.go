package nsenter

import (
	"context"
	"fmt"
	"github.com/Devatoria/go-nsenter"
)

type Client struct {
	config nsenter.Config
}

func NewClient(pid int) Client {
	return Client{
		config: BuildConfig(pid),
	}
}

func BuildConfig(pid int) nsenter.Config {
	return nsenter.Config{
		Net: true,
		IPC: true,
		Target: pid,
	}
}

func (c *Client) Curl(args ...string) string {
	stdout, stderr, err := c.config.ExecuteContext(context.Background(), "curl", args...)
	if err != nil {
		fmt.Println(stderr)
		panic(err)
	}
	return stdout
}