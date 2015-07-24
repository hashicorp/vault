package api

import "fmt"

// SSH is used to return a client to invoke operations on SSH backend.
type SSH struct {
	c *Client
}

// SSH is used to return the client for logical-backend API calls.
func (c *Client) SSH() *SSH {
	return &SSH{c: c}
}

// Invokes the SSH backend API to revoke a key identified by its lease ID.
func (c *SSH) KeyRevoke(id string) error {
	r := c.c.NewRequest("PUT", "/v1/sys/revoke/"+id)
	resp, err := c.c.RawRequest(r)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

// Invokes the SSH backend API to create a dynamic key or an OTP
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
