// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package quotas

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/builtin/logical/pki"
	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/helper/testhelpers"
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

// TestQuotas_RateLimitQuota_GroupByConfig tests the validations imposed on the group_by and secondary_rate fields
func TestQuotas_RateLimitQuota_GroupByConfig(t *testing.T) {
	conf, opts := teststorage.ClusterSetup(coreConfig, nil, nil)
	opts.NoDefaultQuotas = true
	cluster := vault.NewTestCluster(t, conf, opts)
	cluster.Start()
	t.Cleanup(cluster.Cleanup)

	core := cluster.Cores[0].Core
	client := cluster.Cores[0].Client

	vault.TestWaitActive(t, core)

	testCases := map[string]struct {
		reqData              map[string]interface{}
		expectedErr          string
		expectedReadContains map[string]interface{}
		enterpriseOnly       bool
		CeOnly               bool
	}{
		"group_by_defaults_to_ip": {
			reqData: map[string]interface{}{
				"rate": 100,
			},
			expectedReadContains: map[string]interface{}{
				"group_by":       "ip",
				"secondary_rate": json.Number("0"),
			},
		},
		"explicitly_empty_group_by_defaults_to_ip": {
			reqData: map[string]interface{}{
				"rate":     100,
				"group_by": "",
			},
			expectedReadContains: map[string]interface{}{
				"group_by":       "ip",
				"secondary_rate": json.Number("0"),
			},
		},
		"explicit_group_by_ip_allowed_in_ce": {
			reqData: map[string]interface{}{
				"rate":     100,
				"group_by": "ip",
			},
			expectedReadContains: map[string]interface{}{
				"group_by":       "ip",
				"secondary_rate": json.Number("0"),
			},
		},
		"invalid_group_by": {
			reqData: map[string]interface{}{
				"rate":     100,
				"group_by": "invalid",
			},
			expectedErr: `invalid grouping mode "invalid"`,
		},
		"group_by_none_not_allowed_in_ce": {
			CeOnly: true,
			reqData: map[string]interface{}{
				"rate":     100,
				"group_by": "none",
			},
			expectedErr: `grouping mode "none" is only available in Vault Enterprise`,
		},
		"group_by_entity_then_none_not_allowed_in_ce": {
			CeOnly: true,
			reqData: map[string]interface{}{
				"rate":     100,
				"group_by": "entity_then_none",
			},
			expectedErr: `grouping mode "entity_then_none" is only available in Vault Enterprise`,
		},
		"group_by_entity_then_ip_not_allowed_in_ce": {
			CeOnly: true,
			reqData: map[string]interface{}{
				"rate":     100,
				"group_by": "entity_then_ip",
			},
			expectedErr: `grouping mode "entity_then_ip" is only available in Vault Enterprise`,
		},
		"secondary_rate_invalid_with_group_by_none": {
			enterpriseOnly: true,
			reqData: map[string]interface{}{
				"rate":           100,
				"secondary_rate": 1,
				"group_by":       "none",
			},
			expectedErr: "secondary rate is only valid when using entity-based grouping",
		},
		"secondary_rate_invalid_with_group_by_ip": {
			enterpriseOnly: true,
			reqData: map[string]interface{}{
				"rate":           100,
				"secondary_rate": 1,
				"group_by":       "ip",
			},
			expectedErr: "secondary rate is only valid when using entity-based grouping",
		},
		"secondary_rate_defaults_to_rate_with_entity_then_ip": {
			enterpriseOnly: true,
			reqData: map[string]interface{}{
				"rate":           100,
				"group_by":       "entity_then_ip",
				"secondary_rate": 0,
			},
			expectedReadContains: map[string]interface{}{
				"group_by":       "entity_then_ip",
				"secondary_rate": json.Number("100"),
			},
		},
		"secondary_rate_defaults_to_rate_with_entity_then_none": {
			enterpriseOnly: true,
			reqData: map[string]interface{}{
				"rate":     100,
				"group_by": "entity_then_none",
			},
			expectedReadContains: map[string]interface{}{
				"group_by":       "entity_then_none",
				"secondary_rate": json.Number("100"),
			},
		},
		"secondary_rate_defaults_to_zero_on_group_by_none": {
			enterpriseOnly: true,
			reqData: map[string]interface{}{
				"rate":     100,
				"group_by": "none",
			},
			expectedReadContains: map[string]interface{}{
				"group_by":       "none",
				"secondary_rate": json.Number("0"),
			},
		},
		"secondary_rate_defaults_to_zero_on_group_by_ip": {
			reqData: map[string]interface{}{
				"rate":     100,
				"group_by": "ip",
			},
			expectedReadContains: map[string]interface{}{
				"group_by":       "ip",
				"secondary_rate": json.Number("0"),
			},
		},
		"secondary_rate_defaults_cannot_be_negative": {
			enterpriseOnly: true,
			reqData: map[string]interface{}{
				"rate":           100,
				"group_by":       "entity_then_ip",
				"secondary_rate": -1,
			},
			expectedErr: "secondary rate must be greater than or equal to 0",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.enterpriseOnly && !constants.IsEnterprise {
				t.Skip("skipping test because it is only valid in enterprise")
			} else if tc.CeOnly && constants.IsEnterprise {
				t.Skip("skipping test because it is only valid in community edition")
			}

			_, err := client.Logical().Write("sys/quotas/rate-limit/"+name, tc.reqData)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)

				resp, err := client.Logical().Read("sys/quotas/rate-limit/" + name)
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.NotEmpty(t, resp.Data)
				for k, v := range tc.expectedReadContains {
					require.Contains(t, resp.Data, k)
					require.Equal(t, v, resp.Data[k])
				}
			}

			_, err = client.Logical().Delete("sys/quotas/rate-limit/" + name)
			require.NoError(t, err)
		})
	}
}

// TestQuotas_RateLimit_ZeroRetryRegression verifies that the rate limit response
// headers do not return a Retry-After value of 0.
func TestQuotas_RateLimit_ZeroRetryRegression(t *testing.T) {
	conf, opts := teststorage.ClusterSetup(coreConfig, nil, nil)
	cluster := vault.NewTestCluster(t, conf, opts)
	cluster.Start()
	defer cluster.Cleanup()

	testhelpers.WaitForActiveNode(t, cluster)
	client := cluster.Cores[0].Client
	_, err := client.Logical().Write("sys/quotas/config", map[string]interface{}{
		"enable_rate_limit_response_headers": true,
	})
	require.NoError(t, err)
	_, err = client.Logical().Write("sys/quotas/rate-limit/root-rlq", map[string]interface{}{
		"name": "root-rlq",
		"rate": 1,
	})
	require.NoError(t, err)
	failed := atomic.NewBool(false)
	wg := sync.WaitGroup{}
	client = client.WithResponseCallbacks(func(response *api.Response) {
		if response.Header.Get("Retry-After") == "0" {
			failed.Store(true)
		}
	})
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client.Logical().Read("sys/mounts")
		}()
	}
	wg.Wait()
	require.False(t, failed.Load())
}
