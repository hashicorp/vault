// +build !enterprise

package vault

import (
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

func (m *ExpirationManager) leaseView(*namespace.Namespace) *BarrierView {
	return m.idView
}

func (m *ExpirationManager) tokenIndexView(*namespace.Namespace) *BarrierView {
	return m.tokenView
}

func (m *ExpirationManager) collectLeases() (map[*namespace.Namespace][]string, int, error) {
	leaseCount := 0
	existing := make(map[*namespace.Namespace][]string)

	var keys []string
	var err error
	for attempt := uint(0); attempt < maxCollectAttempts; attempt++ {
		keys, err = logical.CollectKeys(m.quitContext, m.leaseView(namespace.RootNamespace))
		if err == nil || m.quitContext.Err() != nil {
			break
		}

		m.logger.Error("failed to scan for leases", "error", err)
		time.Sleep((1 << attempt) * collectRetryBase)
	}
	if err != nil {
		return nil, 0, errwrap.Wrapf("failed to scan for leases: {{err}}", err)
	}
	existing[namespace.RootNamespace] = keys
	leaseCount += len(keys)
	return existing, leaseCount, nil
}
