package api

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

func (c *Sys) ListAuth() (map[string]*AuthMount, error) {
	r := c.c.NewRequest("GET", "/v1/sys/auth")
	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	secret, err := ParseSecret(resp.Body)
	if err != nil {
		return nil, err
	}

	if secret == nil || secret.Data == nil || len(secret.Data) == 0 {
		return nil, nil
	}

	result := map[string]*AuthMount{}
	for k, v := range secret.Data {
		var res AuthMount
		err = mapstructure.Decode(v, &res)
		if err != nil {
			return nil, err
		}
		result[k] = &res
	}

	return result, err
}

func (c *Sys) EnableAuth(path, authType, desc string) error {
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
	r := c.c.NewRequest("DELETE", fmt.Sprintf("/v1/sys/auth/%s", path))
	resp, err := c.c.RawRequest(r)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

// Structures for the requests/resposne are all down here. They aren't
// individually documentd because the map almost directly to the raw HTTP API
// documentation. Please refer to that documentation for more details.

type AuthMount struct {
	Type        string           `json:"type" structs:"type" mapstructure:"type"`
	Description string           `json:"description" structs:"description" mapstructure:"description"`
	Config      AuthConfigOutput `json:"config" structs:"config" mapstructure:"config"`
}

type AuthConfigOutput struct {
	DefaultLeaseTTL int `json:"default_lease_ttl" structs:"default_lease_ttl" mapstructure:"default_lease_ttl"`
	MaxLeaseTTL     int `json:"max_lease_ttl" structs:"max_lease_ttl" mapstructure:"max_lease_ttl"`
}
