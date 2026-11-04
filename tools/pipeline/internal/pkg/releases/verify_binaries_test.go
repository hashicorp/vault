// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package releases

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

var (
	allStandardOS = []string{"darwin", "freebsd", "linux", "netbsd", "openbsd", "solaris", "windows"}
	linuxOnlyOS   = []string{"linux"}
)

// makeVersionData creates a ListBinaryVersionsRes with the specified version and variants.
func makeVersionData(version string, variants []VariantInfo) *ListBinaryVersionsRes {
	return &ListBinaryVersionsRes{
		AllVersions: []string{version},
		ValidVersions: map[string]struct {
			Status   string        `json:"status"`
			Variants []VariantInfo `json:"variants"`
		}{
			version: {
				Status:   "valid",
				Variants: variants,
			},
		},
		InvalidVersions: []string{},
	}
}

// makeValidBinaryData creates a fully valid ListBinaryVersionsRes fixture.
func makeValidBinaryData(version string) *ListBinaryVersionsRes {
	base := strings.TrimSuffix(version, "-lts")
	base = strings.TrimSuffix(base, "-ce")
	base = strings.TrimSuffix(base, "-ent")

	osList := RequiredOSList
	fipsSuffix := "fips1403"
	clean := strings.TrimPrefix(version, "v")
	parts := strings.Split(clean, ".")
	if len(parts) > 0 && parts[0] == "1" {
		fipsSuffix = "fips1402"
	}

	return makeVersionData(version, []VariantInfo{
		{Variant: base, OS: osList},
		{Variant: base + "+ent", OS: osList},
		{Variant: base + "+ent." + fipsSuffix, OS: linuxOnlyOS},
		{Variant: base + "+ent.hsm", OS: linuxOnlyOS},
		{Variant: base + "+ent.hsm." + fipsSuffix, OS: linuxOnlyOS},
	})
}

type mockFetcher struct {
	data *ListBinaryVersionsRes
	err  error
}

func (m *mockFetcher) Fetch(ctx context.Context, versions string) (*ListBinaryVersionsRes, error) {
	return m.data, m.err
}

func TestVerifyBinaries_PositiveAllPresent(t *testing.T) {
	data := makeValidBinaryData("2.2.1")
	req := &VerifyBinariesReq{
		VersionsString: "2.2.1",
		BinaryFetcher:  &mockFetcher{data: data},
	}

	res, err := req.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !res.Success {
		t.Errorf("expected Success=true, got false")
	}
	if res.HasMissing {
		t.Errorf("expected HasMissing=false, got true")
	}
	if len(res.MissingVersions) != 0 {
		t.Errorf("expected 0 missing versions, got %d", len(res.MissingVersions))
	}
	if len(res.MissingVariants) != 0 {
		t.Errorf("expected 0 missing variants, got %d", len(res.MissingVariants))
	}
	if len(res.MissingOS) != 0 {
		t.Errorf("expected 0 missing OS, got %d", len(res.MissingOS))
	}
	if len(res.InvalidOS) != 0 {
		t.Errorf("expected 0 invalid OS, got %d", len(res.InvalidOS))
	}

	report := res.ToReport()
	if !strings.Contains(report, "*Available Versions:*") {
		t.Errorf("expected report to contain Available Versions header, got: %s", report)
	}
	if strings.Contains(report, "*Missing Release Binary") {
		t.Errorf("expected report to NOT contain Missing sections, got: %s", report)
	}

	tableOut := res.ToTable()
	if !strings.Contains(tableOut, "VALID") {
		t.Errorf("expected table to contain VALID, got: %s", tableOut)
	}

	jsonBytes, err := res.ToJSON()
	if err != nil {
		t.Fatalf("unexpected json error: %v", err)
	}
	var jsonRes VerifyBinariesRes
	if err := json.Unmarshal(jsonBytes, &jsonRes); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}
	if !jsonRes.Success || jsonRes.HasMissing {
		t.Errorf("expected JSON success:true, has_missing:false, got: %+v", jsonRes)
	}

	ghOutput, err := res.ToGithubOutput()
	if err != nil {
		t.Fatalf("unexpected error from ToGithubOutput: %v", err)
	}
	if len(ghOutput) == 0 {
		t.Errorf("expected non-empty ToGithubOutput")
	}
}

