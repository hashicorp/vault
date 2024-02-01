// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build testonly

package command_testonly

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command"
	"github.com/hashicorp/vault/limits"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/require"
)

func init() {
	if signed := os.Getenv("VAULT_LICENSE_CI"); signed != "" {
		os.Setenv(command.EnvVaultLicense, signed)
	}
}

const (
	baseHCL = `
		backend "inmem" { }
		disable_mlock = true
		listener "tcp" {
			address     = "127.0.0.1:8209"
			tls_disable = "true"
		}
		api_addr = "http://127.0.0.1:8209"
	`
	requestLimiterDisableHCL = `
  request_limiter {
	disable = true
  }
`
	requestLimiterEnableHCL = `
  request_limiter {
	disable = false
  }
`
)

// TestServer_ReloadRequestLimiter tests a series of reloads and state
// transitions between RequestLimiter enable and disable.
func TestServer_ReloadRequestLimiter(t *testing.T) {
	t.Parallel()

	enabledResponse := &vault.RequestLimiterResponse{
		GlobalDisabled:   false,
		ListenerDisabled: false,
		Limiters: map[string]*vault.LimiterStatus{
			limits.WriteLimiter: {
				Enabled: true,
				Flags:   limits.DefaultLimiterFlags[limits.WriteLimiter],
			},
			limits.SpecialPathLimiter: {
				Enabled: true,
				Flags:   limits.DefaultLimiterFlags[limits.SpecialPathLimiter],
			},
		},
	}

	disabledResponse := &vault.RequestLimiterResponse{
		GlobalDisabled:   true,
		ListenerDisabled: false,
		Limiters: map[string]*vault.LimiterStatus{
			limits.WriteLimiter: {
				Enabled: false,
			},
			limits.SpecialPathLimiter: {
				Enabled: false,
			},
		},
	}

	cases := []struct {
		name             string
		configAfter      string
		expectedResponse *vault.RequestLimiterResponse
	}{
		{
			"enable after default",
			baseHCL + requestLimiterEnableHCL,
			enabledResponse,
		},
		{
			"enable after enable",
			baseHCL + requestLimiterEnableHCL,
			enabledResponse,
		},
		{
			"disable after enable",
			baseHCL + requestLimiterDisableHCL,
			disabledResponse,
		},
		{
			"default after disable",
			baseHCL,
			enabledResponse,
		},
		{
			"default after default",
			baseHCL,
			enabledResponse,
		},
		{
			"disable after default",
			baseHCL + requestLimiterDisableHCL,
			disabledResponse,
		},
		{
			"disable after disable",
			baseHCL + requestLimiterDisableHCL,
			disabledResponse,
		},
	}

	ui, srv := command.TestServerCommand(t)

	f, err := os.CreateTemp(t.TempDir(), "")
	require.NoErrorf(t, err, "error creating temp dir: %v", err)

	_, err = f.WriteString(baseHCL)
	require.NoErrorf(t, err, "cannot write temp file contents")

	configPath := f.Name()

	var output string
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		code := srv.Run([]string{"-config", configPath})
		output = ui.ErrorWriter.String() + ui.OutputWriter.String()
		require.Equal(t, 0, code, output)
	}()

	select {
	case <-srv.StartedCh():
	case <-time.After(5 * time.Second):
		t.Fatalf("timeout")
	}
	defer func() {
		srv.ShutdownCh <- struct{}{}
		wg.Wait()
	}()

	err = f.Close()
	require.NoErrorf(t, err, "unable to close temp file")

	// create a client and unseal vault
	cli, err := srv.Client()
	require.NoError(t, err)
	require.NoError(t, cli.SetAddress("http://127.0.0.1:8209"))
	initResp, err := cli.Sys().Init(&api.InitRequest{SecretShares: 1, SecretThreshold: 1})
	require.NoError(t, err)
	_, err = cli.Sys().Unseal(initResp.Keys[0])
	require.NoError(t, err)
	cli.SetToken(initResp.RootToken)

	output = ui.ErrorWriter.String() + ui.OutputWriter.String()
	require.Contains(t, output, "Request Limiter: enabled")

	verifyLimiters := func(t *testing.T, expectedResponse *vault.RequestLimiterResponse) {
		t.Helper()

		statusResp, err := cli.Logical().Read("/sys/internal/request-limiter/status")
		require.NoError(t, err)
		require.NotNil(t, statusResp)

		limitersResp, ok := statusResp.Data["request_limiter"]
		require.True(t, ok)
		require.NotNil(t, limitersResp)

		var limiters *vault.RequestLimiterResponse
		err = mapstructure.Decode(limitersResp, &limiters)
		require.NoError(t, err)
		require.NotNil(t, limiters)

		require.Equal(t, expectedResponse, limiters)
	}

	// Start off with default enabled
	verifyLimiters(t, enabledResponse)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Write the new contents and reload the server
			f, err = os.OpenFile(configPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
			require.NoError(t, err)
			defer f.Close()

			_, err = f.WriteString(tc.configAfter)
			require.NoErrorf(t, err, "cannot write temp file contents")

			srv.SighupCh <- struct{}{}
			select {
			case <-srv.ReloadedCh():
			case <-time.After(5 * time.Second):
				t.Fatalf("test timed out")
			}

			verifyLimiters(t, tc.expectedResponse)
		})
	}
}
