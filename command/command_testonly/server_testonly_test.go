// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build testonly

// NOTE: we can't use this with HSM. We can't set testing mode on and it's not
// safe to use env vars since that provides an attack vector in the real world.
//
// The server tests have a go-metrics/exp manager race condition :(.

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
		name        string
		configAfter string
		disabled    bool
	}{
		{
			"enable after default",
			baseHCL + requestLimiterEnableHCL,
			false,
		},
		{
			"enable after enable",
			baseHCL + requestLimiterEnableHCL,
			false,
		},
		{
			"disable after enable",
			baseHCL + requestLimiterDisableHCL,
			true,
		},
		{
			"default after disable",
			baseHCL,
			false,
		},
		{
			"default after default",
			baseHCL,
			false,
		},
		{
			"disable after default",
			baseHCL + requestLimiterDisableHCL,
			true,
		},
		{
			"disable after disable",
			baseHCL + requestLimiterDisableHCL,
			true,
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

	verifyLimiters := func(t *testing.T, expectedDisabled bool) {
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

		switch expectedDisabled {
		case true:
			require.Equal(t, disabledResponse, limiters)
		default:
			require.Equal(t, enabledResponse, limiters)
		}
	}

	verifyLimiters(t, false)

	for _, tc := range cases {
		tc := tc
		// Check that we default on
		t.Run(tc.name, func(t *testing.T) {
			// write the new contents and reload the server
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
			verifyLimiters(t, tc.disabled)
		})
	}
}
