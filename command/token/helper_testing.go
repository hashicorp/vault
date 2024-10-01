// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package token

import (
	"sync"

	"github.com/hashicorp/vault/api/tokenhelper"
)

var _ tokenhelper.TokenHelper = (*TestingTokenHelper)(nil)

// TestingTokenHelper implements token.TokenHelper which runs entirely
// in-memory. This should not be used outside of testing.
type TestingTokenHelper struct {
	lock  sync.RWMutex
	token string
}

func NewTestingTokenHelper() *TestingTokenHelper {
	return &TestingTokenHelper{}
}

func (t *TestingTokenHelper) Erase() error {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.token = ""
	return nil
}

func (t *TestingTokenHelper) Get() (string, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.token, nil
}

func (t *TestingTokenHelper) Path() string {
	return ""
}

func (t *TestingTokenHelper) Store(token string) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.token = token
	return nil
}
