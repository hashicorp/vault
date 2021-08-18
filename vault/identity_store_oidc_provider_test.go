package vault

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

// TestOIDC_Path_OIDC_ProviderClient_NoKeyParameter tests that a client cannot
// be created without a key parameter
func TestOIDC_Path_OIDC_ProviderClient_NoKeyParameter(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test client "test-client1" without a key param -- should fail
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client1",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})
	expectError(t, resp, err)
	// validate error message
	expectedStrings := map[string]interface{}{
		"the key parameter is required": true,
	}
	expectStrings(t, []string{resp.Data["error"].(string)}, expectedStrings)
}

// TestOIDC_Path_OIDC_ProviderClient_NilKeyEntry tests that a client cannot be
// created when a key parameter is provided but the key does not exist
func TestOIDC_Path_OIDC_ProviderClient_NilKeyEntry(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test client "test-client1" with a non-existent key -- should fail
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client1",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"key": "test-key",
		},
		Storage: storage,
	})
	expectError(t, resp, err)
	// validate error message
	expectedStrings := map[string]interface{}{
		"cannot find key \"test-key\"": true,
	}
	expectStrings(t, []string{resp.Data["error"].(string)}, expectedStrings)
}

// TestOIDC_Path_OIDC_ProviderClient_AssignmentDoesNotExist tests that a client
// cannot be created with assignments that do not exist
func TestOIDC_Path_OIDC_ProviderClient_AssignmentDoesNotExist(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test key "test-key"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"verification_ttl": "2m",
			"rotation_period":  "2m",
		},
		Storage: storage,
	})

	// Create a test client "test-client" -- should fail
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"key":         "test-key",
			"assignments": "my-assignment",
		},
	})
	expectError(t, resp, err)
	// validate error message
	expectedStrings := map[string]interface{}{
		"cannot find assignment \"my-assignment\"": true,
	}
	expectStrings(t, []string{resp.Data["error"].(string)}, expectedStrings)
}

// TestOIDC_Path_OIDC_ProviderClient tests CRUD operations for clients
func TestOIDC_Path_OIDC_ProviderClient(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test key "test-key"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"verification_ttl": "2m",
			"rotation_period":  "2m",
		},
		Storage: storage,
	})

	// Create a test client "test-client" -- should succeed
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"key": "test-key",
		},
	})
	expectSuccess(t, resp, err)

	// Read "test-client" and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected := map[string]interface{}{
		"redirect_uris":    []string{},
		"assignments":      []string{},
		"key":              "test-key",
		"id_token_ttl":     0,
		"access_token_ttl": 0,
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	// Create a test assignment "my-assignment" -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/my-assignment",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)

	// Update "test-client" -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"redirect_uris":    "http://localhost:3456/callback",
			"assignments":      "my-assignment",
			"key":              "test-key",
			"id_token_ttl":     0,
			"access_token_ttl": 0,
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-client" again and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected = map[string]interface{}{
		"redirect_uris":    []string{"http://localhost:3456/callback"},
		"assignments":      []string{"my-assignment"},
		"key":              "test-key",
		"id_token_ttl":     0,
		"access_token_ttl": 0,
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	// Delete test-client -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.DeleteOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-client" again and validate
	resp, _ = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	if resp != nil {
		t.Fatalf("expected nil but got resp: %#v", resp)
	}
}

// TestOIDC_Path_OIDC_ProviderClient_Update tests Update operations for clients
func TestOIDC_Path_OIDC_ProviderClient_Update(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test key "test-key"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"verification_ttl": "2m",
			"rotation_period":  "2m",
		},
		Storage: storage,
	})

	// Create a test assignment "my-assignment" -- should succeed
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/my-assignment",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)

	// Create a test client "test-client" -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"redirect_uris":    "http://localhost:3456/callback",
			"assignments":      "my-assignment",
			"key":              "test-key",
			"id_token_ttl":     0,
			"access_token_ttl": 0,
		},
	})
	expectSuccess(t, resp, err)

	// Read "test-client" and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected := map[string]interface{}{
		"redirect_uris":    []string{"http://localhost:3456/callback"},
		"assignments":      []string{"my-assignment"},
		"key":              "test-key",
		"id_token_ttl":     0,
		"access_token_ttl": 0,
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	// Update "test-client" -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"redirect_uris": "http://localhost:3456/callback2",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-client" again and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected = map[string]interface{}{
		"redirect_uris":    []string{"http://localhost:3456/callback2"},
		"assignments":      []string{"my-assignment"},
		"key":              "test-key",
		"id_token_ttl":     0,
		"access_token_ttl": 0,
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}
}

