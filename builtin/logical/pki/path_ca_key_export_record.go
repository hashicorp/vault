// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"context"
	"crypto/mlkem"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// caKeyExportRecordPrefix is the storage prefix for CA key export audit records.
// Records are stored at:
//
//	caKeyExportRecordPrefix + <spki-fingerprint> + "/" + <timestamp>
//
// Keying on the SHA-256 fingerprint of the SubjectPublicKeyInfo (SPKI) rather
// than the KeyEntry UUID means records survive a delete-and-reimport cycle: the
// key material, and therefore the fingerprint, is identical after reimport.
const caKeyExportRecordPrefix = "ca-key-export-record/"

// CAKeyExportRecord is persisted once per successful WRITE /pki/keys/:ca-key-uuid/export.
// It is intentionally independent of KeyEntry so that deleting or reimporting a
// CA key does not erase the audit trail.
type CAKeyExportRecord struct {
	// ExportKeyHMAC is the HMAC of the wrapping public key used to encrypt the CA
	// private key material.  It correlates this record with the export key on the
	// destination Vault instance.
	ExportKeyHMAC string `json:"export_key_hmac"`
	// ExportedAt is the UTC timestamp of the export.
	ExportedAt time.Time `json:"exported_at"`
}

// caPublicKeyFingerprint returns a "sha256:<hex>" fingerprint over the
// SubjectPublicKeyInfo (SPKI) DER encoding of the public key embedded in
// privKeyPEM.  The fingerprint is stable across delete-and-reimport cycles
// because it depends only on the key material, not the KeyEntry UUID.
//
// Supported PEM formats: SEC1 EC ("EC PRIVATE KEY"), PKCS#1 RSA ("RSA PRIVATE
// KEY"), and PKCS#8 ("PRIVATE KEY") — matching the formats Vault itself generates.
func caPublicKeyFingerprint(privKeyPEM string) (string, error) {
	block, _ := pem.Decode([]byte(privKeyPEM))
	if block == nil {
		return "", errutil.InternalError{Err: "unable to compute CA key fingerprint: no PEM block found"}
	}

	if block.Type == "ML-KEM PRIVATE KEY" {
		var pubBytes []byte
		alg := block.Headers["Algorithm"]
		switch alg {
		case "ML-KEM-768":
			dk768, err := mlkem.NewDecapsulationKey768(block.Bytes)
			if err != nil {
				return "", errutil.InternalError{Err: fmt.Sprintf("unable to compute CA key fingerprint: invalid ML-KEM-768 private key: %v", err)}
			}
			pubBytes = dk768.EncapsulationKey().Bytes()
		case "ML-KEM-1024":
			dk1024, err := mlkem.NewDecapsulationKey1024(block.Bytes)
			if err != nil {
				return "", errutil.InternalError{Err: fmt.Sprintf("unable to compute CA key fingerprint: invalid ML-KEM-1024 private key: %v", err)}
			}
			pubBytes = dk1024.EncapsulationKey().Bytes()
		default:
			return "", errutil.InternalError{Err: fmt.Sprintf("unable to compute CA key fingerprint: unrecognised ML-KEM Algorithm header %q", alg)}
		}
		sum := sha256.Sum256(pubBytes)
		return "sha256:" + hex.EncodeToString(sum[:]), nil
	}

	// certutil.ParseDERKey handles SEC1 EC, PKCS#1 RSA, and PKCS#8 in one call.
	signer, _, err := certutil.ParseDERKey(block.Bytes)
	if err != nil {
		return "", errutil.InternalError{Err: fmt.Sprintf("unable to compute CA key fingerprint: failed to parse private key: %v", err)}
	}

	spkiDER, err := x509.MarshalPKIXPublicKey(signer.Public())
	if err != nil {
		return "", errutil.InternalError{Err: fmt.Sprintf("unable to compute CA key fingerprint: failed to marshal SPKI: %v", err)}
	}

	sum := sha256.Sum256(spkiDER)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

// writeCAKeyExportRecord appends a new CAKeyExportRecord to storage under the
// given SPKI fingerprint.  Multiple exports of the same key produce distinct,
// chronologically ordered entries.
func writeCAKeyExportRecord(ctx context.Context, s logical.Storage, spkiFingerprint string, record CAKeyExportRecord) error {
	// Use the nanosecond-precision RFC3339 timestamp as the leaf key so records
	// sort chronologically and collisions are practically impossible.
	index := record.ExportedAt.UTC().Format(time.RFC3339Nano)
	path := caKeyExportRecordPrefix + spkiFingerprint + "/" + index

	entry, err := logical.StorageEntryJSON(path, record)
	if err != nil {
		return errutil.InternalError{Err: fmt.Sprintf("unable to marshal CA key export record: %v", err)}
	}

	if err := s.Put(ctx, entry); err != nil {
		return errutil.InternalError{Err: fmt.Sprintf("unable to store CA key export record: %v", err)}
	}

	return nil
}

// listCAKeyExportRecords returns all CAKeyExportRecord values persisted for the
// given SPKI fingerprint.  An empty slice is returned when no exports have been
// recorded.
func listCAKeyExportRecords(ctx context.Context, s logical.Storage, spkiFingerprint string) ([]CAKeyExportRecord, error) {
	prefix := caKeyExportRecordPrefix + spkiFingerprint + "/"

	keys, err := s.List(ctx, prefix)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to list CA key export records: %v", err)}
	}

	records := make([]CAKeyExportRecord, 0, len(keys))
	for _, key := range keys {
		raw, err := s.Get(ctx, prefix+key)
		if err != nil {
			return nil, errutil.InternalError{Err: fmt.Sprintf("unable to read CA key export record %s: %v", key, err)}
		}
		if raw == nil {
			continue
		}

		var rec CAKeyExportRecord
		if err := raw.DecodeJSON(&rec); err != nil {
			return nil, errutil.InternalError{Err: fmt.Sprintf("unable to decode CA key export record %s: %v", key, err)}
		}

		records = append(records, rec)
	}

	return records, nil
}
