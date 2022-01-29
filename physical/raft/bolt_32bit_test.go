//go:build 386 || arm

package raft

import (
	"os"
	"strconv"
	"testing"
)

func Test_BoltOptions(t *testing.T) {
	t.Parallel()
	key := "VAULT_RAFT_INITIAL_MMAP_SIZE"

	testCases := []struct {
		name         string
		env          string
		expectedSize int
	}{
		{"none", "", 0},
		{"5MB", strconv.Itoa(5 * 1024 * 1024), 5 * 1024 * 1024},
		{"negative", "-1", 0},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			if tc.env != "" {
				current := os.Getenv(key)
				defer os.Setenv(key, current)
				os.Setenv(key, tc.env)
			}

			o := boltOptions("")

			if o.InitialMmapSize != tc.expectedSize {
				t.Errorf("expected InitialMmapSize to be %d but it was %d", tc.expectedSize, o.InitialMmapSize)
			}
		})
	}
}
