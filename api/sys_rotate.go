package api

import (
	"time"

	"github.com/mitchellh/mapstructure"
)

func (c *Sys) Rotate() error {
	r := c.c.NewRequest("POST", "/v1/sys/rotate")
	resp, err := c.c.RawRequest(r)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

func (c *Sys) KeyStatus() (*KeyStatus, error) {
	r := c.c.NewRequest("GET", "/v1/sys/key-status")
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

	var result KeyStatus
	err = mapstructure.Decode(secret.Data, &result)
	if err != nil {
		return nil, err
	}

	return &result, err
}

type KeyStatus struct {
	Term        int
	InstallTime time.Time `json:"install_time"`
}
