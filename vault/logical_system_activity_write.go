// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !testonly

package vault

import (
	"github.com/hashicorp/vault/sdk/framework"
)

func (b *SystemBackend) activityWritePath() *framework.Path { return nil }
