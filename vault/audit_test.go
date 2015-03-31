package vault

import (
	"reflect"
	"testing"

	"github.com/hashicorp/vault/audit"
)

type NoopAudit struct{}

func TestCore_EnableAudit(t *testing.T) {
	c, key, _ := TestCoreUnsealed(t)
	c.auditBackends["noop"] = func(map[string]string) (audit.Backend, error) {
		return &NoopAudit{}, nil
	}

	me := &MountEntry{
		Path: "foo",
		Type: "noop",
	}
	err := c.enableAudit(me)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !c.auditBroker.IsRegistered("foo") {
		t.Fatalf("missing audit backend")
	}

	conf := &CoreConfig{
		Physical:      c.physical,
		AuditBackends: make(map[string]audit.Factory),
	}
	conf.AuditBackends["noop"] = func(map[string]string) (audit.Backend, error) {
		return &NoopAudit{}, nil
	}
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

	// Verify matching audit tables
	if !reflect.DeepEqual(c.audit, c2.audit) {
		t.Fatalf("mismatch: %v %v", c.audit, c2.audit)
	}

	// Check for registration
	if !c2.auditBroker.IsRegistered("foo") {
		t.Fatalf("missing audit backend")
	}
}

func TestCore_DisableAudit(t *testing.T) {
	c, key, _ := TestCoreUnsealed(t)
	c.auditBackends["noop"] = func(map[string]string) (audit.Backend, error) {
		return &NoopAudit{}, nil
	}

	err := c.disableAudit("foo")
	if err.Error() != "no matching backend" {
		t.Fatalf("err: %v", err)
	}

	me := &MountEntry{
		Path: "foo",
		Type: "noop",
	}
	err = c.enableAudit(me)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	err = c.disableAudit("foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Check for registration
	if c.auditBroker.IsRegistered("foo") {
		t.Fatalf("audit backend present")
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
	if !reflect.DeepEqual(c.audit, c2.audit) {
		t.Fatalf("mismatch: %v %v", c.audit, c2.audit)
	}
}

func TestCore_DefaultAuditTable(t *testing.T) {
	c, key, _ := TestCoreUnsealed(t)
	verifyDefaultAuditTable(t, c.audit)

	// Verify we have an audit broker
	if c.auditBroker == nil {
		t.Fatalf("missing audit broker")
	}

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
	if !reflect.DeepEqual(c.audit, c2.audit) {
		t.Fatalf("mismatch: %v %v", c.audit, c2.audit)
	}
}

func TestDefaultAuditTable(t *testing.T) {
	table := defaultAuditTable()
	verifyDefaultAuditTable(t, table)
}

func verifyDefaultAuditTable(t *testing.T, table *MountTable) {
	if len(table.Entries) != 0 {
		t.Fatalf("bad: %v", table.Entries)
	}
}
