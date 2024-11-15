// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dependency

import (
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
)

// Ensure implements
var _ Dependency = (*VaultAgentTokenQuery)(nil)

const (
	// VaultAgentTokenSleepTime is the amount of time to sleep between queries, since
	// the fsnotify library is not compatible with solaris and other OSes yet.
	VaultAgentTokenSleepTime = DefaultNonBlockingQuerySleepTime
)

// VaultAgentTokenQuery is the dependency to Vault Agent token
type VaultAgentTokenQuery struct {
	stopCh chan struct{}
	path   string
	stat   os.FileInfo
}

// NewVaultAgentTokenQuery creates a new dependency.
func NewVaultAgentTokenQuery(path string) (*VaultAgentTokenQuery, error) {
	return &VaultAgentTokenQuery{
		stopCh: make(chan struct{}, 1),
		path:   path,
	}, nil
}

// Fetch retrieves this dependency and returns the result or any errors that
// occur in the process.
func (d *VaultAgentTokenQuery) Fetch(clients *ClientSet, opts *QueryOptions) (interface{}, *ResponseMetadata, error) {
	log.Printf("[TRACE] %s: READ %s", d, d.path)

	var token string
	select {
	case <-d.stopCh:
		log.Printf("[TRACE] %s: stopped", d)
		return "", nil, ErrStopped
	case r := <-d.watch(d.stat):
		if r.err != nil {
			return "", nil, errors.Wrap(r.err, d.String())
		}

		log.Printf("[TRACE] %s: reported change", d)

		raw_token, err := os.ReadFile(d.path)
		if err != nil {
			return "", nil, errors.Wrap(err, d.String())
		}
		d.stat = r.stat
		token = string(raw_token)
	}

	return respWithMetadata(token)
}

// CanShare returns if this dependency is sharable.
func (d *VaultAgentTokenQuery) CanShare() bool {
	return false
}

// Stop halts the dependency's fetch function.
func (d *VaultAgentTokenQuery) Stop() {
	close(d.stopCh)
}

// String returns the human-friendly version of this dependency.
func (d *VaultAgentTokenQuery) String() string {
	return "vault-agent.token"
}

// Type returns the type of this dependency.
func (d *VaultAgentTokenQuery) Type() Type {
	return TypeVault
}

// watch watches the file for changes
func (d *VaultAgentTokenQuery) watch(lastStat os.FileInfo) <-chan *watchResult {
	ch := make(chan *watchResult, 1)

	go func(lastStat os.FileInfo) {
		for {
			stat, err := os.Stat(d.path)
			if err != nil {
				select {
				case <-d.stopCh:
					return
				case ch <- &watchResult{err: err}:
					return
				}
			}

			changed := lastStat == nil ||
				lastStat.Size() != stat.Size() ||
				lastStat.ModTime() != stat.ModTime()

			if changed {
				select {
				case <-d.stopCh:
					return
				case ch <- &watchResult{stat: stat}:
					return
				}
			}

			time.Sleep(VaultAgentTokenSleepTime)
		}
	}(lastStat)

	return ch
}
