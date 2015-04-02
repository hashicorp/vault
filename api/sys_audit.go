package api

import (
	"fmt"
)

func (c *Sys) ListAudit() (map[string]*Audit, error) {
	r := c.c.NewRequest("GET", "/v1/sys/audit")
	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]*Audit
	err = resp.DecodeJSON(&result)
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
	Type        string
	Description string
	Options     map[string]string
}
