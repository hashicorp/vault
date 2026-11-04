// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package releases

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/Masterminds/semver"
	slogctx "github.com/veqryn/slog-context"
)

// BinaryVersionsFetcher defines an interface for retrieving binary version metadata.
type BinaryVersionsFetcher interface {
	Fetch(ctx context.Context, versions string) (*ListBinaryVersionsRes, error)
}

type defaultBinaryFetcher struct{}

func (d *defaultBinaryFetcher) Fetch(ctx context.Context, versions string) (*ListBinaryVersionsRes, error) {
	req := NewListBinaryVersionsReq(versions)
	return req.Run(ctx)
}

var (
	// RequiredOSList is the set of OS distros every standard (base and +ent) build must have.
	RequiredOSList = []string{"darwin", "freebsd", "linux", "windows"}
	// OptionalOSList is accepted but not required. These were dropped upstream starting
	// with 1.19.13, 1.20.7, 1.21.2 and 2.0.0 but are still present on older releases.
	OptionalOSList = []string{"netbsd", "openbsd", "solaris"}
)

// VariantOSDetail contains OS verification details for a single variant.
type VariantOSDetail struct {
	Variant   string   `json:"variant"`
	ValidOS   []string `json:"valid_os"`
	MissingOS []string `json:"missing_os"`
	InvalidOS []string `json:"invalid_os"`
}

// VersionVerificationDetail contains verification details for a single version.
type VersionVerificationDetail struct {
	Version         string                      `json:"version"`
	BaseVersion     string                      `json:"base_version"`
	Status          string                      `json:"status"` // "valid", "missing"
	ValidVariants   []string                    `json:"valid_variants"`
	MissingVariants []string                    `json:"missing_variants"`
	VariantsOS      map[string]*VariantOSDetail `json:"variants_os,omitempty"`
}

// VerifyBinariesRes contains the results of verifying production release binaries.
type VerifyBinariesRes struct {
	Success         bool                                  `json:"success"`
	HasMissing      bool                                  `json:"has_missing"`
	AllVersions     []string                              `json:"all_versions"`
	ValidVersions   []string                              `json:"valid_versions"`
	MissingVersions []string                              `json:"missing_versions"`
	MissingVariants []string                              `json:"missing_variants"`
	MissingOS       []string                              `json:"missing_os"`
	InvalidOS       []string                              `json:"invalid_os"`
	Details         map[string]*VersionVerificationDetail `json:"details,omitempty"`
	Report          string                                `json:"report,omitempty"`
}

// VerifyBinariesReq defines the request parameters for verifying production release binaries.
type VerifyBinariesReq struct {
	VersionsString string
	BinaryFetcher  BinaryVersionsFetcher
	// VersionsConfig is the decoded .release/versions.hcl. When it has an entry
	// for a version's release line (e.g. 1.21.x), its ce_active value decides
	// whether that version has a CE build. Optional.
	VersionsConfig *VersionsConfig
}

// NewVerifyBinariesReq creates a new VerifyBinariesReq with default settings.
func NewVerifyBinariesReq() *VerifyBinariesReq {
	return &VerifyBinariesReq{
		BinaryFetcher: &defaultBinaryFetcher{},
	}
}

// Run executes the binary verification logic.
func (r *VerifyBinariesReq) Run(ctx context.Context) (*VerifyBinariesRes, error) {
	if r == nil {
		return nil, errors.New("uninitialized request")
	}

	ctx = slogctx.Append(ctx, slog.String("component", "releases.verify_binaries"))
	slog.Default().DebugContext(ctx, "starting verify binaries check")

	data, err := r.resolveBinaryData(ctx)
	if err != nil {
		return nil, fmt.Errorf("resolving binary data: %w", err)
	}

	res := verifyBinaryData(data)
	res.Report = res.ToReport()

	slog.Default().DebugContext(ctx, "completed verify binaries check",
		slog.Bool("success", res.Success),
		slog.Int("missing_versions", len(res.MissingVersions)),
		slog.Int("missing_variants", len(res.MissingVariants)),
		slog.Int("missing_os", len(res.MissingOS)),
		slog.Int("invalid_os", len(res.InvalidOS)),
	)

	return res, nil
}

