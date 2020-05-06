package ldap

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"testing"
	"time"

	goldap "github.com/go-ldap/ldap/v3"
	"github.com/go-test/deep"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/testhelpers/ldap"
	logicaltest "github.com/hashicorp/vault/sdk/testing/stepwise"
	"github.com/hashicorp/vault/sdk/helper/ldaputil"
	"github.com/hashicorp/vault/sdk/helper/policyutil"
	"github.com/hashicorp/vault/sdk/helper/tokenutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
)

func createBackendWithStorage(t *testing.T) (*backend, logical.Storage) {
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

func TestLdapAuthBackend_Listing(t *testing.T) {
	b, storage := createBackendWithStorage(t)

	// Create group "testgroup"
	resp, err := b.HandleRequest(namespace.RootContext(nil), &logical.Request{
		Path:      "groups/testgroup",
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"policies": []string{"default"},
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}

	// Create group "nested/testgroup"
	resp, err = b.HandleRequest(namespace.RootContext(nil), &logical.Request{
		Path:      "groups/nested/testgroup",
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"policies": []string{"default"},
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}

	// Create user "testuser"
	resp, err = b.HandleRequest(namespace.RootContext(nil), &logical.Request{
		Path:      "users/testuser",
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"policies": []string{"default"},
			"groups":   "testgroup,nested/testgroup",
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}

	// Create user "nested/testuser"
	resp, err = b.HandleRequest(namespace.RootContext(nil), &logical.Request{
		Path:      "users/nested/testuser",
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"policies": []string{"default"},
			"groups":   "testgroup,nested/testgroup",
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}

	// List users
	resp, err = b.HandleRequest(namespace.RootContext(nil), &logical.Request{
		Path:      "users/",
		Operation: logical.ListOperation,
		Storage:   storage,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	expected := []string{"testuser", "nested/testuser"}
	if !reflect.DeepEqual(expected, resp.Data["keys"].([]string)) {
		t.Fatalf("bad: listed users; expected: %#v actual: %#v", expected, resp.Data["keys"].([]string))
	}

	// List groups
	resp, err = b.HandleRequest(namespace.RootContext(nil), &logical.Request{
		Path:      "groups/",
		Operation: logical.ListOperation,
		Storage:   storage,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v", resp, err)
	}
	expected = []string{"testgroup", "nested/testgroup"}
	if !reflect.DeepEqual(expected, resp.Data["keys"].([]string)) {
		t.Fatalf("bad: listed groups; expected: %#v actual: %#v", expected, resp.Data["keys"].([]string))
	}
}

func TestLdapAuthBackend_CaseSensitivity(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

	ctx := context.Background()

	testVals := func(caseSensitive bool) {
		// Clear storage
		userList, err := storage.List(ctx, "user/")
		if err != nil {
			t.Fatal(err)
		}
		for _, user := range userList {
			err = storage.Delete(ctx, "user/"+user)
			if err != nil {
				t.Fatal(err)
			}
		}
		groupList, err := storage.List(ctx, "group/")
		if err != nil {
			t.Fatal(err)
		}
		for _, group := range groupList {
			err = storage.Delete(ctx, "group/"+group)
			if err != nil {
				t.Fatal(err)
			}
		}

		configReq := &logical.Request{
			Path:      "config",
			Operation: logical.ReadOperation,
			Storage:   storage,
		}
		resp, err = b.HandleRequest(ctx, configReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%v resp:%#v", err, resp)
		}
		if resp == nil {
			t.Fatal("nil response")
		}
		if resp.Data["case_sensitive_names"].(bool) != caseSensitive {
			t.Fatalf("expected case sensitivity %t, got %t", caseSensitive, resp.Data["case_sensitive_names"].(bool))
		}

		groupReq := &logical.Request{
			Operation: logical.UpdateOperation,
			Data: map[string]interface{}{
				"policies": "grouppolicy",
			},
			Path:    "groups/EngineerS",
			Storage: storage,
		}
		resp, err = b.HandleRequest(ctx, groupReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%v resp:%#v", err, resp)
		}
		keys, err := storage.List(ctx, "group/")
		if err != nil {
			t.Fatal(err)
		}
		switch caseSensitive {
		case true:
			if keys[0] != "EngineerS" {
				t.Fatalf("bad: %s", keys[0])
			}
		default:
			if keys[0] != "engineers" {
				t.Fatalf("bad: %s", keys[0])
			}
		}

		userReq := &logical.Request{
			Operation: logical.UpdateOperation,
			Data: map[string]interface{}{
				"groups":   "EngineerS",
				"policies": "userpolicy",
			},
			Path:    "users/hermeS conRad",
			Storage: storage,
		}
		resp, err = b.HandleRequest(ctx, userReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%v resp:%#v", err, resp)
		}
		keys, err = storage.List(ctx, "user/")
		if err != nil {
			t.Fatal(err)
		}
		switch caseSensitive {
		case true:
			if keys[0] != "hermeS conRad" {
				t.Fatalf("bad: %s", keys[0])
			}
		default:
			if keys[0] != "hermes conrad" {
				t.Fatalf("bad: %s", keys[0])
			}
		}

		if caseSensitive {
			// The online test server is actually case sensitive so we need to
			// write again so it works
			userReq = &logical.Request{
				Operation: logical.UpdateOperation,
				Data: map[string]interface{}{
					"groups":   "EngineerS",
					"policies": "userpolicy",
				},
				Path:       "users/Hermes Conrad",
				Storage:    storage,
				Connection: &logical.Connection{},
			}
			resp, err = b.HandleRequest(ctx, userReq)
			if err != nil || (resp != nil && resp.IsError()) {
				t.Fatalf("err:%v resp:%#v", err, resp)
			}
		}

		loginReq := &logical.Request{
			Operation: logical.UpdateOperation,
			Path:      "login/Hermes Conrad",
			Data: map[string]interface{}{
				"password": "hermes",
			},
			Storage:    storage,
			Connection: &logical.Connection{},
		}
		resp, err = b.HandleRequest(ctx, loginReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%v resp:%#v", err, resp)
		}
		expected := []string{"grouppolicy", "userpolicy"}
		if !reflect.DeepEqual(expected, resp.Auth.Policies) {
			t.Fatalf("bad: policies: expected: %q, actual: %q", expected, resp.Auth.Policies)
		}
	}

	cleanup, cfg := ldap.PrepareTestContainer(t, "latest")
	defer cleanup()
	configReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data: map[string]interface{}{
			"url":       cfg.Url,
			"userattr":  cfg.UserAttr,
			"userdn":    cfg.UserDN,
			"groupdn":   cfg.GroupDN,
			"groupattr": cfg.GroupAttr,
			"binddn":    cfg.BindDN,
			"bindpass":  cfg.BindPassword,
		},
		Storage: storage,
	}
	resp, err = b.HandleRequest(ctx, configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	testVals(false)

	// Check that if the value is nil, on read it is case sensitive
	configEntry, err := b.Config(ctx, configReq)
	if err != nil {
		t.Fatal(err)
	}
	configEntry.CaseSensitiveNames = nil
	entry, err := logical.StorageEntryJSON("config", configEntry)
	if err != nil {
		t.Fatal(err)
	}
	err = configReq.Storage.Put(ctx, entry)
	if err != nil {
		t.Fatal(err)
	}

	testVals(true)
}

func TestLdapAuthBackend_UserPolicies(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

	cleanup, cfg := ldap.PrepareTestContainer(t, "latest")
	defer cleanup()
	configReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data: map[string]interface{}{
			"url":          cfg.Url,
			"userattr":     cfg.UserAttr,
			"userdn":       cfg.UserDN,
			"groupdn":      cfg.GroupDN,
			"groupattr":    cfg.GroupAttr,
			"binddn":       cfg.BindDN,
			"bindpassword": cfg.BindPassword,
		},
		Storage: storage,
	}
	resp, err = b.HandleRequest(context.Background(), configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	groupReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"policies": "grouppolicy",
		},
		Path:       "groups/engineers",
		Storage:    storage,
		Connection: &logical.Connection{},
	}
	resp, err = b.HandleRequest(context.Background(), groupReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	userReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"groups":   "engineers",
			"policies": "userpolicy",
		},
		Path:       "users/hermes conrad",
		Storage:    storage,
		Connection: &logical.Connection{},
	}

	resp, err = b.HandleRequest(context.Background(), userReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	loginReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "login/hermes conrad",
		Data: map[string]interface{}{
			"password": "hermes",
		},
		Storage:    storage,
		Connection: &logical.Connection{},
	}

	resp, err = b.HandleRequest(context.Background(), loginReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	expected := []string{"grouppolicy", "userpolicy"}
	if !reflect.DeepEqual(expected, resp.Auth.Policies) {
		t.Fatalf("bad: policies: expected: %q, actual: %q", expected, resp.Auth.Policies)
	}
}

/*
 * Acceptance test for LDAP Auth Method
 *
 * The tests here rely on a docker LDAP server:
 * [https://github.com/rroemhild/docker-test-openldap]
 *
 * ...as well as existence of a person object, `cn=Hermes Conrad,dc=example,dc=com`,
 *    which is a member of a group, `cn=admin_staff,ou=people,dc=example,dc=com`
 *
 * Querying the server from the command line:
 *   $ docker run --privileged -d -p 389:389 --name ldap --rm rroemhild/test-openldap
 *   $ ldapsearch -x -H ldap://localhost -b dc=planetexpress,dc=com -s sub uid=hermes
 *   $ ldapsearch -x -H ldap://localhost -b dc=planetexpress,dc=com -s sub \
         'member=cn=Hermes Conrad,ou=people,dc=planetexpress,dc=com'
*/
func factory(t *testing.T) logical.Backend {
	defaultLeaseTTLVal := time.Hour * 24
	maxLeaseTTLVal := time.Hour * 24 * 32
	b, err := Factory(context.Background(), &logical.BackendConfig{
		Logger: hclog.New(&hclog.LoggerOptions{
			Name:  "FactoryLogger",
			Level: hclog.Debug,
		}),
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
	cleanup, cfg := ldap.PrepareTestContainer(t, "latest")
	defer cleanup()

	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfigUrl(t, cfg),
			// Map Admin_staff group (from LDAP server) with foo policy
			testAccStepGroup(t, "admin_staff", "foo"),

			// Map engineers group (local) with bar policy
			testAccStepGroup(t, "engineers", "bar"),

			// Map hermes conrad user with local engineers group
			testAccStepUser(t, "hermes conrad", "engineers"),

			// Authenticate
			testAccStepLogin(t, "hermes conrad", "hermes"),

			// Verify both groups mappings can be listed back
			testAccStepGroupList(t, []string{"engineers", "admin_staff"}),

			// Verify user mapping can be listed back
			testAccStepUserList(t, []string{"hermes conrad"}),
		},
	})
}

func TestBackend_basic_noPolicies(t *testing.T) {
	b := factory(t)
	cleanup, cfg := ldap.PrepareTestContainer(t, "latest")
	defer cleanup()

	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfigUrl(t, cfg),
			// Create LDAP user
			testAccStepUser(t, "hermes conrad", ""),
			// Authenticate
			testAccStepLoginNoAttachedPolicies(t, "hermes conrad", "hermes"),
			testAccStepUserList(t, []string{"hermes conrad"}),
		},
	})
}

func TestBackend_basic_group_noPolicies(t *testing.T) {
	b := factory(t)
	cleanup, cfg := ldap.PrepareTestContainer(t, "latest")
	defer cleanup()

	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfigUrl(t, cfg),
			// Create engineers group with no policies
			testAccStepGroup(t, "engineers", ""),
			// Map hermes conrad user with local engineers group
			testAccStepUser(t, "hermes conrad", "engineers"),
			// Authenticate
			testAccStepLoginNoAttachedPolicies(t, "hermes conrad", "hermes"),
			// Verify group mapping can be listed back
			testAccStepGroupList(t, []string{"engineers"}),
		},
	})
}

