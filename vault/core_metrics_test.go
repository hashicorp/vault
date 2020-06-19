package vault

import (
	"testing"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

func TestCoreMetrics_KvSecretGauge(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)

	// I think I can't test with the real kv-v2 because it's
	// a plugin. But we can fake it by using "metadata" as
	// part of the secret path with a V1 backend.
	testMounts := []struct {
		Path          string
		Version       string
		ExpectedCount int
	}{
		{"secret/", "2", 0},
		{"secret1/", "1", 3},
		{"secret2/", "1", 0},
		{"secret3/", "2", 4},
		{"prefix/secret3/", "2", 0},
		{"prefix/secret4/", "2", 5},
	}
	ctx := namespace.RootContext(nil)

	// skip 0, secret/ is already mounted
	for _, tm := range testMounts[1:] {
		me := &MountEntry{
			Table:   mountTableType,
			Path:    sanitizeMountPath(tm.Path),
			Type:    "kv",
			Options: map[string]string{"version": tm.Version},
		}
		err := core.mount(ctx, me)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
	}

	secrets := []string{
		"secret1/a", // 3
		"secret1/b",
		"secret1/c/d",
		"secret3/metadata/a", // 4
		"secret3/metadata/b",
		"secret3/metadata/c/d",
		"secret3/metadata/c/e",
		"prefix/secret4/metadata/a/secret", // 5
		"prefix/secret4/metadata/a/secret2",
		"prefix/secret4/metadata/a/b/c/secret",
		"prefix/secret4/metadata/a/b/c/secret2",
		"prefix/secret4/metadata/a/b/c/d/secret3",
	}
	for _, p := range secrets {
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
	is, ghAccessor, core := testIdentityStoreWithGithubAuth(ctx, t)

	// Create an entity
	alias1 := &logical.Alias{
		MountType:     "github",
		MountAccessor: ghAccessor,
		Name:          "githubuser",
	}

	entity, err := is.CreateOrFetchEntity(ctx, alias1)
	if err != nil {
		t.Fatal(err)
	}

	// Create a second alias for the same entity
	registerReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "entity-alias",
		Data: map[string]interface{}{
			"name":           "githubuser2",
			"canonical_id":   entity.ID,
			"mount_accessor": ghAccessor,
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

	if len(glv) != 1 {
		t.Fatalf("Wrong number of gauges %v, expected %v", len(glv), 1)
	}

	if glv[0].Value != 2.0 {
		t.Errorf("Alias count %v, expected %v", glv[0].Value, 2.0)
	}

	metricLabelsMatch(t, glv[0].Labels,
		map[string]string{
			"namespace":   "root",
			"auth_method": "github",
			"mount_point": "auth/github/",
		})
}