func TestVerifyBinaries_MissingVersions(t *testing.T) {
	data := &ListBinaryVersionsRes{
		AllVersions:     []string{"2.2.1", "1.21.13", "1.19.24"},
		ValidVersions:   makeValidBinaryData("2.2.1").ValidVersions,
		InvalidVersions: []string{"1.21.13", "1.19.24"},
	}

	req := &VerifyBinariesReq{
		VersionsString: "2.2.1 1.21.13 1.19.24",
		BinaryFetcher:  &mockFetcher{data: data},
	}

	res, err := req.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Success {
		t.Errorf("expected Success=false, got true")
	}
	if !res.HasMissing {
		t.Errorf("expected HasMissing=true, got false")
	}
	if len(res.MissingVersions) != 2 {
		t.Fatalf("expected 2 missing versions, got %d", len(res.MissingVersions))
	}
	if res.MissingVersions[0] != "1.21.13" || res.MissingVersions[1] != "1.19.24" {
		t.Errorf("unexpected missing versions list: %v", res.MissingVersions)
	}

	report := res.ToReport()
	if !strings.Contains(report, "*Missing Release Binary Versions:*") {
		t.Errorf("expected report to contain Missing Versions header, got: %s", report)
	}
	if !strings.Contains(report, "Missing Versions: 1.21.13") {
		t.Errorf("expected report to mention 1.21.13, got: %s", report)
	}

	tableOut := res.ToTable()
	if !strings.Contains(tableOut, "MISSING") {
		t.Errorf("expected table to contain MISSING status, got: %s", tableOut)
	}
	if !strings.Contains(tableOut, "All") {
		t.Errorf("expected table to show 'All' for missing version variants/OS, got: %s", tableOut)
	}

	jsonBytes, err := res.ToJSON()
	if err != nil {
		t.Fatalf("unexpected json error: %v", err)
	}
	var jsonRes VerifyBinariesRes
	if err := json.Unmarshal(jsonBytes, &jsonRes); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}
	if jsonRes.Success || !jsonRes.HasMissing || len(jsonRes.MissingVersions) != 2 {
		t.Errorf("unexpected JSON response for missing versions: %+v", jsonRes)
	}
}

func TestVerifyBinaries_MissingVariants(t *testing.T) {
	// 2.2.1 is missing fips and hsm.fips
	data := makeVersionData("2.2.1", []VariantInfo{
		{Variant: "2.2.1", OS: allStandardOS},
		{Variant: "2.2.1+ent", OS: allStandardOS},
		{Variant: "2.2.1+ent.hsm", OS: linuxOnlyOS},
	})

	req := &VerifyBinariesReq{
		VersionsString: "2.2.1",
		BinaryFetcher:  &mockFetcher{data: data},
	}

	res, err := req.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Success {
		t.Errorf("expected Success=false, got true")
	}
	if len(res.MissingVariants) != 2 {
		t.Fatalf("expected 2 missing variants, got %d: %v", len(res.MissingVariants), res.MissingVariants)
	}

	report := res.ToReport()
	if !strings.Contains(report, "*Missing Release Binary Variants:*") {
		t.Errorf("expected report to contain Missing Variants header, got: %s", report)
	}
	if !strings.Contains(report, "Version 2.2.1 is missing variant: 2.2.1+ent.fips") {
		t.Errorf("expected report to mention missing fips variant, got: %s", report)
	}
	if !strings.Contains(report, "Version 2.2.1 is missing variant: 2.2.1+ent.hsm.fips") {
		t.Errorf("expected report to mention missing hsm.fips variant, got: %s", report)
	}

	tableOut := res.ToTable()
	if !strings.Contains(tableOut, "INCOMPLETE") {
		t.Errorf("expected table to contain INCOMPLETE status, got: %s", tableOut)
	}
	if !strings.Contains(tableOut, "2.2.1+ent.fips") || !strings.Contains(tableOut, "2.2.1+ent.hsm.fips") {
		t.Errorf("expected table to show missing variant names, got: %s", tableOut)
	}

	jsonBytes, err := res.ToJSON()
	if err != nil {
		t.Fatalf("unexpected json error: %v", err)
	}
	var jsonRes VerifyBinariesRes
	if err := json.Unmarshal(jsonBytes, &jsonRes); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}
	if len(jsonRes.MissingVariants) != 2 {
		t.Errorf("expected 2 missing variants in JSON, got: %v", jsonRes.MissingVariants)
	}
}

