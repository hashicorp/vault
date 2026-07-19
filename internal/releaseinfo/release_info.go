// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package releaseinfo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/helper/useragent"
)

const (
	defaultGitHubContentsAPI = "https://api.github.com/repos/hashicorp/web-unified-docs/contents/content/vault/global/partials/important-changes/summary-tables"
	defaultHTTPTimeout       = 30 * time.Second
	cacheTTL                 = 1 * time.Hour
)

var (
	httpClient        = newReleaseInfoHTTPClient()
	githubContentsAPI = defaultGitHubContentsAPI

	cacheMu      sync.Mutex
	cachedResult []VersionInfo
	cacheExpiry  time.Time
)

func newReleaseInfoHTTPClient() *http.Client {
	client := cleanhttp.DefaultPooledClient()
	client.Timeout = defaultHTTPTimeout
	return client
}

// VersionInfo represents release information for a specific version
type VersionInfo struct {
	Version         string       `json:"version"`
	BreakingChanges []ChangeItem `json:"breaking_changes"`
	NewBehavior     []ChangeItem `json:"new_behavior"`
	KnownIssues     []IssueItem  `json:"known_issues"`
}

// ChangeItem represents a breaking change or new behavior
type ChangeItem struct {
	Edition         string `json:"edition"`
	Recommendations bool   `json:"recommendations"`
	Introduced      string `json:"introduced"`
	Change          string `json:"change"`
	Link            string `json:"link"`
}

// KnownIssue represents a known issue
type IssueItem struct {
	Found      string `json:"found"`
	Fixed      string `json:"fixed"`
	Workaround string `json:"workaround"`
	Edition    string `json:"edition"`
	Issue      string `json:"issue"`
	Link       string `json:"link"`
}

type githubFile struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	DownloadURL string `json:"download_url"`
}

// FetchReleaseInformation calls the GitHub API to retrieve and parse
// all .mdx summary table files from the web-unified-docs repository.
func FetchReleaseInformation() ([]VersionInfo, error) {
	ctx := context.Background()
	return FetchReleaseInformationWithContext(ctx)
}

// FetchReleaseInformationWithContext calls the GitHub API to retrieve and parse
// all .mdx summary table files from the web-unified-docs repository with context support.
func FetchReleaseInformationWithContext(ctx context.Context) ([]VersionInfo, error) {
	// Serve a still-valid cached result to avoid re-fetching from GitHub
	// on every call.
	cacheMu.Lock()
	if cachedResult != nil && time.Now().Before(cacheExpiry) {
		result := cachedResult
		cacheMu.Unlock()
		return result, nil
	}
	cacheMu.Unlock()

	// Step 1: list the directory
	files, err := listMDXFiles(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing mdx files: %w", err)
	}

	// Step 2: fetch and parse content from each file
	var versions []VersionInfo
	for _, f := range files {
		versionInfo, err := fetchAndParse(ctx, f)
		if err != nil {
			// Stop on cancellation/timeout; otherwise skip files that
			// don't match the expected format.
			if ctxErr := ctx.Err(); ctxErr != nil {
				return nil, ctxErr
			}
			continue
		}
		if versionInfo != nil {
			versions = append(versions, *versionInfo)
		}
	}

	// Cache the successful result so subsequent reads within the TTL skip GitHub.
	cacheMu.Lock()
	cachedResult = versions
	cacheExpiry = time.Now().Add(cacheTTL)
	cacheMu.Unlock()

	return versions, nil
}

func listMDXFiles(ctx context.Context) ([]githubFile, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, githubContentsAPI, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", useragent.String())

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d from GitHub API", resp.StatusCode)
	}

	var all []githubFile
	if err := json.NewDecoder(resp.Body).Decode(&all); err != nil {
		return nil, err
	}

	// Filter: files only, .mdx extension
	var mdxFiles []githubFile
	for _, f := range all {
		if f.Type == "file" && strings.HasSuffix(f.Name, ".mdx") {
			mdxFiles = append(mdxFiles, f)
		}
	}

	return mdxFiles, nil
}

func fetchAndParse(ctx context.Context, file githubFile) (*VersionInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, file.DownloadURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", useragent.String())

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading: %w", err)
	}

	return parseMDX(string(body), file.Name)
}

