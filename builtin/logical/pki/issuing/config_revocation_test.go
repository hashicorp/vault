// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package issuing

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

// deleteRecordingStorage records every path handed to Delete so a test can
// check the exact storage paths a cleanup routine produces.
type deleteRecordingStorage struct {
	logical.Storage
	deleted []string
}

func (s *deleteRecordingStorage) Delete(ctx context.Context, path string) error {
	s.deleted = append(s.deleted, path)
	return s.Storage.Delete(ctx, path)
}

// Test_cleanupInternalCRLMapping_OrphanedCRLPath checks that a crl left on
// disk without a mapping entry is removed through a well-formed storage path,
// while the full and delta crls of a referenced issuer are kept.
// a path with an empty segment is rejected by some physical backends, such as
// zookeeper, and on backends that key on the literal string, such as raft, it
// matches nothing and the orphan stays behind.
func Test_cleanupInternalCRLMapping_OrphanedCRLPath(t *testing.T) {
	tests := []struct {
		name       string
		configPath string
		base       string
	}{
		{name: "local", configPath: StorageLocalCRLConfig, base: PathCrls},
		{name: "unified", configPath: StorageUnifiedCRLConfig, base: "unified-crls/"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			s := &deleteRecordingStorage{Storage: &logical.InmemStorage{}}

			mapping := &InternalCRLConfigEntry{
				IssuerIDCRLMap:        map[IssuerID]CrlID{"issuer": "referenced"},
				CRLNumberMap:          map[CrlID]int64{"referenced": 1},
				LastCompleteNumberMap: map[CrlID]int64{"referenced": 1},
			}
			for _, id := range []string{"referenced", "referenced" + DeltaCRLPathSuffix, "orphan", "orphan" + DeltaCRLPathSuffix} {
				require.NoError(t, s.Put(ctx, &logical.StorageEntry{Key: tc.base + id, Value: []byte("crl")}))
			}

			require.NoError(t, _cleanupInternalCRLMapping(ctx, s, mapping, tc.configPath))

			for _, path := range s.deleted {
				require.NotContains(t, path, "//", "delete path %q has an empty segment", path)
			}
			require.Contains(t, s.deleted, tc.base+"orphan")
			require.Contains(t, s.deleted, tc.base+"orphan"+DeltaCRLPathSuffix)

			entry, err := s.Get(ctx, tc.base+"orphan")
			require.NoError(t, err)
			require.Nil(t, entry, "orphaned crl was not removed from storage")

			entry, err = s.Get(ctx, tc.base+"referenced")
			require.NoError(t, err)
			require.NotNil(t, entry, "referenced crl must be kept")

			// delta crls live next to their full crl with a suffix and belong to
			// the same issuer, so they are not orphans either
			entry, err = s.Get(ctx, tc.base+"referenced"+DeltaCRLPathSuffix)
			require.NoError(t, err)
			require.NotNil(t, entry, "referenced delta crl must be kept")
		})
	}
}
