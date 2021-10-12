package quotas

import (
	"context"
	"testing"
	"time"

	"github.com/go-test/deep"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/stretchr/testify/require"
)

func TestQuotas_MountPathOverwrite(t *testing.T) {
	qm, err := NewManager(logging.NewVaultLogger(log.Trace), nil, metricsutil.BlackholeSink())
	require.NoError(t, err)

	quota := NewRateLimitQuota("tq", "", "kv1/", 10, time.Second, 0)
	require.NoError(t, qm.SetQuota(context.Background(), TypeRateLimit.String(), quota, false))
	quota = quota.Clone()
	quota.MountPath = "kv2/"
	require.NoError(t, qm.SetQuota(context.Background(), TypeRateLimit.String(), quota, false))

	q, err := qm.QueryQuota(&Request{
		Type:      TypeRateLimit,
		MountPath: "kv1/",
	})
	require.NoError(t, err)
	require.Nil(t, q)

	require.NoError(t, qm.DeleteQuota(context.Background(), TypeRateLimit.String(), "tq"))

	q, err = qm.QueryQuota(&Request{
		Type:      TypeRateLimit,
		MountPath: "kv1/",
	})
	require.NoError(t, err)
	require.Nil(t, q)
}

func TestQuotas_Precedence(t *testing.T) {
	qm, err := NewManager(logging.NewVaultLogger(log.Trace), nil, metricsutil.BlackholeSink())
	require.NoError(t, err)

	setQuotaFunc := func(t *testing.T, name, nsPath, mountPath string) Quota {
		t.Helper()
		quota := NewRateLimitQuota(name, nsPath, mountPath, 10, time.Second, 0)
		require.NoError(t, qm.SetQuota(context.Background(), TypeRateLimit.String(), quota, true))
		return quota
	}

	checkQuotaFunc := func(t *testing.T, nsPath, mountPath string, expected Quota) {
		t.Helper()
		quota, err := qm.QueryQuota(&Request{
			Type:          TypeRateLimit,
			NamespacePath: nsPath,
			MountPath:     mountPath,
		})
		require.NoError(t, err)

		if diff := deep.Equal(expected, quota); len(diff) > 0 {
			t.Fatal(diff)
		}
	}

	// No quota present. Expect nil.
	checkQuotaFunc(t, "", "", nil)

	// Define global quota and expect that to be returned.
	rateLimitGlobalQuota := setQuotaFunc(t, "rateLimitGlobalQuota", "", "")
	checkQuotaFunc(t, "", "", rateLimitGlobalQuota)

	// Define a global mount specific quota and expect that to be returned.
	rateLimitGlobalMountQuota := setQuotaFunc(t, "rateLimitGlobalMountQuota", "", "testmount")
	checkQuotaFunc(t, "", "testmount", rateLimitGlobalMountQuota)

	// Define a namespace quota and expect that to be returned.
	rateLimitNSQuota := setQuotaFunc(t, "rateLimitNSQuota", "testns", "")
	checkQuotaFunc(t, "testns", "", rateLimitNSQuota)

	// Define a namespace mount specific quota and expect that to be returned.
	rateLimitNSMountQuota := setQuotaFunc(t, "rateLimitNSMountQuota", "testns", "testmount")
	checkQuotaFunc(t, "testns", "testmount", rateLimitNSMountQuota)

	// Now that many quota types are defined, verify that the most specific
	// matches are returned per namespace.
	checkQuotaFunc(t, "", "", rateLimitGlobalQuota)
	checkQuotaFunc(t, "testns", "", rateLimitNSQuota)
}
