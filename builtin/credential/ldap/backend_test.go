package ldap

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/mitchellh/mapstructure"
)

func factory(t *testing.T) logical.Backend {
	defaultLeaseTTLVal := time.Hour * 24
	maxLeaseTTLVal := time.Hour * 24 * 30
	b, err := Factory(&logical.BackendConfig{
		Logger: nil,
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: defaultLeaseTTLVal,
			MaxLeaseTTLVal:     maxLeaseTTLVal,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}
	return b
}

func TestBackend_basic(t *testing.T) {
	b := factory(t)

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		Backend:        b,
		Steps: []logicaltest.TestStep{
			testAccStepConfigUrl(t),
			testAccStepGroup(t, "scientists", "foo"),
			testAccStepGroup(t, "engineers", "bar"),
			testAccStepUser(t, "tesla", "engineers"),
			testAccStepLogin(t, "tesla", "password"),
		},
	})
}

func TestBackend_basic_authbind(t *testing.T) {
	b := factory(t)

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		Backend:        b,
		Steps: []logicaltest.TestStep{
			testAccStepConfigUrlWithAuthBind(t),
			testAccStepGroup(t, "scientists", "foo"),
			testAccStepGroup(t, "engineers", "bar"),
			testAccStepUser(t, "tesla", "engineers"),
			testAccStepLogin(t, "tesla", "password"),
		},
	})
}

func TestBackend_basic_discover(t *testing.T) {
	b := factory(t)

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		Backend:        b,
		Steps: []logicaltest.TestStep{
			testAccStepConfigUrlWithDiscover(t),
			testAccStepGroup(t, "scientists", "foo"),
			testAccStepGroup(t, "engineers", "bar"),
			testAccStepUser(t, "tesla", "engineers"),
			testAccStepLogin(t, "tesla", "password"),
		},
	})
}

func TestBackend_basic_nogroupdn(t *testing.T) {
	b := factory(t)

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		Backend:        b,
		Steps: []logicaltest.TestStep{
			testAccStepConfigUrlNoGroupDN(t),
			testAccStepGroup(t, "scientists", "foo"),
			testAccStepGroup(t, "engineers", "bar"),
			testAccStepUser(t, "tesla", "engineers"),
			testAccStepLoginNoGroupDN(t, "tesla", "password"),
		},
	})
}

func TestBackend_groupCrud(t *testing.T) {
	b := factory(t)

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		Backend:        b,
		Steps: []logicaltest.TestStep{
			testAccStepGroup(t, "g1", "foo"),
			testAccStepReadGroup(t, "g1", "foo"),
			testAccStepDeleteGroup(t, "g1"),
			testAccStepReadGroup(t, "g1", ""),
		},
	})
}

func testAccStepConfigUrl(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data: map[string]interface{}{
			// Online LDAP test server
			// http://www.forumsys.com/tutorials/integration-how-to/ldap/online-ldap-test-server/
			"url":      "ldap://ldap.forumsys.com",
			"userattr": "uid",
			"userdn":   "dc=example,dc=com",
			"groupdn":  "dc=example,dc=com",
		},
	}
}

func testAccStepConfigUrlWithAuthBind(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data: map[string]interface{}{
			// Online LDAP test server
			// http://www.forumsys.com/tutorials/integration-how-to/ldap/online-ldap-test-server/
			"url":      "ldap://ldap.forumsys.com",
			"userattr": "uid",
			"userdn":   "dc=example,dc=com",
			"groupdn":  "dc=example,dc=com",
			"binddn":   "cn=read-only-admin,dc=example,dc=com",
			"bindpass": "password",
		},
	}
}

func testAccStepConfigUrlWithDiscover(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data: map[string]interface{}{
			// Online LDAP test server
			// http://www.forumsys.com/tutorials/integration-how-to/ldap/online-ldap-test-server/
			"url":        "ldap://ldap.forumsys.com",
			"userattr":   "uid",
			"userdn":     "dc=example,dc=com",
			"groupdn":    "dc=example,dc=com",
			"discoverdn": true,
		},
	}
}

func testAccStepConfigUrlNoGroupDN(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data: map[string]interface{}{
			// Online LDAP test server
			// http://www.forumsys.com/tutorials/integration-how-to/ldap/online-ldap-test-server/
			"url":        "ldap://ldap.forumsys.com",
			"userattr":   "uid",
			"userdn":     "dc=example,dc=com",
			"discoverdn": true,
		},
	}
}

func testAccStepGroup(t *testing.T, group string, policies string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "groups/" + group,
		Data: map[string]interface{}{
			"policies": policies,
		},
	}
}

func testAccStepReadGroup(t *testing.T, group string, policies string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "groups/" + group,
		Check: func(resp *logical.Response) error {
			if resp == nil {
				if policies == "" {
					return nil
				}
				return fmt.Errorf("bad: %#v", resp)
			}

			var d struct {
				Policies string `mapstructure:"policies"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}

			if d.Policies != policies {
				return fmt.Errorf("bad: %#v", resp)
			}

			return nil
		},
	}
}

func testAccStepDeleteGroup(t *testing.T, group string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      "groups/" + group,
	}
}

func TestBackend_userCrud(t *testing.T) {
	b := Backend()

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		Backend:        b,
		Steps: []logicaltest.TestStep{
			testAccStepUser(t, "g1", "bar"),
			testAccStepReadUser(t, "g1", "bar"),
			testAccStepDeleteUser(t, "g1"),
			testAccStepReadUser(t, "g1", ""),
		},
	})
}

func testAccStepUser(t *testing.T, user string, groups string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "users/" + user,
		Data: map[string]interface{}{
			"groups": groups,
		},
	}
}

func testAccStepReadUser(t *testing.T, user string, groups string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "users/" + user,
		Check: func(resp *logical.Response) error {
			if resp == nil {
				if groups == "" {
					return nil
				}
				return fmt.Errorf("bad: %#v", resp)
			}

			var d struct {
				Groups string `mapstructure:"groups"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}

			if d.Groups != groups {
				return fmt.Errorf("bad: %#v", resp)
			}

			return nil
		},
	}
}

func testAccStepDeleteUser(t *testing.T, user string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      "users/" + user,
	}
}

func testAccStepLogin(t *testing.T, user string, pass string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "login/" + user,
		Data: map[string]interface{}{
			"password": pass,
		},
		Unauthenticated: true,

		Check: logicaltest.TestCheckAuth([]string{"bar", "default", "foo"}),
	}
}

func testAccStepLoginNoGroupDN(t *testing.T, user string, pass string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "login/" + user,
		Data: map[string]interface{}{
			"password": pass,
		},
		Unauthenticated: true,

		Check: func(resp *logical.Response) error {
			if len(resp.Warnings()) != 1 {
				return fmt.Errorf("expected a warning due to no group dn, got: %#v", resp.Warnings())
			}

			return logicaltest.TestCheckAuth([]string{"bar", "default"})(resp)
		},
	}
}

func TestLDAPEscape(t *testing.T) {
	testcases := map[string]string{
		"#test":       "\\#test",
		"test,hello":  "test\\,hello",
		"test,hel+lo": "test\\,hel\\+lo",
		"test\\hello": "test\\\\hello",
		"  test  ":    "\\  test \\ ",
	}

	for test, answer := range testcases {
		res := EscapeLDAPValue(test)
		if res != answer {
			t.Errorf("Failed to escape %s: %s != %s\n", test, res, answer)
		}
	}
}