func TestBackend_basic_authbind(t *testing.T) {
	b := factory(t)
	cleanup, cfg := ldap.PrepareTestContainer(t, "latest")
	defer cleanup()

	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfigUrlWithAuthBind(t, cfg),
			testAccStepGroup(t, "admin_staff", "foo"),
			testAccStepGroup(t, "engineers", "bar"),
			testAccStepUser(t, "hermes conrad", "engineers"),
			testAccStepLogin(t, "hermes conrad", "hermes"),
		},
	})
}

func TestBackend_basic_authbind_upndomain(t *testing.T) {
	b := factory(t)
	cleanup, cfg := ldap.PrepareTestContainer(t, "latest")
	defer cleanup()
	cfg.UPNDomain = "planetexpress.com"

	// Setup connection
	client := &ldaputil.Client{
		Logger: hclog.New(&hclog.LoggerOptions{
			Name:  "LDAPAuthTest",
			Level: hclog.Debug,
		}),
		LDAP: ldaputil.NewLDAP(),
	}
	conn, err := client.DialLDAP(cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	if err := conn.Bind("cn=admin,cn=config", cfg.BindPassword); err != nil {
		t.Fatal(err)
	}

	// Add userPrincipalName attribute type
	userPrincipleNameTypeReq := goldap.NewModifyRequest("cn={0}core,cn=schema,cn=config", nil)
	userPrincipleNameTypeReq.Add("olcAttributetypes", []string{"( 2.25.247072656268950430024439664556757516066 NAME ( 'userPrincipalName' ) SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 EQUALITY caseIgnoreMatch SINGLE-VALUE )"})
	if err := conn.Modify(userPrincipleNameTypeReq); err != nil {
		t.Fatal(err)
	}

	// Add new object class
	userPrincipleNameObjClassReq := goldap.NewModifyRequest("cn={0}core,cn=schema,cn=config", nil)
	userPrincipleNameObjClassReq.Add("olcObjectClasses", []string{"( 1.2.840.113556.6.2.6 NAME 'PrincipalNameClass' AUXILIARY MAY ( userPrincipalName ) )"})
	if err := conn.Modify(userPrincipleNameObjClassReq); err != nil {
		t.Fatal(err)
	}

	// Re-authenticate with the binddn user
	if err := conn.Bind(cfg.BindDN, cfg.BindPassword); err != nil {
		t.Fatal(err)
	}

	// Modify professor user and add userPrincipalName attribute
	profDN := "cn=Hubert J. Farnsworth,ou=people,dc=planetexpress,dc=com"
	modifyUserReq := goldap.NewModifyRequest(profDN, nil)
	modifyUserReq.Add("objectClass", []string{"PrincipalNameClass"})
	modifyUserReq.Add("userPrincipalName", []string{"professor@planetexpress.com"})
	if err := conn.Modify(modifyUserReq); err != nil {
		t.Fatal(err)
	}

	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfigUrlWithAuthBind(t, cfg),
			testAccStepLoginNoAttachedPolicies(t, "professor", "professor"),
		},
	})
}

