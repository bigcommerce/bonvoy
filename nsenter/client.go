package nsenter

import (
	"context"
	"github.com/Devatoria/go-nsenter"
)

type Client struct {
	config nsenter.Config
}

func NewClient(pid int) Client {
	return Client{
		config: nsenter.Config{
			Net: true,
			IPC: true,
			Target: pid,
		},
	}
}

func (c *Client) Exec(program string, args ...string) (string, error) {
	stdout, stderr, err := c.config.ExecuteContext(context.Background(), program, args...)
	if err != nil {
		return stderr, err
	}
	return stdout, nil
}

func (c *Client) Curl(args ...string) (string, error) {
	return c.Exec("curl", args...)
}