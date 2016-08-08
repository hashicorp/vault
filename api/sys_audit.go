package api

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

func (c *Sys) AuditHash(path string, input string) (string, error) {
	body := map[string]interface{}{
		"input": input,
	}

	r := c.c.NewRequest("PUT", fmt.Sprintf("/v1/sys/audit-hash/%s", path))
	if err := r.SetJSONBody(body); err != nil {
		return "", err
	}

	resp, err := c.c.RawRequest(r)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	secret, err := ParseSecret(resp.Body)
	if err != nil {
		return "", err
	}

	if secret == nil || secret.Data == nil || len(secret.Data) == 0 {
		return "", nil
	}

	type d struct {
		Hash string
	}

	var result d
	err = mapstructure.Decode(secret.Data, &result)
	if err != nil {
		return "", err
	}

	return result.Hash, err
}

func (c *Sys) ListAudit() (map[string]*Audit, error) {
	r := c.c.NewRequest("GET", "/v1/sys/audit")
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

	result := map[string]*Audit{}
	for k, v := range secret.Data {
		var res Audit
		err = mapstructure.Decode(v, &res)
		if err != nil {
			return nil, err
		}
		result[k] = &res
	}

	return result, err
}

func (c *Sys) EnableAudit(
	path string, auditType string, desc string, opts map[string]string) error {
	body := map[string]interface{}{
		"type":        auditType,
		"description": desc,
		"options":     opts,
	}

	r := c.c.NewRequest("PUT", fmt.Sprintf("/v1/sys/audit/%s", path))
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

func (c *Sys) DisableAudit(path string) error {
	r := c.c.NewRequest("DELETE", fmt.Sprintf("/v1/sys/audit/%s", path))
	resp, err := c.c.RawRequest(r)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

// Structures for the requests/resposne are all down here. They aren't
// individually documentd because the map almost directly to the raw HTTP API
// documentation. Please refer to that documentation for more details.

type Audit struct {
	Path        string
	Type        string
	Description string
	Options     map[string]string
}
