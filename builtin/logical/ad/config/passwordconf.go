package config

import "github.com/hashicorp/vault/logical/framework"

func newPasswordConfig(fieldData *framework.FieldData) *PasswordConf {
	return &PasswordConf{
		TTL:    fieldData.Get("ttl").(int),
		MaxTTL: fieldData.Get("max_ttl").(int),
		Length: fieldData.Get("password_length").(int),
	}
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
