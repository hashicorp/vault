// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"slices"
	"strings"

	"github.com/hashicorp/vault/helper/locking"
)

//
// This file contains locking helper functions that are related to core.go.
//

// parseDeadlockDetectionSetting takes the detectDeadlockConfigParameter string
// and transforms it to a lowercase version of the string, then splits it into
// a slice of strings by interpreting commas as the element delimiters.
func parseDetectDeadlockConfigParameter(detectDeadlockConfigParameter string) []string {
	if detectDeadlockConfigParameter == "" {
		// This doesn't seem necessary, since the companion functions that use
		// this slice can handle an empty slice just the same as a nil slice,
		// but for the sake of compatibility, this will be introduced for now
		// until all occurrences that rely on Core.detectDeadlocks have been
		// switched to using functions from this file to create their locks.
		return nil
	}

	result := strings.Split(strings.ToLower(detectDeadlockConfigParameter), ",")
	for i := range result {
		result[i] = strings.TrimSpace(result[i])
	}

	return result
}

// createAppropriateRWMutex determines if the specified lock (identifier) should
// use a deadlock detecting implementation (locking.DeadlockRWMutex) or simply a
// sync.RWMutex instance. This is done by splitting the deadlockDetectionLocks
// string into a slice of strings. If the slice contains the specified lock
// (identifier), then the deadlock detecting implementation is used, otherwise a
// sync.Mutex is returned.
func createAppropriateRWMutex(deadlockDetectionLocks []string, lock string) locking.RWMutex {
	if slices.Contains(deadlockDetectionLocks, strings.ToLower(lock)) {
		return &locking.DeadlockRWMutex{}
	}

	return &locking.SyncRWMutex{}
}
