// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cert

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/go-rootcerts"
	"github.com/hashicorp/vault/api"
)

func TestLogin(t *testing.T) {
	certDir := createCertificatePairs(t)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	ln := runTestServer(t, certDir, func(w http.ResponseWriter, r *http.Request) {
		wg.Done()

		if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
			t.Fatalf("no client cert provided")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"auth": {"client_token": "test-token"}}`))
	})

	// Create a new CertAuth struct.
	auth, err := NewCertAuth(
		"role-name",
		filepath.Join(certDir, "client_cert.pem"),
		filepath.Join(certDir, "client_key.pem"),
		WithCACert(filepath.Join(certDir, "ca_cert.pem")),
	)
	if err != nil {
		t.Fatalf("failed to create CertAuth: %v", err)
	}

	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("https://%s", ln.Addr())

	client, err := api.NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	secret, err := auth.Login(context.Background(), client)
	if err != nil {
		t.Fatalf("failed to login: %v", err)
	}

	if secret == nil || secret.Auth == nil || secret.Auth.ClientToken != "test-token" {
		t.Fatalf("unexpected response: %v", secret)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		t.Log("done")
	case <-time.After(time.Second):
		t.Error("timeout")
	}
}

func runTestServer(t *testing.T, certDir string, fn http.HandlerFunc) net.Listener {
	certPool, err := rootcerts.LoadCACerts(&rootcerts.Config{
		CAPath: filepath.Join(certDir, "ca_cert.pem"),
	})
	if err != nil {
		t.Fatalf("failed to load CA certs: %v", err)
	}

	tlsConfig := &tls.Config{
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  certPool,
		RootCAs:    certPool,
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to create listener: %v", err)
	}

	sv := &http.Server{
		TLSConfig: tlsConfig,
		Handler:   http.HandlerFunc(fn),
	}

	t.Cleanup(func() {
		sv.Close()
		ln.Close()
	})

	go func() {
		sv.ServeTLS(ln, filepath.Join(certDir, "server_cert.pem"), filepath.Join(certDir, "server_key.pem"))
	}()

	return ln
}

func createCertificatePairs(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "api_certauth_test")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("test %s, temp dir %s", t.Name(), tempDir)
	caCertTemplate := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: "localhost",
		},
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
		KeyUsage:              x509.KeyUsage(x509.KeyUsageCertSign | x509.KeyUsageCRLSign),
		SerialNumber:          big.NewInt(mathrand.Int63()),
		NotBefore:             time.Now().Add(-30 * time.Second),
		NotAfter:              time.Now().Add(262980 * time.Hour),
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	caBytes, err := x509.CreateCertificate(rand.Reader, caCertTemplate, caCertTemplate, caKey.Public(), caKey)
	if err != nil {
		t.Fatal(err)
	}
	caCert, err := x509.ParseCertificate(caBytes)
	if err != nil {
		t.Fatal(err)
	}
	caCertPEMBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	}
	err = os.WriteFile(filepath.Join(tempDir, "ca_cert.pem"), pem.EncodeToMemory(caCertPEMBlock), 0o755)
	if err != nil {
		t.Fatal(err)
	}
	marshaledCAKey, err := x509.MarshalECPrivateKey(caKey)
	if err != nil {
		t.Fatal(err)
	}
	caKeyPEMBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: marshaledCAKey,
	}
	err = os.WriteFile(filepath.Join(tempDir, "ca_key.pem"), pem.EncodeToMemory(caKeyPEMBlock), 0o755)
	if err != nil {
		t.Fatal(err)
	}

	createCertificate := func(prefix string) {
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			t.Fatal(err)
		}

		template := &x509.Certificate{
			Subject: pkix.Name{
				CommonName: "example.com",
			},
			EmailAddresses: []string{"valid@example.com"},
			IPAddresses:    []net.IP{net.ParseIP("127.0.0.1")},
			ExtKeyUsage: []x509.ExtKeyUsage{
				x509.ExtKeyUsageServerAuth,
				x509.ExtKeyUsageClientAuth,
			},
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement,
			SerialNumber: big.NewInt(mathrand.Int63()),
			NotBefore:    time.Now().Add(-30 * time.Second),
			NotAfter:     time.Now().Add(262980 * time.Hour),
		}
		certBytes, err := x509.CreateCertificate(rand.Reader, template, caCert, key.Public(), caKey)
		if err != nil {
			t.Fatal(err)
		}
		certPEMBlock := &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: certBytes,
		}
		err = os.WriteFile(filepath.Join(tempDir, prefix+"_cert.pem"), pem.EncodeToMemory(certPEMBlock), 0o755)
		if err != nil {
			t.Fatal(err)
		}
		marshaledKey, err := x509.MarshalECPrivateKey(key)
		if err != nil {
			t.Fatal(err)
		}
		keyPEMBlock := &pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: marshaledKey,
		}
		err = os.WriteFile(filepath.Join(tempDir, prefix+"_key.pem"), pem.EncodeToMemory(keyPEMBlock), 0o755)
		if err != nil {
			t.Fatal(err)
		}
	}

	createCertificate("client")
	createCertificate("server")

	return tempDir
}