func TestBackend_basic_discover(t *testing.T) {
	b := factory(t)
	cleanup, cfg := ldap.PrepareTestContainer(t, "latest")
	defer cleanup()

	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfigUrlWithDiscover(t, cfg),
			testAccStepGroup(t, "admin_staff", "foo"),
			testAccStepGroup(t, "engineers", "bar"),
			testAccStepUser(t, "hermes conrad", "engineers"),
			testAccStepLogin(t, "hermes conrad", "hermes"),
		},
	})
}

func TestBackend_basic_nogroupdn(t *testing.T) {
	b := factory(t)
	cleanup, cfg := ldap.PrepareTestContainer(t, "latest")
	defer cleanup()

	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfigUrlNoGroupDN(t, cfg),
			testAccStepGroup(t, "admin_staff", "foo"),
			testAccStepGroup(t, "engineers", "bar"),
			testAccStepUser(t, "hermes conrad", "engineers"),
			testAccStepLoginNoGroupDN(t, "hermes conrad", "hermes"),
		},
	})
}

func TestBackend_groupCrud(t *testing.T) {
	b := factory(t)

	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: b,
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
		CredentialBackend: b,
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
						t.Errorf("Default mismatch: deny_null_bind. Expected: '%t', received :'%s'", defaultDenyNullBind, cfg["deny_null_bind"])
					}

					return nil
				},
			},
		},
	})
}

