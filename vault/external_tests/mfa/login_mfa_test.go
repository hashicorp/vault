package mfa

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

// TestLoginMFA_LoginEnforcement_UniqueNames tests to ensure that 2 different login enforcements can be created with
// the same name, as long as they're in separate namespaces.
func TestLoginMFA_LoginEnforcement_UniqueNames(t *testing.T) {
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)
	client := cluster.Cores[0].Client

	// create a few namespaces
	_, err := client.Logical().Write("sys/namespaces/foo", nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Logical().Write("sys/namespaces/bar", nil)
	if err != nil {
		t.Fatal(err)
	}

	// create some prereq data
	resp, err := client.Logical().Write("identity/mfa/method/totp", map[string]interface{}{
		"issuer":    "fooCorp",
		"period":    10,
		"algorithm": "SHA1",
		"digits":    6,
		"skew":      1,
		"key_size":  uint(10),
		"qr_size":   100,
	})
	if err != nil {
		t.Fatal(err)
	}
	fooConfigId := resp.Data["method_id"].(string)

	resp, err = client.Logical().Write("identity/mfa/method/totp", map[string]interface{}{
		"issuer":    "barCorp",
		"period":    10,
		"algorithm": "SHA1",
		"digits":    6,
		"skew":      1,
		"key_size":  uint(10),
		"qr_size":   100,
	})
	if err != nil {
		t.Fatal(err)
	}
	barConfigId := resp.Data["method_id"].(string)

	resp, err = client.Logical().Write("identity/entity", map[string]interface{}{"name": "bob"})
	if err != nil {
		t.Fatal(err)
	}
	bobId := resp.Data["id"].(string)

	resp, err = client.Logical().Write("identity/entity", map[string]interface{}{"name": "alice"})
	if err != nil {
		t.Fatal(err)
	}
	aliceId := resp.Data["id"].(string)

	myPath := "identity/mfa/login-enforcement/baz"

	// create a login enforcement config in the foo ns
	client.SetNamespace("foo")
	_, err = client.Logical().Write(myPath, map[string]interface{}{
		"mfa_method_ids":      []string{fooConfigId},
		"identity_entity_ids": []string{bobId},
	})
	if err != nil {
		t.Fatal(err)
	}

	// create the same login enforcement config with the same name in the bar ns.
	// this should succeed because enforcement config names are unique per ns,
	// not globally.
	client.SetNamespace("bar")
	_, err = client.Logical().Write(myPath, map[string]interface{}{
		"mfa_method_ids":      []string{barConfigId},
		"identity_entity_ids": []string{aliceId},
	})
	if err != nil {
		t.Fatal(err)
	}

	// when we read the foo login enforcement config back out, it should have fooCorp and bob, not barCorp and alice
	// because both baz login enforcements were stored separately, since they were in separate namespaces. if they
	// weren't stored separately, the second write would've overwritten the first.
	client.SetNamespace("foo")
	resp, err = client.Logical().Read(myPath)
	if err != nil {
		t.Fatal(err)
	}
	if ieid := resp.Data["identity_entity_ids"].([]interface{})[0]; ieid != bobId {
		t.Fatalf("expected bob but got %q", ieid)
	}
	if mmid := resp.Data["mfa_method_ids"].([]interface{})[0]; mmid != fooConfigId {
		t.Fatalf("expected %q but got %q", fooConfigId, mmid)
	}
}

