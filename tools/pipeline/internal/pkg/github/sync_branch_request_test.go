// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"errors"
	"testing"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/changed"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/config"
	"github.com/stretchr/testify/require"
)

// TestSyncBranchReq_Validate verifies that a sync policy requires a decoded
// pipeline config the same way that disallowed groups do. We can't look up the
// policy without it, and we'd rather fail before we clone anything.
func TestSyncBranchReq_Validate(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		disallowedGroups changed.FileGroups
		syncPolicy       string
		decodeRes        *config.DecodeRes
		errHas           string
	}{
		"no changed file checks": {},
		"disallowed groups without config": {
			disallowedGroups: changed.FileGroups{"enterprise"},
			errHas:           "pipeline.hcl decode response is unitialized",
		},
		"sync policy without config": {
			syncPolicy: "ce",
			errHas:     "pipeline.hcl decode response is unitialized",
		},
		"sync policy with a decode error": {
			syncPolicy: "ce",
			decodeRes:  &config.DecodeRes{Err: errors.New("bad pipeline.hcl")},
			errHas:     "bad pipeline.hcl",
		},
		"sync policy with no decoded config": {
			syncPolicy: "ce",
			decodeRes:  &config.DecodeRes{},
			errHas:     "disallowed groups or a sync policy have been provided but no changed files config has been provided",
		},
		"sync policy with config": {
			syncPolicy: "ce",
			decodeRes:  &config.DecodeRes{Config: &config.Config{}},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := &SyncBranchReq{
				FromOwner:        "hashicorp",
				FromRepo:         "vault-enterprise",
				FromOrigin:       "from",
				FromBranch:       "main",
				ToOwner:          "hashicorp",
				ToRepo:           "vault",
				ToOrigin:         "to",
				ToBranch:         "main",
				DisallowedGroups: test.disallowedGroups,
				SyncPolicy:       test.syncPolicy,
				DecodeRes:        test.decodeRes,
			}

			err := req.Validate(t.Context())
			if test.errHas != "" {
				require.ErrorContains(t, err, test.errHas)

				return
			}
			require.NoError(t, err)
		})
	}
}

// TestSyncBranchRes_ToTable verifies that a sync without an error gets an empty
// error cell rather than <nil>, which lets us hide the error column.
func TestSyncBranchRes_ToTable(t *testing.T) {
	t.Parallel()

	res := &SyncBranchRes{Request: &SyncBranchReq{
		FromOwner:  "hashicorp",
		FromRepo:   "vault-enterprise",
		FromBranch: "ce/main",
		ToOwner:    "hashicorp",
		ToRepo:     "vault",
		ToBranch:   "main",
	}}

	for name, test := range map[string]struct {
		err          error
		wantContains []string
		wantMissing  []string
	}{
		"no error": {
			wantContains: []string{"hashicorp/vault-enterprise/ce/main", "hashicorp/vault/main"},
			wantMissing:  []string{"<nil>", "ERROR"},
		},
		"error": {
			err:          errors.New("merging failed"),
			wantContains: []string{"hashicorp/vault-enterprise/ce/main", "ERROR", "merging failed"},
			wantMissing:  []string{"<nil>"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := res.ToTable(test.err).Render()
			for _, want := range test.wantContains {
				require.Contains(t, got, want)
			}
			for _, missing := range test.wantMissing {
				require.NotContains(t, got, missing)
			}
		})
	}
}