func TestVerifyBinaries_HSMExcludesHSMFIPS(t *testing.T) {
	// Regression test: +ent.hsm.fips1403 must NOT satisfy the +ent.hsm variant check.
	data := makeVersionData("2.2.1", []VariantInfo{
		{Variant: "2.2.1", OS: allStandardOS},
		{Variant: "2.2.1+ent", OS: allStandardOS},
		{Variant: "2.2.1+ent.fips1403", OS: linuxOnlyOS},
		{Variant: "2.2.1+ent.hsm.fips1403", OS: linuxOnlyOS}, // only HSM.FIPS, plain HSM is missing!
	})

	req := &VerifyBinariesReq{
		VersionsString: "2.2.1",
		BinaryFetcher:  &mockFetcher{data: data},
	}

	res, err := req.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Success {
		t.Errorf("expected Success=false when plain HSM is missing, got true")
	}

	foundMissingHSM := false
	for _, mv := range res.MissingVariants {
		if strings.Contains(mv, "2.2.1+ent.hsm") && !strings.Contains(mv, "fips") {
			foundMissingHSM = true
			break
		}
	}
	if !foundMissingHSM {
		t.Errorf("expected missing variant 2.2.1+ent.hsm in %v", res.MissingVariants)
	}
}

func TestVerifyBinaries_MissingOS_MissingHSMLinux(t *testing.T) {
	data := makeVersionData("1.18.5", []VariantInfo{
		{Variant: "1.18.5", OS: allStandardOS},
		{Variant: "1.18.5+ent", OS: RequiredOSList},
		{Variant: "1.18.5+ent.fips1402", OS: linuxOnlyOS},
		{Variant: "1.18.5+ent.hsm", OS: []string{}}, // missing linux
		{Variant: "1.18.5+ent.hsm.fips1402", OS: linuxOnlyOS},
	})

	req := &VerifyBinariesReq{
		VersionsString: "1.18.5",
		BinaryFetcher:  &mockFetcher{data: data},
	}

	res, err := req.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Success {
		t.Errorf("expected Success=false when linux is missing on HSM, got true")
	}
	if len(res.MissingOS) != 1 {
		t.Fatalf("expected 1 missing OS item, got %d: %v", len(res.MissingOS), res.MissingOS)
	}

	report := res.ToReport()
	if !strings.Contains(report, "Variant 1.18.5+ent.hsm of version 1.18.5 is missing OS: linux") {
		t.Errorf("expected report to mention missing linux on HSM, got: %s", report)
	}
}

