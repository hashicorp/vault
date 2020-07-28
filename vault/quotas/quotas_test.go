package quotas

import (
	"context"
	"testing"

	"github.com/go-test/deep"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/sdk/helper/logging"
)

func TestQuotas_Precedence(t *testing.T) {
	qm, err := NewManager(logging.NewVaultLogger(log.Trace), nil, metricsutil.BlackholeSink())
	if err != nil {
		t.Fatal(err)
	}

	setQuotaFunc := func(t *testing.T, name, nsPath, mountPath string) Quota {
		t.Helper()
		quota := NewRateLimitQuota(name, nsPath, mountPath, 10)
		err := qm.SetQuota(context.Background(), TypeRateLimit.String(), quota, true)
		if err != nil {
			t.Fatal(err)
		}
		return quota
	}

	checkQuotaFunc := func(t *testing.T, nsPath, mountPath string, expected Quota) {
		t.Helper()
		quota, err := qm.QueryQuota(&Request{
			Type:          TypeRateLimit,
			NamespacePath: nsPath,
			MountPath:     mountPath,
		})
		if err != nil {
			t.Fatal(err)
		}
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
