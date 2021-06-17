package nomad

type Allocations struct {
	c *Client
}

func (c *Client) Allocations() *Allocations {
	return &Allocations{c}
}

func (a *Allocations) Restart(allocationId string) error {
	return a.c.RestartAllocation(allocationId)
}

func (c *Client) RestartAllocation(allocationId string) error {
	req := c.client.Request()
	req.Path("/v1/client/allocation/" + allocationId + "/restart")
	req.SetHeader("Content-Type", "application/json")
	req.Method("POST")
	_, err := req.Send()
	if err != nil { return err }

	return nil
}