// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package config

import (
	"testing"

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
		want       *SyncPolicyConfig
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
			want: &SyncPolicyConfig{
				Name:                       "ce",
				DisallowedGroups:           []string{"enterprise"},
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
			want: &SyncPolicyConfig{
				Name:             "docs",
				DisallowedGroups: []string{"enterprise", "changelog"},
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
			want: &SyncPolicyConfig{
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
