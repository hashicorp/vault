// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package plugincatalog

import (
	"archive/zip"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/stretchr/testify/assert"
)

// TestPluginCatalog_copyExtractedPluginDirectory tests the copyExtractedPluginDirectory method.
func TestPluginCatalog_copyExtractedPluginDirectory(t *testing.T) {
	t.Parallel()

	c := testPluginCatalog(t)
	type args struct {
		pluginName    string
		pluginVersion string
		pluginType    consts.PluginType
		overwrite     bool
		create        bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success without overwriting",
			args: args{
				pluginName:    "vault-plugin-auth-example",
				pluginVersion: "0.1.1+ent",
				pluginType:    consts.PluginTypeCredential,
				create:        true,
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with overwriting",
			args: args{
				pluginName:    "vault-plugin-auth-example",
				pluginVersion: "0.1.1+ent",
				pluginType:    consts.PluginTypeCredential,
				create:        true,
				overwrite:     true,
			},
			wantErr: assert.NoError,
		},
		{
			name: "extracted plugin dir does not exist",
			args: args{
				pluginName:    "vault-plugin-auth-example",
				pluginVersion: "0.1.1+ent",
				pluginType:    consts.PluginTypeCredential,
			},
			wantErr: func(t assert.TestingT, err error, args ...interface{}) bool {
				assert.ErrorContains(t, err, "failed to copy plugin directory")
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()

			relativePluginDir := strings.Trim(zipName(tt.args.pluginName, tt.args.pluginVersion), ".zip")
			pluginDir := filepath.Join(tempDir, relativePluginDir)
			overwriteFile := filepath.Join(c.directory, relativePluginDir, "dummy.txt")

			if tt.args.create {
				err := os.MkdirAll(pluginDir, 0o755)
				assert.NoError(t, err, "expected successful create extracted plugin directory")

				if tt.args.overwrite {
					os.Create(overwriteFile)
				}

				contents := generatePluginArtifactContents(t, tt.args.pluginName,
					tt.args.pluginVersion, tt.args.pluginType, true, nil)

				// Write the files to the extracted plugin directory
				for name, content := range contents {
					err = os.WriteFile(filepath.Join(pluginDir, name), content, 0o644)
					assert.NoError(t, err, "expected successful zip file write")
				}
			}

			tt.wantErr(t, c.copyExtractedPluginDirectory(pluginDir))
			if tt.args.overwrite {
				assert.NoFileExists(t, overwriteFile, "expected file to be removed")
			}
		})
	}
}

