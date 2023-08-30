// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

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

	quota := NewRateLimitQuota("tq", "", "kv1/", "", "", false, time.Second, 0, 10)
	require.NoError(t, qm.SetQuota(context.Background(), TypeRateLimit.String(), quota, false))
	quota = quota.Clone().(*RateLimitQuota)
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

	setQuotaFunc := func(t *testing.T, name, nsPath, mountPath, pathSuffix, role string, inheritable bool) Quota {
		t.Helper()
		quota := NewRateLimitQuota(name, nsPath, mountPath, pathSuffix, role, inheritable, time.Second, 0, 10)
		require.NoError(t, qm.SetQuota(context.Background(), TypeRateLimit.String(), quota, true))
		return quota
	}

	checkQuotaFunc := func(t *testing.T, nsPath, mountPath, pathSuffix, role string, expected Quota) {
		t.Helper()
		quota, err := qm.QueryQuota(&Request{
			Type:          TypeRateLimit,
			NamespacePath: nsPath,
			MountPath:     mountPath,
			Role:          role,
			Path:          nsPath + mountPath + pathSuffix,
		})
		require.NoError(t, err)

		if diff := deep.Equal(expected, quota); len(diff) > 0 {
			t.Fatal(diff)
		}
	}

	// No quota present. Expect nil.
	checkQuotaFunc(t, "", "", "", "", nil)

	// Define global quota and expect that to be returned.
	rateLimitGlobalQuota := setQuotaFunc(t, "rateLimitGlobalQuota", "", "", "", "", true)
	checkQuotaFunc(t, "", "", "", "", rateLimitGlobalQuota)

	// Define a global mount specific quota and expect that to be returned.
	rateLimitGlobalMountQuota := setQuotaFunc(t, "rateLimitGlobalMountQuota", "", "testmount/", "", "", false)
	checkQuotaFunc(t, "", "testmount/", "", "", rateLimitGlobalMountQuota)

	// Define a global mount + path specific quota and expect that to be returned.
	rateLimitGlobalMountPathQuota := setQuotaFunc(t, "rateLimitGlobalMountPathQuota", "", "testmount/", "testpath", "", false)
	checkQuotaFunc(t, "", "testmount/", "testpath", "", rateLimitGlobalMountPathQuota)

	// Define a namespace quota and expect that to be returned.
	rateLimitNSQuota := setQuotaFunc(t, "rateLimitNSQuota", "testns/", "", "", "", false)
	checkQuotaFunc(t, "testns/", "", "", "", rateLimitNSQuota)

	// Define a namespace mount specific quota and expect that to be returned.
	rateLimitNSMountQuota := setQuotaFunc(t, "rateLimitNSMountQuota", "testns/", "testmount/", "", "", false)
	checkQuotaFunc(t, "testns/", "testmount/", "testpath", "", rateLimitNSMountQuota)

	// Define a namespace mount + glob and expect that to be returned.
	rateLimitNSMountGlob := setQuotaFunc(t, "rateLimitNSMountGlob", "testns/", "testmount/", "*", "", false)
	checkQuotaFunc(t, "testns/", "testmount/", "testpath", "", rateLimitNSMountGlob)

	// Define a namespace mount + path specific quota with a glob and expect that to be returned.
	rateLimitNSMountPathSuffixGlob := setQuotaFunc(t, "rateLimitNSMountPathSuffixGlob", "testns/", "testmount/", "test*", "", false)
	checkQuotaFunc(t, "testns/", "testmount/", "testpath", "", rateLimitNSMountPathSuffixGlob)

	// Define a namespace mount + path specific quota with a glob at the end of the path and expect that to be returned.
	rateLimitNSMountPathSuffixGlobAfterPath := setQuotaFunc(t, "rateLimitNSMountPathSuffixGlobAfterPath", "testns/", "testmount/", "testpath*", "", false)
	checkQuotaFunc(t, "testns/", "testmount/", "testpath", "", rateLimitNSMountPathSuffixGlobAfterPath)

	// Define a namespace mount + path specific quota and expect that to be returned.
	rateLimitNSMountPathQuota := setQuotaFunc(t, "rateLimitNSMountPathQuota", "testns/", "testmount/", "testpath", "", false)
	checkQuotaFunc(t, "testns/", "testmount/", "testpath", "", rateLimitNSMountPathQuota)

	// Define a namespace mount + role specific quota and expect that to be returned.
	rateLimitNSMountRoleQuota := setQuotaFunc(t, "rateLimitNSMountPathQuota", "testns/", "testmount/", "", "role", false)
	checkQuotaFunc(t, "testns/", "testmount/", "", "role", rateLimitNSMountRoleQuota)

	// Create an inheritable namespace quota and expect that to be returned on a child namespace
	rateLimitNSInheritableQuota := setQuotaFunc(t, "rateLimitNSInheritableNSQuota", "testns/nested2/", "", "", "", true)
	checkQuotaFunc(t, "testns/nested2/nested3/", "testmount/", "", "", rateLimitNSInheritableQuota)
	checkQuotaFunc(t, "testns/nested2/nested3/nested4/", "testmount/", "", "", rateLimitNSInheritableQuota)

	// Create a non-namespace quota on a nested namespace and make sure it takes precedence over the inherited quota
	rateLimitNonNSNestedQuota := setQuotaFunc(t, "rateLimitNonNSNestedQuota", "testns/nested2/nested3/", "testmount/", "", "", false)
	checkQuotaFunc(t, "testns/nested2/nested3/", "testmount/", "", "", rateLimitNonNSNestedQuota)

	// Create a non-namespace quota on a nested namespace and make sure it takes precedence over the inherited quota
	rateLimitMultiNestedNsInheritableQuota := setQuotaFunc(t, "rateLimitNSInheritableNSQuota", "testns/nested2/nested3/", "", "", "", true)
	checkQuotaFunc(t, "testns/nested2/nested3/nested4/", "testmount/", "", "", rateLimitMultiNestedNsInheritableQuota)

	// Now that many quota types are defined, verify that the most specific
	// matches are returned per namespace.
	checkQuotaFunc(t, "", "", "", "", rateLimitGlobalQuota)
	checkQuotaFunc(t, "testns/", "", "", "", rateLimitNSQuota)
	checkQuotaFunc(t, "testns/nested1/", "", "", "", rateLimitGlobalQuota)
	checkQuotaFunc(t, "testns/nested2/", "", "", "", rateLimitNSInheritableQuota)
	checkQuotaFunc(t, "testns/nested2/nested6/", "", "", "", rateLimitNSInheritableQuota)
	checkQuotaFunc(t, "testns/nested2/nested3/", "", "", "", rateLimitMultiNestedNsInheritableQuota)
	checkQuotaFunc(t, "testns/nested2/nested3/nested4/nested5", "", "", "", rateLimitMultiNestedNsInheritableQuota)
}

