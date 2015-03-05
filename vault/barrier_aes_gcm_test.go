package vault

import (
	"testing"

	"github.com/hashicorp/vault/physical"
)

func TestAESGCMBarrier_Basic(t *testing.T) {
	inm := physical.NewInmem()
	b, err := NewAESGCMBarrier(inm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	testBarrier(t, b)
}

func TestAESGCMBarrier_Confidential(t *testing.T) {
	// TODO: Verify data sent through is encrypted
}

func TestAESGCMBarrier_Integrity(t *testing.T) {
	// TODO: Verify data sent through is cannot be tampered
}
