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
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	mathrand "math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/go-rootcerts"
	"github.com/hashicorp/vault/api"
)

var extractMountRe = regexp.MustCompile(`^/v1/auth/([a-z]+)/login$`)

// TestLogin tests the login method of the CertAuth struct with the various overrides we support,
// the values the client uses to call the server with influence the returned client token so
// the test can validate we sent what we expected to.
func TestLogin(t *testing.T) {
	certs := createCertificates(t)

	// build up a server that will influence the returned token based on the request
	ln := runTestServer(t, certs, func(w http.ResponseWriter, r *http.Request) {
		// Make sure we have a client certificate
		if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("missing peer certificate"))
			return
		}

		// Make sure we get the expected client certificate
		if r.TLS.PeerCertificates[0].SerialNumber.Cmp(certs.parsedClientCert.SerialNumber) != 0 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("incorrect peer certificate provided"))
			return
		}

		// Extract the mount path from the URL we were called with
		var mount string
		if matches := extractMountRe.FindStringSubmatch(r.URL.Path); len(matches) == 2 {
			mount = matches[1]
		}

		// Parse out the role name from the request body if provided
		rawBody, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("failed reading body: " + err.Error()))
			return
		}
		defer r.Body.Close()

		data := make(map[string]interface{})
		if err = json.Unmarshal(rawBody, &data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("failed parsing body: " + err.Error()))
			return
		}

		roleName := "none"
		if role, ok := data["name"].(string); ok {
			roleName = role
		}

		// Build up our token that includes the mount of our cert auth and the role name
		token := fmt.Sprintf("%s-%s-token", mount, roleName)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"auth": {"client_token": "%s"}}`, token)))
	})
	url := fmt.Sprintf("https://%s", ln.Addr())
	defaultAuthClient, err := NewDefaultCertAuthClient(url, &api.TLSConfig{
		CACert:     certs.caCert,
		ClientCert: certs.clientCert,
		ClientKey:  certs.clientKey,
	})
	if err != nil {
		t.Fatalf("failed to create default CertAuth: %v", err)
	}

	defConfig := api.DefaultConfig()
	defConfig.Address = url
	baseClient, err := api.NewClient(defConfig)
	if err != nil {
		t.Fatalf("failed to create CertAuth client based on base client: %v", err)
	}

	basedOnClient, err := NewCertAuthClient(baseClient, &api.TLSConfig{
		CACert:     certs.caCert,
		ClientCert: certs.clientCert,
		ClientKey:  certs.clientKey,
	})
	if err != nil {
		t.Fatalf("failed to create CertAuth based on another client: %v", err)
	}

	tests := []struct {
		name          string
		opts          []LoginOption
		expectedToken string
	}{
		{"default", []LoginOption{}, "cert-none-token"},
		{"with mount path", []LoginOption{WithMountPath("mycert")}, "mycert-none-token"},
		{"with role", []LoginOption{WithRole("myrole")}, "cert-myrole-token"},
		{"with mount and role", []LoginOption{WithMountPath("mycert"), WithRole("myrole")}, "mycert-myrole-token"},
	}
	for name, certAuthClient := range map[string]*api.Client{"default": defaultAuthClient, "based-on-client": basedOnClient} {
		for _, test := range tests {
			t.Run(fmt.Sprintf("%s-%s", name, test.name), func(t *testing.T) {
				auth, err := NewCertAuth(test.opts...)
				if err != nil {
					t.Fatalf("failed to create CertAuth: %v", err)
				}

				secret, err := auth.Login(context.Background(), certAuthClient)
				if err != nil {
					t.Fatalf("failed to login: %v", err)
				}

				if secret == nil || secret.Auth == nil || secret.Auth.ClientToken != test.expectedToken {
					t.Fatalf("unexpected response: %v", secret)
				}
			})
		}
	}
}

// TestNewDefaultCertAuthClient tests the NewDefaultCertAuthClient function validates inputs properly
func TestNewDefaultCertAuthClient(t *testing.T) {
	certs := createCertificates(t)
	addr := "https://127.0.0.1:8200"
	tlsConfig := &api.TLSConfig{
		CAPath:     certs.caCert,
		ClientCert: certs.clientCert,
		ClientKey:  certs.clientKey,
	}
	type args struct {
		address   string
		tlsConfig *api.TLSConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"missing address", args{"", tlsConfig}, true},
		{"missing config", args{addr, nil}, true},
		{"missing client cert", args{addr, &api.TLSConfig{}}, true},
		{"ok", args{addr, tlsConfig}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDefaultCertAuthClient(tt.args.address, tt.args.tlsConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDefaultCertAuthClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("NewDefaultCertAuthClient() returned nil client")
					return
				}
				if got.Address() != tt.args.address {
					t.Errorf("NewDefaultCertAuthClient() returned client with address %s, expected %s", got.Address(), tt.args.address)
				}

				// The Login test will validate that our TLS config was passed in properly.
			}
		})
	}
}

// TestNewCertAuthClient tests the NewCertAuthClient function validates inputs properly
func TestNewCertAuthClient(t *testing.T) {
	certs := createCertificates(t)
	addr := "https://127.0.0.1:8200"
	tlsConfig := &api.TLSConfig{
		CAPath:     certs.caCert,
		ClientCert: certs.clientCert,
		ClientKey:  certs.clientKey,
	}
	defConfig := api.DefaultConfig()
	defConfig.Address = addr
	client, err := api.NewClient(defConfig)
	if err != nil {
		t.Fatalf("failed to create CertAuth: %v", err)
	}

	type args struct {
		client    *api.Client
		tlsConfig *api.TLSConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"missing address", args{nil, tlsConfig}, true},
		{"missing config", args{client, nil}, true},
		{"missing client cert", args{client, &api.TLSConfig{}}, true},
		{"ok", args{client, tlsConfig}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCertAuthClient(tt.args.client, tt.args.tlsConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDefaultCertAuthClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("NewDefaultCertAuthClient() returned nil client")
					return
				}
				if got.Address() != addr {
					t.Errorf("NewDefaultCertAuthClient() returned client with address %s, expected %s", got.Address(), addr)
				}

				// The Login test will validate that our TLS config was passed in properly.
			}
		})
	}
}

func runTestServer(t *testing.T, cert testCerts, fn http.HandlerFunc) net.Listener {
	certPool, err := rootcerts.LoadCACerts(&rootcerts.Config{
		CAPath: cert.caCert,
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
		Handler:   fn,
	}

	t.Cleanup(func() {
		_ = sv.Close()
		_ = ln.Close()
	})

	go func() {
		_ = sv.ServeTLS(ln, cert.serverCert, cert.serverKey)
	}()

	return ln
}

type testCerts struct {
	caCert           string
	caKey            string
	serverCert       string
	serverKey        string
	clientCert       string
	clientKey        string
	parsedClientCert *x509.Certificate
}

func createCertificates(t *testing.T) testCerts {
	tempDir := t.TempDir()

	caCertTemplate := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: "localhost",
		},
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
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
	caCertFile := filepath.Join(tempDir, "ca_cert.pem")
	err = os.WriteFile(caCertFile, pem.EncodeToMemory(caCertPEMBlock), 0o755)
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
	caKeyFile := filepath.Join(tempDir, "ca_key.pem")
	err = os.WriteFile(caKeyFile, pem.EncodeToMemory(caKeyPEMBlock), 0o755)
	if err != nil {
		t.Fatal(err)
	}

	createCertificate := func(prefix string) *x509.Certificate {
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

		parsedCert, err := x509.ParseCertificate(certBytes)
		if err != nil {
			t.Fatal(err)
		}
		return parsedCert
	}

	clientCert := createCertificate("client")
	createCertificate("server")

	filepath.Join(tempDir, "client_cert.pem")

	return testCerts{
		caCert:           caCertFile,
		caKey:            caKeyFile,
		serverCert:       filepath.Join(tempDir, "server_cert.pem"),
		serverKey:        filepath.Join(tempDir, "server_key.pem"),
		clientCert:       filepath.Join(tempDir, "client_cert.pem"),
		clientKey:        filepath.Join(tempDir, "client_key.pem"),
		parsedClientCert: clientCert,
	}
}
