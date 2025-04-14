// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"fmt"
	"testing"

	"github.com/hashicorp/vault/command/server"
	"github.com/stretchr/testify/require"
)

// TestFormatDiscoveredAddr validates that the string returned by formatDiscoveredAddr always respect the format `host:port`.
func TestFormatDiscoveredAddr(t *testing.T) {
	type TestCase struct {
		addr string
		port uint
		res  string
	}
	cases := []TestCase{
		{addr: "127.0.0.1", port: uint(8200), res: "127.0.0.1:8200"},
		{addr: "192.168.137.1:8201", port: uint(8200), res: "192.168.137.1:8201"},
		{addr: "fe80::aa5e:45ff:fe54:c6ce", port: uint(8200), res: "[fe80::aa5e:45ff:fe54:c6ce]:8200"},
		{addr: "::1", port: uint(8200), res: "[::1]:8200"},
		{addr: "[::1]", port: uint(8200), res: "[::1]:8200"},
		{addr: "[::1]:8201", port: uint(8200), res: "[::1]:8201"},
		{addr: "[fe80::aa5e:45ff:fe54:c6ce]", port: uint(8200), res: "[fe80::aa5e:45ff:fe54:c6ce]:8200"},
		{addr: "[fe80::aa5e:45ff:fe54:c6ce]:8201", port: uint(8200), res: "[fe80::aa5e:45ff:fe54:c6ce]:8201"},
	}
	for i, c := range cases {
		res := formatDiscoveredAddr(c.addr, c.port)
		if res != c.res {
			t.Errorf("case %d result shoud be \"%s\" but is \"%s\"", i, c.res, res)
		}
	}
}

// TestRaftDirPath verifies that the Raft data directory path is correctly extracted from the storage configuration.
// It exercises all execution paths in the RaftDataDirPath to ensure that:
//   - The path is correctly extracted on a happy path.
//   - In case of an empty path in config, it returns an empty string and true.
//   - If there are any issues getting the path, it returns an empty string and false.
func TestRaftDirPath(t *testing.T) {
	testRaftPath := "/storage/path/raft"

	testCases := map[string]struct {
		config           *CoreConfig
		expectedRaftPath string
		shouldError      bool
	}{
		"happy-path": {
			config: &CoreConfig{
				RawConfig: &server.Config{
					Storage: &server.Storage{
						Type: "raft",
						Config: map[string]string{
							"path": testRaftPath,
						},
					},
				},
			},
			expectedRaftPath: testRaftPath,
		},
		"empty-raft-data-dir-path": {
			config: &CoreConfig{
				RawConfig: &server.Config{
					Storage: &server.Storage{
						Type: "raft",
						Config: map[string]string{
							"path": "",
						},
					},
				},
			},
			expectedRaftPath: "",
			shouldError:      true,
		},
		"no-config": {
			config: &CoreConfig{
				RawConfig: nil,
			},
			expectedRaftPath: "",
			shouldError:      true,
		},
		"no-storage": {
			config: &CoreConfig{
				RawConfig: &server.Config{
					Storage: nil,
				},
			},
			expectedRaftPath: "",
			shouldError:      true,
		},
		"no-storage-config": {
			config: &CoreConfig{
				RawConfig: &server.Config{
					Storage: &server.Storage{
						Type:   "raft",
						Config: nil,
					},
				},
			},
			expectedRaftPath: "",
			shouldError:      true,
		},
		"no-storage-type": {
			config: &CoreConfig{
				RawConfig: &server.Config{
					Storage: &server.Storage{
						Type: "",
						Config: map[string]string{
							"path": testRaftPath,
						},
					},
				},
			},
			expectedRaftPath: "",
			shouldError:      true,
		},
		"no-raft-data-dir-path-in-config": {
			config: &CoreConfig{
				RawConfig: &server.Config{
					Storage: &server.Storage{
						Type:   "raft",
						Config: map[string]string{},
					},
				},
			},
			expectedRaftPath: "",
			shouldError:      true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			fmt.Println("name: ", name)
			core, _, _ := TestCoreUnsealedWithConfig(t, tc.config)
			raftPath, ok := core.RaftDataDirPath()
			if tc.shouldError {
				require.False(t, ok)
				require.Empty(t, raftPath)
			} else {
				require.True(t, ok)
				require.Equal(t, tc.expectedRaftPath, raftPath)
			}
		})
	}
}
