// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

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

type PartialSealWrapError struct {
	Err error
}

func (p *PartialSealWrapError) Error() string {
	return p.Err.Error()
}

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

	// Configured indicates the wrapper was successfully configured at initialization
	Configured bool

	// hcLock protects lastHealthy, lastSeenHealthy, and healthy.
	// Do not modify those fields directly, use setHealth instead.
	// Do not access these fields directly, use getHealth instead.
	hcLock          sync.RWMutex
	lastHealthCheck time.Time
	lastSeenHealthy time.Time
	healthy         bool
}

func NewSealWrapper(wrapper wrapping.Wrapper, priority int, name string, sealConfigType string, disabled bool, configured bool) *SealWrapper {
	ret := &SealWrapper{
		Wrapper:         wrapper,
		Priority:        priority,
		Name:            name,
		SealConfigType:  sealConfigType,
		Disabled:        disabled,
		Configured:      configured,
		lastSeenHealthy: time.Now(),
		healthy:         false,
	}

	if configured {
		ret.healthy = true
	}

	return ret
}

func (sw *SealWrapper) SetHealthy(healthy bool, checkTime time.Time) {
	sw.hcLock.Lock()
	defer sw.hcLock.Unlock()

	sw.healthy = healthy
	sw.lastHealthCheck = checkTime

	if healthy {
		sw.lastSeenHealthy = checkTime
	}
}

func (sw *SealWrapper) IsHealthy() bool {
	healthy, _, _ := getHealth(sw)

	return healthy
}

func (sw *SealWrapper) LastSeenHealthy() time.Time {
	_, lastSeenHealthy, _ := getHealth(sw)

	return lastSeenHealthy
}

func (sw *SealWrapper) LastHealthCheck() time.Time {
	_, _, lastHealthCheck := getHealth(sw)

	return lastHealthCheck
}

var (
	// vars for unit testing
	HealthTestIntervalNominal   = 10 * time.Minute
	HealthTestIntervalUnhealthy = 1 * time.Minute
	HealthTestTimeout           = 1 * time.Minute
)

func (sw *SealWrapper) CheckHealth(ctx context.Context, checkTime time.Time) error {
	testVal := fmt.Sprintf("Heartbeat %d", mathrand.Intn(1000))
	ciphertext, err := sw.Wrapper.Encrypt(ctx, []byte(testVal), nil)
	if err != nil {
		sw.SetHealthy(false, checkTime)
		return fmt.Errorf("failed to encrypt test value, seal wrapper may be unreachable: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, HealthTestTimeout)
	defer cancel()
	plaintext, err := sw.Wrapper.Decrypt(ctx, ciphertext, nil)
	if err != nil && !IsOldKeyError(err) {
		sw.SetHealthy(false, checkTime)
		return fmt.Errorf("failed to decrypt test value, seal wrapper may be unreachable: %w", err)
	}
	if !bytes.Equal([]byte(testVal), plaintext) {
		sw.SetHealthy(false, checkTime)
		return errors.New("failed to decrypt health test value to expected result")
	}

	sw.SetHealthy(true, checkTime)

	return nil
}

// getHealth is the only function allowed to inspect the health fields directly
func getHealth(sw *SealWrapper) (healthy bool, lastSeenHealthy time.Time, lastHealthCheck time.Time) {
	sw.hcLock.RLock()
	defer sw.hcLock.RUnlock()

	return sw.healthy, sw.lastSeenHealthy, sw.lastHealthCheck
}
