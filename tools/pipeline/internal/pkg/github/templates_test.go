// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"errors"
	"io"
	"testing"

	libgithub "github.com/google/go-github/v81/github"
	"github.com/stretchr/testify/require"
)

func Test_renderEmbeddedTemplate_backportPRMessage(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		expectContains    []string
		expectNotContains []string
		origin            *libgithub.PullRequest
		attempt           *CreateBackportAttempt
	}{
		"no error": {
			expectContains:    []string{"original body"},
			expectNotContains: []string{"error body"},
			origin: &libgithub.PullRequest{
				Body:     libgithub.Ptr("original body"),
				Number:   libgithub.Ptr(1234),
				HTMLURL:  libgithub.Ptr("https://github.com/hashicorp/vault-enterprise/pull/1234"),
				MergedBy: &libgithub.User{Login: libgithub.Ptr("my-login")},
			},
			attempt: &CreateBackportAttempt{
				TargetRef: "release/1.19.x",
			},
		},
		"error": {
			expectContains: []string{"original body", "error body"},
			origin: &libgithub.PullRequest{
				Body:     libgithub.Ptr("original body"),
				Number:   libgithub.Ptr(1234),
				HTMLURL:  libgithub.Ptr("https://github.com/hashicorp/vault-enterprise/pull/1234"),
				MergedBy: &libgithub.User{Login: libgithub.Ptr("my-login")},
			},
			attempt: &CreateBackportAttempt{
				TargetRef: "release/1.19.x",
				Error:     errors.New("error body"),
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, err := renderEmbeddedTemplate("backport-pr-message.tmpl", struct {
				OriginPullRequest *libgithub.PullRequest
				Attempt           *CreateBackportAttempt
			}{test.origin, test.attempt})
			require.NoError(t, err)
			for _, c := range test.expectContains {
				require.Containsf(t, got, c, got)
			}
			for _, nc := range test.expectNotContains {
				require.NotContainsf(t, got, nc, got)
			}
		})
	}
}

func Test_renderEmbeddedTemplate_copyPRMessage(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		expectContains    []string
		expectNotContains []string
		origin            *libgithub.PullRequest
		error             error
		targetRef         string
	}{
		"no error": {
			expectContains:    []string{"original body"},
			expectNotContains: []string{"error body"},
			origin: &libgithub.PullRequest{
				Body:     libgithub.Ptr("original body"),
				Number:   libgithub.Ptr(1234),
				HTMLURL:  libgithub.Ptr("https://github.com/hashicorp/vault-enterprise/pull/1234"),
				MergedBy: &libgithub.User{Login: libgithub.Ptr("my-login")},
			},
			targetRef: "release/1.19.x",
			error:     nil,
		},
		"error": {
			expectContains: []string{"original body", "error body"},
			origin: &libgithub.PullRequest{
				Body:     libgithub.Ptr("original body"),
				Number:   libgithub.Ptr(1234),
				HTMLURL:  libgithub.Ptr("https://github.com/hashicorp/vault-enterprise/pull/1234"),
				MergedBy: &libgithub.User{Login: libgithub.Ptr("my-login")},
			},
			targetRef: "release/1.19.x",
			error:     errors.New("error body"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, err := renderEmbeddedTemplate("copy-pr-message.tmpl", struct {
				OriginPullRequest *libgithub.PullRequest
				TargetRef         string
				Error             error
			}{
				test.origin,
				test.targetRef,
				test.error,
			})
			require.NoError(t, err)
			for _, c := range test.expectContains {
				require.Containsf(t, got, c, got)
			}
			for _, nc := range test.expectNotContains {
				require.NotContainsf(t, got, nc, got)
			}
		})
	}
}

func Test_renderEmbeddedTemplateToTmpFile_copyPRComment(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		coAuthoredByTrailers []string
		originPullRequest    *libgithub.PullRequest
		targetRef            string
	}{
		"Co-Authored-By": {
			coAuthoredByTrailers: []string{
				"Co-Authored-By: Jane Doe <jane@example.com>",
				"Co-Authored-By: John Doe <john@example.com>",
			},
			originPullRequest: &libgithub.PullRequest{
				Body:     libgithub.Ptr("original body"),
				Number:   libgithub.Ptr(1234),
				HTMLURL:  libgithub.Ptr("https://github.com/hashicorp/vault-enterprise/pull/1234"),
				MergedBy: &libgithub.User{Login: libgithub.Ptr("my-login")},
			},
			targetRef: "release/1.19.x+ent",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			file, err := renderEmbeddedTemplateToTmpFile("copy-pr-commit-message.tmpl", struct {
				CoAuthoredByTrailers []string
				OriginPullRequest    *libgithub.PullRequest
				TargetRef            string
			}{
				test.coAuthoredByTrailers,
				test.originPullRequest,
				test.targetRef,
			})

			require.NoError(t, err)
			defer file.Close()
			bytes, err := io.ReadAll(file)
			require.NoError(t, err)
			require.NotEmpty(t, bytes)
			for _, c := range test.coAuthoredByTrailers {
				require.Contains(t, string(bytes), c)
			}
		})
	}
}
