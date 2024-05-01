// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !testonly

package vault

import (
	"github.com/hashicorp/vault/sdk/framework"
)

func (b *SystemBackend) requestLimiterReadPath() *framework.Path { return nil }
