package vault

import (
	"encoding/base64"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/physical"
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

	bc1, err := defSeal.BarrierConfig()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(*bc, *bc1) {
		t.Fatal("config mismatch")
	}

	// Now, test without the benefit of the cached value in the seal
	defSeal = &DefaultSeal{}
	defSeal.SetCore(core)
	bc2, err := defSeal.BarrierConfig()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(*bc, *bc2) {
		t.Fatal("config mismatch")
	}

	defSeal = &DefaultSeal{}
	defSeal.SetCore(core)
	if err := defSeal.checkCore(); err != nil {
		t.Fatal(err)
	}
	bc.Type = defSeal.BarrierType()
	buf, err := json.Marshal(bc)
	if err != nil {
		t.Fatalf("failed to encode seal configuration: %v", err)
	}

	pe := &physical.Entry{
		Key: barrierSealConfigPath,
		// Store the base64 encoded barrier configuration
		Value: []byte(base64.StdEncoding.EncodeToString(buf)),
	}
	if err := defSeal.core.physical.Put(pe); err != nil {
		t.Fatalf("failed to write seal configuration: %v", err)
	}
	defSeal.config = bc.Clone()

	bc3, err := defSeal.BarrierConfig()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(*bc, *bc3) {
		t.Fatal("config mismatch")
	}
}
