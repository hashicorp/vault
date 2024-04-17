// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package userpass

import (
	"context"
	"crypto/tls"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/go-sockaddr"
	logicaltest "github.com/hashicorp/vault/helper/testhelpers/logical"
	"github.com/hashicorp/vault/sdk/helper/policyutil"
	"github.com/hashicorp/vault/sdk/helper/tokenutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
)

const (
	testSysTTL    = time.Hour * 10
	testSysMaxTTL = time.Hour * 20
)

func TestBackend_CRUD(t *testing.T) {
	var resp *logical.Response
	var err error

	storage := &logical.InmemStorage{}

	config := logical.TestBackendConfig()
	config.StorageView = storage

	ctx := context.Background()

	b, err := Factory(ctx, config)
	if err != nil {
		t.Fatal(err)
	}
	if b == nil {
		t.Fatalf("failed to create backend")
	}

	localhostSockAddr, err := sockaddr.NewSockAddr("127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}

	// Use new token_ forms
	resp, err = b.HandleRequest(ctx, &logical.Request{
		Path:      "users/testuser",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"password":          "testpassword",
			"token_ttl":         5,
			"token_max_ttl":     10,
			"token_policies":    []string{"foo"},
			"token_bound_cidrs": []string{"127.0.0.1"},
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v\n", resp, err)
	}

	resp, err = b.HandleRequest(ctx, &logical.Request{
		Path:      "users/testuser",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v\n", resp, err)
	}
	if resp.Data["token_ttl"].(int64) != 5 && resp.Data["token_max_ttl"].(int64) != 10 {
		t.Fatalf("bad: token_ttl and token_max_ttl are not set correctly")
	}
	if diff := deep.Equal(resp.Data["token_policies"], []string{"foo"}); diff != nil {
		t.Fatal(diff)
	}
	if diff := deep.Equal(resp.Data["token_bound_cidrs"], []*sockaddr.SockAddrMarshaler{{localhostSockAddr}}); diff != nil {
		t.Fatal(diff)
	}

	localhostSockAddr, err = sockaddr.NewSockAddr("127.0.1.1")
	if err != nil {
		t.Fatal(err)
	}

	// Use the old forms and verify that they zero out the new ones and then
	// the new ones read with the expected value
	resp, err = b.HandleRequest(ctx, &logical.Request{
		Path:      "users/testuser",
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"ttl":         "5m",
			"max_ttl":     "10m",
			"policies":    []string{"bar"},
			"bound_cidrs": []string{"127.0.1.1"},
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v\n", resp, err)
	}

	resp, err = b.HandleRequest(ctx, &logical.Request{
		Path:      "users/testuser",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr: %v\n", resp, err)
	}
	if resp.Data["ttl"].(int64) != 300 && resp.Data["max_ttl"].(int64) != 600 {
		t.Fatalf("bad: ttl and max_ttl are not set correctly")
	}
	if resp.Data["token_ttl"].(int64) != 300 && resp.Data["token_max_ttl"].(int64) != 600 {
		t.Fatalf("bad: token_ttl and token_max_ttl are not set correctly")
	}
	if diff := deep.Equal(resp.Data["policies"], []string{"bar"}); diff != nil {
		t.Fatal(diff)
	}
	if diff := deep.Equal(resp.Data["token_policies"], []string{"bar"}); diff != nil {
		t.Fatal(diff)
	}
	if diff := deep.Equal(resp.Data["bound_cidrs"], []*sockaddr.SockAddrMarshaler{{localhostSockAddr}}); diff != nil {
		t.Fatal(diff)
	}
	if diff := deep.Equal(resp.Data["token_bound_cidrs"], []*sockaddr.SockAddrMarshaler{{localhostSockAddr}}); diff != nil {
		t.Fatal(diff)
	}
}

func TestBackend_basic(t *testing.T) {
	b, err := Factory(context.Background(), &logical.BackendConfig{
		Logger: nil,
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: testSysTTL,
			MaxLeaseTTLVal:     testSysMaxTTL,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}
	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepUser(t, "web", "password", "foo"),
			testAccStepUser(t, "web2", "password", "foo"),
			testAccStepUser(t, "web3", "password", "foo"),
			testAccStepList(t, []string{"web", "web2", "web3"}),
			testAccStepLogin(t, "web", "password", []string{"default", "foo"}),
		},
	})
}

func TestBackend_userCrud(t *testing.T) {
	b, err := Factory(context.Background(), &logical.BackendConfig{
		Logger: nil,
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: testSysTTL,
			MaxLeaseTTLVal:     testSysMaxTTL,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}

	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepUser(t, "web", "password", "foo"),
			testAccStepReadUser(t, "web", "foo"),
			testAccStepDeleteUser(t, "web"),
			testAccStepReadUser(t, "web", ""),
		},
	})
}

func TestBackend_userCreateOperation(t *testing.T) {
	b, err := Factory(context.Background(), &logical.BackendConfig{
		Logger: nil,
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: testSysTTL,
			MaxLeaseTTLVal:     testSysMaxTTL,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}

	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: b,
		Steps: []logicaltest.TestStep{
			testUserCreateOperation(t, "web", "password", "foo"),
			testAccStepLogin(t, "web", "password", []string{"default", "foo"}),
		},
	})
}

