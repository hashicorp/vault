// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package tokenutil

import (
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func entTokenFields(fields map[string]*framework.FieldSchema) map[string]*framework.FieldSchema {
	return fields
}

func (t *TokenParams) entParseTokenFields(d *framework.FieldData) {}

func (t *TokenParams) entPopulateTokenData(m map[string]any) {}

func (t *TokenParams) entPopulateTokenAuth(auth *logical.Auth) {}
