// +build !race,!hsm

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
	"net"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/vault/physical"
	physInmem "github.com/hashicorp/vault/physical/inmem"
	"github.com/mitchellh/cli"
)

func testRandomPort(tb testing.TB) int {
	tb.Helper()

	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		tb.Fatal(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		tb.Fatal(err)
	}
	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port
}

func testBaseHCL(tb testing.TB) string {
	tb.Helper()

	return strings.TrimSpace(fmt.Sprintf(`
		disable_mlock = true
		listener "tcp" {
		  address     = "127.0.0.1:%d"
		  tls_disable = "true"
		}
	`, testRandomPort(tb)))
}

const (
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
		startedCh:  make(chan struct{}, 5),
		reloadedCh: make(chan struct{}, 5),
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
	ioutil.WriteFile(td+"/reload_cert.pem", inBytes, 0777)
	inBytes, _ = ioutil.ReadFile(wd + "reload_foo.key")
	ioutil.WriteFile(td+"/reload_key.pem", inBytes, 0777)

	relhcl := strings.Replace(reloadHCL, "TMPDIR", td, -1)
	ioutil.WriteFile(td+"/reload.hcl", []byte(relhcl), 0777)

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

	relhcl = strings.Replace(reloadHCL, "TMPDIR", td, -1)
	inBytes, _ = ioutil.ReadFile(wd + "reload_bar.pem")
	ioutil.WriteFile(td+"/reload_cert.pem", inBytes, 0777)
	inBytes, _ = ioutil.ReadFile(wd + "reload_bar.key")
	ioutil.WriteFile(td+"/reload_key.pem", inBytes, 0777)
	ioutil.WriteFile(td+"/reload.hcl", []byte(relhcl), 0777)

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
	}{
		{
			"common_ha",
			testBaseHCL(t) + inmemHCL,
			"(HA available)",
			0,
		},
		{
			"separate_ha",
			testBaseHCL(t) + inmemHCL + haInmemHCL,
			"HA Storage:",
			0,
		},
		{
			"bad_separate_ha",
			testBaseHCL(t) + inmemHCL + badHAInmemHCL,
			"Specified HA storage does not support HA",
			1,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ui, cmd := testServerCommand(t)
			f, err := ioutil.TempFile("", "")
			if err != nil {
				t.Fatalf("error creating temp dir: %v", err)
			}
			f.WriteString(tc.contents)
			f.Close()
			defer os.Remove(f.Name())

			code := cmd.Run([]string{
				"-config", f.Name(),
				"-test-verify-only",
			})
			output := ui.ErrorWriter.String() + ui.OutputWriter.String()
			if code != tc.code {
				t.Errorf("expected %d to be %d: %s", code, tc.code, output)
			}

			if !strings.Contains(output, tc.exp) {
				t.Fatalf("expected %q to contain %q", output, tc.exp)
			}
		})
	}
}
