package api

import (
	"fmt"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/vault"
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

func (c *Sys) Remount(from, to string, config *vault.MountConfig) error {
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
	if config != nil {
		body["config"] = *config
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

func (c *Sys) checkMountPath(path string) error {
	if path[0] == '/' {
		return fmt.Errorf("path must not start with /: %s", path)
	}

	return nil
}

type Mount struct {
	Type        string             `json:"type" structs:"type"`
	Description string             `json:"description" structs:"description"`
	Config      *vault.MountConfig `json:"config" structs:"config"`
}
