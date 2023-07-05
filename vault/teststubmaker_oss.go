// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package vault

import (
	"github.com/hashicorp/vault/vault/cluster"
	"github.com/hashicorp/vault/vault/seal"
)

//go:generate go run github.com/hashicorp/vault/tools/stubmaker
func myfunc(cluster.Listener) *seal.Envelope { return nil }
func myfunc2() *seal.Envelope                { return nil }
