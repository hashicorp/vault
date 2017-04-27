package pki

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/logical"
)

func TestPki_CAGenerateRoot(t *testing.T) {
	storage := &logical.InmemStorage{}
	config := logical.TestBackendConfig()
	config.StorageView = storage

	b := Backend()
	_, err := b.Setup(config)
	if err != nil {
		t.Fatal(err)
	}

	req := logical.TestRequest(t, logical.UpdateOperation, "root/generate/internal")
	req.Storage = storage
	req.Data["common_name"] = "test.example.com"

	// resp, err := b.pathCAGenerateRoot(req, fd)
	resp, err := b.HandleRequest(req)
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

func TestPki_CASignIntermediate(t *testing.T) {
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

	csrPEM, err := ioutil.ReadFile("test-fixtures/root/csr.pem")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	req := logical.TestRequest(t, logical.UpdateOperation, "root/sign-intermediate")
	req.Storage = storage
	req.Data["csr"] = string(csrPEM)
	req.Data["common_name"] = "test.example.com"

	// resp, err := b.pathCASignIntermediate(req, fd)
	resp, err := b.HandleRequest(req)
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
