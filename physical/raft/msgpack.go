// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package raft

// If we downgrade msgpack from v1.1.5 to v0.5.5, everything will still
// work, but any pre-existing raft clusters will break on upgrade.
// This file exists so that the Vault project has an explicit dependency
// on the library, which allows us to pin the version in go.mod.

import (
	_ "github.com/hashicorp/go-msgpack/codec"
)
