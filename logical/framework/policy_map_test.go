package framework

import (
	"reflect"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestPolicyMap(t *testing.T) {
	p := &PolicyMap{}
	p.PathMap.Name = "foo"
	s := new(logical.InmemStorage)

	p.Put(s, "foo", map[string]interface{}{"value": "bar"})
	p.Put(s, "bar", map[string]interface{}{"value": "foo,baz "})

	// Read via API
	actual, err := p.Policies(s, "foo", "bar")
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}

	expected := []string{"bar", "baz", "foo"}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v", actual)
	}
}
