// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package config

import "fmt"

// SyncConfig holds our policies for synchronizing branches.
type SyncConfig struct {
	Policies []*SyncPolicyConfig `hcl:"policy,block" json:"policies,omitempty"`
}

// SyncPolicyConfig is a named sync policy. Commands that support it select it
// with --sync-policy <name>.
type SyncPolicyConfig struct {
	Name string `hcl:"name,label" json:"name"`
	// DisallowedGroups are changed file groups that must not exist in the
	// branch. We add them to any groups passed with --disallowed-groups.
	DisallowedGroups []string `hcl:"disallowed_groups,optional" json:"disallowed_groups,omitempty"`
	// ExemptFromDisallowedGroups are exact file paths that we won't fail on
	// even if they're in a disallowed group. We use it for files that already
	// exist in the history of the branch we sync to and can't be removed.
	ExemptFromDisallowedGroups []string `hcl:"exempt_from_disallowed_groups,optional" json:"exempt_from_disallowed_groups,omitempty"`
}

// SyncPolicy returns the sync policy with the given name.
func (c *Config) SyncPolicy(name string) (*SyncPolicyConfig, error) {
	if c == nil || c.Sync == nil {
		return nil, fmt.Errorf("sync policy %s was requested but no sync config was found in pipeline.hcl", name)
	}

	for _, policy := range c.Sync.Policies {
		if policy.Name == name {
			return policy, nil
		}
	}

	return nil, fmt.Errorf("no sync policy named %s was found in pipeline.hcl", name)
}

// validate checks the things the HCL decoder can't. The decoder is happy to
// give us more than one policy with the same name, which would make
// --sync-policy ambiguous.
func (s *SyncConfig) validate() error {
	if s == nil {
		return nil
	}

	names := map[string]struct{}{}
	for _, policy := range s.Policies {
		if _, ok := names[policy.Name]; ok {
			return fmt.Errorf("sync policy %s has been defined more than once", policy.Name)
		}
		names[policy.Name] = struct{}{}
	}

	return nil
}