func TestVerifyBinaries_OptionalOSAcceptedButNotRequired(t *testing.T) {
	for name, tc := range map[string]struct {
		version string
		os      []string
	}{
		"1.x with all 7 OS":    {"1.21.1", allStandardOS},
		"1.x with required OS": {"1.21.13", RequiredOSList},
		"2.x with required OS": {"2.2.1", RequiredOSList},
		"2.x with all 7 OS":    {"2.2.1", allStandardOS},
	} {
		t.Run(name, func(t *testing.T) {
			data := makeVersionData(tc.version, []VariantInfo{
				{Variant: tc.version, OS: tc.os},
				{Variant: tc.version + "+ent", OS: tc.os},
				{Variant: tc.version + "+ent.fips1403", OS: linuxOnlyOS},
				{Variant: tc.version + "+ent.hsm", OS: linuxOnlyOS},
				{Variant: tc.version + "+ent.hsm.fips1403", OS: linuxOnlyOS},
			})
			req := &VerifyBinariesReq{VersionsString: tc.version, BinaryFetcher: &mockFetcher{data: data}}

			res, err := req.Run(context.Background())
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !res.Success || res.HasMissing {
				t.Errorf("expected success, got missing=%v invalid=%v", res.MissingOS, res.InvalidOS)
			}
		})
	}
}

func TestVerifyBinaries_UnknownStandardOSIsInvalid(t *testing.T) {
	data := makeVersionData("2.2.1", []VariantInfo{
		{Variant: "2.2.1", OS: append([]string{"plan9"}, RequiredOSList...)},
		{Variant: "2.2.1+ent", OS: RequiredOSList},
		{Variant: "2.2.1+ent.fips1403", OS: linuxOnlyOS},
		{Variant: "2.2.1+ent.hsm", OS: linuxOnlyOS},
		{Variant: "2.2.1+ent.hsm.fips1403", OS: linuxOnlyOS},
	})
	req := &VerifyBinariesReq{VersionsString: "2.2.1", BinaryFetcher: &mockFetcher{data: data}}

	res, err := req.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Success {
		t.Errorf("expected Success=false for unknown OS plan9")
	}
	if !strings.Contains(res.ToReport(), "Variant 2.2.1 of version 2.2.1 has invalid OS: plan9") {
		t.Errorf("expected report to mention invalid plan9, got: %s", res.ToReport())
	}
}

func TestVerifyBinaries_MissingOS_V2_MissingWindows(t *testing.T) {
	// Required OS: darwin, freebsd, linux, windows; missing "windows" on CE
	data := makeVersionData("2.2.1", []VariantInfo{
		{Variant: "2.2.1", OS: []string{"darwin", "freebsd", "linux"}}, // missing windows
		{Variant: "2.2.1+ent", OS: RequiredOSList},
		{Variant: "2.2.1+ent.fips1403", OS: linuxOnlyOS},
		{Variant: "2.2.1+ent.hsm", OS: linuxOnlyOS},
		{Variant: "2.2.1+ent.hsm.fips1403", OS: linuxOnlyOS},
	})

	req := &VerifyBinariesReq{
		VersionsString: "2.2.1",
		BinaryFetcher:  &mockFetcher{data: data},
	}

	res, err := req.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Success {
		t.Errorf("expected Success=false when windows is missing on 2.x, got true")
	}
	if len(res.MissingOS) != 1 {
		t.Fatalf("expected 1 missing OS item, got %d: %v", len(res.MissingOS), res.MissingOS)
	}

	report := res.ToReport()
	if !strings.Contains(report, "Variant 2.2.1 of version 2.2.1 is missing OS: windows") {
		t.Errorf("expected report to mention missing windows, got: %s", report)
	}
}

