package pki

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/logical"
)

func TestPki_SetSignedIntermediate(t *testing.T) {
	storage := &logical.InmemStorage{}
	config := logical.TestBackendConfig()
	config.StorageView = storage

	b := Backend()
	_, err := b.Setup(config)
	if err != nil {
		t.Fatal(err)
	}

	// Put cert bundle in inmem storage
	privateCertPEM, err := ioutil.ReadFile("test-fixtures/cakey.pem")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	cb := &certutil.CertBundle{}
	cb.PrivateKey = string(privateCertPEM)
	cb.PrivateKeyType = certutil.RSAPrivateKey

	bundleEntry, err := logical.StorageEntryJSON("config/ca_bundle", cb)
	if err != nil {
		t.Fatal(err)
	}
	err = storage.Put(bundleEntry)
	if err != nil {
		t.Fatal(err)
	}

	certValue, err := ioutil.ReadFile("test-fixtures/cacert.pem")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	req := logical.TestRequest(t, logical.UpdateOperation, "intermediate/set-signed")
	req.Data["certificate"] = certValue
	req.Storage = storage

	resp, err := b.HandleRequest(req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	// Verify that value was written to storage
	serial := "5e:21:03:b9:e7:30:b9:af:7e:8f:55:c7:2e:77:28:9f:14:3f:24:17"
	storageKey := "certs/" + strings.ToLower(strings.Replace(serial, ":", "-", -1))
	entry, err := storage.Get(storageKey)
	if err != nil {
		t.Fatal(err)
	}
	if entry == nil {
		t.Fatal("update operation unsucessful, data not written to storage")
	}
}
