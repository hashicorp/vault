package pki

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"math/big"
	mathrand "math/rand"
	"strings"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/certutil"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func TestBackend_CA_Steps(t *testing.T) {
	var b *backend

	factory := func(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
		be, err := Factory(ctx, conf)
		if err == nil {
			b = be.(*backend)
		}
		return be, err
	}

	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"pki": factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client

	// Set RSA/EC CA certificates
	var rsaCAKey, rsaCACert, ecCAKey, ecCACert string
	{
		cak, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			panic(err)
		}
		marshaledKey, err := x509.MarshalECPrivateKey(cak)
		if err != nil {
			panic(err)
		}
		keyPEMBlock := &pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: marshaledKey,
		}
		ecCAKey = strings.TrimSpace(string(pem.EncodeToMemory(keyPEMBlock)))
		if err != nil {
			panic(err)
		}
		subjKeyID, err := certutil.GetSubjKeyID(cak)
		if err != nil {
			panic(err)
		}
		caCertTemplate := &x509.Certificate{
			Subject: pkix.Name{
				CommonName: "root.localhost",
			},
			SubjectKeyId:          subjKeyID,
			DNSNames:              []string{"root.localhost"},
			KeyUsage:              x509.KeyUsage(x509.KeyUsageCertSign | x509.KeyUsageCRLSign),
			SerialNumber:          big.NewInt(mathrand.Int63()),
			NotAfter:              time.Now().Add(262980 * time.Hour),
			BasicConstraintsValid: true,
			IsCA:                  true,
		}
		caBytes, err := x509.CreateCertificate(rand.Reader, caCertTemplate, caCertTemplate, cak.Public(), cak)
		if err != nil {
			panic(err)
		}
		caCertPEMBlock := &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: caBytes,
		}
		ecCACert = strings.TrimSpace(string(pem.EncodeToMemory(caCertPEMBlock)))

		rak, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			panic(err)
		}
		marshaledKey = x509.MarshalPKCS1PrivateKey(rak)
		keyPEMBlock = &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: marshaledKey,
		}
		rsaCAKey = strings.TrimSpace(string(pem.EncodeToMemory(keyPEMBlock)))
		if err != nil {
			panic(err)
		}
		_, err = certutil.GetSubjKeyID(rak)
		if err != nil {
			panic(err)
		}
		caBytes, err = x509.CreateCertificate(rand.Reader, caCertTemplate, caCertTemplate, rak.Public(), rak)
		if err != nil {
			panic(err)
		}
		caCertPEMBlock = &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: caBytes,
		}
		rsaCACert = strings.TrimSpace(string(pem.EncodeToMemory(caCertPEMBlock)))
	}

	// Setup backends
	var rsaRoot, rsaInt, ecRoot, ecInt *backend
	{
		if err := client.Sys().Mount("rsaroot", &api.MountInput{
			Type: "pki",
			Config: api.MountConfigInput{
				DefaultLeaseTTL: "16h",
				MaxLeaseTTL:     "60h",
			},
		}); err != nil {
			t.Fatal(err)
		}
		rsaRoot = b

		if err := client.Sys().Mount("rsaint", &api.MountInput{
			Type: "pki",
			Config: api.MountConfigInput{
				DefaultLeaseTTL: "16h",
				MaxLeaseTTL:     "60h",
			},
		}); err != nil {
			t.Fatal(err)
		}
		rsaInt = b

		if err := client.Sys().Mount("ecroot", &api.MountInput{
			Type: "pki",
			Config: api.MountConfigInput{
				DefaultLeaseTTL: "16h",
				MaxLeaseTTL:     "60h",
			},
		}); err != nil {
			t.Fatal(err)
		}
		ecRoot = b

		if err := client.Sys().Mount("ecint", &api.MountInput{
			Type: "pki",
			Config: api.MountConfigInput{
				DefaultLeaseTTL: "16h",
				MaxLeaseTTL:     "60h",
			},
		}); err != nil {
			t.Fatal(err)
		}
		ecInt = b
	}

	t.Run("teststeps", func(t *testing.T) {
		t.Run("rsa", func(t *testing.T) {
			t.Parallel()
			subClient, err := client.Clone()
			if err != nil {
				t.Fatal(err)
			}
			subClient.SetToken(client.Token())
			runSteps(t, rsaRoot, rsaInt, subClient, "rsaroot/", "rsaint/", rsaCACert, rsaCAKey)
		})
		t.Run("ec", func(t *testing.T) {
			t.Parallel()
			subClient, err := client.Clone()
			if err != nil {
				t.Fatal(err)
			}
			subClient.SetToken(client.Token())
			runSteps(t, ecRoot, ecInt, subClient, "ecroot/", "ecint/", ecCACert, ecCAKey)
		})
	})
}

