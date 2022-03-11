package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
