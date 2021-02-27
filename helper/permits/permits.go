package permits

import (
	"github.com/armon/go-metrics"
	"github.com/hashicorp/vault/sdk/physical"
)

type InstrumentedPermitPool struct {
	*physical.PermitPool
	permitCountMetricName []string
}

func NewInstrumentedPermitPool(permits int, keyPrefix ...string) *InstrumentedPermitPool {
	limitMetricName := append(keyPrefix, "permits-limit")
	metrics.SetGauge(limitMetricName, float32(permits))

	pool := &InstrumentedPermitPool{
		PermitPool:            physical.NewPermitPool(permits),
		permitCountMetricName: append(keyPrefix, "permits"),
	}

	// initialize gauge to 0
	pool.updatePermitGauge()
	return pool
}

// Acquire returns when a permit has been acquired
func (pool *InstrumentedPermitPool) Acquire() {
	pool.PermitPool.Acquire()
	pool.updatePermitGauge()
}

// Release returns a permit to the pool
func (pool *InstrumentedPermitPool) Release() {
	pool.PermitPool.Release()
	pool.updatePermitGauge()
}

// Get number of requests in the permit pool
func (pool *InstrumentedPermitPool) CurrentPermits() int {
	return pool.PermitPool.CurrentPermits()
}

func (pool *InstrumentedPermitPool) updatePermitGauge() {
	metrics.SetGauge(pool.permitCountMetricName, float32(pool.PermitPool.CurrentPermits()))
}
