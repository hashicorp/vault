// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package http

import (
	"net/http"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

//go:generate go run github.com/hashicorp/vault/tools/stubmaker

func entHandleEventsSubscribe(core *vault.Core, req *logical.Request) http.Handler {
	return nil
}
