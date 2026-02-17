// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

//go:build !testonly && !ent

package vault

import (
	"github.com/hashicorp/vault/sdk/framework"
)

func (b *SystemBackend) consumptionBillingWritePath() *framework.Path { return nil }