func TestVerifyBinaries_InvalidOS(t *testing.T) {
	// FIPS variant unexpectedly has windows and darwin OS distros
	data := makeVersionData("2.2.1", []VariantInfo{
		{Variant: "2.2.1", OS: allStandardOS},
		{Variant: "2.2.1+ent", OS: allStandardOS},
		{Variant: "2.2.1+ent.fips1403", OS: []string{"linux", "windows", "darwin"}},
		{Variant: "2.2.1+ent.hsm", OS: linuxOnlyOS},
		{Variant: "2.2.1+ent.hsm.fips1403", OS: linuxOnlyOS},
	})

	req := &VerifyBinariesReq{
		VersionsString: "2.2.1",
		BinaryFetcher:  &mockFetcher{data: data},
	}

	res, err := req.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Success {
		t.Errorf("expected Success=false, got true")
	}
	if len(res.InvalidOS) != 2 {
		t.Fatalf("expected 2 invalid OS items, got %d: %v", len(res.InvalidOS), res.InvalidOS)
	}

	report := res.ToReport()
	if !strings.Contains(report, "Variant 2.2.1+ent.fips1403 of version 2.2.1 has invalid OS: windows") {
		t.Errorf("expected report to mention invalid windows OS, got: %s", report)
	}
	if !strings.Contains(report, "Variant 2.2.1+ent.fips1403 of version 2.2.1 has invalid OS: darwin") {
		t.Errorf("expected report to mention invalid darwin OS, got: %s", report)
	}

	tableOut := res.ToTable()
	if !strings.Contains(tableOut, "invalid [windows,darwin]") {
		t.Errorf("expected table to show invalid OS details, got: %s", tableOut)
	}

	jsonBytes, err := res.ToJSON()
	if err != nil {
		t.Fatalf("unexpected json error: %v", err)
	}
	var jsonRes VerifyBinariesRes
	if err := json.Unmarshal(jsonBytes, &jsonRes); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}
	if len(jsonRes.InvalidOS) != 2 {
		t.Errorf("expected 2 invalid OS items in JSON, got: %v", jsonRes.InvalidOS)
	}
}

func TestVerifyBinaries_EmptyVersionsError(t *testing.T) {
	req := &VerifyBinariesReq{
		VersionsString: "",
		BinaryFetcher:  &mockFetcher{},
	}

	_, err := req.Run(context.Background())
	if err == nil {
		t.Fatalf("expected error when no versions provided, got nil")
	}
}

func TestVerifyBinaries_LatestVersionRequiresCEVariant(t *testing.T) {
	// The latest version must have a CE folder, so a missing 2.2.1 fails verification.
	data := makeVersionData("2.2.1", []VariantInfo{
		{Variant: "2.2.1+ent", OS: allStandardOS},
		{Variant: "2.2.1+ent.fips1403", OS: linuxOnlyOS},
		{Variant: "2.2.1+ent.hsm", OS: linuxOnlyOS},
		{Variant: "2.2.1+ent.hsm.fips1403", OS: linuxOnlyOS},
	})

	req := &VerifyBinariesReq{
		VersionsString: "2.2.1",
		BinaryFetcher:  &mockFetcher{data: data},
	}

	res, err := req.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Success {
		t.Errorf("expected Success=false when CE variant is missing, got true")
	}
	if len(res.MissingVariants) != 1 || !strings.Contains(res.MissingVariants[0], "missing variant: 2.2.1") {
		t.Errorf("expected missing CE variant 2.2.1, got: %v", res.MissingVariants)
	}
}

func TestVerifyBinaries_ExplicitEntLabelHasNoCE(t *testing.T) {
	// An enterprise-only label needs only the 4 enterprise variants and no CE variant.
	data := makeVersionData("1.21.13-ent", []VariantInfo{
		{Variant: "1.21.13+ent", OS: allStandardOS},
		{Variant: "1.21.13+ent.fips1403", OS: linuxOnlyOS},
		{Variant: "1.21.13+ent.hsm", OS: linuxOnlyOS},
		{Variant: "1.21.13+ent.hsm.fips1403", OS: linuxOnlyOS},
	})

	req := &VerifyBinariesReq{
		VersionsString: "1.21.13-ent",
		BinaryFetcher:  &mockFetcher{data: data},
	}

	res, err := req.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Success {
		t.Errorf("expected Success=true for enterprise-only version, got missing variants: %v", res.MissingVariants)
	}
}

