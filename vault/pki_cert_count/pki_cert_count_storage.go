// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package pki_cert_count

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/logical"
)

// PkiCertificateCountSubPath is the subpath under which PKI certificate counts
// are stored. It should be configured to use local storage, since each cluster
// needs to keep its own count.
const PkiCertificateCountSubPath = "pki-certificate-counts/"

// storageFormatPkiCertificateCount is the format of the storage path used to store counts of PKI certificates
const storageFormatPkiCertificateCount = PkiCertificateCountSubPath + "%d/%02d"

type PkiCertificateCount struct {
	IssuedCertificateCountsByDay []uint64 `json:"issuedCertificateCountsByDay"`
	StoredCertificateCountsByDay []uint64 `json:"storedCertificateCountsByDay"`
}

func IncrementStoredCounts(ctx context.Context, storage logical.Storage, issuedCount, storedCount uint64) error {
	year, month, day := time.Now().Date()

	storagePath := getStoragePath(year, month)
	var currentMonthCounts PkiCertificateCount

	entry, err := storage.Get(ctx, storagePath)
	if err != nil {
		return fmt.Errorf("error reading from storage: %w", err)
	}

	if entry != nil {
		err := json.Unmarshal(entry.Value, &currentMonthCounts)
		if err != nil {
			return fmt.Errorf("error unmarshalling storage entry for PKI certificate counts: %w", err)
		}
	} else {
		daysInMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
		currentMonthCounts.IssuedCertificateCountsByDay = make([]uint64, daysInMonth+1)
		currentMonthCounts.StoredCertificateCountsByDay = make([]uint64, daysInMonth+1)
	}

	currentMonthCounts.IssuedCertificateCountsByDay[day] += issuedCount
	currentMonthCounts.StoredCertificateCountsByDay[day] += storedCount

	countBytes, err := json.Marshal(currentMonthCounts)
	if err != nil {
		return fmt.Errorf("error marshalling certificate counts for storage: %w", err)
	}

	newEntry := &logical.StorageEntry{
		Key:   storagePath,
		Value: countBytes,
	}

	err = storage.Put(ctx, newEntry)
	if err != nil {
		return fmt.Errorf("error writing certificate counts to storage: %w", err)
	}

	return nil
}

func ReadStoredCounts(ctx context.Context, storage logical.Storage, date time.Time) (issuedCount uint64, storedCount uint64, err error) {
	issuedCount = 0
	storedCount = 0

	year, month, _ := date.Date()

	entry, err := storage.Get(ctx, getStoragePath(year, month))
	if err != nil {
		return 0, 0, fmt.Errorf("error reading from storage: %w", err)
	}

	if entry == nil {
		return 0, 0, fmt.Errorf("certificate counts not found for %d-%02d", year, month)
	}

	var certificateCounts PkiCertificateCount
	err = json.Unmarshal(entry.Value, &certificateCounts)
	if err != nil {
		return 0, 0, fmt.Errorf("error unmarshalling certificate counts from storage: %w", err)
	}

	for i := range certificateCounts.IssuedCertificateCountsByDay {
		issuedCount += certificateCounts.IssuedCertificateCountsByDay[i]
		storedCount += certificateCounts.StoredCertificateCountsByDay[i]
	}

	return issuedCount, storedCount, nil
}

func getStoragePath(year int, month time.Month) string {
	return fmt.Sprintf(storageFormatPkiCertificateCount, year, month)
}
