// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package listenerutil

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	osuser "os/user"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/vault/internalshared/configutil"
)

func TestTLSConfigReloadsClientCA(t *testing.T) {
	caPath := filepath.Join(t.TempDir(), "client-ca.pem")
	initialCA, err := os.ReadFile("../../helper/serverconfig/test-fixtures/reload/reload_ca.pem")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(caPath, initialCA, 0o600); err != nil {
		t.Fatal(err)
	}

	config, reload, err := TLSConfig(&configutil.Listener{
		TLSCertFile:                   "../../helper/serverconfig/test-fixtures/reload/reload_foo.pem",
		TLSKeyFile:                    "../../helper/serverconfig/test-fixtures/reload/reload_foo.key",
		TLSRequireAndVerifyClientCert: true,
		TLSClientCAFile:               caPath,
	}, make(map[string]string), cli.NewMockUi())
	if err != nil {
		t.Fatal(err)
	}
	if config.GetConfigForClient == nil {
		t.Fatal("expected TLS config to resolve the current client CA pool for each handshake")
	}

	initialConfig, err := config.GetConfigForClient(&tls.ClientHelloInfo{})
	if err != nil {
		t.Fatal(err)
	}
	initialSubjects := initialConfig.ClientCAs.Subjects()
	if len(initialSubjects) != 1 {
		t.Fatalf("expected one initial client CA, got %d", len(initialSubjects))
	}
	clientCertPEM, err := os.ReadFile("../../helper/serverconfig/test-fixtures/reload/reload_foo.pem")
	if err != nil {
		t.Fatal(err)
	}
	block, _ := pem.Decode(clientCertPEM)
	if block == nil {
		t.Fatal("failed to decode client certificate fixture")
	}
	clientCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatal(err)
	}
	verifyClientCert := func(pool *x509.CertPool) error {
		_, err := clientCert.Verify(x509.VerifyOptions{
			Roots:     pool,
			KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		})
		return err
	}
	if err := verifyClientCert(initialConfig.ClientCAs); err != nil {
		t.Fatalf("initial client CA did not trust the fixture certificate: %v", err)
	}

	replacementCA, err := os.ReadFile("../../api/test-fixtures/root/rootcacert.pem")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(caPath, replacementCA, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := reload(); err != nil {
		t.Fatal(err)
	}

	reloadedConfig, err := config.GetConfigForClient(&tls.ClientHelloInfo{})
	if err != nil {
		t.Fatal(err)
	}
	reloadedSubjects := reloadedConfig.ClientCAs.Subjects()
	if len(reloadedSubjects) != 1 {
		t.Fatalf("expected one reloaded client CA, got %d", len(reloadedSubjects))
	}
	if err := verifyClientCert(reloadedConfig.ClientCAs); err == nil {
		t.Fatal("reloaded client CA pool still trusted a certificate signed by the previous CA")
	}

	if err := os.WriteFile(caPath, []byte("invalid CA"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := reload(); err == nil {
		t.Fatal("expected reload to reject an invalid client CA file")
	}
	afterFailure, err := config.GetConfigForClient(&tls.ClientHelloInfo{})
	if err != nil {
		t.Fatal(err)
	}
	if err := verifyClientCert(afterFailure.ClientCAs); err == nil {
		t.Fatal("failed CA reload changed the active client CA pool")
	}
}

func TestUnixSocketListener(t *testing.T) {
	t.Run("ids", func(t *testing.T) {
		socket, err := ioutil.TempFile("", "socket")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(socket.Name())

		uid, gid := os.Getuid(), os.Getgid()

		u, err := osuser.LookupId(strconv.Itoa(uid))
		if err != nil {
			t.Fatal(err)
		}
		user := u.Username

		g, err := osuser.LookupGroupId(strconv.Itoa(gid))
		if err != nil {
			t.Fatal(err)
		}
		group := g.Name

		l, err := UnixSocketListener(socket.Name(), &UnixSocketsConfig{
			User:  user,
			Group: group,
			Mode:  "644",
		})
		if err != nil {
			t.Fatal(err)
		}
		defer l.Close()

		fi, err := os.Stat(socket.Name())
		if err != nil {
			t.Fatal(err)
		}

		mode, err := strconv.ParseUint("644", 8, 32)
		if err != nil {
			t.Fatal(err)
		}
		if fi.Mode().Perm() != os.FileMode(mode) {
			t.Fatalf("failed to set permissions on the socket file")
		}
	})
	t.Run("names", func(t *testing.T) {
		socket, err := ioutil.TempFile("", "socket")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(socket.Name())

		uid, gid := os.Getuid(), os.Getgid()
		l, err := UnixSocketListener(socket.Name(), &UnixSocketsConfig{
			User:  strconv.Itoa(uid),
			Group: strconv.Itoa(gid),
			Mode:  "644",
		})
		if err != nil {
			t.Fatal(err)
		}
		defer l.Close()

		fi, err := os.Stat(socket.Name())
		if err != nil {
			t.Fatal(err)
		}

		mode, err := strconv.ParseUint("644", 8, 32)
		if err != nil {
			t.Fatal(err)
		}
		if fi.Mode().Perm() != os.FileMode(mode) {
			t.Fatalf("failed to set permissions on the socket file")
		}
	})
}
