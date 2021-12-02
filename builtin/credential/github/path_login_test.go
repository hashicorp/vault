package github

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
)

func TestGitHub_Login(t *testing.T) {
	b, s := createBackendWithStorage(t)

	// use a test server to return our mock GH org info
	ts := setupTestServer(t)
	defer ts.Close()

	// Write the config
	b.HandleRequest(context.Background(), &logical.Request{
		Path:      "config",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"organization": "foo-org",
			"base_url":     ts.URL, // base_url will call the test server
		},
		Storage: s,
	})

	// Read the config
	b.HandleRequest(context.Background(), &logical.Request{
		Path:      "config",
		Operation: logical.ReadOperation,
		Storage:   s,
	})

	// Read the config
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Path:      "login",
		Operation: logical.UpdateOperation,
		Storage:   s,
	})
	assert.NoError(t, err)
	assert.NoError(t, resp.Error())
	expectedMetaData := map[string]string{
		"org":      "foo-org",
		"username": "user-foo",
	}
	assert.Equal(t, expectedMetaData, resp.Auth.Metadata)
}
