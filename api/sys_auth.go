package api

import (
	"fmt"
)

func (c *Sys) ListAuth() ([]*AuthResponse, error) {
	r := c.c.NewRequest("GET", "/sys/auth")
	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result []*AuthResponse
	err = resp.DecodeJSON(&result)
	return result, err
}

func (c *Sys) EnableAuth(id string, opts *AuthRequest) error {
	body := make(map[string]string)
	for k, v := range opts.Config {
		body[k] = v
	}
	body["type"] = opts.Type

	r := c.c.NewRequest("PUT", fmt.Sprintf("/sys/auth/%s", id))
	if err := r.SetJSONBody(body); err != nil {
		return err
	}

	resp, err := c.c.RawRequest(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *Sys) DisableAuth(id string) error {
	r := c.c.NewRequest("DELETE", fmt.Sprintf("/sys/auth/%s", id))
	resp, err := c.c.RawRequest(r)
	defer resp.Body.Close()
	return err
}

// Structures for the requests/resposne are all down here. They aren't
// individually documentd because the map almost directly to the raw HTTP API
// documentation. Please refer to that documentation for more details.

type AuthRequest struct {
	Type   string
	Config map[string]string
}

type AuthResponse struct {
	ID   string
	Type string
	Help string
	Keys []string
}
