// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"
	"errors"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/snapshots"
)

//go:generate go run github.com/hashicorp/vault/tools/stubmaker

type entSystemBackend struct{}

func entUnauthenticatedPaths() []string {
	return []string{}
}

func (s *SystemBackend) entInit() {}

func (s *SystemBackend) makeSnapshotSource(ctx context.Context, _ *framework.FieldData) (snapshots.Source, error) {
	body, ok := logical.ContextOriginalBodyValue(ctx)
	if !ok {
		return nil, errors.New("no reader for request")
	}
	return snapshots.NewManualSnapshotSource(body), nil
}

// mountInfo returns a map of information about the given mount entry
// Enterprise-specific fields are added in the enterprise version of this method.
func (b *SystemBackend) mountInfo(ctx context.Context, entry *MountEntry, legacyTTLFormat bool) map[string]interface{} {
	info := map[string]interface{}{
		"type":                    entry.Type,
		"description":             entry.Description,
		"accessor":                entry.Accessor,
		"local":                   entry.Local,
		"seal_wrap":               entry.SealWrap,
		"external_entropy_access": entry.ExternalEntropyAccess,
		"options":                 entry.Options,
		"uuid":                    entry.UUID,
		"plugin_version":          entry.Version,
		"running_plugin_version":  entry.RunningVersion,
		"running_sha256":          entry.RunningSha256,
	}
	coreDefTTL := int64(b.Core.defaultLeaseTTL.Seconds())
	coreMaxTTL := int64(b.Core.maxLeaseTTL.Seconds())
	entDefTTL := int64(entry.Config.DefaultLeaseTTL.Seconds())
	entMaxTTL := int64(entry.Config.MaxLeaseTTL.Seconds())
	entryConfig := map[string]interface{}{
		"default_lease_ttl": entDefTTL,
		"max_lease_ttl":     entMaxTTL,
		"force_no_cache":    entry.Config.ForceNoCache,
	}
	if !legacyTTLFormat {
		if entDefTTL == 0 {
			entryConfig["default_lease_ttl"] = coreDefTTL
		}
		if entMaxTTL == 0 {
			entryConfig["max_lease_ttl"] = coreMaxTTL
		}
	}
	if entry.Config.TrimRequestTrailingSlashes {
		entryConfig["trim_request_trailing_slashes"] = true
	}
	if rawVal, ok := entry.synthesizedConfigCache.Load("audit_non_hmac_request_keys"); ok {
		entryConfig["audit_non_hmac_request_keys"] = rawVal.([]string)
	}
	if rawVal, ok := entry.synthesizedConfigCache.Load("audit_non_hmac_response_keys"); ok {
		entryConfig["audit_non_hmac_response_keys"] = rawVal.([]string)
	}
	// Even though empty value is valid for ListingVisibility, we can ignore
	// this case during mount since there's nothing to unset/hide.
	if len(entry.Config.ListingVisibility) > 0 {
		entryConfig["listing_visibility"] = entry.Config.ListingVisibility
	}
	if rawVal, ok := entry.synthesizedConfigCache.Load("passthrough_request_headers"); ok {
		entryConfig["passthrough_request_headers"] = rawVal.([]string)
	}
	if rawVal, ok := entry.synthesizedConfigCache.Load("allowed_response_headers"); ok {
		entryConfig["allowed_response_headers"] = rawVal.([]string)
	}
	if rawVal, ok := entry.synthesizedConfigCache.Load("allowed_managed_keys"); ok {
		entryConfig["allowed_managed_keys"] = rawVal.([]string)
	}
	if rawVal, ok := entry.synthesizedConfigCache.Load("identity_token_key"); ok {
		entryConfig["identity_token_key"] = rawVal.(string)
	}
	if entry.Table == credentialTableType {
		entryConfig["token_type"] = entry.Config.TokenType.String()
	}
	if entry.Config.UserLockoutConfig != nil {
		userLockoutConfig := map[string]interface{}{
			"user_lockout_counter_reset_duration": int64(entry.Config.UserLockoutConfig.LockoutCounterReset.Seconds()),
			"user_lockout_threshold":              entry.Config.UserLockoutConfig.LockoutThreshold,
			"user_lockout_duration":               int64(entry.Config.UserLockoutConfig.LockoutDuration.Seconds()),
			"user_lockout_disable":                entry.Config.UserLockoutConfig.DisableLockout,
		}
		entryConfig["user_lockout_config"] = userLockoutConfig
	}
	if rawVal, ok := entry.synthesizedConfigCache.Load("delegated_auth_accessors"); ok {
		entryConfig["delegated_auth_accessors"] = rawVal.([]string)
	}

	// Add deprecation status only if it exists
	builtinType := b.Core.builtinTypeFromMountEntry(ctx, entry)
	if status, ok := b.Core.builtinRegistry.DeprecationStatus(entry.Type, builtinType); ok {
		info["deprecation_status"] = status.String()
	}
	info["config"] = entryConfig

	return info
}

func (b *SystemBackend) callUnsyncMountHelper(ctx context.Context, path string) error {
	return nil
}