// resolveBinaryData loads or fetches the ListBinaryVersionsRes data.
func (r *VerifyBinariesReq) resolveBinaryData(ctx context.Context) (*ListBinaryVersionsRes, error) {
	versionsStr := strings.TrimSpace(r.VersionsString)
	if versionsStr == "" {
		return nil, errors.New("no versions provided to verify")
	}

	// A version has a CE build only if its release line has ce_active = true in
	// versions.hcl. Lines missing from the config fall back to "only the highest
	// version in the list has CE". CE-less versions are enterprise-only and have
	// no plain folder upstream, only "+ent" ones.
	labels := strings.Fields(versionsStr)
	latest := latestVersion(labels)
	for i, label := range labels {
		if strings.HasSuffix(label, "-ent") || strings.HasSuffix(label, "+ent") {
			continue
		}
		clean := stripLabel(label)
		hasCE := clean == latest
		if ce, ok := r.ceActiveForVersion(clean); ok {
			hasCE = ce
		}
		if hasCE || latest == "" {
			continue
		}
		labels[i] = strings.TrimSuffix(strings.TrimSuffix(label, "-lts"), "-ce") + "-ent"
	}
	versionsStr = strings.Join(labels, " ")

	fetcher := r.BinaryFetcher
	if fetcher == nil {
		fetcher = &defaultBinaryFetcher{}
	}

	slog.Default().DebugContext(ctx, "fetching binary versions upstream", slog.String("versions", versionsStr))
	return fetcher.Fetch(ctx, versionsStr)
}

// ceActiveForVersion looks up the ce_active value for the release line of a
// version (e.g. 1.21.5 -> 1.21.x). ok is false when the config has no such line.
func (r *VerifyBinariesReq) ceActiveForVersion(version string) (ce bool, ok bool) {
	if r.VersionsConfig == nil || r.VersionsConfig.ActiveVersion == nil {
		return false, false
	}
	v, err := semver.NewVersion(version)
	if err != nil {
		return false, false
	}
	entry, found := r.VersionsConfig.ActiveVersion.Versions[fmt.Sprintf("%d.%d.x", v.Major(), v.Minor())]
	if !found || entry == nil {
		return false, false
	}
	return entry.CEActive, true
}

// stripLabel removes the label suffixes (-lts, -ce, -ent, +ent) and any "v" prefix.
func stripLabel(label string) string {
	label = strings.TrimPrefix(label, "v")
	for _, suffix := range []string{"-lts", "-ce", "-ent", "+ent"} {
		label = strings.TrimSuffix(label, suffix)
	}
	return label
}

// latestVersion returns the highest version among the labels, or "" if none parse.
func latestVersion(labels []string) string {
	var latest *semver.Version
	var latestStr string
	for _, label := range labels {
		clean := stripLabel(label)
		v, err := semver.NewVersion(clean)
		if err != nil {
			continue
		}
		if latest == nil || v.GreaterThan(latest) {
			latest, latestStr = v, clean
		}
	}
	return latestStr
}

