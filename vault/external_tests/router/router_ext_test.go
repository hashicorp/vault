// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package router

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func TestRouter_MountSubpath_Checks(t *testing.T) {
	testRouter_MountSubpath(t, []string{"a/abcd/123", "abcd/123"})
	testRouter_MountSubpath(t, []string{"abcd/123", "a/abcd/123"})
	testRouter_MountSubpath(t, []string{"a/abcd/123", "abcd/123"})
}

func testRouter_MountSubpath(t *testing.T, mountPoints []string) {
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	// Test auth
	authInput := &api.EnableAuthOptions{
		Type: "userpass",
	}
	for _, mp := range mountPoints {
		t.Logf("mounting %s", "auth/"+mp)
		var err error
		err = client.Sys().EnableAuthWithOptions("auth/"+mp, authInput)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
	}

	// Test secrets
	mountInput := &api.MountInput{
		Type: "pki",
	}
	for _, mp := range mountPoints {
		t.Logf("mounting %s", "s/"+mp)
		var err error
		err = client.Sys().Mount("s/"+mp, mountInput)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
	}

	cluster.EnsureCoresSealed(t)
	cluster.UnsealCores(t)
	t.Logf("Done: %#v", mountPoints)
}

func TestRouter_UnmountRollbackIsntFatal(t *testing.T) {
	cluster := minimal.NewTestSoloCluster(t, &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"noop": vault.NoopBackendRollbackErrFactory,
		},
	})
	client := cluster.Cores[0].Client

	if err := client.Sys().Mount("noop", &api.MountInput{
		Type: "noop",
	}); err != nil {
		t.Fatalf("failed to mount PKI: %v", err)
	}

	if _, err := client.Logical().Write("sys/plugins/reload/backend", map[string]interface{}{
		"mounts": "noop",
	}); err != nil {
		t.Fatalf("expected reload of noop with broken periodic func to succeed; got err=%v", err)
	}

	if _, err := client.Logical().Write("sys/remount", map[string]interface{}{
		"from": "noop",
		"to":   "noop-to",
	}); err != nil {
		t.Fatalf("expected remount of noop with broken periodic func to succeed; got err=%v", err)
	}

	cluster.EnsureCoresSealed(t)
	cluster.UnsealCores(t)
}

func TestWellKnownRedirect_HA(t *testing.T) {
	var records *[][]byte
	cluster := vault.NewTestCluster(t, &vault.CoreConfig{
		AuditBackends: map[string]audit.Factory{
			"noop": audit.NoopAuditFactory(&records),
		},
		DisablePerformanceStandby: true,
		LogicalBackends: map[string]logical.Factory{
			"noop": func(_ context.Context, _ *logical.BackendConfig) (logical.Backend, error) {
				return &vault.NoopBackend{
					RequestHandler: func(context.Context, *logical.Request) (*logical.Response, error) {
						// Return something for any request
						return &logical.Response{
							Data: map[string]interface{}{
								"good": "very",
							},
						}, nil
					},
				}, nil
			},
		},
	}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	testhelpers.WaitForActiveNodeAndStandbys(t, cluster)
	active := testhelpers.DeriveActiveCore(t, cluster)
	standbys := testhelpers.DeriveStandbyCores(t, cluster)
	standby := standbys[0].Client

	if err := active.Client.Sys().EnableAuditWithOptions("noop", &api.EnableAuditOptions{
		Type: "noop",
	}); err != nil {
		t.Fatalf("failed to enable audit: %v", err)
	}

	if err := active.Client.Sys().Mount("noop", &api.MountInput{
		Type: "noop",
	}); err != nil {
		t.Fatalf("failed to mount PKI: %v", err)
	}

	resp, err := active.Client.Logical().Read("sys/mounts")
	if err != nil {
		t.Fatalf("failed to fetch new mount: %v", err)
	}
	var mountUUID string
	for k, m := range resp.Data {
		if k == "noop/" {
			mountUUID = m.(map[string]interface{})["uuid"].(string)
			break
		}
	}

	if err := active.Core.WellKnownRedirects.TryRegister(context.Background(), active.Core, mountUUID, "foo", "bar"); err != nil {
		t.Fatal(err)
	}

	standby.SetCheckRedirect(nil)
	resp2, err := standby.RawRequest(standby.NewRequest(http.MethodGet, "/.well-known/foo/baz"))
	if err != nil {
		t.Fatal(err)
	} else if resp2.StatusCode != http.StatusOK {
		t.Fatal("did not get expected response from noop backend after redirect")
	}

	if len(*records) < 2 {
		t.Fatal("audit entries not populated")
	} else {
		rs := *records
		// Make sure RequestURI is present in the redirect audit entries
		if !strings.Contains(string(rs[len(rs)-1]), "request_uri\":\"/.well-known/foo/baz") ||
			!strings.Contains(string(rs[len(rs)-2]), "request_uri\":\"/.well-known/foo/baz") {
			t.Fatal("did not find request_uri in audit entries")
		}
	}
}
