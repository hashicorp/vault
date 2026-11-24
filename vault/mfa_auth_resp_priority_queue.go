// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/hashicorp/vault/sdk/queue"
)

// NewLoginMFAPriorityQueue initializes the internal data structures and returns a new
// PriorityQueue with the default pending limits.
func NewLoginMFAPriorityQueue() *LoginMFAPriorityQueue {
	pq := queue.New()
	loginPQ := &LoginMFAPriorityQueue{
		wrapped:  pq,
		byEntity: make(map[string]int),
		bySource: make(map[string]int),
		limits: mfaPendingLimits{
			perEntity:   maxPendingMFAPerEntity,
			perSourceIP: maxPendingMFAPerSourceIP,
			global:      maxPendingMFAGlobal,
		},
	}
	return loginPQ
}

// mfaPendingLimits holds the maximum number of live pending MFA entries
// allowed per entity, per source IP and across the whole queue.
type mfaPendingLimits struct {
	perEntity   int
	perSourceIP int
	global      int
}

// LoginMFAPriorityQueue is a priority queue of pending MFA auth responses,
// ordered by storage time. It also enforces the per-entity, per-source-IP and
// global limits on live pending entries.
type LoginMFAPriorityQueue struct {
	wrapped *queue.PriorityQueue

	// limits is set when the queue is created and not changed afterwards, so
	// it needs no locking. Tests may lower it on a fresh queue to avoid
	// pushing the full production caps.
	limits mfaPendingLimits

	// byEntity and bySource count queued entries per entity ID and per source
	// IP. They are guarded by l and must be updated whenever an entry enters or
	// leaves the queue. Keys are deleted when their count reaches zero.
	byEntity map[string]int
	bySource map[string]int

	// Here is a scenarios in which the lock is needed. For example, suppose
	// RemoveExpiredMfaAuthResponse function pops an item to check if the item
	// has been expired or not and assume that the item is still valid. Then,
	// if in the meantime, an MFA validation request comes in for the same
	// item, the /sys/mfa/validate endpoint will return invalid request ID
	// which is not true.
	l sync.RWMutex
}

// Len returns the count of items in the Priority Queue
func (pq *LoginMFAPriorityQueue) Len() int {
	pq.l.Lock()
	defer pq.l.Unlock()
	return pq.wrapped.Len()
}

// Push pushes an item on to the queue. This is a wrapper/convenience
// method that calls heap.Push, so consumers do not need to invoke heap
// functions directly. Items must have unique Keys, and Items in the queue
// cannot be updated. To modify an Item, users must first remove it and re-push
// it after modifications.
//
// Push enforces the per-entity, per-source-IP and global pending limits. If
// adding the entry would exceed a limit, it first removes expired entries and
// re-checks, and returns an *mfaPendingLimitError if a limit is still reached.
func (pq *LoginMFAPriorityQueue) Push(resp *MFACachedAuthResponse) error {
	pq.l.Lock()
	defer pq.l.Unlock()

	item := &queue.Item{
		Key:      resp.RequestID,
		Value:    resp,
		Priority: resp.TimeOfStorage.Unix(),
	}

	entityID, sourceIP := mfaEntryScopes(resp)
	if err := pq.checkLimitsLocked(entityID, sourceIP); err != nil {
		// Cleanup may lag behind expiry, so drop expired entries and re-check
		// before rejecting.
		if rmErr := pq.removeExpiredLocked(defaultMFAAuthResponseTTL, time.Now(), maxExpiredMFAAuthResponsesPerCleanup); rmErr != nil {
			return rmErr
		}
		if err := pq.checkLimitsLocked(entityID, sourceIP); err != nil {
			return err
		}
	}

	if err := pq.wrapped.Push(item); err != nil {
		return err
	}
	pq.adjustCountsLocked(entityID, sourceIP, 1)
	return nil
}

// mfaPendingLimitError is returned when a pending MFA entry is rejected
// because a limit was reached. The message is deliberately generic.
type mfaPendingLimitError struct {
	scope string
}

func (e *mfaPendingLimitError) Error() string {
	return "too many pending MFA requests, please try again later"
}

// mfaEntryScopes returns the entity ID and source IP (without port) that an
// entry counts against. Either may be empty, in which case that scope is not
// enforced for the entry.
func mfaEntryScopes(resp *MFACachedAuthResponse) (entityID, sourceIP string) {
	if resp.CachedAuth != nil {
		entityID = resp.CachedAuth.EntityID
	}
	sourceIP = resp.RequestConnRemoteAddr
	if host, _, err := net.SplitHostPort(sourceIP); err == nil {
		sourceIP = host
	}
	return entityID, sourceIP
}