// TestPluginCatalog_unpackArtifact tests the unpackArtifact method.
func TestPluginCatalog_unpackArtifact(t *testing.T) {
	// Cannot run in parallel due to global defaultPGPPubKey
	c := testPluginCatalog(t)
	type args struct {
		pluginName    string
		pluginVersion string
		pluginType    consts.PluginType
	}
	type wants struct {
		tier                    consts.PluginTier
		unpackedArtifactDirName string
		command                 string
		err                     assert.ErrorAssertionFunc
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name: "success",
			args: args{
				pluginName:    "vault-plugin-auth-example",
				pluginVersion: "0.1.1+ent",
				pluginType:    consts.PluginTypeCredential,
			},
			wants: wants{
				tier:                    consts.PluginTierOfficial,
				unpackedArtifactDirName: fmt.Sprintf("vault-plugin-auth-example_0.1.1+ent_%s_%s", runtime.GOOS, runtime.GOARCH),
				command: fmt.Sprintf("vault-plugin-auth-example_0.1.1+ent_%s_%s/vault-plugin-auth-example",
					runtime.GOOS, runtime.GOARCH),
				err: assert.NoError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			privKey, pubKey := generatePGPKeyPair(t)
			defaultPGPPubKey = pubKey

			zipPath := filepath.Join(c.directory, zipName(tt.args.pluginName, tt.args.pluginVersion))

			zipFile, err := os.Create(zipPath)
			assert.NoError(t, err, "expected successful create zip file", zipPath)

			// Create a new ZIP writer
			zipWriter := zip.NewWriter(zipFile)

			// File contents
			contents := generatePluginArtifactContents(t, tt.args.pluginName,
				tt.args.pluginVersion, tt.args.pluginType, true, privKey)

			// Add files to the ZIP
			for name, content := range contents {
				writeFileToZip(t, zipWriter, name, content)
			}

			zipWriter.Close()
			zipFile.Close()

			pluginTier, unpackedArtifactDirName, command, sha256sum, err := c.unpackArtifact(pluginutil.SetPluginInput{
				Name:    tt.args.pluginName,
				Type:    tt.args.pluginType,
				Version: tt.args.pluginVersion,
				Command: tt.args.pluginName,
			})
			if !tt.wants.err(t, err) {
				return
			}

			if err != nil {
				return
			}

			assert.Equalf(t, tt.wants.tier, pluginTier, "tier")
			assert.Equalf(t, tt.wants.unpackedArtifactDirName, unpackedArtifactDirName, "unpacked artifact dir name")
			assert.Equalf(t, tt.wants.command, command, "command")

			hash := sha256.New()
			_, err = hash.Write(contents[tt.args.pluginName])
			assert.NoError(t, err, "expected successful write to hash")
			assert.Equalf(t, hash.Sum(nil), sha256sum, "hash")
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

// Test_unzip tests the unzip function.
func Test_unzip(t *testing.T) {
	t.Parallel()

	pluginCatalog := testPluginCatalog(t)

	type args struct {
		pluginName             string
		pluginVersion          string
		extractFilePathInvalid bool
		srcExists              bool
		dstExists              bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				pluginName:    "vault-plugin-auth-example",
				pluginVersion: "0.1.1+ent",
				srcExists:     true,
				dstExists:     true,
			},
			wantErr: assert.NoError,
		},
		{
			name: "src does not exist",
			args: args{
				pluginName:    "vault-plugin-auth-example",
				pluginVersion: "0.1.2+ent",
				srcExists:     false,
				dstExists:     true,
			},
			wantErr: func(t assert.TestingT, err error, args ...interface{}) bool {
				assert.ErrorContains(t, err, "no such file or directory")
				return true
			},
		},
		{
			name: "dst does not exist",
			args: args{
				pluginName:    "vault-plugin-auth-example",
				pluginVersion: "0.1.3+ent",
				srcExists:     true,
				dstExists:     false,
			},
			wantErr: func(t assert.TestingT, err error, args ...interface{}) bool {
				assert.True(t, errors.Is(err, errPluginArtifactUnzipDstNotFound))
				return true
			},
		},
		{
			name: "src and dst do not exist",
			args: args{
				pluginName:    "vault-plugin-auth-example",
				pluginVersion: "0.1.4+ent",
				srcExists:     false,
				dstExists:     false,
			},
			wantErr: assert.Error,
		},
		{
			name: "invalid extract file path",
			args: args{
				pluginName:             "vault-plugin-auth-example",
				pluginVersion:          "0.1.5+ent",
				extractFilePathInvalid: true,
				srcExists:              true,
				dstExists:              true,
			},
			wantErr: func(t assert.TestingT, err error, args ...interface{}) bool {
				assert.ErrorIs(t, err, errInvalidExtractPath)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zipPath := filepath.Join(pluginCatalog.directory, zipName(tt.args.pluginName, tt.args.pluginVersion))
			if tt.args.srcExists {
				zipFile, err := os.Create(zipPath)
				assert.NoError(t, err, "expected successful create zip file", zipPath)

				zipWriter := zip.NewWriter(zipFile)

				// File contents
				contents := generatePluginArtifactContents(t, tt.args.pluginName,
					tt.args.pluginVersion, consts.PluginTypeCredential, true, nil)

				// Add files to the ZIP
				for name, content := range contents {
					if tt.args.extractFilePathInvalid {
						name = filepath.Join("../../../..", name)
					}
					writeFileToZip(t, zipWriter, name, content)
				}

				zipWriter.Close()
				zipFile.Close()
			}

			var (
				tempDst string
				err     error
			)
			if tt.args.dstExists {
				tempDst, err = os.MkdirTemp(pluginCatalog.directory, "unzip-")
				assert.NoError(t, err, "expected os.MkdirTemp for unzip destination to succeed")
			}

			got, err := unzip(tt.args.pluginName, zipPath, tempDst)
			if !tt.wantErr(t, err,
				fmt.Sprintf("unzip(%v, %v, %v)", tt.args.pluginName, zipPath, tempDst)) {
				return
			}

			if err != nil {
				return
			}

			expected := filepath.Join(tempDst,
				strings.Trim(zipName(tt.args.pluginName, tt.args.pluginVersion), ".zip"))
			assert.Equal(t, expected, got)

			assert.DirExists(t, expected)
			assert.FileExists(t, filepath.Join(expected, tt.args.pluginName))
			assert.FileExists(t, filepath.Join(expected, metadataFile))
			assert.FileExists(t, filepath.Join(expected, metadataSig))
		})
	}
}

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			privKey, pubKeyArmored := generatePGPKeyPair(t)

			contents := generatePluginArtifactContents(t, tt.args.pluginName,
				tt.args.pluginVersion, tt.args.pluginType, !errors.Is(tt.expectedErr, errReadPluginPGPSig), privKey)

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
			default:
				t.Fatalf("unexpected error: %v", tt.expectedErr)
			}

			pluginDir := filepath.Join(tempDir,
				strings.Trim(zipName(tt.args.pluginName, tt.args.pluginVersion), ".zip"))
			err := os.MkdirAll(pluginDir, 0o755)
			assert.NoError(t, err, "expected successful create extracted plugin directory")

			// Write the files to the extracted plugin directory
			for name, content := range contents {
				err = os.WriteFile(filepath.Join(pluginDir, name), content, 0o644)
				assert.NoError(t, err, "expected successful zip file write")
			}

			err = verifyPlugin(pluginDir, tt.args.pluginName, pubKeyArmored)
			assert.ErrorIs(t, err, tt.expectedErr, "expected verify plugin error to match")
		})
	}
}

// Test_zipName tests the zipName function.
func Test_zipName(t *testing.T) {
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
			want: fmt.Sprintf("vault-plugin-auth-aws_0.18.0+ent_%s_%s.zip", runtime.GOOS, runtime.GOARCH),
		},
		{
			name: "un-prefixed version",
			args: args{"vault-plugin-auth-aws", "0.18.0+ent"},
			want: fmt.Sprintf("vault-plugin-auth-aws_0.18.0+ent_%s_%s.zip", runtime.GOOS, runtime.GOARCH),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, zipName(tt.args.command, tt.args.version))
		})
	}
}
