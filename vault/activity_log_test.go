package vault

import (
	"testing"
	"time"
)

func TestActivityLog_Creation(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)

	a := core.activityLog

	if a == nil {
		t.Fatal("no activity log found")
	}
	if a.logger == nil || a.view == nil {
		t.Fatal("activity log not initialized")
	}
	if a.fragment != nil {
		t.Fatal("activity log already has fragment")
	}

	const entity_id = "entity_id_75432"
	const namespace_id = "ns123"
	ts := time.Now()

	a.AddEntityToFragment(entity_id, namespace_id, ts)
	if a.fragment == nil {
		t.Fatal("no fragment created")
	}

	if a.fragment.OriginatingNode != a.nodeID {
		t.Errorf("mismatched node ID, %q vs %q", a.fragment.OriginatingNode, a.nodeID)
	}

	if a.fragment.Entities == nil {
		t.Fatal("no fragment entity slice")
	}

	if a.fragment.NonEntityTokens == nil {
		t.Fatal("no fragment token map")
	}

	if len(a.fragment.Entities) != 1 {
		t.Fatalf("wrong number of entities %v", len(a.fragment.Entities))
	}

	er := a.fragment.Entities[0]
	if er.EntityID != entity_id {
		t.Errorf("mimatched entity ID, %q vs %q", er.EntityID, entity_id)
	}
	if er.NamespaceID != namespace_id {
		t.Errorf("mimatched namespace ID, %q vs %q", er.NamespaceID, namespace_id)
	}
	if er.Timestamp != ts.UnixNano() {
		t.Errorf("mimatched timestamp, %v vs %v", er.Timestamp, ts.UnixNano())
	}

	// Reset and test the other code path
	a.fragment = nil
	a.AddTokenToFragment(namespace_id)

	if a.fragment == nil {
		t.Fatal("no fragment created")
	}

	if a.fragment.NonEntityTokens == nil {
		t.Fatal("no fragment token map")
	}

	actual := a.fragment.NonEntityTokens[namespace_id]
	if actual != 1 {
		t.Errorf("mismatched number of tokens, %v vs %v", actual, 1)
	}

}
