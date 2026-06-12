// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package plugincatalog

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			assert.Equal(t, tt.want, GetExtractedArtifactDir(tt.args.command, tt.args.version))
		})
	}
}

// TestPluginCatalog_load tests that we can successfully load the HashiCorp PGP public key into our global verifier.
func TestPluginCatalog_load(t *testing.T) {
	err := load()
	assert.NoError(t, err, "expected successful load of PGP public key")
}

// TestHashiCorpPGPPubKey2030_ExpirationWarning tests that the hashiCorpPGPPubKey2030
// key has more than 1 year remaining before expiration. This test will fail when
// the key has 1 year or less remaining, serving as an early warning to rotate the key.
func TestHashiCorpPGPPubKey2030_ExpirationWarning(t *testing.T) {
	t.Parallel()

	// Parse the HashiCorp 2030 PGP key
	key, err := crypto.NewKeyFromArmored(hashiCorpPGPPubKey2030)
	require.NoError(t, err, "failed to parse hashiCorpPGPPubKey2030")

	// Get the underlying OpenPGP entity to access expiration details
	entity := key.GetEntity()
	require.NotNil(t, entity, "key entity should not be nil")

	// The key expires on 2030-03-01 according to the comment in plugin_artifact.go
	// We need to check if we're within 1 year of that expiration date
	expectedExpiration := time.Date(2030, 3, 1, 0, 0, 0, 0, time.UTC)
	now := time.Now()
	oneYearFromNow := now.AddDate(1, 0, 0)

	oneYearFromNowUnix := oneYearFromNow.Unix()
	if key.IsExpired(oneYearFromNowUnix) {
		timeRemaining := expectedExpiration.Sub(now)
		t.Fatalf("hashiCorpPGPPubKey2030 will be expired in 1 year (current time remaining: %v, expires: %v). "+
			"This key needs to be rotated!",
			timeRemaining.Round(24*time.Hour), expectedExpiration.Format("2006-01-02"))
	}

	timeRemaining := expectedExpiration.Sub(now)
	t.Logf("hashiCorpPGPPubKey2030 has %v remaining before expiration (expires: %v)",
		timeRemaining.Round(24*time.Hour), expectedExpiration.Format("2006-01-02"))
}
