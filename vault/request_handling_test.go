// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"strings"
	"testing"
	"time"

	"github.com/armon/go-metrics"
	"github.com/go-test/deep"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/builtin/credential/approle"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

func TestRequestHandling_Wrapping(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)

	core.logicalBackends["kv"] = PassthroughBackendFactory

	meUUID, _ := uuid.GenerateUUID()
	err := core.mount(namespace.RootContext(nil), &MountEntry{
		Table: mountTableType,
		UUID:  meUUID,
		Path:  "wraptest",
		Type:  "kv",
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// No duration specified
	req := &logical.Request{
		Path:        "wraptest/foo",
		ClientToken: root,
		Operation:   logical.UpdateOperation,
		Data: map[string]interface{}{
			"zip": "zap",
		},
	}
	resp, err := core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	req = &logical.Request{
		Path:        "wraptest/foo",
		ClientToken: root,
		Operation:   logical.ReadOperation,
		WrapInfo: &logical.RequestWrapInfo{
			TTL: time.Duration(15 * time.Second),
		},
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil {
		t.Fatalf("bad: %v", resp)
	}
	if resp.WrapInfo == nil || resp.WrapInfo.TTL != time.Duration(15*time.Second) {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestRequestHandling_LoginWrapping(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)

	if err := core.loadMounts(namespace.RootContext(nil)); err != nil {
		t.Fatalf("err: %v", err)
	}

	core.credentialBackends["userpass"] = credUserpass.Factory

	// No duration specified
	req := &logical.Request{
		Path:        "sys/auth/userpass",
		ClientToken: root,
		Operation:   logical.UpdateOperation,
		Data: map[string]interface{}{
			"type": "userpass",
		},
		Connection: &logical.Connection{},
	}
	resp, err := core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	req.Path = "auth/userpass/users/test"
	req.Data = map[string]interface{}{
		"password": "foo",
		"policies": "default",
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	req = &logical.Request{
		Path:      "auth/userpass/login/test",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"password": "foo",
		},
		Connection: &logical.Connection{},
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil {
		t.Fatalf("bad: %v", resp)
	}
	if resp.WrapInfo != nil {
		t.Fatalf("bad: %#v", resp)
	}

	req = &logical.Request{
		Path:      "auth/userpass/login/test",
		Operation: logical.UpdateOperation,
		WrapInfo: &logical.RequestWrapInfo{
			TTL: time.Duration(15 * time.Second),
		},
		Data: map[string]interface{}{
			"password": "foo",
		},
		Connection: &logical.Connection{},
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil {
		t.Fatalf("bad: %v", resp)
	}
	if resp.WrapInfo == nil || resp.WrapInfo.TTL != time.Duration(15*time.Second) {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestRequestHandling_Login_PeriodicToken(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)

	if err := core.loadMounts(namespace.RootContext(nil)); err != nil {
		t.Fatalf("err: %v", err)
	}

	core.credentialBackends["approle"] = approle.Factory

	// Enable approle
	req := &logical.Request{
		Path:        "sys/auth/approle",
		ClientToken: root,
		Operation:   logical.UpdateOperation,
		Data: map[string]interface{}{
			"type": "approle",
		},
		Connection: &logical.Connection{},
	}
	resp, err := core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	// Create role
	req.Path = "auth/approle/role/role-period"
	req.Data = map[string]interface{}{
		"period": "5s",
	}
	_, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Get role ID
	req.Path = "auth/approle/role/role-period/role-id"
	req.Operation = logical.ReadOperation
	req.Data = nil
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.Data == nil {
		t.Fatalf("bad: %#v", resp)
	}
	roleID := resp.Data["role_id"]

	// Get secret ID
	req.Path = "auth/approle/role/role-period/secret-id"
	req.Operation = logical.UpdateOperation
	req.Data = nil
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.Data == nil {
		t.Fatalf("bad: %#v", resp)
	}
	secretID := resp.Data["secret_id"]

	// Perform login
	req = &logical.Request{
		Path:      "auth/approle/login",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"role_id":   roleID,
			"secret_id": secretID,
		},
		Connection: &logical.Connection{},
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.Auth == nil {
		t.Fatalf("bad: %v", resp)
	}
	loginToken := resp.Auth.ClientToken
	entityID := resp.Auth.EntityID
	accessor := resp.Auth.Accessor

	// Perform token lookup on the generated token
	req = &logical.Request{
		Path:        "auth/token/lookup",
		Operation:   logical.UpdateOperation,
		ClientToken: root,
		Data: map[string]interface{}{
			"token": loginToken,
		},
		Connection: &logical.Connection{},
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil {
		t.Fatalf("bad: %v", resp)
	}
	if resp.Data == nil {
		t.Fatalf("bad: %#v", resp)
	}

	if resp.Data["creation_time"].(int64) == 0 {
		t.Fatal("creation time was zero")
	}

	// Depending on timing of the test this may have ticked down, so reset it
	// back to the original value as long as it's not expired.
	if resp.Data["ttl"].(int64) > 0 && resp.Data["ttl"].(int64) < 5 {
		resp.Data["ttl"] = int64(5)
	}

	exp := map[string]interface{}{
		"accessor":         accessor,
		"creation_time":    resp.Data["creation_time"].(int64),
		"creation_ttl":     int64(5),
		"display_name":     "approle",
		"entity_id":        entityID,
		"expire_time":      resp.Data["expire_time"].(time.Time),
		"explicit_max_ttl": int64(0),
		"id":               loginToken,
		"issue_time":       resp.Data["issue_time"].(time.Time),
		"meta":             map[string]string{"role_name": "role-period"},
		"num_uses":         0,
		"orphan":           true,
		"path":             "auth/approle/login",
		"period":           int64(5),
		"policies":         []string{"default"},
		"renewable":        true,
		"ttl":              int64(5),
		"type":             "service",
	}

	if diff := deep.Equal(resp.Data, exp); diff != nil {
		t.Fatal(diff)
	}
}

func labelsMatch(actual, expected map[string]string) bool {
	for expected_label, expected_val := range expected {
		if v, ok := actual[expected_label]; ok {
			if v != expected_val {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func checkCounter(t *testing.T, inmemSink *metrics.InmemSink, keyPrefix string, expectedLabels map[string]string) {
	t.Helper()

	intervals := inmemSink.Data()
	if len(intervals) > 1 {
		t.Skip("Detected interval crossing.")
	}

	var counter *metrics.SampledValue = nil
	var labels map[string]string
	for _, c := range intervals[0].Counters {
		if !strings.HasPrefix(c.Name, keyPrefix) {
			continue
		}
		counter = &c

		labels = make(map[string]string)
		for _, l := range counter.Labels {
			labels[l.Name] = l.Value
		}

		// Distinguish between different label sets
		if labelsMatch(labels, expectedLabels) {
			break
		}
	}
	if counter == nil {
		t.Fatalf("No %q counter found with matching labels", keyPrefix)
	}

	if !labelsMatch(labels, expectedLabels) {
		t.Errorf("No matching label set, found %v", labels)
	}

	if counter.Count != 1 {
		t.Errorf("Counter number of samples %v is not 1.", counter.Count)
	}

	if counter.Sum != 1.0 {
		t.Errorf("Counter sum %v is not 1.", counter.Sum)
	}
}

func TestRequestHandling_LoginMetric(t *testing.T) {
	core, _, root, sink := TestCoreUnsealedWithMetrics(t)

	if err := core.loadMounts(namespace.RootContext(nil)); err != nil {
		t.Fatalf("err: %v", err)
	}

	core.credentialBackends["userpass"] = credUserpass.Factory

	// Setup mount
	req := &logical.Request{
		Path:        "sys/auth/userpass",
		ClientToken: root,
		Operation:   logical.UpdateOperation,
		Data: map[string]interface{}{
			"type": "userpass",
		},
		Connection: &logical.Connection{},
	}
	resp, err := core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	// Create user
	req.Path = "auth/userpass/users/test"
	req.Data = map[string]interface{}{
		"password": "foo",
		"policies": "default",
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	// Login with response wrapping
	req = &logical.Request{
		Path:      "auth/userpass/login/test",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"password": "foo",
		},
		WrapInfo: &logical.RequestWrapInfo{
			TTL: time.Duration(15 * time.Second),
		},
		Connection: &logical.Connection{},
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil {
		t.Fatalf("bad: %v", resp)
	}

	// There should be two counters
	checkCounter(t, sink, "token.creation",
		map[string]string{
			"cluster":      "test-cluster",
			"namespace":    "root",
			"auth_method":  "userpass",
			"mount_point":  "auth/userpass/",
			"creation_ttl": "+Inf",
			"token_type":   "service",
		},
	)
	checkCounter(t, sink, "token.creation",
		map[string]string{
			"cluster":      "test-cluster",
			"namespace":    "root",
			"auth_method":  "response_wrapping",
			"mount_point":  "auth/userpass/",
			"creation_ttl": "1m",
			"token_type":   "service",
		},
	)
}

func TestRequestHandling_SecretLeaseMetric(t *testing.T) {
	core, _, root, sink := TestCoreUnsealedWithMetrics(t)

	// Create a key with a lease
	req := logical.TestRequest(t, logical.UpdateOperation, "secret/foo")
	req.Data["foo"] = "bar"
	req.ClientToken = root
	resp, err := core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	// Read a key with a LeaseID
	req = logical.TestRequest(t, logical.ReadOperation, "secret/foo")
	req.ClientToken = root
	err = core.PopulateTokenEntry(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.Secret == nil || resp.Secret.LeaseID == "" {
		t.Fatalf("bad: %#v", resp)
	}

	checkCounter(t, sink, "secret.lease.creation",
		map[string]string{
			"cluster":       "test-cluster",
			"namespace":     "root",
			"secret_engine": "kv",
			"mount_point":   "secret/",
			"creation_ttl":  "+Inf",
		},
	)
}

func TestRequestHandling_DelegatedAuth(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)

	if err := core.loadMounts(namespace.RootContext(nil)); err != nil {
		t.Fatalf("err: %v", err)
	}
