// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"reflect"
	"testing"
)

// TestDefaultSeal_Config exercises Shamir SetBarrierConfig and BarrierConfig.
// Note that this is a little questionable, because we're doing an init and
// unseal, then changing the barrier config using an internal function instead
// of an API.  In other words if your change break this test, it might be more
// the test's fault than your changes.
func TestDefaultSeal_Config(t *testing.T) {
	bc := &SealConfig{
		SecretShares:    4,
		SecretThreshold: 2,
	}
	core, _, _ := TestCoreUnsealed(t)

	defSeal := NewDefaultSeal(nil)
	defSeal.SetCore(core)
	err := defSeal.SetBarrierConfig(context.Background(), bc)
	if err != nil {
		t.Fatal(err)
	}

	newBc, err := defSeal.BarrierConfig(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(*bc, *newBc) {
		t.Fatal("config mismatch")
	}

	// Now, test without the benefit of the cached value in the seal
	defSeal = NewDefaultSeal(nil)
	defSeal.SetCore(core)
	newBc, err = defSeal.BarrierConfig(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(*bc, *newBc) {
		t.Fatal("config mismatch")
	}
}
