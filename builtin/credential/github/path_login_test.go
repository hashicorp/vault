package github

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/vault/helper/namespace"
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

// TestGitHub_Login_OrgInvalid tests that we cannot login with an ID other than
// what is set in the config
func TestGitHub_Login_OrgInvalid(t *testing.T) {
	b, s := createBackendWithStorage(t)
	ctx := namespace.RootContext(nil)

	// use a test server to return our mock GH org info
	ts := setupTestServer(t)
	defer ts.Close()

	// write and store config
	config := config{
		Organization:   "foo-org",
		OrganizationID: 9999,
		BaseURL:        ts.URL + "/", // base_url will call the test server
	}
	entry, err := logical.StorageEntryJSON("config", config)
	if err != nil {
		t.Fatalf("failed creating storage entry")
	}
	if err := s.Put(ctx, entry); err != nil {
		t.Fatalf("writing to in mem storage failed")
	}

	// attempt a login
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Path:      "login",
		Operation: logical.UpdateOperation,
		Storage:   s,
	})

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, errors.New("user is not part of required org"), err)
}

// TestGitHub_Login_OrgNameChanged tests that we can successfully login with the
// given config and emit a warning when the organization name has changed
func TestGitHub_Login_OrgNameChanged(t *testing.T) {
	b, s := createBackendWithStorage(t)
	ctx := namespace.RootContext(nil)

	// use a test server to return our mock GH org info
	ts := setupTestServer(t)
	defer ts.Close()

	// write and store config
	// the name does not match what the API will return but the ID does
	config := config{
		Organization:   "old-name",
		OrganizationID: 12345,
		BaseURL:        ts.URL + "/", // base_url will call the test server
	}
	entry, err := logical.StorageEntryJSON("config", config)
	if err != nil {
		t.Fatalf("failed creating storage entry")
	}
	if err := s.Put(ctx, entry); err != nil {
		t.Fatalf("writing to in mem storage failed")
	}

	// attempt a login
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Path:      "login",
		Operation: logical.UpdateOperation,
		Storage:   s,
	})

	assert.NoError(t, err)
	assert.Nil(t, resp.Error())
	assert.Equal(
		t,
		[]string{"the organization name has changed to \"foo-org\". It is recommended to verify and update the organization name in the config: organization_id=12345"},
		resp.Warnings,
	)
}

// TestGitHub_Login_NoOrgID tests that we can successfully login with the given
// config when no organization ID is present and write the fetched ID to the
// config
func TestGitHub_Login_NoOrgID(t *testing.T) {
	b, s := createBackendWithStorage(t)
	ctx := namespace.RootContext(nil)

	// use a test server to return our mock GH org info
	ts := setupTestServer(t)
	defer ts.Close()

	// write and store config without Org ID
	config := config{
		Organization: "foo-org",
		BaseURL:      ts.URL + "/", // base_url will call the test server
	}
	entry, err := logical.StorageEntryJSON("config", config)
	if err != nil {
		t.Fatalf("failed creating storage entry")
	}
	if err := s.Put(ctx, entry); err != nil {
		t.Fatalf("writing to in mem storage failed")
	}

	// attempt a login
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
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

	// Read the config
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Path:      "config",
		Operation: logical.ReadOperation,
		Storage:   s,
	})
	assert.NoError(t, err)
	assert.NoError(t, resp.Error())

	// the ID should be set, we grab it from the GET /orgs API
	assert.Equal(t, int64(12345), resp.Data["organization_id"])
}
