package oauth2

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/logformat"
	"github.com/hashicorp/vault/helper/policyutil"
	log "github.com/mgutz/logxi/v1"

	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
)

func TestBackend_Config(t *testing.T) {
	t.Logf("Running oauth2 tests...\n")
	b, err := Factory(&logical.BackendConfig{
		Logger: logformat.NewVaultLogger(log.LevelTrace),
		System: &logical.StaticSystemView{},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}

	username := os.Getenv("VT_OAUTH_USERNAME")
	password := os.Getenv("VT_OAUTH_PASSWORD")

	configData := map[string]interface{}{
		"client_id":     os.Getenv("VT_OAUTH_CLIENT_ID"),
		"client_secret": os.Getenv("VT_OAUTH_CLIENT_SECRET"),
		"provider_url":  os.Getenv("VT_OAUTH_PROVIDER_URL"),
		"userinfo_url":  os.Getenv("VT_OAUTH_USERINFO_URL"),
		"scope":         os.Getenv("VT_OAUTH_SCOPE"),
	}

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck:       func() { testAccPreCheck(t) },
		Backend:        b,
		Steps: []logicaltest.TestStep{
			// Write config
			testConfigCreate(t, configData),
			// Read config back out
			testConfigRead(t, configData),
			// Login should fail with bad creds
			testLoginWrite(t, username, "wrong", "auth failed", nil),
			// Login should fail with successful creds but no policy assignments
			testLoginWrite(t, username, password, "user is not a member of any authorized policy", nil),
			// Update user with local group membership
			testAccUserGroups(t, username, "local_group,local_group2"),
			// Add policy to group
			testAccGroups(t, "local_group", "local_group_policy"),
			// Test user has group policy
			testLoginWrite(t, username, password, "", []string{"local_group_policy"}),
			// Add policy to second local group
			testAccGroups(t, "local_group2", "local_group2_policy"),
			// Test user got policy from second local group
			testLoginWrite(t, username, password, "", []string{"local_group_policy", "local_group2_policy"}),
			// Add two policies to the "Everyone" group
			testAccGroups(t, "local_group3", "local_group3_policy"),
			// Test user didn't get policies of group3
			testLoginWrite(t, username, password, "", []string{"local_group_policy", "local_group2_policy"}),
		},
	})
}

func testLoginWrite(t *testing.T, username, password, reason string, policies []string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "login/" + username,
		ErrorOk:   true,
		Data: map[string]interface{}{
			"username": username,
			"password": password,
		},
		Check: func(resp *logical.Response) error {
			if resp.IsError() {
				if reason == "" {
					return resp.Error()
				} else if !strings.Contains(resp.Error().Error(), reason) {
					return fmt.Errorf("expected '%s' but got '%s'", reason, resp.Error().Error())
				}
			}

			if resp.Auth != nil {
				if !policyutil.EquivalentPolicies(resp.Auth.Policies, policies) {
					return fmt.Errorf("policy mismatch; expected %v but got %v", policies, resp.Auth.Policies)
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

func testConfigRead(t *testing.T, d map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "config",
		Check: func(resp *logical.Response) error {
			if resp.IsError() {
				return resp.Error()
			}

			if resp.Data["ProviderURL"] != d["provider_url"] {
				return fmt.Errorf("ProviderURL mismatch expected %s but got %s", d["provider_url"], resp.Data["ProviderURL"])
			}

			if resp.Data["UserInfoURL"] != d["userinfo_url"] {
				return fmt.Errorf("UserInfoURL mismatch expected %s but got %s", d["userinfo_url"], resp.Data["UserInfoURL"])
			}

			if resp.Data["ClientID"] != d["client_id"] {
				return fmt.Errorf("ClientID mismatch expected %s but got %s", d["client_id"], resp.Data["ClientID"])
			}

			if resp.Data["Scope"] != d["scope"] {
				return fmt.Errorf("Scope mismatch expected %s but got %s", d["scope"], resp.Data["Scope"])
			}

			if _, present := resp.Data["ClientSecret"]; present {
				return fmt.Errorf("ClientSecret should not be readable")
			}

			return nil
		},
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("VT_OAUTH_USERNAME"); v == "" {
		t.Fatal("VT_OAUTH_USERNAME must be set for acceptance tests")
	}

	if v := os.Getenv("VT_OAUTH_PASSWORD"); v == "" {
		t.Fatal("VT_OAUTH_PASSWORD must be set for acceptance tests")
	}

	if v := os.Getenv("VT_OAUTH_PROVIDER_URL"); v == "" {
		t.Fatal("VT_OAUTH_PROVIDER_URL must be set for acceptance tests")
	}
}

func testAccUserGroups(t *testing.T, user string, groups string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "users/" + user,
		Data: map[string]interface{}{
			"groups": groups,
		},
	}
}

func testAccGroups(t *testing.T, group string, policies string) logicaltest.TestStep {
	t.Logf("[testAccGroups] - Registering group %s, policy %s", group, policies)
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "groups/" + group,
		Data: map[string]interface{}{
			"policies": policies,
		},
	}
}
