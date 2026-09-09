// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package sarif

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	gitclient "github.com/hashicorp/vault/tools/pipeline/internal/pkg/git/client"
	"github.com/jedib0t/go-pretty/v6/table"
)

const (
	pushRetryBaseDelay = 2 * time.Second
	pushRetryJitterMax = 2 * time.Second
)

// UploadUSTArchiveReq converts a SARIF file to Concert format, packages it as
// a UST dynamic-scan archive, and pushes it to the UST destination GitHub
// repository.
//
// Authentication comes from the GITHUB_TOKEN or GH_TOKEN environment variable.
// Use actions/create-github-app-token in your workflow to mint a short-lived
// installation token and export it via one of those variables before calling Run.
type UploadUSTArchiveReq struct {
	SarifPath      string // path to the input SARIF file to convert and upload
	OfferingID     string // UST offering identifier written into metadata.json
	ProductID      string // UST product identifier; also the destination repo name under Org
	SquadID        string // UST squad identifier; used as the final path segment under results/<product>/<squad>/
	SourceRepoURL  string // URL of the repository that was scanned; written into repo_metadata.json
	SourceBranch   string // branch of the scanned repository; written into repo_metadata.json
	SourceCommit   string // commit SHA of the scanned repository; written into repo_metadata.json
	Host           string // GitHub host for the UST destination repository (e.g. github.ibm.com)
	Org            string // GitHub organization that owns the UST destination repository
	DestBranch     string // branch to push to in the UST destination repository
	GitName        string // git committer name used when pushing the results commit
	GitEmail       string // git committer email used when pushing the results commit
	ConcertOutPath string // if non-empty, the Concert JSON is also written to this path on disk
	UploadRetries  int    // max retry attempts when pushing the results
}

// UploadUSTArchiveRes is the response returned by UploadUSTArchiveReq.Run().
type UploadUSTArchiveRes struct {
	DestinationURL string `json:"destination_url"`
	SourceBranch   string `json:"source_branch"`
	SourceCommit   string `json:"source_commit"`
	ArchiveName    string `json:"archive_name"`
}

// Run executes the full conversion of the SARIF file into Concert JSON, builds
// the UST archive, and pushes it to the UST destination repo.
//
// A GitHub token must be present in GITHUB_TOKEN or GH_TOKEN before calling
// Run. Use actions/create-github-app-token in your workflow to mint a
// short-lived installation token and export it via one of those variables.
func (r *UploadUSTArchiveReq) Run(ctx context.Context) (*UploadUSTArchiveRes, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	// Convert Zap SARIF to Concert JSON
	convertRes, err := (&ConvertZapReq{SarifPath: r.SarifPath}).Run(ctx)
	if err != nil {
		return nil, fmt.Errorf("converting SARIF to Concert format: %w", err)
	}

	concertBytes, err := convertRes.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("serializing Concert JSON: %w", err)
	}

	if r.ConcertOutPath != "" {
		slog.Default().DebugContext(ctx, "writing concert JSON to output path", slog.String("path", r.ConcertOutPath))
		if err := os.WriteFile(r.ConcertOutPath, concertBytes, 0o644); err != nil {
			return nil, fmt.Errorf("writing concert JSON to %s: %w", r.ConcertOutPath, err)
		}
	}

	// Create the UST archive with our Concert data.
	archiveReq := &CreateUSTArchiveReq{
		ConcertData:   concertBytes,
		OfferingID:    r.OfferingID,
		ProductID:     r.ProductID,
		SquadID:       r.SquadID,
		SourceRepoURL: r.SourceRepoURL,
		SourceBranch:  r.SourceBranch,
		SourceCommit:  r.SourceCommit,
	}
	archiveRes, err := archiveReq.Run(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating UST archive: %w", err)
	}
	defer os.RemoveAll(archiveRes.WorkDir)

	// Clone, stage, commit, push the archive.
	// Token is read from GITHUB_TOKEN or GH_TOKEN by WithLoadTokenFromEnv.
	gitClient := gitclient.NewClient(
		gitclient.WithLoadTokenFromEnv(),
		gitclient.WithHost(r.Host),
		gitclient.WithConfig(r.gitConfig()),
	)

	pushedSHA, err := r.upload(ctx, gitClient, archiveRes)
	if err != nil {
		return nil, fmt.Errorf("delivering UST archive: %w", err)
	}

	destURL := fmt.Sprintf("https://%s/%s/%s/commit/%s",
		r.Host, r.Org, r.ProductID, pushedSHA)

	return &UploadUSTArchiveRes{
		DestinationURL: destURL,
		SourceBranch:   r.SourceBranch,
		SourceCommit:   r.SourceCommit,
		ArchiveName:    "dynamic-scan.tar.gz",
	}, nil
}

