// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package alicloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
)

/*
testEnv allows us to reuse the same requests and response-checking
for both integration tests that don't hit Alibaba's real API, and
for acceptance tests that do hit their real API.
*/
type testEnv struct {
	AccessKey string
	SecretKey string
	RoleARN   string

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
			"access_key": e.AccessKey,
			"secret_key": e.SecretKey,
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

func (e *testEnv) ReadFirstConfig(t *testing.T) {
	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "config",
		Storage:   e.Storage,
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%v", resp, err)
	}
	if resp == nil {
		t.Fatal("expected a response")
	}
	if resp.Data["access_key"] != e.AccessKey {
		t.Fatal("expected access_key of " + e.AccessKey)
	}
	if resp.Data["secret_key"] != nil {
		t.Fatal("secret_key should not be returned")
	}
}

func (e *testEnv) UpdateConfig(t *testing.T) {
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Storage:   e.Storage,
		Data: map[string]interface{}{
			"access_key": "foo",
			"secret_key": "bar",
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

func (e *testEnv) ReadSecondConfig(t *testing.T) {
	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "config",
		Storage:   e.Storage,
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%v", resp, err)
	}
	if resp == nil {
		t.Fatal("expected a response")
	}
	if resp.Data["access_key"] != "foo" {
		t.Fatal("expected access_key of foo")
	}
	if resp.Data["secret_key"] != nil {
		t.Fatal("secret_key should not be returned")
	}
}

func (e *testEnv) DeleteConfig(t *testing.T) {
	req := &logical.Request{
		Operation: logical.DeleteOperation,
		Path:      "config",
		Storage:   e.Storage,
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%v", resp, err)
	}
	if resp != nil {
		t.Fatal("expected nil response to represent a 204")
	}
}

func (e *testEnv) ReadEmptyConfig(t *testing.T) {
	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "config",
		Storage:   e.Storage,
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%v", resp, err)
	}
	if resp != nil {
		t.Fatal("expected nil response to represent a 204")
	}
}