func TestBackend_passwordUpdate(t *testing.T) {
	b, err := Factory(context.Background(), &logical.BackendConfig{
		Logger: nil,
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: testSysTTL,
			MaxLeaseTTLVal:     testSysMaxTTL,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}

	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepUser(t, "web", "password", "foo"),
			testAccStepReadUser(t, "web", "foo"),
			testAccStepLogin(t, "web", "password", []string{"default", "foo"}),
			testUpdatePassword(t, "web", "newpassword"),
			testAccStepLogin(t, "web", "newpassword", []string{"default", "foo"}),
		},
	})
}

func TestBackend_policiesUpdate(t *testing.T) {
	b, err := Factory(context.Background(), &logical.BackendConfig{
		Logger: nil,
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: testSysTTL,
			MaxLeaseTTLVal:     testSysMaxTTL,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}

	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepUser(t, "web", "password", "foo"),
			testAccStepReadUser(t, "web", "foo"),
			testAccStepLogin(t, "web", "password", []string{"default", "foo"}),
			testUpdatePolicies(t, "web", "foo,bar"),
			testAccStepReadUser(t, "web", "bar,foo"),
			testAccStepLogin(t, "web", "password", []string{"bar", "default", "foo"}),
		},
	})
}

func testUpdatePassword(t *testing.T, user, password string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "users/" + user + "/password",
		Data: map[string]interface{}{
			"password": password,
		},
	}
}

func testUpdatePolicies(t *testing.T, user, policies string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "users/" + user + "/policies",
		Data: map[string]interface{}{
			"policies": policies,
		},
	}
}

func testAccStepList(t *testing.T, users []string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ListOperation,
		Path:      "users",
		Check: func(resp *logical.Response) error {
			if resp.IsError() {
				return fmt.Errorf("got error response: %#v", *resp)
			}

			exp := []string{"web", "web2", "web3"}
			if !reflect.DeepEqual(exp, resp.Data["keys"].([]string)) {
				return fmt.Errorf("expected:\n%#v\ngot:\n%#v\n", exp, resp.Data["keys"])
			}
			return nil
		},
	}
}

func testAccStepLogin(t *testing.T, user string, pass string, policies []string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "login/" + user,
		Data: map[string]interface{}{
			"password": pass,
		},
		Unauthenticated: true,

		Check:     logicaltest.TestCheckAuth(policies),
		ConnState: &tls.ConnectionState{},
	}
}

func testUserCreateOperation(
	t *testing.T, name string, password string, policies string,
) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.CreateOperation,
		Path:      "users/" + name,
		Data: map[string]interface{}{
			"password": password,
			"policies": policies,
		},
	}
}

func testAccStepUser(
	t *testing.T, name string, password string, policies string,
) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "users/" + name,
		Data: map[string]interface{}{
			"password": password,
			"policies": policies,
		},
	}
}

func testAccStepDeleteUser(t *testing.T, n string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      "users/" + n,
	}
}

func testAccStepReadUser(t *testing.T, name string, policies string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "users/" + name,
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

func TestBackend_UserUpgrade(t *testing.T) {
	s := &logical.InmemStorage{}

	config := logical.TestBackendConfig()
	config.StorageView = s

	ctx := context.Background()

	b := Backend()
	if b == nil {
		t.Fatalf("failed to create backend")
	}
	if err := b.Setup(ctx, config); err != nil {
		t.Fatal(err)
	}

	foo := &UserEntry{
		Policies:   []string{"foo"},
		TTL:        time.Second,
		MaxTTL:     time.Second,
		BoundCIDRs: []*sockaddr.SockAddrMarshaler{{SockAddr: sockaddr.MustIPAddr("127.0.0.1")}},
	}

	entry, err := logical.StorageEntryJSON("user/foo", foo)
	if err != nil {
		t.Fatal(err)
	}
	err = s.Put(ctx, entry)
	if err != nil {
		t.Fatal(err)
	}

	userEntry, err := b.user(ctx, s, "foo")
	if err != nil {
		t.Fatal(err)
	}

	exp := &UserEntry{
		Policies:   []string{"foo"},
		TTL:        time.Second,
		MaxTTL:     time.Second,
		BoundCIDRs: []*sockaddr.SockAddrMarshaler{{SockAddr: sockaddr.MustIPAddr("127.0.0.1")}},
		TokenParams: tokenutil.TokenParams{
			TokenPolicies:   []string{"foo"},
			TokenTTL:        time.Second,
			TokenMaxTTL:     time.Second,
			TokenBoundCIDRs: []*sockaddr.SockAddrMarshaler{{SockAddr: sockaddr.MustIPAddr("127.0.0.1")}},
		},
	}
	if diff := deep.Equal(userEntry, exp); diff != nil {
		t.Fatal(diff)
	}
}
