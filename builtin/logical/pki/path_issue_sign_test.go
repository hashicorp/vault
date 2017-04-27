package pki

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func TestPki_IssueSignCert(t *testing.T) {
	storage := &logical.InmemStorage{}
	config := logical.TestBackendConfig()
	config.StorageView = storage

	b := Backend()
	_, err := b.Setup(config)
	if err != nil {
		t.Fatal(err)
	}

	// Place CA cert in storage
	rootCAKeyPEM, err := ioutil.ReadFile("test-fixtures/root/rootcakey.pem")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	rootCACertPEM, err := ioutil.ReadFile("test-fixtures/root/rootcacert.pem")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	cb := &certutil.CertBundle{}
	cb.PrivateKey = string(rootCAKeyPEM)
	cb.PrivateKeyType = certutil.RSAPrivateKey
	cb.Certificate = string(rootCACertPEM)

	bundleEntry, err := logical.StorageEntryJSON("config/ca_bundle", cb)
	if err != nil {
		t.Fatal(err)
	}
	err = storage.Put(bundleEntry)
	if err != nil {
		t.Fatal(err)
	}

	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   storage,
	}

	ttl := b.System().DefaultLeaseTTL()
	role := &roleEntry{
		TTL:              ttl.String(),
		AllowLocalhost:   true,
		AllowAnyName:     true,
		AllowIPSANs:      true,
		EnforceHostnames: false,
		GenerateLease:    new(bool),
		KeyType:          "rsa",
		KeyBits:          2048,
		UseCSRCommonName: false,
		UseCSRSANs:       false,
	}
	*role.GenerateLease = false

	fd := &framework.FieldData{
		Raw: map[string]interface{}{
			"format":      "pem",
			"common_name": "test.example.com",
		},
		Schema: map[string]*framework.FieldSchema{
			"format":               &framework.FieldSchema{Type: framework.TypeString},
			"common_name":          &framework.FieldSchema{Type: framework.TypeString},
			"exclude_cn_from_sans": &framework.FieldSchema{Type: framework.TypeBool},
		},
	}

	resp, err := b.pathIssueSignCert(req, fd, role, false, false)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	// Verify that value was written to storage
	serial := resp.Data["serial_number"].(string)
	storageKey := "certs/" + strings.ToLower(strings.Replace(serial, ":", "-", -1))
	entry, err := storage.Get(storageKey)
	if err != nil {
		t.Fatal(err)
	}
	if entry == nil {
		t.Fatal("update operation unsucessful, data not written to storage")
	}
}