// ToTable returns a table.Writer with the upload result.
func (r *UploadUSTArchiveRes) ToTable() table.Writer {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"DESTINATION", "BRANCH", "COMMIT", "ARCHIVE"})
	if r != nil {
		t.AppendRow(table.Row{r.DestinationURL, r.SourceBranch, r.SourceCommit, r.ArchiveName})
	}
	return t
}

// ToJSON marshals the response to indented JSON.
func (r *UploadUSTArchiveRes) ToJSON() ([]byte, error) {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshaling upload result to JSON: %w", err)
	}
	return append(b, '\n'), nil
}

// ToMarkdown returns a Markdown table of the result.
func (r *UploadUSTArchiveRes) ToMarkdown() string {
	return r.ToTable().RenderMarkdown()
}

// String returns a one-line summary.
func (r *UploadUSTArchiveRes) String() string {
	if r == nil {
		return "<nil>"
	}
	return fmt.Sprintf("uploaded %s to %s (branch %s, commit %s)",
		r.ArchiveName, r.DestinationURL, r.SourceBranch, r.SourceCommit)
}

// upload clones the destination repo, copies the archive and metadata files,
// commits, and delegates the push to doUpload. It saves and restores the
// working directory so callers don't have to worry about it.
func (r *UploadUSTArchiveReq) upload(
	ctx context.Context,
	git *gitclient.Client,
	archiveRes *CreateUSTArchiveRes,
) (string, error) {
	originalDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getting working directory: %w", err)
	}
	defer func() { _ = os.Chdir(originalDir) }()

	cloneURL := fmt.Sprintf("https://%s/%s/%s", r.Host, r.Org, r.ProductID)
	cloneDir, err := os.MkdirTemp("", "pipeline-ust-clone-")
	if err != nil {
		return "", fmt.Errorf("creating clone temp dir: %w", err)
	}
	defer os.RemoveAll(cloneDir)

	slog.Default().DebugContext(ctx, "cloning destination repo",
		slog.String("url", cloneURL),
		slog.String("branch", r.DestBranch),
	)
	if _, err := git.Clone(ctx, &gitclient.CloneOpts{
		Depth:        1,
		Branch:       r.DestBranch,
		SingleBranch: true,
		Repository:   cloneURL,
		Directory:    cloneDir,
	}); err != nil {
		return "", fmt.Errorf("cloning destination repo: %w", err)
	}

	if err := os.Chdir(cloneDir); err != nil {
		return "", fmt.Errorf("changing to clone dir: %w", err)
	}

	resultsPath := filepath.Join("results", r.ProductID, r.SquadID)
	if err := os.MkdirAll(resultsPath, 0o755); err != nil {
		return "", fmt.Errorf("creating results path %s: %w", resultsPath, err)
	}

	for _, src := range []string{archiveRes.ArchivePath, archiveRes.MetadataPath} {
		dst := filepath.Join(resultsPath, filepath.Base(src))
		slog.Default().DebugContext(ctx, "copying file to clone", slog.String("src", src), slog.String("dst", dst))
		if err := copyFile(src, dst); err != nil {
			return "", fmt.Errorf("copying %s: %w", filepath.Base(src), err)
		}
	}

	if _, err := git.Add(ctx, &gitclient.AddOpts{PathSpec: []string{resultsPath}}); err != nil {
		return "", fmt.Errorf("adding results to git: %w", err)
	}

	commitMsg := r.commitMessage()
	if _, err := git.Commit(ctx, &gitclient.CommitOpts{Message: commitMsg}); err != nil {
		return "", fmt.Errorf("commit results: %w", err)
	}

	pushOpts := &gitclient.PushOpts{Repository: "origin", Refspec: []string{r.DestBranch}}
	pullFn := func(ctx context.Context) error {
		_, err := git.Pull(ctx, &gitclient.PullOpts{
			Rebase:     gitclient.RebaseStrategyTrue,
			Repository: "origin",
			Refspec:    []string{r.DestBranch},
		})
		return err
	}
	pushFn := func(ctx context.Context) error {
		_, err := git.Push(ctx, pushOpts)
		return err
	}
	if err := r.doUpload(ctx, pushFn, pullFn); err != nil {
		return "", err
	}

	revRes, err := git.RevParse(ctx, &gitclient.RevParseOpts{Args: []string{"HEAD"}})
	if err != nil {
		return "", fmt.Errorf("git rev-parse HEAD: %w", err)
	}
	return strings.TrimSpace(string(revRes.Stdout)), nil
}

