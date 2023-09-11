// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"sync"
	"time"

	"github.com/hashicorp/vault/sdk/queue"
)

// NewLoginMFAPriorityQueue initializes the internal data structures and returns a new
// PriorityQueue
func NewLoginMFAPriorityQueue() *LoginMFAPriorityQueue {
	pq := queue.New()
	loginPQ := &LoginMFAPriorityQueue{
		wrapped: pq,
	}
	return loginPQ
}

type LoginMFAPriorityQueue struct {
	wrapped *queue.PriorityQueue

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
// it after modifications
func (pq *LoginMFAPriorityQueue) Push(resp *MFACachedAuthResponse) error {
	pq.l.Lock()
	defer pq.l.Unlock()

	item := &queue.Item{
		Key:      resp.RequestID,
		Value:    resp,
		Priority: resp.TimeOfStorage.Unix(),
	}

	return pq.wrapped.Push(item)
}

// PopByKey searches the queue for an item with the given key and removes it
// from the queue if found. Returns nil if not found.
func (pq *LoginMFAPriorityQueue) PopByKey(reqID string) (*MFACachedAuthResponse, error) {
	pq.l.Lock()
	defer pq.l.Unlock()

	item, err := pq.wrapped.PopByKey(reqID)
	if err != nil || item == nil {
		return nil, err
	}

	return item.Value.(*MFACachedAuthResponse), nil
}

// RemoveExpiredMfaAuthResponse pops elements of the queue and check
// if the entry has expired or not. If the entry has not expired, it pushes
// back the entry to the queue. It returns false if there is no expired element
// left to be removed, true otherwise.
// cutoffTime should normally be time.Now() except for tests.
func (pq *LoginMFAPriorityQueue) RemoveExpiredMfaAuthResponse(expiryTime time.Duration, cutoffTime time.Time) error {
	pq.l.Lock()
	defer pq.l.Unlock()

	item, err := pq.wrapped.Pop()
	if err != nil && err != queue.ErrEmpty {
		return err
	}
	if err == queue.ErrEmpty {
		return nil
	}

	mfaResp := item.Value.(*MFACachedAuthResponse)

	storageTime := mfaResp.TimeOfStorage
	if cutoffTime.Before(storageTime.Add(expiryTime)) {
		// the highest priority entry has not been expired yet, pushing it
		// back and return
		err := pq.wrapped.Push(item)
		if err != nil {
			return err
		}
	}
	return nil
}
