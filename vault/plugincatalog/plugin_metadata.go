// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package plugincatalog

import (
	"encoding/json"
	"os"

	"github.com/hashicorp/vault/sdk/helper/consts"
)

// PluginMetadata represents metadata.json in the plugin artifact
type PluginMetadata struct {
	Version string `json:"version"`
	Plugin  Plugin `json:"plugin"`
}

// Plugin represents the metadata of a plugin as defined
// under "plugin" object in the metadata.json
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
	// Sha256 is the SHA256 checksum of the plugin binary.
	// TODO: Not currently present in today's plugin release metadata, but will be added in the future.
	Sha256 string `json:"sha256"`
}

func readPluginMetadata(metadataPath string) (*PluginMetadata, error) {
	metadataBytes, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil, err
	}

	metadata := PluginMetadata{}
	if err = json.Unmarshal(metadataBytes, &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}
