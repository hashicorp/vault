package config

import "github.com/hashicorp/vault/logical/framework"

func newPasswordConfig(fieldData *framework.FieldData) *PasswordConf {
	return &PasswordConf{
		DefaultPasswordTTL: fieldData.Get("default_password_ttl").(int),
		MaxPasswordTTL:     fieldData.Get("max_password_ttl").(int),
		PasswordLength:     fieldData.Get("password_length").(int),
	}
}

type PasswordConf struct {
	DefaultPasswordTTL int `json:"default_password_ttl"`
	MaxPasswordTTL     int `json:"max_password_ttl"`
	PasswordLength     int `json:"password_length"`
}

func (c *PasswordConf) Map() map[string]interface{} {
	return map[string]interface{}{
		"default_password_ttl": c.DefaultPasswordTTL,
		"max_password_ttl":     c.MaxPasswordTTL,
		"password_length":      c.PasswordLength,
	}
}