// checkLimitsLocked reports whether adding one entry would exceed a limit.
// The caller must hold pq.l.
func (pq *LoginMFAPriorityQueue) checkLimitsLocked(entityID, sourceIP string) error {
	if pq.wrapped.Len() >= pq.limits.global {
		return &mfaPendingLimitError{scope: "global"}
	}
	if entityID != "" && pq.byEntity[entityID] >= pq.limits.perEntity {
		return &mfaPendingLimitError{scope: "entity"}
	}
	if sourceIP != "" && pq.bySource[sourceIP] >= pq.limits.perSourceIP {
		return &mfaPendingLimitError{scope: "source_ip"}
	}
	return nil
}

// adjustCountsLocked applies delta to the per-entity and per-source counters,
// deleting keys that reach zero. The caller must hold pq.l.
func (pq *LoginMFAPriorityQueue) adjustCountsLocked(entityID, sourceIP string, delta int) {
	adjust := func(m map[string]int, key string) {
		if key == "" {
			return
		}
		if n := m[key] + delta; n > 0 {
			m[key] = n
		} else {
			delete(m, key)
		}
	}
	adjust(pq.byEntity, entityID)
	adjust(pq.bySource, sourceIP)
}

// PopByKey searches the queue for an item with the given key and removes it
// from the queue if found. Returns nil if not found. The per-entity and
// per-source counters are decremented for the removed entry.
func (pq *LoginMFAPriorityQueue) PopByKey(reqID string) (*MFACachedAuthResponse, error) {
	pq.l.Lock()
	defer pq.l.Unlock()

	item, err := pq.wrapped.PopByKey(reqID)
	if err != nil || item == nil {
		return nil, err
	}

	resp := item.Value.(*MFACachedAuthResponse)
	entityID, sourceIP := mfaEntryScopes(resp)
	pq.adjustCountsLocked(entityID, sourceIP, -1)
	return resp, nil
}

// maxExpiredMFAAuthResponsesPerCleanup bounds how many expired entries one
// sweep removes while holding the queue lock, whether from a cleanup tick or
// an inline sweep in Push.
const maxExpiredMFAAuthResponsesPerCleanup = 10000

// RemoveExpiredMfaAuthResponse removes expired entries from the queue, oldest
// first, while holding the queue lock. It pops entries until it reaches one
// that has not expired (which it pushes back) or the queue is empty. Each call
// removes at most maxExpiredMFAAuthResponsesPerCleanup entries, so that a large
// backlog cannot hold the lock for too long. Any remaining expired entries are
// removed by later calls. Counters for the removed entries are updated. Push
// runs the same sweep, over at most the same number of entries, before it
// rejects an entry for being at a limit. It returns an error if popping or
// re-pushing an entry fails. cutoffTime should normally be time.Now() except
// for tests.
func (pq *LoginMFAPriorityQueue) RemoveExpiredMfaAuthResponse(expiryTime time.Duration, cutoffTime time.Time) error {
	pq.l.Lock()
	defer pq.l.Unlock()
	return pq.removeExpiredLocked(expiryTime, cutoffTime, maxExpiredMFAAuthResponsesPerCleanup)
}

// removeExpiredLocked removes up to limit expired entries, oldest first, and
// updates the counters for each removed entry. The caller must hold pq.l.
func (pq *LoginMFAPriorityQueue) removeExpiredLocked(expiryTime time.Duration, cutoffTime time.Time, limit int) error {
	for removed := 0; removed < limit; removed++ {
		item, err := pq.wrapped.Pop()
		if errors.Is(err, queue.ErrEmpty) {
			return nil
		}
		if err != nil {
			return err
		}

		resp := item.Value.(*MFACachedAuthResponse)
		entityID, sourceIP := mfaEntryScopes(resp)
		if !resp.isExpired(expiryTime, cutoffTime) {
			// The oldest remaining entry is live; put it back without
			// touching the counters since it never left the queue logically.
			if err := pq.wrapped.Push(item); err != nil {
				pq.adjustCountsLocked(entityID, sourceIP, -1)
				return fmt.Errorf("failed to re-push live MFA entry: %w", err)
			}
			return nil
		}
		pq.adjustCountsLocked(entityID, sourceIP, -1)
	}

	return nil
}
