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
func TestCoreInit(t *testing.T, core *Core) [][]byte {
	result, err := core.Initialize(&SealConfig{
		SecretShares:    1,
		SecretThreshold: 1,
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return result.SecretShares
}
