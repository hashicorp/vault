package api

import (
	"encoding/json"
	"fmt"
)

type Ssh struct {
	c *Client
}

func (c *Client) Ssh() *Ssh {
	return &Ssh{c: c}
}

func (c *Ssh) KeyCreate(role string, data map[string]interface{}) (*Secret, error) {
	r := c.c.NewRequest("PUT", fmt.Sprintf("/v1/ssh/creds/"+role))
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

func (c *Ssh) Lookup(data map[string]interface{}) (*SshRoles, error) {
	r := c.c.NewRequest("PUT", "/v1/ssh/lookup")
	if err := r.SetJSONBody(data); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var roles SshRoles
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&roles); err != nil {
		return nil, err
	}
	return &roles, nil
}

type SshRoles struct {
	Data map[string]interface{} `json:"data"`
}
