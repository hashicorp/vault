// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package sarif

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var testConcertData = []byte(`{"test":true}`)

// TestCreateUSTArchiveReq_Run verifies that Run() produces all three expected
// output files with the correct content and structure when given valid inputs.
func TestCreateUSTArchiveReq_Run(t *testing.T) {
	t.Parallel()

	req := &CreateUSTArchiveReq{
		ConcertData:   testConcertData,
		OfferingID:    "vault",
		ProductID:     "vault-enterprise",
		SquadID:       "vault-ent",
		SourceRepoURL: "https://github.com/hashicorp/vault-enterprise",
		SourceBranch:  "main",
		SourceCommit:  "abc123def456",
	}

	res, err := req.Run(t.Context())
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(res.WorkDir) })

	t.Run("all output files exist", func(t *testing.T) {
		require.FileExists(t, res.MetadataPath)
		require.FileExists(t, res.RepoMetadataPath)
		require.FileExists(t, res.ArchivePath)
	})

	t.Run("metadata.json content", func(t *testing.T) {
		data, readErr := os.ReadFile(res.MetadataPath)
		require.NoError(t, readErr)

		var m ustMetadata
		require.NoError(t, json.Unmarshal(data, &m))
		require.Len(t, m.Metadata, 1)
		entry := m.Metadata[0]
		require.Equal(t, "vault", entry.Offering)
		require.Equal(t, "vault-enterprise", entry.Product)
		require.Equal(t, "vault-ent", entry.Squad)
		require.Equal(t, "latest", entry.OfferingVersion)
		require.Equal(t, "latest", entry.ProductVersion)
		require.Equal(t, "latest", entry.SquadVersion)
	})

	t.Run("repo_metadata.json content", func(t *testing.T) {
		data, readErr := os.ReadFile(res.RepoMetadataPath)
		require.NoError(t, readErr)

		var entries []ustRepoMetadataEntry
		require.NoError(t, json.Unmarshal(data, &entries))
		require.Len(t, entries, 1)
		require.Equal(t, "https://github.com/hashicorp/vault-enterprise", entries[0].GitRepoURL)
		require.Equal(t, "main", entries[0].GitBranch)
		require.Equal(t, "abc123def456", entries[0].GitCommitSHA)
	})

	t.Run("tarball entry name is zap/report-concert.json", func(t *testing.T) {
		f, openErr := os.Open(res.ArchivePath)
		require.NoError(t, openErr)
		defer f.Close()

		gr, gzErr := gzip.NewReader(f)
		require.NoError(t, gzErr)
		defer gr.Close()

		tr := tar.NewReader(gr)
		hdr, nextErr := tr.Next()
		require.NoError(t, nextErr)
		require.Equal(t, "zap/report-concert.json", hdr.Name)

		content, readErr := io.ReadAll(tr)
		require.NoError(t, readErr)
		require.Equal(t, `{"test":true}`, string(content))

		_, eofErr := tr.Next()
		require.ErrorIs(t, eofErr, io.EOF)
	})
}

// TestCreateUSTArchiveReq_Run_ValidationErrors verifies that Run() returns
// descriptive errors for each missing required field.
func TestCreateUSTArchiveReq_Run_ValidationErrors(t *testing.T) {
	t.Parallel()

	for name, req := range map[string]*CreateUSTArchiveReq{
		"empty ConcertData": {
			OfferingID: "o", ProductID: "p", SquadID: "s",
			SourceRepoURL: "u", SourceBranch: "b", SourceCommit: "c",
		},
		"missing OfferingID": {
			ConcertData: testConcertData, ProductID: "p", SquadID: "s",
			SourceRepoURL: "u", SourceBranch: "b", SourceCommit: "c",
		},
		"missing ProductID": {
			ConcertData: testConcertData, OfferingID: "o", SquadID: "s",
			SourceRepoURL: "u", SourceBranch: "b", SourceCommit: "c",
		},
		"missing SquadID": {
			ConcertData: testConcertData, OfferingID: "o", ProductID: "p",
			SourceRepoURL: "u", SourceBranch: "b", SourceCommit: "c",
		},
		"missing SourceRepoURL": {
			ConcertData: testConcertData, OfferingID: "o", ProductID: "p", SquadID: "s",
			SourceBranch: "b", SourceCommit: "c",
		},
		"missing SourceBranch": {
			ConcertData: testConcertData, OfferingID: "o", ProductID: "p", SquadID: "s",
			SourceRepoURL: "u", SourceCommit: "c",
		},
		"missing SourceCommit": {
			ConcertData: testConcertData, OfferingID: "o", ProductID: "p", SquadID: "s",
			SourceRepoURL: "u", SourceBranch: "b",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, err := req.Run(t.Context())
			require.Error(t, err)
		})
	}
}
