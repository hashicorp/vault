package awsiam

import (
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
)

func TestBackend_pathRole(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b := Backend()
	_, err := b.Setup(config)
	if err != nil {
		t.Fatal(err)
	}

	// make sure we start with empty roles, which gives us confidence that the read later
	// actually is the two roles we created
	resp, err := b.HandleRequest(&logical.Request{
		Operation: logical.ListOperation,
		Path:      "roles",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Data == nil || resp.IsError() {
		t.Fatalf("failed to list role entries")
	}
	if resp.Data["keys"] != nil {
		t.Fatalf("Received roles when expected none")
	}

	data := map[string]interface{}{
		"policies":            "p,q,r,s",
		"max_ttl":             "2h",
		"bound_iam_principal": "n:aws:iam::123456789012:user/MyUserName",
	}
	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/MyRoleName",
		Data:      data,
		Storage:   storage,
	})

	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.IsError() {
		t.Fatal("failed to create the role entry")
	}

	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.ReadOperation,
		Path:      "role/MyRoleName",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.IsError() {
		t.Fatal("failed to read the role entry")
	}
	if !policyutil.EquivalentPolicies(strings.Split(data["policies"].(string), ","), resp.Data["policies"].([]string)) {
		t.Fatalf("bad: policies: expected %#v\ngot: %#v\n", data, resp.Data)
	}

	// generate a second role, ensure we're able to list both
	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/MyOtherRoleName",
		Data:      data,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to create additional role")
	}

	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.ListOperation,
		Path:      "roles",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Data == nil || resp.IsError() {
		t.Fatalf("failed to list role entries")
	}
	keys := resp.Data["keys"].([]string)
	if len(keys) != 2 {
		t.Fatalf("bad: keys %#v\n", keys)
	}

	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.DeleteOperation,
		Path:      "role/MyOtherRoleName",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.ReadOperation,
		Path:      "role/MyOtherRoleName",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil {
		t.Fatalf("bad: response: expected: nil actual:%3v\n", resp)
	}
}
