// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

// ConfigGen is used by `vault agent generate-config` only
type ConfigGen struct {
	AutoAuth     *AutoAuthGen      `hcl:"auto_auth,block"`
	Vault        *VaultGen         `hcl:"vault,block"`
	EnvTemplates []*EnvTemplateGen `hcl:"env_template,block"`
	Exec         *ExecConfig       `hcl:"exec,block"`
}

type EnvTemplateGen struct {
	Name              string  `hcl:"name,label"`
	Contents          string  `hcl:"contents,attr"`
	ErrorOnMissingKey *bool   `hcl:"error_on_missing_key,optional"`
	Group             *string `hcl:"group,optional"`
}

type VaultGen struct {
	Address string `hcl:"address"`
}

type AutoAuthGen struct {
	Method *AutoAuthMethodGen `hcl:"method,block"`
}

type AutoAuthMethodGen struct {
	Type   string                  `hcl:"type"`
	Config AutoAuthMethodConfigGen `hcl:"config,block"`
}

type AutoAuthMethodConfigGen struct {
	TokenFilePath string `hcl:"token_file_path"`
}
