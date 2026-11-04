// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package releases

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"testing"
)

func TestNewListBinaryVersionsReq(t *testing.T) {
	// Verify that the constructor returns a valid object
	req := NewListBinaryVersionsReq("1.2.3 4.5.6-ent")
	if req == nil {
		t.Fatal("returned nil")
	}

	// Ensure the input string is stored without modification
	if req.VersionsString != "1.2.3 4.5.6-ent" {
		t.Errorf("got %q", req.VersionsString)
	}
}

// Unit test: JSON output and String() summary formatting
func TestListBinaryVersionsRes_ToJSON_and_String(t *testing.T) {
	// Create a fake response object
	res := &ListBinaryVersionsRes{
		ValidVersions: map[string]struct {
			Status   string        `json:"status"`
			Variants []VariantInfo `json:"variants"`
		}{
			"1.17.0-ent": {
				Status: "valid",
				Variants: []VariantInfo{
					{Variant: "1.17.0+ent", OS: []string{"linux", "darwin"}},
					{Variant: "1.17.0+ent.hsm", OS: []string{"linux"}},
				},
			},
		},
		InvalidVersions: []string{"9.9.9"},
		AllVersions:     []string{"1.17.0-ent", "9.9.9"},
	}

	// Test: ToJSON output
	b, err := res.ToJSON()
	if err != nil {
		t.Fatal(err)
	}

	// Ensure the JSON contains known variant values
	if !strings.Contains(string(b), "1.17.0+ent.hsm") {
		t.Error("JSON missing variant")
	}

	// -------------------------
	// Test: String() summary
	// -------------------------
	// Expected values based on the test fixture:
	//   total requested: 2
	//   valid versions: 1
	//   total variants: 2
	//   missing: 1
	expected := "2 → 1 valid (2 variants), 1 missing"

	// Check that the summary contains the correct computed text.
	if !strings.Contains(res.String(), expected) {
		t.Errorf("String() = %q, want substring %q", res.String(), expected)
	}
}

type fakeReleasesTransport struct {
	pages map[string]string
}

func (f *fakeReleasesTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	body, ok := f.pages[req.URL.String()]
	code := http.StatusOK
	if !ok {
		code = http.StatusNotFound
	}
	return &http.Response{
		StatusCode: code,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
		Request:    req,
	}, nil
}

func releasesListing(links ...string) string {
	var sb strings.Builder
	sb.WriteString("<html><body><ul>")
	for _, l := range links {
		fmt.Fprintf(&sb, `<li><a href="%s">%s</a></li>`, l, l)
	}
	sb.WriteString("</ul></body></html>")
	return sb.String()
}

func TestListBinaryVersionsReq_EntLabelReturnsAllEnterpriseVariants(t *testing.T) {
	zip := func(variant string, osName string) string {
		return fmt.Sprintf("vault_%s_%s_amd64.zip", variant, osName)
	}

	pages := map[string]string{
		baseURL: releasesListing(
			"/vault/1.21.13+ent/",
			"/vault/1.21.13+ent.fips1403/",
			"/vault/1.21.13+ent.hsm/",
			"/vault/1.21.13+ent.hsm.fips1403/",
			"/vault/1.21.130+ent/",
		),
		baseURL + "1.21.13+ent/":              releasesListing(zip("1.21.13+ent", "linux"), zip("1.21.13+ent", "windows")),
		baseURL + "1.21.13+ent.fips1403/":     releasesListing(zip("1.21.13+ent.fips1403", "linux")),
		baseURL + "1.21.13+ent.hsm/":          releasesListing(zip("1.21.13+ent.hsm", "linux")),
		baseURL + "1.21.13+ent.hsm.fips1403/": releasesListing(zip("1.21.13+ent.hsm.fips1403", "linux")),
	}

	oldClient := http.DefaultClient
	http.DefaultClient = &http.Client{Transport: &fakeReleasesTransport{pages: pages}}
	t.Cleanup(func() { http.DefaultClient = oldClient })

	res, err := NewListBinaryVersionsReq("1.21.13-ent 1.19.99-ent").Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entry, ok := res.ValidVersions["1.21.13-ent"]
	if !ok {
		t.Fatalf("expected 1.21.13-ent to be valid, got valid=%v invalid=%v", res.ValidVersions, res.InvalidVersions)
	}

	var got []string
	for _, v := range entry.Variants {
		got = append(got, v.Variant)
	}
	slices.Sort(got)
	want := []string{"1.21.13+ent", "1.21.13+ent.fips1403", "1.21.13+ent.hsm", "1.21.13+ent.hsm.fips1403"}
	if !slices.Equal(got, want) {
		t.Errorf("expected variants %v, got %v", want, got)
	}

	if !slices.Contains(res.InvalidVersions, "1.19.99-ent") {
		t.Errorf("expected 1.19.99-ent to be invalid, got %v", res.InvalidVersions)
	}
}
