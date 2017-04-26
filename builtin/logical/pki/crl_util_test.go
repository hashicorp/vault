package pki

import (
	"encoding/pem"
	"io/ioutil"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/logical"
)

func TestRevokeCert(t *testing.T) {
	storage := &logical.InmemStorage{}
	config := logical.TestBackendConfig()
	config.StorageView = storage

	b := Backend()
	_, err := b.Setup(config)
	if err != nil {
		t.Fatal(err)
	}

	certValue, err := ioutil.ReadFile("test-fixtures/keys/cert.pem")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	block, _ := pem.Decode(certValue)
	if block == nil {
		t.Fatal("failed to decode PEM cert into DER")
	}

	var revInfo revocationInfo
	currTime := time.Now()
	revInfo.CertificateBytes = block.Bytes
	revInfo.RevocationTime = currTime.Unix()
	revInfo.RevocationTimeUTC = currTime.UTC()
	encodedCertValue, err := jsonutil.EncodeJSON(revInfo)
	if err != nil {
		t.Fatalf("error encoding pseudo cert value: %s", err)
	}

	cases := map[string]struct {
		Req          *logical.Request
		StorageKey   string
		StorageValue []byte
	}{
		"cert, valid colon": {
			&logical.Request{
				Operation: logical.UpdateOperation,
				Path:      "certs/7f:e8:e1:29:31:41:9e:a4:ac:df:82:08:d1:64:b5:2f:84:2c:6d:b0",
				Storage:   storage,
			},
			"7f:e8:e1:29:31:41:9e:a4:ac:df:82:08:d1:64:b5:2f:84:2c:6d:b0",
			certValue,
		},
		"cert, revoked colon": {
			&logical.Request{
				Operation: logical.UpdateOperation,
				Path:      "revoked/7f:e8:e1:29:31:41:9e:a4:ac:df:82:08:d1:64:b5:2f:84:2c:6d:b0",
				Storage:   storage,
			},
			"7f:e8:e1:29:31:41:9e:a4:ac:df:82:08:d1:64:b5:2f:84:2c:6d:b0",
			encodedCertValue,
		},
		"cert, valid hyphen": {
			&logical.Request{
				Operation: logical.UpdateOperation,
				Path:      "certs/7f:e8:e1:29:31:41:9e:a4:ac:df:82:08:d1:64:b5:2f:84:2c:6d:b0",
				Storage:   storage,
			},
			"7f-e8-e1-29-31-41-9e-a4-ac-df-82-08-d1-64-b5-2f-84-2c-6d-b0",
			certValue,
		},
		"cert, revoked hyphen": {
			&logical.Request{
				Operation: logical.UpdateOperation,
				Path:      "revoked/7f:e8:e1:29:31:41:9e:a4:ac:df:82:08:d1:64:b5:2f:84:2c:6d:b0",
				Storage:   storage,
			},
			"7f-e8-e1-29-31-41-9e-a4-ac-df-82-08-d1-64-b5-2f-84-2c-6d-b0",
			encodedCertValue,
		},
	}

	for name, tc := range cases {
		// Put pseudo-cert in inmem storage
		err := storage.Put(&logical.StorageEntry{
			Key:   tc.Req.Path,
			Value: tc.StorageValue,
		})
		if err != nil {
			t.Fatalf("error writing to storage on %s: %s", name, err)
		}

		_, err = revokeCert(b, tc.Req, tc.StorageKey, false)
		if err != nil {
			t.Fatalf("revokeCert error on %s: %s", name, err)
		}
	}
}
