package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
)

// ListAuth lists auth methods.
//
// Deprecated: ListAuth is deprecated. Use ListAuths instead.
func (c *Sys) ListAuth() (map[string]*AuthMount, error) {
	return c.ListAuthsWithContext(context.Background())
}

// ListAuths lists the enabled auth methods.
func (c *Sys) ListAuths() (map[string]*AuthMount, error) {
	return c.ListAuthsWithContext(context.Background())
}

// ListAuthsWithContext lists the enabled auth methods, with a context.
func (c *Sys) ListAuthsWithContext(ctx context.Context) (map[string]*AuthMount, error) {
	req := c.c.NewRequest(http.MethodGet, "/v1/sys/auth")
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

	mounts := map[string]*AuthMount{}
	for k, v := range result {
		switch v.(type) {
		case map[string]interface{}:
		default:
			continue
		}
		var res AuthMount
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

// EnableAuth enables an authentication at the given path, with the given type
// and description.
//
// Deprecated: EnableAuth is deprecated. Use EnableAuthWithOptions instead.
func (c *Sys) EnableAuth(path, authType, desc string) error {
	return c.EnableAuthWithOptions(path, &EnableAuthOptions{
		Type:        authType,
		Description: desc,
	})
}

// EnableAuthWithOptions enables an authentication at the given path with the
// given options.
func (c *Sys) EnableAuthWithOptions(path string, opts *EnableAuthOptions) error {
	return c.EnableAuthWithOptionsWithContext(context.Background(), path, opts)
}

// EnableAuthWithOptionsWithContext enables an authentication at the given path
// with the given options, with a context.
func (c *Sys) EnableAuthWithOptionsWithContext(ctx context.Context, path string, opts *EnableAuthOptions) error {
	req := c.c.NewRequest(http.MethodPost, fmt.Sprintf("/v1/sys/auth/%s", path))
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

// DisableAuth disables an auth at the given path.
func (c *Sys) DisableAuth(path string) error {
	return c.DisableAuthWithContext(context.Background(), path)
}

// DisableAuthWithContext disables an auth at the given path, with a context.
func (c *Sys) DisableAuthWithContext(ctx context.Context, path string) error {
	req := c.c.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/sys/auth/%s", path))
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

// EnableAuthOptions is used as input to the EnableAuthWithOptions function.
type EnableAuthOptions struct {
	// Type is the type of auth method to enable.
	Type string `json:"type" structs:"type"`

	// Description is a human-friendly description of the auth method.
	Description string `json:"description" structs:"description"`

	// Config is configuration about the auth method.
	Config AuthConfigInput `json:"config" structs:"config"`

	// Local is a boolean indicating the auth method should not be replicated.
	Local bool `json:"local" structs:"local"`

	// PluginName is the name of the plugin (if Type == "plugin").
	PluginName string `json:"plugin_name,omitempty" structs:"plugin_name,omitempty"`
}

// AuthConfigInput is the input auth config.
type AuthConfigInput struct {
	// PluginName is the name of the plugin (if Type == "plugin").
	PluginName string `json:"plugin_name,omitempty" structs:"plugin_name,omitempty" mapstructure:"plugin_name"`
}

// AuthMount is information about an auth mount.
type AuthMount struct {
	// Type is the type of auth method.
	Type string `json:"type" structs:"type" mapstructure:"type"`

	// Description is a human-friendly description of the auth method.
	Description string `json:"description" structs:"description" mapstructure:"description"`

	// Accessor is the accessor of the auth method.
	Accessor string `json:"accessor" structs:"accessor" mapstructure:"accessor"`

	// Config is configuration about the auth method.
	Config AuthConfigOutput `json:"config" structs:"config" mapstructure:"config"`

	// Local is a boolean indicating the auth method should not be replicated.
	Local bool `json:"local" structs:"local" mapstructure:"local"`
}

// AuthConfigOutput is output auth config.
type AuthConfigOutput struct {
	// DefaultLeaseTTL is the default minimum lease.
	DefaultLeaseTTL int `json:"default_lease_ttl" structs:"default_lease_ttl" mapstructure:"default_lease_ttl"`

	// MaxLeaseTTL is the maximum lease.
	MaxLeaseTTL int `json:"max_lease_ttl" structs:"max_lease_ttl" mapstructure:"max_lease_ttl"`

	// PluginName is the name of the plugin (if Type == "plugin")
	PluginName string `json:"plugin_name,omitempty" structs:"plugin_name,omitempty" mapstructure:"plugin_name"`
}
