// +build !race

package command

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/vault/meta"
	"github.com/mitchellh/cli"
)

var (
	basehcl = `
disable_mlock = true

listener "tcp" {
  address = "127.0.0.1:8200"
  tls_disable = "true"
}
`

	consulhcl = `
backend "consul" {
    prefix = "foo/"
    advertise_addr = "http://127.0.0.1:8200"
    disable_registration = "true"
}
`
	haconsulhcl = `
ha_backend "consul" {
    prefix = "bar/"
    redirect_addr = "http://127.0.0.1:8200"
    disable_registration = "true"
}
`

	badhaconsulhcl = `
ha_backend "file" {
    path = "/dev/null"
}
`

	reloadhcl = `
backend "file" {
    path = "/dev/null"
}

disable_mlock = true

listener "tcp" {
    address = "127.0.0.1:8203"
    tls_cert_file = "TMPDIR/reload_FILE.pem"
    tls_key_file = "TMPDIR/reload_FILE.key"
}
`
)

// The following tests have a go-metrics/exp manager race condition
func TestServer_ReloadListener(t *testing.T) {
	wd, _ := os.Getwd()
	wd += "/server/test-fixtures/reload/"

	td, err := ioutil.TempDir("", fmt.Sprintf("vault-test-%d", rand.New(rand.NewSource(time.Now().Unix())).Int63))
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(td)

	wg := &sync.WaitGroup{}

	// Setup initial certs
	inBytes, _ := ioutil.ReadFile(wd + "reload_foo.pem")
	ioutil.WriteFile(td+"/reload_foo.pem", inBytes, 0777)
	inBytes, _ = ioutil.ReadFile(wd + "reload_foo.key")
	ioutil.WriteFile(td+"/reload_foo.key", inBytes, 0777)
	inBytes, _ = ioutil.ReadFile(wd + "reload_bar.pem")
	ioutil.WriteFile(td+"/reload_bar.pem", inBytes, 0777)
	inBytes, _ = ioutil.ReadFile(wd + "reload_bar.key")
	ioutil.WriteFile(td+"/reload_bar.key", inBytes, 0777)

	relhcl := strings.Replace(strings.Replace(reloadhcl, "TMPDIR", td, -1), "FILE", "foo", -1)
	ioutil.WriteFile(td+"/reload.hcl", []byte(relhcl), 0777)

	inBytes, _ = ioutil.ReadFile(wd + "reload_ca.pem")
	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(inBytes)
	if !ok {
		t.Fatal("not ok when appending CA cert")
	}

	ui := new(cli.MockUi)
	c := &ServerCommand{
		Meta: meta.Meta{
			Ui: ui,
		},
		ShutdownCh: MakeShutdownCh(),
		SighupCh:   MakeSighupCh(),
	}

	finished := false
	finishedMutex := sync.Mutex{}

	wg.Add(1)
	args := []string{"-config", td + "/reload.hcl"}
	go func() {
		if code := c.Run(args); code != 0 {
			t.Error("got a non-zero exit status")
		}
		finishedMutex.Lock()
		finished = true
		finishedMutex.Unlock()
		wg.Done()
	}()

	checkFinished := func() {
		finishedMutex.Lock()
		if finished {
			t.Fatalf(fmt.Sprintf("finished early; relhcl was\n%s\nstdout was\n%s\nstderr was\n%s\n", relhcl, ui.OutputWriter.String(), ui.ErrorWriter.String()))
		}
		finishedMutex.Unlock()
	}

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

	checkFinished()
	time.Sleep(5 * time.Second)
	checkFinished()

	if err := testCertificateName("foo.example.com"); err != nil {
		t.Fatalf("certificate name didn't check out: %s", err)
	}

	relhcl = strings.Replace(strings.Replace(reloadhcl, "TMPDIR", td, -1), "FILE", "bar", -1)
	ioutil.WriteFile(td+"/reload.hcl", []byte(relhcl), 0777)

	c.SighupCh <- struct{}{}
	checkFinished()
	time.Sleep(2 * time.Second)
	checkFinished()

	if err := testCertificateName("bar.example.com"); err != nil {
		t.Fatalf("certificate name didn't check out: %s", err)
	}

	c.ShutdownCh <- struct{}{}

	wg.Wait()
}
