// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package config

import (
	"testing"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/changed"
	"github.com/stretchr/testify/require"
)

// TestConfig_SyncPolicy verifies that we can decode sync policies and look them
// up by name. Asking for a policy that doesn't exist has to be an error so
// that a typo in --sync-policy can't skip the check.
func TestConfig_SyncPolicy(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		hcl        string
		policy     string
		want       *SyncPolicy
		errHas     string
		wantNoSync bool
	}{
		"one policy": {
			hcl: `
sync {
  policy "ce" {
    disallowed_groups             = ["enterprise"]
    exempt_from_disallowed_groups = [joinpath(".github", "copilot-instructions.md")]
  }
}
`,
			policy: "ce",
			want: &SyncPolicy{
				Name:                       "ce",
				DisallowedGroups:           changed.FileGroups{"enterprise"},
				ExemptFromDisallowedGroups: []string{".github/copilot-instructions.md"},
			},
		},
		"many policies": {
			hcl: `
sync {
  policy "ce" {
    disallowed_groups = ["enterprise"]
  }

  policy "docs" {
    disallowed_groups = ["enterprise", "changelog"]
  }
}
`,
			policy: "docs",
			want: &SyncPolicy{
				Name:             "docs",
				DisallowedGroups: changed.FileGroups{"enterprise", "changelog"},
			},
		},
		"policy with only exempt files": {
			hcl: `
sync {
  policy "ce" {
    exempt_from_disallowed_groups = ["AGENTS.md"]
  }
}
`,
			policy: "ce",
			want: &SyncPolicy{
				Name:                       "ce",
				ExemptFromDisallowedGroups: []string{"AGENTS.md"},
			},
		},
		"unknown policy": {
			hcl: `
sync {
  policy "ce" {
    disallowed_groups = ["enterprise"]
  }
}
`,
			policy: "community",
			errHas: "no sync policy named community",
		},
		"no sync block": {
			hcl: `
changed_files {
  group "enterprise" {
    match {
      extension = [".go"]
    }
  }
}
`,
			policy:     "ce",
			errHas:     "no sync config was found",
			wantNoSync: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cfg, err := DecodeBytes([]byte(test.hcl))
			require.NoError(t, err)
			if test.wantNoSync {
				require.Nil(t, cfg.Sync)
			}

			policy, err := cfg.SyncPolicy(test.policy)
			if test.errHas != "" {
				require.ErrorContains(t, err, test.errHas)
				require.Nil(t, policy)

				return
			}
			require.NoError(t, err)
			require.Equal(t, test.want, policy)
		})
	}
}

// TestConfig_DisallowedGroups tests that we combine the groups we're given with
// those of the sync policies, and that we fail if any of them aren't changed
// file groups since a typo would never match anything and quietly pass.
func TestConfig_DisallowedGroups(t *testing.T) {
	t.Parallel()

	const pipelineHCL = `
changed_files {
  group "changelog" {
    match {
      base_dir = ["changelog"]
    }
  }

  group "enterprise" {
    match {
      base_dir = ["vault_ent"]
    }
  }
}
`

	ce := &SyncPolicy{
		Name:                       "ce",
		DisallowedGroups:           changed.FileGroups{"enterprise"},
		ExemptFromDisallowedGroups: []string{"AGENTS.md"},
	}
	exemptOnly := &SyncPolicy{
		Name:                       "exempt-only",
		ExemptFromDisallowedGroups: []string{"AGENTS.md"},
	}

	for name, test := range map[string]struct {
		hcl      string
		groups   changed.FileGroups
		policies []*SyncPolicy
		want     changed.FileGroups
		errHas   string
	}{
		"groups only": {
			groups: changed.FileGroups{"enterprise", "changelog"},
			want:   changed.FileGroups{"changelog", "enterprise"},
		},
		"sync policy only": {
			policies: []*SyncPolicy{ce},
			want:     changed.FileGroups{"enterprise"},
		},
		"groups are added to the sync policy groups": {
			groups:   changed.FileGroups{"changelog", "enterprise"},
			policies: []*SyncPolicy{ce},
			want:     changed.FileGroups{"changelog", "enterprise"},
		},
		"many sync policies": {
			policies: []*SyncPolicy{
				ce,
				{Name: "changelog", DisallowedGroups: changed.FileGroups{"changelog", "enterprise"}},
			},
			want: changed.FileGroups{"changelog", "enterprise"},
		},
		"nil sync policies are ignored": {
			groups:   changed.FileGroups{"enterprise"},
			policies: []*SyncPolicy{nil},
			want:     changed.FileGroups{"enterprise"},
		},
		"sync policy with only exempt files and groups": {
			groups:   changed.FileGroups{"enterprise"},
			policies: []*SyncPolicy{exemptOnly},
			want:     changed.FileGroups{"enterprise"},
		},
		"sync policy with only exempt files": {
			policies: []*SyncPolicy{exemptOnly},
			errHas:   "no disallowed groups have been configured",
		},
		"nothing": {
			errHas: "no disallowed groups have been configured",
		},
		"unknown group": {
			groups: changed.FileGroups{"enterprse"},
			errHas: "disallowed group enterprse is not a changed file group in pipeline.hcl",
		},
		"unknown group in the sync policy": {
			policies: []*SyncPolicy{{Name: "typo", DisallowedGroups: changed.FileGroups{"enterprse"}}},
			errHas:   "disallowed group enterprse is not a changed file group in pipeline.hcl",
		},
		"no changed files config": {
			hcl:    `sync {}`,
			groups: changed.FileGroups{"enterprise"},
			errHas: "no changed file grouping config was found in pipeline.hcl",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			hcl := pipelineHCL
			if test.hcl != "" {
				hcl = test.hcl
			}
			cfg, err := DecodeBytes([]byte(hcl))
			require.NoError(t, err)

			groups, err := cfg.DisallowedGroups(
				WithDisallowedGroupsGroups(test.groups),
				WithDisallowedGroupsSyncPolicies(test.policies...),
			)
			if test.errHas != "" {
				require.ErrorContains(t, err, test.errHas)
				require.Nil(t, groups)

				return
			}
			require.NoError(t, err)
			require.Equal(t, test.want, groups)
		})
	}
}

// TestDecodeBytes_InvalidSync verifies that we reject sync config that would
// make --sync-policy ambiguous or that has attributes we don't support. HCL
// will decode more than one policy with the same name so we have to catch
// that ourselves.
func TestDecodeBytes_InvalidSync(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		hcl    string
		errHas string
	}{
		"duplicate policy names": {
			hcl: `
sync {
  policy "ce" {
    disallowed_groups = ["enterprise"]
  }

  policy "ce" {
    disallowed_groups = ["changelog"]
  }
}
`,
			errHas: "sync policy ce has been defined more than once",
		},
		"more than one sync block": {
			hcl: `
sync {
  policy "ce" {}
}

sync {
  policy "docs" {}
}
`,
			errHas: "Duplicate sync block",
		},
		"policy without a name": {
			hcl: `
sync {
  policy {
    disallowed_groups = ["enterprise"]
  }
}
`,
			errHas: "Missing name for policy",
		},
		"unsupported attribute": {
			hcl: `
sync {
  policy "ce" {
    allow = ["AGENTS.md"]
  }
}
`,
			errHas: "Unsupported argument",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cfg, err := DecodeBytes([]byte(test.hcl))
			require.ErrorContains(t, err, test.errHas)
			require.Nil(t, cfg)
		})
	}
}
