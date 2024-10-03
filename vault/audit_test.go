// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/copystructure"
	"github.com/stretchr/testify/require"
)

func TestAudit_ReadOnlyViewDuringMount(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	c.auditBackends["noop"] = func(config *audit.BackendConfig, _ audit.HeaderFormatter) (audit.Backend, error) {
		err := config.SaltView.Put(context.Background(), &logical.StorageEntry{
			Key:   "bar",
			Value: []byte("baz"),
		})
		if err == nil || !strings.Contains(err.Error(), logical.ErrSetupReadOnly.Error()) {
			t.Fatalf("expected a read-only error")
		}
		factory := audit.NoopAuditFactory(nil)
		return factory(config, nil)
	}

	me := &MountEntry{
		Table: auditTableType,
		Path:  "foo",
		Type:  "noop",
	}
	err := c.enableAudit(namespace.RootContext(context.Background()), me, true)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
}

func TestCore_EnableAudit(t *testing.T) {
	c, keys, _ := TestCoreUnsealed(t)

	me := &MountEntry{
		Table: auditTableType,
		Path:  "foo",
		Type:  audit.TypeFile,
		Options: map[string]string{
			"file_path": "discard",
		},
	}
	err := c.enableAudit(namespace.RootContext(context.Background()), me, true)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !c.auditBroker.IsRegistered("foo/") {
		t.Fatalf("missing audit backend")
	}

	conf := &CoreConfig{
		AuditBackends: map[string]audit.Factory{
			audit.TypeFile:   audit.NewFileBackend,
			audit.TypeSocket: audit.NewSocketBackend,
			audit.TypeSyslog: audit.NewSyslogBackend,
		},
		DisableMlock: true,
		Physical:     c.physical,
	}

	c2, err := NewCore(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer c2.Shutdown()
	for i, key := range keys {
		unseal, err := TestCoreUnseal(c2, key)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if i+1 == len(keys) && !unseal {
			t.Fatalf("should be unsealed")
		}
	}

	// Verify matching audit tables
	if !reflect.DeepEqual(c.audit, c2.audit) {
		t.Fatalf("mismatch: %v %v", c.audit, c2.audit)
	}

	// Check for registration
	if !c2.auditBroker.IsRegistered("foo/") {
		t.Fatalf("missing audit backend")
	}
}

func TestCore_EnableAudit_MixedFailures(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)

	// Additional audit backend type that will fail.
	c.auditBackends["fail"] = func(config *audit.BackendConfig, _ audit.HeaderFormatter) (audit.Backend, error) {
		return nil, fmt.Errorf("failing enabling")
	}

	c.audit = &MountTable{
		Type: auditTableType,
		Entries: []*MountEntry{
			{
				Table: auditTableType,
				Path:  "noop/",
				Type:  audit.TypeFile,
				UUID:  "abcd",
				Options: map[string]string{
					"file_path": "discard",
				},
			},
			{
				Table: auditTableType,
				Path:  "noop2/",
				Type:  audit.TypeFile,
				UUID:  "bcde",
				Options: map[string]string{
					"file_path": "discard",
				},
			},
		},
	}

	// Both should set up successfully
	err := c.setupAudits(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// We expect this to work because the other entry is still valid
	c.audit.Entries[0].Type = "fail"
	err = c.setupAudits(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// No audit backend set up successfully, so expect error
	c.audit.Entries[1].Type = "fail"
	err = c.setupAudits(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

// Test that the local table actually gets populated as expected with local
// entries, and that upon reading the entries from both are recombined
// correctly
func TestCore_EnableAudit_Local(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)

	// Additional audit backend that will always fail.
	c.auditBackends["fail"] = func(config *audit.BackendConfig, _ audit.HeaderFormatter) (audit.Backend, error) {
		return nil, fmt.Errorf("failing enabling")
	}

	c.audit = &MountTable{
		Type: auditTableType,
		Entries: []*MountEntry{
			{
				Table:       auditTableType,
				Path:        "noop/",
				Type:        audit.TypeFile,
				UUID:        "abcd",
				Accessor:    "noop-abcd",
				NamespaceID: namespace.RootNamespaceID,
				namespace:   namespace.RootNamespace,
				Options: map[string]string{
					"file_path": "discard",
				},
			},
			{
				Table:       auditTableType,
				Path:        "noop2/",
				Type:        audit.TypeFile,
				UUID:        "bcde",
				Accessor:    "noop-bcde",
				NamespaceID: namespace.RootNamespaceID,
				namespace:   namespace.RootNamespace,
				Options: map[string]string{
					"file_path": "discard",
				},
			},
		},
	}

	// Both should set up successfully
	err := c.setupAudits(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	rawLocal, err := c.barrier.Get(context.Background(), coreLocalAuditConfigPath)
	if err != nil {
		t.Fatal(err)
	}
	if rawLocal == nil {
		t.Fatal("expected non-nil local audit")
	}
	localAuditTable := &MountTable{}
	if err := jsonutil.DecodeJSON(rawLocal.Value, localAuditTable); err != nil {
		t.Fatal(err)
	}
	if len(localAuditTable.Entries) > 0 {
		t.Fatalf("expected no entries in local audit table, got %#v", localAuditTable)
	}

	c.audit.Entries[1].Local = true
	if err := c.persistAudit(context.Background(), c.audit, false); err != nil {
		t.Fatal(err)
	}

	rawLocal, err = c.barrier.Get(context.Background(), coreLocalAuditConfigPath)
	if err != nil {
		t.Fatal(err)
	}
	if rawLocal == nil {
		t.Fatal("expected non-nil local audit")
	}
	localAuditTable = &MountTable{}
	if err := jsonutil.DecodeJSON(rawLocal.Value, localAuditTable); err != nil {
		t.Fatal(err)
	}
	if len(localAuditTable.Entries) != 1 {
		t.Fatalf("expected one entry in local audit table, got %#v", localAuditTable)
	}

	oldAudit := c.audit
	if err := c.loadAudits(context.Background()); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(oldAudit, c.audit) {
		t.Fatalf("expected\n%#v\ngot\n%#v\n", oldAudit, c.audit)
	}

	if len(c.audit.Entries) != 2 {
		t.Fatalf("expected two audit entries, got %#v", localAuditTable)
	}
}

func TestCore_DisableAudit(t *testing.T) {
	c, keys, _ := TestCoreUnsealed(t)

	existed, err := c.disableAudit(namespace.RootContext(context.Background()), "foo", true)
	if existed && err != nil {
		t.Fatalf("existed: %v; err: %v", existed, err)
	}

	me := &MountEntry{
		Table: auditTableType,
		Path:  "foo",
		Type:  audit.TypeFile,
		Options: map[string]string{
			"file_path": "discard",
		},
	}
	err = c.enableAudit(namespace.RootContext(context.Background()), me, true)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	existed, err = c.disableAudit(namespace.RootContext(context.Background()), "foo", true)
	if !existed || err != nil {
		t.Fatalf("existed: %v; err: %v", existed, err)
	}

	// Check for registration
	if c.auditBroker.IsRegistered("foo") {
		t.Fatalf("audit backend present")
	}

	conf := &CoreConfig{
		Physical:     c.physical,
		DisableMlock: true,
	}
	c2, err := NewCore(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer c2.Shutdown()
	for i, key := range keys {
		unseal, err := TestCoreUnseal(c2, key)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if i+1 == len(keys) && !unseal {
			t.Fatalf("should be unsealed")
		}
	}

	// Verify matching mount tables
	if !reflect.DeepEqual(c.audit, c2.audit) {
		t.Fatalf("mismatch:\n%#v\n%#v", c.audit, c2.audit)
	}
}

func TestCore_DefaultAuditTable(t *testing.T) {
	c, keys, _ := TestCoreUnsealed(t)
	verifyDefaultAuditTable(t, c.audit)

	// Verify we have an audit broker
	if c.auditBroker == nil {
		t.Fatalf("missing audit broker")
	}

	// Start a second core with same physical
	conf := &CoreConfig{
		Physical:     c.physical,
		DisableMlock: true,
	}
	c2, err := NewCore(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer c2.Shutdown()
	for i, key := range keys {
		unseal, err := TestCoreUnseal(c2, key)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if i+1 == len(keys) && !unseal {
			t.Fatalf("should be unsealed")
		}
	}

	// Verify matching mount tables
	if !reflect.DeepEqual(c.audit, c2.audit) {
		t.Fatalf("mismatch: %v %v", c.audit, c2.audit)
	}
}

func TestDefaultAuditTable(t *testing.T) {
	table := defaultAuditTable()
	verifyDefaultAuditTable(t, table)
}

func verifyDefaultAuditTable(t *testing.T, table *MountTable) {
	if len(table.Entries) != 0 {
		t.Fatalf("bad: %v", table.Entries)
	}
	if table.Type != auditTableType {
		t.Fatalf("bad: %v", *table)
	}
}

func TestAuditBroker_LogRequest(t *testing.T) {
	l := logging.NewVaultLogger(log.Trace)
	b, err := audit.NewBroker(l)
	if err != nil {
		t.Fatal(err)
	}

	a1 := audit.TestNoopAudit(t, "foo", nil)
	a2 := audit.TestNoopAudit(t, "bar", nil)
	err = b.Register(a1, false)
	require.NoError(t, err)
	err = b.Register(a2, false)
	require.NoError(t, err)

	auth := &logical.Auth{
		ClientToken: "foo",
		Policies:    []string{"dev", "ops"},
		Metadata: map[string]string{
			"user":   "armon",
			"source": "github",
		},
	}
	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "sys/mounts",
	}

	// Copy so we can verify nothing changed
	authCopyRaw, err := copystructure.Copy(auth)
	if err != nil {
		t.Fatal(err)
	}
	authCopy := authCopyRaw.(*logical.Auth)

	reqCopyRaw, err := copystructure.Copy(req)
	if err != nil {
		t.Fatal(err)
	}
	reqCopy := reqCopyRaw.(*logical.Request)

	// Create an identifier for the request to verify against
	req.ID, err = uuid.GenerateUUID()
	if err != nil {
		t.Fatalf("failed to generate identifier for the request: path%s err: %v", req.Path, err)
	}
	reqCopy.ID = req.ID

	reqErrs := errors.New("errs")

	logInput := &logical.LogInput{
		Auth:     authCopy,
		Request:  reqCopy,
		OuterErr: reqErrs,
	}
	ctx := namespace.RootContext(context.Background())
	err = b.LogRequest(ctx, logInput)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	for _, a := range []*audit.NoopAudit{a1, a2} {
		if !reflect.DeepEqual(a.ReqAuth[0], auth) {
			t.Fatalf("Bad: %#v", a.ReqAuth[0])
		}
		if !reflect.DeepEqual(a.Req[0], req) {
			t.Fatalf("Bad: %#v\n wanted %#v", a.Req[0], req)
		}
		if a.ReqErrs[0].Error() != reqErrs.Error() {
			t.Fatalf("expected %v, got %v", a.ReqErrs[0].Error(), reqErrs.Error())
		}
	}

	// Should still work with one failing backend
	a1.ReqErr = fmt.Errorf("failed")
	logInput = &logical.LogInput{
		Auth:    auth,
		Request: req,
	}
	if err := b.LogRequest(ctx, logInput); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should FAIL when both backends fail
	a2.ReqErr = fmt.Errorf("failed")
	if err = b.LogRequest(ctx, logInput); err != nil && !strings.HasPrefix(err.Error(), "event not processed by enough 'sink' nodes") {
		t.Fatalf("err: %v", err)
	}
}

func TestAuditBroker_LogResponse(t *testing.T) {
	l := logging.NewVaultLogger(log.Trace)
	b, err := audit.NewBroker(l)
	if err != nil {
		t.Fatal(err)
	}

	a1 := audit.TestNoopAudit(t, "foo", nil)
	a2 := audit.TestNoopAudit(t, "bar", nil)
	err = b.Register(a1, false)
	require.NoError(t, err)
	err = b.Register(a2, false)
	require.NoError(t, err)

	auth := &logical.Auth{
		NumUses:     10,
		ClientToken: "foo",
		Policies:    []string{"dev", "ops"},
		Metadata: map[string]string{
			"user":   "armon",
			"source": "github",
		},
	}
	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "sys/mounts",
	}
	resp := &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL: 1 * time.Hour,
			},
		},
		Data: map[string]interface{}{
			"user":     "root",
			"password": "password",
		},
	}
	respErr := fmt.Errorf("permission denied")

	// Copy so we can verify nothing changed
	authCopyRaw, err := copystructure.Copy(auth)
	if err != nil {
		t.Fatal(err)
	}
	authCopy := authCopyRaw.(*logical.Auth)

	reqCopyRaw, err := copystructure.Copy(req)
	if err != nil {
		t.Fatal(err)
	}
	reqCopy := reqCopyRaw.(*logical.Request)

	respCopyRaw, err := copystructure.Copy(resp)
	if err != nil {
		t.Fatal(err)
	}
	respCopy := respCopyRaw.(*logical.Response)

	logInput := &logical.LogInput{
		Auth:     authCopy,
		Request:  reqCopy,
		Response: respCopy,
		OuterErr: respErr,
	}
	ctx := namespace.RootContext(context.Background())
	err = b.LogResponse(ctx, logInput)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	for _, a := range []*audit.NoopAudit{a1, a2} {
		if !reflect.DeepEqual(a.RespAuth[0], auth) {
			t.Fatalf("Bad: %#v", a.ReqAuth[0])
		}
		if !reflect.DeepEqual(a.RespReq[0], req) {
			t.Fatalf("Bad: %#v", a.Req[0])
		}
		if !reflect.DeepEqual(a.Resp[0], resp) {
			t.Fatalf("Bad: %#v", a.Resp[0])
		}
		if a.RespErrs[0].Error() != respErr.Error() {
			t.Fatalf("expected %v, got %v", respErr, a.RespErrs[0])
		}
	}

	// Should still work with one failing backend
	a1.RespErr = fmt.Errorf("failed")
	logInput = &logical.LogInput{
		Auth:     auth,
		Request:  req,
		Response: resp,
		OuterErr: respErr,
	}
	err = b.LogResponse(ctx, logInput)
	require.NoError(t, err)

	// Should FAIL work with both failing backends
	a2.RespErr = fmt.Errorf("failed")
	err = b.LogResponse(ctx, logInput)
	require.Error(t, err)
	require.ErrorContains(t, err, "event not processed by enough 'sink' nodes")
}

func TestAuditBroker_AuditHeaders(t *testing.T) {
	logger := logging.NewVaultLogger(log.Trace)

	b, err := audit.NewBroker(logger)
	if err != nil {
		t.Fatal(err)
	}

	a1 := audit.TestNoopAudit(t, "foo", nil)
	a2 := audit.TestNoopAudit(t, "bar", nil)

	err = b.Register(a1, false)
	require.NoError(t, err)
	err = b.Register(a2, false)
	require.NoError(t, err)

	auth := &logical.Auth{
		ClientToken: "foo",
		Policies:    []string{"dev", "ops"},
		Metadata: map[string]string{
			"user":   "armon",
			"source": "github",
		},
	}
	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "sys/mounts",
		Headers: map[string][]string{
			"X-Test-Header":  {"foo"},
			"X-Vault-Header": {"bar"},
		},
	}
	respErr := fmt.Errorf("permission denied")

	// Copy so we can verify nothing changed
	reqCopyRaw, err := copystructure.Copy(req)
	if err != nil {
		t.Fatal(err)
	}
	reqCopy := reqCopyRaw.(*logical.Request)

	logInput := &logical.LogInput{
		Auth:     auth,
		Request:  reqCopy,
		OuterErr: respErr,
	}
	ctx := namespace.RootContext(context.Background())
	err = b.LogRequest(ctx, logInput)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	expected := map[string][]string{
		"x-test-header":  {"foo"},
		"x-vault-header": {"bar"},
	}

	for _, a := range []*audit.NoopAudit{a1, a2} {
		if !reflect.DeepEqual(a.ReqHeaders[0], expected) {
			t.Fatalf("Bad audited headers: %#v", a.ReqHeaders[0])
		}
	}

	// Should still work with one failing backend
	a1.ReqErr = fmt.Errorf("failed")
	logInput = &logical.LogInput{
		Auth:     auth,
		Request:  req,
		OuterErr: respErr,
	}
	err = b.LogRequest(ctx, logInput)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should FAIL work with both failing backends
	a2.ReqErr = fmt.Errorf("failed")
	if err = b.LogRequest(ctx, logInput); err != nil && !strings.HasPrefix(err.Error(), "event not processed by enough 'sink' nodes") {
		t.Fatalf("err: %v", err)
	}
}

