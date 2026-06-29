// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package transit

import (
	"context"
	"errors"

	"github.com/hashicorp/vault/sdk/helper/keysutil"
)

var errEntOnly = errors.New("managed keys are supported within enterprise edition only")

func GetManagedKeyUUID(ctx context.Context, b *backend, keyName string, keyId string) (uuid string, err error) {
	return "", errEntOnly
}

func getFormattedManagedKeyPublicKey(_ context.Context, _ *backend, _ *keysutil.Policy) (map[string]map[string]interface{}, error) {
	return nil, nil
}
