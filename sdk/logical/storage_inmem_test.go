// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logical

import (
	"testing"
)

func TestInmemStorage(t *testing.T) {
	TestStorage(t, new(InmemStorage))
}
