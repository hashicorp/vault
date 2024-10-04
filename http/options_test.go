// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/version"
	"github.com/stretchr/testify/require"
)

// TestOptions_Default ensures that the default values are as expected.
func TestOptions_Default(t *testing.T) {
	opts := getDefaultOptions()
	require.NotNil(t, opts)
	require.Equal(t, "", opts.withRedactionValue)
}

// TestOptions_WithRedactionValue ensures that we set the correct value to use for
// redaction when required.
func TestOptions_WithRedactionValue(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		Value           string
		ExpectedValue   string
		IsErrorExpected bool
	}{
		"empty": {
			Value:           "",
			ExpectedValue:   "",
			IsErrorExpected: false,
		},
		"whitespace": {
			Value:           "     ",
			ExpectedValue:   "     ",
			IsErrorExpected: false,
		},
		"value": {
			Value:           "*****",
			ExpectedValue:   "*****",
			IsErrorExpected: false,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			opts := &listenerConfigOptions{}
			applyOption := WithRedactionValue(tc.Value)
			err := applyOption(opts)
			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
			default:
				require.NoError(t, err)
				require.Equal(t, tc.ExpectedValue, opts.withRedactionValue)
			}
		})
	}
}

// TestOptions_WithRedactAddresses ensures that the option works as intended.
func TestOptions_WithRedactAddresses(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		Value         bool
		ExpectedValue bool
	}{
		"true": {
			Value:         true,
			ExpectedValue: true,
		},
		"false": {
			Value:         false,
			ExpectedValue: false,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			opts := &listenerConfigOptions{}
			applyOption := WithRedactAddresses(tc.Value)
			err := applyOption(opts)
			require.NoError(t, err)
			require.Equal(t, tc.ExpectedValue, opts.withRedactAddresses)
		})
	}
}

// TestOptions_WithRedactClusterName ensures that the option works as intended.
func TestOptions_WithRedactClusterName(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		Value         bool
		ExpectedValue bool
	}{
		"true": {
			Value:         true,
			ExpectedValue: true,
		},
		"false": {
			Value:         false,
			ExpectedValue: false,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			opts := &listenerConfigOptions{}
			applyOption := WithRedactClusterName(tc.Value)
			err := applyOption(opts)
			require.NoError(t, err)
			require.Equal(t, tc.ExpectedValue, opts.withRedactClusterName)
		})
	}
}

// TestOptions_WithRedactVersion ensures that the option works as intended.
func TestOptions_WithRedactVersion(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		Value         bool
		ExpectedValue bool
	}{
		"true": {
			Value:         true,
			ExpectedValue: true,
		},
		"false": {
			Value:         false,
			ExpectedValue: false,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			opts := &listenerConfigOptions{}
			applyOption := WithRedactVersion(tc.Value)
			err := applyOption(opts)
			require.NoError(t, err)
			require.Equal(t, tc.ExpectedValue, opts.withRedactVersion)
		})
	}
}

// TestRedactVersionListener tests that the version will be redacted
// from e.g. sys/health and the OpenAPI response if `redact_version`
// is set on the listener.
func TestRedactVersionListener(t *testing.T) {
	conf := &vault.CoreConfig{
		EnableUI:        false,
		EnableRaw:       true,
		BuiltinRegistry: corehelpers.NewMockBuiltinRegistry(),
	}
	core, _, token := vault.TestCoreUnsealedWithConfig(t, conf)

	// Setup listener without redaction
	ln, addr := TestListener(t)
	props := &vault.HandlerProperties{
		Core: core,
		ListenerConfig: &configutil.Listener{
			RedactVersion: false,
		},
	}
	TestServerWithListenerAndProperties(t, ln, addr, core, props)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	testRedactVersionEndpoints(t, addr, token, version.Version)

	// Setup listener with redaction
	ln, addr = TestListener(t)
	props.ListenerConfig.RedactVersion = true
	TestServerWithListenerAndProperties(t, ln, addr, core, props)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	testRedactVersionEndpoints(t, addr, token, "")
}

// testRedactVersionEndpoints tests the endpoints containing versions
// contain the expected version
func testRedactVersionEndpoints(t *testing.T, addr, token, expectedVersion string) {
	client := cleanhttp.DefaultClient()
	req, err := http.NewRequest("GET", addr+"/v1/auth/token?help=1", nil)
	require.NoError(t, err)

	req.Header.Set(consts.AuthHeaderName, token)
	resp, err := client.Do(req)
	require.NoError(t, err)

	testResponseStatus(t, resp, 200)

	var actual map[string]interface{}
	testResponseBody(t, resp, &actual)

	require.NotNil(t, actual["openapi"])
	openAPI, ok := actual["openapi"].(map[string]interface{})
	require.True(t, ok)

	require.NotNil(t, openAPI["info"])
	info, ok := openAPI["info"].(map[string]interface{})
	require.True(t, ok)

	require.NotNil(t, info["version"])
	version, ok := info["version"].(string)
	require.True(t, ok)
	require.Equal(t, expectedVersion, version)

	req, err = http.NewRequest("GET", addr+"/v1/sys/internal/specs/openapi", nil)
	require.NoError(t, err)

	req.Header.Set(consts.AuthHeaderName, "")
	resp, err = client.Do(req)
	require.NoError(t, err)

	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)

	require.NotNil(t, actual["info"])
	info, ok = openAPI["info"].(map[string]interface{})
	require.True(t, ok)

	require.NotNil(t, info["version"])
	version, ok = info["version"].(string)
	require.True(t, ok)
	require.Equal(t, expectedVersion, version)

	req, err = http.NewRequest("GET", addr+"/v1/sys/health", nil)
	require.NoError(t, err)

	req.Header.Set(consts.AuthHeaderName, "")
	resp, err = client.Do(req)
	require.NoError(t, err)

	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)

	require.NotNil(t, actual["version"])
	version, ok = actual["version"].(string)
	require.True(t, ok)

	// sys/health is special and uses a different format to the OpenAPI
	// version.GetVersion().VersionNumber() instead of version.Version
	// We use substring to make sure the check works anyway.
	// In practice, version.GetVersion().VersionNumber() will give something like 1.17.0-beta1
	// and version.Version gives something like 1.17.0
	require.Truef(t, strings.HasPrefix(version, expectedVersion), "version was not as expected, version=%s, expectedVersion=%s",
		version, expectedVersion)
}