// TestOIDC_Path_OIDC_ProviderClient_List tests the List operation for clients
func TestOIDC_Path_OIDC_ProviderClient_List(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test key "test-key"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"verification_ttl": "2m",
			"rotation_period":  "2m",
		},
		Storage: storage,
	})

	// Prepare two clients, test-client1 and test-client2
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client1",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"key": "test-key",
		},
	})

	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client2",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"key": "test-key",
		},
	})

	// list clients
	respListClients, listErr := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client",
		Operation: logical.ListOperation,
		Storage:   storage,
	})
	expectSuccess(t, respListClients, listErr)

	// validate list response
	expectedStrings := map[string]interface{}{"test-client1": true, "test-client2": true}
	expectStrings(t, respListClients.Data["keys"].([]string), expectedStrings)

	// delete test-client2
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client2",
		Operation: logical.DeleteOperation,
		Storage:   storage,
	})

	// list clients again and validate response
	respListClientAfterDelete, listErrAfterDelete := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client",
		Operation: logical.ListOperation,
		Storage:   storage,
	})
	expectSuccess(t, respListClientAfterDelete, listErrAfterDelete)

	// validate list response
	delete(expectedStrings, "test-client2")
	expectStrings(t, respListClientAfterDelete.Data["keys"].([]string), expectedStrings)
}

// TestOIDC_Path_OIDC_ProviderScope_ReservedName tests that the reserved name
// "openid" cannot be used when creating a scope
func TestOIDC_Path_OIDC_ProviderScope_ReservedName(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test scope "test-scope" -- should succeed
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/openid",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})
	expectError(t, resp, err)
	// validate error message
	expectedStrings := map[string]interface{}{
		"the \"openid\" scope name is reserved": true,
	}
	expectStrings(t, []string{resp.Data["error"].(string)}, expectedStrings)
}

// TestOIDC_Path_OIDC_ProviderScope tests CRUD operations for scopes
func TestOIDC_Path_OIDC_ProviderScope(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test scope "test-scope" -- should succeed
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-scope" and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected := map[string]interface{}{
		"template":    "",
		"description": "",
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	// Update "test-scope" -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"template":    "eyAiZ3JvdXBzIjoge3tpZGVudGl0eS5lbnRpdHkuZ3JvdXBzLm5hbWVzfX0gfQ==",
			"description": "my-description",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-scope" again and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected = map[string]interface{}{
		"template":    "{ \"groups\": {{identity.entity.groups.names}} }",
		"description": "my-description",
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	// Delete test-scope -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope",
		Operation: logical.DeleteOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-scope" again and validate
	resp, _ = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	if resp != nil {
		t.Fatalf("expected nil but got resp: %#v", resp)
	}
}

// TestOIDC_Path_OIDC_ProviderScope_Update tests Update operations for scopes
func TestOIDC_Path_OIDC_ProviderScope_Update(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test scope "test-scope" -- should succeed
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"template":    "eyAiZ3JvdXBzIjoge3tpZGVudGl0eS5lbnRpdHkuZ3JvdXBzLm5hbWVzfX0gfQ==",
			"description": "my-description",
		},
	})
	expectSuccess(t, resp, err)

	// Read "test-scope" and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected := map[string]interface{}{
		"template":    "{ \"groups\": {{identity.entity.groups.names}} }",
		"description": "my-description",
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	// Update "test-scope" -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"template":    "eyAiZ3JvdXBzIjoge3tpZGVudGl0eS5lbnRpdHkuZ3JvdXBzLm5hbWVzfX0gfQ==",
			"description": "my-description-2",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-scope" again and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected = map[string]interface{}{
		"template":    "{ \"groups\": {{identity.entity.groups.names}} }",
		"description": "my-description-2",
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}
}

