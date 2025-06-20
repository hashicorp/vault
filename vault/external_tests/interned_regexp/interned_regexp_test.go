// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package loadedsnapshots

import (
	"fmt"
	"strconv"
	"sync"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/stretchr/testify/require"
)

// TestInternedRegexpConcurrentAccess tests that multiple goroutines can
// concurrently access and create multiple instances of the same type of mount
// without causing any issues. This is important to ensure that the interned
// regular expressions used in the mount paths do not cause have concurrency
// faults.
func TestInternedRegexpConcurrentAccess(t *testing.T) {
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	// Build mount input
	mountInput := &api.MountInput{
		Type: "pki",
	}

	// Mount 10 PKI secrets engine mounts.
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		mountPath := "pki" + strconv.Itoa(i)
		wg.Add(1)
		go func() {
			err := client.Sys().Mount(mountPath, mountInput)
			require.NoError(t, err)

			// Verify the mount was created
			_, err = client.Sys().GetMount(mountPath)
			require.NoError(t, err)

			_, err = client.Logical().List(fmt.Sprintf("/%s/roles", mountPath))
			require.NoError(t, err)
			wg.Done()
		}()
	}
	wg.Wait()
}
