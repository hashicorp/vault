package vault

import (
	"context"
	"encoding/base64"
	"fmt"
	"reflect"
	"sort"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	ldapcred "github.com/hashicorp/vault/builtin/credential/ldap"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

var (
	testPolicyName        = "testpolicy"
	rawTestPasswordPolicy = `
length = 20
rule "charset" {
	charset = "abcdefghijklmnopqrstuvwxyz"
	min_chars = 1
}
rule "charset" {
	charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	min_chars = 1
}
rule "charset" {
	charset = "0123456789"
	min_chars = 1
}`
)

func TestIdentity_BackendTemplating(t *testing.T) {
	var err error
	coreConfig := &CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       log.NewNullLogger(),
		CredentialBackends: map[string]logical.Factory{
			"ldap": ldapcred.Factory,
		},
	}

	cluster := NewTestCluster(t, coreConfig, &TestClusterOptions{})

	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0].Core

	TestWaitActive(t, core)

	req := logical.TestRequest(t, logical.UpdateOperation, "sys/auth/ldap")
	req.ClientToken = cluster.RootToken
	req.Data["type"] = "ldap"
	resp, err := core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}

	req = logical.TestRequest(t, logical.ReadOperation, "sys/auth")
	req.ClientToken = cluster.RootToken
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	accessor := resp.Data["ldap/"].(map[string]interface{})["accessor"].(string)

	// Create an entity
	req = logical.TestRequest(t, logical.UpdateOperation, "identity/entity")
	req.ClientToken = cluster.RootToken
	req.Data["name"] = "entity1"
	req.Data["metadata"] = map[string]string{
		"organization": "hashicorp",
		"team":         "vault",
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatal(err)
	}

	entityID := resp.Data["id"].(string)

	// Create an alias
	req = logical.TestRequest(t, logical.UpdateOperation, "identity/entity-alias")
	req.ClientToken = cluster.RootToken
	req.Data["name"] = "alias1"
	req.Data["canonical_id"] = entityID
	req.Data["mount_accessor"] = accessor
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatal(err)
	}

	aliasID := resp.Data["id"].(string)

	// Create a group
	req = logical.TestRequest(t, logical.UpdateOperation, "identity/group")
	req.ClientToken = cluster.RootToken
	req.Data["name"] = "group1"
	req.Data["member_entity_ids"] = []string{entityID}
	req.Data["metadata"] = map[string]string{
		"group": "vault",
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatal(err)
	}

	groupID := resp.Data["id"].(string)

	// Get the ldap mount
	sysView := core.router.MatchingSystemView(namespace.RootContext(nil), "auth/ldap/")

	tCases := []struct {
		tpl      string
		expected string
	}{
		{
			tpl:      "{{identity.entity.id}}",
			expected: entityID,
		},
		{
			tpl:      "{{identity.entity.name}}",
			expected: "entity1",
		},
		{
			tpl:      "{{identity.entity.metadata.organization}}",
			expected: "hashicorp",
		},
		{
			tpl:      "{{identity.entity.aliases." + accessor + ".id}}",
			expected: aliasID,
		},
		{
			tpl:      "{{identity.entity.aliases." + accessor + ".name}}",
			expected: "alias1",
		},
		{
			tpl:      "{{identity.groups.ids." + groupID + ".name}}",
			expected: "group1",
		},
		{
			tpl:      "{{identity.groups.names.group1.id}}",
			expected: groupID,
		},
		{
			tpl:      "{{identity.groups.names.group1.metadata.group}}",
			expected: "vault",
		},
		{
			tpl:      "{{identity.groups.ids." + groupID + ".metadata.group}}",
			expected: "vault",
		},
	}

	for _, tCase := range tCases {
		out, err := framework.PopulateIdentityTemplate(tCase.tpl, entityID, sysView)
		if err != nil {
			t.Fatal(err)
		}

		if out != tCase.expected {
			t.Fatalf("got %q, expected %q", out, tCase.expected)
		}
	}
}

