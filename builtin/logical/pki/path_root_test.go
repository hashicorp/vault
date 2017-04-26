package pki

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func TestCAGenerateRoot(t *testing.T) {
	storage := &logical.InmemStorage{}
	config := logical.TestBackendConfig()
	config.StorageView = storage

	b := Backend()
	_, err := b.Setup(config)
	if err != nil {
		t.Fatal(err)
	}

	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/generate/internal",
		Storage:   storage,
	}

	fd := &framework.FieldData{
		Raw: map[string]interface{}{
			"exported":    "internal",
			"common_name": "test.example.com",
		},
		Schema: pathGenerateRoot(b).Fields,
	}

	resp, err := b.pathCAGenerateRoot(req, fd)
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if resp.Error() != nil {
		t.Fatalf("logical.Response error: %s", resp.Error())
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

func TestCASignIntermediate(t *testing.T) {
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
		Path:      "root/sign-intermediate",
		Storage:   storage,
	}

	csrPEM, err := ioutil.ReadFile("test-fixtures/root/csr.pem")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	fd := &framework.FieldData{
		Raw: map[string]interface{}{
			"common_name": "test.example.com",
			"csr":         string(csrPEM),
		},
		Schema: pathSignIntermediate(b).Fields,
	}

	resp, err := b.pathCASignIntermediate(req, fd)
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