func runSteps(t *testing.T, rootB, intB *backend, client *api.Client, rootName, intName, caCert, caKey string) {
	//  Load CA cert/key in and ensure we can fetch it back in various formats,
	//  unauthenticated
	{
		// Attempt import but only provide one the cert
		{
			_, err := client.Logical().Write(rootName+"config/ca", map[string]interface{}{
				"pem_bundle": caCert,
			})
			if err == nil {
				t.Fatal("expected error")
			}
		}

		// Same but with only the key
		{
			_, err := client.Logical().Write(rootName+"config/ca", map[string]interface{}{
				"pem_bundle": caKey,
			})
			if err == nil {
				t.Fatal("expected error")
			}
		}

		// Import CA bundle
		{
			_, err := client.Logical().Write(rootName+"config/ca", map[string]interface{}{
				"pem_bundle": strings.Join([]string{caKey, caCert}, "\n"),
			})
			if err != nil {
				t.Fatal(err)
			}
		}

		prevToken := client.Token()
		client.SetToken("")

		// cert/ca path
		{
			resp, err := client.Logical().Read(rootName + "cert/ca")
			if err != nil {
				t.Fatal(err)
			}
			if resp == nil {
				t.Fatal("nil response")
			}
			if diff := deep.Equal(resp.Data["certificate"].(string), caCert); diff != nil {
				t.Fatal(diff)
			}
		}
		// ca/pem path (raw string)
		{
			req := &logical.Request{
				Path:      "ca/pem",
				Operation: logical.ReadOperation,
				Storage:   rootB.storage,
			}
			resp, err := rootB.HandleRequest(context.Background(), req)
			if err != nil {
				t.Fatal(err)
			}
			if resp == nil {
				t.Fatal("nil response")
			}
			if diff := deep.Equal(resp.Data["http_raw_body"].([]byte), []byte(caCert)); diff != nil {
				t.Fatal(diff)
			}
			if resp.Data["http_content_type"].(string) != "application/pkix-cert" {
				t.Fatal("wrong content type")
			}
		}

		// ca (raw DER bytes)
		{
			req := &logical.Request{
				Path:      "ca",
				Operation: logical.ReadOperation,
				Storage:   rootB.storage,
			}
			resp, err := rootB.HandleRequest(context.Background(), req)
			if err != nil {
				t.Fatal(err)
			}
			if resp == nil {
				t.Fatal("nil response")
			}
			rawBytes := resp.Data["http_raw_body"].([]byte)
			pemBytes := strings.TrimSpace(string(pem.EncodeToMemory(&pem.Block{
				Type:  "CERTIFICATE",
				Bytes: rawBytes,
			})))
			if diff := deep.Equal(pemBytes, caCert); diff != nil {
				t.Fatal(diff)
			}
			if resp.Data["http_content_type"].(string) != "application/pkix-cert" {
				t.Fatal("wrong content type")
			}
		}

		client.SetToken(prevToken)
	}

	// Configure an expiry on the CRL and verify what comes back
	{
		// Set CRL config
		{
			_, err := client.Logical().Write(rootName+"config/crl", map[string]interface{}{
				"expiry": "16h",
			})
			if err != nil {
				t.Fatal(err)
			}
		}

		// Verify it
		{
			resp, err := client.Logical().Read(rootName + "config/crl")
			if err != nil {
				t.Fatal(err)
			}
			if resp == nil {
				t.Fatal("nil response")
			}
			if resp.Data["expiry"].(string) != "16h" {
				t.Fatal("expected a 16 hour expiry")
			}
		}
	}

	// Test generating a root, an intermediate, signing it, setting signed, and
	// revoking it

	// We'll need this later
	var intSerialNumber string
	{
		// First, delete the existing CA info
		{
			_, err := client.Logical().Delete(rootName + "root")
			if err != nil {
				t.Fatal(err)
			}
		}

		var rootPEM, rootKey, rootPEMBundle string
		// Test exported root generation
		{
			resp, err := client.Logical().Write(rootName+"root/generate/exported", map[string]interface{}{
				"common_name": "Root Cert",
				"ttl":         "180h",
			})
			if err != nil {
				t.Fatal(err)
			}
			if resp == nil {
				t.Fatal("nil response")
			}
			rootPEM = resp.Data["certificate"].(string)
			rootKey = resp.Data["private_key"].(string)
			rootPEMBundle = strings.Join([]string{rootPEM, rootKey}, "\n")
			// This is really here to keep the use checker happy
			if rootPEMBundle == "" {
				t.Fatal("bad root pem bundle")
			}
		}

		var intPEM, intCSR, intKey string
		// Test exported intermediate CSR generation
		{
			resp, err := client.Logical().Write(intName+"intermediate/generate/exported", map[string]interface{}{
				"common_name": "intermediate.cert.com",
				"ttl":         "180h",
			})
			if err != nil {
				t.Fatal(err)
			}
			if resp == nil {
				t.Fatal("nil response")
			}
			intCSR = resp.Data["csr"].(string)
			intKey = resp.Data["private_key"].(string)
			// This is really here to keep the use checker happy
			if intCSR == "" || intKey == "" {
				t.Fatal("int csr or key empty")
			}
		}

		// Test signing
		{
			resp, err := client.Logical().Write(rootName+"root/sign-intermediate", map[string]interface{}{
				"common_name": "intermediate.cert.com",
				"ttl":         "10s",
				"csr":         intCSR,
			})
			if err != nil {
				t.Fatal(err)
			}
			if resp == nil {
				t.Fatal("nil response")
			}
			intPEM = resp.Data["certificate"].(string)
			intSerialNumber = resp.Data["serial_number"].(string)
		}

		// Test setting signed
		{
			resp, err := client.Logical().Write(intName+"intermediate/set-signed", map[string]interface{}{
				"certificate": intPEM,
			})
			if err != nil {
				t.Fatal(err)
			}
			if resp != nil {
				t.Fatal("expected nil response")
			}
		}

		// Verify we can find it via the root
		{
			resp, err := client.Logical().Read(rootName + "cert/" + intSerialNumber)
			if err != nil {
				t.Fatal(err)
			}
			if resp == nil {
				t.Fatal("nil response")
			}
			if resp.Data["revocation_time"].(json.Number).String() != "0" {
				t.Fatal("expected a zero revocation time")
			}
		}

		// Revoke the intermediate
		{
			resp, err := client.Logical().Write(rootName+"revoke", map[string]interface{}{
				"serial_number": intSerialNumber,
			})
			if err != nil {
				t.Fatal(err)
			}
			if resp == nil {
				t.Fatal("nil response")
			}
		}
	}

	verifyRevocation := func(t *testing.T, serial string, shouldFind bool) {
		t.Helper()
		// Verify it is now revoked
		{
			resp, err := client.Logical().Read(rootName + "cert/" + intSerialNumber)
			if err != nil {
				t.Fatal(err)
			}
			switch shouldFind {
			case true:
				if resp == nil {
					t.Fatal("nil response")
				}
				if resp.Data["revocation_time"].(json.Number).String() == "0" {
					t.Fatal("expected a non-zero revocation time")
				}
			default:
				if resp != nil {
					t.Fatalf("expected nil response, got %#v", *resp)
				}
			}
		}

		// Fetch the CRL and make sure it shows up
		{
			req := &logical.Request{
				Path:      "crl",
				Operation: logical.ReadOperation,
				Storage:   rootB.storage,
			}
			resp, err := rootB.HandleRequest(context.Background(), req)
			if err != nil {
				t.Fatal(err)
			}
			if resp == nil {
				t.Fatal("nil response")
			}
			crlBytes := resp.Data["http_raw_body"].([]byte)
			certList, err := x509.ParseCRL(crlBytes)
			if err != nil {
				t.Fatal(err)
			}
			switch shouldFind {
			case true:
				revokedList := certList.TBSCertList.RevokedCertificates
				if len(revokedList) != 1 {
					t.Fatalf("bad length of revoked list: %d", len(revokedList))
				}
				revokedString := certutil.GetHexFormatted(revokedList[0].SerialNumber.Bytes(), ":")
				if revokedString != intSerialNumber {
					t.Fatalf("bad revoked serial: %s", revokedString)
				}
			default:
				revokedList := certList.TBSCertList.RevokedCertificates
				if len(revokedList) != 0 {
					t.Fatalf("bad length of revoked list: %d", len(revokedList))
				}
			}
		}
	}

	// Validate current state of revoked certificates
	verifyRevocation(t, intSerialNumber, true)

	// Give time for the safety buffer to pass before tidying
	time.Sleep(10 * time.Second)

	// Test tidying
	{
		// Run with a high safety buffer, nothing should happen
		{
			resp, err := client.Logical().Write(rootName+"tidy", map[string]interface{}{
				"safety_buffer":      "3h",
				"tidy_cert_store":    true,
				"tidy_revoked_certs": true,
			})
			if err != nil {
				t.Fatal(err)
			}
			if resp == nil {
				t.Fatal("expected warnings")
			}

			// Wait a few seconds as it runs in a goroutine
			time.Sleep(5 * time.Second)

			// Check to make sure we still find the cert and see it on the CRL
			verifyRevocation(t, intSerialNumber, true)
		}

		// Run with both values set false, nothing should happen
		{
			resp, err := client.Logical().Write(rootName+"tidy", map[string]interface{}{
				"safety_buffer":      "1s",
				"tidy_cert_store":    false,
				"tidy_revoked_certs": false,
			})
			if err != nil {
				t.Fatal(err)
			}
			if resp == nil {
				t.Fatal("expected warnings")
			}

			// Wait a few seconds as it runs in a goroutine
			time.Sleep(5 * time.Second)

			// Check to make sure we still find the cert and see it on the CRL
			verifyRevocation(t, intSerialNumber, true)
		}

		// Run with a short safety buffer and both set to true, both should be cleared
		{
			resp, err := client.Logical().Write(rootName+"tidy", map[string]interface{}{
				"safety_buffer":      "1s",
				"tidy_cert_store":    true,
				"tidy_revoked_certs": true,
			})
			if err != nil {
				t.Fatal(err)
			}
			if resp == nil {
				t.Fatal("expected warnings")
			}

			// Wait a few seconds as it runs in a goroutine
			time.Sleep(5 * time.Second)

			// Check to make sure we still find the cert and see it on the CRL
			verifyRevocation(t, intSerialNumber, false)
		}
	}
}
