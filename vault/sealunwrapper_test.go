// +build !ent
// +build !prem
// +build !pro
// +build !hsm

package vault

import (
	"bytes"
	"context"
	"sync"
	"testing"

	proto "github.com/golang/protobuf/proto"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/logbridge"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/physical/inmem"
)

func TestSealUnwrapper(t *testing.T) {
	logger := logbridge.NewLogger(hclog.New(&hclog.LoggerOptions{
		Mutex: &sync.Mutex{},
	}))

	// Test without transactions
	phys, err := inmem.NewInmemHA(nil, logger.LogxiLogger())
	if err != nil {
		t.Fatal(err)
	}
	performTestSealUnwrapper(t, phys, logger)

	// Test with transactions
	tPhys, err := inmem.NewTransactionalInmemHA(nil, logger.LogxiLogger())
	if err != nil {
		t.Fatal(err)
	}
	performTestSealUnwrapper(t, tPhys, logger)
}

func performTestSealUnwrapper(t *testing.T, phys physical.Backend, logger *logbridge.Logger) {
	ctx := context.Background()
	base := &CoreConfig{
		Physical: phys,
	}
	cluster := NewTestCluster(t, base, &TestClusterOptions{
		RawLogger: logger,
	})
	cluster.Start()
	defer cluster.Cleanup()

	// Read a value and then save it back in a proto message
	entry, err := phys.Get(ctx, "core/master")
	if err != nil {
		t.Fatal(err)
	}
	if len(entry.Value) == 0 {
		t.Fatal("got no value for master")
	}
	// Save the original for comparison later
	origBytes := make([]byte, len(entry.Value))
	copy(origBytes, entry.Value)
	se := &physical.SealWrapEntry{
		Ciphertext: entry.Value,
	}
	seb, err := proto.Marshal(se)
	if err != nil {
		t.Fatal(err)
	}
	// Write the canary
	entry.Value = append(seb, 's')
	// Save the protobuf value for comparison later
	pBytes := make([]byte, len(entry.Value))
	copy(pBytes, entry.Value)
	if err = phys.Put(ctx, entry); err != nil {
		t.Fatal(err)
	}

	// At this point we should be able to read through the standby cores,
	// successfully decode it, but be able to unmarshal it when read back from
	// the underlying physical store. When we read from active, it should both
	// successfully decode it and persist it back.
	checkValue := func(core *Core, wrapped bool) {
		entry, err := core.physical.Get(ctx, "core/master")
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(entry.Value, origBytes) {
			t.Fatalf("mismatched original bytes and unwrapped entry bytes: %v vs %v", entry.Value, origBytes)
		}
		underlyingEntry, err := phys.Get(ctx, "core/master")
		if err != nil {
			t.Fatal(err)
		}
		switch wrapped {
		case true:
			if !bytes.Equal(underlyingEntry.Value, pBytes) {
				t.Fatalf("mismatched original bytes and proto entry bytes: %v vs %v", underlyingEntry.Value, pBytes)
			}
		default:
			if !bytes.Equal(underlyingEntry.Value, origBytes) {
				t.Fatalf("mismatched original bytes and unwrapped entry bytes: %v vs %v", underlyingEntry.Value, origBytes)
			}
		}
	}

	TestWaitActive(t, cluster.Cores[0].Core)
	checkValue(cluster.Cores[2].Core, true)
	checkValue(cluster.Cores[1].Core, true)
	checkValue(cluster.Cores[0].Core, false)
}