// TestLoginMFA_Method_CRUD tests creating/reading/updating/deleting a method config for all of the MFA providers
func TestLoginMFA_Method_CRUD(t *testing.T) {
	cluster := vault.NewTestCluster(t, &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": userpass.Factory,
		},
	}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)
	client := cluster.Cores[0].Client

	// Enable userpass authentication
	err := client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatalf("failed to enable userpass auth: %v", err)
	}

	auths, err := client.Sys().ListAuth()
	if err != nil {
		t.Fatal(err)
	}
	mountAccessor := auths["userpass/"].Accessor

	testCases := []struct {
		methodName    string
		configData    map[string]interface{}
		keyToUpdate   string
		valueToUpdate string
		keyToCheck    string
		updatedValue  string
	}{
		{
			"totp",
			map[string]interface{}{
				"issuer":    "yCorp",
				"period":    10,
				"algorithm": "SHA1",
				"digits":    6,
				"skew":      1,
				"key_size":  uint(10),
				"qr_size":   100,
			},
			"issuer",
			"zCorp",
			"",
			"",
		},
		{
			"duo",
			map[string]interface{}{
				"mount_accessor":  mountAccessor,
				"secret_key":      "lol-secret",
				"integration_key": "integration-key",
				"api_hostname":    "some-hostname",
			},
			"api_hostname",
			"api-updated.duosecurity.com",
			"",
			"",
		},
		{
			"okta",
			map[string]interface{}{
				"mount_accessor": mountAccessor,
				"base_url":       "example.com",
				"org_name":       "my-org",
				"api_token":      "lol-token",
			},
			"org_name",
			"dev-62954466-updated",
			"",
			"",
		},
		{
			"pingid",
			map[string]interface{}{
				"mount_accessor":       mountAccessor,
				"settings_file_base64": "I0F1dG8tR2VuZXJhdGVkIGZyb20gUGluZ09uZSwgZG93bmxvYWRlZCBieSBpZD1bU1NPXSBlbWFpbD1baGFtaWRAaGFzaGljb3JwLmNvbV0KI1dlZCBEZWMgMTUgMTM6MDg6NDQgTVNUIDIwMjEKdXNlX2Jhc2U2NF9rZXk9YlhrdGMyVmpjbVYwTFd0bGVRPT0KdXNlX3NpZ25hdHVyZT10cnVlCnRva2VuPWxvbC10b2tlbgppZHBfdXJsPWh0dHBzOi8vaWRweG55bDNtLnBpbmdpZGVudGl0eS5jb20vcGluZ2lkCm9yZ19hbGlhcz1sb2wtb3JnLWFsaWFzCmFkbWluX3VybD1odHRwczovL2lkcHhueWwzbS5waW5naWRlbnRpdHkuY29tL3BpbmdpZAphdXRoZW50aWNhdG9yX3VybD1odHRwczovL2F1dGhlbnRpY2F0b3IucGluZ29uZS5jb20vcGluZ2lkL3BwbQ==",
			},
			"settings_file_base64",
			"I0F1dG8tR2VuZXJhdGVkIGZyb20gUGluZ09uZSwgZG93bmxvYWRlZCBieSBpZD1bU1NPXSBlbWFpbD1baGFtaWRAaGFzaGljb3JwLmNvbV0KI1dlZCBEZWMgMTUgMTM6MDg6NDQgTVNUIDIwMjEKdXNlX2Jhc2U2NF9rZXk9YlhrdGMyVmpjbVYwTFd0bGVRPT0KdXNlX3NpZ25hdHVyZT10cnVlCnRva2VuPWxvbC10b2tlbgppZHBfdXJsPWh0dHBzOi8vaWRweG55bDNtLnBpbmdpZGVudGl0eS5jb20vcGluZ2lkL3VwZGF0ZWQKb3JnX2FsaWFzPWxvbC1vcmctYWxpYXMKYWRtaW5fdXJsPWh0dHBzOi8vaWRweG55bDNtLnBpbmdpZGVudGl0eS5jb20vcGluZ2lkCmF1dGhlbnRpY2F0b3JfdXJsPWh0dHBzOi8vYXV0aGVudGljYXRvci5waW5nb25lLmNvbS9waW5naWQvcHBt",
			"idp_url",
			"https://idpxnyl3m.pingidentity.com/pingid/updated",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.methodName, func(t *testing.T) {
			// create a new method config
			myPath := fmt.Sprintf("identity/mfa/method/%s", tc.methodName)
			resp, err := client.Logical().Write(myPath, tc.configData)
			if err != nil {
				t.Fatal(err)
			}

			methodId := resp.Data["method_id"]
			if methodId == "" {
				t.Fatal("method id is empty")
			}

			myNewPath := fmt.Sprintf("%s/%s", myPath, methodId)

			// read it back
			resp, err = client.Logical().Read(myNewPath)
			if err != nil {
				t.Fatal(err)
			}

			if resp.Data["id"] != methodId {
				t.Fatal("expected response id to match existing method id but it didn't")
			}

			// listing should show it
			resp, err = client.Logical().List(myPath)
			if err != nil {
				t.Fatal(err)
			}
			if resp.Data["keys"].([]interface{})[0] != methodId {
				t.Fatalf("expected %q in the list of method ids but it wasn't there", methodId)
			}

			// update it
			tc.configData[tc.keyToUpdate] = tc.valueToUpdate
			_, err = client.Logical().Write(myNewPath, tc.configData)
			if err != nil {
				t.Fatal(err)
			}

			resp, err = client.Logical().Read(myNewPath)
			if err != nil {
				t.Fatal(err)
			}

			// these shenanigans are to work around the arcane way that pingid does updates
			if tc.keyToCheck != "" && tc.updatedValue != "" {
				if resp.Data[tc.keyToCheck] != tc.updatedValue {
					t.Fatalf("expected config to update but it didn't: %v != %v", resp.Data[tc.keyToCheck], tc.updatedValue)
				}
			} else {
				if resp.Data[tc.keyToUpdate] != tc.valueToUpdate {
					t.Fatalf("expected config to update but it didn't: %v != %v", resp.Data[tc.keyToUpdate], tc.valueToUpdate)
				}
			}

			// delete it
			_, err = client.Logical().Delete(myNewPath)
			if err != nil {
				t.Fatal(err)
			}

			// try to read it again - should 404
			resp, err = client.Logical().Read(myNewPath)
			if !(resp == nil && err == nil) {
				t.Fatal("expected a 404 but didn't get one")
			}
		})
	}
}

