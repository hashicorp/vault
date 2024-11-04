// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package pki

import (
	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/sdk/framework"
)

//go:generate go run github.com/hashicorp/vault/tools/stubmaker

func addEntPathIssuerFields(fields map[string]*framework.FieldSchema)         {}
func addEntPathIssuerResponseFields(fields map[string]*framework.FieldSchema) {}
func setEntIssuerData(data map[string]any, issuer *issuing.IssuerEntry)       {}

func updateEntIssuerFields(issuer *issuing.IssuerEntry, data *framework.FieldData, ignoreNotPresent bool) bool {
	return false
}
