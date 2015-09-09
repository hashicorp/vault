package api

import (
	"fmt"
	"time"

	"github.com/fatih/structs"
)

func (c *Sys) ListMounts() (map[string]*Mount, error) {
	r := c.c.NewRequest("GET", "/v1/sys/mounts")
	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]*Mount
	err = resp.DecodeJSON(&result)
	return result, err
}

func (c *Sys) Mount(path string, mountInfo *Mount) error {
	if err := c.checkMountPath(path); err != nil {
		return err
	}

	body := structs.Map(mountInfo)

	r := c.c.NewRequest("POST", fmt.Sprintf("/v1/sys/mounts/%s", path))
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

func (c *Sys) Unmount(path string) error {
	if err := c.checkMountPath(path); err != nil {
		return err
	}

	r := c.c.NewRequest("DELETE", fmt.Sprintf("/v1/sys/mounts/%s", path))
	resp, err := c.c.RawRequest(r)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

func (c *Sys) Remount(from, to string) error {
	if err := c.checkMountPath(from); err != nil {
		return err
	}
	if err := c.checkMountPath(to); err != nil {
		return err
	}

	body := map[string]interface{}{
		"from": from,
		"to":   to,
	}

	r := c.c.NewRequest("POST", "/v1/sys/remount")
	if err := r.SetJSONBody(body); err != nil {
		return err
	}

	resp, err := c.c.RawRequest(r)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

func (c *Sys) TuneMount(path string, config APIMountConfig) error {
	if err := c.checkMountPath(path); err != nil {
		return err
	}

	body := structs.Map(config)
	r := c.c.NewRequest("POST", fmt.Sprintf("/v1/sys/mounts/%s/tune", path))
	if err := r.SetJSONBody(body); err != nil {
		return err
	}

	resp, err := c.c.RawRequest(r)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

func (c *Sys) MountConfig(path string) error {
	if err := c.checkMountPath(path); err != nil {
		return err
	}

	r := c.c.NewRequest("GET", fmt.Sprintf("/v1/sys/mounts/%s/tune", path))

	resp, err := c.c.RawRequest(r)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

func (c *Sys) checkMountPath(path string) error {
	if path[0] == '/' {
		return fmt.Errorf("path must not start with /: %s", path)
	}

	return nil
}

type Mount struct {
	Type        string         `json:"type" structs:"type"`
	Description string         `json:"description" structs:"description"`
	Config      APIMountConfig `json:"config" structs:"config"`
}

type APIMountConfig struct {
	DefaultLeaseTTL *time.Duration `json:"default_lease_ttl" structs:"default_lease_ttl" mapstructure:"default_lease_ttl"`
	MaxLeaseTTL     *time.Duration `json:"max_lease_ttl" structs:"max_lease_ttl" mapstructure:"max_lease_ttl"`
}
