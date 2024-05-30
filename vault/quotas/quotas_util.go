// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package quotas

import (
	"context"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-memdb"
	"github.com/hashicorp/vault/helper/metricsutil"
)

func quotaTypes() []string {
	return []string{
		TypeRateLimit.String(),
	}
}

func (m *Manager) init(walkFunc leaseWalkFunc) {}

func (m *Manager) recomputeLeaseCounts(ctx context.Context, txn *memdb.Txn) error {
	return nil
}

func (m *Manager) setIsPerfStandby(quota Quota) {}

func (m *Manager) inLeasePathCache(path string) bool {
	return false
}

func (m *Manager) setupDefaultLeaseCountQuotaInStorage(_ctx context.Context) error {
	return nil
}

type entManager struct {
	isPerfStandby bool
	isDRSecondary bool
	isNewInstall  bool
}

func (*entManager) Reset() error {
	return nil
}

type LeaseCountQuota struct{}

func (l LeaseCountQuota) GetNamespacePath() string {
	panic("implement me")
}

func (l LeaseCountQuota) IsInheritable() bool {
	panic("implement me")
}

func (l LeaseCountQuota) allow(_ context.Context, _ *Request) (Response, error) {
	panic("implement me")
}

func (l LeaseCountQuota) quotaID() string {
	panic("implement me")
}

func (l LeaseCountQuota) QuotaName() string {
	panic("implement me")
}

func (l LeaseCountQuota) initialize(logger log.Logger, sink *metricsutil.ClusterMetricSink) error {
	panic("implement me")
}

func (l LeaseCountQuota) close(_ context.Context) error {
	panic("implement me")
}

func (l LeaseCountQuota) Clone() Quota {
	panic("implement me")
}

func (l LeaseCountQuota) handleRemount(mountPath, nsPath string) {
	panic("implement me")
}
