package api

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/vault/helper/jsonutil"
)

const (
	wrappedResponseLocation = "cubbyhole/response"
)

// Logical is used to perform logical backend operations on Vault.
type Logical struct {
	c *Client
}

// Logical is used to return the client for logical-backend API calls.
func (c *Client) Logical() *Logical {
	return &Logical{c: c}
}

func (c *Logical) Read(path string) (*Secret, error) {
	r := c.c.NewRequest("GET", "/v1/"+path)
	resp, err := c.c.RawRequest(r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp != nil && resp.StatusCode == 404 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return ParseSecret(resp.Body)
}

func (c *Logical) List(path string) (*Secret, error) {
	r := c.c.NewRequest("GET", "/v1/"+path)
	r.Params.Set("list", "true")
	resp, err := c.c.RawRequest(r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp != nil && resp.StatusCode == 404 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return ParseSecret(resp.Body)
}

func (c *Logical) Write(path string, data map[string]interface{}) (*Secret, error) {
	r := c.c.NewRequest("PUT", "/v1/"+path)
	if err := r.SetJSONBody(data); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		return ParseSecret(resp.Body)
	}

	return nil, nil
}

func (c *Logical) Delete(path string) (*Secret, error) {
	r := c.c.NewRequest("DELETE", "/v1/"+path)
	resp, err := c.c.RawRequest(r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		return ParseSecret(resp.Body)
	}

	return nil, nil
}

func (c *Logical) Unwrap(wrappingToken string) (*Secret, error) {
	origToken := c.c.Token()
	defer c.c.SetToken(origToken)

	c.c.SetToken(wrappingToken)

	secret, err := c.Read(wrappedResponseLocation)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %s", wrappedResponseLocation, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no value found at %s", wrappedResponseLocation)
	}
	if secret.Data == nil {
		return nil, fmt.Errorf("\"data\" not found in wrapping response")
	}
	if _, ok := secret.Data["response"]; !ok {
		return nil, fmt.Errorf("\"response\" not found in wrapping response \"data\" map")
	}

	wrappedSecret := new(Secret)
	buf := bytes.NewBufferString(secret.Data["response"].(string))
	if err := jsonutil.DecodeJSONFromReader(buf, wrappedSecret); err != nil {
		return nil, fmt.Errorf("error unmarshaling wrapped secret: %s", err)
	}

	return wrappedSecret, nil
}
