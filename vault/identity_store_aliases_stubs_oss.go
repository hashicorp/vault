// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import "context"

func (c *Core) validateSyntheticAliasAccessor(context.Context, string) (bool, bool, error) {
	return false, false, nil
}

func (c *Core) generateSyntheticAliasAccessor(context.Context, string) (string, bool, error) {
	return "", false, nil
}

func (c *Core) resolveProfileNameToAccessor(_ context.Context, _ string) (resolvedProfileRef, bool, error) {
	return resolvedProfileRef{}, false, errProfileRefUnsupported
}

func (c *Core) resolveConfigIDToAccessor(_ context.Context, _ string) (resolvedProfileRef, bool, error) {
	return resolvedProfileRef{}, false, errProfileRefUnsupported
}

// resolveAccessorToProfileRef reports not-found rather than
// errProfileRefUnsupported: unlike 'profile_name' and 'config_id', which a caller
// can only have supplied deliberately, this is asked of every 'mount_accessor',
// and on CE no accessor is ever a synthetic OAuth RS accessor.
func (c *Core) resolveAccessorToProfileRef(_ context.Context, _ string) (resolvedProfileRef, bool, error) {
	return resolvedProfileRef{}, false, nil
}
