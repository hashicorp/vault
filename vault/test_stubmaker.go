// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"github.com/hashicorp/vault/vault/cluster"
)

func test() {
	testStubmaker(cluster.Listener{})
	testStubmaker2()
}
