// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mongodbatlas

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
)

type testEnv struct {
	PublicKey      string
	PrivateKey     string
	ProjectID      string
	OrganizationID string

	Backend logical.Backend
	Context context.Context
	Storage logical.Storage

	MostRecentSecret *logical.Secret
}

func (e *testEnv) AddConfig(t *testing.T) {
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Storage:   e.Storage,
		Data: map[string]interface{}{
			"public_key":  e.PublicKey,
			"private_key": e.PrivateKey,
		},
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%v", resp, err)
	}
	if resp != nil {
		t.Fatal("expected nil response to represent a 204")
	}
}

func (e *testEnv) AddLeaseConfig(t *testing.T) {
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/lease",
		Storage:   e.Storage,
		Data: map[string]interface{}{
			"ttl":     "80s",
			"max_ttl": "160s",
		},
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%v", resp, err)
	}
	if resp != nil {
		t.Fatal("expected nil response to represent a 204")
	}
}

func (e *testEnv) AddProgrammaticAPIKeyRole(t *testing.T) {
	roles := []string{"ORG_MEMBER"}
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/test-programmatic-key",
		Storage:   e.Storage,
		Data: map[string]interface{}{
			"organization_id": e.OrganizationID,
			"roles":           roles,
		},
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%v", resp, err)
	}
}

func (e *testEnv) AddProgrammaticAPIKeyRoleWithProjectIDAndOrgID(t *testing.T) {
	roles := []string{"ORG_MEMBER"}
	projectRoles := []string{"GROUP_READ_ONLY"}
	ips := []string{"192.168.1.1", "192.168.1.2"}
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/test-programmatic-key",
		Storage:   e.Storage,
		Data: map[string]interface{}{
			"organization_id": e.OrganizationID,
			"project_id":      e.ProjectID,
			"roles":           roles,
			"project_roles":   projectRoles,
			"ip_addresses":    ips,
		},
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%v", resp, err)
	}
}

func (e *testEnv) AddProgrammaticAPIKeyRoleWithTTL(t *testing.T) {
	roles := []string{"ORG_MEMBER"}
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/test-programmatic-key",
		Storage:   e.Storage,
		Data: map[string]interface{}{
			"organization_id": e.OrganizationID,
			"roles":           roles,
			"ttl":             "20s",
			"max_ttl":         "60s",
		},
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%v", resp, err)
	}
}

func (e *testEnv) AddProgrammaticAPIKeyRoleWithIP(t *testing.T) {
	roles := []string{"ORG_MEMBER"}
	ips := []string{"192.168.1.1", "192.168.1.2"}
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/test-programmatic-key",
		Storage:   e.Storage,
		Data: map[string]interface{}{
			"organization_id": e.OrganizationID,
			"roles":           roles,
			"ip_addresses":    ips,
		},
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%v", resp, err)
	}
}

func (e *testEnv) AddProgrammaticAPIKeyRoleProjectWithIP(t *testing.T) {
	roles := []string{"ORG_MEMBER"}
	ips := []string{"192.168.1.1", "192.168.1.2"}
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/test-programmatic-key",
		Storage:   e.Storage,
		Data: map[string]interface{}{
			"project_id":   e.ProjectID,
			"roles":        roles,
			"ip_addresses": ips,
		},
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%v", resp, err)
	}
}

func (e *testEnv) AddProgrammaticAPIKeyRoleWithCIDR(t *testing.T) {
	roles := []string{"ORG_MEMBER"}
	cidrBlocks := []string{"179.154.224.2/32"}
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/test-programmatic-key",
		Storage:   e.Storage,
		Data: map[string]interface{}{
			"organization_id": e.OrganizationID,
			"roles":           roles,
			"cidr_blocks":     cidrBlocks,
		},
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%v", resp, err)
	}
}

func (e *testEnv) AddProgrammaticAPIKeyRoleWithCIDRAndIP(t *testing.T) {
	roles := []string{"ORG_MEMBER"}
	cidrBlocks := []string{"179.154.224.2/32"}
	ips := []string{"192.168.1.1", "192.168.1.2"}
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/test-programmatic-key",
		Storage:   e.Storage,
		Data: map[string]interface{}{
			"organization_id": e.OrganizationID,
			"roles":           roles,
			"cidr_blocks":     cidrBlocks,
			"ip_addresses":    ips,
		},
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%v", resp, err)
	}
}

