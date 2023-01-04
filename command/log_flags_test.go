package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogFlags_ValuesProvider(t *testing.T) {
	cases := map[string]struct {
		flagKey   string
		envVarKey string
		current   string
		fallback  string
		want      string
	}{
		"only-fallback": {
			flagKey:   "invalid",
			envVarKey: "invalid",
			current:   "",
			fallback:  "foo",
			want:      "foo",
		},
		"only-config": {
			flagKey:   "invalid",
			envVarKey: "invalid",
			current:   "bar",
			fallback:  "",
			want:      "bar",
		},
		"flag-missing": {
			flagKey:   "invalid",
			envVarKey: "valid-env-var",
			current:   "my-config-value1",
			fallback:  "",
			want:      "envVarValue",
		},
		"envVar-missing": {
			flagKey:   "valid-flag",
			envVarKey: "invalid",
			current:   "my-config-value1",
			fallback:  "",
			want:      "flagValue",
		},
		"all-present": {
			flagKey:   "valid-flag",
			envVarKey: "valid-env-var",
			current:   "my-config-value1",
			fallback:  "foo",
			want:      "flagValue",
		},
	}

	// Sneaky little fake provider
	fakeProvider := func(key string) (string, bool) {
		switch key {
		case "valid-flag":
			return "flagValue", true
		case "valid-env-var":
			return "envVarValue", true
		}

		return "", false
	}

	vp := valuesProvider{
		flagProvider:   fakeProvider,
		envVarProvider: fakeProvider,
	}

	for _, tc := range cases {
		got := vp.getAggregatedConfigValue(tc.flagKey, tc.envVarKey, tc.current, tc.fallback)
		assert.Equal(t, tc.want, got)
	}
}
