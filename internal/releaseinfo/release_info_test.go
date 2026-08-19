// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package releaseinfo

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// mockGitHubServer creates a test HTTP server that simulates GitHub API responses
func mockGitHubServer(t *testing.T, files []githubFile, fileContents map[string]string) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock the directory listing endpoint
		if r.URL.Path == "/repos/hashicorp/web-unified-docs/contents/content/vault/global/partials/important-changes/summary-tables" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(files)
			return
		}

		// Mock file download endpoints
		for filename, content := range fileContents {
			if r.URL.Path == "/download/"+filename {
				w.Header().Set("Content-Type", "text/plain")
				w.Write([]byte(content))
				return
			}
		}

		http.NotFound(w, r)
	}))
}

// TestFetchReleaseInformation tests the FetchReleaseInformation function with mock GitHub API responses.
func TestFetchReleaseInformation(t *testing.T) {
	// Save original values
	originalClient := httpClient
	originalAPI := githubContentsAPI

	// Restore after test
	defer func() {
		httpClient = originalClient
		githubContentsAPI = originalAPI
		cacheMu.Lock()
		cachedResult = nil
		cacheExpiry = time.Time{}
		cacheMu.Unlock()
	}()

	// Create mock files
	mockFiles := []githubFile{
		{
			Name:        "1-21-x.mdx",
			Type:        "file",
			DownloadURL: "",
		},
		{
			Name:        "1_16.mdx",
			Type:        "file",
			DownloadURL: "",
		},
	}

	mockContent121 := `# Vault 1.21.x

## Breaking changes

| Edition | Recommendations | Introduced | Change |
|---------|----------------|------------|--------|
| Enterprise | **Yes** | 1.21.0 | [Database secrets engine change](https://example.com/change1) |
| All | Yes | 1.21.1 | [API endpoint deprecated](https://example.com/change2) |

## New behavior

| Edition | Recommendations | Introduced | Change |
|---------|----------------|------------|--------|
| CE | No | 1.21.0 | [New feature added](https://example.com/feature1) |

## Known issues

| Found | Fixed | Workaround | Edition | Issue |
|-------|-------|------------|---------|-------|
| 1.21.0 | 1.21.2 | Yes | All | [Memory leak in plugin](https://example.com/issue1) |
`

	mockContent116 := `# Vault 1.16

## Breaking changes

| Edition | Recommendations | Introduced | Change |
|---------|----------------|------------|--------|
| All | **Yes** | 1.16.0 | [Auth method removed](https://example.com/change3) |
`

	fileContents := map[string]string{
		"1-21-x.mdx": mockContent121,
		"1_16.mdx":   mockContent116,
	}

	// Create mock server
	server := mockGitHubServer(t, mockFiles, fileContents)
	defer server.Close()

	// Update mock files with server URLs
	mockFiles[0].DownloadURL = server.URL + "/download/1-21-x.mdx"
	mockFiles[1].DownloadURL = server.URL + "/download/1_16.mdx"

	// Override global variables
	githubContentsAPI = server.URL + "/repos/hashicorp/web-unified-docs/contents/content/vault/global/partials/important-changes/summary-tables"
	httpClient = &http.Client{Timeout: 5 * time.Second}

	// Test the function
	versions, err := FetchReleaseInformation()
	if err != nil {
		t.Fatalf("FetchReleaseInformation() error = %v", err)
	}

	if len(versions) != 2 {
		t.Errorf("Expected 2 versions, got %d", len(versions))
	}

	// Verify version 1.21
	var v121 *VersionInfo
	for i := range versions {
		if versions[i].Version == "1.21" {
			v121 = &versions[i]
			break
		}
	}

	if v121 == nil {
		t.Fatal("Version 1.21 not found")
	}

	if len(v121.BreakingChanges) != 2 {
		t.Errorf("Expected 2 breaking changes for 1.21, got %d", len(v121.BreakingChanges))
	}

	if len(v121.NewBehavior) != 1 {
		t.Errorf("Expected 1 new behavior for 1.21, got %d", len(v121.NewBehavior))
	}

	if len(v121.KnownIssues) != 1 {
		t.Errorf("Expected 1 known issue for 1.21, got %d", len(v121.KnownIssues))
	}

	// Verify first breaking change
	bc := v121.BreakingChanges[0]
	if bc.Edition != "Enterprise" {
		t.Errorf("Expected edition 'Enterprise', got '%s'", bc.Edition)
	}
	if !bc.Recommendations {
		t.Error("Expected recommendations to be true")
	}
	if bc.Introduced != "1.21.0" {
		t.Errorf("Expected introduced '1.21.0', got '%s'", bc.Introduced)
	}
	if bc.Change != "Database secrets engine change" {
		t.Errorf("Expected change text, got '%s'", bc.Change)
	}
	if bc.Link != "https://example.com/change1" {
		t.Errorf("Expected link, got '%s'", bc.Link)
	}
}

