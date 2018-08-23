package pki

import (
	"crypto/x509"
	"testing"

	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
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

	var serials = make(map[int]string)
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
		resp, err := client.Logical().Read("pki/cert/crl")
		if err != nil {
			t.Fatal(err)
		}
		crlPem := resp.Data["certificate"].(string)
		certList, err := x509.ParseCRL([]byte(crlPem))
		if err != nil {
			t.Fatal(err)
		}
		lenList := len(certList.TBSCertList.RevokedCertificates)
		if lenList != num {
			t.Fatalf("expected %d, found %d", num, lenList)
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
}
