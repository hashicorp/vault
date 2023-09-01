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
type SealWrapper struct {
	Wrapper  wrapping.Wrapper
	Priority int
	Name     string

	// sealConfigType is the KMS.Type of this wrapper. It is a string rather than a SealConfigType
	// to avoid a circular go package depency
	SealConfigType string

	// Disabled indicates, when true indicates that this wrapper should only be used for decryption.
	Disabled bool

	HcLock          sync.RWMutex
	LastHealthCheck time.Time
	LastSeenHealthy time.Time
	Healthy         bool
}

func (sw *SealWrapper) IsHealthy() bool {
	sw.HcLock.RLock()
	defer sw.HcLock.RUnlock()

	return sw.Healthy
}

var (
	// vars for unit testing
	HealthTestIntervalNominal   = 10 * time.Minute
	HealthTestIntervalUnhealthy = 1 * time.Minute
	HealthTestTimeout           = 1 * time.Minute
)

func (sw *SealWrapper) CheckHealth(ctx context.Context, checkTime time.Time) error {
	sw.HcLock.Lock()
	defer sw.HcLock.Unlock()

	sw.LastHealthCheck = checkTime

	// Assume the wrapper is unhealthy, if we make it to the end we'll set it to true
	sw.Healthy = false

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

	sw.LastSeenHealthy = checkTime
	sw.Healthy = true

	return nil
}
