package okta

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	logicaltest "github.com/hashicorp/vault/helper/testhelpers/logical"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/helper/policyutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func TestBackend_Config(t *testing.T) {
	defaultLeaseTTLVal := time.Hour * 12
	maxLeaseTTLVal := time.Hour * 24
	b, err := Factory(context.Background(), &logical.BackendConfig{
		Logger: logging.NewVaultLogger(log.Trace),
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: defaultLeaseTTLVal,
			MaxLeaseTTLVal:     maxLeaseTTLVal,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}

	username := os.Getenv("OKTA_USERNAME")
	password := os.Getenv("OKTA_PASSWORD")
	token := os.Getenv("OKTA_API_TOKEN")

	configData := map[string]interface{}{
		"org_name": os.Getenv("OKTA_ORG"),
		"base_url": "oktapreview.com",
	}

	updatedDuration := time.Hour * 1
	configDataToken := map[string]interface{}{
		"api_token": token,
		"token_ttl": "1h",
	}

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest:    true,
		PreCheck:          func() { testAccPreCheck(t) },
		CredentialBackend: b,
		Steps: []logicaltest.TestStep{
			testConfigCreate(t, configData),
			// 2. Login with bad password, expect failure (E0000004=okta auth failure).
			testLoginWrite(t, username, "wrong", "E0000004", 0, nil),
			// 3. Make our user belong to two groups and have one user-specific policy.
			testAccUserGroups(t, username, "local_grouP,lOcal_group2", []string{"user_policy"}),
			// 4. Create the group local_group, assign it a single policy.
			testAccGroups(t, "local_groUp", "loCal_group_policy"),
			// 5. Login with good password, expect user to have their user-specific
			// policy and the policy of the one valid group they belong to.
			testLoginWrite(t, username, password, "", defaultLeaseTTLVal, []string{"local_group_policy", "user_policy"}),
			// 6. Create the group everyone, assign it two policies.  This is a
			// magic group name in okta that always exists and which every
			// user automatically belongs to.
			testAccGroups(t, "everyoNe", "everyone_grouP_policy,eveRy_group_policy2"),
			// 7. Login as before, expect same result
			testLoginWrite(t, username, password, "", defaultLeaseTTLVal, []string{"local_group_policy", "user_policy"}),
			// 8. Add API token so we can lookup groups
			testConfigUpdate(t, configDataToken),
			testConfigRead(t, token, configData),
			// 10. Login should now lookup okta groups; since all okta users are
			// in the "everyone" group, that should be returned; since we
			// defined policies attached to the everyone group, we should now
			// see those policies attached to returned vault token.
			testLoginWrite(t, username, password, "", updatedDuration, []string{"everyone_group_policy", "every_group_policy2", "local_group_policy", "user_policy"}),
			testAccGroups(t, "locAl_group2", "testgroup_group_policy"),
			testLoginWrite(t, username, password, "", updatedDuration, []string{"everyone_group_policy", "every_group_policy2", "local_group_policy", "testgroup_group_policy", "user_policy"}),
		},
	})
}

func testLoginWrite(t *testing.T, username, password, reason string, expectedTTL time.Duration, policies []string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "login/" + username,
		ErrorOk:   true,
		Data: map[string]interface{}{
			"password": password,
		},
		Check: func(resp *logical.Response) error {
			if resp.IsError() {
				if reason == "" || !strings.Contains(resp.Error().Error(), reason) {
					return resp.Error()
				}
			} else if reason != "" {
				return fmt.Errorf("expected error containing %q, got no error", reason)

			}

			if resp.Auth != nil {
				if !policyutil.EquivalentPolicies(resp.Auth.Policies, policies) {
					return fmt.Errorf("policy mismatch expected %v but got %v", policies, resp.Auth.Policies)
				}

				actualTTL := resp.Auth.LeaseOptions.TTL
				if actualTTL != expectedTTL {
					return fmt.Errorf("TTL mismatch expected %v but got %v", expectedTTL, actualTTL)
				}
			}

			return nil
		},
	}
}

func testConfigCreate(t *testing.T, d map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.CreateOperation,
		Path:      "config",
		Data:      d,
	}
}

func testConfigUpdate(t *testing.T, d map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data:      d,
	}
}

func testConfigRead(t *testing.T, token string, d map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "config",
		Check: func(resp *logical.Response) error {
			if resp.IsError() {
				return resp.Error()
			}

			if resp.Data["org_name"] != d["org_name"] {
				return fmt.Errorf("org mismatch expected %s but got %s", d["organization"], resp.Data["Org"])
			}

			if resp.Data["base_url"] != d["base_url"] {
				return fmt.Errorf("BaseURL mismatch expected %s but got %s", d["base_url"], resp.Data["BaseURL"])
			}

			for _, value := range resp.Data {
				if value == token {
					return fmt.Errorf("token should not be returned on a read request")
				}
			}

			return nil
		},
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("OKTA_USERNAME"); v == "" {
		t.Fatal("OKTA_USERNAME must be set for acceptance tests")
	}

	if v := os.Getenv("OKTA_PASSWORD"); v == "" {
		t.Fatal("OKTA_PASSWORD must be set for acceptance tests")
	}

	if v := os.Getenv("OKTA_ORG"); v == "" {
		t.Fatal("OKTA_ORG must be set for acceptance tests")
	}

	if v := os.Getenv("OKTA_API_TOKEN"); v == "" {
		t.Fatal("OKTA_API_TOKEN must be set for acceptance tests")
	}
}

func testAccUserGroups(t *testing.T, user string, groups interface{}, policies interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "users/" + user,
		Data: map[string]interface{}{
			"groups":   groups,
			"policies": policies,
		},
	}
}

func testAccGroups(t *testing.T, group string, policies interface{}) logicaltest.TestStep {
	t.Logf("[testAccGroups] - Registering group %s, policy %s", group, policies)
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "groups/" + group,
		Data: map[string]interface{}{
			"policies": policies,
		},
	}
}

func testAccLogin(t *testing.T, user, password string, keys []string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "login/" + user,
		Data: map[string]interface{}{
			"password": password,
		},
		Unauthenticated: true,

		Check: logicaltest.TestCheckAuth(keys),
	}
}