// TestOIDC_Path_OIDC_ProviderScope_List tests the List operation for scopes
func TestOIDC_Path_OIDC_ProviderScope_List(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Prepare two scopes, test-scope1 and test-scope2
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope1",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})

	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope2",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})

	// list scopes
	respListScopes, listErr := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope",
		Operation: logical.ListOperation,
		Storage:   storage,
	})
	expectSuccess(t, respListScopes, listErr)

	// validate list response
	expectedStrings := map[string]interface{}{"test-scope1": true, "test-scope2": true}
	expectStrings(t, respListScopes.Data["keys"].([]string), expectedStrings)

	// delete test-scope2
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope2",
		Operation: logical.DeleteOperation,
		Storage:   storage,
	})

	// list scopes again and validate response
	respListScopeAfterDelete, listErrAfterDelete := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope",
		Operation: logical.ListOperation,
		Storage:   storage,
	})
	expectSuccess(t, respListScopeAfterDelete, listErrAfterDelete)

	// validate list response
	delete(expectedStrings, "test-scope2")
	expectStrings(t, respListScopeAfterDelete.Data["keys"].([]string), expectedStrings)
}

// TestOIDC_Path_OIDC_ProviderAssignment tests CRUD operations for assignments
func TestOIDC_Path_OIDC_ProviderAssignment(t *testing.T) {
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

	// Read "test-assignment" again and validate
	resp, _ = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/test-assignment",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	if resp != nil {
		t.Fatalf("expected nil but got resp: %#v", resp)
	}
}

// TestOIDC_Path_OIDC_ProviderAssignment_DeleteWithExistingClient tests that an
// assignment cannot be deleted when it is referenced by a client
func TestOIDC_Path_OIDC_ProviderAssignment_DeleteWithExistingClient(t *testing.T) {
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

	// Create a test key "test-key"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"verification_ttl": "2m",
			"rotation_period":  "2m",
		},
		Storage: storage,
	})

	// Create a test client "test-client" -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"key":         "test-key",
			"assignments": []string{"test-assignment"},
		},
	})
	expectSuccess(t, resp, err)

	// Delete test-assignment -- should fail
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/test-assignment",
		Operation: logical.DeleteOperation,
		Storage:   storage,
	})
	expectError(t, resp, err)
	// validate error message
	expectedStrings := map[string]interface{}{
		"unable to delete assignment \"test-assignment\" because it is currently referenced by these clients: test-client": true,
	}
	expectStrings(t, []string{resp.Data["error"].(string)}, expectedStrings)

	// Read "test-assignment" again and validate
	resp, _ = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/test-assignment",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	if resp != nil {
		t.Fatalf("expected nil but got resp: %#v", resp)
	}
}

// TestOIDC_Path_OIDC_ProviderAssignment_Update tests Update operations for assignments
func TestOIDC_Path_OIDC_ProviderAssignment_Update(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test assignment "test-assignment" -- should succeed
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/test-assignment",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"groups":   "my-group",
			"entities": "my-entity",
		},
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
		"groups":   []string{"my-group"},
		"entities": []string{"my-entity"},
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	// Update "test-assignment" -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/test-assignment",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"groups": "my-group2",
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
		"groups":   []string{"my-group2"},
		"entities": []string{"my-entity"},
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}
}

// TestOIDC_Path_OIDC_ProviderAssignment_List tests the List operation for assignments
func TestOIDC_Path_OIDC_ProviderAssignment_List(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Prepare two assignments, test-assignment1 and test-assignment2
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/test-assignment1",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})

	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/test-assignment2",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})

	// list assignments
	respListAssignments, listErr := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment",
		Operation: logical.ListOperation,
		Storage:   storage,
	})
	expectSuccess(t, respListAssignments, listErr)

	// validate list response
	expectedStrings := map[string]interface{}{"test-assignment1": true, "test-assignment2": true}
	expectStrings(t, respListAssignments.Data["keys"].([]string), expectedStrings)

	// delete test-assignment2
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/test-assignment2",
		Operation: logical.DeleteOperation,
		Storage:   storage,
	})

	// list assignments again and validate response
	respListAssignmentAfterDelete, listErrAfterDelete := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment",
		Operation: logical.ListOperation,
		Storage:   storage,
	})
	expectSuccess(t, respListAssignmentAfterDelete, listErrAfterDelete)

	// validate list response
	delete(expectedStrings, "test-assignment2")
	expectStrings(t, respListAssignmentAfterDelete.Data["keys"].([]string), expectedStrings)
}
