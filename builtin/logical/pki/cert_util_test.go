package pki

import (
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestFetchCertBySerial(t *testing.T) {
	storage := &logical.InmemStorage{}

	cases := map[string]struct {
		Req        *logical.Request
		StorageKey string
	}{
		"cert, valid colon": {
			&logical.Request{
				Operation: logical.ReadOperation,
				Path:      "certs/10:e6:fc:62:b7:41:8a:d5:00:5e:45:b6",
				Storage:   storage,
			},
			"10:e6:fc:62:b7:41:8a:d5:00:5e:45:b6",
		},
		"cert, revoked colon": {
			&logical.Request{
				Operation: logical.ReadOperation,
				Path:      "revoked/10:e6:fc:62:b7:41:8a:d5:00:5e:45:b6",
				Storage:   storage,
			},
			"10:e6:fc:62:b7:41:8a:d5:00:5e:45:b6",
		},
		"cert, valid hyphen": {
			&logical.Request{
				Operation: logical.ReadOperation,
				Path:      "certs/10:e6:fc:62:b7:41:8a:d5:00:5e:45:b6",
				Storage:   storage,
			},
			"10-e6-fc-62-b7-41-8a-d5-00-5e-45-b6",
		},
		"cert, revoked hyphen": {
			&logical.Request{
				Operation: logical.ReadOperation,
				Path:      "revoked/10:e6:fc:62:b7:41:8a:d5:00:5e:45:b6",
				Storage:   storage,
			},
			"10-e6-fc-62-b7-41-8a-d5-00-5e-45-b6",
		},
		"cert, ca": {
			&logical.Request{
				Operation: logical.ReadOperation,
				Path:      "ca",
				Storage:   storage,
			},
			"",
		},
		"cert, crl": {
			&logical.Request{
				Operation: logical.ReadOperation,
				Path:      "crl",
				Storage:   storage,
			},
			"",
		},
	}

	for name, tc := range cases {
		// Put pseudo-cert in inmem storage
		err := storage.Put(&logical.StorageEntry{
			Key:   tc.Req.Path,
			Value: []byte("some data"),
		})
		if err != nil {
			t.Fatalf("error writing to storage on %s: %s", name, err)
		}

		certEntry, err := fetchCertBySerial(tc.Req, tc.Req.Path, tc.StorageKey)
		if err != nil {
			t.Fatalf("fetchBySerial error on %s: %s", name, err)
		}

		// Check for non-nil on valid/revoked certs
		if certEntry == nil && tc.Req.Path != "ca" && tc.Req.Path != "crl" { // if true
			t.Fatalf("fetchBySerial returned nil on %s", name)
		}
	}
}
