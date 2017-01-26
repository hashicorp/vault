package vault

import (
	"reflect"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func mockConfigStore(t *testing.T) *ConfigStore {
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "foo/")
	c := NewConfigStore(view, logical.TestSystemView())
	return c
}

func mockConfigStoreNoCache(t *testing.T) *ConfigStore {
	sysView := logical.TestSystemView()
	sysView.CachingDisabledVal = true
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "foo/")
	c := NewConfigStore(view, sysView)
	return c
}

func TestConfigStore_CRUD(t *testing.T) {
	cs := mockConfigStore(t)
	testConfigStore_CRUD(t, cs)

	cs = mockConfigStoreNoCache(t)
	testConfigStore_CRUD(t, cs)
}

func testConfigStore_CRUD(t *testing.T, cs *ConfigStore) {
	// Get should return nothing
	c, err := cs.GetConfig("dev")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if c != nil {
		t.Fatalf("bad: %v", c)
	}

	// Delete should be no-op
	err = cs.DeleteConfig("dev")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// List should be blank
	out, err := cs.ListConfigs()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("bad: %v", out)
	}

	config := &Config{
		Name: "cors",
		Settings: map[string]string{
			"allowed_origins": "http://www.example.com http://localhost",
			"enabled":         "true",
		},
	}

	err = cs.SetConfig(config)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	// Get should work
	c, err = cs.GetConfig("cors")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !reflect.DeepEqual(c, config) {
		t.Fatalf("bad: %v", c)
	}

	// List should be one element
	out, err = cs.ListConfigs()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(out) != 1 || out[0] != "cors" {
		t.Fatalf("bad: %v", out)
	}

	// Delete should be clear the entry
	err = cs.DeleteConfig("cors")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Get should fail
	c, err = cs.GetConfig("cors")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if c != nil {
		t.Fatalf("bad: %v", c)
	}
}
