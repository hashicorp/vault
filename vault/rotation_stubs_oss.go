// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1package vault

//go:build !enterprise

package vault

//go:generate go run github.com/hashicorp/vault/tools/stubmaker

import (
	"context"

	"github.com/hashicorp/vault/sdk/helper/automatedrotationutil"
	"github.com/hashicorp/vault/sdk/rotation"
)

type RotationManager struct{}

func (c *Core) startRotation() error {
	return nil
}

func (c *Core) stopRotation() error {
	return nil
}

func (c *Core) GetRotationInformation(_ context.Context, _ string, _ *rotation.RotationInfoRequest) (*rotation.RotationInfoResponse, error) {
	return nil, automatedrotationutil.ErrRotationManagerUnsupported
}

func (c *Core) RegisterRotationJob(_ context.Context, _ *rotation.RotationJob) (string, error) {
	return "", automatedrotationutil.ErrRotationManagerUnsupported
}

// The DeregisterRotationJob stub returns nil instead of an error because it is intended to be valid to send a deregister
// request for a non-existent job. As a result, the plugin sends a deregister request whenever the relevant rotation
// values are unset. This means that for a plugin running in CE Vault, it will _always_ try to send a deregister request.
func (c *Core) DeregisterRotationJob(_ context.Context, _ *rotation.RotationJobDeregisterRequest) error {
	return nil
}
