// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pb

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/helper/wrapping"
	"github.com/hashicorp/vault/sdk/logical"
)

func TestTranslation_Errors(t *testing.T) {
	errs := []error{
		nil,
		errors.New("test"),
		errutil.UserError{Err: "test"},
		errutil.InternalError{Err: "test"},
		logical.CodedError(403, "test"),
		&logical.StatusBadRequest{Err: "test"},
		logical.ErrUnsupportedOperation,
		logical.ErrUnsupportedPath,
		logical.ErrInvalidRequest,
		logical.ErrPermissionDenied,
		logical.ErrMultiAuthzPending,
	}

	for _, err := range errs {
		pe := ErrToProtoErr(err)
		e := ProtoErrToErr(pe)

		if !reflect.DeepEqual(e, err) {
			t.Fatalf("Errs did not match: %#v, %#v", e, err)
		}
	}
}

func TestTranslation_StorageEntry(t *testing.T) {
	tCases := []*logical.StorageEntry{
		nil,
		{Key: "key", Value: []byte("value")},
		{Key: "key1", Value: []byte("value1"), SealWrap: true},
		{Key: "key1", SealWrap: true},
	}

	for _, c := range tCases {
		p := LogicalStorageEntryToProtoStorageEntry(c)
		e := ProtoStorageEntryToLogicalStorageEntry(p)

		if !reflect.DeepEqual(c, e) {
			t.Fatalf("Entries did not match: %#v, %#v", e, c)
		}
	}
}

func TestTranslation_Request(t *testing.T) {
	certs, err := peerCertificates()
	if err != nil {
		t.Logf("No test certificates were generated: %v", err)
	}

	tCases := []*logical.Request{
		nil,
		{
			ID:                       "ID",
			ReplicationCluster:       "RID",
			Operation:                logical.CreateOperation,
			Path:                     "test/foo",
			ClientToken:              "token",
			ClientTokenAccessor:      "accessor",
			DisplayName:              "display",
			MountPoint:               "test",
			MountType:                "secret",
			MountAccessor:            "test-231234",
			ClientTokenRemainingUses: 1,
			EntityID:                 "tester",
			PolicyOverride:           true,
			Unauthenticated:          true,
			Connection: &logical.Connection{
				RemoteAddr: "localhost",
				ConnState: &tls.ConnectionState{
					Version:           tls.VersionTLS12,
					HandshakeComplete: true,
					PeerCertificates:  certs,
				},
			},
		},
		{
			ID:                 "ID",
			ReplicationCluster: "RID",
			Operation:          logical.CreateOperation,
			Path:               "test/foo",
			Data: map[string]interface{}{
				"string": "string",
				"bool":   true,
				"array":  []interface{}{"1", "2"},
				"map": map[string]interface{}{
					"key": "value",
				},
			},
			Secret: &logical.Secret{
				LeaseOptions: logical.LeaseOptions{
					TTL:       time.Second,
					MaxTTL:    time.Second,
					Renewable: true,
					Increment: time.Second,
					IssueTime: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
				},
				InternalData: map[string]interface{}{
					"role": "test",
				},
				LeaseID: "LeaseID",
			},
			Auth: &logical.Auth{
				LeaseOptions: logical.LeaseOptions{
					TTL:       time.Second,
					MaxTTL:    time.Second,
					Renewable: true,
					Increment: time.Second,
					IssueTime: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
				},
				InternalData: map[string]interface{}{
					"role": "test",
				},
				DisplayName: "test",
				Policies:    []string{"test", "Test"},
				Metadata: map[string]string{
					"test": "test",
				},
				ClientToken: "token",
				Accessor:    "accessor",
				Period:      5 * time.Second,
				NumUses:     1,
				EntityID:    "id",
				Alias: &logical.Alias{
					MountType:     "type",
					MountAccessor: "accessor",
					Name:          "name",
				},
				GroupAliases: []*logical.Alias{
					{
						MountType:     "type",
						MountAccessor: "accessor",
						Name:          "name",
					},
				},
			},
			Headers: map[string][]string{
				"X-Vault-Test": {"test"},
			},
			ClientToken:         "token",
			ClientTokenAccessor: "accessor",
			DisplayName:         "display",
			MountPoint:          "test",
			MountType:           "secret",
			MountAccessor:       "test-231234",
			WrapInfo: &logical.RequestWrapInfo{
				TTL:      time.Second,
				Format:   "token",
				SealWrap: true,
			},
			ClientTokenRemainingUses: 1,
			EntityID:                 "tester",
			PolicyOverride:           true,
			Unauthenticated:          true,
		},
	}

	for _, c := range tCases {
		p, err := LogicalRequestToProtoRequest(c)
		if err != nil {
			t.Fatal(err)
		}
		r, err := ProtoRequestToLogicalRequest(p)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(c, r) {
			t.Fatalf("Requests did not match: \n%#v, \n%#v", c, r)
		}
	}
}

