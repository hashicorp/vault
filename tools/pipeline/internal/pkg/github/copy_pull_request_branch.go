// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// copyBranchPrefix is the prefix for copied pull request branches. It is used
// both for copying pull requests and for closing copied pull requests.
const copyBranchPrefix = "copy"

// encodeCopyPullRequestBranch returns the branch name to use for our PR copy operation.
// The branch name is important because we trigger post-merge automation based
// on the branch prefix. We also encode origin PR information for use in
// post-merge operations.
//
// The format is: <copy-branch-prefix>/<owner>/<repo>/<pr>/<pr-branch-name>
// eg: copy/hashicorp/vault/99999/my-feature-branch
//
// This informs the post-merge close operations that the Pull Request was originally
// copied from hashicorp/vault's Pull Request number 99999 and that the branch
// name of the pull request is my-feature-branch.
//
// It's imporant that both the encodeCopyPullRequestBranch and decodeCopyPullRequestBranch
// stay in sync.
func encodeCopyPullRequestBranch(
	owner string,
	repo string,
	number uint,
	prBranch string,
) string {
	prNumber := strconv.Itoa(int(number))
	name := strings.Join([]string{
		copyBranchPrefix,
		owner,
		repo,
		prNumber,
		prBranch,
	}, "/")
	if len(name) > 250 {
		// Handle Githubs branch name max length
		name = name[:250]
	}

	return name
}

// decodeCopyPullRequestBranch returns the encoded origin PR details from the
// copied Pull Request branch name.
//
// The format must be the same as what is described in encodeCopyPullRequestBranch()
func decodeCopyPullRequestBranch(ref string) (
	owner string,
	repo string,
	number uint,
	branch string,
	err error,
) {
	if ref == "" {
		err = errors.New("no copy branch provided")
		return owner, repo, number, branch, err
	}

	if !strings.HasPrefix(ref, copyBranchPrefix) {
		err = fmt.Errorf("invalid copy branch: branch does not start with %s", copyBranchPrefix)
		return owner, repo, number, branch, err
	}

	// eg: copy/hashicorp/vault/99999/my-feature-branch
	parts := strings.SplitN(ref, "/", 5)
	if len(parts) < 5 {
		err = fmt.Errorf("invalid copy branch: expected 5 parts, got %d", len(parts))
		return owner, repo, number, branch, err
	}
	owner = parts[1]
	repo = parts[2]

	var signedNumber int
	signedNumber, err = strconv.Atoi(parts[3])
	if err != nil {
		err = fmt.Errorf("invalid copy branch: pull request number is not a number: %w", err)
		return owner, repo, number, branch, err
	}
	if signedNumber < 0 {
		err = fmt.Errorf("invalid copy branch: number must be positive, got %d", signedNumber)
		return owner, repo, number, branch, err
	}
	number = uint(signedNumber)

	branch = parts[4]

	return owner, repo, number, branch, err
}

func closingIssueRefEqual(a, b *ClosingIssueRef) bool {
	if a == nil && b == nil {
		return true
	}

	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}

	if a.URL != b.URL {
		return false
	}

	if a.Number != b.Number {
		return false
	}

	if a.Title != b.Title {
		return false
	}

	if a.Repository.NameWithOwner != b.Repository.NameWithOwner {
		return false
	}

	return true
}
