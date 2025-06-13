// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package plugincatalog

import (
	"encoding/json"
	"fmt"
	"runtime"
	"testing"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/google/uuid"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/stretchr/testify/assert"
)

// gereratePGPKeyPair generates a PGP key pair for testing purposes
func generatePGPKeyPair(t *testing.T) (*crypto.Key, string) {
	pgp := crypto.PGP()
	user := "test" + uuid.NewString()[:5]
	privKey, err := pgp.KeyGeneration().AddUserId(user, fmt.Sprintf("%s@hashicorp.com", user)).New().GenerateKey()
	assert.NoError(t, err)

	pubkey, err := privKey.ToPublic()
	assert.NoError(t, err)

	armored, err := pubkey.Armor()
	assert.NoError(t, err)

	return privKey, armored
}

// generatePluginArtifactContents generates file contents for a plugin artifact for testing purposes
// If key is nil, signatures will carry a placeholder value
func generatePluginArtifactContents(t *testing.T, pluginName, pluginVersion string, pluginType consts.PluginType,
	includeBinarySig bool, privKey *crypto.Key,
) map[string][]byte {
	t.Helper()

	metadata := PluginMetadata{
		Version: "v0",
		Plugin: Plugin{
			Name:     pluginName,
			Type:     pluginType,
			Tier:     consts.PluginTierOfficial,
			By:       "hashicorp",
			Version:  pluginVersion,
			Platform: runtime.GOOS,
			Arch:     runtime.GOARCH,
			PGPSig:   "signature-placeholder",
		},
	}
	metadataBytes, err := json.Marshal(metadata)
	assert.NoError(t, err)

	pluginBytes := []byte("plugin-binary-placeholder")
	metadataSigBytes := []byte("signature-placeholder")
	pgp := crypto.PGP()
	if privKey != nil {
		signer, err := pgp.Sign().SigningKey(privKey).Detached().New()
		defer signer.ClearPrivateParams()
		assert.NoError(t, err)

		// exclude binary signature for bad plugin signature read test
		metadata.Plugin.PGPSig = ""
		if includeBinarySig {
			signature, err := signer.Sign(pluginBytes, crypto.Armor)
			assert.NoError(t, err)

			metadata.Plugin.PGPSig = string(signature)
		}

		metadataBytes, err = json.Marshal(metadata)
		assert.NoError(t, err)

		metadataSigBytes, err = signer.Sign(metadataBytes, crypto.Armor)
		assert.NoError(t, err)
	}

	return map[string][]byte{
		metadataFile: metadataBytes,
		metadataSig:  metadataSigBytes,
		pluginName:   pluginBytes,
		"LICENSE":    []byte("license-placeholder"),
	}
}
