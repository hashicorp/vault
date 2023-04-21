package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
)

func createBackendWithStorage(t *testing.T) (*backend, logical.Storage) {
	t.Helper()
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	b := Backend()
	if b == nil {
		t.Fatalf("failed to create backend")
	}
	err := b.Backend.Setup(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	return b, config.StorageView
}

// setupTestServer configures httptest server to intercept and respond to the
// request to base_url
func setupTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var resp string
		if strings.Contains(r.URL.String(), "/user/orgs") {
			resp = string(listOrgResponse)
		} else if strings.Contains(r.URL.String(), "/user/teams") {
			resp = string(listUserTeamsResponse)
		} else if strings.Contains(r.URL.String(), "/user") {
			resp = getUserResponse
		} else if strings.Contains(r.URL.String(), "/orgs/") {
			resp = getOrgResponse
		}

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintln(w, resp)
	}))
}

// TestGitHub_WriteReadConfig tests that we can successfully read and write
// the github auth config
func TestGitHub_WriteReadConfig(t *testing.T) {
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
	assert.Nil(t, resp)
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
	assert.Equal(t, "foo-org", resp.Data["organization"])
}

// TestGitHub_WriteReadConfig_OrgID tests that we can successfully read and
// write the github auth config with an organization_id param
func TestGitHub_WriteReadConfig_OrgID(t *testing.T) {
	b, s := createBackendWithStorage(t)

	// Write the config and pass in organization_id
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Path:      "config",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"organization":    "foo-org",
			"organization_id": 98765,
		},
		Storage: s,
	})
	assert.NoError(t, err)
	assert.Nil(t, resp)
	assert.NoError(t, resp.Error())

	// Read the config
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Path:      "config",
		Operation: logical.ReadOperation,
		Storage:   s,
	})
	assert.NoError(t, err)
	assert.NoError(t, resp.Error())

	// the ID should be set to what was written in the config
	assert.Equal(t, int64(98765), resp.Data["organization_id"])
	assert.Equal(t, "foo-org", resp.Data["organization"])
}

// TestGitHub_WriteReadConfig_Token tests that we can successfully read and
// write the github auth config with a token environment variable
func TestGitHub_WriteReadConfig_Token(t *testing.T) {
	b, s := createBackendWithStorage(t)
	// use a test server to return our mock GH org info
	ts := setupTestServer(t)
	defer ts.Close()

	err := os.Setenv("VAULT_AUTH_CONFIG_GITHUB_TOKEN", "foobar")
	assert.NoError(t, err)

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
	assert.Nil(t, resp)
	assert.NoError(t, resp.Error())

	// Read the config
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Path:      "config",
		Operation: logical.ReadOperation,
		Storage:   s,
	})
	assert.NoError(t, err)
	assert.NoError(t, resp.Error())

	// the token should not be returned in the read config response.
	assert.Nil(t, resp.Data["token"])
}

// TestGitHub_ErrorNoOrgID tests that an error is returned when we cannot fetch
// the org ID for the given org name
func TestGitHub_ErrorNoOrgID(t *testing.T) {
	b, s := createBackendWithStorage(t)
	// use a test server to return our mock GH org info
	ts := func() *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			resp := `{ "id": 0 }`
			fmt.Fprintln(w, resp)
		}))
	}

	defer ts().Close()

	// Write the config
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Path:      "config",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"organization": "foo-org",
			"base_url":     ts().URL, // base_url will call the test server
		},
		Storage: s,
	})
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, errors.New(
		"unable to fetch the organization_id, you must manually set it in the config: organization_id not found for foo-org",
	), err)
}

// TestGitHub_WriteConfig_ErrorNoOrg tests that an error is returned when the
// required "organization" parameter is not provided
func TestGitHub_WriteConfig_ErrorNoOrg(t *testing.T) {
	b, s := createBackendWithStorage(t)

	// Write the config
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Path:      "config",
		Operation: logical.UpdateOperation,
		Data:      map[string]interface{}{},
		Storage:   s,
	})

	assert.NoError(t, err)
	assert.Error(t, resp.Error())
	assert.Equal(t, errors.New("organization is a required parameter"), resp.Error())
}

// https://docs.github.com/en/rest/reference/users#get-the-authenticated-user
// Note: many of the fields have been omitted
var getUserResponse = `
{
	"login": "user-foo",
	"id": 6789,
	"description": "A great user. The very best user.",
	"name": "foo name",
	"company": "foo-company",
	"type": "User"
}
`

// https://docs.github.com/en/rest/reference/orgs#get-an-organization
// Note: many of the fields have been omitted, we only care about 'login' and 'id'
var getOrgResponse = `
{
	"login": "foo-org",
	"id": 12345,
	"description": "A great org. The very best org.",
	"name": "foo-display-name",
	"company": "foo-company",
	"type": "Organization"
}
`

// https://docs.github.com/en/rest/reference/orgs#list-organizations-for-the-authenticated-user
var listOrgResponse = []byte(fmt.Sprintf(`[%v]`, getOrgResponse))

// https://docs.github.com/en/rest/reference/teams#list-teams-for-the-authenticated-user
// Note: many of the fields have been omitted
var listUserTeamsResponse = []byte(fmt.Sprintf(`[
{
    "id": 1,
    "node_id": "MDQ6VGVhbTE=",
    "name": "Foo team",
    "slug": "foo-team",
    "description": "A great team. The very best team.",
    "permission": "admin",
    "organization": %v
  }
]`, getOrgResponse))
