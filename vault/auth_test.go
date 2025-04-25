// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/helper/versions"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func TestAuth_ReadOnlyViewDuringMount(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	c.credentialBackends["noop"] = func(ctx context.Context, config *logical.BackendConfig) (logical.Backend, error) {
		err := config.StorageView.Put(ctx, &logical.StorageEntry{
			Key:   "bar",
			Value: []byte("baz"),
		})
		if err == nil || !strings.Contains(err.Error(), logical.ErrSetupReadOnly.Error()) {
			t.Fatalf("expected a read-only error")
		}
		return &NoopBackend{
			BackendType: logical.TypeCredential,
		}, nil
	}

	me := &MountEntry{
		Table: credentialTableType,
		Path:  "foo",
		Type:  "noop",
	}
	err := c.enableCredential(namespace.RootContext(nil), me)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
}

func TestAuthMountMetrics(t *testing.T) {
	c, _, _, _ := TestCoreUnsealedWithMetrics(t)
	c.credentialBackends["noop"] = func(ctx context.Context, config *logical.BackendConfig) (logical.Backend, error) {
		return &NoopBackend{
			BackendType: logical.TypeCredential,
		}, nil
	}
	mountKeyName := "core.mount_table.num_entries.type|auth||local|false||"
	mountMetrics := &c.metricsHelper.LoopMetrics.Metrics
	loadMetric, ok := mountMetrics.Load(mountKeyName)
	var numEntriesMetric metricsutil.GaugeMetric = loadMetric.(metricsutil.GaugeMetric)

	// 1 default nonlocal auth backend
	if !ok || numEntriesMetric.Value != 1 {
		t.Fatalf("Auth values should be: %+v", numEntriesMetric)
	}

	me := &MountEntry{
		Table: credentialTableType,
		Path:  "foo",
		Type:  "noop",
	}
	err := c.enableCredential(namespace.RootContext(nil), me)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	mountMetrics = &c.metricsHelper.LoopMetrics.Metrics
	loadMetric, ok = mountMetrics.Load(mountKeyName)
	numEntriesMetric = loadMetric.(metricsutil.GaugeMetric)
	if !ok || numEntriesMetric.Value != 2 {
		t.Fatalf("mount metrics for num entries do not match true values")
	}
	if len(numEntriesMetric.Key) != 3 ||
		numEntriesMetric.Key[0] != "core" ||
		numEntriesMetric.Key[1] != "mount_table" ||
		numEntriesMetric.Key[2] != "num_entries" {
		t.Fatalf("mount metrics for num entries have wrong key")
	}
	if len(numEntriesMetric.Labels) != 2 ||
		numEntriesMetric.Labels[0].Name != "type" ||
		numEntriesMetric.Labels[0].Value != "auth" ||
		numEntriesMetric.Labels[1].Name != "local" ||
		numEntriesMetric.Labels[1].Value != "false" {
		t.Fatalf("mount metrics for num entries have wrong labels")
	}
	mountSizeKeyName := "core.mount_table.size.type|auth||local|false||"
	loadMetric, ok = mountMetrics.Load(mountSizeKeyName)
	sizeMetric := loadMetric.(metricsutil.GaugeMetric)

	if !ok {
		t.Fatalf("mount metrics for size do not match exist")
	}
	if len(sizeMetric.Key) != 3 ||
		sizeMetric.Key[0] != "core" ||
		sizeMetric.Key[1] != "mount_table" ||
		sizeMetric.Key[2] != "size" {
		t.Fatalf("mount metrics for size have wrong key")
	}
	if len(sizeMetric.Labels) != 2 ||
		sizeMetric.Labels[0].Name != "type" ||
		sizeMetric.Labels[0].Value != "auth" ||
		sizeMetric.Labels[1].Name != "local" ||
		sizeMetric.Labels[1].Value != "false" {
		t.Fatalf("mount metrics for size have wrong labels")
	}
}