// TestFetchReleaseInformationWithContext tests the FetchReleaseInformationWithContext function with a valid context.
func TestFetchReleaseInformationWithContext(t *testing.T) {
	// Save original values
	originalClient := httpClient
	originalAPI := githubContentsAPI

	defer func() {
		httpClient = originalClient
		githubContentsAPI = originalAPI
		cacheMu.Lock()
		cachedResult = nil
		cacheExpiry = time.Time{}
		cacheMu.Unlock()
	}()

	mockFiles := []githubFile{
		{
			Name:        "1-21-x.mdx",
			Type:        "file",
			DownloadURL: "",
		},
	}

	mockContent := `# Vault 1.21.x

## Breaking changes

| Edition | Recommendations | Introduced | Change |
|---------|----------------|------------|--------|
| All | Yes | 1.21.0 | Test change |
`

	fileContents := map[string]string{
		"1-21-x.mdx": mockContent,
	}

	server := mockGitHubServer(t, mockFiles, fileContents)
	defer server.Close()

	mockFiles[0].DownloadURL = server.URL + "/download/1-21-x.mdx"
	githubContentsAPI = server.URL + "/repos/hashicorp/web-unified-docs/contents/content/vault/global/partials/important-changes/summary-tables"
	httpClient = &http.Client{Timeout: 5 * time.Second}

	// Test with context
	ctx := context.Background()
	versions, err := FetchReleaseInformationWithContext(ctx)
	if err != nil {
		t.Fatalf("FetchReleaseInformationWithContext() error = %v", err)
	}

	if len(versions) != 1 {
		t.Errorf("Expected 1 version, got %d", len(versions))
	}
}

// TestFetchReleaseInformationWithContextCancellation tests that FetchReleaseInformationWithContext properly handles context cancellation.
func TestFetchReleaseInformationWithContextCancellation(t *testing.T) {
	// Save original values
	originalClient := httpClient
	originalAPI := githubContentsAPI

	defer func() {
		httpClient = originalClient
		githubContentsAPI = originalAPI
	}()

	// Create a server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	githubContentsAPI = server.URL + "/test"
	httpClient = &http.Client{Timeout: 5 * time.Second}

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := FetchReleaseInformationWithContext(ctx)
	if err == nil {
		t.Error("Expected error with cancelled context, got nil")
	}
}

// TestListMDXFiles tests the listMDXFiles function to ensure it correctly filters and returns only .mdx files.
func TestListMDXFiles(t *testing.T) {
	originalClient := httpClient
	originalAPI := githubContentsAPI

	defer func() {
		httpClient = originalClient
		githubContentsAPI = originalAPI
	}()

	mockFiles := []githubFile{
		{Name: "1-21-x.mdx", Type: "file"},
		{Name: "1_16.mdx", Type: "file"},
		{Name: "README.md", Type: "file"}, // Should be filtered out
		{Name: "subdir", Type: "dir"},     // Should be filtered out
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockFiles)
	}))
	defer server.Close()

	githubContentsAPI = server.URL
	httpClient = &http.Client{Timeout: 5 * time.Second}

	files, err := listMDXFiles(context.Background())
	if err != nil {
		t.Fatalf("listMDXFiles() error = %v", err)
	}

	if len(files) != 2 {
		t.Errorf("Expected 2 .mdx files, got %d", len(files))
	}

	for _, f := range files {
		if f.Type != "file" {
			t.Errorf("Expected type 'file', got '%s'", f.Type)
		}
		if !hasSuffix(f.Name, ".mdx") {
			t.Errorf("Expected .mdx extension, got '%s'", f.Name)
		}
	}
}

func hasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

// TestListMDXFilesError tests that listMDXFiles properly handles HTTP errors from the GitHub API.
func TestListMDXFilesError(t *testing.T) {
	originalClient := httpClient
	originalAPI := githubContentsAPI

	defer func() {
		httpClient = originalClient
		githubContentsAPI = originalAPI
	}()

	// Server returns error status
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	githubContentsAPI = server.URL
	httpClient = &http.Client{Timeout: 5 * time.Second}

	_, err := listMDXFiles(context.Background())
	if err == nil {
		t.Error("Expected error for non-200 status, got nil")
	}
}

