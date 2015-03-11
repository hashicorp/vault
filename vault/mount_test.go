package vault

import (
	"reflect"
	"testing"
)

func TestCore_DefaultMountTable(t *testing.T) {
	c, key := testUnsealedCore(t)
	verifyDefaultTable(t, c.mounts)

	// Start a second core with same physical
	conf := &CoreConfig{physical: c.physical}
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
	if !reflect.DeepEqual(c.mounts, c2.mounts) {
		t.Fatalf("mismatch: %v %v", c.mounts, c2.mounts)
	}
}

func TestDefaultMountTable(t *testing.T) {
	table := defaultMountTable()
	verifyDefaultTable(t, table)
}

func verifyDefaultTable(t *testing.T, table *MountTable) {
	if len(table.Entries) != 2 {
		t.Fatalf("bad: %v", table.Entries)
	}
	for idx, entry := range table.Entries {
		switch idx {
		case 0:
			if entry.Path != "secret/" {
				t.Fatalf("bad: %v", entry)
			}
			if entry.Type != "generic" {
				t.Fatalf("bad: %v", entry)
			}
		case 1:
			if entry.Path != "sys/" {
				t.Fatalf("bad: %v", entry)
			}
			if entry.Type != "system" {
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
