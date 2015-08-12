package api

import "fmt"

const SSHDefaultPath = "ssh"

// SSH is used to return a client to invoke operations on SSH backend.
type SSH struct {
	c    *Client
	Path string
}

// SSH is used to return the client for logical-backend API calls.
func (c *Client) SSH() *SSH {
	return c.SSHWithPath(SSHDefaultPath)
}

func (c *Client) SSHWithPath(path string) *SSH {
	return &SSH{
		c:    c,
		Path: path,
	}
}

// Invokes the SSH backend API to create a dynamic key or an OTP
func (c *SSH) Credential(role string, data map[string]interface{}) (*Secret, error) {
	r := c.c.NewRequest("PUT", fmt.Sprintf("/v1/%s/creds/%s", c.Path, role))
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
