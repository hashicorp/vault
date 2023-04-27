// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package benchhelpers

import (
	"testing"

	testinginterface "github.com/mitchellh/go-testing-interface"
)

type tbWrapper struct {
	testing.TB
}

func (b tbWrapper) Parallel() {
	// no-op
}

func TBtoT(tb testing.TB) testinginterface.T {
	return tbWrapper{tb}
}
