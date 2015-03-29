package vault

import (
	"reflect"
	"testing"
)

func TestCore_DefaultMountTable(t *testing.T) {
	c, key, _ := TestCoreUnsealed(t)
	verifyDefaultTable(t, c.mounts)

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
	if !reflect.DeepEqual(c.mounts, c2.mounts) {
		t.Fatalf("mismatch: %v %v", c.mounts, c2.mounts)
	}
}

func TestCore_Mount(t *testing.T) {
	c, key, _ := TestCoreUnsealed(t)
	me := &MountEntry{
		Path: "foo",
		Type: "generic",
	}
	err := c.mount(me)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	match := c.router.MatchingMount("foo/bar")
	if match != "foo/" {
		t.Fatalf("missing mount")
	}

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
	if !reflect.DeepEqual(c.mounts, c2.mounts) {
		t.Fatalf("mismatch: %v %v", c.mounts, c2.mounts)
	}
}

func TestCore_Unmount(t *testing.T) {
	c, key, _ := TestCoreUnsealed(t)
	err := c.unmount("secret")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	match := c.router.MatchingMount("secret/foo")
	if match != "" {
		t.Fatalf("backend present")
	}

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
	if !reflect.DeepEqual(c.mounts, c2.mounts) {
		t.Fatalf("mismatch: %v %v", c.mounts, c2.mounts)
	}
}

func TestCore_Remount(t *testing.T) {
	c, key, _ := TestCoreUnsealed(t)
	err := c.remount("secret", "foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	match := c.router.MatchingMount("foo/bar")
	if match != "foo/" {
		t.Fatalf("failed remount")
	}

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
	if !reflect.DeepEqual(c.mounts, c2.mounts) {
		t.Fatalf("mismatch: %v %v", c.mounts, c2.mounts)
	}
}

func TestCore_Remount_Protected(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	err := c.remount("sys", "foo")
	if err.Error() != "cannot remount 'sys/'" {
		t.Fatalf("err: %v", err)
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