func (e *testEnv) AddProgrammaticAPIKeyRoleWithProjectID(t *testing.T) {
	roles := []string{"ORG_MEMBER"}
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/test-programmatic-key",
		Storage:   e.Storage,
		Data: map[string]interface{}{
			"roles":      roles,
			"project_id": e.ProjectID,
		},
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%v", resp, err)
	}
}

func (e *testEnv) AddProgrammaticAPIKeyRoleWithProjectIDWithTTL(t *testing.T) {
	roles := []string{"ORG_MEMBER"}
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/test-programmatic-key",
		Storage:   e.Storage,
		Data: map[string]interface{}{
			"roles":      roles,
			"project_id": e.ProjectID,
			"ttl":        "20s",
			"max_ttl":    "60s",
		},
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%v", resp, err)
	}
}

func (e *testEnv) ReadProgrammaticAPIKeyRule(t *testing.T) {
	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "creds/test-programmatic-key",
		Storage:   e.Storage,
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%v", resp, err)
	}
	if resp == nil {
		t.Fatal("expected a response")
	}

	if resp.Data["public_key"] == "" {
		t.Fatal("failed to receive access_key")
	}
	if resp.Data["private_key"] == "" {
		t.Fatal("failed to receive secret_key")
	}
	e.MostRecentSecret = resp.Secret
}

func (e *testEnv) CheckLease(t *testing.T) {
	ttl := int(e.MostRecentSecret.TTL.Seconds())
	wantedTTL := 20
	maxTTL := int(e.MostRecentSecret.MaxTTL.Seconds())
	wantedMaxTTL := 60

	if ttl != wantedTTL {
		t.Fatal(fmt.Sprintf("ttl=%d, wanted=%d", ttl, wantedTTL))
	}
	if maxTTL != wantedMaxTTL {
		t.Fatal(fmt.Sprintf("maxTTL=%d, wanted=%d", ttl, wantedMaxTTL))
	}
}

func (e *testEnv) CheckExtendedLease(t *testing.T) {
	ttl := int(e.MostRecentSecret.TTL.Seconds())
	maxTTL := int(e.MostRecentSecret.MaxTTL.Seconds())
	wantedMaxTTL := 60

	if ttl != wantedMaxTTL {
		t.Fatal(fmt.Sprintf("ttl=%d, wanted=%d", ttl, wantedMaxTTL))
	}
	if maxTTL != wantedMaxTTL {
		t.Fatal(fmt.Sprintf("maxTTL=%d, wanted=%d", ttl, wantedMaxTTL))
	}
}

func (e *testEnv) RenewProgrammaticAPIKeys(t *testing.T) {
	req := &logical.Request{
		Operation: logical.RenewOperation,
		Storage:   e.Storage,
		Secret:    e.MostRecentSecret,
		Data: map[string]interface{}{
			"lease_id": "foo",
		},
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%v", resp, err)
	}
	if resp == nil {
		t.Fatal("expected a response")
	}
	if resp.Secret != e.MostRecentSecret {
		t.Fatalf("expected %+v but got %+v", e.MostRecentSecret, resp.Secret)
	}
}

func (e *testEnv) RenewProgrammaticAPIKeysWithExtendedLease(t *testing.T) {
	req := &logical.Request{
		Operation: logical.RenewOperation,
		Storage:   e.Storage,
		Secret:    e.MostRecentSecret,
		Data: map[string]interface{}{
			"lease_id":  "foo",
			"increment": "180s",
		},
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%v", resp, err)
	}
	if resp == nil {
		t.Fatal("expected a response")
	}
	if resp.Secret != e.MostRecentSecret {
		t.Fatalf("expected %+v but got %+v", e.MostRecentSecret, resp.Secret)
	}
}

func (e *testEnv) RevokeProgrammaticAPIKeys(t *testing.T) {
	req := &logical.Request{
		Operation: logical.RevokeOperation,
		Storage:   e.Storage,
		Secret:    e.MostRecentSecret,
		Data: map[string]interface{}{
			"lease_id": "foo",
		},
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%v", resp, err)
	}
	if resp != nil {
		t.Fatal("expected nil response to represent a 204")
	}
}
