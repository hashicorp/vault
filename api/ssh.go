package api

import "fmt"

type Ssh struct {
	c *Client
}

func (c *Client) Ssh() *Ssh {
	return &Ssh{c: c}
}

func (c *Ssh) KeyCreate(data map[string]interface{}) (*Secret, error) {
	r := c.c.NewRequest("PUT", fmt.Sprintf("/v1/ssh/creds/web"))
	if err := r.SetJSONBody(data); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseSecret(resp.Body)
}