func testAccStepConfigUrl(t *testing.T, cfg *ldaputil.ConfigEntry) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data: map[string]interface{}{
			"url":                  cfg.Url,
			"userattr":             cfg.UserAttr,
			"userdn":               cfg.UserDN,
			"groupdn":              cfg.GroupDN,
			"groupattr":            cfg.GroupAttr,
			"binddn":               cfg.BindDN,
			"bindpass":             cfg.BindPassword,
			"case_sensitive_names": true,
			"token_policies":       "abc,xyz",
			"request_timeout":      cfg.RequestTimeout,
		},
	}
}

func testAccStepConfigUrlWithAuthBind(t *testing.T, cfg *ldaputil.ConfigEntry) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data: map[string]interface{}{
			// In this test we also exercise multiple URL support
			"url":                  "foobar://ldap.example.com," + cfg.Url,
			"userattr":             cfg.UserAttr,
			"userdn":               cfg.UserDN,
			"groupdn":              cfg.GroupDN,
			"groupattr":            cfg.GroupAttr,
			"binddn":               cfg.BindDN,
			"bindpass":             cfg.BindPassword,
			"upndomain":            cfg.UPNDomain,
			"case_sensitive_names": true,
			"token_policies":       "abc,xyz",
			"request_timeout":      cfg.RequestTimeout,
		},
	}
}

