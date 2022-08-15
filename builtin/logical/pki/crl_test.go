package pki

import (
	"context"
	"encoding/asn1"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func TestBackend_CRL_EnableDisableRoot(t *testing.T) {
	b, s := createBackendWithStorage(t)

	resp, err := CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"ttl":         "40h",
		"common_name": "myvault.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	caSerial := resp.Data["serial_number"].(string)

	crlEnableDisableTestForBackend(t, b, s, []string{caSerial})
}

func TestBackend_CRL_AllKeyTypeSigAlgos(t *testing.T) {
	type testCase struct {
		KeyType string
		KeyBits int
		SigAlgo string
	}

	testCases := []testCase{
		{"rsa", 2048, "SHA256WithRSA"},
		{"rsa", 2048, "SHA384WithRSA"},
		{"rsa", 2048, "SHA512WithRSA"},
		{"rsa", 2048, "SHA256WithRSAPSS"},
		{"rsa", 2048, "SHA384WithRSAPSS"},
		{"rsa", 2048, "SHA512WithRSAPSS"},
		{"ec", 256, "ECDSAWithSHA256"},
		{"ec", 384, "ECDSAWithSHA384"},
		{"ec", 521, "ECDSAWithSHA512"},
		{"ed25519", 0, "PureEd25519"},
	}

	for index, tc := range testCases {
		t.Logf("tv %v", index)
		b, s := createBackendWithStorage(t)

		resp, err := CBWrite(b, s, "root/generate/internal", map[string]interface{}{
			"ttl":         "40h",
			"common_name": "myvault.com",
			"key_type":    tc.KeyType,
			"key_bits":    tc.KeyBits,
		})
		if err != nil {
			t.Fatalf("tc %v: %v", index, err)
		}
		caSerial := resp.Data["serial_number"].(string)

		_, err = CBPatch(b, s, "issuer/default", map[string]interface{}{
			"revocation_signature_algorithm": tc.SigAlgo,
		})
		if err != nil {
			t.Fatalf("tc %v: %v", index, err)
		}

		crlEnableDisableTestForBackend(t, b, s, []string{caSerial})

		crl := getParsedCrlFromBackend(t, b, s, "crl")
		if strings.HasSuffix(tc.SigAlgo, "PSS") {
			algo := crl.SignatureAlgorithm
			pssOid := asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 10}
			if !algo.Algorithm.Equal(pssOid) {
				t.Fatalf("tc %v failed: expected sig-alg to be %v / got %v", index, pssOid, algo)
			}
		}
	}
}

func TestBackend_CRL_EnableDisableIntermediateWithRoot(t *testing.T) {
	crlEnableDisableIntermediateTestForBackend(t, true)
}

func TestBackend_CRL_EnableDisableIntermediateWithoutRoot(t *testing.T) {
	crlEnableDisableIntermediateTestForBackend(t, false)
}

