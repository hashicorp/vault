// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package configutil

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	sdkResource "github.com/hashicorp/hcp-sdk-go/resource"
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

	return nil
}
