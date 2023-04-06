// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

// ConfigGen is used by `vault agent generate-config` only
type ConfigGen struct {
	AutoAuth     *AutoAuthGen         `hcl:"auto_auth,block"`
	Vault        *VaultGen            `hcl:"vault,block"`
	EnvTemplates []*EnvTemplateConfig `hcl:"env_template,block"`
	Exec         *ExecConfig          `hcl:"exec,block"`
}

type VaultGen struct {
	Address string `hcl:"address"`
}

type AutoAuthGen struct {
	Method *AutoAuthMethodGen `hcl:"method,block"`
}

type AutoAuthMethodGen struct {
	Type   string
	Config AutoAuthMethodConfigGen `hcl:"config,block"`
}

type AutoAuthMethodConfigGen struct {
	TokenFileLocation string `hcl:"token_file_location"`
}