func TestDynamicSystemView_GeneratePasswordFromPolicy_successful(t *testing.T) {
	var err error
	coreConfig := &CoreConfig{
		DisableMlock:       true,
		DisableCache:       true,
		Logger:             log.NewNullLogger(),
		CredentialBackends: map[string]logical.Factory{},
	}

	cluster := NewTestCluster(t, coreConfig, &TestClusterOptions{})

	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0].Core
	TestWaitActive(t, core)

	b64Policy := base64.StdEncoding.EncodeToString([]byte(rawTestPasswordPolicy))

	path := fmt.Sprintf("sys/policies/password/%s", testPolicyName)
	req := logical.TestRequest(t, logical.CreateOperation, path)
	req.ClientToken = cluster.RootToken
	req.Data["policy"] = b64Policy

	_, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	ctx = namespace.RootContext(ctx)
	dsv := TestDynamicSystemView(cluster.Cores[0].Core, nil)

	runeset := map[rune]bool{}
	runesFound := []rune{}

	for i := 0; i < 100; i++ {
		actual, err := dsv.GeneratePasswordFromPolicy(ctx, testPolicyName)
		if err != nil {
			t.Fatalf("no error expected, but got: %s", err)
		}
		for _, r := range actual {
			if runeset[r] {
				continue
			}
			runeset[r] = true
			runesFound = append(runesFound, r)
		}
	}

	sort.Sort(runes(runesFound))

	expectedRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	sort.Sort(runes(expectedRunes)) // Sort it so they can be compared

	if !reflect.DeepEqual(runesFound, expectedRunes) {
		t.Fatalf("Didn't find all characters from the charset\nActual  : [%s]\nExpected: [%s]", string(runesFound), string(expectedRunes))
	}
}

func TestDynamicSystemView_GeneratePasswordFromPolicy_failed(t *testing.T) {
	type testCase struct {
		policyName string
		getEntry   *logical.StorageEntry
		getErr     error
	}

	tests := map[string]testCase{
		"no policy name": {
			policyName: "",
		},
		"no policy found": {
			policyName: "testpolicy",
			getEntry:   nil,
			getErr:     nil,
		},
		"error retrieving policy": {
			policyName: "testpolicy",
			getEntry:   nil,
			getErr:     fmt.Errorf("a test error"),
		},
		"saved policy is malformed": {
			policyName: "testpolicy",
			getEntry: &logical.StorageEntry{
				Key:   getPasswordPolicyKey("testpolicy"),
				Value: []byte(`{"policy":"asdfahsdfasdf"}`),
			},
			getErr: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			testStorage := fakeBarrier{
				getEntry: test.getEntry,
				getErr:   test.getErr,
			}

			core := &Core{
				systemBarrierView: NewBarrierView(testStorage, "sys/"),
			}
			dsv := TestDynamicSystemView(core, nil)

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			actualPassword, err := dsv.GeneratePasswordFromPolicy(ctx, test.policyName)
			if err == nil {
				t.Fatalf("err expected, got nil")
			}
			if actualPassword != "" {
				t.Fatalf("no password expected, got %s", actualPassword)
			}
		})
	}
}

type runes []rune

func (r runes) Len() int           { return len(r) }
func (r runes) Less(i, j int) bool { return r[i] < r[j] }
func (r runes) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

type fakeBarrier struct {
	getEntry *logical.StorageEntry
	getErr   error
}

func (b fakeBarrier) Get(context.Context, string) (*logical.StorageEntry, error) {
	return b.getEntry, b.getErr
}

func (b fakeBarrier) List(context.Context, string) ([]string, error) {
	return nil, fmt.Errorf("not implemented")
}

func (b fakeBarrier) Put(context.Context, *logical.StorageEntry) error {
	return fmt.Errorf("not implemented")
}

func (b fakeBarrier) Delete(context.Context, string) error {
	return fmt.Errorf("not implemented")
}
