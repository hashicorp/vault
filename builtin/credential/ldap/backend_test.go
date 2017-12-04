package ldap

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/mitchellh/mapstructure"
)

func createBackendWithStorage(t *testing.T) (*backend, logical.Storage) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	b := Backend()
	if b == nil {
		t.Fatalf("failed to create backend")
	}

	err := b.Backend.Setup(config)
	if err != nil {
		t.Fatal(err)
	}

	return b, config.StorageView
}

func TestLdapAuthBackend_UserPolicies(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

	configReq := &logical.Request{
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
		},
		Storage: storage,
	}
	resp, err = b.HandleRequest(configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	groupReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"policies": "grouppolicy",
		},
		Path:    "groups/engineers",
		Storage: storage,
	}
	resp, err = b.HandleRequest(groupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	userReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"groups":   "engineers",
			"policies": "userpolicy",
		},
		Path:    "users/tesla",
		Storage: storage,
	}

	resp, err = b.HandleRequest(userReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	loginReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "login/tesla",
		Data: map[string]interface{}{
			"password": "password",
		},
		Storage: storage,
	}

	resp, err = b.HandleRequest(loginReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	expected := []string{"grouppolicy", "userpolicy"}
	if !reflect.DeepEqual(expected, resp.Auth.Policies) {
		t.Fatalf("bad: policies: expected: %q, actual: %q", expected, resp.Auth.Policies)
	}
}

/*
 * Acceptance test for LDAP Auth Backend
 *
 * The tests here rely on a public LDAP server:
 * [http://www.forumsys.com/tutorials/integration-how-to/ldap/online-ldap-test-server/]
 *
 * ...as well as existence of a person object, `uid=tesla,dc=example,dc=com`,
 *    which is a member of a group, `ou=scientists,dc=example,dc=com`
 *
 * Querying the server from the command line:
 *   $ ldapsearch -x -H ldap://ldap.forumsys.com -b dc=example,dc=com -s sub \
 *       '(&(objectClass=groupOfUniqueNames)(uniqueMember=uid=tesla,dc=example,dc=com))'
 *
 *   $ ldapsearch -x -H ldap://ldap.forumsys.com -b dc=example,dc=com -s sub uid=tesla
 */
func factory(t *testing.T) logical.Backend {
	defaultLeaseTTLVal := time.Hour * 24
	maxLeaseTTLVal := time.Hour * 24 * 32
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
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfigUrl(t),
			// Map Scientists group (from LDAP server) with foo policy
			testAccStepGroup(t, "Scientists", "foo"),

			// Map engineers group (local) with bar policy
			testAccStepGroup(t, "engineers", "bar"),

			// Map tesla user with local engineers group
			testAccStepUser(t, "tesla", "engineers"),

			// Authenticate
			testAccStepLogin(t, "tesla", "password"),

			// Verify both groups mappings can be listed back
			testAccStepGroupList(t, []string{"engineers", "Scientists"}),

			// Verify user mapping can be listed back
			testAccStepUserList(t, []string{"tesla"}),
		},
	})
}

func TestBackend_basic_authbind(t *testing.T) {
	b := factory(t)

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfigUrlWithAuthBind(t),
			testAccStepGroup(t, "Scientists", "foo"),
			testAccStepGroup(t, "engineers", "bar"),
			testAccStepUser(t, "tesla", "engineers"),
			testAccStepLogin(t, "tesla", "password"),
		},
	})
}

func TestBackend_basic_discover(t *testing.T) {
	b := factory(t)

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfigUrlWithDiscover(t),
			testAccStepGroup(t, "Scientists", "foo"),
			testAccStepGroup(t, "engineers", "bar"),
			testAccStepUser(t, "tesla", "engineers"),
			testAccStepLogin(t, "tesla", "password"),
		},
	})
}

func TestBackend_basic_nogroupdn(t *testing.T) {
	b := factory(t)

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfigUrlNoGroupDN(t),
			testAccStepGroup(t, "Scientists", "foo"),
			testAccStepGroup(t, "engineers", "bar"),
			testAccStepUser(t, "tesla", "engineers"),
			testAccStepLoginNoGroupDN(t, "tesla", "password"),
		},
	})
}

