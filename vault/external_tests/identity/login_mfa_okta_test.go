// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package identity

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
)

// To run these tests, set the following env variables:
// VAULT_ACC=1
// OKTA_ORG=dev-219337
// OKTA_API_TOKEN=<generate via web UI, see Confluence for login details>
// OKTA_USERNAME=test3@example.com
// OKTA_PASSWORD=<find in 1password>
//
// You will need to install the Okta client app on your mobile device and
// setup MFA in order to use the Okta web UI.  This test does not exercise
// MFA however (which is an enterprise feature), and therefore the test
// user in OKTA_USERNAME should not be configured with it.  Currently
// test3@example.com is not a member of testgroup, which is the group with
// the profile that requires MFA. If you need to use a different group name
// for the test group, you can set:
// OKTA_TEST_GROUP=alttestgroup

func TestOktaEngineMFA(t *testing.T) {
	if os.Getenv("VAULT_ACC") == "" {
		t.Skip("This test requires manual intervention and OKTA verify on cellphone is needed")
	}

	// Ensure each cred is populated.
	credNames := []string{
		"OKTA_ORG",
		"OKTA_API_TOKEN",
		"OKTA_USERNAME",
		"OKTA_PASSWORD",
	}
	testhelpers.SkipUnlessEnvVarsSet(t, credNames)

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	// Enable Okta engine
	err := client.Sys().EnableAuthWithOptions("okta", &api.EnableAuthOptions{
		Type: "okta",
	})
	if err != nil {
		t.Fatalf("failed to enable okta auth: %v", err)
	}

	_, err = client.Logical().Write("auth/okta/config", map[string]interface{}{
		"base_url":  "okta.com",
		"org_name":  os.Getenv("OKTA_ORG"),
		"api_token": os.Getenv("OKTA_API_TOKEN"),
	})
	if err != nil {
		t.Fatalf("error configuring okta mount: %v", err)
	}

	testGroup := os.Getenv("OKTA_TEST_GROUP")
	if len(testGroup) == 0 {
		testGroup = "testgroup"
	}

	_, err = client.Logical().Write("auth/okta/groups/"+testGroup, map[string]interface{}{
		"policies": "default",
	})
	if err != nil {
		t.Fatalf("error configuring okta group, %v", err)
	}

	_, err = client.Logical().Write("auth/okta/login/"+os.Getenv("OKTA_USERNAME"), map[string]interface{}{
		"password": os.Getenv("OKTA_PASSWORD"),
	})
	if err != nil {
		t.Fatalf("error configuring okta group, %v", err)
	}
}

func TestInteg_PolicyMFAOkta(t *testing.T) {
	if os.Getenv("VAULT_ACC") == "" {
		t.Skip("This test requires manual intervention and OKTA verify on cellphone is needed")
	}

	// Ensure each cred is populated.
	credNames := []string{
		"OKTA_ORG",
		"OKTA_API_TOKEN",
		"OKTA_USERNAME",
	}
	testhelpers.SkipUnlessEnvVarsSet(t, credNames)

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	// Enable Userpass authentication
	mountAccessor := testhelpers.SetupUserpassMountAccessor(t, client)
	entityClient, entityID, _ := testhelpers.CreateCustomEntityAndAliasWithinMount(t,
		client, mountAccessor, "userpass", "testuser",
		map[string]interface{}{
			"name":     "test-entity",
			"policies": "mfa_policy",
			"metadata": map[string]string{
				"email": os.Getenv("OKTA_USERNAME"),
			},
		})

	err := mfaGenerateOktaPolicyMFATest(entityClient, mountAccessor, entityID)
	if err != nil {
		t.Fatalf("Okta failed: %s", err)
	}
}

func mfaGenerateOktaPolicyMFATest(client *api.Client, mountAccessor, entityID string) error {
	var err error

	rules := `
path "secret/foo" {
	capabilities = ["read"]
	mfa_methods = ["my_okta"]
}
	`

	err = client.Sys().PutPolicy("mfa_policy", rules)
	if err != nil {
		return fmt.Errorf("failed to create mfa_policy: %v", err)
	}

	mfaConfigData := map[string]interface{}{
		"mount_accessor":  mountAccessor,
		"org_name":        os.Getenv("OKTA_ORG"),
		"api_token":       os.Getenv("OKTA_API_TOKEN"),
		"primary_email":   true,
		"username_format": "{{identity.entity.metadata.email}}",
	}
	_, err = client.Logical().Write("sys/mfa/method/okta/my_okta", mfaConfigData)
	if err != nil {
		return fmt.Errorf("failed to persist TOTP MFA configuration: %v", err)
	}

	// Write some data in the path that requires TOTP MFA
	genericData := map[string]interface{}{
		"somedata": "which can only be read if MFA succeeds",
	}
	_, err = client.Logical().Write("secret/foo", genericData)
	if err != nil {
		return fmt.Errorf("failed to store data in generic backend: %v", err)
	}

	// Replace the token in client with the one that has access to MFA
	// required path
	originalToken := client.Token()
	defer client.SetToken(originalToken)

	// login to the testuser
	secret, err := client.Logical().Write("auth/userpass/login/testuser", map[string]interface{}{
		"password": "testpassword",
	})
	if err != nil {
		return fmt.Errorf("failed to login using userpass auth: %v", err)
	}

	userpassToken := secret.Auth.ClientToken
	client.SetToken(userpassToken)

	secret, err = client.Logical().Read("secret/foo")
	if err != nil {
		return fmt.Errorf("failed to read the secret: %v", err)
	}

	// It should be possible to access the secret
	// secret, err = api.ParseSecret(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to parse the secret: %v", err)
	}
	if !reflect.DeepEqual(secret.Data, genericData) {
		return fmt.Errorf("bad: generic data; expected: %#v\nactual: %#v", genericData, secret.Data)
	}
	return nil
}

