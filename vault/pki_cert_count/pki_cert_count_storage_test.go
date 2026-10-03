// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package pki_cert_count

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical/inmem"
	"github.com/stretchr/testify/require"
)

func TestGetStoragePath(t *testing.T) {
	require.Equal(t, "pki-certificate-counts/2025/09", getStoragePath(2025, 9))
}

func TestGetCertificateCount(t *testing.T) {
	backend, err := inmem.NewInmem(nil, hclog.NewNullLogger())
	require.NoError(t, err)

	storage := logical.NewLogicalStorage(backend)
	currentTime := time.Now()

	writeCertCounts(t, storage, currentTime)
	writeCertCounts(t, storage, time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC))

	testCases := map[string]struct {
		date           time.Time
		expectErr      bool
		expectedIssued uint64
		expectedStored uint64
	}{
		"current month counts": {
			date:           currentTime,
			expectedIssued: 7,
			expectedStored: 4,
		},
		"past month counts": {
			date:           time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
			expectedIssued: 7,
			expectedStored: 4,
		},
		"counts not found": {
			date:      time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC),
			expectErr: true,
		},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			issuedCount, storedCount, err := ReadStoredCounts(context.Background(), storage, tt.date)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedIssued, issuedCount)
				require.Equal(t, tt.expectedStored, storedCount)
			}
		})
	}
}

func TestStoreCertificateCounts(t *testing.T) {
	backend, err := inmem.NewInmem(nil, hclog.NewNullLogger())
	require.NoError(t, err)

	storage := logical.NewLogicalStorage(backend)

	var expectedIssuedCount uint64 = 5
	var expectedStoredCount uint64 = 3

	err = IncrementStoredCounts(context.Background(), storage, expectedIssuedCount, expectedStoredCount)
	require.NoError(t, err)

	year, month, day := time.Now().Date()
	counts := retrieveCertificateCountsFromStorage(t, storage, year, month)

	require.Equal(t, expectedIssuedCount, counts.IssuedCertificateCountsByDay[day])
	require.Equal(t, expectedStoredCount, counts.StoredCertificateCountsByDay[day])
}

func TestReadAfterStore(t *testing.T) {
	backend, err := inmem.NewInmem(nil, hclog.NewNullLogger())
	require.NoError(t, err)

	storage := logical.NewLogicalStorage(backend)

	var expectedIssuedCount uint64 = 5
	var expectedStoredCount uint64 = 3

	err = IncrementStoredCounts(context.Background(), storage, expectedIssuedCount, expectedStoredCount)
	require.NoError(t, err)

	issued, stored, err := ReadStoredCounts(context.Background(), storage, time.Now())
	require.NoError(t, err)
	require.Equal(t, expectedIssuedCount, issued)
	require.Equal(t, expectedStoredCount, stored)
}

func retrieveCertificateCountsFromStorage(t *testing.T, storage logical.Storage, year int, month time.Month) *PkiCertificateCount {
	storagePath := getStoragePath(year, month)

	entry, err := storage.Get(context.Background(), storagePath)
	require.NoError(t, err)

	require.NotNil(t, entry)

	var counts PkiCertificateCount
	err = json.Unmarshal(entry.Value, &counts)
	require.NoError(t, err)

	return &counts
}

func writeCertCounts(t *testing.T, storage logical.Storage, date time.Time) {
	storagePath := getStoragePath(date.Year(), date.Month())
	daysInMonth := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()

	issuedCounts := make([]uint64, daysInMonth+1)
	storedCounts := make([]uint64, daysInMonth+1)

	issuedCounts[2] = 3
	storedCounts[2] = 3

	issuedCounts[18] = 4
	storedCounts[18] = 1

	counts := PkiCertificateCount{
		issuedCounts,
		storedCounts,
	}

	countsBytes, err := json.Marshal(counts)
	require.NoError(t, err)

	entry := logical.StorageEntry{
		Key:   storagePath,
		Value: countsBytes,
	}

	err = storage.Put(context.Background(), &entry)
	require.NoError(t, err)
}
