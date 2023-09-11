// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"encoding/base64"
	"errors"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/armon/go-metrics"
	logicalKv "github.com/hashicorp/vault-plugin-secrets-kv"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

func TestCoreMetrics_KvSecretGauge(t *testing.T) {
	// Use the real KV implementation instead of Passthrough
	AddTestLogicalBackend("kv", logicalKv.Factory)
	// Clean up for the next test-- is there a better way?
	defer func() {
		delete(testLogicalBackends, "kv")
	}()
	core, _, root := TestCoreUnsealed(t)

	testMounts := []struct {
		Path          string
		Type          string
		Version       string
		ExpectedCount int
	}{
		{"secret/", "kv", "2", 0},
		{"secret1/", "kv", "1", 3},
		{"secret2/", "kv", "1", 0},
		{"secret3/", "kv", "2", 4},
		{"prefix/secret3/", "kv", "2", 0},
		{"prefix/secret4/", "kv", "2", 5},
		{"generic/", "generic", "1", 3},
	}
	ctx := namespace.RootContext(nil)

	// skip 0, secret/ is already mounted
	for _, tm := range testMounts[1:] {
		me := &MountEntry{
			Table:   mountTableType,
			Path:    sanitizePath(tm.Path),
			Type:    tm.Type,
			Options: map[string]string{"version": tm.Version},
		}
		err := core.mount(ctx, me)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
	}

	v1secrets := []string{
		"secret1/a", // 3
		"secret1/b",
		"secret1/c/d",
		"generic/a",
		"generic/b",
		"generic/c",
	}
	v2secrets := []string{
		"secret3/data/a", // 4
		"secret3/data/b",
		"secret3/data/c/d",
		"secret3/data/c/e",
		"prefix/secret4/data/a/secret", // 5
		"prefix/secret4/data/a/secret2",
		"prefix/secret4/data/a/b/c/secret",
		"prefix/secret4/data/a/b/c/secret2",
		"prefix/secret4/data/a/b/c/d/secret3",
	}
	for _, p := range v1secrets {
		req := logical.TestRequest(t, logical.CreateOperation, p)
		req.Data["foo"] = "bar"
		req.ClientToken = root
		resp, err := core.HandleRequest(ctx, req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if resp != nil {
			t.Fatalf("bad: %#v", resp)
		}
	}
	for _, p := range v2secrets {
		for i := 0; i < 50; i++ {
			req := logical.TestRequest(t, logical.CreateOperation, p)
			req.Data["data"] = map[string]interface{}{"foo": "bar"}
			req.ClientToken = root
			resp, err := core.HandleRequest(ctx, req)
			if err != nil {
				if errors.Is(err, logical.ErrInvalidRequest) {
					// Handle scenario where KVv2 upgrade is ongoing
					time.Sleep(100 * time.Millisecond)
					continue
				}
				t.Fatalf("err: %v", err)
			}
			if resp.Error() != nil {
				t.Fatalf("bad: %#v", resp)
			}
			break
		}
	}

	values, err := core.kvSecretGaugeCollector(ctx)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(values) != len(testMounts) {
		t.Errorf("Got %v values but expected %v mounts", len(values), len(testMounts))
	}

	for _, glv := range values {
		mountPoint := ""
		for _, l := range glv.Labels {
			if l.Name == "mount_point" {
				mountPoint = l.Value
			} else if l.Name == "namespace" {
				if l.Value != "root" {
					t.Errorf("Namespace is %v, not root", l.Value)
				}
			} else {
				t.Errorf("Unexpected label %v", l.Name)
			}
		}
		if mountPoint == "" {
			t.Errorf("No mount point in labels %v", glv.Labels)
			continue
		}
		found := false
		for _, tm := range testMounts {
			if tm.Path == mountPoint {
				found = true
				if glv.Value != float32(tm.ExpectedCount) {
					t.Errorf("Mount %v reported %v, not %v",
						tm.Path, glv.Value, tm.ExpectedCount)
				}
				break
			}
		}
		if !found {
			t.Errorf("Unexpected mount point %v", mountPoint)
		}
	}
}

func TestCoreMetrics_KvSecretGauge_BadPath(t *testing.T) {
	// Use the real KV implementation instead of Passthrough
	AddTestLogicalBackend("kv", logicalKv.Factory)
	// Clean up for the next test.
	defer func() {
		delete(testLogicalBackends, "kv")
	}()
	core, _, _ := TestCoreUnsealed(t)

	me := &MountEntry{
		Table:   mountTableType,
		Path:    sanitizePath("kv1"),
		Type:    "kv",
		Options: map[string]string{"version": "1"},
	}
	ctx := namespace.RootContext(nil)
	err := core.mount(ctx, me)
	if err != nil {
		t.Fatalf("mount error: %v", err)
	}

	// I don't think there's any remaining way to create a zero-length
	// key via the API, so we'll fake it by talking to the storage layer directly.
	fake_entry := &logical.StorageEntry{
		Key:   "logical/" + me.UUID + "/foo/",
		Value: []byte{1},
	}
	err = core.barrier.Put(ctx, fake_entry)
	if err != nil {
		t.Fatalf("put error: %v", err)
	}

	values, err := core.kvSecretGaugeCollector(ctx)
	if err != nil {
		t.Fatalf("collector error: %v", err)
	}
	t.Logf("Values: %v", values)
	found := false
	var count float32 = -1
	for _, glv := range values {
		for _, l := range glv.Labels {
			if l.Name == "mount_point" && l.Value == "kv1/" {
				found = true
				count = glv.Value
				break
			}
		}
	}
	if found {
		if count != 1.0 {
			t.Errorf("bad secret count for kv1/")
		}
	} else {
		t.Errorf("no secret count for kv1/")
	}
}

func TestCoreMetrics_KvSecretGaugeError(t *testing.T) {
	core, _, _, sink := TestCoreUnsealedWithMetrics(t)
	ctx := namespace.RootContext(nil)

	badKvMount := &kvMount{
		Namespace:  namespace.RootNamespace,
		MountPoint: "bad/path",
		Version:    "1",
		NumSecrets: 0,
	}

	core.walkKvMountSecrets(ctx, badKvMount)

	intervals := sink.Data()
	// Test crossed an interval boundary, don't try to deal with it.
	if len(intervals) > 1 {
		t.Skip("Detected interval crossing.")
	}

	// Should be an error
	keyPrefix := "metrics.collection.error"
	var counter *metrics.SampledValue = nil

	for _, c := range intervals[0].Counters {
		if strings.HasPrefix(c.Name, keyPrefix) {
			counter = &c
			break
		}
	}
	if counter == nil {
		t.Fatal("No metrics.collection.error counter found.")
	}
	if counter.Count != 1 {
		t.Errorf("Counter number of samples %v is not 1.", counter.Count)
	}
}

func metricLabelsMatch(t *testing.T, actual []metrics.Label, expected map[string]string) {
	t.Helper()

	if len(actual) != len(expected) {
		t.Errorf("Expected %v labels, got %v: %v", len(expected), len(actual), actual)
	}

	for _, l := range actual {
		if v, ok := expected[l.Name]; ok {
			if v != l.Value {
				t.Errorf("Mismatched value %v=%v, expected %v", l.Name, l.Value, v)
			}
		} else {
			t.Errorf("Unexpected label %v", l.Name)
		}
	}
}

func TestCoreMetrics_EntityGauges(t *testing.T) {
	ctx := namespace.RootContext(nil)
	is, ghAccessor, upAccessor, core := testIdentityStoreWithGithubUserpassAuth(ctx, t)

	// Create an entity
	alias1 := &logical.Alias{
		MountType:     "github",
		MountAccessor: ghAccessor,
		Name:          "githubuser",
	}

	entity, _, err := is.CreateOrFetchEntity(ctx, alias1)
	if err != nil {
		t.Fatal(err)
	}

	// Create a second alias for the same entity
	registerReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "entity-alias",
		Data: map[string]interface{}{
			"name":           "userpassuser",
			"canonical_id":   entity.ID,
			"mount_accessor": upAccessor,
		},
	}
	resp, err := is.HandleRequest(ctx, registerReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	glv, err := core.entityGaugeCollector(ctx)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if len(glv) != 1 {
		t.Fatalf("Wrong number of gauges %v, expected %v", len(glv), 1)
	}

	if glv[0].Value != 1.0 {
		t.Errorf("Entity count %v, expected %v", glv[0].Value, 1.0)
	}

	metricLabelsMatch(t, glv[0].Labels,
		map[string]string{
			"namespace": "root",
		})

	glv, err = core.entityGaugeCollectorByMount(ctx)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if len(glv) != 2 {
		t.Fatalf("Wrong number of gauges %v, expected %v", len(glv), 1)
	}

	if glv[0].Value != 1.0 {
		t.Errorf("Alias count %v, expected %v", glv[0].Value, 1.0)
	}

	if glv[1].Value != 1.0 {
		t.Errorf("Alias count %v, expected %v", glv[0].Value, 1.0)
	}

	// Sort both metrics.Label slices by Name, causing the Label
	// with Name auth_method to be first in both arrays
	sort.Slice(glv[0].Labels, func(i, j int) bool { return glv[0].Labels[i].Name < glv[0].Labels[j].Name })
	sort.Slice(glv[1].Labels, func(i, j int) bool { return glv[1].Labels[i].Name < glv[1].Labels[j].Name })

	// Sort the GaugeLabelValues slice by the Value of the first metric,
	// in this case auth_method, in each metrics.Label slice
	sort.Slice(glv, func(i, j int) bool { return glv[i].Labels[0].Value < glv[j].Labels[0].Value })

	metricLabelsMatch(t, glv[0].Labels,
		map[string]string{
			"namespace":   "root",
			"auth_method": "github",
			"mount_point": "auth/github/",
		})

	metricLabelsMatch(t, glv[1].Labels,
		map[string]string{
			"namespace":   "root",
			"auth_method": "userpass",
			"mount_point": "auth/userpass/",
		})
}

// TestCoreMetrics_AvailablePolicies tests the that available metrics are getting correctly collected when the availablePoliciesGaugeCollector function is invoked
func TestCoreMetrics_AvailablePolicies(t *testing.T) {
	aclPolicy := map[string]interface{}{
		"policy": base64.StdEncoding.EncodeToString([]byte(`path "ns1/secret/foo/*" {
    capabilities = ["create", "read", "update", "delete", "list"]
}`)),
		"name": "secret",
	}

	type pathPolicy struct {
		Path   string
		Policy map[string]interface{}
	}

	tests := map[string]struct {
		Policies       []pathPolicy
		ExpectedValues map[string]float32
	}{
		"single acl": {
			Policies: []pathPolicy{
				{
					"sys/policy/secret", aclPolicy,
				},
			},
			ExpectedValues: map[string]float32{
				// The "default" policy will always be included
				"acl": 2,
				"egp": 0,
				"rgp": 0,
			},
		},
		"multiple acl": {
			Policies: []pathPolicy{
				{
					"sys/policy/secret", aclPolicy,
				},
				{
					"sys/policy/secret2", aclPolicy,
				},
			},
			ExpectedValues: map[string]float32{
				// The "default" policy will always be included
				"acl": 3,
				"egp": 0,
				"rgp": 0,
			},
		},
	}

	for name, tst := range tests {
		t.Run(name, func(t *testing.T) {
			core, _, root := TestCoreUnsealed(t)

			ctxRoot := namespace.RootContext(context.Background())

			// Create policies
			for _, p := range tst.Policies {
				req := logical.TestRequest(t, logical.UpdateOperation, p.Path)
				req.Data = p.Policy
				req.ClientToken = root

				resp, err := core.HandleRequest(ctxRoot, req)
				if err != nil {
					t.Fatalf("err: %v", err)
				}
				if resp != nil {
					logger.Info("expected nil response", resp)
					t.Fatalf("expected nil response")
				}
			}

			gValues, err := core.configuredPoliciesGaugeCollector(ctxRoot)
			if err != nil {
				t.Fatalf("err: %v", err)
			}

			// Check the metrics values match the expected values
			mgValues := make(map[string]float32, len(gValues))
			for _, v := range gValues {
				mgValues[v.Labels[0].Value] = v.Value
			}

			assert.EqualValues(t, tst.ExpectedValues, mgValues)
		})
	}
}