func TestInteg_LoginMFAOkta(t *testing.T) {
	if os.Getenv("VAULT_ACC") == "" {
		t.Skip("This test requires manual intervention and OKTA verify on cellphone is needed")
	}

	// Ensure each cred is populated.
	credNames := []string{
		"OKTA_ORG",
		"OKTA_API_TOKEN",
		"OKTA_USERNAME",
	}
	testhelpers.SkipUnlessEnvVarsSet(t, credNames)

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	// Enable Userpass authentication
	mountAccessor := testhelpers.SetupUserpassMountAccessor(t, client)

	// Create testuser entity and alias
	entityClient, entityID, _ := testhelpers.CreateCustomEntityAndAliasWithinMount(t,
		client, mountAccessor, "userpass", "testuser",
		map[string]interface{}{
			"name": "test-entity",
			"metadata": map[string]string{
				"email": os.Getenv("OKTA_USERNAME"),
			},
		})

	err := mfaGenerateOktaLoginMFATest(entityClient, mountAccessor, entityID, t.Log)
	if err != nil {
		t.Fatalf("Okta failed: %s", err)
	}
}

func mfaGenerateOktaLoginMFATest(client *api.Client, mountAccessor, entityID string, log func(...any)) error {
	var methodID string
	var userpassToken string
	// login MFA
	{
		// create a config
		mfaConfigData := map[string]interface{}{
			"mount_accessor":  mountAccessor,
			"org_name":        os.Getenv("OKTA_ORG"),
			"api_token":       os.Getenv("OKTA_API_TOKEN"),
			"primary_email":   true,
			"username_format": "{{identity.entity.metadata.email}}",
		}
		resp, err := client.Logical().Write("identity/mfa/method/okta", mfaConfigData)

		if err != nil || (resp == nil) {
			return fmt.Errorf("bad: resp: %#v\n err: %v", resp, err)
		}

		methodID = resp.Data["method_id"].(string)
		if methodID == "" {
			return fmt.Errorf("method ID is empty")
		}
		// creating MFAEnforcementConfig
		_, err = client.Logical().Write("identity/mfa/login-enforcement/randomName", map[string]interface{}{
			"auth_method_accessors": []string{mountAccessor},
			"auth_method_types":     []string{"userpass"},
			"identity_entity_ids":   []string{entityID},
			"name":                  "randomName",
			"mfa_method_ids":        []string{methodID},
		})
		if err != nil {
			return fmt.Errorf("failed to configure MFAEnforcementConfig: %v", err)
		}
	}

	secret, err := client.Logical().Write("auth/userpass/login/testuser", map[string]interface{}{
		"password": "testpassword",
	})
	if err != nil {
		return fmt.Errorf("failed to login using userpass auth: %v", err)
	}

	if secret.Auth == nil || secret.Auth.MFARequirement == nil {
		return fmt.Errorf("two phase login returned nil MFARequirement")
	}
	if secret.Auth.MFARequirement.MFARequestID == "" {
		return fmt.Errorf("MFARequirement contains empty MFARequestID")
	}
	if secret.Auth.MFARequirement.MFAConstraints == nil || len(secret.Auth.MFARequirement.MFAConstraints) == 0 {
		return fmt.Errorf("MFAConstraints is nil or empty")
	}
	mfaConstraints, ok := secret.Auth.MFARequirement.MFAConstraints["randomName"]
	if !ok {
		return fmt.Errorf("failed to find the mfaConstraints")
	}
	if mfaConstraints.Any == nil || len(mfaConstraints.Any) == 0 {
		return fmt.Errorf("")
	}
	for _, mfaAny := range mfaConstraints.Any {
		if mfaAny.ID != methodID || mfaAny.Type != "okta" {
			return fmt.Errorf("invalid mfa constraints")
		}
	}

	// validation
	secret, err = client.Sys().MFAValidateWithContext(context.Background(),
		secret.Auth.MFARequirement.MFARequestID,
		map[string]interface{}{
			methodID: []string{},
		},
	)
	if err != nil {
		return fmt.Errorf("MFA failed: %v", err)
	}

	userpassToken = secret.Auth.ClientToken
	if secret.Auth.ClientToken == "" {
		return fmt.Errorf("MFA was not enforced")
	}

	client.SetToken(client.Token())
	secret, err = client.Logical().Write("auth/token/lookup", map[string]interface{}{
		"token": userpassToken,
	})
	if err != nil {
		return fmt.Errorf("failed to lookup userpass authenticated token: %v", err)
	}

	entityIDCheck := secret.Data["entity_id"].(string)
	if entityIDCheck != entityID {
		return fmt.Errorf("different entityID assigned")
	}

	return nil
}