func (e *testEnv) AddPolicyBasedRole(t *testing.T) {
	req := &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/policy-based",
		Storage:   e.Storage,
		Data: map[string]interface{}{
			"remote_policies": []string{
				"name:AliyunOSSReadOnlyAccess,type:System",
				"name:AliyunRDSReadOnlyAccess,type:System",
			},
			"inline_policies": rawInlinePolicies,
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

func (e *testEnv) ReadPolicyBasedRole(t *testing.T) {
	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "role/policy-based",
		Storage:   e.Storage,
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%v", resp, err)
	}
	if resp == nil {
		t.Fatal("expected a response")
	}

	if resp.Data["role_arn"] != "" {
		t.Fatalf("expected no role_arn but received %s", resp.Data["role_arn"])
	}

	inlinePolicies := resp.Data["inline_policies"].([]*inlinePolicy)
	for i, inlinePolicy := range inlinePolicies {
		if inlinePolicy.PolicyDocument["Version"] != "1" {
			t.Fatalf("expected version of 1 but received %s", inlinePolicy.PolicyDocument["Version"])
		}
		stmts := inlinePolicy.PolicyDocument["Statement"].([]interface{})
		if len(stmts) != 1 {
			t.Fatalf("expected 1 statement but received %d", len(stmts))
		}
		stmt := stmts[0].(map[string]interface{})
		action := stmt["Action"].([]interface{})[0].(string)
		if stmt["Effect"] != "Allow" {
			t.Fatalf("expected Allow statement but received %s", stmt["Effect"])
		}
		resource := stmt["Resource"].([]interface{})[0].(string)
		if resource != "acs:oss:*:*:*" {
			t.Fatalf("received incorrect resource: %s", resource)
		}
		switch i {
		case 0:
			if action != "rds:*" {
				t.Fatalf("expected rds:* but received %s", action)
			}
		case 1:
			if action != "oss:*" {
				t.Fatalf("expected oss:* but received %s", action)
			}
		}
	}

	remotePolicies := resp.Data["remote_policies"].([]*remotePolicy)
	for i, remotePol := range remotePolicies {
		switch i {
		case 0:
			if remotePol.Name != "AliyunOSSReadOnlyAccess" {
				t.Fatalf("received unexpected policy name of %s", remotePol.Name)
			}
			if remotePol.Type != "System" {
				t.Fatalf("received unexpected policy type of %s", remotePol.Type)
			}
		case 1:
			if remotePol.Name != "AliyunRDSReadOnlyAccess" {
				t.Fatalf("received unexpected policy name of %s", remotePol.Name)
			}
			if remotePol.Type != "System" {
				t.Fatalf("received unexpected policy type of %s", remotePol.Type)
			}
		}
	}

	ttl := fmt.Sprintf("%d", resp.Data["ttl"])
	if ttl != "0" {
		t.Fatalf("expected ttl of 0 but received %s", ttl)
	}

	maxTTL := fmt.Sprintf("%d", resp.Data["max_ttl"])
	if maxTTL != "0" {
		t.Fatalf("expected max_ttl of 0 but received %s", maxTTL)
	}
}

func (e *testEnv) AddARNBasedRole(t *testing.T) {
	req := &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/role-based",
		Storage:   e.Storage,
		Data: map[string]interface{}{
			"role_arn": e.RoleARN,
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

func (e *testEnv) ReadARNBasedRole(t *testing.T) {
	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "role/role-based",
		Storage:   e.Storage,
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%v", resp, err)
	}
	if resp == nil {
		t.Fatal("expected a response")
	}

	if resp.Data["role_arn"] != e.RoleARN {
		t.Fatalf("received unexpected role_arn of %s", resp.Data["role_arn"])
	}

	inlinePolicies := resp.Data["inline_policies"].([]*inlinePolicy)
	if len(inlinePolicies) != 0 {
		t.Fatalf("expected no inline policies but received %+v", inlinePolicies)
	}

	remotePolicies := resp.Data["remote_policies"].([]*remotePolicy)
	if len(remotePolicies) != 0 {
		t.Fatalf("expected no remote policies but received %+v", remotePolicies)
	}

	ttl := fmt.Sprintf("%d", resp.Data["ttl"])
	if ttl != "0" {
		t.Fatalf("expected ttl of 0 but received %s", ttl)
	}

	maxTTL := fmt.Sprintf("%d", resp.Data["max_ttl"])
	if maxTTL != "0" {
		t.Fatalf("expected max_ttl of 0 but received %s", maxTTL)
	}
}

func (e *testEnv) UpdateARNBasedRole(t *testing.T) {
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "role/role-based",
		Storage:   e.Storage,
		Data: map[string]interface{}{
			"role_arn": "acs:ram::5138828231865461:role/notrustedactors",
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

func (e *testEnv) ReadUpdatedRole(t *testing.T) {
	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "role/role-based",
		Storage:   e.Storage,
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%v", resp, err)
	}
	if resp == nil {
		t.Fatal("expected a response")
	}

	if resp.Data["role_arn"] != "acs:ram::5138828231865461:role/notrustedactors" {
		t.Fatalf("received unexpected role_arn of %s", resp.Data["role_arn"])
	}

	inlinePolicies := resp.Data["inline_policies"].([]*inlinePolicy)
	if len(inlinePolicies) != 0 {
		t.Fatalf("expected no inline policies but received %+v", inlinePolicies)
	}

	remotePolicies := resp.Data["remote_policies"].([]*remotePolicy)
	if len(remotePolicies) != 0 {
		t.Fatalf("expected no remote policies but received %+v", remotePolicies)
	}

	ttl := fmt.Sprintf("%d", resp.Data["ttl"])
	if ttl != "0" {
		t.Fatalf("expected ttl of 100 but received %s", ttl)
	}

	maxTTL := fmt.Sprintf("%d", resp.Data["max_ttl"])
	if maxTTL != "0" {
		t.Fatalf("expected max_ttl of 1000 but received %s", maxTTL)
	}
}

func (e *testEnv) ListTwoRoles(t *testing.T) {
	req := &logical.Request{
		Operation: logical.ListOperation,
		Path:      "role",
		Storage:   e.Storage,
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%v", resp, err)
	}
	if resp == nil {
		t.Fatal("expected a response")
	}
	keys := resp.Data["keys"].([]string)
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys but received %d", len(keys))
	}
	if keys[0] != "policy-based" {
		t.Fatalf("expectied policy-based role name but received %s", keys[0])
	}
	if keys[1] != "role-based" {
		t.Fatalf("expected role-based role name but received %s", keys[1])
	}
}

func (e *testEnv) DeleteARNBasedRole(t *testing.T) {
	req := &logical.Request{
		Operation: logical.DeleteOperation,
		Path:      "role/role-based",
		Storage:   e.Storage,
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%v", resp, err)
	}
	if resp != nil {
		t.Fatal("expected nil response to represent a 204")
	}
}

func (e *testEnv) ListOneRole(t *testing.T) {
	req := &logical.Request{
		Operation: logical.ListOperation,
		Path:      "role",
		Storage:   e.Storage,
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%v", resp, err)
	}
	if resp == nil {
		t.Fatal("expected a response")
	}
	keys := resp.Data["keys"].([]string)
	if len(keys) != 1 {
		t.Fatalf("expected 2 keys but received %d", len(keys))
	}
	if keys[0] != "policy-based" {
		t.Fatalf("expectied policy-based role name but received %s", keys[0])
	}
}

func (e *testEnv) ReadPolicyBasedCreds(t *testing.T) {
	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "creds/policy-based",
		Storage:   e.Storage,
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%v", resp, err)
	}
	if resp == nil {
		t.Fatal("expected a response")
	}

	if resp.Data["access_key"] == "" {
		t.Fatal("failed to receive access_key")
	}
	if resp.Data["secret_key"] == "" {
		t.Fatal("failed to receive secret_key")
	}
	e.MostRecentSecret = resp.Secret
}

func (e *testEnv) RenewPolicyBasedCreds(t *testing.T) {
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

func (e *testEnv) RevokePolicyBasedCreds(t *testing.T) {
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

func (e *testEnv) ReadARNBasedCreds(t *testing.T) {
	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "creds/role-based",
		Storage:   e.Storage,
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%v", resp, err)
	}
	if resp == nil {
		t.Fatal("expected a response")
	}

	if resp.Data["access_key"] == "" {
		t.Fatal("received blank access_key")
	}
	if resp.Data["secret_key"] == "" {
		t.Fatal("received blank secret_key")
	}
	if fmt.Sprintf("%s", resp.Data["expiration"]) == "" {
		t.Fatal("received blank expiration")
	}
	if resp.Data["security_token"] == "" {
		t.Fatal("received blank security_token")
	}
	e.MostRecentSecret = resp.Secret
}

func (e *testEnv) RenewARNBasedCreds(t *testing.T) {
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
	if resp != nil {
		t.Fatal("expected nil response to represent a 204")
	}
}

func (e *testEnv) RevokeARNBasedCreds(t *testing.T) {
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

const rawInlinePolicies = `[
	{
		"Statement": [{
			"Action": ["rds:*"],
			"Effect": "Allow",
			"Resource": ["acs:oss:*:*:*"]
		}],
		"Version": "1"
	},
	{
		"Statement": [{
			"Action": ["oss:*"],
			"Effect": "Allow",
			"Resource": ["acs:oss:*:*:*"]
		}],
		"Version": "1"
	}
]
`
