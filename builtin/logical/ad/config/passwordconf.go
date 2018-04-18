package config

import (
	"fmt"

	"github.com/go-errors/errors"
	"github.com/hashicorp/vault/builtin/logical/ad/util"
	"github.com/hashicorp/vault/logical/framework"
)

func newPasswordConfig(fieldData *framework.FieldData) (*PasswordConf, error) {

	ttl := fieldData.Get("ttl").(int)
	maxTTL := fieldData.Get("max_ttl").(int)
	length := fieldData.Get("password_length").(int)

	if ttl > maxTTL {
		return nil, errors.New("ttl must be smaller than or equal to max_ttl")
	}

	if ttl < 1 {
		return nil, errors.New("ttl must be positive")
	}

	if maxTTL < 1 {
		return nil, errors.New("max_ttl must be positive")
	}

	if length < util.MinimumPasswordLength {
		return nil, fmt.Errorf("minimum password length is %d for sufficient complexity to be secure, though Vault recommends a higher length", util.MinimumPasswordLength)
	}

	return &PasswordConf{
		TTL:    ttl,
		MaxTTL: maxTTL,
		Length: length,
	}, nil
}

type PasswordConf struct {
	TTL    int `json:"ttl"`
	MaxTTL int `json:"max_ttl"`
	Length int `json:"password_length"`
}

func (c *PasswordConf) Map() map[string]interface{} {
	return map[string]interface{}{
		"ttl":             c.TTL,
		"max_ttl":         c.MaxTTL,
		"password_length": c.Length,
	}
}
