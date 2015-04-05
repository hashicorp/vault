package api

import (
	"fmt"
)

func (c *Sys) ListAuth() (map[string]*AuthMount, error) {
	r := c.c.NewRequest("GET", "/v1/sys/auth")
	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]*AuthMount
	err = resp.DecodeJSON(&result)
	return result, err
}

func (c *Sys) EnableAuth(path, authType, desc string) error {
	if err := c.checkAuthPath(path); err != nil {
		return err
	}

	body := map[string]string{
		"type":        authType,
		"description": desc,
	}

	r := c.c.NewRequest("POST", fmt.Sprintf("/v1/sys/auth/%s", path))
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

func (c *Sys) DisableAuth(path string) error {
	if err := c.checkAuthPath(path); err != nil {
		return err
	}

	r := c.c.NewRequest("DELETE", fmt.Sprintf("/v1/sys/auth/%s", path))
	resp, err := c.c.RawRequest(r)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

func (c *Sys) checkAuthPath(path string) error {
	if path[0] == '/' {
		return fmt.Errorf("path must not start with /: %s", path)
	}

	return nil
}

// Structures for the requests/resposne are all down here. They aren't
// individually documentd because the map almost directly to the raw HTTP API
// documentation. Please refer to that documentation for more details.

type AuthMount struct {
	Type        string
	Description string
}
