package salt

import (
	"crypto/sha1"
	"crypto/sha256"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/logical"
)

func TestSalt(t *testing.T) {
	inm := &logical.InmemStorage{}
	conf := &Config{}

	salt, err := NewSalt(inm, conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !salt.DidGenerate() {
		t.Fatalf("expected generation")
	}

	// Verify the salt exists
	out, err := inm.Get(DefaultLocation)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out == nil {
		t.Fatalf("missing salt")
	}

	// Create a new salt, should restore
	salt2, err := NewSalt(inm, conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if salt2.DidGenerate() {
		t.Fatalf("unexpected generation")
	}

	// Check for a match
	if salt.salt != salt2.salt {
		t.Fatalf("salt mismatch: %s %s", salt.salt, salt2.salt)
	}

	// Verify a match
	id := "foobarbaz"
	sid1 := salt.SaltID(id)
	sid2 := salt2.SaltID(id)

	if sid1 != sid2 {
		t.Fatalf("mismatch")
	}
}

func TestSaltID(t *testing.T) {
	salt, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	id := "foobarbaz"

	sid1 := SaltID(salt, id, SHA1Hash)
	sid2 := SaltID(salt, id, SHA1Hash)

	if len(sid1) != sha1.Size*2 {
		t.Fatalf("Bad len: %d %s", len(sid1), sid1)
	}

	if sid1 != sid2 {
		t.Fatalf("mismatch")
	}

	sid1 = SaltID(salt, id, SHA256Hash)
	sid2 = SaltID(salt, id, SHA256Hash)

	if len(sid1) != sha256.Size*2 {
		t.Fatalf("Bad len: %d", len(sid1))
	}

	if sid1 != sid2 {
		t.Fatalf("mismatch")
	}
}
