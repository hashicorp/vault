// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricFilterConfigs(t *testing.T) {
	t.Parallel()
	cases := []struct {
		configFile            string
		expectedFilterDefault *bool
		expectedPrefixFilter  []string
	}{
		{
			"./test-fixtures/telemetry/valid_prefix_filter.hcl",
			nil,
			[]string{"-vault.expire", "-vault.audit", "+vault.expire.num_irrevocable_leases"},
		},
		{
			"./test-fixtures/telemetry/filter_default_override.hcl",
			boolPointer(false),
			[]string(nil),
		},
	}
	t.Run("validate metric filter configs", func(t *testing.T) {
		t.Parallel()

		for _, tc := range cases {
			config, err := LoadConfigFile(tc.configFile)
			if err != nil {
				t.Fatalf("Error encountered when loading config %+v", err)
			}

			assert.Equal(t, tc.expectedFilterDefault, config.SharedConfig.Telemetry.FilterDefault)
			assert.Equal(t, tc.expectedPrefixFilter, config.SharedConfig.Telemetry.PrefixFilter)
		}
	})
}

// TestRollbackMountPointMetricsConfig verifies that the add_mount_point_rollback_metrics
// config option is parsed correctly, when it is set to true. Also verifies that
// the default for this setting is false
func TestRollbackMountPointMetricsConfig(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		configFile     string
		wantMountPoint bool
	}{
		{
			name:           "include mount point",
			configFile:     "./test-fixtures/telemetry/rollback_mount_point.hcl",
			wantMountPoint: true,
		},
		{
			name:           "exclude mount point",
			configFile:     "./test-fixtures/telemetry/valid_prefix_filter.hcl",
			wantMountPoint: false,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			config, err := LoadConfigFile(tc.configFile)
			require.NoError(t, err)
			require.Equal(t, tc.wantMountPoint, config.Telemetry.RollbackMetricsIncludeMountPoint)
		})
	}
}
