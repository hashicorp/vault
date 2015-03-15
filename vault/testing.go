package vault

import (
	"testing"

	"github.com/hashicorp/vault/physical"
)

// This file contains a number of methods that are useful for unit
// tests within other packages.

// TestCore returns a pure in-memory, uninitialized core for testing.
func TestCore(t *testing.T) *Core {
	physicalBackend := physical.NewInmem()
	c, err := NewCore(&CoreConfig{
		Physical: physicalBackend,
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return c
}

// TestCoreInit initializes the core with a single key, and returns
// the list of keys that must be used to unseal the core.
func TestCoreInit(t *testing.T, core *Core) []byte {
	result, err := core.Initialize(&SealConfig{
		SecretShares:    1,
		SecretThreshold: 1,
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return result.SecretShares[0]
}

// TestCoreUnsealed returns a pure in-memory core that is already
// initialized and unsealed.
func TestCoreUnsealed(t *testing.T) (*Core, []byte) {
	core := TestCore(t)
	key := TestCoreInit(t, core)
	if _, err := core.Unseal(key); err != nil {
		t.Fatalf("unseal err: %s", err)
	}

	sealed, err := core.Sealed()
	if err != nil {
		t.Fatalf("err checking seal status: %s", err)
	}
	if sealed {
		t.Fatal("should not be sealed")
	}

	return core, key
}