// TestLoginMFA_LoginEnforcement_CRUD tests creating/reading/updating/deleting a login enforcement config
func TestLoginMFA_LoginEnforcement_CRUD(t *testing.T) {
	cluster := vault.NewTestCluster(t, &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": userpass.Factory,
		},
	}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)
	client := cluster.Cores[0].Client

	// first create a few configs
	configIDs := make([]string, 0)

	for i := 0; i < 2; i++ {
		resp, err := client.Logical().Write("identity/mfa/method/totp", map[string]interface{}{
			"issuer":    fmt.Sprintf("fooCorp%d", i),
			"period":    10,
			"algorithm": "SHA1",
			"digits":    6,
			"skew":      1,
			"key_size":  uint(10),
			"qr_size":   100 + i,
		})
		if err != nil {
			t.Fatal(err)
		}

		configIDs = append(configIDs, resp.Data["method_id"].(string))
	}

	// enable userpass auth
	err := client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatal(err)
	}

	auths, err := client.Sys().ListAuth()
	if err != nil {
		t.Fatal(err)
	}

	var mountAccessor string
	if auths != nil && auths["userpass/"] != nil {
		mountAccessor = auths["userpass/"].Accessor
	}

	// create a few entities
	resp, err := client.Logical().Write("identity/entity", map[string]interface{}{"name": "bob"})
	if err != nil {
		t.Fatal(err)
	}
	bobId := resp.Data["id"].(string)
	resp, err = client.Logical().Write("identity/entity", map[string]interface{}{"name": "alice"})
	if err != nil {
		t.Fatal(err)
	}
	aliceId := resp.Data["id"].(string)

	// create a few groups
	resp, err = client.Logical().Write("identity/group", map[string]interface{}{
		"metadata":          map[string]interface{}{"rad": true},
		"member_entity_ids": []string{aliceId},
	})
	if err != nil {
		t.Fatal(err)
	}
	radGroupId := resp.Data["id"].(string)

	resp, err = client.Logical().Write("identity/group", map[string]interface{}{
		"metadata":          map[string]interface{}{"sad": true},
		"member_entity_ids": []string{bobId},
	})
	if err != nil {
		t.Fatal(err)
	}
	sadGroupId := resp.Data["id"].(string)

	myPath := "identity/mfa/login-enforcement/foo"
	data := map[string]interface{}{
		"mfa_method_ids":        []string{configIDs[0], configIDs[1]},
		"auth_method_accessors": []string{mountAccessor},
	}

	// create a login enforcement config
	_, err = client.Logical().Write(myPath, data)
	if err != nil {
		t.Fatal(err)
	}

	// read it back
	resp, err = client.Logical().Read(myPath)
	if err != nil {
		t.Fatal(err)
	}

	equal := strutil.EquivalentSlices(data["mfa_method_ids"].([]string), stringSliceFromInterfaceSlice(resp.Data["mfa_method_ids"].([]interface{})))
	if !equal {
		t.Fatal("expected input mfa method ids to equal output mfa method ids")
	}
	equal = strutil.EquivalentSlices(data["auth_method_accessors"].([]string), stringSliceFromInterfaceSlice(resp.Data["auth_method_accessors"].([]interface{})))
	if !equal {
		t.Fatal("expected input auth method accessors to equal output auth method accessors")
	}

	// update it
	data["identity_group_ids"] = []string{radGroupId, sadGroupId}
	data["identity_entity_ids"] = []string{bobId, aliceId}
	_, err = client.Logical().Write(myPath, data)
	if err != nil {
		t.Fatal(err)
	}

	// read it back
	resp, err = client.Logical().Read(myPath)
	if err != nil {
		t.Fatal(err)
	}

	equal = strutil.EquivalentSlices(data["identity_group_ids"].([]string), stringSliceFromInterfaceSlice(resp.Data["identity_group_ids"].([]interface{})))
	if !equal {
		t.Fatal("expected input identity group ids to equal output identity group ids")
	}
	equal = strutil.EquivalentSlices(data["identity_entity_ids"].([]string), stringSliceFromInterfaceSlice(resp.Data["identity_entity_ids"].([]interface{})))
	if !equal {
		t.Fatal("expected input identity entity ids to equal output identity entity ids")
	}

	// delete it
	_, err = client.Logical().Delete(myPath)
	if err != nil {
		t.Fatal(err)
	}

	// try to read it back again - should 404
	resp, err = client.Logical().Read(myPath)

	// when both the response and the error are nil on a read request, that gets translated into a 404
	if !(resp == nil && err == nil) {
		t.Fatal("expected the read to 404 but it didn't")
	}
}

