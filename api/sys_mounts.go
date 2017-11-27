package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
)

// ListMounts lists all available mounts.
func (c *Sys) ListMounts() (map[string]*MountOutput, error) {
	return c.ListMountsWithContext(context.Background())
}

// ListMountsWithContext lists all available mounts, with a context.
func (c *Sys) ListMountsWithContext(ctx context.Context) (map[string]*MountOutput, error) {
	req := c.c.NewRequest(http.MethodGet, "/v1/sys/mounts")
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

	mounts := map[string]*MountOutput{}
	for k, v := range result {
		switch v.(type) {
		case map[string]interface{}:
		default:
			continue
		}
		var res MountOutput
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

// Mount mounts a new backend at the given path.
func (c *Sys) Mount(path string, mountInfo *MountInput) error {
	return c.MountWithContext(context.Background(), path, mountInfo)
}

// MountWithContext mounts a new backend at the given path, with a context.
func (c *Sys) MountWithContext(ctx context.Context, path string, mountInfo *MountInput) error {
	req := c.c.NewRequest(http.MethodPost, fmt.Sprintf("/v1/sys/mounts/%s", path))
	req = req.WithContext(ctx)

	body := structs.Map(mountInfo)
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

// Unmount unmounts the backend at the given path, if it exits.
func (c *Sys) Unmount(path string) error {
	return c.UnmountWithContext(context.Background(), path)
}

// UnmountWithContext unmounts the backend at the given path, if it exits, with
// a context.
func (c *Sys) UnmountWithContext(ctx context.Context, path string) error {
	req := c.c.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/sys/mounts/%s", path))
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

// Remount moves the backend at the "from" to the "to" path.
func (c *Sys) Remount(from, to string) error {
	return c.RemountWithContext(context.Background(), from, to)
}

// RemountWithContext moves the backend at the "from" to the "to" path, with a
// context.
func (c *Sys) RemountWithContext(ctx context.Context, from, to string) error {
	req := c.c.NewRequest(http.MethodPost, "/v1/sys/remount")
	req = req.WithContext(ctx)

	body := map[string]interface{}{
		"from": from,
		"to":   to,
	}
	if err := req.SetJSONBody(body); err != nil {
		return err
	}

	resp, err := c.c.RawRequest(req)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

// TuneMount tunes the backend mounted at the given path.
func (c *Sys) TuneMount(path string, config MountConfigInput) error {
	return c.TuneMountWithContext(context.Background(), path, config)
}

// TuneMountWithContext tunes the backend mounted at the given path, with a
// context.
func (c *Sys) TuneMountWithContext(ctx context.Context, path string, config MountConfigInput) error {
	req := c.c.NewRequest(http.MethodPost, fmt.Sprintf("/v1/sys/mounts/%s/tune", path))
	req = req.WithContext(ctx)

	body := structs.Map(config)
	if err := req.SetJSONBody(body); err != nil {
		return err
	}

	resp, err := c.c.RawRequest(req)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

// MountConfig returns the configuration for the mount at the given path.
func (c *Sys) MountConfig(path string) (*MountConfigOutput, error) {
	return c.MountConfigWithContext(context.Background(), path)
}

// MountConfigWithContext returns the configuration for the mount at the given
// path, with a context.
func (c *Sys) MountConfigWithContext(ctx context.Context, path string) (*MountConfigOutput, error) {
	req := c.c.NewRequest(http.MethodGet, fmt.Sprintf("/v1/sys/mounts/%s/tune", path))
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result MountConfigOutput
	err = resp.DecodeJSON(&result)
	if err != nil {
		return nil, err
	}

	return &result, err
}

// MountInput is used as input to the mount function for mounting a new backend.
type MountInput struct {
	// Type is the type of backend to mount.
	Type string `json:"type" structs:"type"`

	// Description is an optional human-friendly description of the backend.
	Description string `json:"description" structs:"description"`

	// Config is configuration information such as the default leases and
	// replication behavior.
	Config MountConfigInput `json:"config" structs:"config"`

	// Local is a boolean representing that this mount is local-only (should not
	// replicate).
	Local bool `json:"local" structs:"local"`

	// PluginName is the name of the plugin to use for this mount (if type ==
	// "plugin").
	PluginName string `json:"plugin_name,omitempty" structs:"plugin_name"`

	// SealWrap is a boolean indicating whether to use a seal wrap.
	SealWrap bool `json:"seal_wrap" structs:"seal_wrap" mapstructure:"seal_wrap"`
}

// MountConfigInput is used as input about the mount configuration.
type MountConfigInput struct {
	// DefaultLeaseTTL is the default minimum lease.
	DefaultLeaseTTL string `json:"default_lease_ttl" structs:"default_lease_ttl" mapstructure:"default_lease_ttl"`

	// MaxLeaseTTL is the default maximum lease.
	MaxLeaseTTL string `json:"max_lease_ttl" structs:"max_lease_ttl" mapstructure:"max_lease_ttl"`

	// ForceNoCache forces the mount to not use caching.
	ForceNoCache bool `json:"force_no_cache" structs:"force_no_cache" mapstructure:"force_no_cache"`

	// PluginName is the name of the plugin to use for the mount (if type ==
	// "plugin").
	PluginName string `json:"plugin_name,omitempty" structs:"plugin_name,omitempty" mapstructure:"plugin_name"`
}

// MountOutput is specific output about a mount.
type MountOutput struct {
	// Type is the type of mount.
	Type string `json:"type" structs:"type"`

	// Description is the human-friendly description of the mount, if one was
	// given.
	Description string `json:"description" structs:"description"`

	// Accessor is the mount accessor.
	Accessor string `json:"accessor" structs:"accessor"`

	// Config is the configuration information for the mount.
	Config MountConfigOutput `json:"config" structs:"config"`

	// Local is a boolean representing whether this mount is replicated.
	Local bool `json:"local" structs:"local"`

	// SealWrap is a boolean indicating whether this mount has a seal wrap.
	SealWrap bool `json:"seal_wrap" structs:"seal_wrap" mapstructure:"seal_wrap"`
}

// MountConfigOutput is the mount configuration.
type MountConfigOutput struct {
	// DefaultLeaseTTL is the default minimum lease.
	DefaultLeaseTTL int `json:"default_lease_ttl" structs:"default_lease_ttl" mapstructure:"default_lease_ttl"`

	// MaxLeaseTTL is the default maximum lease.
	MaxLeaseTTL int `json:"max_lease_ttl" structs:"max_lease_ttl" mapstructure:"max_lease_ttl"`

	// ForceNoCache forces the mount to not use caching.
	ForceNoCache bool `json:"force_no_cache" structs:"force_no_cache" mapstructure:"force_no_cache"`

	// PluginName is the name of the plugin to use for the mount (if type ==
	// "plugin").
	PluginName string `json:"plugin_name,omitempty" structs:"plugin_name,omitempty" mapstructure:"plugin_name"`
}
