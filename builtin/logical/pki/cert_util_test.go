package pki

import (
	"context"
	"fmt"
	"testing"

	"strings"

	"github.com/hashicorp/vault/logical"
)

func TestPki_FetchCertBySerial(t *testing.T) {
	storage := &logical.InmemStorage{}

	cases := map[string]struct {
		Req    *logical.Request
		Prefix string
		Serial string
	}{
		"valid cert": {
			&logical.Request{
				Storage: storage,
			},
			"certs/",
			"00:00:00:00:00:00:00:00",
		},
		"revoked cert": {
			&logical.Request{
				Storage: storage,
			},
			"revoked/",
			"11:11:11:11:11:11:11:11",
		},
	}

	// Test for colon-based paths in storage
	for name, tc := range cases {
		storageKey := fmt.Sprintf("%s%s", tc.Prefix, tc.Serial)
		err := storage.Put(context.Background(), &logical.StorageEntry{
			Key:   storageKey,
			Value: []byte("some data"),
		})
		if err != nil {
			t.Fatalf("error writing to storage on %s colon-based storage path: %s", name, err)
		}

		certEntry, err := fetchCertBySerial(context.Background(), tc.Req, tc.Prefix, tc.Serial)
		if err != nil {
			t.Fatalf("error on %s for colon-based storage path: %s", name, err)
		}

		// Check for non-nil on valid/revoked certs
		if certEntry == nil {
			t.Fatalf("nil on %s for colon-based storage path", name)
		}

		// Ensure that cert serials are converted/updated after fetch
		expectedKey := tc.Prefix + normalizeSerial(tc.Serial)
		se, err := storage.Get(context.Background(), expectedKey)
		if err != nil {
			t.Fatalf("error on %s for colon-based storage path:%s", name, err)
		}
		if strings.Compare(expectedKey, se.Key) != 0 {
			t.Fatalf("expected: %s, got: %s", expectedKey, certEntry.Key)
		}
	}

	// Reset storage
	storage = &logical.InmemStorage{}

	// Test for hyphen-base paths in storage
	for name, tc := range cases {
		storageKey := tc.Prefix + normalizeSerial(tc.Serial)
		err := storage.Put(context.Background(), &logical.StorageEntry{
			Key:   storageKey,
			Value: []byte("some data"),
		})
		if err != nil {
			t.Fatalf("error writing to storage on %s hyphen-based storage path: %s", name, err)
		}

		certEntry, err := fetchCertBySerial(context.Background(), tc.Req, tc.Prefix, tc.Serial)
		if err != nil || certEntry == nil {
			t.Fatalf("error on %s for hyphen-based storage path: err: %v, entry: %v", name, err, certEntry)
		}
	}

	noConvCases := map[string]struct {
		Req    *logical.Request
		Prefix string
		Serial string
	}{
		"ca": {
			&logical.Request{
				Storage: storage,
			},
			"",
			"ca",
		},
		"crl": {
			&logical.Request{
				Storage: storage,
			},
			"",
			"crl",
		},
	}

	// Test for ca and crl case
	for name, tc := range noConvCases {
		err := storage.Put(context.Background(), &logical.StorageEntry{
			Key:   tc.Serial,
			Value: []byte("some data"),
		})
		if err != nil {
			t.Fatalf("error writing to storage on %s: %s", name, err)
		}

		certEntry, err := fetchCertBySerial(context.Background(), tc.Req, tc.Prefix, tc.Serial)
		if err != nil || certEntry == nil {
			t.Fatalf("error on %s: err: %v, entry: %v", name, err, certEntry)
		}
	}
}
