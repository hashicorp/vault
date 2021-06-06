// +build !enterprise

package quotas

import (
	"context"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/metricsutil"

	memdb "github.com/hashicorp/go-memdb"
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

type entManager struct {
	isPerfStandby bool
	isDRSecondary bool
}

func (*entManager) Reset() error {
	return nil
}

type LeaseCountQuota struct{}

func (l LeaseCountQuota) allow(request *Request) (Response, error) {
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

func (l LeaseCountQuota) close() error {
	panic("implement me")
}

func (l LeaseCountQuota) handleRemount(s string) {
	panic("implement me")
}
