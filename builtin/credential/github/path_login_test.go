package github

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
)

// TestGitHub_Login tests that we can successfully login with the given config
func TestGitHub_Login(t *testing.T) {
	b, s := createBackendWithStorage(t)

	// use a test server to return our mock GH org info
	ts := setupTestServer(t)
	defer ts.Close()

	// Write the config
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Path:      "config",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"organization": "foo-org",
			"base_url":     ts.URL, // base_url will call the test server
		},
		Storage: s,
	})
	assert.NoError(t, err)
	assert.NoError(t, resp.Error())

	// Read the config
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Path:      "config",
		Operation: logical.ReadOperation,
		Storage:   s,
	})
	assert.NoError(t, err)
	assert.NoError(t, resp.Error())

	// attempt a login
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Path:      "login",
		Operation: logical.UpdateOperation,
		Storage:   s,
	})

	expectedMetaData := map[string]string{
		"org":      "foo-org",
		"username": "user-foo",
	}
	assert.Equal(t, expectedMetaData, resp.Auth.Metadata)
	assert.NoError(t, err)
	assert.NoError(t, resp.Error())
}
