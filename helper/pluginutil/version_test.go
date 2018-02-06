package pluginutil

import (
	"os"
	"testing"
)

func TestGRPCSupport(t *testing.T) {
	cases := []struct {
		envVersion string
		expected   bool
	}{
		{
			"0.8.3",
			false,
		},
		{
			"0.9.2",
			false,
		},
		{
			"0.9.3",
			false,
		},
		{
			"0.9.4+ent",
			true,
		},
		{
			"0.9.4-beta",
			false,
		},
		{
			"0.9.4",
			true,
		},
		{
			"unknown",
			true,
		},
		{
			"",
			false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.envVersion, func(t *testing.T) {
			err := os.Setenv(PluginVaultVersionEnv, tc.envVersion)
			if err != nil {
				t.Fatal(err)
			}

			result := GRPCSupport()

			if result != tc.expected {
				t.Fatalf("got: %t, expected: %t", result, tc.expected)
			}
		})
	}
}
