// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package sarif

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

// CreateUSTArchiveReq creates the three files UST expects for a dynamic scan
// upload: metadata.json, repo_metadata.json, and dynamic-scan.tar.gz.
type CreateUSTArchiveReq struct {
	ConcertData   []byte // Concert-format JSON to package into the archive
	OfferingID    string
	ProductID     string
	SquadID       string
	SourceRepoURL string
	SourceBranch  string
	SourceCommit  string
}

// CreateUSTArchiveRes holds paths to the files produced by
// CreateUSTArchiveReq.Run(). The caller may call os.RemoveAll(res.WorkDir)
// to clean up files when/if they are no longer needed.
type CreateUSTArchiveRes struct {
	WorkDir          string // temp dir created by Run()
	MetadataPath     string // path to metadata.json inside WorkDir
	RepoMetadataPath string // path to repo_metadata.json inside WorkDir
	ArchivePath      string // path to dynamic-scan.tar.gz inside WorkDir
}

type ustMetadataEntry struct {
	Offering string `json:"offering"`
	Product  string `json:"product"`
	Squad    string `json:"squad"`

	// OfferingVersion, ProductVersion, and SquadVersion are always "latest".
	//
	// These fields represent the software-component version that was scanned,
	// not the source-code revision. For a DAST scan running against a live
	// deployment there is no meaningful semver to record here: the authoritative
	// traceability information (branch and commit SHA) lives in repo_metadata.json
	// as gitBranch and gitCommitSha, which we do populate dynamically.
	OfferingVersion string `json:"offeringVersion"`
	ProductVersion  string `json:"productVersion"`
	SquadVersion    string `json:"squadVersion"`
}

type ustMetadata struct {
	Metadata []ustMetadataEntry `json:"metadata"`
}

type ustRepoMetadataEntry struct {
	GitRepoURL   string `json:"gitRepoUrl"`
	GitBranch    string `json:"gitBranch"`
	GitCommitSHA string `json:"gitCommitSha"`
}

// Run validates the request, creates a temp directory, writes metadata.json,
// repo_metadata.json, and dynamic-scan.tar.gz into it, and returns the paths.
func (r *CreateUSTArchiveReq) Run(ctx context.Context) (*CreateUSTArchiveRes, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	workDir, err := os.MkdirTemp("", "pipeline-ust-archive-")
	if err != nil {
		return nil, fmt.Errorf("creating temp work dir: %w", err)
	}

	res := &CreateUSTArchiveRes{WorkDir: workDir}

	metadataPath, err := r.writeMetadata(ctx, workDir)
	if err != nil {
		return nil, err
	}
	res.MetadataPath = metadataPath

	repoMetadataPath, err := r.writeRepoMetadata(ctx, workDir)
	if err != nil {
		return nil, err
	}
	res.RepoMetadataPath = repoMetadataPath

	archivePath, err := r.writeArchive(ctx, workDir, r.ConcertData)
	if err != nil {
		return nil, err
	}
	res.ArchivePath = archivePath

	return res, nil
}

func (r *CreateUSTArchiveReq) validate() error {
	if r == nil {
		return errors.New("uninitialized request")
	}
	switch {
	case len(r.ConcertData) == 0:
		return errors.New("ConcertData is required")
	case r.OfferingID == "":
		return errors.New("OfferingID is required")
	case r.ProductID == "":
		return errors.New("ProductID is required")
	case r.SquadID == "":
		return errors.New("SquadID is required")
	case r.SourceRepoURL == "":
		return errors.New("SourceRepoURL is required")
	case r.SourceBranch == "":
		return errors.New("SourceBranch is required")
	case r.SourceCommit == "":
		return errors.New("SourceCommit is required")
	}
	return nil
}

func (r *CreateUSTArchiveReq) writeMetadata(ctx context.Context, workDir string) (string, error) {
	m := ustMetadata{
		Metadata: []ustMetadataEntry{{
			Offering:        r.OfferingID,
			Product:         r.ProductID,
			Squad:           r.SquadID,
			OfferingVersion: "latest",
			ProductVersion:  "latest",
			SquadVersion:    "latest",
		}},
	}

	b, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("marshaling metadata.json: %w", err)
	}

	path := filepath.Join(workDir, "metadata.json")
	slog.Default().DebugContext(ctx, "writing metadata.json", slog.String("path", path))
	if err := os.WriteFile(path, b, 0o644); err != nil {
		return "", fmt.Errorf("writing metadata.json: %w", err)
	}
	return path, nil
}

func (r *CreateUSTArchiveReq) writeRepoMetadata(ctx context.Context, workDir string) (string, error) {
	entries := []ustRepoMetadataEntry{{
		GitRepoURL:   r.SourceRepoURL,
		GitBranch:    r.SourceBranch,
		GitCommitSHA: r.SourceCommit,
	}}

	b, err := json.Marshal(entries)
	if err != nil {
		return "", fmt.Errorf("marshaling repo_metadata.json: %w", err)
	}

	path := filepath.Join(workDir, "repo_metadata.json")
	slog.Default().DebugContext(ctx, "writing repo_metadata.json", slog.String("path", path))
	if err := os.WriteFile(path, b, 0o644); err != nil {
		return "", fmt.Errorf("writing repo_metadata.json: %w", err)
	}
	return path, nil
}

func (r *CreateUSTArchiveReq) writeArchive(ctx context.Context, workDir string, concertData []byte) (archivePath string, err error) {
	archivePath = filepath.Join(workDir, "dynamic-scan.tar.gz")
	slog.Default().DebugContext(ctx, "writing dynamic-scan.tar.gz", slog.String("path", archivePath))

	f, err := os.Create(archivePath)
	if err != nil {
		return "", fmt.Errorf("creating dynamic-scan.tar.gz: %w", err)
	}
	isClosed := false
	defer func() {
		if isClosed {
			return
		}
		err = errors.Join(err, f.Close())
	}()

	gw := gzip.NewWriter(f)
	tw := tar.NewWriter(gw)

	hdr := &tar.Header{
		Name: "zap/report-concert.json",
		Mode: 0o644,
		Size: int64(len(concertData)),
	}
	if err = tw.WriteHeader(hdr); err != nil {
		return "", fmt.Errorf("writing tar header: %w", err)
	}
	if _, err = io.Copy(tw, bytes.NewReader(concertData)); err != nil {
		return "", fmt.Errorf("writing concert JSON to tar: %w", err)
	}

	// Close in reverse order: tar → gzip → file.
	// Each layer must be closed before the next to flush its trailer bytes.
	// Errors here mean the archive is incomplete or corrupt, so we must
	// surface them rather than silently swallow them via defer.
	// The deferred f.Close() is guarded by isClosed to prevent a double-close
	// on the success path.
	if err = tw.Close(); err != nil {
		return "", fmt.Errorf("closing tar writer: %w", err)
	}
	if err = gw.Close(); err != nil {
		return "", fmt.Errorf("closing gzip writer: %w", err)
	}
	if err = f.Close(); err != nil {
		return "", fmt.Errorf("closing archive file: %w", err)
	}
	isClosed = true

	return archivePath, nil
}
