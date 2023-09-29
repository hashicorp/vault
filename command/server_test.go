// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !race && !hsm && !fips_140_3

// NOTE: we can't use this with HSM. We can't set testing mode on and it's not
// safe to use env vars since that provides an attack vector in the real world.
//
// The server tests have a go-metrics/exp manager race condition :(.

package command

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/physical"
	physInmem "github.com/hashicorp/vault/sdk/physical/inmem"
	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	if signed := os.Getenv("VAULT_LICENSE_CI"); signed != "" {
		os.Setenv(EnvVaultLicense, signed)
	}
}

func testBaseHCL(tb testing.TB, listenerExtras string) string {
	tb.Helper()

	return strings.TrimSpace(fmt.Sprintf(`
		disable_mlock = true
		listener "tcp" {
			address     = "127.0.0.1:%d"
			tls_disable = "true"
			%s
		}
	`, 0, listenerExtras))
}

const (
	goodListenerTimeouts = `http_read_header_timeout = 12
			http_read_timeout = "34s"
			http_write_timeout = "56m"
			http_idle_timeout = "78h"`

	badListenerReadHeaderTimeout = `http_read_header_timeout = "12km"`
	badListenerReadTimeout       = `http_read_timeout = "34日"`
	badListenerWriteTimeout      = `http_write_timeout = "56lbs"`
	badListenerIdleTimeout       = `http_idle_timeout = "78gophers"`

	inmemHCL = `
backend "inmem_ha" {
  advertise_addr       = "http://127.0.0.1:8200"
}
`
	haInmemHCL = `
ha_backend "inmem_ha" {
  redirect_addr        = "http://127.0.0.1:8200"
}
`

	badHAInmemHCL = `
ha_backend "inmem" {}
`

	reloadHCL = `
backend "inmem" {}
disable_mlock = true
listener "tcp" {
  address       = "127.0.0.1:8203"
  tls_cert_file = "TMPDIR/reload_cert.pem"
  tls_key_file  = "TMPDIR/reload_key.pem"
}
`
	cloudHCL = `
cloud {
      resource_id = "organization/bc58b3d0-2eab-4ab8-abf4-f61d3c9975ff/project/1c78e888-2142-4000-8918-f933bbbc7690/hashicorp.example.resource/example"
    client_id = "J2TtcSYOyPUkPV2z0mSyDtvitxLVjJmu"
    client_secret = "N9JtHZyOnHrIvJZs82pqa54vd4jnkyU3xCcqhFXuQKJZZuxqxxbP1xCfBZVB82vY"
}
`
)

func testServerCommand(tb testing.TB) (*cli.MockUi, *ServerCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &ServerCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
		ShutdownCh: MakeShutdownCh(),
		SighupCh:   MakeSighupCh(),
		SigUSR2Ch:  MakeSigUSR2Ch(),
		PhysicalBackends: map[string]physical.Factory{
			"inmem":    physInmem.NewInmem,
			"inmem_ha": physInmem.NewInmemHA,
		},

		// These prevent us from random sleep guessing...
		startedCh:         make(chan struct{}, 5),
		reloadedCh:        make(chan struct{}, 5),
		licenseReloadedCh: make(chan error),
	}
}

