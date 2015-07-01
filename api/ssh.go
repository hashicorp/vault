package api

import (
	"encoding/json"
	"fmt"
)

type SSH struct {
	c *Client
}

func (c *Client) SSH() *SSH {
	return &SSH{c: c}
}

func (c *SSH) KeyCreate(role string, data map[string]interface{}) (*Secret, error) {
	r := c.c.NewRequest("PUT", fmt.Sprintf("/v1/ssh/creds/%s", role))
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

func (c *SSH) Lookup(data map[string]interface{}) (*SSHRoles, error) {
	r := c.c.NewRequest("PUT", "/v1/ssh/lookup")
	if err := r.SetJSONBody(data); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var roles SSHRoles
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&roles); err != nil {
		return nil, err
	}
	return &roles, nil
}

type SSHRoles struct {
	Data map[string]interface{} `json:"data"`
}
