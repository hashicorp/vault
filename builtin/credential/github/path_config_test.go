package github

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
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

// setupTestServer configures httptest server to intercept and respond to the ...
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

		w.Header().Add("Content-Type", "applicaion/json")
		fmt.Fprintln(w, resp)
	}))
}

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

	var expectedOrgID int64 = 12345
	// the ID should be set, we grab it from the GET /orgs API
	assert.Equal(t, expectedOrgID, resp.Data["organization_id"])
	assert.Equal(t, "foo-org", resp.Data["organization"])
}

// many of the fields have been omitted
// https://docs.github.com/en/rest/reference/users#get-the-authenticated-user
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

// many of the fields have been omitted, we only care about 'login' and 'id'
// https://docs.github.com/en/rest/reference/orgs#get-an-organization
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