func crlEnableDisableIntermediateTestForBackend(t *testing.T, withRoot bool) {
	b_root, s_root := createBackendWithStorage(t)

	resp, err := CBWrite(b_root, s_root, "root/generate/internal", map[string]interface{}{
		"ttl":         "40h",
		"common_name": "myvault.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	rootSerial := resp.Data["serial_number"].(string)

	b_int, s_int := createBackendWithStorage(t)

	resp, err = CBWrite(b_int, s_int, "intermediate/generate/internal", map[string]interface{}{
		"common_name": "intermediate myvault.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected intermediate CSR info")
	}
	intermediateData := resp.Data

	resp, err = CBWrite(b_root, s_root, "root/sign-intermediate", map[string]interface{}{
		"ttl": "30h",
		"csr": intermediateData["csr"],
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected signed intermediate info")
	}
	intermediateSignedData := resp.Data
	var certs string = intermediateSignedData["certificate"].(string)
	caSerial := intermediateSignedData["serial_number"].(string)
	caSerials := []string{caSerial}
	if withRoot {
		intermediateAndRootCert := intermediateSignedData["ca_chain"].([]string)
		certs = strings.Join(intermediateAndRootCert, "\n")
		caSerials = append(caSerials, rootSerial)
	}

	resp, err = CBWrite(b_int, s_int, "intermediate/set-signed", map[string]interface{}{
		"certificate": certs,
	})

	crlEnableDisableTestForBackend(t, b_int, s_int, caSerials)
}

func crlEnableDisableTestForBackend(t *testing.T, b *backend, s logical.Storage, caSerials []string) {
	var err error

	_, err = CBWrite(b, s, "roles/test", map[string]interface{}{
		"allow_bare_domains": true,
		"allow_subdomains":   true,
		"allowed_domains":    "foobar.com",
		"generate_lease":     true,
	})
	if err != nil {
		t.Fatal(err)
	}

	serials := make(map[int]string)
	for i := 0; i < 6; i++ {
		resp, err := CBWrite(b, s, "issue/test", map[string]interface{}{
			"common_name": "test.foobar.com",
		})
		if err != nil {
			t.Fatal(err)
		}
		serials[i] = resp.Data["serial_number"].(string)
	}

	test := func(numRevokedExpected int, expectedSerials ...string) {
		certList := getParsedCrlFromBackend(t, b, s, "crl").TBSCertList
		lenList := len(certList.RevokedCertificates)
		if lenList != numRevokedExpected {
			t.Fatalf("expected %d revoked certificates, found %d", numRevokedExpected, lenList)
		}

		for _, serialNum := range expectedSerials {
			requireSerialNumberInCRL(t, certList, serialNum)
		}
	}

	revoke := func(serialIndex int) {
		_, err = CBWrite(b, s, "revoke", map[string]interface{}{
			"serial_number": serials[serialIndex],
		})
		if err != nil {
			t.Fatal(err)
		}

		for _, caSerial := range caSerials {
			_, err = CBWrite(b, s, "revoke", map[string]interface{}{
				"serial_number": caSerial,
			})
			if err == nil {
				t.Fatal("expected error")
			}
		}
	}

	toggle := func(disabled bool) {
		_, err = CBWrite(b, s, "config/crl", map[string]interface{}{
			"disable": disabled,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	test(0)
	revoke(0)
	revoke(1)
	test(2, serials[0], serials[1])
	toggle(true)
	test(0)
	revoke(2)
	revoke(3)
	test(0)
	toggle(false)
	test(4, serials[0], serials[1], serials[2], serials[3])
	revoke(4)
	revoke(5)
	test(6)
	toggle(true)
	test(0)
	toggle(false)
	test(6)

	// The rotate command should reset the update time of the CRL.
	crlCreationTime1 := getParsedCrlFromBackend(t, b, s, "crl").TBSCertList.ThisUpdate
	time.Sleep(1 * time.Second)
	_, err = CBRead(b, s, "crl/rotate")
	require.NoError(t, err)

	crlCreationTime2 := getParsedCrlFromBackend(t, b, s, "crl").TBSCertList.ThisUpdate
	require.NotEqual(t, crlCreationTime1, crlCreationTime2)
}

func TestBackend_Secondary_CRL_Rebuilding(t *testing.T) {
	ctx := context.Background()
	b, s := createBackendWithStorage(t)
	sc := b.makeStorageContext(ctx, s)

	// Write out the issuer/key to storage without going through the api call as replication would.
	bundle := genCertBundle(t, b, s)
	issuer, _, err := sc.writeCaBundle(bundle, "", "")
	require.NoError(t, err)

	// Just to validate, before we call the invalidate function, make sure our CRL has not been generated
	// and we get a nil response
	resp := requestCrlFromBackend(t, s, b)
	require.Nil(t, resp.Data["http_raw_body"])

	// This should force any calls from now on to rebuild our CRL even a read
	b.invalidate(ctx, issuerPrefix+issuer.ID.String())

	// Perform the read operation again, we should have a valid CRL now...
	resp = requestCrlFromBackend(t, s, b)
	crl := parseCrlPemBytes(t, resp.Data["http_raw_body"].([]byte))
	require.Equal(t, 0, len(crl.RevokedCertificates))
}

func TestCrlRebuilder(t *testing.T) {
	ctx := context.Background()
	b, s := createBackendWithStorage(t)
	sc := b.makeStorageContext(ctx, s)

	// Write out the issuer/key to storage without going through the api call as replication would.
	bundle := genCertBundle(t, b, s)
	_, _, err := sc.writeCaBundle(bundle, "", "")
	require.NoError(t, err)

	req := &logical.Request{Storage: s}
	cb := crlBuilder{}

	// Force an initial build
	err = cb.rebuild(ctx, b, req, true)
	require.NoError(t, err, "Failed to rebuild CRL")

	resp := requestCrlFromBackend(t, s, b)
	crl1 := parseCrlPemBytes(t, resp.Data["http_raw_body"].([]byte))

	// We shouldn't rebuild within this call.
	err = cb.rebuildIfForced(ctx, b, req)
	require.NoError(t, err, "Failed to rebuild if forced CRL")
	resp = requestCrlFromBackend(t, s, b)
	crl2 := parseCrlPemBytes(t, resp.Data["http_raw_body"].([]byte))
	require.Equal(t, crl1.ThisUpdate, crl2.ThisUpdate, "According to the update field, we rebuilt the CRL")

	// Make sure we have ticked over to the next second
	for {
		diff := time.Now().Sub(crl1.ThisUpdate)
		if diff.Seconds() >= 1 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	// This should rebuild the CRL
	cb.requestRebuildIfActiveNode(b)
	err = cb.rebuildIfForced(ctx, b, req)
	require.NoError(t, err, "Failed to rebuild if forced CRL")
	resp = requestCrlFromBackend(t, s, b)
	crl3 := parseCrlPemBytes(t, resp.Data["http_raw_body"].([]byte))
	require.True(t, crl1.ThisUpdate.Before(crl3.ThisUpdate),
		"initial crl time: %#v not before next crl rebuild time: %#v", crl1.ThisUpdate, crl3.ThisUpdate)
}

func TestBYOC(t *testing.T) {
	b, s := createBackendWithStorage(t)

	// Create a root CA.
	resp, err := CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "root example.com",
		"issuer_name": "root",
		"key_type":    "ec",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Data["certificate"])
	oldRoot := resp.Data["certificate"].(string)

	// Create a role for issuance.
	_, err = CBWrite(b, s, "roles/local-testing", map[string]interface{}{
		"allow_any_name":    true,
		"enforce_hostnames": false,
		"key_type":          "ec",
		"ttl":               "75s",
		"no_store":          "true",
	})
	require.NoError(t, err)

	// Issue a leaf cert and ensure we can revoke it.
	resp, err = CBWrite(b, s, "issue/local-testing", map[string]interface{}{
		"common_name": "testing",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Data["certificate"])

	_, err = CBWrite(b, s, "revoke", map[string]interface{}{
		"certificate": resp.Data["certificate"],
	})
	require.NoError(t, err)

	// Issue a second leaf, but hold onto it for now.
	resp, err = CBWrite(b, s, "issue/local-testing", map[string]interface{}{
		"common_name": "testing2",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Data["certificate"])
	notStoredCert := resp.Data["certificate"].(string)

	// Update the role to make things stored and issue another cert.
	_, err = CBWrite(b, s, "roles/stored-testing", map[string]interface{}{
		"allow_any_name":    true,
		"enforce_hostnames": false,
		"key_type":          "ec",
		"ttl":               "75s",
		"no_store":          "false",
	})
	require.NoError(t, err)

	// Issue a leaf cert and ensure we can revoke it.
	resp, err = CBWrite(b, s, "issue/stored-testing", map[string]interface{}{
		"common_name": "testing",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Data["certificate"])
	storedCert := resp.Data["certificate"].(string)

	// Delete the root and regenerate a new one.
	_, err = CBDelete(b, s, "issuer/default")
	require.NoError(t, err)

	resp, err = CBList(b, s, "issuers")
	require.NoError(t, err)
	require.Equal(t, len(resp.Data), 0)

	_, err = CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "root2 example.com",
		"issuer_name": "root2",
		"key_type":    "ec",
	})
	require.NoError(t, err)

	// Issue a new leaf and revoke that one.
	resp, err = CBWrite(b, s, "issue/local-testing", map[string]interface{}{
		"common_name": "testing3",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Data["certificate"])

	_, err = CBWrite(b, s, "revoke", map[string]interface{}{
		"certificate": resp.Data["certificate"],
	})
	require.NoError(t, err)

	// Now attempt to revoke the earlier leaves. The first should fail since
	// we deleted its issuer, but the stored one should succeed.
	_, err = CBWrite(b, s, "revoke", map[string]interface{}{
		"certificate": notStoredCert,
	})
	require.Error(t, err)

	_, err = CBWrite(b, s, "revoke", map[string]interface{}{
		"certificate": storedCert,
	})
	require.NoError(t, err)

	// Import the old root again and revoke the no stored leaf should work.
	_, err = CBWrite(b, s, "issuers/import/bundle", map[string]interface{}{
		"pem_bundle": oldRoot,
	})
	require.NoError(t, err)

	_, err = CBWrite(b, s, "revoke", map[string]interface{}{
		"certificate": notStoredCert,
	})
	require.NoError(t, err)
}

func requestCrlFromBackend(t *testing.T, s logical.Storage, b *backend) *logical.Response {
	crlReq := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "crl/pem",
		Storage:   s,
	}
	resp, err := b.HandleRequest(context.Background(), crlReq)
	require.NoError(t, err, "crl req failed with an error")
	require.NotNil(t, resp, "crl response was nil with no error")
	require.False(t, resp.IsError(), "crl error response: %v", resp)
	return resp
}
