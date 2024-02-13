// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package quotas

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/builtin/logical/pki"
	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
	"github.com/hashicorp/vault/sdk/helper/testhelpers/schema"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
)

const (
	testLookupOnlyPolicy = `
path "/auth/token/lookup" {
	capabilities = [ "create", "update"]
}
`
)

var coreConfig = &vault.CoreConfig{
	LogicalBackends: map[string]logical.Factory{
		"pki": pki.Factory,
	},
	CredentialBackends: map[string]logical.Factory{
		"userpass": userpass.Factory,
	},
}

func setupMounts(t *testing.T, client *api.Client) {
	t.Helper()

	err := client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("auth/userpass/users/foo", map[string]interface{}{
		"password": "bar",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = client.Sys().Mount("pki", &api.MountInput{
		Type: "pki",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("pki/root/generate/internal", map[string]interface{}{
		"common_name": "testvault.com",
		"ttl":         "200h",
		"ip_sans":     "127.0.0.1",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("pki/roles/test", map[string]interface{}{
		"require_cn":       false,
		"allowed_domains":  "testvault.com",
		"allow_subdomains": true,
		"max_ttl":          "2h",
		"generate_lease":   true,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func teardownMounts(t *testing.T, client *api.Client) {
	t.Helper()
	if err := client.Sys().Unmount("pki"); err != nil {
		t.Fatal(err)
	}
	if err := client.Sys().DisableAuth("userpass"); err != nil {
		t.Fatal(err)
	}
	if err := client.Sys().DisableAuth("approle"); err != nil {
		t.Fatal(err)
	}
}

func testRPS(reqFunc func(numSuccess, numFail *atomic.Int32), d time.Duration) (int32, int32, time.Duration) {
	numSuccess := atomic.NewInt32(0)
	numFail := atomic.NewInt32(0)

	start := time.Now()
	end := start.Add(d)
	for time.Now().Before(end) {
		reqFunc(numSuccess, numFail)
	}

	return numSuccess.Load(), numFail.Load(), time.Since(start)
}

func waitForRemovalOrTimeout(c *api.Client, path string, tick, to time.Duration) error {
	ticker := time.Tick(tick)
	timeout := time.After(to)

	// wait for the resource to be removed
	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout exceeding waiting for resource to be deleted: %s", path)

		case <-ticker:
			resp, err := c.Logical().Read(path)
			if err != nil {
				return err
			}

			if resp == nil {
				return nil
			}
		}
	}
}

func TestQuotas_RateLimit_DupName(t *testing.T) {
	conf, opts := teststorage.ClusterSetup(coreConfig, nil, nil)
	opts.NoDefaultQuotas = true
	opts.RequestResponseCallback = schema.ResponseValidatingCallback(t)
	cluster := vault.NewTestCluster(t, conf, opts)
	cluster.Start()
	defer cluster.Cleanup()
	core := cluster.Cores[0].Core
	client := cluster.Cores[0].Client
	vault.TestWaitActive(t, core)

	// create a rate limit quota w/ 'secret' path
	_, err := client.Logical().Write("sys/quotas/rate-limit/secret-rlq", map[string]interface{}{
		"rate": 7.7,
		"path": "secret",
	})
	require.NoError(t, err)

	s, err := client.Logical().Read("sys/quotas/rate-limit/secret-rlq")
	require.NoError(t, err)
	require.NotEmpty(t, s.Data)

	// create a rate limit quota w/ empty path (same name)
	_, err = client.Logical().Write("sys/quotas/rate-limit/secret-rlq", map[string]interface{}{
		"rate": 7.7,
		"path": "",
	})
	require.NoError(t, err)

	// list again and verify that only 1 item is returned
	s, err = client.Logical().List("sys/quotas/rate-limit")
	require.NoError(t, err)

	require.Len(t, s.Data, 1, "incorrect number of quotas")
}

func TestQuotas_RateLimit_DupPath(t *testing.T) {
	conf, opts := teststorage.ClusterSetup(coreConfig, nil, nil)
	opts.NoDefaultQuotas = true
	opts.RequestResponseCallback = schema.ResponseValidatingCallback(t)
	cluster := vault.NewTestCluster(t, conf, opts)
	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0].Core
	client := cluster.Cores[0].Client
	vault.TestWaitActive(t, core)
	// create a global rate limit quota
	_, err := client.Logical().Write("sys/quotas/rate-limit/global-rlq", map[string]interface{}{
		"rate": 10,
		"path": "",
	})
	require.NoError(t, err)

	// create a rate limit quota w/ 'secret' path
	_, err = client.Logical().Write("sys/quotas/rate-limit/secret-rlq", map[string]interface{}{
		"rate": 7.7,
		"path": "secret",
	})
	require.NoError(t, err)

	s, err := client.Logical().Read("sys/quotas/rate-limit/secret-rlq")
	require.NoError(t, err)
	require.NotEmpty(t, s.Data)

	// create a rate limit quota w/ empty path (same name)
	_, err = client.Logical().Write("sys/quotas/rate-limit/secret-rlq", map[string]interface{}{
		"rate": 7.7,
		"path": "",
	})

	if err == nil {
		t.Fatal("Duplicated paths were accepted")
	}
}

func TestQuotas_RateLimitQuota_ExemptPaths(t *testing.T) {
	conf, opts := teststorage.ClusterSetup(coreConfig, nil, nil)
	opts.NoDefaultQuotas = true
	opts.RequestResponseCallback = schema.ResponseValidatingCallback(t)
	cluster := vault.NewTestCluster(t, conf, opts)
	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0].Core
	client := cluster.Cores[0].Client
	vault.TestWaitActive(t, core)

	_, err := client.Logical().Write("sys/quotas/rate-limit/rlq", map[string]interface{}{
		"rate": 7.7,
	})
	require.NoError(t, err)

	// ensure exempt paths are not empty by default
	resp, err := client.Logical().Read("sys/quotas/config")
	require.NoError(t, err)
	require.NotEmpty(t, resp.Data["rate_limit_exempt_paths"].([]interface{}), "expected no exempt paths by default")

	reqFunc := func(numSuccess, numFail *atomic.Int32) {
		_, err := client.Logical().Read("sys/quotas/rate-limit/rlq")

		if err != nil {
			numFail.Add(1)
		} else {
			numSuccess.Add(1)
		}
	}

	numSuccess, numFail, elapsed := testRPS(reqFunc, 5*time.Second)
	ideal := 8 + (7.7 * float64(elapsed) / float64(time.Second))
	want := int32(ideal + 1)
	require.NotZerof(t, numFail, "expected some requests to fail; numSuccess: %d, elapsed: %d", numSuccess, elapsed)
	require.LessOrEqualf(t, numSuccess, want, "too many successful requests;numSuccess: %d, numFail: %d, elapsed: %d", numSuccess, numFail, elapsed)

	// allow time (1s) for rate limit to refill before updating the quota config
	time.Sleep(time.Second)

	_, err = client.Logical().Write("sys/quotas/config", map[string]interface{}{
		"rate_limit_exempt_paths": []string{"sys/quotas/rate-limit"},
	})
	require.NoError(t, err)

	// all requests should success
	numSuccess, numFail, _ = testRPS(reqFunc, 5*time.Second)
	require.NotZero(t, numSuccess)
	require.Zero(t, numFail)
}

func TestQuotas_RateLimitQuota_DefaultExemptPaths(t *testing.T) {
	conf, opts := teststorage.ClusterSetup(coreConfig, nil, nil)
	opts.NoDefaultQuotas = true
	opts.RequestResponseCallback = schema.ResponseValidatingCallback(t)
	cluster := vault.NewTestCluster(t, conf, opts)
	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0].Core
	client := cluster.Cores[0].Client
	vault.TestWaitActive(t, core)

	_, err := client.Logical().Write("sys/quotas/rate-limit/rlq", map[string]interface{}{
		"rate": 1,
	})
	require.NoError(t, err)

	resp, err := client.Logical().Read("sys/health")
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)

	// The second sys/health call should not fail as /v1/sys/health is
	// part of the default exempt paths
	resp, err = client.Logical().Read("sys/health")
	require.NoError(t, err)
	// If the response is nil, then we are being rate limited
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
}

func TestQuotas_RateLimitQuota_Mount(t *testing.T) {
	conf, opts := teststorage.ClusterSetup(coreConfig, nil, nil)
	cluster := vault.NewTestCluster(t, conf, opts)
	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0].Core
	client := cluster.Cores[0].Client
	vault.TestWaitActive(t, core)

	err := client.Sys().Mount("pki", &api.MountInput{
		Type: "pki",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("pki/root/generate/internal", map[string]interface{}{
		"common_name": "testvault.com",
		"ttl":         "200h",
		"ip_sans":     "127.0.0.1",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("pki/roles/test", map[string]interface{}{
		"require_cn":       false,
		"allowed_domains":  "testvault.com",
		"allow_subdomains": true,
		"max_ttl":          "2h",
		"generate_lease":   true,
	})
	if err != nil {
		t.Fatal(err)
	}

	reqFunc := func(numSuccess, numFail *atomic.Int32) {
		_, err := client.Logical().Read("pki/cert/ca_chain")

		if err != nil {
			numFail.Add(1)
		} else {
			numSuccess.Add(1)
		}
	}

	// Create a rate limit quota with a low RPS of 7.7, which means we can process
	// ⌈7.7⌉*2 requests in the span of roughly a second -- 8 initially, followed
	// by a refill rate of 7.7 per-second.
	_, err = client.Logical().Write("sys/quotas/rate-limit/rlq", map[string]interface{}{
		"rate": 7.7,
		"path": "pki/",
	})
	if err != nil {
		t.Fatal(err)
	}

	numSuccess, numFail, elapsed := testRPS(reqFunc, 5*time.Second)

	// evaluate the ideal RPS as (ceil(RPS) + (RPS * totalSeconds))
	ideal := 8 + (7.7 * float64(elapsed) / float64(time.Second))

	// ensure there were some failed requests
	if numFail == 0 {
		t.Fatalf("expected some requests to fail; numSuccess: %d, numFail: %d, elapsed: %d", numSuccess, numFail, elapsed)
	}

	// ensure that we should never get more requests than allowed
	if want := int32(ideal + 1); numSuccess > want {
		t.Fatalf("too many successful requests; want: %d, numSuccess: %d, numFail: %d, elapsed: %d", want, numSuccess, numFail, elapsed)
	}

	// update the rate limit quota with a high RPS such that no requests should fail
	_, err = client.Logical().Write("sys/quotas/rate-limit/rlq", map[string]interface{}{
		"rate": 10000.0,
		"path": "pki/",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, numFail, _ = testRPS(reqFunc, 5*time.Second)
	if numFail > 0 {
		t.Fatalf("unexpected number of failed requests: %d", numFail)
	}
}

func TestQuotas_RateLimitQuota_MountPrecedence(t *testing.T) {
	conf, opts := teststorage.ClusterSetup(coreConfig, nil, nil)
	opts.NoDefaultQuotas = true
	cluster := vault.NewTestCluster(t, conf, opts)
	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0].Core
	client := cluster.Cores[0].Client

	vault.TestWaitActive(t, core)

	// create PKI mount
	err := client.Sys().Mount("pki", &api.MountInput{
		Type: "pki",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("pki/root/generate/internal", map[string]interface{}{
		"common_name": "testvault.com",
		"ttl":         "200h",
		"ip_sans":     "127.0.0.1",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("pki/roles/test", map[string]interface{}{
		"require_cn":       false,
		"allowed_domains":  "testvault.com",
		"allow_subdomains": true,
		"max_ttl":          "2h",
		"generate_lease":   true,
	})
	if err != nil {
		t.Fatal(err)
	}

	// create a root rate limit quota
	_, err = client.Logical().Write("sys/quotas/rate-limit/root-rlq", map[string]interface{}{
		"name": "root-rlq",
		"rate": 14.7,
	})
	if err != nil {
		t.Fatal(err)
	}

	// create a mount rate limit quota with a lower RPS than the root rate limit quota
	_, err = client.Logical().Write("sys/quotas/rate-limit/mount-rlq", map[string]interface{}{
		"name": "mount-rlq",
		"rate": 7.7,
		"path": "pki/",
	})
	if err != nil {
		t.Fatal(err)
	}

	// ensure mount rate limit quota takes precedence over root rate limit quota
	reqFunc := func(numSuccess, numFail *atomic.Int32) {
		_, err := client.Logical().Read("pki/cert/ca_chain")

		if err != nil {
			numFail.Add(1)
		} else {
			numSuccess.Add(1)
		}
	}

	// ensure mount rate limit quota takes precedence over root rate limit quota
	numSuccess, numFail, elapsed := testRPS(reqFunc, 5*time.Second)

	// evaluate the ideal RPS as (ceil(RPS) + (RPS * totalSeconds))
	ideal := 8 + (7.7 * float64(elapsed) / float64(time.Second))

	// ensure there were some failed requests
	if numFail == 0 {
		t.Fatalf("expected some requests to fail; numSuccess: %d, numFail: %d, elapsed: %d", numSuccess, numFail, elapsed)
	}

	// ensure that we should never get more requests than allowed
	if want := int32(ideal + 1); numSuccess > want {
		t.Fatalf("too many successful requests; want: %d, numSuccess: %d, numFail: %d, elapsed: %d", want, numSuccess, numFail, elapsed)
	}
}

func TestQuotas_RateLimitQuota(t *testing.T) {
	conf, opts := teststorage.ClusterSetup(coreConfig, nil, nil)
	opts.NoDefaultQuotas = true
	cluster := vault.NewTestCluster(t, conf, opts)
	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0].Core
	client := cluster.Cores[0].Client

	vault.TestWaitActive(t, core)

	err := client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("auth/userpass/users/foo", map[string]interface{}{
		"password": "bar",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a rate limit quota with a low RPS of 7.7, which means we can process
	// ⌈7.7⌉*2 requests in the span of roughly a second -- 8 initially, followed
	// by a refill rate of 7.7 per-second.
	_, err = client.Logical().Write("sys/quotas/rate-limit/rlq", map[string]interface{}{
		"rate": 7.7,
	})
	if err != nil {
		t.Fatal(err)
	}

	reqFunc := func(numSuccess, numFail *atomic.Int32) {
		_, err := client.Logical().Read("sys/quotas/rate-limit/rlq")

		if err != nil {
			numFail.Add(1)
		} else {
			numSuccess.Add(1)
		}
	}

	numSuccess, numFail, elapsed := testRPS(reqFunc, 5*time.Second)

	// evaluate the ideal RPS as (ceil(RPS) + (RPS * totalSeconds))
	ideal := 8 + (7.7 * float64(elapsed) / float64(time.Second))

	// ensure there were some failed requests
	if numFail == 0 {
		t.Fatalf("expected some requests to fail; numSuccess: %d, numFail: %d, elapsed: %d", numSuccess, numFail, elapsed)
	}

	// ensure that we should never get more requests than allowed
	if want := int32(ideal + 1); numSuccess > want {
		t.Fatalf("too many successful requests; want: %d, numSuccess: %d, numFail: %d, elapsed: %d", want, numSuccess, numFail, elapsed)
	}

	// allow time (1s) for rate limit to refill before updating the quota
	time.Sleep(time.Second)

	// update the rate limit quota with a high RPS such that no requests should fail
	_, err = client.Logical().Write("sys/quotas/rate-limit/rlq", map[string]interface{}{
		"rate": 10000.0,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, numFail, _ = testRPS(reqFunc, 5*time.Second)
	if numFail > 0 {
		t.Fatalf("unexpected number of failed requests: %d", numFail)
	}
}