func parseMDX(content string, filename string) (*VersionInfo, error) {
	// Extract version from filename
	// Supports: "1-21-x.mdx", "1_16.mdx" -> "1.21", "1.16"
	// Also supports: "2xx.mdx" -> "2.x"
	var version string

	versionRegex := regexp.MustCompile(`(\d+)[-_](\d+)(?:[-_]x)?\.mdx`)
	matches := versionRegex.FindStringSubmatch(filename)
	if len(matches) >= 3 {
		version = fmt.Sprintf("%s.%s", matches[1], matches[2])
	} else {
		// Try pattern for "2xx.mdx" format
		xxRegex := regexp.MustCompile(`(\d+)xx\.mdx`)
		xxMatches := xxRegex.FindStringSubmatch(filename)
		if len(xxMatches) >= 2 {
			version = fmt.Sprintf("%s.x", xxMatches[1])
		} else {
			return nil, fmt.Errorf("could not extract version from filename: %s", filename)
		}
	}

	versionInfo := &VersionInfo{
		Version:         version,
		BreakingChanges: []ChangeItem{},
		NewBehavior:     []ChangeItem{},
		KnownIssues:     []IssueItem{},
	}

	// Parse tables from MDX content
	// Look for HTML table tags or markdown tables
	lines := strings.Split(content, "\n")

	var currentSection string
	var inTable bool
	var headers []string

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		// Detect section headers (markdown headings like "### Breaking changes")
		if strings.HasPrefix(line, "#") {
			lowerLine := strings.ToLower(line)
			if strings.Contains(lowerLine, "breaking changes") {
				currentSection = "breaking_changes"
				inTable = false
				continue
			} else if strings.Contains(lowerLine, "new behavior") {
				currentSection = "new_behavior"
				inTable = false
				continue
			} else if strings.Contains(lowerLine, "known issues") {
				currentSection = "known_issues"
				inTable = false
				continue
			}
		}

		// Skip empty lines
		if line == "" {
			continue
		}

		// Detect table start (header row)
		if strings.Contains(line, "|") && !inTable && currentSection != "" {
			// Check if this looks like a table header (not a separator line)
			if !strings.Contains(line, "---") {
				inTable = true
				headers = parseTableRow(line)
				// Skip separator line (next line should contain ---)
				if i+1 < len(lines) {
					nextLine := strings.TrimSpace(lines[i+1])
					if strings.Contains(nextLine, "---") {
						i++
					}
				}
				continue
			}
		}

		// Parse table rows
		if inTable && strings.Contains(line, "|") {
			// Skip separator lines
			if strings.Contains(line, "---") {
				continue
			}

			cells := parseTableRow(line)
			// Check if this is a data row (not empty)
			if len(cells) > 0 && cells[0] != "" {
				switch currentSection {
				case "breaking_changes", "new_behavior":
					item := parseChangeItem(headers, cells)
					if currentSection == "breaking_changes" {
						versionInfo.BreakingChanges = append(versionInfo.BreakingChanges, item)
					} else {
						versionInfo.NewBehavior = append(versionInfo.NewBehavior, item)
					}
				case "known_issues":
					issue := parseKnownIssue(headers, cells)
					versionInfo.KnownIssues = append(versionInfo.KnownIssues, issue)
				}
			}
		} else if inTable && !strings.Contains(line, "|") && line != "" {
			// End of table
			inTable = false
		}
	}

	return versionInfo, nil
}

func parseTableRow(line string) []string {
	// Remove leading/trailing pipes if present, then split
	line = strings.Trim(line, " \t")
	line = strings.Trim(line, "|")
	parts := strings.Split(line, "|")
	var cells []string
	for _, part := range parts {
		cells = append(cells, strings.TrimSpace(part))
	}
	return cells
}

func parseChangeItem(headers, cells []string) ChangeItem {
	item := ChangeItem{
		Edition:         "All",
		Recommendations: true,
	}

	for i, header := range headers {
		if i >= len(cells) {
			break
		}
		value := cells[i]

		switch strings.ToLower(header) {
		case "edition":
			item.Edition = mapEdition(value)
		case "recommendations":
			// Check for "Yes" or "**Yes**"
			cleanValue := strings.ToLower(strings.Trim(value, "*"))
			item.Recommendations = cleanValue == "yes" || cleanValue == "true"
		case "introduced":
			item.Introduced = value
		case "change", "description":
			// Extract both text and link from the markdown link format
			item.Change = extractText(value)
			item.Link = extractLink(value)
		}
	}

	return item
}

func parseKnownIssue(headers, cells []string) IssueItem {
	issue := IssueItem{
		Edition: "All",
	}

	for i, header := range headers {
		if i >= len(cells) {
			break
		}
		value := cells[i]

		switch strings.ToLower(header) {
		case "found":
			issue.Found = value
		case "fixed":
			issue.Fixed = value
		case "workaround":
			// Check for "Yes" or "**Yes**"
			cleanValue := strings.ToLower(strings.Trim(value, "*"))
			if cleanValue == "yes" || cleanValue == "no" {
				issue.Workaround = cleanValue
			} else {
				issue.Workaround = value
			}
		case "edition":
			issue.Edition = mapEdition(value)
		case "issue", "description":
			// Extract both text and link from the markdown link format
			issue.Issue = extractText(value)
			issue.Link = extractLink(value)
		}
	}

	return issue
}

func extractText(mdx string) string {
	// Remove markdown link syntax: [text](url) -> text
	linkRegex := regexp.MustCompile(`\[([^\]]+)\]\([^\)]+\)`)
	text := linkRegex.ReplaceAllString(mdx, "$1")
	return strings.TrimSpace(text)
}

func extractLink(mdx string) string {
	// Extract URL from markdown link: [text](url) -> url
	linkRegex := regexp.MustCompile(`\[([^\]]+)\]\(([^\)]+)\)`)
	matches := linkRegex.FindStringSubmatch(mdx)
	if len(matches) > 2 {
		return matches[2]
	}
	// If no markdown link, return as-is (might be plain URL)
	return strings.TrimSpace(mdx)
}

func mapEdition(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	switch value {
	case "enterprise", "ent":
		return "Enterprise"
	case "ce", "community":
		return "CE"
	default:
		return "All"
	}
}
