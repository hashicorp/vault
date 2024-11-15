// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// The implementation of this package is a simple sync.Map to allow creation
// and redemption of nonces. Creating nonces consumes memory and sync.Map is
// not easily shrinkable without an expensive copy operation (which has not
// be implemented here). However, redemption of nonces (including invalid
// or reused nonces) is as fast as a single sync.Map lookup.
//
// As such, this implementation should not be used.

package nonceutil

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

type syncMapNonceService struct {
	validity   time.Duration
	issued     *atomic.Uint64
	nextExpiry *atomic.Int64
	nonces     *sync.Map // map[string]time.Time
}

var _ NonceService = &syncMapNonceService{}

func newSyncMapNonceService(validity time.Duration) *syncMapNonceService {
	return &syncMapNonceService{
		validity:   validity,
		issued:     new(atomic.Uint64),
		nextExpiry: new(atomic.Int64),
		nonces:     new(sync.Map),
	}
}

func (a *syncMapNonceService) Initialize() error { return nil }
func (a *syncMapNonceService) IsStrict() bool    { return true }
func (a *syncMapNonceService) IsCrossNode() bool { return false }

func generateNonce() (string, error) {
	return generateRandomBase64(21)
}

func generateRandomBase64(srcBytes int) (string, error) {
	data := make([]byte, 21)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(data), nil
}

func (a *syncMapNonceService) Get() (string, time.Time, error) {
	now := time.Now()
	nonce, err := generateNonce()
	if err != nil {
		return "", now, err
	}

	then := now.Add(a.validity)
	a.nonces.Store(nonce, then)

	nextExpiry := a.nextExpiry.Load()
	next := time.Unix(nextExpiry, 0)
	if then.Before(next) {
		a.nextExpiry.Store(then.Unix())
	}

	a.issued.Add(1)

	return nonce, then, nil
}

func (a *syncMapNonceService) Redeem(nonce string) bool {
	rawTimeout, present := a.nonces.LoadAndDelete(nonce)
	if !present {
		return false
	}

	timeout := rawTimeout.(time.Time)
	if time.Now().After(timeout) {
		return false
	}

	return true
}

func (a *syncMapNonceService) Tidy() *NonceStatus {
	now := time.Now()
	nextRun := now.Add(a.validity)
	var outstanding uint64
	a.nonces.Range(func(key, value any) bool {
		timeout := value.(time.Time)
		if now.After(timeout) {
			a.nonces.Delete(key)
		} else {
			outstanding += 1
		}

		if timeout.Before(nextRun) {
			nextRun = timeout
		}

		return false /* don't quit looping */
	})

	a.nextExpiry.Store(nextRun.Unix())

	return &NonceStatus{
		Issued:      a.issued.Load(),
		Outstanding: outstanding,
	}
}