// TestAudit_enableAudit ensures that we do not enable an audit device with Enterprise
// only options on a non-Enterprise version of Vault.
func TestAudit_enableAudit(t *testing.T) {
	cluster := NewTestCluster(t, nil, &TestClusterOptions{NumCores: 1})
	defer cluster.Cleanup()

	c := cluster.Cores[0]

	me := &MountEntry{
		Table: auditTableType,
		Path:  "foo",
		Type:  audit.TypeFile,
		Options: map[string]string{
			"file_path": "discard",
			"fallback":  "true",
		},
	}
	err := c.enableAudit(namespace.RootContext(context.Background()), me, true)

	if constants.IsEnterprise {
		require.NoError(t, err)
	} else {
		require.Error(t, err)
	}
}

// TestAudit_newAuditBackend ensures that we do not create a new audit device
// with Enterprise only options on a non-Enterprise version of Vault.
func TestAudit_newAuditBackend(t *testing.T) {
	cluster := NewTestCluster(t, nil, &TestClusterOptions{NumCores: 1})
	defer cluster.Cleanup()

	c := cluster.Cores[0]

	me := &MountEntry{
		Table: auditTableType,
		Path:  "foo",
		Type:  audit.TypeFile,
		Options: map[string]string{
			"file_path": "discard",
			"fallback":  "true",
		},
	}

	_, err := c.newAuditBackend(me, &logical.InmemStorage{}, me.Options)

	if constants.IsEnterprise {
		require.NoError(t, err)
	} else {
		require.Error(t, err)
	}
}
