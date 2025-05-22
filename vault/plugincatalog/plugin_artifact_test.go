// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package plugincatalog

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/stretchr/testify/assert"
)

// Test_verifyPlugin tests the verifyPlugin function.
func Test_verifyPlugin(t *testing.T) {
	t.Parallel()

	type args struct {
		pluginName    string
		pluginVersion string
		pluginType    consts.PluginType
	}
	tests := []struct {
		name        string
		args        args
		expectedErr error
	}{
		{
			name: "success",
			args: args{
				pluginName:    "vault-plugin-auth-example",
				pluginVersion: "0.1.1+ent",
				pluginType:    consts.PluginTypeCredential,
			},
			expectedErr: nil,
		},
		{
			name: "missing metadata",
			args: args{
				pluginName:    "vault-plugin-auth-example",
				pluginVersion: "0.1.2+ent",
				pluginType:    consts.PluginTypeCredential,
			},
			expectedErr: errReadMetadata,
		},
		{
			name: "missing metadata signature",
			args: args{
				pluginName:    "vault-plugin-auth-example",
				pluginVersion: "0.1.3+ent",
				pluginType:    consts.PluginTypeCredential,
			},
			expectedErr: errReadMetadataSig,
		},
		{
			name: "bad metadata signature verify",
			args: args{
				pluginName:    "vault-plugin-secret-example",
				pluginVersion: "0.1.4+ent",
				pluginType:    consts.PluginTypeSecrets,
			},
			expectedErr: errVerifyMetadataSig,
		},
		{
			name: "missing plugin binary",
			args: args{
				pluginName:    "vault-plugin-database-example",
				pluginVersion: "0.1.5+ent",
				pluginType:    consts.PluginTypeDatabase,
			},
			expectedErr: errReadPlugin,
		},
		{
			name: "bad plugin binary signature verify",
			args: args{
				pluginName:    "vault-plugin-database-example",
				pluginVersion: "0.1.6+ent",
				pluginType:    consts.PluginTypeDatabase,
			},
			expectedErr: errVerifyPluginSig,
		},
		{
			name: "bad extracted artifact directory",
			args: args{
				pluginName:    "vault-plugin-database-example",
				pluginVersion: "0.1.6+ent",
				pluginType:    consts.PluginTypeDatabase,
			},
			expectedErr: errExtractedArtifactDirNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			privKey, pubKeyArmored := generatePGPKeyPair(t)

			contents := generatePluginArtifactContents(t, tt.args.pluginName,
				tt.args.pluginVersion, tt.args.pluginType, !errors.Is(tt.expectedErr, errReadPluginPGPSig), privKey)

			actualExtractedArtifactDir := getExtractedArtifactDir(tt.args.pluginName, tt.args.pluginVersion)
			switch {
			case tt.expectedErr == nil:
			case errors.Is(tt.expectedErr, errReadPluginPGPSig):
				// no-op
			case errors.Is(tt.expectedErr, errReadMetadata):
				delete(contents, metadataFile)
			case errors.Is(tt.expectedErr, errReadMetadataSig):
				delete(contents, metadataSig)
			case errors.Is(tt.expectedErr, errVerifyMetadataSig):
				contents[metadataFile] = []byte(`{"will not" : "match signature"}`)
			case errors.Is(tt.expectedErr, errReadPlugin):
				delete(contents, tt.args.pluginName)
			case errors.Is(tt.expectedErr, errVerifyPluginSig):
				contents[tt.args.pluginName] = []byte("will not match signature")
			case errors.Is(tt.expectedErr, errExtractedArtifactDirNotFound):
				actualExtractedArtifactDir += "not_found"
			default:
				t.Fatalf("unexpected error: %v", tt.expectedErr)
			}

			tempDir := t.TempDir()
			actualPluginDir := filepath.Join(tempDir, actualExtractedArtifactDir)
			err := os.MkdirAll(actualPluginDir, 0o755)
			assert.NoError(t, err, "expected successful create extracted plugin directory")

			// Write the files to the extracted plugin directory
			for name, content := range contents {
				err = os.WriteFile(filepath.Join(actualPluginDir, name), content, 0o644)
				assert.NoError(t, err, "expected successful file write")
			}
			var metadata *pluginMetadata
			metadata, err = verifyPlugin(path.Join(tempDir, getExtractedArtifactDir(tt.args.pluginName, tt.args.pluginVersion)),
				tt.args.pluginName, pubKeyArmored)
			assert.ErrorIs(t, err, tt.expectedErr, "expected verify plugin error to match")

			if tt.expectedErr == nil {
				assert.NotNil(t, metadata)
				assert.Equal(t, tt.args.pluginName, metadata.Plugin.Name)
				assert.Equal(t, tt.args.pluginVersion, metadata.Plugin.Version)
			}
		})
	}
}

// Test_getExtractedArtifactDir tests the getExtractedArtifactDir function.
func Test_getExtractedArtifactDir(t *testing.T) {
	t.Parallel()

	type args struct {
		command string
		version string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "v-prefixed version",
			args: args{"vault-plugin-auth-aws", "v0.18.0+ent"},
			want: fmt.Sprintf("vault-plugin-auth-aws_0.18.0+ent_%s_%s", runtime.GOOS, runtime.GOARCH),
		},
		{
			name: "un-prefixed version",
			args: args{"vault-plugin-auth-aws", "0.18.0+ent"},
			want: fmt.Sprintf("vault-plugin-auth-aws_0.18.0+ent_%s_%s", runtime.GOOS, runtime.GOARCH),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, getExtractedArtifactDir(tt.args.command, tt.args.version))
		})
	}
}

// TestPluginCatalog_hashiCorpPubPGPKey tests hashiCorpPubPGPKey read
// and verification key creation.
func TestPluginCatalog_hashiCorpPubPGPKey(t *testing.T) {
	pgp := crypto.PGP()
	key, err := crypto.NewKeyFromArmored(hashiCorpPGPPubKey)
	assert.NoError(t, err)

	_, err = pgp.Verify().VerificationKey(key).New()
	assert.NoError(t, err)
}