// doUpload pushes our results to the remote host. Because we run scans in
// parallel, several jobs often race to push at the same time. When a push
// fails we pull --rebase and try again, up to UploadRetries times with a
// short randomised delay between attempts. Both push and pull failures are
// retryable. We extracted this from upload so we can test the retry logic
// without needing a real repo.
func (r *UploadUSTArchiveReq) doUpload(
	ctx context.Context,
	pushFn func(context.Context) error,
	pullFn func(context.Context) error,
) error {
	if err := pushFn(ctx); err == nil {
		return nil
	}
	maxRetries := r.UploadRetries
	if maxRetries < 1 {
		maxRetries = 1
	}
	var err error
	for attempt := range maxRetries {
		attempt++
		delay := pushRetryBaseDelay + time.Duration(rand.Int63n(int64(pushRetryJitterMax)))
		slog.Default().DebugContext(ctx, "push failed, retrying",
			slog.Int("attempt", attempt),
			slog.Int("max_retries", maxRetries),
			slog.Duration("delay", delay),
		)
		time.Sleep(delay)
		if pullErr := pullFn(ctx); pullErr != nil {
			slog.Default().DebugContext(ctx, "pull --rebase failed, will retry",
				slog.Int("attempt", attempt),
				slog.Int("max_retries", maxRetries),
				slog.String("error", pullErr.Error()),
			)
			err = errors.Join(err, pullErr)
			continue
		}
		if pushErr := pushFn(ctx); pushErr != nil {
			err = errors.Join(err, pushErr)
			continue
		}

		return nil
	}

	return fmt.Errorf("failed to upload results: %w", err)
}

// commitMessage builds the commit message using CI environment variables.
func (r *UploadUSTArchiveReq) commitMessage() string {
	repo := os.Getenv("GITHUB_REPOSITORY")
	if repo == "" {
		repo = "unknown"
	}
	runID := os.Getenv("GITHUB_RUN_ID")
	if runID == "" {
		runID = "?"
	}
	runAttempt := os.Getenv("GITHUB_RUN_ATTEMPT")
	if runAttempt == "" {
		runAttempt = "?"
	}
	shortSHA := r.SourceCommit
	if len(shortSHA) > 8 {
		shortSHA = shortSHA[:8]
	}
	return fmt.Sprintf("scan results upload for %s@%s (GHA run %s/%s)",
		repo, shortSHA, runID, runAttempt)
}

// gitConfig returns the user.name / user.email overrides for the git client.
func (r *UploadUSTArchiveReq) gitConfig() map[string]string {
	return map[string]string{
		"user.name":  r.GitName,
		"user.email": r.GitEmail,
	}
}

func (r *UploadUSTArchiveReq) validate() error {
	if r == nil {
		return errors.New("uninitialized request")
	}
	switch {
	case r.SarifPath == "":
		return errors.New("SarifPath is required")
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
	case r.Host == "":
		return errors.New("Host is required")
	case r.Org == "":
		return errors.New("Org is required")
	case r.DestBranch == "":
		return errors.New("DestBranch is required")
	}
	return nil
}

// copyFile copies src to dst, creating or truncating dst.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
