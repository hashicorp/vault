package identity

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/okta"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

var (
	org_name  = "<okta org name>"
	api_token = "<okta api token>"
)

var identityOktaMFACoreConfig = &vault.CoreConfig{
	CredentialBackends: map[string]logical.Factory{
		"userpass": userpass.Factory,
		"okta":     okta.Factory,
	},
}

func TestOktaEngineMFA(t *testing.T) {
	t.Skip("This test requires manual intervention and OKTA verify on cellphone is needed")
	cluster := vault.NewTestCluster(t, identityOktaMFACoreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

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
		"org_name":  org_name,
		"api_token": api_token,
	})
	if err != nil {
		t.Fatalf("error configuring okta mount: %v", err)
	}

	_, err = client.Logical().Write("auth/okta/groups/testgroup", map[string]interface{}{
		"policies": "default",
	})
	if err != nil {
		t.Fatalf("error configuring okta group, %v", err)
	}

	_, err = client.Logical().Write("auth/okta/login/<okta username>", map[string]interface{}{
		"password": "<okta password>",
	})
	if err != nil {
		t.Fatalf("error configuring okta group, %v", err)
	}
}

func TestInteg_PolicyMFAOkta(t *testing.T) {
	t.Skip("This test requires manual intervention and OKTA verify on cellphone is needed")
	cluster := vault.NewTestCluster(t, identityOktaMFACoreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client

	// Enable Userpass authentication
	err := client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatalf("failed to enable userpass auth: %v", err)
	}

	err = mfaGenerateOktaPolicyMFATest(client)
	if err != nil {
		t.Fatalf("Okta failed: %s", err)
	}
}

func mfaGenerateOktaPolicyMFATest(client *api.Client) error {
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

	// listing auth mounts to find the mount accessor for the userpass
	auths, err := client.Sys().ListAuth()
	if err != nil {
		return fmt.Errorf("error listing auth mounts")
	}
	mountAccessor := auths["userpass/"].Accessor

	// creating a user in userpass
	_, err = client.Logical().Write("auth/userpass/users/testuser", map[string]interface{}{
		"password": "testpassword",
	})
	if err != nil {
		return fmt.Errorf("failed to configure userpass backend: %v", err)
	}

	// creating an identity with email metadata to be used for MFA validation
	secret, err := client.Logical().Write("identity/entity", map[string]interface{}{
		"name":     "test-entity",
		"policies": "mfa_policy",
		"metadata": map[string]string{
			"email": "<okta username>",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create an entity")
	}
	entityID := secret.Data["id"].(string)

	// assigning the entity ID to the testuser alias
	_, err = client.Logical().Write("identity/entity-alias", map[string]interface{}{
		"name":           "testuser",
		"canonical_id":   entityID,
		"mount_accessor": mountAccessor,
	})
	if err != nil {
		return fmt.Errorf("failed to create an entity alias")
	}

	mfaConfigData := map[string]interface{}{
		"mount_accessor":  mountAccessor,
		"org_name":        org_name,
		"api_token":       api_token,
		"primary_email":   true,
		"username_format": "{{entity.metadata.email}}",
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
	secret, err = client.Logical().Write("auth/userpass/login/testuser", map[string]interface{}{
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
	t.Skip("This test requires manual intervention and OKTA verify on cellphone is needed")
	cluster := vault.NewTestCluster(t, identityOktaMFACoreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client

	// Enable Userpass authentication
	err := client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatalf("failed to enable userpass auth: %v", err)
	}

	err = mfaGenerateOktaLoginMFATest(client)
	if err != nil {
		t.Fatalf("Okta failed: %s", err)
	}
}

func mfaGenerateOktaLoginMFATest(client *api.Client) error {
	var err error

	auths, err := client.Sys().ListAuth()
	if err != nil {
		return fmt.Errorf("failed to list auth mounts")
	}
	mountAccessor := auths["userpass/"].Accessor

	_, err = client.Logical().Write("auth/userpass/users/testuser", map[string]interface{}{
		"password": "testpassword",
	})
	if err != nil {
		return fmt.Errorf("failed to configure userpass backend: %v", err)
	}

	secret, err := client.Logical().Write("identity/entity", map[string]interface{}{
		"name": "test-entity",
		"metadata": map[string]string{
			"email": "<okta username>",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create an entity")
	}
	entityID := secret.Data["id"].(string)

	_, err = client.Logical().Write("identity/entity-alias", map[string]interface{}{
		"name":           "testuser",
		"canonical_id":   entityID,
		"mount_accessor": mountAccessor,
	})
	if err != nil {
		return fmt.Errorf("failed to create an entity alias")
	}

	var methodID string
	var userpassToken string
	// login MFA
	{
		// create a config
		mfaConfigData := map[string]interface{}{
			"mount_accessor":  mountAccessor,
			"org_name":        org_name,
			"api_token":       api_token,
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

	secret, err = client.Logical().Write("auth/userpass/login/testuser", map[string]interface{}{
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
		return fmt.Errorf("failed to find the mfaConstrains")
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
