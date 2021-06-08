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
		config: nsenter.Config{
			Net: true,
			IPC: true,
			Target: pid,
		},
	}
}

func (c *Client) Exec(program string, args ...string) string {
	stdout, stderr, err := c.config.ExecuteContext(context.Background(), program, args...)
	if err != nil {
		fmt.Println(stderr)
		panic(err)
	}
	return stdout
}

func (c *Client) Curl(args ...string) string {
	return c.Exec("curl", args...)
}