// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package pki

import (
	"crypto/x509/pkix"
	"errors"
)

func entParseNameConstraintsJson(nameConstraintsString string) (*pkix.Extension, error) {
	return nil, errors.New("name_constraints is an ent-only field")
}
