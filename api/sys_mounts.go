package api

import (
	"fmt"

	"github.com/fatih/structs"
)

func (c *Sys) ListMounts() (map[string]*MountOutput, error) {
	r := c.c.NewRequest("GET", "/v1/sys/mounts")
	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]*MountOutput
	err = resp.DecodeJSON(&result)
	return result, err
}

func (c *Sys) Mount(path string, mountInfo *MountInput) error {
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

func (c *Sys) TuneMount(path string, config MountConfigInput) error {
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

type MountInput struct {
	Type        string           `json:"type" structs:"type"`
	Description string           `json:"description" structs:"description"`
	Config      MountConfigInput `json:"config" structs:"config"`
}

type MountConfigInput struct {
	DefaultLeaseTTL string `json:"default_lease_ttl" structs:"default_lease_ttl" mapstructure:"default_lease_ttl"`
	MaxLeaseTTL     string `json:"max_lease_ttl" structs:"max_lease_ttl" mapstructure:"max_lease_ttl"`
}

type MountOutput struct {
	Type        string            `json:"type" structs:"type"`
	Description string            `json:"description" structs:"description"`
	Config      MountConfigOutput `json:"config" structs:"config"`
}

type MountConfigOutput struct {
	DefaultLeaseTTL int `json:"default_lease_ttl" structs:"default_lease_ttl" mapstructure:"default_lease_ttl"`
	MaxLeaseTTL     int `json:"max_lease_ttl" structs:"max_lease_ttl" mapstructure:"max_lease_ttl"`
}
