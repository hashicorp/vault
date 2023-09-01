package seal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	mathrand "math/rand"
	"sync"
	"time"

	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
)

// SealWrapper contains a Wrapper and related information needed by the seal that uses it.
// Use NewSealWrapper to construct new instances, do not do it directly.
type SealWrapper struct {
	Wrapper  wrapping.Wrapper
	Priority int
	Name     string

	// sealConfigType is the KMS.Type of this wrapper. It is a string rather than a SealConfigType
	// to avoid a circular go package depency
	SealConfigType string

	// Disabled indicates, when true indicates that this wrapper should only be used for decryption.
	Disabled bool

	// hcLock protects lastHealthy, lastSeenHealthy, and healthy. Do not modify those fields directly, use setHealth instead.
	hcLock          sync.RWMutex
	lastHealthCheck time.Time
	lastSeenHealthy time.Time
	healthy         bool
}

func NewSealWrapper(wrapper wrapping.Wrapper, priority int, name string, sealConfigType string, disabled bool) *SealWrapper {
	ret := &SealWrapper{
		Wrapper:        wrapper,
		Priority:       priority,
		Name:           name,
		SealConfigType: sealConfigType,
		Disabled:       disabled,
	}

	ret.setHealth(true, time.Now(), ret.lastHealthCheck)

	return ret
}

func (sw *SealWrapper) rlock() func() {
	sw.hcLock.RLock()
	return sw.hcLock.RUnlock
}

func (sw *SealWrapper) lock() func() {
	sw.hcLock.Lock()
	return sw.hcLock.Unlock
}

func (sw *SealWrapper) SetHealthy(healthy bool, checkTime time.Time) {
	unlock := sw.lock()
	defer unlock()

	wasHealthy := sw.healthy
	lastHealthy := sw.lastSeenHealthy
	if !wasHealthy && healthy {
		lastHealthy = checkTime
	}

	sw.setHealth(healthy, lastHealthy, checkTime)
}

func (sw *SealWrapper) IsHealthy() bool {
	unlock := sw.rlock()
	defer unlock()

	return sw.healthy
}

func (sw *SealWrapper) LastSeenHealthy() time.Time {
	unlock := sw.rlock()
	defer unlock()

	return sw.lastSeenHealthy
}

func (sw *SealWrapper) LastHealthCheck() time.Time {
	unlock := sw.rlock()
	defer unlock()

	return sw.lastHealthCheck
}

var (
	// vars for unit testing
	HealthTestIntervalNominal   = 10 * time.Minute
	HealthTestIntervalUnhealthy = 1 * time.Minute
	HealthTestTimeout           = 1 * time.Minute
)

func (sw *SealWrapper) CheckHealth(ctx context.Context, checkTime time.Time) error {
	unlock := sw.lock()
	defer unlock()

	// Assume the wrapper is unhealthy, if we make it to the end we'll set it to true
	sw.setHealth(false, sw.lastSeenHealthy, checkTime)

	testVal := fmt.Sprintf("Heartbeat %d", mathrand.Intn(1000))
	ciphertext, err := sw.Wrapper.Encrypt(ctx, []byte(testVal), nil)
	if err != nil {
		return fmt.Errorf("failed to encrypt test value, seal wrapper may be unreachable: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, HealthTestTimeout)
	defer cancel()
	plaintext, err := sw.Wrapper.Decrypt(ctx, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("failed to decrypt test value, seal wrapper may be unreachable: %w", err)
	}
	if !bytes.Equal([]byte(testVal), plaintext) {
		return errors.New("failed to decrypt health test value to expected result")
	}

	sw.setHealth(true, checkTime, checkTime)

	return nil
}

// setHealth sets the fields protected by sw.hcLock, callers *must* hold the write lock.
func (sw *SealWrapper) setHealth(healthy bool, lastSeenHealthy, lastHealthCheck time.Time) {
	sw.healthy = healthy
	sw.lastSeenHealthy = lastSeenHealthy
	sw.lastHealthCheck = lastHealthCheck
}