// TestLoginMFA_LoginEnforcement_MethodIdsIsRequired ensures that login enforcements have method ids attached
func TestLoginMFA_LoginEnforcement_MethodIdsIsRequired(t *testing.T) {
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)
	client := cluster.Cores[0].Client

	// create a login enforcement config, which should fail
	_, err := client.Logical().Write("identity/mfa/login-enforcement/foo", map[string]interface{}{})
	if err == nil {
		t.Fatal("expected an error but didn't get one")
	}

	if !strings.Contains(err.Error(), "missing method ids") {
		t.Fatal("should have received an error about missing method ids but didn't")
	}
}

// TestLoginMFA_LoginEnforcement_RequiredParameters validates that all of the required fields must be present
func TestLoginMFA_LoginEnforcement_RequiredParameters(t *testing.T) {
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)
	client := cluster.Cores[0].Client

	// first create a few configs
	configIDs := make([]string, 0)

	for i := 0; i < 2; i++ {
		resp, err := client.Logical().Write("identity/mfa/method/totp", map[string]interface{}{
			"issuer":    fmt.Sprintf("fooCorp%d", i),
			"period":    10,
			"algorithm": "SHA1",
			"digits":    6,
			"skew":      1,
			"key_size":  uint(10),
			"qr_size":   100 + i,
		})
		if err != nil {
			t.Fatal(err)
		}

		configIDs = append(configIDs, resp.Data["method_id"].(string))
	}

	// create a login enforcement config, which should fail
	_, err := client.Logical().Write("identity/mfa/login-enforcement/foo", map[string]interface{}{
		"mfa_method_ids": []string{configIDs[0], configIDs[1]},
	})
	if err == nil {
		t.Fatal("expected an error but didn't get one")
	}
	if !strings.Contains(err.Error(), "One of auth_method_accessors, auth_method_types, identity_group_ids, identity_entity_ids must be specified") {
		t.Fatal("expected an error about required fields but didn't get one")
	}
}

