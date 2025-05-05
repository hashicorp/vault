// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package configutil

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestMerge tests SharedConfig#Merge
func TestMerge(t *testing.T) {
	tests := []struct {
		name   string
		c1     *SharedConfig
		c2     *SharedConfig
		expect *SharedConfig
	}{
		{
			"nil second config",
			&SharedConfig{},
			nil,
			&SharedConfig{},
		},
		{
			"blank configs",
			&SharedConfig{},
			&SharedConfig{},
			&SharedConfig{},
		},
		{
			"combined FoundKeys",
			&SharedConfig{
				FoundKeys: []string{
					"DisableMlock",
					"ClusterName",
				},
			},
			&SharedConfig{
				FoundKeys: []string{
					"PidFile",
					"DefaultMaxRequestDuration",
				},
			},
			&SharedConfig{
				FoundKeys: []string{
					"DisableMlock",
					"ClusterName",
					"PidFile",
					"DefaultMaxRequestDuration",
				},
			},
		},
		{
			"Entropy Overwrite",
			&SharedConfig{
				Entropy: &Entropy{
					Mode:     EntropyAugmentation,
					SealName: "seal-to-be-overwritten",
				},
			},
			&SharedConfig{
				Entropy: &Entropy{
					Mode:     EntropyAugmentation,
					SealName: "new-seal",
				},
			},
			&SharedConfig{
				Entropy: &Entropy{
					Mode:     EntropyAugmentation,
					SealName: "new-seal",
				},
			},
		},
		{
			"DisableMlock true overrides false",
			&SharedConfig{
				DisableMlock: false,
			},
			&SharedConfig{
				DisableMlock: true,
			},
			&SharedConfig{
				DisableMlock: true,
			},
		},
		{
			"longer duration overrides shorter",
			&SharedConfig{
				DefaultMaxRequestDuration: time.Duration(5 * time.Second),
			},
			&SharedConfig{
				DefaultMaxRequestDuration: time.Duration(10 * time.Second),
			},
			&SharedConfig{
				DefaultMaxRequestDuration: time.Duration(10 * time.Second),
			},
		},
		{
			"combined listeners",
			&SharedConfig{
				Listeners: []*Listener{
					{
						Type:    TCP,
						Address: "127.0.0.1",
					},
					{
						Type:    Unix,
						Address: "mnt/listener",
					},
				},
			},
			&SharedConfig{
				Listeners: []*Listener{
					{
						Type:    TCP,
						Address: "127.0.0.3",
					},
					{
						Type:    Unix,
						Address: "mnt/listener2",
					},
				},
			},
			&SharedConfig{
				Listeners: []*Listener{
					{
						Type:    TCP,
						Address: "127.0.0.1",
					},
					{
						Type:    Unix,
						Address: "mnt/listener",
					},
					{
						Type:    TCP,
						Address: "127.0.0.3",
					},
					{
						Type:    Unix,
						Address: "mnt/listener2",
					},
				},
			},
		},
		{
			"combined user lockouts",
			&SharedConfig{
				UserLockouts: []*UserLockout{
					{
						Type: "lockout1",
					},
				},
			},
			&SharedConfig{
				UserLockouts: []*UserLockout{
					{
						Type: "lockout2",
					},
				},
			},
			&SharedConfig{
				UserLockouts: []*UserLockout{
					{
						Type: "lockout1",
					},
					{
						Type: "lockout2",
					},
				},
			},
		},
		{
			"combined seals",
			&SharedConfig{
				Seals: []*KMS{
					{
						Purpose: []string{"purpose1"},
					},
				},
			},
			&SharedConfig{
				Seals: []*KMS{
					{
						Purpose: []string{"purpose2"},
					},
				},
			},
			&SharedConfig{
				Seals: []*KMS{
					{
						Purpose: []string{"purpose1"},
					},
					{
						Purpose: []string{"purpose2"},
					},
				},
			},
		},
		{
			"telemetry overwrite",
			&SharedConfig{
				Telemetry: &Telemetry{
					StatsiteAddr: "https://example.com",
				},
			},
			&SharedConfig{
				Telemetry: &Telemetry{},
			},
			&SharedConfig{
				Telemetry: &Telemetry{},
			},
		},
		{
			"HCPLinkConf overwrite",
			&SharedConfig{},
			&SharedConfig{},
			&SharedConfig{},
		},
		{
			"log fields overwrite",
			&SharedConfig{
				LogFile:           "file1.log",
				LogFormat:         "json",
				LogLevel:          "warn",
				LogRotateBytes:    2048,
				LogRotateDuration: "24h",
				LogRotateMaxFiles: 32,
			},
			&SharedConfig{
				LogFile:              "file2.log",
				LogFormat:            "txt",
				LogLevel:             "error",
				LogRotateBytes:       1024,
				LogRotateBytesRaw:    1024,
				LogRotateDuration:    "12h",
				LogRotateMaxFiles:    8,
				LogRotateMaxFilesRaw: 8,
			},
			&SharedConfig{
				LogFile:              "file2.log",
				LogFormat:            "txt",
				LogLevel:             "error",
				LogRotateBytes:       1024,
				LogRotateBytesRaw:    1024,
				LogRotateDuration:    "12h",
				LogRotateMaxFiles:    8,
				LogRotateMaxFilesRaw: 8,
			},
		},
		{
			"log fields raw overwrite",
			&SharedConfig{
				LogFile:           "file1.log",
				LogFormat:         "json",
				LogLevel:          "warn",
				LogRotateBytes:    2048,
				LogRotateDuration: "24h",
				LogRotateMaxFiles: 32,
			},
			&SharedConfig{
				LogRotateBytes:       1024,
				LogRotateBytesRaw:    1024,
				LogRotateMaxFiles:    8,
				LogRotateMaxFilesRaw: 8,
			},
			&SharedConfig{
				LogFile:              "file1.log",
				LogFormat:            "json",
				LogLevel:             "warn",
				LogRotateDuration:    "24h",
				LogRotateBytes:       1024,
				LogRotateBytesRaw:    1024,
				LogRotateMaxFiles:    8,
				LogRotateMaxFilesRaw: 8,
			},
		},
		{
			"pidfile overwrite",
			&SharedConfig{
				PidFile: "file1.pid",
			},
			&SharedConfig{
				PidFile: "file2.pid",
			},
			&SharedConfig{
				PidFile: "file2.pid",
			},
		},
		{
			"cluster name overwrite",
			&SharedConfig{
				ClusterName: "vault-cluster1",
			},
			&SharedConfig{
				ClusterName: "vault-cluster2",
			},
			&SharedConfig{
				ClusterName: "vault-cluster2",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.True(t, reflect.DeepEqual(tc.expect, tc.c1.Merge(tc.c2)))
		})
	}
}
