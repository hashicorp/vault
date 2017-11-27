package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
)

// AuditHash calculates the audit hash for the given input.
func (c *Sys) AuditHash(path string, input string) (string, error) {
	return c.AuditHashWithContext(context.Background(), path, input)
}

// AuditHashWithContext calculates the audit hash for the given input, with a
// context.
func (c *Sys) AuditHashWithContext(ctx context.Context, path string, input string) (string, error) {
	req := c.c.NewRequest(http.MethodPut, fmt.Sprintf("/v1/sys/audit-hash/%s", path))
	req = req.WithContext(ctx)

	body := map[string]interface{}{
		"input": input,
	}
	if err := req.SetJSONBody(body); err != nil {
		return "", err
	}

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result hashResp
	err = resp.DecodeJSON(&result)
	if err != nil {
		return "", err
	}

	return result.Hash, err
}

// ListAudit lists audit backends.
//
// Deprecated: ListAudit is deprecated. Use ListAudits instead.
func (c *Sys) ListAudit() (map[string]*Audit, error) {
	return c.ListAuditsWithContext(context.Background())
}

// ListAudits returns all enabled audit backends.
func (c *Sys) ListAudits() (map[string]*Audit, error) {
	return c.ListAuditsWithContext(context.Background())
}

// ListAuditsWithContext returns all enabled audit backends, with a context.
func (c *Sys) ListAuditsWithContext(ctx context.Context) (map[string]*Audit, error) {
	req := c.c.NewRequest(http.MethodGet, "/v1/sys/audit")
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = resp.DecodeJSON(&result)
	if err != nil {
		return nil, err
	}

	mounts := map[string]*Audit{}
	for k, v := range result {
		switch v.(type) {
		case map[string]interface{}:
		default:
			continue
		}
		var res Audit
		err = mapstructure.Decode(v, &res)
		if err != nil {
			return nil, err
		}
		// Not a mount, some other api.Secret data
		if res.Type == "" {
			continue
		}
		mounts[k] = &res
	}

	return mounts, nil
}

// EnableAudit enables an audit backend at the given path, with the given type
// and description.
//
// Deprecated: EnableAudit is deprecated. Use EnableAuditWithOptions instead.
func (c *Sys) EnableAudit(path string, auditType string, desc string, opts map[string]string) error {
	return c.EnableAuditWithOptions(path, &EnableAuditOptions{
		Type:        auditType,
		Description: desc,
		Options:     opts,
	})
}

// EnableAuditWithOptions enables the given audit backend with the provided
// options.
func (c *Sys) EnableAuditWithOptions(path string, opts *EnableAuditOptions) error {
	return c.EnableAuditWithOptionsWithContext(context.Background(), path, opts)
}

// EnableAuditWithOptionsWithContext enables the given audit backend with the
// provided options, with a context.
func (c *Sys) EnableAuditWithOptionsWithContext(ctx context.Context, path string, opts *EnableAuditOptions) error {
	req := c.c.NewRequest(http.MethodPut, fmt.Sprintf("/v1/sys/audit/%s", path))
	req = req.WithContext(ctx)

	body := structs.Map(opts)
	if err := req.SetJSONBody(body); err != nil {
		return err
	}

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// DisableAudit disables the audit at the given path.
func (c *Sys) DisableAudit(path string) error {
	return c.DisableAuditWithContext(context.Background(), path)
}

// DisableAuditWithContext disables the audit at the given path, with a context.
func (c *Sys) DisableAuditWithContext(ctx context.Context, path string) error {
	req := c.c.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/sys/audit/%s", path))
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

// EnableAuditOptions is used as input to the EnableAuditWithOptions function.
type EnableAuditOptions struct {
	// Type is the type of audit backend to enable.
	Type string `json:"type" structs:"type"`

	// Description is a human-friendly description of the audit backend.
	Description string `json:"description" structs:"description"`

	// Options is a list of options to pass directly to the audit backend.
	Options map[string]string `json:"options" structs:"options"`

	// Local is a boolean indicating the audit backend should not be replicated.
	Local bool `json:"local" structs:"local"`
}

// Audit is information about an audit backend.
type Audit struct {
	// Path is the path where the audit backend is mounted.
	Path string

	// Type is the type of audit backend.
	Type string

	// Description is the human-friendly description of the audit backend.
	Description string

	// Options are the list of configurations for the audit backend.
	Options map[string]string

	// Local is a boolean indicating the audit backend is not be replicated.
	Local bool
}

type hashResp struct {
	Hash string `json:"hash"`
}
