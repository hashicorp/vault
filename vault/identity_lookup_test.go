package vault

import (
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestIdentityStore_Lookup_Entity(t *testing.T) {
	var err error
	var resp *logical.Response

	i, _, _ := testIdentityStoreWithGithubAuth(t)

	entityReq := &logical.Request{
		Path:      "entity",
		Operation: logical.UpdateOperation,
	}
	resp, err = i.HandleRequest(entityReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %#v\nresp: %v", err, resp)
	}
	entityID := resp.Data["id"].(string)

	entity, err := i.memDBEntityByID(entityID, false)
	if err != nil {
		t.Fatal(err)
	}

	lookupReq := &logical.Request{
		Path:      "lookup/entity",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"type": "id",
			"id":   entityID,
		},
	}
	resp, err = i.HandleRequest(lookupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %#v\nresp: %v", err, resp)
	}

	if resp.Data["id"].(string) != entityID {
		t.Fatalf("bad: entity: %#v", resp.Data)
	}

	lookupReq.Data = map[string]interface{}{
		"type": "name",
		"name": entity.Name,
	}

	resp, err = i.HandleRequest(lookupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %#v\nresp: %v", err, resp)
	}

	if resp.Data["id"].(string) != entityID {
		t.Fatalf("bad: entity: %#v", resp.Data)
	}
}