func TestVerifyBinaries_CEFromVersionsConfig(t *testing.T) {
	cfg := &VersionsConfig{ActiveVersion: &ActiveVersion{Versions: map[string]*Version{
		"2.2.x":  {CEActive: true},
		"1.21.x": {CEActive: false},
	}}}

	tests := map[string]struct {
		versions string
		want     string
	}{
		"lone hotfix on a CE-less line is enterprise-only": {"1.21.5", "1.21.5-ent"},
		"lone release on a CE line keeps CE":               {"2.2.3", "2.2.3"},
		"mixed list decides per line":                      {"2.2.1 1.21.13", "2.2.1 1.21.13-ent"},
		"line missing from config falls back to highest":   {"3.0.0 2.1.4", "3.0.0 2.1.4-ent"},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			f := &capturingFetcher{data: makeVersionData("x", nil)}
			req := &VerifyBinariesReq{VersionsString: tc.versions, BinaryFetcher: f, VersionsConfig: cfg}
			_, _ = req.Run(context.Background())
			if f.got != tc.want {
				t.Errorf("fetched %q, want %q", f.got, tc.want)
			}
		})
	}
}

type capturingFetcher struct {
	got  string
	data *ListBinaryVersionsRes
}

func (c *capturingFetcher) Fetch(ctx context.Context, versions string) (*ListBinaryVersionsRes, error) {
	c.got = versions
	return c.data, nil
}

func TestVerifyBinaries_EnterpriseOnlySeriesRequestsEntFolder(t *testing.T) {
	// The upstream fetcher keys enterprise-only versions by their "-ent" label.
	data := makeVersionData("1.21.13-ent", []VariantInfo{
		{Variant: "1.21.13+ent", OS: RequiredOSList},
		{Variant: "1.21.13+ent.fips1403", OS: linuxOnlyOS},
		{Variant: "1.21.13+ent.hsm", OS: linuxOnlyOS},
		{Variant: "1.21.13+ent.hsm.fips1403", OS: linuxOnlyOS},
	})
	fetcher := &capturingFetcher{data: data}
	req := &VerifyBinariesReq{
		VersionsString: "1.21.13 2.2.1",
		BinaryFetcher:  fetcher,
	}

	ltsFetcher := &capturingFetcher{data: data}
	ltsReq := &VerifyBinariesReq{VersionsString: "1.21.13-lts 2.2.1", BinaryFetcher: ltsFetcher}
	if _, err := ltsReq.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ltsFetcher.got != "1.21.13-ent 2.2.1" {
		t.Errorf("expected -lts label to become %q, got %q", "1.21.13-ent 2.2.1", ltsFetcher.got)
	}

	res, err := req.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fetcher.got != "1.21.13-ent 2.2.1" {
		t.Errorf("expected fetch labels %q, got %q", "1.21.13-ent 2.2.1", fetcher.got)
	}
	if !res.Success {
		t.Errorf("expected success, got variants=%v os=%v invalid=%v", res.MissingVariants, res.MissingOS, res.InvalidOS)
	}
}

func TestVerifyBinaries_OnlyLatestVersionKeepsCELabel(t *testing.T) {
	for name, tc := range map[string]struct{ in, want string }{
		"single version is latest":        {"1.21.13", "1.21.13"},
		"order does not matter":           {"1.19.24 2.2.1 1.21.13", "1.19.24-ent 2.2.1 1.21.13-ent"},
		"rc is lower than its release":    {"2.2.0-rc1 2.2.0", "2.2.0-rc1-ent 2.2.0"},
		"explicit ent labels are kept":    {"2.2.1 1.21.13+ent", "2.2.1 1.21.13+ent"},
		"duplicate labels share a latest": {"2.2.1 2.2.1-ce", "2.2.1 2.2.1-ce"},
	} {
		t.Run(name, func(t *testing.T) {
			f := &capturingFetcher{data: makeVersionData("x", nil)}
			_, err := (&VerifyBinariesReq{VersionsString: tc.in, BinaryFetcher: f}).Run(context.Background())
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if f.got != tc.want {
				t.Errorf("labels for %q: want %q, got %q", tc.in, tc.want, f.got)
			}
		})
	}
}
