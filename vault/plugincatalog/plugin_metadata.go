// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package plugincatalog

import (
	"encoding/json"
	"os"

	"github.com/hashicorp/vault/sdk/helper/consts"
)

type pluginMetadata struct {
	Version string `json:"version"`
	Plugin  Plugin `json:"plugin"`
}

type Plugin struct {
	Name string            `json:"name"`
	Type consts.PluginType `json:"type"`
	Tier consts.PluginTier `json:"tier,omitempty"`
	// By is the plugin author's GitHub account, following Terraform Registry's convention
	By       string `json:"by"`
	Version  string `json:"version"`
	Platform string `json:"platform"`
	Arch     string `json:"arch"`
	// PGPSig is PGP ASCII armored detached signature
	PGPSig string `json:"pgp_sig"`
}

func readPluginMetadata(metadataPath string) (*pluginMetadata, error) {
	metadataBytes, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil, err
	}

	metadata := pluginMetadata{}
	if err = json.Unmarshal(metadataBytes, &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}
