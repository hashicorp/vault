package proxy

import (
	"net/http"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	credCert "github.com/hashicorp/vault/builtin/credential/cert"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func TestLogin(t *testing.T) {
	b := newTestBackend(t)

	header := "REMOTE_USER"
	req := createConfigRequest(map[string]interface{}{"user_header": header})
	b.AssertHandleRequest(req)

	roleName := "test1"
	requiredHeaders := map[string]string{
		"REQ_HDR_1": "REQ_VAL_1",
		"REQ_HDR_2": "REQ_VAL_2",
	}
	roleData := map[string]interface{}{
		"allowed_users":    "user1,user2,user3*x",
		"required_headers": requiredHeaders,
	}
	req = createRoleRequest(roleName, roleData)
	b.AssertHandleRequest(req)

	incHeaders := map[string][]string{
		"REQ_HDR_1": []string{"REQ_VAL_1"},
		"REQ_HDR_2": []string{"REQ_VAL_2"},
	}

	tests := []struct {
		role          string
		username      string
		userHeader    string
		headers       map[string][]string
		expectSuccess bool
		expectError   bool
	}{
		{roleName, "user1", header, incHeaders, true, false},
		{roleName, "user2", header, incHeaders, true, false},

		// user that matches glob => success
		{roleName, "user3foobarx", header, incHeaders, true, false},

		// user that does not glob => fail
		{roleName, "johndoe", header, incHeaders, false, false},
		{roleName, "user3foobar", header, incHeaders, false, false},

		// user doesn't match any existing role => fail
		{"unknownrole", "user1", header, incHeaders, false, false},

		// missing a required header => fail
		{roleName, "user1", header, map[string][]string{
			"REQ_HDR_1": []string{"REQ_VAL_1"},
		}, false, false},

		// required header has wrong value => fail
		{roleName, "user1", header, map[string][]string{
			"REQ_HDR_1": []string{"REQ_VAL_1"},
			"REQ_HDR_2": []string{"WRONG_VALUE"},
		}, false, false},

		// required header has two values => fail
		{roleName, "user1", header, map[string][]string{
			"REQ_HDR_1": []string{"REQ_VAL_1"},
			"REQ_HDR_2": []string{"REQ_VAL_2", "WRONG_VALUE"},
		}, false, false},
	}

	for idx, test := range tests {
		t.Logf("performing login test idx=%d, user=%s, role=%s", idx, test.username, test.role)
		req := loginRequest(test.username, test.role, test.userHeader, test.headers)

		if test.expectSuccess {
			b.AssertHandleRequest(req)
		} else {
			resp, err := b.HandleRequest(req)
			if test.expectError && err == nil {
				t.Fatal("unexpectedly got nil error")
			}

			if !test.expectSuccess && !resp.IsError() {
				t.Fatal("expected failed response but got success")
			}
		}
	}
}

// TestLoginWithCertClient verifies we can login to the proxy engine using the
// cert engines cli; allowing clients that lack proxy auth support to
// authenticate
func TestLoginWithCertClient(t *testing.T) {
	// start backend
	coreConfig := &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       log.NewNullLogger(),
		CredentialBackends: map[string]logical.Factory{
			"proxy": Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()
	cores := cluster.Cores
	vault.TestWaitActive(t, cores[0].Core)
	client := cores[0].Client

	proxyHeader := "remote_user"
	authOpts := &api.EnableAuthOptions{
		Type:        "proxy",
		Description: "enables entities that have authenticated to AB apache to obtain a token",
		Config: api.AuthConfigInput{
			DefaultLeaseTTL:           "24h",
			PassthroughRequestHeaders: []string{proxyHeader},
		},
	}

	// mount the proxy auth engine
	mountPoint := "proxy_auth/"
	if err := client.Sys().EnableAuthWithOptions(mountPoint, authOpts); err != nil {
		t.Fatalf("err mounting proxy auth engine: %+v", err)
	}

	if _, err := client.Logical().Write("auth/"+mountPoint+"config", map[string]interface{}{"user_header": proxyHeader}); err != nil {
		t.Fatalf("Error configuring proxy auth engine: %+v", err)
	}

	// setup a role
	role := "test_role"
	allowedUser := "user1"
	if _, err := client.Logical().Write("auth/"+mountPoint+"role/"+role, map[string]interface{}{"allowed_users": allowedUser}); err != nil {
		t.Fatalf("Error configuring proxy auth role: %+v", err)
	}

	// login to role with cert cli client
	headers := client.Headers()
	if headers == nil {
		headers = make(http.Header)
	}

	headers.Set(proxyHeader, allowedUser)
	client.SetHeaders(headers)
	cli := credCert.CLIHandler{}
	secret, err := cli.Auth(client, map[string]string{
		"name":  role,
		"mount": mountPoint,
	})
	if err != nil {
		t.Fatalf("login failed: %+v", err)
	}
	if secret.Auth == nil {
		t.Fatalf("login returned nil Auth")
	}
}
