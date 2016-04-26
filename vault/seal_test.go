package vault

import (
	"reflect"
	"testing"
)

func TestDefaultSeal_Config(t *testing.T) {
	bc, _ := TestSealDefConfigs()
	// Change these to non-default values to ensure we are seeing the real
	// config we set
	bc.SecretShares = 4
	bc.SecretThreshold = 2

	core, _, _ := TestCoreUnsealed(t)

	defSeal := &DefaultSeal{}
	defSeal.SetCore(core)
	err := defSeal.SetBarrierConfig(bc)
	if err != nil {
		t.Fatal(err)
	}

	newBc, err := defSeal.BarrierConfig()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(*bc, *newBc) {
		t.Fatal("config mismatch")
	}

	// Now, test without the benefit of the cached value in the seal
	defSeal = &DefaultSeal{}
	defSeal.SetCore(core)
	newBc, err = defSeal.BarrierConfig()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(*bc, *newBc) {
		t.Fatal("config mismatch")
	}
}