func testAccStepConfigUrlWithDiscover(t *testing.T, cfg *ldaputil.ConfigEntry) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data: map[string]interface{}{
			"url":                  cfg.Url,
			"userattr":             cfg.UserAttr,
			"userdn":               cfg.UserDN,
			"groupdn":              cfg.GroupDN,
			"groupattr":            cfg.GroupAttr,
			"binddn":               cfg.BindDN,
			"bindpass":             cfg.BindPassword,
			"discoverdn":           true,
			"case_sensitive_names": true,
			"token_policies":       "abc,xyz",
			"request_timeout":      cfg.RequestTimeout,
		},
	}
}

func testAccStepConfigUrlNoGroupDN(t *testing.T, cfg *ldaputil.ConfigEntry) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data: map[string]interface{}{
			"url":                  cfg.Url,
			"userattr":             cfg.UserAttr,
			"userdn":               cfg.UserDN,
			"binddn":               cfg.BindDN,
			"bindpass":             cfg.BindPassword,
			"discoverdn":           true,
			"case_sensitive_names": true,
			"request_timeout":      cfg.RequestTimeout,
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
		CredentialBackend: b,
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

		// Verifies user hermes conrad maps to groups via local group (engineers) as well as remote group (Scientists)
		Check: logicaltest.TestCheckAuth([]string{"abc", "bar", "default", "foo", "xyz"}),
	}
}

