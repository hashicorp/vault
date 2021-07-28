package vault

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

// TestOIDC_Path_OIDCAssignment tests CRUD operations for assignments
func TestOIDC_Path_OIDCAssignment(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test assignment "test-assignment" -- should succeed
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/test-assignment",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-assignment" and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/test-assignment",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected := map[string]interface{}{
		"groups":   []string{},
		"entities": []string{},
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	// Update "test-assignment" -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/test-assignment",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"groups":   "my-group",
			"entities": "my-entity",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-assignment" again and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/test-assignment",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected = map[string]interface{}{
		"groups":   []string{"my-group"},
		"entities": []string{"my-entity"},
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	// Delete test-assignment -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/test-assignment",
		Operation: logical.DeleteOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
}

// // TestOIDC_Path_OIDCKey tests the List operation for keys
// func TestOIDC_Path_OIDCKey(t *testing.T) {
// 	c, _, _ := TestCoreUnsealed(t)
// 	ctx := namespace.RootContext(nil)
// 	storage := &logical.InmemStorage{}

// 	// Prepare two keys, test-key1 and test-key2
// 	c.identityStore.HandleRequest(ctx, &logical.Request{
// 		Path:      "oidc/key/test-key1",
// 		Operation: logical.CreateOperation,
// 		Storage:   storage,
// 	})

// 	c.identityStore.HandleRequest(ctx, &logical.Request{
// 		Path:      "oidc/key/test-key2",
// 		Operation: logical.CreateOperation,
// 		Storage:   storage,
// 	})

// 	// list keys
// 	respListKey, listErr := c.identityStore.HandleRequest(ctx, &logical.Request{
// 		Path:      "oidc/key",
// 		Operation: logical.ListOperation,
// 		Storage:   storage,
// 	})
// 	expectSuccess(t, respListKey, listErr)

// 	// validate list response
// 	expectedStrings := map[string]interface{}{"test-key1": true, "test-key2": true}
// 	expectStrings(t, respListKey.Data["keys"].([]string), expectedStrings)

// 	// delete test-key2
// 	c.identityStore.HandleRequest(ctx, &logical.Request{
// 		Path:      "oidc/key/test-key2",
// 		Operation: logical.DeleteOperation,
// 		Storage:   storage,
// 	})

// 	// list keyes again and validate response
// 	respListKeyAfterDelete, listErrAfterDelete := c.identityStore.HandleRequest(ctx, &logical.Request{
// 		Path:      "oidc/key",
// 		Operation: logical.ListOperation,
// 		Storage:   storage,
// 	})
// 	expectSuccess(t, respListKeyAfterDelete, listErrAfterDelete)

// 	// validate list response
// 	delete(expectedStrings, "test-key2")
// 	expectStrings(t, respListKeyAfterDelete.Data["keys"].([]string), expectedStrings)
// }

// // some helpers
// func expectSuccess(t *testing.T, resp *logical.Response, err error) {
// 	t.Helper()
// 	if err != nil || (resp != nil && resp.IsError()) {
// 		t.Fatalf("expected success but got error:\n%v\nresp: %#v", err, resp)
// 	}
// }

// func expectError(t *testing.T, resp *logical.Response, err error) {
// 	if err == nil {
// 		if resp == nil || !resp.IsError() {
// 			t.Fatalf("expected error but got success; error:\n%v\nresp: %#v", err, resp)
// 		}
// 	}
// }

// // expectString fails unless every string in actualStrings is also included in expectedStrings and
// // the length of actualStrings and expectedStrings are the same
// func expectStrings(t *testing.T, actualStrings []string, expectedStrings map[string]interface{}) {
// 	if len(actualStrings) != len(expectedStrings) {
// 		t.Fatalf("expectStrings mismatch:\nactual strings:\n%#v\nexpected strings:\n%#v\n", actualStrings, expectedStrings)
// 	}
// 	for _, actualString := range actualStrings {
// 		_, ok := expectedStrings[actualString]
// 		if !ok {
// 			t.Fatalf("the string %q was not expected", actualString)
// 		}
// 	}
// }