func TestServer_ReloadListener(t *testing.T) {
	t.Parallel()

	wd, _ := os.Getwd()
	wd += "/server/test-fixtures/reload/"

	td, err := ioutil.TempDir("", "vault-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(td)

	wg := &sync.WaitGroup{}
	// Setup initial certs
	inBytes, _ := ioutil.ReadFile(wd + "reload_foo.pem")
	ioutil.WriteFile(td+"/reload_cert.pem", inBytes, 0o777)
	inBytes, _ = ioutil.ReadFile(wd + "reload_foo.key")
	ioutil.WriteFile(td+"/reload_key.pem", inBytes, 0o777)

	relhcl := strings.ReplaceAll(reloadHCL, "TMPDIR", td)
	ioutil.WriteFile(td+"/reload.hcl", []byte(relhcl), 0o777)

	inBytes, _ = ioutil.ReadFile(wd + "reload_ca.pem")
	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(inBytes)
	if !ok {
		t.Fatal("not ok when appending CA cert")
	}

	ui, cmd := testServerCommand(t)
	_ = ui

	wg.Add(1)
	args := []string{"-config", td + "/reload.hcl"}
	go func() {
		if code := cmd.Run(args); code != 0 {
			output := ui.ErrorWriter.String() + ui.OutputWriter.String()
			t.Errorf("got a non-zero exit status: %s", output)
		}
		wg.Done()
	}()

	testCertificateName := func(cn string) error {
		conn, err := tls.Dial("tcp", "127.0.0.1:8203", &tls.Config{
			RootCAs: certPool,
		})
		if err != nil {
			return err
		}
		defer conn.Close()
		if err = conn.Handshake(); err != nil {
			return err
		}
		servName := conn.ConnectionState().PeerCertificates[0].Subject.CommonName
		if servName != cn {
			return fmt.Errorf("expected %s, got %s", cn, servName)
		}
		return nil
	}

	select {
	case <-cmd.startedCh:
	case <-time.After(5 * time.Second):
		t.Fatalf("timeout")
	}

	if err := testCertificateName("foo.example.com"); err != nil {
		t.Fatalf("certificate name didn't check out: %s", err)
	}

	relhcl = strings.ReplaceAll(reloadHCL, "TMPDIR", td)
	inBytes, _ = ioutil.ReadFile(wd + "reload_bar.pem")
	ioutil.WriteFile(td+"/reload_cert.pem", inBytes, 0o777)
	inBytes, _ = ioutil.ReadFile(wd + "reload_bar.key")
	ioutil.WriteFile(td+"/reload_key.pem", inBytes, 0o777)
	ioutil.WriteFile(td+"/reload.hcl", []byte(relhcl), 0o777)

	cmd.SighupCh <- struct{}{}
	select {
	case <-cmd.reloadedCh:
	case <-time.After(5 * time.Second):
		t.Fatalf("timeout")
	}

	if err := testCertificateName("bar.example.com"); err != nil {
		t.Fatalf("certificate name didn't check out: %s", err)
	}

	cmd.ShutdownCh <- struct{}{}

	wg.Wait()
}

func TestServer(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		contents string
		exp      string
		code     int
		args     []string
	}{
		{
			"common_ha",
			testBaseHCL(t, "") + inmemHCL,
			"(HA available)",
			0,
			[]string{"-test-verify-only"},
		},
		{
			"separate_ha",
			testBaseHCL(t, "") + inmemHCL + haInmemHCL,
			"HA Storage:",
			0,
			[]string{"-test-verify-only"},
		},
		{
			"bad_separate_ha",
			testBaseHCL(t, "") + inmemHCL + badHAInmemHCL,
			"Specified HA storage does not support HA",
			1,
			[]string{"-test-verify-only"},
		},
		{
			"good_listener_timeout_config",
			testBaseHCL(t, goodListenerTimeouts) + inmemHCL,
			"",
			0,
			[]string{"-test-server-config"},
		},
		{
			"bad_listener_read_header_timeout_config",
			testBaseHCL(t, badListenerReadHeaderTimeout) + inmemHCL,
			"unknown unit \"km\" in duration \"12km\"",
			1,
			[]string{"-test-server-config"},
		},
		{
			"bad_listener_read_timeout_config",
			testBaseHCL(t, badListenerReadTimeout) + inmemHCL,
			"unknown unit \"\\xe6\\x97\\xa5\" in duration",
			1,
			[]string{"-test-server-config"},
		},
		{
			"bad_listener_write_timeout_config",
			testBaseHCL(t, badListenerWriteTimeout) + inmemHCL,
			"unknown unit \"lbs\" in duration \"56lbs\"",
			1,
			[]string{"-test-server-config"},
		},
		{
			"bad_listener_idle_timeout_config",
			testBaseHCL(t, badListenerIdleTimeout) + inmemHCL,
			"unknown unit \"gophers\" in duration \"78gophers\"",
			1,
			[]string{"-test-server-config"},
		},
		{
			"environment_variables_logged",
			testBaseHCL(t, "") + inmemHCL,
			"Environment Variables",
			0,
			[]string{"-test-verify-only"},
		},
		{
			"cloud_config",
			testBaseHCL(t, "") + inmemHCL + cloudHCL,
			"HCP Organization: bc58b3d0-2eab-4ab8-abf4-f61d3c9975ff",
			0,
			[]string{"-test-verify-only"},
		},
		{
			"recovery_mode",
			testBaseHCL(t, "") + inmemHCL,
			"",
			0,
			[]string{"-test-verify-only", "-recovery"},
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ui, cmd := testServerCommand(t)

			f, err := os.CreateTemp(t.TempDir(), "")
			require.NoErrorf(t, err, "error creating temp dir: %v", err)

			_, err = f.WriteString(tc.contents)
			require.NoErrorf(t, err, "cannot write temp file contents")

			err = f.Close()
			require.NoErrorf(t, err, "unable to close temp file")

			args := append(tc.args, "-config", f.Name())
			code := cmd.Run(args)
			output := ui.ErrorWriter.String() + ui.OutputWriter.String()
			require.Equal(t, tc.code, code, "expected %d to be %d: %s", code, tc.code, output)
			require.Contains(t, output, tc.exp, "expected %q to contain %q", output, tc.exp)
		})
	}
}

