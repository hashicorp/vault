// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pki

import (
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
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
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func TestBackend_CA_Steps(t *testing.T) {
	t.Parallel()
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
	var rsaCAKey, rsaCACert, ecCAKey, ecCACert, edCAKey, edCACert string
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

		_, edk, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			panic(err)
		}
		marshaledKey, err = x509.MarshalPKCS8PrivateKey(edk)
		if err != nil {
			panic(err)
		}
		keyPEMBlock = &pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: marshaledKey,
		}
		edCAKey = strings.TrimSpace(string(pem.EncodeToMemory(keyPEMBlock)))
		if err != nil {
			panic(err)
		}
		_, err = certutil.GetSubjKeyID(edk)
		if err != nil {
			panic(err)
		}
		caBytes, err = x509.CreateCertificate(rand.Reader, caCertTemplate, caCertTemplate, edk.Public(), edk)
		if err != nil {
			panic(err)
		}
		caCertPEMBlock = &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: caBytes,
		}
		edCACert = strings.TrimSpace(string(pem.EncodeToMemory(caCertPEMBlock)))
	}

	// Setup backends
	var rsaRoot, rsaInt, ecRoot, ecInt, edRoot, edInt *backend
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

		if err := client.Sys().Mount("ed25519root", &api.MountInput{
			Type: "pki",
			Config: api.MountConfigInput{
				DefaultLeaseTTL: "16h",
				MaxLeaseTTL:     "60h",
			},
		}); err != nil {
			t.Fatal(err)
		}
		edRoot = b

		if err := client.Sys().Mount("ed25519int", &api.MountInput{
			Type: "pki",
			Config: api.MountConfigInput{
				DefaultLeaseTTL: "16h",
				MaxLeaseTTL:     "60h",
			},
		}); err != nil {
			t.Fatal(err)
		}
		edInt = b
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
		t.Run("ed25519", func(t *testing.T) {
			t.Parallel()
			subClient, err := client.Clone()
			if err != nil {
				t.Fatal(err)
			}
			subClient.SetToken(client.Token())
			runSteps(t, edRoot, edInt, subClient, "ed25519root/", "ed25519int/", edCACert, edCAKey)
		})
	})
}

func runSteps(t *testing.T, rootB, intB *backend, client *api.Client, rootName, intName, caCert, caKey string) {
	//  Load CA cert/key in and ensure we can fetch it back in various formats,
	//  unauthenticated
	{
		// Attempt import but only provide one the cert; this should work.
		{
			_, err := client.Logical().Write(rootName+"config/ca", map[string]interface{}{
				"pem_bundle": caCert,
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}

		// Same but with only the key
		{
			_, err := client.Logical().Write(rootName+"config/ca", map[string]interface{}{
				"pem_bundle": caKey,
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}

		// Import entire CA bundle; this should work as well
		{
			_, err := client.Logical().Write(rootName+"config/ca", map[string]interface{}{
				"pem_bundle": strings.Join([]string{caKey, caCert}, "\n"),
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}

		prevToken := client.Token()
		client.SetToken("")

		// cert/ca and issuer/default/json path
		for _, path := range []string{"cert/ca", "issuer/default/json"} {
			resp, err := client.Logical().Read(rootName + path)
			if err != nil {
				t.Fatal(err)
			}
			if resp == nil {
				t.Fatal("nil response")
			}
			expected := caCert
			if path == "issuer/default/json" {
				// Preserves the new line.
				expected += "\n"
				_, present := resp.Data["issuer_id"]
				if !present {
					t.Fatalf("expected issuer/default/json to include issuer_id")
				}
				_, present = resp.Data["issuer_name"]
				if !present {
					t.Fatalf("expected issuer/default/json to include issuer_name")
				}
			}
			if diff := deep.Equal(resp.Data["certificate"].(string), expected); diff != nil {
				t.Fatal(diff)
			}
		}

		// ca/pem and issuer/default/pem path (raw string)
		for _, path := range []string{"ca/pem", "issuer/default/pem"} {
			req := &logical.Request{
				Path:      path,
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
			expected := []byte(caCert)
			if path == "issuer/default/pem" {
				// Preserves the new line.
				expected = []byte(caCert + "\n")
			}
			if diff := deep.Equal(resp.Data["http_raw_body"].([]byte), expected); diff != nil {
				t.Fatal(diff)
			}
			if resp.Data["http_content_type"].(string) != "application/pem-certificate-chain" {
				t.Fatal("wrong content type")
			}
		}

		// ca and issuer/default/der (raw DER bytes)
		for _, path := range []string{"ca", "issuer/default/der"} {
			req := &logical.Request{
				Path:      path,
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
			if resp == nil {
				t.Fatal("nil response")
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
		for path, derPemOrJSON := range map[string]int{
			"crl":                    0,
			"issuer/default/crl/der": 0,
			"crl/pem":                1,
			"issuer/default/crl/pem": 1,
			"cert/crl":               2,
			"issuer/default/crl":     3,
		} {
			req := &logical.Request{
				Path:      path,
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

			var crlBytes []byte
			if derPemOrJSON == 2 {
				// Old endpoint
				crlBytes = []byte(resp.Data["certificate"].(string))
			} else if derPemOrJSON == 3 {
				// New endpoint
				crlBytes = []byte(resp.Data["crl"].(string))
			} else {
				// DER or PEM
				crlBytes = resp.Data["http_raw_body"].([]byte)
			}

			if derPemOrJSON >= 1 {
				// Do for both PEM and JSON endpoints
				pemBlock, _ := pem.Decode(crlBytes)
				crlBytes = pemBlock.Bytes
			}

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

	verifyTidyStatus := func(expectedCertStoreDeleteCount int, expectedRevokedCertDeletedCount int) {
		tidyStatus, err := client.Logical().Read(rootName + "tidy-status")
		if err != nil {
			t.Fatal(err)
		}

		if tidyStatus.Data["state"] != "Finished" {
			t.Fatalf("Expected tidy operation to be finished, but tidy-status reports its state is %v", tidyStatus.Data)
		}

		var count int64
		if count, err = tidyStatus.Data["cert_store_deleted_count"].(json.Number).Int64(); err != nil {
			t.Fatal(err)
		}
		if int64(expectedCertStoreDeleteCount) != count {
			t.Fatalf("Expected %d for cert_store_deleted_count, but got %d", expectedCertStoreDeleteCount, count)
		}

		if count, err = tidyStatus.Data["revoked_cert_deleted_count"].(json.Number).Int64(); err != nil {
			t.Fatal(err)
		}
		if int64(expectedRevokedCertDeletedCount) != count {
			t.Fatalf("Expected %d for revoked_cert_deleted_count, but got %d", expectedRevokedCertDeletedCount, count)
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

			verifyTidyStatus(0, 0)
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

			verifyTidyStatus(0, 0)
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

			verifyTidyStatus(1, 1)
		}
	}
}