// TestLoginMFA_Method_Namespaces tests to ensure that namespace rules are followed when operating on method configs
func TestLoginMFA_Method_Namespaces(t *testing.T) {
	cluster := vault.NewTestCluster(t, &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": userpass.Factory,
		},
	}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)
	client := cluster.Cores[0].Client

	// Enable userpass authentication
	err := client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatalf("failed to enable userpass auth: %v", err)
	}

	auths, err := client.Sys().ListAuth()
	if err != nil {
		t.Fatal(err)
	}
	mountAccessor := auths["userpass/"].Accessor

	// create a few namespaces - foo, foo/bar, foo/bar/baz, quux
	_, err = client.Logical().Write("sys/namespaces/foo", nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Logical().Write("sys/namespaces/quux", nil)
	if err != nil {
		t.Fatal(err)
	}
	client.SetNamespace("foo")
	_, err = client.Logical().Write("sys/namespaces/bar", nil)
	if err != nil {
		t.Fatal(err)
	}
	client.SetNamespace("foo/bar")
	_, err = client.Logical().Write("sys/namespaces/baz", nil)
	if err != nil {
		t.Fatal(err)
	}

	// create a method config in ns foo/bar
	resp, err := client.Logical().Write("identity/mfa/method/totp", map[string]interface{}{
		"issuer":    "yCorp",
		"period":    10,
		"algorithm": "SHA1",
		"digits":    6,
		"skew":      1,
		"key_size":  uint(10),
		"qr_size":   100,
	})
	if err != nil {
		t.Fatal(err)
	}
	fooBarMethodId := resp.Data["method_id"].(string)
	fooBarPath := fmt.Sprintf("identity/mfa/method/totp/%s", fooBarMethodId)

	// create 2 additional method configs in ns foo
	client.SetNamespace("foo")
	resp, err = client.Logical().Write("identity/mfa/method/duo", map[string]interface{}{
		"mount_accessor":  mountAccessor,
		"secret_key":      "oIiQkWhGZw3r5gV1cRSUQ9dwiUv4atW4vdTCx2v9",
		"integration_key": "DI6XBJ2S2VEDGW8KZ2BH",
		"api_hostname":    "api-52ae179c.duosecurity.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err = client.Logical().Write("identity/mfa/method/totp", map[string]interface{}{
		"issuer":    "aCorp",
		"period":    10,
		"algorithm": "SHA1",
		"digits":    6,
		"skew":      1,
		"key_size":  uint(10),
		"qr_size":   100,
	})
	if err != nil {
		t.Fatal(err)
	}
	fooMethodId := resp.Data["method_id"].(string)

	// create another method config in the root ns
	client.ClearNamespace()
	resp, err = client.Logical().Write("identity/mfa/method/totp", map[string]interface{}{
		"issuer":    "zCorp",
		"period":    10,
		"algorithm": "SHA1",
		"digits":    6,
		"skew":      1,
		"key_size":  uint(10),
		"qr_size":   100,
	})
	if err != nil {
		t.Fatal(err)
	}
	rootMethodId := resp.Data["method_id"].(string)
	rootPath := fmt.Sprintf("identity/mfa/method/totp/%s", rootMethodId)

	successCallback := func(r *api.Secret, e error) {
		if e != nil {
			t.Fatal(e)
		}
		if r != nil && r.Data["error"] != nil {
			t.Fatal(r.Data["error"])
		}
	}

	failureCallback := func(r *api.Secret, e error) {
		if e == nil {
			t.Fatal("expected to get an error but didn't get one")
		}
	}

	testCases := []struct {
		name      string
		action    string
		namespace string
		path      string
		data      map[string]interface{}
		callback  func(*api.Secret, error)
	}{
		{
			"read foo/bar from foo/bar",
			"read",
			"foo/bar",
			fooBarPath,
			nil,
			successCallback,
		},
		{
			"update foo/bar from foo/bar",
			"update",
			"foo/bar",
			fooBarPath,
			map[string]interface{}{"issuer": "lolCorp"},
			successCallback,
		},
		{
			"read foo/bar from root",
			"read",
			"",
			fooBarPath,
			nil,
			successCallback,
		},
		{
			"update foo/bar from root",
			"update",
			"",
			fooBarPath,
			map[string]interface{}{"issuer": "lolCorp"},
			failureCallback,
		},
		{
			"read foo/bar from quux",
			"read",
			"quux",
			fooBarPath,
			nil,
			failureCallback,
		},
		{
			"update foo/bar from quux",
			"update",
			"quux",
			fooBarPath,
			map[string]interface{}{"issuer": "lolCorp"},
			failureCallback,
		},
		{
			"read foo/bar from foo/bar/baz",
			"read",
			"foo/bar/baz",
			fooBarPath,
			nil,
			successCallback,
		},
		{
			"update foo/bar from foo/bar/baz",
			"update",
			"foo/bar/baz",
			fooBarPath,
			map[string]interface{}{"issuer": "lolCorp"},
			failureCallback,
		},
		{
			"read foo/bar from foo",
			"read",
			"foo",
			fooBarPath,
			nil,
			successCallback,
		},
		{
			"update foo/bar from foo",
			"update",
			"foo",
			fooBarPath,
			map[string]interface{}{"issuer": "lolCorp"},
			failureCallback,
		},
		{
			"read root from root",
			"read",
			"",
			rootPath,
			nil,
			successCallback,
		},
		{
			"update root from root",
			"update",
			"",
			rootPath,
			map[string]interface{}{"issuer": "lolCorp"},
			successCallback,
		},
		{
			"read root from foo",
			"read",
			"foo",
			rootPath,
			nil,
			successCallback,
		},
		{
			"update root from foo",
			"update",
			"foo",
			rootPath,
			map[string]interface{}{"issuer": "lolCorp"},
			failureCallback,
		},
		{
			"list foo/bar from foo/bar",
			"list",
			"foo/bar",
			"identity/mfa/method/totp",
			nil,
			func(s *api.Secret, e error) {
				if e != nil {
					t.Fatal(e)
				}
				if s != nil && s.Data["error"] != nil {
					t.Fatal(s.Data["error"])
				}

				// we should get 3 results back when listing foo/bar from foo/bar:
				// one from foo/bar itself, one from foo, and one from root.
				// note that there are 2 method configs defined in foo, 1 in foo/bar, 1 in root, so 4 total,
				// but foo has one totp and one duo. we're listing totp here, so we should not get
				// the duo one back.
				if k := len(s.Data["keys"].([]interface{})); k != 3 {
					t.Fatalf("expected 3 keys but got %d", k)
				}
				expectedKeys := []string{fooBarMethodId, fooMethodId, rootMethodId}
				actualKeys := stringSliceFromInterfaceSlice(s.Data["keys"].([]interface{}))

				if !strutil.EquivalentSlices(actualKeys, expectedKeys) {
					t.Fatalf("expected %v to be equivalent to %v but it wasn't", actualKeys, expectedKeys)
				}
			},
		},
	}

	for _, testCase := range testCases {
		name := fmt.Sprintf("%s %s", testCase.action, testCase.name)
		t.Run(name, func(t *testing.T) {
			if testCase.namespace == "" {
				client.ClearNamespace()
			} else {
				client.SetNamespace(testCase.namespace)
			}

			var err error
			var resp *api.Secret

			switch testCase.action {
			case "read":
				resp, err = client.Logical().Read(testCase.path)
			case "update":
				resp, err = client.Logical().Write(testCase.path, testCase.data)
			case "list":
				resp, err = client.Logical().List(testCase.path)
			}

			testCase.callback(resp, err)
		})
	}
}

