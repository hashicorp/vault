// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package config

import (
	"errors"
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/changed"
)

// SyncConfig holds our policies for synchronizing branches.
type SyncConfig struct {
	Policies []*SyncPolicy `hcl:"policy,block" json:"policies,omitempty"`
}

// SyncPolicy is a named sync policy. Commands that support it select it with
// --sync-policy <name>.
type SyncPolicy struct {
	Name string `hcl:"name,label" json:"name"`
	// DisallowedGroups are changed file groups that must not exist in the
	// branch. We add them to any groups passed with --disallowed-groups.
	DisallowedGroups changed.FileGroups `hcl:"disallowed_groups,optional" json:"disallowed_groups,omitempty"`
	// ExemptFromDisallowedGroups are exact file paths that we won't fail on
	// even if they're in a disallowed group. We use it for files that already
	// exist in the history of the branch we sync to and can't be removed.
	ExemptFromDisallowedGroups []string `hcl:"exempt_from_disallowed_groups,optional" json:"exempt_from_disallowed_groups,omitempty"`
}

// SyncPolicy returns the sync policy with the given name.
func (c *Config) SyncPolicy(name string) (*SyncPolicy, error) {
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

// DisallowedGroupsOpt is a functional option for DisallowedGroups()
type DisallowedGroupsOpt func(*disallowedGroupsOpts)

type disallowedGroupsOpts struct {
	groups changed.FileGroups
}

// WithDisallowedGroupsGroups adds the groups to the disallowed groups.
func WithDisallowedGroupsGroups(groups changed.FileGroups) DisallowedGroupsOpt {
	return func(opts *disallowedGroupsOpts) {
		opts.groups = opts.groups.Add(groups...)
	}
}

// WithDisallowedGroupsSyncPolicies adds the disallowed groups of each sync
// policy to the disallowed groups.
func WithDisallowedGroupsSyncPolicies(policies ...*SyncPolicy) DisallowedGroupsOpt {
	return func(opts *disallowedGroupsOpts) {
		for _, policy := range policies {
			if policy == nil {
				continue
			}

			opts.groups = opts.groups.Add(policy.DisallowedGroups...)
		}
	}
}

// DisallowedGroups takes variable options and returns all of the groups that we
// disallow.
func (c *Config) DisallowedGroups(opts ...DisallowedGroupsOpt) (changed.FileGroups, error) {
	if c == nil || c.ChangedFiles == nil {
		return nil, errors.New("no changed file grouping config was found in pipeline.hcl")
	}

	o := &disallowedGroupsOpts{groups: changed.FileGroups{}}
	for _, opt := range opts {
		opt(o)
	}

	// A sync policy doesn't have to disallow any groups, but if it doesn't and
	// we weren't given any either then we've got nothing to check.
	if len(o.groups) < 1 {
		return nil, errors.New("no disallowed groups have been configured")
	}

	// Make sure every disallowed group exists. A misspelled group would never
	// match anything and would quietly pass the check.
	configuredGroups := changed.FileGroups{}
	for _, group := range c.ChangedFiles.Groups {
		configuredGroups = configuredGroups.Add(changed.FileGroup(group.Name))
	}
	for _, group := range o.groups {
		if _, ok := configuredGroups.In(group); !ok {
			return nil, fmt.Errorf("disallowed group %s is not a changed file group in pipeline.hcl", group)
		}
	}

	return o.groups, nil
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
