// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"

	libgithub "github.com/google/go-github/v83/github"
	slogctx "github.com/veqryn/slog-context"
)

// branchVersionPattern matches release branch names of the form release/X.Y.x
var branchVersionPattern = regexp.MustCompile(`release/([0-9]+\.[0-9]+\.x)(\+ent)?`)

const (
	CERepo  = "hashicorp/vault"
	ENTRepo = "hashicorp/vault-enterprise"
)

// mergeCommitRetryWait is the duration to wait between merge commit SHA fetch
// retries. Exposed as a variable so tests can set it to zero.
var mergeCommitRetryWait = 60 * time.Second

// PullRequestEntry represents one element in the pull_requests JSON array
// accepted by the collect-artifact-shas workflow.
type PullRequestEntry struct {
	TargetBranch string `json:"target_branch"`
	URL          string `json:"url"`
}

// VersionSHAEntry is a single version → SHA pair in the output sha_map.
type VersionSHAEntry struct {
	Version string `json:"key"`
	SHA     string `json:"value"`
}

// set of release pull requests and build the accompanying Slack message.
type CollectArtifactShasReq struct {
	// PullRequests is the list of pull-request entries parsed from the workflow input.
	PullRequests []PullRequestEntry
	// If we're not in live mode, we want to at least generate some dummy shas to "go through the motions"
	DryRunShas   bool
	SlackChannel string
	// We want to pass our sorted_versions list in so that we can replace the branches from our release-build-ent output
	Versions []string
}

type CollectArtifactShasRes struct {
	// We take the merge shas corresponding to the PRs we passed in, and key them off of the release version
	// The versions make a much cleaner key than the release branch
	SHAMap   []VersionSHAEntry `json:"sha_map"`
	SlackMsg string            `json:"slack_msg"`
}

// Run executes the SHA-collection logic against the GitHub API.
func (r *CollectArtifactShasReq) Run(ctx context.Context, client *libgithub.Client) (*CollectArtifactShasRes, error) {
	ctx = slogctx.Append(
		ctx,
		slog.String("ce-repo", CERepo),
		slog.String("ent-repo", ENTRepo),
		slog.Bool("use-real-shas", !r.DryRunShas),
		slog.Int("pull-request-count", len(r.PullRequests)),
		slog.Int("version-count", len(r.Versions)),
	)
	slog.Default().DebugContext(ctx, "collecting artifact SHAs")

	if err := r.validate(); err != nil {
		return nil, fmt.Errorf("validating request: %w", err)
	}

	// let's take use the sorted versions to replace the release branch strings from release-build-ent
	branches := make([]string, len(r.PullRequests))
	for i, pr := range r.PullRequests {
		branches[i] = pr.TargetBranch
	}

	slog.Default().DebugContext(ctx, "mapping branches to versions")
	branchToVersion, err := mapBranchesToVersions(ctx, branches, r.Versions)
	if err != nil {
		return nil, fmt.Errorf("mapping branches to versions: %w", err)
	}

	// Using the PR number from our pr url, let's get the SHA from the merged commit used to create our artifact
	versionToSHA := make(map[string]string, len(r.PullRequests))
	for _, prEntry := range r.PullRequests {
		branch := prEntry.TargetBranch
		version, ok := branchToVersion[branch]
		if !ok {
			// mapBranchesToVersions already returns an error for unmapped branches;
			// this is just a safety net.
			continue
		}

		prCtx := slogctx.Append(
			ctx,
			slog.String("branch", branch),
			slog.String("version", version),
			slog.String("pr-url", prEntry.URL),
		)

		var sha string
		if r.DryRunShas {
			// Test mode: no call to github api, we just create clearly fake sha
			sha = fmt.Sprintf("fakeshaforbranch_%s", strings.ReplaceAll(branch, "/", "_"))
			slog.Default().DebugContext(prCtx, "test mode: generated dummy SHA", slog.String("sha", sha))
		} else {
			prNumber, err := parsePRNumber(prEntry.URL)
			if err != nil {
				return nil, fmt.Errorf("parsing PR number from %q: %w", prEntry.URL, err)
			}

			repo := repoForBranch(branch)
			owner, repoName, err := splitOwnerRepo(repo)
			if err != nil {
				return nil, err
			}

			prCtx = slogctx.Append(
				prCtx,
				slog.Int("pr-number", prNumber),
				slog.String("repo", repo),
			)
			slog.Default().DebugContext(prCtx, "fetching pull request details")

			pr, err := getPullRequest(ctx, client, owner, repoName, prNumber)
			if err != nil {
				return nil, fmt.Errorf("getting PR #%d from %s: %w", prNumber, repo, err)
			}

			// The merge commit SHA is only populated by GitHub after the merge event is
			// fully processed. If it is empty the PR may have just merged and the field
			// hasn't propagated yet — wait and retry before giving up.
			const mergeCommitRetries = 3
			for attempt := range mergeCommitRetries {
				sha = pr.GetMergeCommitSHA()
				if sha != "" {
					break
				}
				if attempt < mergeCommitRetries-1 {
					slog.Default().DebugContext(
						prCtx, "merge commit SHA not yet available, retrying",
						slog.Int("attempt", attempt+1),
						slog.Int("max-attempts", mergeCommitRetries),
						slog.Duration("wait", mergeCommitRetryWait),
					)
					time.Sleep(mergeCommitRetryWait)
					pr, err = getPullRequest(ctx, client, owner, repoName, prNumber)
					if err != nil {
						return nil, fmt.Errorf("getting PR #%d from %s on retry: %w", prNumber, repo, err)
					}
				}
			}
			if sha == "" {
				return nil, fmt.Errorf("merge commit SHA not available for PR #%d in %s after %d attempts", prNumber, repo, mergeCommitRetries)
			}
			slog.Default().DebugContext(prCtx, "resolved SHA", slog.String("sha", sha))
		}

		versionToSHA[version] = sha
	}

	// --- Step 3: build the ordered result and Slack message ---
	slog.Default().DebugContext(ctx, "building result and slack message")
	res := &CollectArtifactShasRes{}

	var sb strings.Builder
	if r.DryRunShas {
		sb.WriteString("*Merged PRs and SHAs (Test Mode):*\n")
	} else {
		sb.WriteString("*Merged PRs and SHAs:*\n")
	}

	// Iterate in the same order as the input versions to keep output stable.
	for _, version := range r.Versions {
		sha, ok := versionToSHA[version]
		if !ok {
			continue
		}
		res.SHAMap = append(res.SHAMap, VersionSHAEntry{Version: version, SHA: sha})
		sb.WriteString(fmt.Sprintf(":check_h2: %s → %s\n", version, sha))
		slog.Default().DebugContext(
			ctx, "mapped version to SHA",
			slog.String("version", version),
			slog.String("sha", sha),
		)
	}

	if len(res.SHAMap) == 0 {
		return nil, errors.New("no SHAs were collected: verify that the pull request URLs and versions provided are correct and correspond to one another")
	}

	sb.WriteString("\nPRs and SHAs pulled successfully :greenlightgo:\n")
	res.SlackMsg = sb.String()

	slog.Default().DebugContext(ctx, "artifact SHA collection complete", slog.Int("sha-map-count", len(res.SHAMap)))
	return res, nil
}

