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
		invalidType   string
		configData    map[string]interface{}
		keyToUpdate   string
		valueToUpdate string
		keyToCheck    string
		updatedValue  string
	}{
		{
			"totp",
			"duo",
			map[string]interface{}{
				"issuer":                  "yCorp",
				"period":                  10,
				"algorithm":               "SHA1",
				"digits":                  6,
				"skew":                    1,
				"key_size":                uint(10),
				"qr_size":                 100,
				"max_validation_attempts": 1,
			},
			"issuer",
			"zCorp",
			"",
			"",
		},
		{
			"duo",
			"totp",
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
			"pingid",
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
			"okta",
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

			// read the id on another MFA type endpoint should fail
			invalidPath := fmt.Sprintf("identity/mfa/method/%s/%s", tc.invalidType, methodId)
			resp, err = client.Logical().Read(invalidPath)
			if err == nil {
				t.Fatal(err)
			}

			// read the id globally should succeed
			globalPath := fmt.Sprintf("identity/mfa/method/%s", methodId)
			resp, err = client.Logical().Read(globalPath)
			if err != nil {
				t.Fatal(err)
			}
			if resp.Data["id"] != methodId {
				t.Fatal("expected response id to match existing method id but it didn't")
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

// TestLoginMFA_ListAllMFAConfigs tests listing all configs globally
func TestLoginMFA_ListAllMFAConfigsGlobally(t *testing.T) {
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

	mfaConfigs := []struct {
		methodType string
		configData map[string]interface{}
	}{
		{
			"totp",
			map[string]interface{}{
				"issuer":                  "yCorp",
				"period":                  10,
				"algorithm":               "SHA1",
				"digits":                  6,
				"skew":                    1,
				"key_size":                uint(10),
				"qr_size":                 100,
				"max_validation_attempts": 1,
			},
		},
		{
			"duo",
			map[string]interface{}{
				"mount_accessor":  mountAccessor,
				"secret_key":      "lol-secret",
				"integration_key": "integration-key",
				"api_hostname":    "some-hostname",
			},
		},
		{
			"okta",
			map[string]interface{}{
				"mount_accessor": mountAccessor,
				"base_url":       "example.com",
				"org_name":       "my-org",
				"api_token":      "lol-token",
			},
		},
		{
			"pingid",
			map[string]interface{}{
				"mount_accessor":       mountAccessor,
				"settings_file_base64": "I0F1dG8tR2VuZXJhdGVkIGZyb20gUGluZ09uZSwgZG93bmxvYWRlZCBieSBpZD1bU1NPXSBlbWFpbD1baGFtaWRAaGFzaGljb3JwLmNvbV0KI1dlZCBEZWMgMTUgMTM6MDg6NDQgTVNUIDIwMjEKdXNlX2Jhc2U2NF9rZXk9YlhrdGMyVmpjbVYwTFd0bGVRPT0KdXNlX3NpZ25hdHVyZT10cnVlCnRva2VuPWxvbC10b2tlbgppZHBfdXJsPWh0dHBzOi8vaWRweG55bDNtLnBpbmdpZGVudGl0eS5jb20vcGluZ2lkCm9yZ19hbGlhcz1sb2wtb3JnLWFsaWFzCmFkbWluX3VybD1odHRwczovL2lkcHhueWwzbS5waW5naWRlbnRpdHkuY29tL3BpbmdpZAphdXRoZW50aWNhdG9yX3VybD1odHRwczovL2F1dGhlbnRpY2F0b3IucGluZ29uZS5jb20vcGluZ2lkL3BwbQ==",
			},
		},
	}

	var methodIDs []interface{}
	for _, method := range mfaConfigs {
		// create a new method config
		myPath := fmt.Sprintf("identity/mfa/method/%s", method.methodType)
		resp, err := client.Logical().Write(myPath, method.configData)
		if err != nil {
			t.Fatal(err)
		}

		methodId := resp.Data["method_id"]
		if methodId == "" {
			t.Fatal("method id is empty")
		}
		methodIDs = append(methodIDs, methodId)
	}
	// listing should show it
	resp, err := client.Logical().List("identity/mfa/method")
	if err != nil || resp == nil {
		t.Fatal(err)
	}

	if len(resp.Data["keys"].([]interface{})) != len(methodIDs) {
		t.Fatalf("global list request did not return all MFA method IDs")
	}
	if len(resp.Data["key_info"].(map[string]interface{})) != len(methodIDs) {
		t.Fatal("global list request did not return all MFA method configurations")
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

	// listing should show it
	resp, err = client.Logical().List("identity/mfa/login-enforcement")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Data["keys"].([]interface{})[0] != "foo" {
		t.Fatal("expected foo in the list of enforcement names but it wasn't there")
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
