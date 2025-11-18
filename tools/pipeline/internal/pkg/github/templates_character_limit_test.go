// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_limitCharacters(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected func(string) bool
	}{
		{
			name:    "short content is unchanged",
			content: "This is a short message",
			expected: func(result string) bool {
				return result == "This is a short message"
			},
		},
		{
			name:    "content at limit is unchanged",
			content: strings.Repeat("a", maxGitHubMessageChars),
			expected: func(result string) bool {
				return len(result) == maxGitHubMessageChars && strings.HasPrefix(result, "aaa")
			},
		},
		{
			name:    "content over limit is truncated",
			content: strings.Repeat("a", maxGitHubMessageChars+1000),
			expected: func(result string) bool {
				return len(result) <= maxGitHubMessageChars &&
					strings.Contains(result, "Message truncated due to GitHub's character limit")
			},
		},
		{
			name:    "markdown content with newlines truncates at newline",
			content: strings.Repeat("line\n", maxGitHubMessageChars/5+100),
			expected: func(result string) bool {
				return len(result) <= maxGitHubMessageChars &&
					strings.Contains(result, "Message truncated due to GitHub's character limit") &&
					strings.HasPrefix(result, "line\n")
			},
		},
		{
			name:    "very long line without newlines gets hard truncated",
			content: strings.Repeat("a", maxGitHubMessageChars+100),
			expected: func(result string) bool {
				return len(result) <= maxGitHubMessageChars &&
					strings.Contains(result, "Message truncated due to GitHub's character limit")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := limitCharacters(tt.content)
			assert.True(t, tt.expected(result), "Result does not match expected criteria")
			assert.LessOrEqual(t, len(result), maxGitHubMessageChars, "Result exceeds maximum character limit")
		})
	}
}

func Test_limitCharacters_preserves_original_under_limit(t *testing.T) {
	content := "# PR Title\n\nThis is the body content.\n\n## Details\n\nSome details here."
	result := limitCharacters(content)
	assert.Equal(t, content, result, "Short content should be unchanged")
}

func Test_limitCharacters_adds_truncation_notice(t *testing.T) {
	longContent := strings.Repeat("This is a very long line that will exceed the GitHub character limit. ", maxGitHubMessageChars/70+10)
	result := limitCharacters(longContent)

	assert.LessOrEqual(t, len(result), maxGitHubMessageChars, "Result should not exceed max chars")
	assert.Contains(t, result, ":scissors: **Message truncated due to GitHub's character limit**", "Should contain truncation notice")
	assert.Contains(t, result, fmt.Sprintf("%d character limit", maxGitHubMessageChars), "Should mention the specific limit")
}