func TestLoginMFA_LoginEnforcement_Namespaces(t *testing.T) {
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)
	client := cluster.Cores[0].Client

	// create a few namespaces - foo, foo/bar, foo/bar/baz, quux
	_, err := client.Logical().Write("sys/namespaces/foo", nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Logical().Write("sys/namespaces/quux", nil)
	if err != nil {
		t.Fatal(err)
	}
	client.SetNamespace("foo")
	_, err = client.Logical().Write("sys/namespaces/bar", nil)
	if err != nil {
		t.Fatal(err)
	}
	client.SetNamespace("foo/bar")
	_, err = client.Logical().Write("sys/namespaces/baz", nil)
	if err != nil {
		t.Fatal(err)
	}

	// create a method config in ns foo/bar
	resp, err := client.Logical().Write("identity/mfa/method/totp", map[string]interface{}{
		"issuer":    "yCorp",
		"period":    10,
		"algorithm": "SHA1",
		"digits":    6,
		"skew":      1,
		"key_size":  uint(10),
		"qr_size":   100,
	})
	if err != nil {
		t.Fatal(err)
	}
	fooBarMethodId := resp.Data["method_id"].(string)

	// create an entity in ns foo/bar
	resp, err = client.Logical().Write("identity/entity", map[string]interface{}{"name": "alice"})
	if err != nil {
		t.Fatal(err)
	}
	aliceId := resp.Data["id"].(string)

	// create a login enforcement config in ns foo/bar
	data := map[string]interface{}{
		"mfa_method_ids":      []string{fooBarMethodId},
		"identity_entity_ids": []string{aliceId},
	}

	fooBarPath := "identity/mfa/login-enforcement/lol"
	resp, err = client.Logical().Write(fooBarPath, data)
	if err != nil {
		t.Fatal(err)
	}

	// create a method config in the root ns
	client.ClearNamespace()
	resp, err = client.Logical().Write("identity/mfa/method/totp", map[string]interface{}{
		"issuer":    "zCorp",
		"period":    10,
		"algorithm": "SHA1",
		"digits":    6,
		"skew":      1,
		"key_size":  uint(10),
		"qr_size":   100,
	})
	if err != nil {
		t.Fatal(err)
	}
	rootMethodId := resp.Data["method_id"].(string)

	// create an entity in the root ns
	resp, err = client.Logical().Write("identity/entity", map[string]interface{}{"name": "bob"})
	if err != nil {
		t.Fatal(err)
	}
	bobId := resp.Data["id"].(string)

	// create another login enforcement config in the root ns
	rootPath := "identity/mfa/login-enforcement/lawl"
	data = map[string]interface{}{
		"mfa_method_ids":      []string{rootMethodId},
		"identity_entity_ids": []string{bobId},
	}
	resp, err = client.Logical().Write(rootPath, data)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name      string
		action    string
		namespace string
		path      string
		succeed   bool
		data      map[string]interface{}
	}{
		{
			"read foo/bar from foo/bar",
			"read",
			"foo/bar",
			fooBarPath,
			true,
			nil,
		},
		{
			"update foo/bar from foo/bar",
			"update",
			"foo/bar",
			fooBarPath,
			true,
			map[string]interface{}{
				"mfa_method_ids":      []string{fooBarMethodId},
				"identity_entity_ids": []string{aliceId},
			},
		},
		{
			"read foo/bar from root",
			"read",
			"",
			fooBarPath,
			true,
			nil,
		},
		{
			"update foo/bar from root",
			"update",
			"",
			fooBarPath,
			true,
			map[string]interface{}{
				"mfa_method_ids":      []string{rootMethodId},
				"identity_entity_ids": []string{bobId},
			},
		},
		{
			"read foo/bar from quux",
			"read",
			"quux",
			fooBarPath,
			false,
			nil,
		},
		{
			"update foo/bar from quux",
			"update",
			"quux",
			fooBarPath,
			false,
			map[string]interface{}{
				"mfa_method_ids":      []string{fooBarMethodId},
				"identity_entity_ids": []string{aliceId},
			},
		},
		{
			"read foo/bar from foo/bar/baz",
			"read",
			"foo/bar/baz",
			fooBarPath,
			false,
			nil,
		},
		{
			"update foo/bar from foo/bar/baz",
			"update",
			"foo/bar/baz",
			fooBarPath,
			false,
			map[string]interface{}{
				"mfa_method_ids":      []string{fooBarMethodId},
				"identity_entity_ids": []string{aliceId},
			},
		},
		{
			"read foo/bar from foo",
			"read",
			"foo",
			fooBarPath,
			true,
			nil,
		},
		{
			"update foo/bar from foo",
			"update",
			"foo",
			fooBarPath,
			false,
			map[string]interface{}{
				"mfa_method_ids":      []string{fooBarMethodId},
				"identity_entity_ids": []string{aliceId},
			},
		},
		{
			"read root from root",
			"read",
			"",
			rootPath,
			true,
			nil,
		},
		{
			"update root from root",
			"update",
			"",
			rootPath,
			true,
			map[string]interface{}{
				"mfa_method_ids":      []string{rootMethodId},
				"identity_entity_ids": []string{bobId},
			},
		},
		{
			"read root from foo",
			"read",
			"foo",
			rootPath,
			false,
			nil,
		},
		{
			"update root from foo",
			"update",
			"foo",
			rootPath,
			false,
			map[string]interface{}{
				"mfa_method_ids":      []string{rootMethodId},
				"identity_entity_ids": []string{bobId},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.namespace == "" {
				client.ClearNamespace()
			} else {
				client.SetNamespace(testCase.namespace)
			}

			var err error
			var resp *api.Secret

			switch testCase.action {
			case "read":
				resp, err = client.Logical().Read(testCase.path)
			case "update":
				resp, err = client.Logical().Write(testCase.path, testCase.data)
			}

			if testCase.succeed {
				if err != nil {
					t.Fatal(err)
				}
				if resp != nil && resp.Data["error"] != nil {
					t.Fatal(resp.Data["error"])
				}
			} else {
				if err == nil && resp != nil {
					t.Fatal("expected to get an error but didn't get one")
				}
			}
		})
	}
}

