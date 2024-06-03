// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build minimal

package command

import (
	_ "github.com/hashicorp/vault/helper/builtinplugins"
	physConsul "github.com/hashicorp/vault/physical/consul"
	physRaft "github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/physical"
)

var (
	physicalBackends = map[string]physical.Factory{
		"consul": physConsul.NewConsulBackend,
		"raft":   physRaft.NewRaftBackend,
	}

	loginHandlers = map[string]LoginHandler{}
)