// validate ensures the request is complete before any work begins.
func (r *CollectArtifactShasReq) validate() error {
	if r == nil {
		return errors.New("failed to initialize request")
	}
	if len(r.PullRequests) == 0 {
		return errors.New("no pull requests provided")
	}
	if len(r.Versions) == 0 {
		return errors.New("no versions provided")
	}
	if r.SlackChannel == "" {
		return errors.New("no slack channel provided")
	}
	return nil
}

// mapBranchesToVersions builds a branch → version map by extracting the
// version pattern from each branch name and matching it against the supplied
// versions list.  +ent suffix consistency is enforced: an ENT branch must map
// to an ENT version and a CE branch must map to a CE version.
//
// Returns an error if any branch cannot be mapped.
func mapBranchesToVersions(ctx context.Context, branches, versions []string) (map[string]string, error) {
	result := make(map[string]string, len(branches))

	for _, branch := range branches {
		branchCtx := slogctx.Append(ctx, slog.String("branch", branch))

		m := branchVersionPattern.FindStringSubmatch(branch)
		if m == nil {
			return nil, fmt.Errorf("could not extract version pattern from branch %q", branch)
		}
		// m[1] = version pattern (e.g. "2.0.x")
		// m[2] = "+ent" or ""
		versionPattern := m[1]
		branchHasEnt := m[2] == "+ent"
		slog.Default().DebugContext(
			branchCtx, "extracted version pattern from branch",
			slog.String("version-pattern", versionPattern),
			slog.Bool("has-ent", branchHasEnt),
		)

		matched := false
		for _, version := range versions {
			versionHasEnt := strings.HasSuffix(version, "+ent")
			if branchHasEnt != versionHasEnt {
				continue
			}

			// Match on major.minor only: strip the .x suffix from the branch pattern,
			// then strip everything from the second dot onward in the version
			// (handles both plain patch versions like "2.0.3" and RC versions like "2.1.0-rc1").
			branchMM := strings.TrimSuffix(versionPattern, ".x")
			versionBase := strings.TrimSuffix(version, "+ent")
			secondDot := strings.Index(versionBase[strings.Index(versionBase, ".")+1:], ".")
			if secondDot < 0 {
				continue
			}
			// secondDot is relative to the slice starting after the first dot; adjust to absolute index
			majorMinor := versionBase[:strings.Index(versionBase, ".")+1+secondDot]
			if branchMM == majorMinor {
				slog.Default().DebugContext(branchCtx, "matched .x branch to version", slog.String("version", version))
				result[branch] = version
				matched = true
				break
			}
		}

		if !matched {
			return nil, fmt.Errorf("could not find matching version for branch %q (pattern: %s)", branch, versionPattern)
		}
	}

	return result, nil
}

func repoForBranch(branch string) string {
	if strings.HasSuffix(branch, "+ent") {
		return ENTRepo
	}
	return CERepo
}

// parsePRNumber extracts the integer PR number from a GitHub pull-request URL
// of the form https://github.com/owner/repo/pull/NNN.
func parsePRNumber(url string) (int, error) {
	var prNumber int
	_, err := fmt.Sscanf(url[strings.LastIndex(url, "/")+1:], "%d", &prNumber)
	if err != nil {
		return 0, fmt.Errorf("expected a numeric PR number at the end of %q: %w", url, err)
	}
	return prNumber, nil
}

// splitOwnerRepo splits an "owner/repo" string into its two components.
func splitOwnerRepo(ownerRepo string) (string, string, error) {
	parts := strings.SplitN(ownerRepo, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("expected owner/repo format, got %q", ownerRepo)
	}
	return parts[0], parts[1], nil
}
