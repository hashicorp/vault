// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package plugincatalog

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
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