func TestTranslation_Response(t *testing.T) {
	tCases := []*logical.Response{
		nil,
		{
			Data: map[string]interface{}{
				"data": "blah",
			},
			Warnings: []string{"warning"},
		},
		{
			Data: map[string]interface{}{
				"string": "string",
				"bool":   true,
				"array":  []interface{}{"1", "2"},
				"map": map[string]interface{}{
					"key": "value",
				},
			},
			Secret: &logical.Secret{
				LeaseOptions: logical.LeaseOptions{
					TTL:       time.Second,
					MaxTTL:    time.Second,
					Renewable: true,
					Increment: time.Second,
					IssueTime: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
				},
				InternalData: map[string]interface{}{
					"role": "test",
				},
				LeaseID: "LeaseID",
			},
			Auth: &logical.Auth{
				LeaseOptions: logical.LeaseOptions{
					TTL:       time.Second,
					MaxTTL:    time.Second,
					Renewable: true,
					Increment: time.Second,
					IssueTime: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
				},
				InternalData: map[string]interface{}{
					"role": "test",
				},
				DisplayName: "test",
				Policies:    []string{"test", "Test"},
				Metadata: map[string]string{
					"test": "test",
				},
				ClientToken: "token",
				Accessor:    "accessor",
				Period:      5 * time.Second,
				NumUses:     1,
				EntityID:    "id",
				Alias: &logical.Alias{
					MountType:     "type",
					MountAccessor: "accessor",
					Name:          "name",
				},
				GroupAliases: []*logical.Alias{
					{
						MountType:     "type",
						MountAccessor: "accessor",
						Name:          "name",
					},
				},
			},
			WrapInfo: &wrapping.ResponseWrapInfo{
				TTL:             time.Second,
				Token:           "token",
				Accessor:        "accessor",
				CreationTime:    time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
				WrappedAccessor: "wrapped-accessor",
				WrappedEntityID: "id",
				Format:          "token",
				CreationPath:    "test/foo",
				SealWrap:        true,
			},
			MountType: "mountType",
		},
	}

	for _, c := range tCases {
		p, err := LogicalResponseToProtoResponse(c)
		if err != nil {
			t.Fatal(err)
		}
		r, err := ProtoResponseToLogicalResponse(p)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(c, r) {
			t.Fatalf("Requests did not match: \n%#v, \n%#v", c, r)
		}
	}
}

// This is the contents of $GOROOT/src/crypto/tls/testdata/example-cert.pem
// If it's good enough for testing the crypto/tls package it's good enough
// for Vault.
const exampleCert = `
-----BEGIN CERTIFICATE-----
MIIBhTCCASugAwIBAgIQIRi6zePL6mKjOipn+dNuaTAKBggqhkjOPQQDAjASMRAw
DgYDVQQKEwdBY21lIENvMB4XDTE3MTAyMDE5NDMwNloXDTE4MTAyMDE5NDMwNlow
EjEQMA4GA1UEChMHQWNtZSBDbzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABD0d
7VNhbWvZLWPuj/RtHFjvtJBEwOkhbN/BnnE8rnZR8+sbwnc/KhCk3FhnpHZnQz7B
5aETbbIgmuvewdjvSBSjYzBhMA4GA1UdDwEB/wQEAwICpDATBgNVHSUEDDAKBggr
BgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MCkGA1UdEQQiMCCCDmxvY2FsaG9zdDo1
NDUzgg4xMjcuMC4wLjE6NTQ1MzAKBggqhkjOPQQDAgNIADBFAiEA2zpJEPQyz6/l
Wf86aX6PepsntZv2GYlA5UpabfT2EZICICpJ5h/iI+i341gBmLiAFQOyTDT+/wQc
6MF9+Yw1Yy0t
-----END CERTIFICATE-----`

func peerCertificates() ([]*x509.Certificate, error) {
	blk, _ := pem.Decode([]byte(exampleCert))
	if blk == nil {
		return nil, errors.New("cannot decode example certificate")
	}

	cert, err := x509.ParseCertificate(blk.Bytes)
	if err != nil {
		return nil, err
	}

	return []*x509.Certificate{cert}, nil
}