// TestLoginMFA_LoginEnforcement_ConfigNamespaces tests that a login enforcement config should be able to access method
// ids configured in its own namespace or any of its ancestor namespaces.
func TestLoginMFA_LoginEnforcement_ConfigNamespaces(t *testing.T) {
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)
	client := cluster.Cores[0].Client

	// create a few namespaces - foo, foo/bar
	_, err := client.Logical().Write("sys/namespaces/foo", nil)
	if err != nil {
		t.Fatal(err)
	}
	client.SetNamespace("foo")
	_, err = client.Logical().Write("sys/namespaces/bar", nil)
	if err != nil {
		t.Fatal(err)
	}

	// create a method config in the root ns
	client.ClearNamespace()
	resp, err := client.Logical().Write("identity/mfa/method/totp", map[string]interface{}{
		"issuer":    "rootCorp",
		"period":    10,
		"algorithm": "SHA1",
		"digits":    6,
		"skew":      1,
		"key_size":  uint(10),
		"qr_size":   100,
	})
	if err != nil {
		t.Fatal(err)
	}
	rootMethodId := resp.Data["method_id"].(string)

	// create a method config in ns foo
	client.SetNamespace("foo")
	resp, err = client.Logical().Write("identity/mfa/method/totp", map[string]interface{}{
		"issuer":    "fooCorp",
		"period":    10,
		"algorithm": "SHA1",
		"digits":    6,
		"skew":      1,
		"key_size":  uint(10),
		"qr_size":   100,
	})
	if err != nil {
		t.Fatal(err)
	}
	fooMethodId := resp.Data["method_id"].(string)

	// create a method config in ns foo/bar
	client.SetNamespace("foo/bar")
	resp, err = client.Logical().Write("identity/mfa/method/totp", map[string]interface{}{
		"issuer":    "fooBarCorp",
		"period":    10,
		"algorithm": "SHA1",
		"digits":    6,
		"skew":      1,
		"key_size":  uint(10),
		"qr_size":   100,
	})
	if err != nil {
		t.Fatal(err)
	}
	fooBarMethodId := resp.Data["method_id"].(string)

	// create an entity in ns foo/bar
	resp, err = client.Logical().Write("identity/entity", map[string]interface{}{"name": "bob"})
	if err != nil {
		t.Fatal(err)
	}
	bobId := resp.Data["id"].(string)

	// from the foo/bar ns, login enforcement configs should be able to reference any of the method configs
	// that were created, since they're all either in foo/bar or are an ancestor of foo/bar
	for _, id := range []string{rootMethodId, fooMethodId, fooBarMethodId} {
		data := map[string]interface{}{
			"mfa_method_ids":      []string{id},
			"identity_entity_ids": []string{bobId},
		}

		resp, err = client.Logical().Write("identity/mfa/login-enforcement/lol", data)
		if err != nil {
			t.Fatal(err)
		}
	}
}

