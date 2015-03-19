package vault

import (
	"reflect"
	"testing"
)

func TestCore_DefaultAuthTable(t *testing.T) {
	c, key := TestCoreUnsealed(t)
	verifyDefaultAuthTable(t, c.auth)

	// Start a second core with same physical
	conf := &CoreConfig{Physical: c.physical}
	c2, err := NewCore(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	unseal, err := c2.Unseal(key)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !unseal {
		t.Fatalf("should be unsealed")
	}

	// Verify matching mount tables
	if !reflect.DeepEqual(c.auth, c2.auth) {
		t.Fatalf("mismatch: %v %v", c.auth, c2.auth)
	}
}

func TestDefaultAuthTable(t *testing.T) {
	table := defaultAuthTable()
	verifyDefaultAuthTable(t, table)
}

func verifyDefaultAuthTable(t *testing.T, table *MountTable) {
	if len(table.Entries) != 1 {
		t.Fatalf("bad: %v", table.Entries)
	}
	for idx, entry := range table.Entries {
		switch idx {
		case 0:
			if entry.Path != "token" {
				t.Fatalf("bad: %v", entry)
			}
			if entry.Type != "token" {
				t.Fatalf("bad: %v", entry)
			}
		}
		if entry.Description == "" {
			t.Fatalf("bad: %v", entry)
		}
		if entry.UUID == "" {
			t.Fatalf("bad: %v", entry)
		}
	}
}
