// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package configutil

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/go-secure-stdlib/reloadutil"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	sdkResource "github.com/hashicorp/hcp-sdk-go/resource"
	"github.com/mitchellh/cli"
)

// HCPLinkConfig is the HCP Link configuration for the server.
type HCPLinkConfig struct {
	UnusedKeys UnusedKeyMap `hcl:",unusedKeyPositions"`

	ResourceIDRaw               string                `hcl:"resource_id"`
	Resource                    *sdkResource.Resource `hcl:"-"`
	EnableAPICapability         bool                  `hcl:"enable_api_capability"`
	EnablePassThroughCapability bool                  `hcl:"enable_passthrough_capability"`
	ClientID                    string                `hcl:"client_id"`
	ClientSecret                string                `hcl:"client_secret"`

	TLSDisable    bool        `hcl:"-"`
	TLSDisableRaw interface{} `hcl:"tls_disable"`
	TLSCertFile   string      `hcl:"tls_cert_file"`
	TLSKeyFile    string      `hcl:"tls_key_file"`
	TLSConfig     *tls.Config `hcl:"-"`
}

func parseCloud(result *SharedConfig, list *ast.ObjectList) error {
	if len(list.Items) > 1 {
		return fmt.Errorf("only one 'cloud' block is permitted")
	}

	// Get our one item
	item := list.Items[0]

	if result.HCPLinkConf == nil {
		result.HCPLinkConf = &HCPLinkConfig{}
	}

	if err := hcl.DecodeObject(&result.HCPLinkConf, item.Val); err != nil {
		return multierror.Prefix(err, "cloud:")
	}

	// let's check if the Client ID and Secret are set in the environment
	if envClientID := os.Getenv("HCP_CLIENT_ID"); envClientID != "" {
		result.HCPLinkConf.ClientID = envClientID
	}
	if envClientSecret := os.Getenv("HCP_CLIENT_SECRET"); envClientSecret != "" {
		result.HCPLinkConf.ClientSecret = envClientSecret
	}

	// three pieces are necessary if the cloud stanza is configured
	if result.HCPLinkConf.ResourceIDRaw == "" || result.HCPLinkConf.ClientID == "" || result.HCPLinkConf.ClientSecret == "" {
		return multierror.Prefix(fmt.Errorf("failed to find the required cloud stanza configurations. all resource ID, client ID and client secret are required"), "cloud:")
	}

	res, err := sdkResource.FromString(result.HCPLinkConf.ResourceIDRaw)
	if err != nil {
		return multierror.Prefix(fmt.Errorf("failed to parse resource_id for HCP Link"), "cloud:")
	}
	result.HCPLinkConf.Resource = &res

	// ENV var takes precedence over the config value
	if apiCapEnv := os.Getenv("HCP_LINK_ENABLE_API_CAPABILITY"); apiCapEnv != "" {
		result.HCPLinkConf.EnableAPICapability = true
	}

	if passthroughCapEnv := os.Getenv("HCP_LINK_ENABLE_PASSTHROUGH_CAPABILITY"); passthroughCapEnv != "" {
		result.HCPLinkConf.EnablePassThroughCapability = true
	}

	if result.HCPLinkConf.TLSDisableRaw != nil {
		if result.HCPLinkConf.TLSDisable, err = parseutil.ParseBool(result.HCPLinkConf.TLSDisableRaw); err != nil {
			return multierror.Prefix(fmt.Errorf("invalid value for tls_disable: %w", err), "cloud:")
		}

		result.HCPLinkConf.TLSDisableRaw = nil
	}

	if !result.HCPLinkConf.TLSDisable && (result.HCPLinkConf.TLSCertFile == "" || result.HCPLinkConf.TLSKeyFile == "") {
		return multierror.Prefix(fmt.Errorf("TLS is enabled but failed to find TLS cert and key file"), "cloud:")
	}

	return nil
}

func (h *HCPLinkConfig) ParseTLSConfig(ui cli.Ui) error {
	if !h.TLSDisable {
		// We try the key without a passphrase first and if we get an incorrect
		// passphrase response, try again after prompting for a passphrase
		cg := reloadutil.NewCertificateGetter(h.TLSCertFile, h.TLSKeyFile, "")
		if err := cg.Reload(); err != nil {
			if errwrap.Contains(err, x509.IncorrectPasswordError.Error()) {
				var passphrase string
				passphrase, err = ui.AskSecret(fmt.Sprintf("Enter passphrase for cloud TLS key file %s:", h.TLSKeyFile))

				if err == nil {
					cg = reloadutil.NewCertificateGetter(h.TLSCertFile, h.TLSKeyFile, passphrase)
					if err = cg.Reload(); err != nil {
						return fmt.Errorf("error loading cloud config TLS cert with provided passphrase: %w", err)
					}

					h.TLSConfig = &tls.Config{
						GetCertificate: cg.GetCertificate,

						// Prefer sensible defaults based on the defaults we use for "tcp" listener config
						MinVersion: tls.VersionTLS12,
						MaxVersion: tls.VersionTLS13,
						ClientAuth: tls.RequestClientCert,
					}

					return nil
				}
			}

			return fmt.Errorf("error loading cloud config TLS cert: %w", err)
		}
	}

	return nil
}