// TestQuotas_QueryRoleQuotas checks to see if quota creation on a mount
// requires a call to ResolveRoleOperation.
func TestQuotas_QueryResolveRole_RateLimitQuotas(t *testing.T) {
	leaseWalkFunc := func(context.Context, func(request *Request) bool) error {
		return nil
	}
	qm, err := NewManager(logging.NewVaultLogger(log.Trace), leaseWalkFunc, metricsutil.BlackholeSink())
	require.NoError(t, err)

	rlqReq := &Request{
		Type:          TypeRateLimit,
		Path:          "",
		MountPath:     "mount1/",
		NamespacePath: "",
		ClientAddress: "127.0.0.1",
	}
	// Check that we have no quotas requiring role resolution on mount1/
	required, err := qm.QueryResolveRoleQuotas(rlqReq)
	require.NoError(t, err)
	require.False(t, required)

	// Create a non-role-based RLQ on mount1/ and make sure it doesn't require role resolution
	rlq := NewRateLimitQuota("tq", rlqReq.NamespacePath, rlqReq.MountPath, rlqReq.Path, rlqReq.Role, false, 1*time.Minute, 10*time.Second, 10)
	require.NoError(t, qm.SetQuota(context.Background(), TypeRateLimit.String(), rlq, false))

	required, err = qm.QueryResolveRoleQuotas(rlqReq)
	require.NoError(t, err)
	require.False(t, required)

	// Create a role-based RLQ on mount1/ and make sure it requires role resolution
	rlqReq.Role = "test"
	rlq = NewRateLimitQuota("tq", rlqReq.NamespacePath, rlqReq.MountPath, rlqReq.Path, rlqReq.Role, false, 1*time.Minute, 10*time.Second, 10)
	require.NoError(t, qm.SetQuota(context.Background(), TypeRateLimit.String(), rlq, false))

	required, err = qm.QueryResolveRoleQuotas(rlqReq)
	require.NoError(t, err)
	require.True(t, required)

	// Check that we have no quotas requiring role resolution on mount2/
	rlqReq.MountPath = "mount2/"
	required, err = qm.QueryResolveRoleQuotas(rlqReq)
	require.NoError(t, err)
	require.False(t, required)
}
