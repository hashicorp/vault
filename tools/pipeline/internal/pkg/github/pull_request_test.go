// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_encodeCopyPullRequest_Roundtrip(t *testing.T) {
	t.Parallel()

	for desc, test := range map[string]struct {
		owner            string
		repo             string
		number           uint
		prBranch         string
		expectedPRBranch string
	}{
		"standard": {
			owner:            "hashicorp",
			repo:             "vault",
			number:           49689,
			prBranch:         "my-feature-branch",
			expectedPRBranch: "my-feature-branch",
		},
		"super long": {
			owner:    "hashicorp",
			repo:     "vault",
			number:   49689,
			prBranch: "Lorem-ipsum-dolor-sit-amet--consectetur-adipiscing-elit--Praesent-accumsan-metus-sed-vehicula-accumsan--Nunc-semper-vehicula-tempor--Vestibulum-fringilla-enim-nec-ipsum-tincidunt-viverra--Etiam-iaculis-metus-ultricies-risus-rutrum--et-lobortis-orci-al",
			// truncated the branch name as we'll hit our max char limit
			expectedPRBranch: "Lorem-ipsum-dolor-sit-amet--consectetur-adipiscing-elit--Praesent-accumsan-metus-sed-vehicula-accumsan--Nunc-semper-vehicula-tempor--Vestibulum-fringilla-enim-nec-ipsum-tincidunt-viverra--Etiam-iaculis-metus-ultricies-risus",
		},
	} {
		t.Run(desc, func(t *testing.T) {
			t.Parallel()

			encodedBranch := encodeCopyPullRequestBranch(test.owner, test.repo, test.number, test.prBranch)
			owner, repo, number, branch, err := decodeCopyPullRequestBranch(encodedBranch)
			require.NoError(t, err)
			require.Equal(t, test.owner, owner)
			require.Equal(t, test.repo, repo)
			require.Equal(t, test.number, number)
			require.Equal(t, test.expectedPRBranch, branch)
		})
	}
}