// verifyBinaryData performs the core verification rules against ListBinaryVersionsRes.
func verifyBinaryData(data *ListBinaryVersionsRes) *VerifyBinariesRes {
	res := &VerifyBinariesRes{
		AllVersions:     data.AllVersions,
		ValidVersions:   make([]string, 0),
		MissingVersions: make([]string, 0),
		MissingVariants: make([]string, 0),
		MissingOS:       make([]string, 0),
		InvalidOS:       make([]string, 0),
		Details:         make(map[string]*VersionVerificationDetail),
	}

	// 1. Process invalid versions
	for _, inv := range data.InvalidVersions {
		res.MissingVersions = append(res.MissingVersions, inv)
		res.Details[inv] = &VersionVerificationDetail{
			Version: inv,
			Status:  "missing",
		}
	}

	// 2. Process valid versions
	var validKeys []string
	if len(data.AllVersions) > 0 {
		for _, v := range data.AllVersions {
			if _, ok := data.ValidVersions[v]; ok && !slices.Contains(validKeys, v) {
				validKeys = append(validKeys, v)
			}
		}
	}
	for v := range data.ValidVersions {
		if !slices.Contains(validKeys, v) {
			validKeys = append(validKeys, v)
		}
	}

	for _, v := range validKeys {
		entry := data.ValidVersions[v]
		res.ValidVersions = append(res.ValidVersions, v)

		baseVersion := strings.TrimSuffix(v, "-lts")
		baseVersion = strings.TrimSuffix(baseVersion, "-ce")
		baseVersion = strings.TrimSuffix(baseVersion, "-ent")

		detail := &VersionVerificationDetail{
			Version:         v,
			BaseVersion:     baseVersion,
			Status:          "valid",
			ValidVariants:   make([]string, 0),
			MissingVariants: make([]string, 0),
			VariantsOS:      make(map[string]*VariantOSDetail),
		}

		// Enterprise-only versions are labeled "-ent" by resolveBinaryData.
		ceActive := !strings.HasSuffix(v, "-ent") && !strings.HasSuffix(v, "+ent")

		type variantCheck struct {
			name    string
			prefix  string
			isExact bool
		}

		var reqVariants []variantCheck
		if ceActive {
			reqVariants = append(reqVariants, variantCheck{name: baseVersion, prefix: baseVersion, isExact: true})
		}
		reqVariants = append(reqVariants,
			variantCheck{name: baseVersion + "+ent", prefix: baseVersion + "+ent", isExact: true},
			variantCheck{name: baseVersion + "+ent.fips", prefix: baseVersion + "+ent.fips", isExact: false},
			variantCheck{name: baseVersion + "+ent.hsm", prefix: baseVersion + "+ent.hsm", isExact: false},
			variantCheck{name: baseVersion + "+ent.hsm.fips", prefix: baseVersion + "+ent.hsm.fips", isExact: false},
		)

		for _, req := range reqVariants {
			var matchedVariant string
			for _, act := range entry.Variants {
				if req.isExact {
					if act.Variant == req.prefix {
						matchedVariant = act.Variant
						break
					}
				} else if req.name == baseVersion+"+ent.hsm" {
					// ENT.HSM must not be satisfied by ENT.HSM.FIPS
					if strings.HasPrefix(act.Variant, req.prefix) && !strings.Contains(act.Variant, ".fips") {
						matchedVariant = act.Variant
						break
					}
				} else {
					if strings.HasPrefix(act.Variant, req.prefix) {
						matchedVariant = act.Variant
						break
					}
				}
			}

			if matchedVariant != "" {
				detail.ValidVariants = append(detail.ValidVariants, matchedVariant)
			} else {
				detail.MissingVariants = append(detail.MissingVariants, req.name)
				res.MissingVariants = append(res.MissingVariants, fmt.Sprintf("Version %s is missing variant: %s", v, req.name))
			}
		}

		// 3. Verify OS for each actual variant present
		for _, act := range entry.Variants {
			osDetail := &VariantOSDetail{
				Variant:   act.Variant,
				ValidOS:   make([]string, 0),
				MissingOS: make([]string, 0),
				InvalidOS: make([]string, 0),
			}

			if act.Variant == baseVersion || act.Variant == baseVersion+"+ent" {
				for _, requiredOS := range RequiredOSList {
					if slices.Contains(act.OS, requiredOS) {
						osDetail.ValidOS = append(osDetail.ValidOS, requiredOS)
					} else {
						osDetail.MissingOS = append(osDetail.MissingOS, requiredOS)
						res.MissingOS = append(res.MissingOS, fmt.Sprintf("Variant %s of version %s is missing OS: %s", act.Variant, v, requiredOS))
					}
				}

				for _, osName := range act.OS {
					switch {
					case slices.Contains(RequiredOSList, osName):
					case slices.Contains(OptionalOSList, osName):
						osDetail.ValidOS = append(osDetail.ValidOS, osName)
					default:
						osDetail.InvalidOS = append(osDetail.InvalidOS, osName)
						res.InvalidOS = append(res.InvalidOS, fmt.Sprintf("Variant %s of version %s has invalid OS: %s", act.Variant, v, osName))
					}
				}
			} else {
				// HSM / FIPS variants should only have linux
				if slices.Contains(act.OS, "linux") {
					osDetail.ValidOS = append(osDetail.ValidOS, "linux")
				} else {
					osDetail.MissingOS = append(osDetail.MissingOS, "linux")
					res.MissingOS = append(res.MissingOS, fmt.Sprintf("Variant %s of version %s is missing OS: linux", act.Variant, v))
				}

				for _, osName := range act.OS {
					if osName != "linux" {
						osDetail.InvalidOS = append(osDetail.InvalidOS, osName)
						res.InvalidOS = append(res.InvalidOS, fmt.Sprintf("Variant %s of version %s has invalid OS: %s", act.Variant, v, osName))
					}
				}
			}

			detail.VariantsOS[act.Variant] = osDetail
		}

		res.Details[v] = detail
	}

	res.HasMissing = len(res.MissingVersions) > 0 ||
		len(res.MissingVariants) > 0 ||
		len(res.MissingOS) > 0 ||
		len(res.InvalidOS) > 0

	res.Success = !res.HasMissing
	return res
}
