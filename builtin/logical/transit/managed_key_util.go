// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package transit

import (
	"context"
	"errors"
)

var errEntOnly = errors.New("managed keys are supported within enterprise edition only")

func GetManagedKeyUUID(ctx context.Context, b *backend, keyName string, keyId string) (uuid string, err error) {
	return "", errEntOnly
}