func TestCore_DefaultAuthTable(t *testing.T) {
	c, keys, _ := TestCoreUnsealed(t)
	verifyDefaultAuthTable(t, c.auth)

	// Start a second core with same physical
	inmemSink := metrics.NewInmemSink(1000000*time.Hour, 2000000*time.Hour)
	conf := &CoreConfig{
		Physical:        c.physical,
		DisableMlock:    true,
		BuiltinRegistry: corehelpers.NewMockBuiltinRegistry(),
		MetricSink:      metricsutil.NewClusterMetricSink("test-cluster", inmemSink),
		MetricsHelper:   metricsutil.NewMetricsHelper(inmemSink, false),
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
	if !reflect.DeepEqual(c.auth, c2.auth) {
		t.Fatalf("mismatch: %v %v", c.auth, c2.auth)
	}
}

func TestCore_BuiltinRegistry(t *testing.T) {
	conf := &CoreConfig{
		// set PluginDirectory and ensure that vault doesn't expect approle to
		// be there when we are mounting the builtin approle
		PluginDirectory: "/Users/foo",

		DisableMlock:    true,
		BuiltinRegistry: corehelpers.NewMockBuiltinRegistry(),
	}
	c, _, _ := TestCoreUnsealedWithConfig(t, conf)

	for _, me := range []*MountEntry{
		{
			Table: credentialTableType,
			Path:  "approle/",
			Type:  "approle",
		},
		{
			Table:   credentialTableType,
			Path:    "approle2/",
			Type:    "approle",
			Version: versions.GetBuiltinVersion(consts.PluginTypeCredential, "approle"),
		},
	} {
		err := c.enableCredential(namespace.RootContext(nil), me)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
	}
}

func TestCore_EnableCredential(t *testing.T) {
	c, keys, _ := TestCoreUnsealed(t)
	c.credentialBackends["noop"] = func(context.Context, *logical.BackendConfig) (logical.Backend, error) {
		return &NoopBackend{
			BackendType: logical.TypeCredential,
		}, nil
	}

	me := &MountEntry{
		Table: credentialTableType,
		Path:  "foo",
		Type:  "noop",
	}
	err := c.enableCredential(namespace.RootContext(nil), me)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	match := c.router.MatchingMount(namespace.RootContext(nil), "auth/foo/bar")
	if match != "auth/foo/" {
		t.Fatalf("missing mount, match: %q", match)
	}

	inmemSink := metrics.NewInmemSink(1000000*time.Hour, 2000000*time.Hour)
	conf := &CoreConfig{
		Physical:        c.physical,
		DisableMlock:    true,
		BuiltinRegistry: corehelpers.NewMockBuiltinRegistry(),
		MetricSink:      metricsutil.NewClusterMetricSink("test-cluster", inmemSink),
		MetricsHelper:   metricsutil.NewMetricsHelper(inmemSink, false),
	}
	c2, err := NewCore(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer c2.Shutdown()
	c2.credentialBackends["noop"] = func(context.Context, *logical.BackendConfig) (logical.Backend, error) {
		return &NoopBackend{
			BackendType: logical.TypeCredential,
		}, nil
	}
	for i, key := range keys {
		unseal, err := TestCoreUnseal(c2, key)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if i+1 == len(keys) && !unseal {
			t.Fatalf("should be unsealed")
		}
	}

	// Verify matching auth tables
	if !reflect.DeepEqual(c.auth, c2.auth) {
		t.Fatalf("mismatch: %v %v", c.auth, c2.auth)
	}
}

// TestCore_EnableCredential_aws_ec2 tests that we can successfully mount aws
// auth using the alias "aws-ec2"
func TestCore_EnableCredential_aws_ec2(t *testing.T) {
	c, keys, _ := TestCoreUnsealed(t)
	c.credentialBackends["aws"] = func(context.Context, *logical.BackendConfig) (logical.Backend, error) {
		return &NoopBackend{
			BackendType: logical.TypeCredential,
		}, nil
	}

	me := &MountEntry{
		Table: credentialTableType,
		Path:  "foo",
		Type:  "aws-ec2",
	}
	err := c.enableCredential(namespace.RootContext(nil), me)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	match := c.router.MatchingMount(namespace.RootContext(nil), "auth/foo/bar")
	if match != "auth/foo/" {
		t.Fatalf("missing mount, match: %q", match)
	}

	inmemSink := metrics.NewInmemSink(1000000*time.Hour, 2000000*time.Hour)
	conf := &CoreConfig{
		Physical:        c.physical,
		DisableMlock:    true,
		BuiltinRegistry: corehelpers.NewMockBuiltinRegistry(),
		MetricSink:      metricsutil.NewClusterMetricSink("test-cluster", inmemSink),
		MetricsHelper:   metricsutil.NewMetricsHelper(inmemSink, false),
	}
	c2, err := NewCore(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer c2.Shutdown()
	c2.credentialBackends["noop"] = func(context.Context, *logical.BackendConfig) (logical.Backend, error) {
		return &NoopBackend{
			BackendType: logical.TypeCredential,
		}, nil
	}
	for i, key := range keys {
		unseal, err := TestCoreUnseal(c2, key)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if i+1 == len(keys) && !unseal {
			t.Fatalf("should be unsealed")
		}
	}

	// Verify matching auth tables
	if !reflect.DeepEqual(c.auth, c2.auth) {
		t.Fatalf("mismatch: %v %v", c.auth, c2.auth)
	}
}

// Test that the local table actually gets populated as expected with local
// entries, and that upon reading the entries from both are recombined
// correctly
func TestCore_EnableCredential_Local(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	c.credentialBackends["noop"] = func(context.Context, *logical.BackendConfig) (logical.Backend, error) {
		return &NoopBackend{
			BackendType: logical.TypeCredential,
		}, nil
	}

	c.auth = &MountTable{
		Type: credentialTableType,
		Entries: []*MountEntry{
			{
				Table:            credentialTableType,
				Path:             "noop/",
				Type:             "noop",
				UUID:             "abcd",
				Accessor:         "noop-abcd",
				BackendAwareUUID: "abcde",
				NamespaceID:      namespace.RootNamespaceID,
				namespace:        namespace.RootNamespace,
			},
			{
				Table:            credentialTableType,
				Path:             "noop2/",
				Type:             "noop",
				UUID:             "bcde",
				Accessor:         "noop-bcde",
				BackendAwareUUID: "bcdea",
				NamespaceID:      namespace.RootNamespaceID,
				namespace:        namespace.RootNamespace,
			},
		},
	}

	// Both should set up successfully
	err := c.setupCredentials(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	rawLocal, err := c.barrier.Get(context.Background(), coreLocalAuthConfigPath)
	if err != nil {
		t.Fatal(err)
	}
	if rawLocal == nil {
		t.Fatal("expected non-nil local credential")
	}
	localCredentialTable := &MountTable{}
	if err := jsonutil.DecodeJSON(rawLocal.Value, localCredentialTable); err != nil {
		t.Fatal(err)
	}
	if len(localCredentialTable.Entries) > 0 {
		t.Fatalf("expected no entries in local credential table, got %#v", localCredentialTable)
	}

	c.auth.Entries[1].Local = true
	if err := c.persistAuth(context.Background(), c.auth, nil); err != nil {
		t.Fatal(err)
	}

	rawLocal, err = c.barrier.Get(context.Background(), coreLocalAuthConfigPath)
	if err != nil {
		t.Fatal(err)
	}
	if rawLocal == nil {
		t.Fatal("expected non-nil local credential")
	}
	localCredentialTable = &MountTable{}
	if err := jsonutil.DecodeJSON(rawLocal.Value, localCredentialTable); err != nil {
		t.Fatal(err)
	}
	if len(localCredentialTable.Entries) != 1 {
		t.Fatalf("expected one entry in local credential table, got %#v", localCredentialTable)
	}

	oldCredential := c.auth
	if err := c.loadCredentials(context.Background()); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(oldCredential, c.auth) {
		t.Fatalf("expected\n%#v\ngot\n%#v\n", oldCredential, c.auth)
	}

	if len(c.auth.Entries) != 2 {
		t.Fatalf("expected two credential entries, got %#v", localCredentialTable)
	}
}

func TestCore_EnableCredential_twice_409(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	c.credentialBackends["noop"] = func(context.Context, *logical.BackendConfig) (logical.Backend, error) {
		return &NoopBackend{
			BackendType: logical.TypeCredential,
		}, nil
	}

	me := &MountEntry{
		Table: credentialTableType,
		Path:  "foo",
		Type:  "noop",
	}
	err := c.enableCredential(namespace.RootContext(nil), me)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// 2nd should be a 409 error
	err2 := c.enableCredential(namespace.RootContext(nil), me)
	switch err2.(type) {
	case logical.HTTPCodedError:
		if err2.(logical.HTTPCodedError).Code() != 409 {
			t.Fatalf("invalid code given")
		}
	default:
		t.Fatalf("expected a different error type")
	}
}

func TestCore_EnableCredential_Token(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	me := &MountEntry{
		Table: credentialTableType,
		Path:  "foo",
		Type:  "token",
	}
	err := c.enableCredential(namespace.RootContext(nil), me)
	if err.Error() != "token credential backend cannot be instantiated" {
		t.Fatalf("err: %v", err)
	}
}

func TestCore_DisableCredential(t *testing.T) {
	c, keys, _ := TestCoreUnsealed(t)
	c.credentialBackends["noop"] = func(context.Context, *logical.BackendConfig) (logical.Backend, error) {
		return &NoopBackend{
			BackendType: logical.TypeCredential,
		}, nil
	}

	err := c.disableCredential(namespace.RootContext(nil), "foo")
	if err != nil && !strings.HasPrefix(err.Error(), "no matching mount") {
		t.Fatal(err)
	}

	me := &MountEntry{
		Table: credentialTableType,
		Path:  "foo",
		Type:  "noop",
	}
	err = c.enableCredential(namespace.RootContext(nil), me)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	err = c.disableCredential(namespace.RootContext(nil), "foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	match := c.router.MatchingMount(namespace.RootContext(nil), "auth/foo/bar")
	if match != "" {
		t.Fatalf("backend present")
	}

	inmemSink := metrics.NewInmemSink(1000000*time.Hour, 2000000*time.Hour)
	conf := &CoreConfig{
		Physical:        c.physical,
		DisableMlock:    true,
		BuiltinRegistry: corehelpers.NewMockBuiltinRegistry(),
		MetricSink:      metricsutil.NewClusterMetricSink("test-cluster", inmemSink),
		MetricsHelper:   metricsutil.NewMetricsHelper(inmemSink, false),
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
	if !reflect.DeepEqual(c.auth, c2.auth) {
		t.Fatalf("mismatch: %v %v", c.auth, c2.auth)
	}
}

func TestCore_DisableCredential_Protected(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	err := c.disableCredential(namespace.RootContext(nil), "token")
	if err.Error() != "token credential backend cannot be disabled" {
		t.Fatalf("err: %v", err)
	}
}

func TestCore_DisableCredential_Cleanup(t *testing.T) {
	noop := &NoopBackend{
		Login:       []string{"login"},
		BackendType: logical.TypeCredential,
	}
	c, _, _ := TestCoreUnsealed(t)
	c.credentialBackends["noop"] = func(context.Context, *logical.BackendConfig) (logical.Backend, error) {
		return noop, nil
	}

	me := &MountEntry{
		Table: credentialTableType,
		Path:  "foo",
		Type:  "noop",
	}
	err := c.enableCredential(namespace.RootContext(nil), me)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Store the view
	view := c.router.MatchingStorageByAPIPath(namespace.RootContext(nil), "auth/foo/")

	// Inject data
	se := &logical.StorageEntry{
		Key:   "plstodelete",
		Value: []byte("test"),
	}
	if err := view.Put(context.Background(), se); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Generate a new token auth
	noop.Response = &logical.Response{
		Auth: &logical.Auth{
			Policies: []string{"foo"},
		},
	}
	r := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "auth/foo/login",
	}
	resp, err := c.HandleRequest(namespace.RootContext(nil), r)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}

	// Disable should cleanup
	err = c.disableCredential(namespace.RootContext(nil), "foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Token should be revoked
	te, err := c.tokenStore.Lookup(namespace.RootContext(nil), resp.Auth.ClientToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if te != nil {
		t.Fatalf("bad: %#v", te)
	}

	// View should be empty
	out, err := logical.CollectKeys(context.Background(), view)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("bad: %#v", out)
	}
}

func TestDefaultAuthTable(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	table := c.defaultAuthTable()
	verifyDefaultAuthTable(t, table)
}

func verifyDefaultAuthTable(t *testing.T, table *MountTable) {
	if len(table.Entries) != 1 {
		t.Fatalf("bad: %v", table.Entries)
	}
	if table.Type != credentialTableType {
		t.Fatalf("bad: %v", *table)
	}
	for idx, entry := range table.Entries {
		switch idx {
		case 0:
			if entry.Path != "token/" {
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

func TestCore_CredentialInitialize(t *testing.T) {
	{
		backend := &InitializableBackend{
			&NoopBackend{
				BackendType: logical.TypeCredential,
			}, false,
		}

		c, _, _ := TestCoreUnsealed(t)
		c.credentialBackends["initable"] = func(context.Context, *logical.BackendConfig) (logical.Backend, error) {
			return backend, nil
		}

		me := &MountEntry{
			Table: credentialTableType,
			Path:  "foo/",
			Type:  "initable",
		}
		err := c.enableCredential(namespace.RootContext(nil), me)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		if !backend.isInitialized {
			t.Fatal("backend is not initialized")
		}
	}
	{
		backend := &InitializableBackend{
			&NoopBackend{
				BackendType: logical.TypeCredential,
			}, false,
		}

		c, _, _ := TestCoreUnsealed(t)
		c.credentialBackends["initable"] = func(context.Context, *logical.BackendConfig) (logical.Backend, error) {
			return backend, nil
		}

		c.auth = &MountTable{
			Type: credentialTableType,
			Entries: []*MountEntry{
				{
					Table:            credentialTableType,
					Path:             "foo/",
					Type:             "initable",
					UUID:             "abcd",
					Accessor:         "initable-abcd",
					BackendAwareUUID: "abcde",
					NamespaceID:      namespace.RootNamespaceID,
					namespace:        namespace.RootNamespace,
				},
			},
		}

		err := c.setupCredentials(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		// run the postUnseal funcs, so that the backend will be inited
		for _, f := range c.postUnsealFuncs {
			f()
		}

		if !backend.isInitialized {
			t.Fatal("backend is not initialized")
		}
	}
}

func remountCredentialFromRoot(c *Core, src, dst string, updateStorage bool) error {
	srcPathDetails := c.splitNamespaceAndMountFromPath("", src)
	dstPathDetails := c.splitNamespaceAndMountFromPath("", dst)
	return c.remountCredential(namespace.RootContext(nil), srcPathDetails, dstPathDetails, updateStorage)
}

func TestCore_RemountCredential(t *testing.T) {
	c, keys, _ := TestCoreUnsealed(t)
	c.credentialBackends["noop"] = func(context.Context, *logical.BackendConfig) (logical.Backend, error) {
		return &NoopBackend{
			BackendType: logical.TypeCredential,
		}, nil
	}

	me := &MountEntry{
		Table: credentialTableType,
		Path:  "foo",
		Type:  "noop",
	}
	err := c.enableCredential(namespace.RootContext(nil), me)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	match := c.router.MatchingMount(namespace.RootContext(nil), "auth/foo/bar")
	if match != "auth/foo/" {
		t.Fatalf("missing mount, match: %q", match)
	}

	err = remountCredentialFromRoot(c, "auth/foo", "auth/bar", true)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	match = c.router.MatchingMount(namespace.RootContext(nil), "auth/bar/baz")
	if match != "auth/bar/" {
		t.Fatalf("auth method not at new location, match: %q", match)
	}

	c.sealInternal()
	for i, key := range keys {
		unseal, err := TestCoreUnseal(c, key)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if i+1 == len(keys) && !unseal {
			t.Fatalf("should be unsealed")
		}
	}

	match = c.router.MatchingMount(namespace.RootContext(nil), "auth/bar/baz")
	if match != "auth/bar/" {
		t.Fatalf("auth method not at new location after unseal, match: %q", match)
	}
}

func TestCore_RemountCredential_Cleanup(t *testing.T) {
	noop := &NoopBackend{
		Login:       []string{"login"},
		BackendType: logical.TypeCredential,
	}
	c, _, _ := TestCoreUnsealed(t)
	c.credentialBackends["noop"] = func(context.Context, *logical.BackendConfig) (logical.Backend, error) {
		return noop, nil
	}

	me := &MountEntry{
		Table: credentialTableType,
		Path:  "foo",
		Type:  "noop",
	}
	err := c.enableCredential(namespace.RootContext(nil), me)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Store the view
	view := c.router.MatchingStorageByAPIPath(namespace.RootContext(nil), "auth/foo/")

	// Inject data
	se := &logical.StorageEntry{
		Key:   "plstodelete",
		Value: []byte("test"),
	}
	if err := view.Put(context.Background(), se); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Generate a new token auth
	noop.Response = &logical.Response{
		Auth: &logical.Auth{
			Policies: []string{"foo"},
		},
	}
	r := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "auth/foo/login",
	}
	resp, err := c.HandleRequest(namespace.RootContext(nil), r)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.Auth.ClientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}

	// Disable should cleanup
	err = remountCredentialFromRoot(c, "auth/foo", "auth/bar", true)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Token should be revoked
	te, err := c.tokenStore.Lookup(namespace.RootContext(nil), resp.Auth.ClientToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if te != nil {
		t.Fatalf("bad: %#v", te)
	}

	// View should be empty
	out, err := logical.CollectKeys(context.Background(), view)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(out) != 1 && out[0] != "plstokeep" {
		t.Fatalf("bad: %#v", out)
	}
}

func TestCore_RemountCredential_InvalidSource(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	err := remountCredentialFromRoot(c, "foo", "auth/bar", true)
	if err.Error() != `cannot remount non-auth mount "foo/"` {
		t.Fatalf("err: %v", err)
	}
}

func TestCore_RemountCredential_InvalidDestination(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	err := remountCredentialFromRoot(c, "auth/foo", "bar", true)
	if err.Error() != `cannot remount auth mount to non-auth mount "bar/"` {
		t.Fatalf("err: %v", err)
	}
}

func TestCore_RemountCredential_ProtectedSource(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	err := remountCredentialFromRoot(c, "auth/token", "auth/bar", true)
	if err.Error() != `cannot remount "auth/token/"` {
		t.Fatalf("err: %v", err)
	}
}

func TestCore_RemountCredential_ProtectedDestination(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	err := remountCredentialFromRoot(c, "auth/foo", "auth/token", true)
	if err.Error() != `cannot remount to "auth/token/"` {
		t.Fatalf("err: %v", err)
	}
}