// TestParseMDX tests the parseMDX function with various MDX content formats and edge cases.
func TestParseMDX(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		filename string
		want     *VersionInfo
		wantErr  bool
	}{
		{
			name:     "version from 1-21-x.mdx",
			filename: "1-21-x.mdx",
			content: `# Vault 1.21.x

## Breaking changes

| Edition | Recommendations | Introduced | Change |
|---------|----------------|------------|--------|
| Enterprise | **Yes** | 1.21.0 | [Test change](https://example.com) |
`,
			want: &VersionInfo{
				Version: "1.21",
				BreakingChanges: []ChangeItem{
					{
						Edition:         "Enterprise",
						Recommendations: true,
						Introduced:      "1.21.0",
						Change:          "Test change",
						Link:            "https://example.com",
					},
				},
				NewBehavior: []ChangeItem{},
				KnownIssues: []IssueItem{},
			},
			wantErr: false,
		},
		{
			name:     "version from 1_16.mdx",
			filename: "1_16.mdx",
			content:  `# Vault 1.16`,
			want: &VersionInfo{
				Version:         "1.16",
				BreakingChanges: []ChangeItem{},
				NewBehavior:     []ChangeItem{},
				KnownIssues:     []IssueItem{},
			},
			wantErr: false,
		},
		{
			name:     "version from 2xx.mdx",
			filename: "2xx.mdx",
			content:  `# Vault 2.x`,
			want: &VersionInfo{
				Version:         "2.x",
				BreakingChanges: []ChangeItem{},
				NewBehavior:     []ChangeItem{},
				KnownIssues:     []IssueItem{},
			},
			wantErr: false,
		},
		{
			name:     "invalid filename",
			filename: "invalid.mdx",
			content:  `# Vault`,
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "new behavior section",
			filename: "1-21-x.mdx",
			content: `# Vault 1.21.x

## New behavior

| Edition | Recommendations | Introduced | Change |
|---------|----------------|------------|--------|
| CE | No | 1.21.0 | New feature |
`,
			want: &VersionInfo{
				Version:         "1.21",
				BreakingChanges: []ChangeItem{},
				NewBehavior: []ChangeItem{
					{
						Edition:         "CE",
						Recommendations: false,
						Introduced:      "1.21.0",
						Change:          "New feature",
						Link:            "New feature",
					},
				},
				KnownIssues: []IssueItem{},
			},
			wantErr: false,
		},
		{
			name:     "known issues section",
			filename: "1-21-x.mdx",
			content: `# Vault 1.21.x

## Known issues

| Found | Fixed | Workaround | Edition | Issue |
|-------|-------|------------|---------|-------|
| 1.21.0 | 1.21.2 | Yes | All | [Memory leak](https://example.com/issue) |
`,
			want: &VersionInfo{
				Version:         "1.21",
				BreakingChanges: []ChangeItem{},
				NewBehavior:     []ChangeItem{},
				KnownIssues: []IssueItem{
					{
						Found:      "1.21.0",
						Fixed:      "1.21.2",
						Workaround: "yes",
						Edition:    "All",
						Issue:      "Memory leak",
						Link:       "https://example.com/issue",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseMDX(tt.content, tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseMDX() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if got.Version != tt.want.Version {
				t.Errorf("Version = %v, want %v", got.Version, tt.want.Version)
			}
			if len(got.BreakingChanges) != len(tt.want.BreakingChanges) {
				t.Errorf("BreakingChanges count = %v, want %v", len(got.BreakingChanges), len(tt.want.BreakingChanges))
			}
			if len(got.NewBehavior) != len(tt.want.NewBehavior) {
				t.Errorf("NewBehavior count = %v, want %v", len(got.NewBehavior), len(tt.want.NewBehavior))
			}
			if len(got.KnownIssues) != len(tt.want.KnownIssues) {
				t.Errorf("KnownIssues count = %v, want %v", len(got.KnownIssues), len(tt.want.KnownIssues))
			}
		})
	}
}

// TestParseTableRow tests the parseTableRow function to ensure it correctly splits markdown table rows.
func TestParseTableRow(t *testing.T) {
	tests := []struct {
		name string
		line string
		want []string
	}{
		{
			name: "basic row",
			line: "| col1 | col2 | col3 |",
			want: []string{"col1", "col2", "col3"},
		},
		{
			name: "row without leading/trailing pipes",
			line: "col1 | col2 | col3",
			want: []string{"col1", "col2", "col3"},
		},
		{
			name: "row with extra spaces",
			line: "|  col1  |  col2  |  col3  |",
			want: []string{"col1", "col2", "col3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseTableRow(tt.line)
			if len(got) != len(tt.want) {
				t.Errorf("parseTableRow() length = %v, want %v", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("parseTableRow()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

// TestExtractText tests the extractText function to ensure it correctly extracts text from markdown links.
func TestExtractText(t *testing.T) {
	tests := []struct {
		name string
		mdx  string
		want string
	}{
		{
			name: "markdown link",
			mdx:  "[Link text](https://example.com)",
			want: "Link text",
		},
		{
			name: "plain text",
			mdx:  "Plain text",
			want: "Plain text",
		},
		{
			name: "multiple links",
			mdx:  "[Link1](url1) and [Link2](url2)",
			want: "Link1 and Link2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractText(tt.mdx); got != tt.want {
				t.Errorf("extractText() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestExtractLink tests the extractLink function to ensure it correctly extracts URLs from markdown links.
func TestExtractLink(t *testing.T) {
	tests := []struct {
		name string
		mdx  string
		want string
	}{
		{
			name: "markdown link",
			mdx:  "[Link text](https://example.com)",
			want: "https://example.com",
		},
		{
			name: "plain text",
			mdx:  "https://example.com",
			want: "https://example.com",
		},
		{
			name: "no link",
			mdx:  "Plain text",
			want: "Plain text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractLink(tt.mdx); got != tt.want {
				t.Errorf("extractLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestMapEdition tests the mapEdition function to ensure it correctly normalizes edition values.
func TestMapEdition(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "enterprise", value: "enterprise", want: "Enterprise"},
		{name: "ent", value: "ent", want: "Enterprise"},
		{name: "Enterprise", value: "Enterprise", want: "Enterprise"},
		{name: "ce", value: "ce", want: "CE"},
		{name: "community", value: "community", want: "CE"},
		{name: "all", value: "all", want: "All"},
		{name: "unknown", value: "unknown", want: "All"},
		{name: "empty", value: "", want: "All"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mapEdition(tt.value); got != tt.want {
				t.Errorf("mapEdition() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestParseChangeItem tests the parseChangeItem function to ensure it correctly parses change items from table rows.
func TestParseChangeItem(t *testing.T) {
	tests := []struct {
		name    string
		headers []string
		cells   []string
		want    ChangeItem
	}{
		{
			name:    "full change item",
			headers: []string{"Edition", "Recommendations", "Introduced", "Change"},
			cells:   []string{"Enterprise", "**Yes**", "1.21.0", "[Test change](https://example.com)"},
			want: ChangeItem{
				Edition:         "Enterprise",
				Recommendations: true,
				Introduced:      "1.21.0",
				Change:          "Test change",
				Link:            "https://example.com",
			},
		},
		{
			name:    "recommendations no",
			headers: []string{"Edition", "Recommendations", "Introduced", "Change"},
			cells:   []string{"All", "No", "1.21.0", "Test"},
			want: ChangeItem{
				Edition:         "All",
				Recommendations: false,
				Introduced:      "1.21.0",
				Change:          "Test",
				Link:            "Test",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseChangeItem(tt.headers, tt.cells)
			if got.Edition != tt.want.Edition {
				t.Errorf("Edition = %v, want %v", got.Edition, tt.want.Edition)
			}
			if got.Recommendations != tt.want.Recommendations {
				t.Errorf("Recommendations = %v, want %v", got.Recommendations, tt.want.Recommendations)
			}
			if got.Introduced != tt.want.Introduced {
				t.Errorf("Introduced = %v, want %v", got.Introduced, tt.want.Introduced)
			}
			if got.Change != tt.want.Change {
				t.Errorf("Change = %v, want %v", got.Change, tt.want.Change)
			}
		})
	}
}

// TestParseKnownIssue tests the parseKnownIssue function to ensure it correctly parses known issues from table rows.
func TestParseKnownIssue(t *testing.T) {
	tests := []struct {
		name    string
		headers []string
		cells   []string
		want    IssueItem
	}{
		{
			name:    "full issue",
			headers: []string{"Found", "Fixed", "Workaround", "Edition", "Issue"},
			cells:   []string{"1.21.0", "1.21.2", "Yes", "All", "[Memory leak](https://example.com)"},
			want: IssueItem{
				Found:      "1.21.0",
				Fixed:      "1.21.2",
				Workaround: "yes",
				Edition:    "All",
				Issue:      "Memory leak",
				Link:       "https://example.com",
			},
		},
		{
			name:    "workaround no",
			headers: []string{"Found", "Fixed", "Workaround", "Edition", "Issue"},
			cells:   []string{"1.21.0", "1.21.2", "No", "Enterprise", "Test issue"},
			want: IssueItem{
				Found:      "1.21.0",
				Fixed:      "1.21.2",
				Workaround: "no",
				Edition:    "Enterprise",
				Issue:      "Test issue",
				Link:       "Test issue",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseKnownIssue(tt.headers, tt.cells)
			if got.Found != tt.want.Found {
				t.Errorf("Found = %v, want %v", got.Found, tt.want.Found)
			}
			if got.Fixed != tt.want.Fixed {
				t.Errorf("Fixed = %v, want %v", got.Fixed, tt.want.Fixed)
			}
			if got.Workaround != tt.want.Workaround {
				t.Errorf("Workaround = %v, want %v", got.Workaround, tt.want.Workaround)
			}
			if got.Edition != tt.want.Edition {
				t.Errorf("Edition = %v, want %v", got.Edition, tt.want.Edition)
			}
		})
	}
}
