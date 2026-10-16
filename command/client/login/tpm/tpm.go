// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package tpm

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	sdktpm "github.com/hashicorp/vault/sdk/helper/tpm"
)

type CLIHandler struct{}

// AuthWithFlags satisfies command.LoginHandlerWithFlags. It reads tpm-state-dir
// and tpm-device-path from the parsed CLI flags on the login command, falling
// back to the k=v map only for role_name and mount.
func (h *CLIHandler) AuthWithFlags(c *api.Client, m map[string]string, flags pluginutil.FlagLookup) (*api.Secret, error) {
	stateDir := flags.Lookup("tpm-state-dir").Value.String()
	devicePath := flags.Lookup("tpm-device-path").Value.String()

	roleName := m["role_name"]
	if roleName == "" {
		return nil, fmt.Errorf("role_name is required")
	}

	mount := m["mount"]
	if mount == "" {
		mount = "tpm"
	}

	result, err := sdktpm.LoadFromDirectory(stateDir)
	if err != nil {
		return nil, err
	}

	tpmCfg := sdktpm.TPMConfig{DevicePath: devicePath}
	tlsCfg, err := sdktpm.BuildTLSConfig(result, "", "", tpmCfg)
	if err != nil {
		return nil, err
	}

	loginConfig := api.DefaultConfig()
	loginConfig.Address = c.Address()
	existingTLS := c.CloneConfig().TLSConfig()
	if existingTLS != nil {
		loginConfig.HttpClient.Transport.(*http.Transport).TLSClientConfig = existingTLS
	}
	loginConfig.HttpClient.Transport.(*http.Transport).TLSClientConfig.GetClientCertificate = tlsCfg.GetClientCertificate

	loginClient, err := api.NewClient(loginConfig)
	if err != nil {
		return nil, err
	}

	options := map[string]interface{}{
		"role_name": roleName,
	}
	path := fmt.Sprintf("auth/%s/login", mount)
	secret, err := loginClient.Logical().Write(path, options)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, fmt.Errorf("empty response from credential provider")
	}
	return secret, nil
}

func (h *CLIHandler) Help() string {
	help := `
Usage: vault login -method=tpm [CONFIG K=V...]

  The TPM auth method allows machines to authenticate using a TPM-backed TLS certificate. 
  The certificate and associated state files must have been generated in advance by running 
  "vault tpm attest".

  Authenticate using a TPM-backed certificate:

      $ vault login -method=tpm -tpm-state-dir=/path/to/state role_name=my-role

  Specify a custom TPM device:

      $ vault login -method=tpm -tpm-state-dir=/path/to/state -tpm-device-path=/dev/tpm0 role_name=my-role

Configuration:

  role_name=<string>
	  Name of the TPM auth role to authenticate against. Required

  mount=<string>
      The mount path of the TPM auth method. Defaults to "tpm".
`

	return strings.TrimSpace(help)
}
