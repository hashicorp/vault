package command

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogFlags_ValuesProvider(t *testing.T) {
	cases := map[string]struct {
		flagKey   string
		envVarKey string
		wantValue string
		wantFound bool
	}{
		"flag-missing": {
			flagKey:   "invalid",
			envVarKey: "valid-env-var",
			wantValue: "envVarValue",
			wantFound: true,
		},
		"envVar-missing": {
			flagKey:   "valid-flag",
			envVarKey: "invalid",
			wantValue: "flagValue",
			wantFound: true,
		},
		"all-present": {
			flagKey:   "valid-flag",
			envVarKey: "valid-env-var",
			wantValue: "flagValue",
			wantFound: true,
		},
		"all-missing": {
			flagKey:   "invalid",
			envVarKey: "invalid",
			wantValue: "",
			wantFound: false,
		},
	}

	// Sneaky little fake providers
	flagFaker := func(key string) (flag.Value, bool) {
		var result fakeFlag
		var found bool

		if key == "valid-flag" {
			result.Set("flagValue")
			found = true
		}

		return &result, found
	}

	envFaker := func(key string) (string, bool) {
		var found bool
		var result string

		if key == "valid-env-var" {
			result = "envVarValue"
			found = true
		}

		return result, found
	}

	vp := valuesProvider{
		flagProvider:   flagFaker,
		envVarProvider: envFaker,
	}

	for name, tc := range cases {
		val, found := vp.overrideValue(tc.flagKey, tc.envVarKey)
		assert.Equal(t, tc.wantFound, found, name)
		assert.Equal(t, tc.wantValue, val, name)
	}
}

type fakeFlag struct {
	value string
}

func (v *fakeFlag) String() string {
	return v.value
}

func (v *fakeFlag) Set(raw string) error {
	v.value = raw
	return nil
}