func testAccStepLoginNoAttachedPolicies(t *testing.T, user string, pass string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "login/" + user,
		Data: map[string]interface{}{
			"password": pass,
		},
		Unauthenticated: true,

		// Verifies user hermes conrad maps to groups via local group (engineers) as well as remote group (Scientists)
		Check: logicaltest.TestCheckAuth([]string{"abc", "default", "xyz"}),
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

		// Verifies a search without defined GroupDN returns a warning rather than failing
		Check: func(resp *logical.Response) error {
			if len(resp.Warnings) != 1 {
				return fmt.Errorf("expected a warning due to no group dn, got: %#v", resp.Warnings)
			}

			return logicaltest.TestCheckAuth([]string{"bar", "default"})(resp)
		},
	}
}

func testAccStepGroupList(t *testing.T, groups []string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ListOperation,
		Path:      "groups",
		Check: func(resp *logical.Response) error {
			if resp.IsError() {
				return fmt.Errorf("got error response: %#v", *resp)
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
				return fmt.Errorf("got error response: %#v", *resp)
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

func TestLdapAuthBackend_ConfigUpgrade(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

	ctx := context.Background()

	cleanup, cfg := ldap.PrepareTestContainer(t, "latest")
	defer cleanup()
	configReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data: map[string]interface{}{
			"url":                    cfg.Url,
			"userattr":               cfg.UserAttr,
			"userdn":                 cfg.UserDN,
			"groupdn":                cfg.GroupDN,
			"groupattr":              cfg.GroupAttr,
			"binddn":                 cfg.BindDN,
			"bindpass":               cfg.BindPassword,
			"token_period":           "5m",
			"token_explicit_max_ttl": "24h",
			"request_timeout":        cfg.RequestTimeout,
		},
		Storage:    storage,
		Connection: &logical.Connection{},
	}
	resp, err = b.HandleRequest(ctx, configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	fd, err := b.getConfigFieldData()
	if err != nil {
		t.Fatal(err)
	}
	defParams, err := ldaputil.NewConfigEntry(nil, fd)
	if err != nil {
		t.Fatal(err)
	}
	falseBool := new(bool)
	*falseBool = false

	exp := &ldapConfigEntry{
		TokenParams: tokenutil.TokenParams{
			TokenPeriod:         5 * time.Minute,
			TokenExplicitMaxTTL: 24 * time.Hour,
		},
		ConfigEntry: &ldaputil.ConfigEntry{
			Url:                      cfg.Url,
			UserAttr:                 cfg.UserAttr,
			UserDN:                   cfg.UserDN,
			GroupDN:                  cfg.GroupDN,
			GroupAttr:                cfg.GroupAttr,
			BindDN:                   cfg.BindDN,
			BindPassword:             cfg.BindPassword,
			GroupFilter:              defParams.GroupFilter,
			DenyNullBind:             defParams.DenyNullBind,
			TLSMinVersion:            defParams.TLSMinVersion,
			TLSMaxVersion:            defParams.TLSMaxVersion,
			CaseSensitiveNames:       falseBool,
			UsePre111GroupCNBehavior: new(bool),
			RequestTimeout:           cfg.RequestTimeout,
		},
	}

	configEntry, err := b.Config(ctx, configReq)
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(exp, configEntry); diff != nil {
		t.Fatal(diff)
	}

	// Store just the config entry portion, for upgrade testing
	entry, err := logical.StorageEntryJSON("config", configEntry.ConfigEntry)
	if err != nil {
		t.Fatal(err)
	}
	err = configReq.Storage.Put(ctx, entry)
	if err != nil {
		t.Fatal(err)
	}

	configEntry, err = b.Config(ctx, configReq)
	if err != nil {
		t.Fatal(err)
	}
	// We won't have token params anymore so nil those out
	exp.TokenParams = tokenutil.TokenParams{}
	if diff := deep.Equal(exp, configEntry); diff != nil {
		t.Fatal(diff)
	}

}