func TestBackend_groupCrud(t *testing.T) {
	b := factory(t)

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepGroup(t, "g1", "foo"),
			testAccStepReadGroup(t, "g1", "foo"),
			testAccStepDeleteGroup(t, "g1"),
			testAccStepReadGroup(t, "g1", ""),
		},
	})
}

/*
 * Test backend configuration defaults are successfully read.
 */
func TestBackend_configDefaultsAfterUpdate(t *testing.T) {
	b := factory(t)

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			logicaltest.TestStep{
				Operation: logical.UpdateOperation,
				Path:      "config",
				Data:      map[string]interface{}{},
			},
			logicaltest.TestStep{
				Operation: logical.ReadOperation,
				Path:      "config",
				Check: func(resp *logical.Response) error {
					if resp == nil {
						return fmt.Errorf("bad: %#v", resp)
					}

					// Test well-known defaults
					cfg := resp.Data
					defaultGroupFilter := "(|(memberUid={{.Username}})(member={{.UserDN}})(uniqueMember={{.UserDN}}))"
					if cfg["groupfilter"] != defaultGroupFilter {
						t.Errorf("Default mismatch: groupfilter. Expected: '%s', received :'%s'", defaultGroupFilter, cfg["groupfilter"])
					}

					defaultGroupAttr := "cn"
					if cfg["groupattr"] != defaultGroupAttr {
						t.Errorf("Default mismatch: groupattr. Expected: '%s', received :'%s'", defaultGroupAttr, cfg["groupattr"])
					}

					defaultUserAttr := "cn"
					if cfg["userattr"] != defaultUserAttr {
						t.Errorf("Default mismatch: userattr. Expected: '%s', received :'%s'", defaultUserAttr, cfg["userattr"])
					}

					defaultDenyNullBind := true
					if cfg["deny_null_bind"] != defaultDenyNullBind {
						t.Errorf("Default mismatch: deny_null_bind. Expected: '%s', received :'%s'", defaultDenyNullBind, cfg["deny_null_bind"])
					}

					return nil
				},
			},
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
			// In this test we also exercise multiple URL support
			"url":      "foobar://ldap.example.com,ldap://ldap.forumsys.com",
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
	t.Logf("[testAccStepGroup] - Registering group %s, policy %s", group, policies)
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
				Policies []string `mapstructure:"policies"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}

			if !reflect.DeepEqual(d.Policies, policyutil.ParsePolicies(policies)) {
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
		Backend: b,
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

		// Verifies user tesla maps to groups via local group (engineers) as well as remote group (Scientiests)
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

		// Verifies a search without defined GroupDN returns a warnting rather than failing
		Check: func(resp *logical.Response) error {
			if len(resp.Warnings) != 1 {
				return fmt.Errorf("expected a warning due to no group dn, got: %#v", resp.Warnings)
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

func testAccStepGroupList(t *testing.T, groups []string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ListOperation,
		Path:      "groups",
		Check: func(resp *logical.Response) error {
			if resp.IsError() {
				return fmt.Errorf("Got error response: %#v", *resp)
			}

			expected := make([]string, len(groups))
			copy(expected, groups)
			sort.Strings(expected)

			sortedResponse := make([]string, len(resp.Data["keys"].([]string)))
			copy(sortedResponse, resp.Data["keys"].([]string))
			sort.Strings(sortedResponse)

			if !reflect.DeepEqual(expected, sortedResponse) {
				return fmt.Errorf("expected:\n%#v\ngot:\n%#v\n", expected, sortedResponse)
			}
			return nil
		},
	}
}

func testAccStepUserList(t *testing.T, users []string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ListOperation,
		Path:      "users",
		Check: func(resp *logical.Response) error {
			if resp.IsError() {
				return fmt.Errorf("Got error response: %#v", *resp)
			}

			expected := make([]string, len(users))
			copy(expected, users)
			sort.Strings(expected)

			sortedResponse := make([]string, len(resp.Data["keys"].([]string)))
			copy(sortedResponse, resp.Data["keys"].([]string))
			sort.Strings(sortedResponse)

			if !reflect.DeepEqual(expected, sortedResponse) {
				return fmt.Errorf("expected:\n%#v\ngot:\n%#v\n", expected, sortedResponse)
			}
			return nil
		},
	}
}
