package pki

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

func TestBackend_CRL_EnableDisable(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"pki": Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client
	var err error
	err = client.Sys().Mount("pki", &api.MountInput{
		Type: "pki",
		Config: api.MountConfigInput{
			DefaultLeaseTTL: "16h",
			MaxLeaseTTL:     "60h",
		},
	})

	resp, err := client.Logical().Write("pki/root/generate/internal", map[string]interface{}{
		"ttl":         "40h",
		"common_name": "myvault.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	caSerial := resp.Data["serial_number"]

	_, err = client.Logical().Write("pki/roles/test", map[string]interface{}{
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
		resp, err := client.Logical().Write("pki/issue/test", map[string]interface{}{
			"common_name": "test.foobar.com",
		})
		if err != nil {
			t.Fatal(err)
		}
		serials[i] = resp.Data["serial_number"].(string)
	}

	test := func(num int) {
		certList := getCrlCertificateList(t, client)
		lenList := len(certList.RevokedCertificates)
		if lenList != num {
			t.Fatalf("expected %d revoked certificates, found %d", num, lenList)
		}
	}

	revoke := func(num int) {
		resp, err = client.Logical().Write("pki/revoke", map[string]interface{}{
			"serial_number": serials[num],
		})
		if err != nil {
			t.Fatal(err)
		}

		resp, err = client.Logical().Write("pki/revoke", map[string]interface{}{
			"serial_number": caSerial,
		})
		if err == nil {
			t.Fatal("expected error")
		}
	}

	toggle := func(disabled bool) {
		_, err = client.Logical().Write("pki/config/crl", map[string]interface{}{
			"disable": disabled,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	test(0)
	revoke(0)
	revoke(1)
	test(2)
	toggle(true)
	test(0)
	revoke(2)
	revoke(3)
	test(0)
	toggle(false)
	test(4)
	revoke(4)
	revoke(5)
	test(6)
	toggle(true)
	test(0)
	toggle(false)
	test(6)

	// The rotate command should reset the update time of the CRL.
	crlCreationTime1 := getCrlCertificateList(t, client).ThisUpdate
	time.Sleep(1 * time.Second)
	_, err = client.Logical().Read("pki/crl/rotate")
	require.NoError(t, err)

	crlCreationTime2 := getCrlCertificateList(t, client).ThisUpdate
	require.NotEqual(t, crlCreationTime1, crlCreationTime2)
}

func TestBackend_Secondary_CRL_Rebuilding(t *testing.T) {
	ctx := context.Background()
	b, s := createBackendWithStorage(t)

	// Write out the issuer/key to storage without going through the api call as replication would.
	bundle := genCertBundle(t, b)
	issuer, _, err := writeCaBundle(ctx, s, bundle, "", "")
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

func getCrlCertificateList(t *testing.T, client *api.Client) pkix.TBSCertificateList {
	resp, err := client.Logical().ReadWithContext(context.Background(), "pki/cert/crl")
	require.NoError(t, err, "crl req failed with an error")
	require.NotNil(t, resp, "crl response was nil with no error")

	crlPem := resp.Data["certificate"].(string)
	return parseCrlPemString(t, crlPem)
}

func parseCrlPemString(t *testing.T, crlPem string) pkix.TBSCertificateList {
	return parseCrlPemBytes(t, []byte(crlPem))
}

func parseCrlPemBytes(t *testing.T, crlPem []byte) pkix.TBSCertificateList {
	certList, err := x509.ParseCRL(crlPem)
	require.NoError(t, err)
	return certList.TBSCertList
}