// TestLoginMFA_LoginEnforcement_Validation tests that all of the parameters provided to a login enforcement config
// exist within Vault and aren't just random values.
func TestLoginMFA_LoginEnforcement_Validation(t *testing.T) {
	cluster := vault.NewTestCluster(t, &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": userpass.Factory,
		},
	}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)
	client := cluster.Cores[0].Client

	// create a few namespaces - foo, bar
	_, err := client.Logical().Write("sys/namespaces/foo", nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Logical().Write("sys/namespaces/bar", nil)
	if err != nil {
		t.Fatal(err)
	}

	// create a config in ns foo
	client.SetNamespace("foo")
	resp, err := client.Logical().Write("identity/mfa/method/totp", map[string]interface{}{
		"issuer":    "fooCorp",
		"period":    10,
		"algorithm": "SHA1",
		"digits":    6,
		"skew":      1,
		"key_size":  uint(10),
		"qr_size":   100,
	})
	if err != nil {
		t.Fatal(err)
	}
	fooConfigId := resp.Data["method_id"].(string)

	// enable userpass auth in ns foo
	err = client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatal(err)
	}

	auths, err := client.Sys().ListAuth()
	if err != nil {
		t.Fatal(err)
	}

	var mountAccessor string
	var mountType string
	if auths != nil && auths["userpass/"] != nil {
		mountAccessor = auths["userpass/"].Accessor
		mountType = auths["userpass/"].Type
	}

	// create an entity in ns foo
	resp, err = client.Logical().Write("identity/entity", map[string]interface{}{"name": "alice"})
	if err != nil {
		t.Fatal(err)
	}
	aliceId := resp.Data["id"].(string)

	// create a group in ns foo
	resp, err = client.Logical().Write("identity/group", map[string]interface{}{
		"metadata":          map[string]interface{}{"rad": true},
		"member_entity_ids": []string{aliceId},
	})
	if err != nil {
		t.Fatal(err)
	}
	radGroupId := resp.Data["id"].(string)

	// create an entity in ns bar
	client.SetNamespace("bar")
	resp, err = client.Logical().Write("identity/entity", map[string]interface{}{"name": "bob"})
	if err != nil {
		t.Fatal(err)
	}
	bobId := resp.Data["id"].(string)

	// create an entity in root ns
	client.ClearNamespace()
	resp, err = client.Logical().Write("identity/entity", map[string]interface{}{"name": "cynthia"})
	if err != nil {
		t.Fatal(err)
	}
	cynthiaId := resp.Data["id"].(string)

	myPath := "identity/mfa/login-enforcement/lol"

	// try to create a login enforcement config with a non-existant method id - should fail
	_, err = client.Logical().Write(myPath, map[string]interface{}{
		"mfa_method_ids":      []string{"wrong"},
		"identity_entity_ids": []string{cynthiaId},
	})
	if err == nil {
		t.Fatal("expected an error but didn't get one")
	}

	// try to create a login enforcement config using a method id from a different namespace that's not an ancestor
	// - should fail
	client.SetNamespace("bar")
	_, err = client.Logical().Write(myPath, map[string]interface{}{
		"mfa_method_ids":      []string{fooConfigId},
		"identity_entity_ids": []string{bobId},
	})
	if err == nil {
		t.Fatal("expected an error but didn't get one")
	}

	// try to create a login enforcement config with a group id for a non-existant group - should fail
	client.SetNamespace("foo")
	_, err = client.Logical().Write(myPath, map[string]interface{}{
		"mfa_method_ids":     [][]string{{fooConfigId}},
		"identity_group_ids": []string{"nope"},
	})
	if err == nil {
		t.Fatal("expected an error but didn't get one")
	}

	// try to create a login enforcement config with an entity id for a non-existant entity - should fail
	_, err = client.Logical().Write(myPath, map[string]interface{}{
		"mfa_method_ids":      []string{fooConfigId},
		"identity_entity_ids": []string{"nope"},
	})
	if err == nil {
		t.Fatal("expected an error but didn't get one")
	}

	// try to create a login enforcement config with a non-existant auth method accessor - should fail
	_, err = client.Logical().Write(myPath, map[string]interface{}{
		"mfa_method_ids":        []string{fooConfigId},
		"auth_method_accessors": []string{"wrong"},
	})
	if err == nil {
		t.Fatal("expected an error but didn't get one")
	}

	// try to create a login enforcement config with a non-existant auth method type - should fail
	_, err = client.Logical().Write(myPath, map[string]interface{}{
		"mfa_method_ids":    []string{fooConfigId},
		"auth_method_types": []string{"wrong"},
	})
	if err == nil {
		t.Fatal("expected an error but didn't get one")
	}

	// try to create a login enforcement config using a method id in the correct namespace with valid
	// data - should succeed
	data := map[string]interface{}{
		"mfa_method_ids":        []string{fooConfigId},
		"identity_group_ids":    []string{radGroupId},
		"identity_entity_ids":   []string{aliceId},
		"auth_method_accessors": []string{mountAccessor},
		"auth_method_types":     []string{mountType},
	}
	_, err = client.Logical().Write(myPath, data)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoginMFA_UpdateNonExistentConfig(t *testing.T) {
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)
	client := cluster.Cores[0].Client

	_, err := client.Logical().Write("mfa/method/totp/a51884c6-51f2-bdc3-f4c5-0da64fe4d061", map[string]interface{}{
		"issuer":    "yCorp",
		"period":    10,
		"algorithm": "SHA1",
		"digits":    6,
		"skew":      1,
		"key_size":  uint(10),
		"qr_size":   100,
	})
	if err == nil {
		t.Fatal("expected to get an error but didn't")
	}
	if !strings.Contains(err.Error(), "Code: 404") {
		t.Fatal("expected to get a 404 but didn't")
	}
}

// This is for converting []interface{} that you know holds all strings into []string
func stringSliceFromInterfaceSlice(input []interface{}) []string {
	result := make([]string, 0, len(input))
	for _, x := range input {
		if val, ok := x.(string); ok {
			result = append(result, val)
		}
	}
	return result
}
