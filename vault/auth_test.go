package vault

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/credential"
	"github.com/hashicorp/vault/logical"
)

type NoopCred struct {
	Root     []string
	Login    []string
	Paths    []string
	Requests []*logical.Request
	Response *logical.Response
}

func (n *NoopCred) HandleRequest(req *logical.Request) (*logical.Response, error) {
	n.Paths = append(n.Paths, req.Path)
	n.Requests = append(n.Requests, req)
	if req.Storage == nil {
		return nil, fmt.Errorf("missing view")
	}
	return n.Response, nil
}

func (n *NoopCred) RootPaths() []string {
	return n.Root
}

func (n *NoopCred) LoginPaths() []string {
	return n.Login
}

func (n *NoopCred) HandleLogin(req *credential.Request) (*credential.Response, error) {
	return nil, nil
}

func TestCore_DefaultAuthTable(t *testing.T) {
	c, key := TestCoreUnsealed(t)
	verifyDefaultAuthTable(t, c.auth)

	// Start a second core with same physical
	conf := &CoreConfig{Physical: c.physical}
	c2, err := NewCore(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	unseal, err := c2.Unseal(key)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !unseal {
		t.Fatalf("should be unsealed")
	}

	// Verify matching mount tables
	if !reflect.DeepEqual(c.auth, c2.auth) {
		t.Fatalf("mismatch: %v %v", c.auth, c2.auth)
	}
}

func TestCore_EnableCredential(t *testing.T) {
	c, key := TestCoreUnsealed(t)
	c.credentialBackends["noop"] = func(map[string]string) (credential.Backend, error) {
		return &NoopCred{}, nil
	}

	me := &MountEntry{
		Path: "foo",
		Type: "noop",
	}
	err := c.enableCredential(me)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	match := c.router.MatchingMount("auth/foo/bar")
	if match != "auth/foo/" {
		t.Fatalf("missing mount")
	}

	conf := &CoreConfig{Physical: c.physical}
	c2, err := NewCore(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	unseal, err := c2.Unseal(key)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !unseal {
		t.Fatalf("should be unsealed")
	}

	// Verify matching auth tables
	if !reflect.DeepEqual(c.auth, c2.auth) {
		t.Fatalf("mismatch: %v %v", c.auth, c2.auth)
	}
}

func TestCore_EnableCredential_Token(t *testing.T) {
	c, _ := TestCoreUnsealed(t)
	me := &MountEntry{
		Path: "foo",
		Type: "token",
	}
	err := c.enableCredential(me)
	if err.Error() != "token credential backend cannot be instantiated" {
		t.Fatalf("err: %v", err)
	}
}

func TestCore_DisableCredential(t *testing.T) {
	c, key := TestCoreUnsealed(t)
	c.credentialBackends["noop"] = func(map[string]string) (credential.Backend, error) {
		return &NoopCred{}, nil
	}

	err := c.disableCredential("foo")
	if err.Error() != "no matching backend" {
		t.Fatalf("err: %v", err)
	}

	me := &MountEntry{
		Path: "foo",
		Type: "noop",
	}
	err = c.enableCredential(me)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	err = c.disableCredential("foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	match := c.router.MatchingMount("auth/foo/bar")
	if match != "" {
		t.Fatalf("backend present")
	}

	conf := &CoreConfig{Physical: c.physical}
	c2, err := NewCore(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	unseal, err := c2.Unseal(key)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !unseal {
		t.Fatalf("should be unsealed")
	}

	// Verify matching mount tables
	if !reflect.DeepEqual(c.auth, c2.auth) {
		t.Fatalf("mismatch: %v %v", c.auth, c2.auth)
	}
}

func TestCore_DisableCredential_Protected(t *testing.T) {
	c, _ := TestCoreUnsealed(t)
	err := c.disableCredential("token")
	if err.Error() != "token credential backend cannot be disabled" {
		t.Fatalf("err: %v", err)
	}
}

func TestDefaultAuthTable(t *testing.T) {
	table := defaultAuthTable()
	verifyDefaultAuthTable(t, table)
}

func verifyDefaultAuthTable(t *testing.T, table *MountTable) {
	if len(table.Entries) != 1 {
		t.Fatalf("bad: %v", table.Entries)
	}
	for idx, entry := range table.Entries {
		switch idx {
		case 0:
			if entry.Path != "token" {
				t.Fatalf("bad: %v", entry)
			}
			if entry.Type != "token" {
				t.Fatalf("bad: %v", entry)
			}
		}
		if entry.Description == "" {
			t.Fatalf("bad: %v", entry)
		}
		if entry.UUID == "" {
			t.Fatalf("bad: %v", entry)
		}
	}
}