// TestServer_DevTLS verifies that a vault server starts up correctly with the -dev-tls flag
func TestServer_DevTLS(t *testing.T) {
	ui, cmd := testServerCommand(t)
	args := []string{"-dev-tls", "-dev-listen-address=127.0.0.1:0", "-test-server-config"}
	retCode := cmd.Run(args)
	output := ui.ErrorWriter.String() + ui.OutputWriter.String()
	require.Equal(t, 0, retCode, output)
	require.Contains(t, output, `tls: "enabled"`)
}

// TestConfigureDevTLS verifies the various logic paths that flow through the
// configureDevTLS function.
func TestConfigureDevTLS(t *testing.T) {
	testcases := []struct {
		ServerCommand   *ServerCommand
		DeferFuncNotNil bool
		ConfigNotNil    bool
		TLSDisable      bool
		CertPathEmpty   bool
		ErrNotNil       bool
		TestDescription string
	}{
		{
			ServerCommand: &ServerCommand{
				flagDevTLS: false,
			},
			ConfigNotNil:    true,
			TLSDisable:      true,
			CertPathEmpty:   true,
			ErrNotNil:       false,
			TestDescription: "flagDev is false, nothing will be configured",
		},
		{
			ServerCommand: &ServerCommand{
				flagDevTLS:        true,
				flagDevTLSCertDir: "",
			},
			DeferFuncNotNil: true,
			ConfigNotNil:    true,
			ErrNotNil:       false,
			TestDescription: "flagDevTLSCertDir is empty",
		},
		{
			ServerCommand: &ServerCommand{
				flagDevTLS:        true,
				flagDevTLSCertDir: "@/#",
			},
			CertPathEmpty:   true,
			ErrNotNil:       true,
			TestDescription: "flagDevTLSCertDir is set to something invalid",
		},
	}

	for _, testcase := range testcases {
		fun, cfg, certPath, err := configureDevTLS(testcase.ServerCommand)
		if fun != nil {
			// If a function is returned, call it right away to clean up
			// files created in the temporary directory before anything else has
			// a chance to fail this test.
			fun()
		}

		t.Run(testcase.TestDescription, func(t *testing.T) {
			assert.Equal(t, testcase.DeferFuncNotNil, (fun != nil))
			assert.Equal(t, testcase.ConfigNotNil, cfg != nil)
			if testcase.ConfigNotNil && cfg != nil {
				assert.True(t, len(cfg.Listeners) > 0)
				assert.Equal(t, testcase.TLSDisable, cfg.Listeners[0].TLSDisable)
			}
			assert.Equal(t, testcase.CertPathEmpty, len(certPath) == 0)
			if testcase.ErrNotNil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
