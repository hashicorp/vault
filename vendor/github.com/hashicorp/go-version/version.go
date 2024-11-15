// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package version

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// The compiled regular expression used to test the validity of a version.
var (
	versionRegexp *regexp.Regexp
	semverRegexp  *regexp.Regexp
)

// The raw regular expression string used for testing the validity
// of a version.
const (
	VersionRegexpRaw string = `v?([0-9]+(\.[0-9]+)*?)` +
		`(-([0-9]+[0-9A-Za-z\-~]*(\.[0-9A-Za-z\-~]+)*)|(-?([A-Za-z\-~]+[0-9A-Za-z\-~]*(\.[0-9A-Za-z\-~]+)*)))?` +
		`(\+([0-9A-Za-z\-~]+(\.[0-9A-Za-z\-~]+)*))?` +
		`?`

	// SemverRegexpRaw requires a separator between version and prerelease
	SemverRegexpRaw string = `v?([0-9]+(\.[0-9]+)*?)` +
		`(-([0-9]+[0-9A-Za-z\-~]*(\.[0-9A-Za-z\-~]+)*)|(-([A-Za-z\-~]+[0-9A-Za-z\-~]*(\.[0-9A-Za-z\-~]+)*)))?` +
		`(\+([0-9A-Za-z\-~]+(\.[0-9A-Za-z\-~]+)*))?` +
		`?`
)

// Version represents a single version.
type Version struct {
	metadata string
	pre      string
	segments []int64
	si       int
	original string
}

func init() {
	versionRegexp = regexp.MustCompile("^" + VersionRegexpRaw + "$")
	semverRegexp = regexp.MustCompile("^" + SemverRegexpRaw + "$")
}

// NewVersion parses the given version and returns a new
// Version.
func NewVersion(v string) (*Version, error) {
	return newVersion(v, versionRegexp)
}

// NewSemver parses the given version and returns a new
// Version that adheres strictly to SemVer specs
// https://semver.org/
func NewSemver(v string) (*Version, error) {
	return newVersion(v, semverRegexp)
}

func newVersion(v string, pattern *regexp.Regexp) (*Version, error) {
	matches := pattern.FindStringSubmatch(v)
	if matches == nil {
		return nil, fmt.Errorf("Malformed version: %s", v)
	}
	segmentsStr := strings.Split(matches[1], ".")
	segments := make([]int64, len(segmentsStr))
	for i, str := range segmentsStr {
		val, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, fmt.Errorf(
				"Error parsing version: %s", err)
		}

		segments[i] = val
	}

	// Even though we could support more than three segments, if we
	// got less than three, pad it with 0s. This is to cover the basic
	// default usecase of semver, which is MAJOR.MINOR.PATCH at the minimum
	for i := len(segments); i < 3; i++ {
		segments = append(segments, 0)
	}

	pre := matches[7]
	if pre == "" {
		pre = matches[4]
	}

	return &Version{
		metadata: matches[10],
		pre:      pre,
		segments: segments,
		si:       len(segmentsStr),
		original: v,
	}, nil
}

// Must is a helper that wraps a call to a function returning (*Version, error)
// and panics if error is non-nil.
func Must(v *Version, err error) *Version {
	if err != nil {
		panic(err)
	}

	return v
}

// Compare compares this version to another version. This
// returns -1, 0, or 1 if this version is smaller, equal,
// or larger than the other version, respectively.
//
// If you want boolean results, use the LessThan, Equal,
// GreaterThan, GreaterThanOrEqual or LessThanOrEqual methods.
func (v *Version) Compare(other *Version) int {
	// A quick, efficient equality check
	if v.String() == other.String() {
		return 0
	}

	// If the segments are the same, we must compare on prerelease info
	if v.equalSegments(other) {
		preSelf := v.Prerelease()
		preOther := other.Prerelease()
		if preSelf == "" && preOther == "" {
			return 0
		}
		if preSelf == "" {
			return 1
		}
		if preOther == "" {
			return -1
		}

		return comparePrereleases(preSelf, preOther)
	}

	segmentsSelf := v.Segments64()
	segmentsOther := other.Segments64()
	// Get the highest specificity (hS), or if they're equal, just use segmentSelf length
	lenSelf := len(segmentsSelf)
	lenOther := len(segmentsOther)
	hS := lenSelf
	if lenSelf < lenOther {
		hS = lenOther
	}
	// Compare the segments
	// Because a constraint could have more/less specificity than the version it's
	// checking, we need to account for a lopsided or jagged comparison
	for i := 0; i < hS; i++ {
		if i > lenSelf-1 {
			// This means Self had the lower specificity
			// Check to see if the remaining segments in Other are all zeros
			if !allZero(segmentsOther[i:]) {
				// if not, it means that Other has to be greater than Self
				return -1
			}
			break
		} else if i > lenOther-1 {
			// this means Other had the lower specificity
			// Check to see if the remaining segments in Self are all zeros -
			if !allZero(segmentsSelf[i:]) {
				// if not, it means that Self has to be greater than Other
				return 1
			}
			break
		}
		lhs := segmentsSelf[i]
		rhs := segmentsOther[i]
		if lhs == rhs {
			continue
		} else if lhs < rhs {
			return -1
		}
		// Otherwis, rhs was > lhs, they're not equal
		return 1
	}

	// if we got this far, they're equal
	return 0
}

func (v *Version) equalSegments(other *Version) bool {
	segmentsSelf := v.Segments64()
	segmentsOther := other.Segments64()

	if len(segmentsSelf) != len(segmentsOther) {
		return false
	}
	for i, v := range segmentsSelf {
		if v != segmentsOther[i] {
			return false
		}
	}
	return true
}

func allZero(segs []int64) bool {
	for _, s := range segs {
		if s != 0 {
			return false
		}
	}
	return true
}

func comparePart(preSelf string, preOther string) int {
	if preSelf == preOther {
		return 0
	}

	var selfInt int64
	selfNumeric := true
	selfInt, err := strconv.ParseInt(preSelf, 10, 64)
	if err != nil {
		selfNumeric = false
	}

	var otherInt int64
	otherNumeric := true
	otherInt, err = strconv.ParseInt(preOther, 10, 64)
	if err != nil {
		otherNumeric = false
	}

	// if a part is empty, we use the other to decide
	if preSelf == "" {
		if otherNumeric {
			return -1
		}
		return 1
	}

	if preOther == "" {
		if selfNumeric {
			return 1
		}
		return -1
	}

	if selfNumeric && !otherNumeric {
		return -1
	} else if !selfNumeric && otherNumeric {
		return 1
	} else if !selfNumeric && !otherNumeric && preSelf > preOther {
		return 1
	} else if selfInt > otherInt {
		return 1
	}

	return -1
}

func comparePrereleases(v string, other string) int {
	// the same pre release!
	if v == other {
		return 0
	}

	// split both pre releases for analyse their parts
	selfPreReleaseMeta := strings.Split(v, ".")
	otherPreReleaseMeta := strings.Split(other, ".")

	selfPreReleaseLen := len(selfPreReleaseMeta)
	otherPreReleaseLen := len(otherPreReleaseMeta)

	biggestLen := otherPreReleaseLen
	if selfPreReleaseLen > otherPreReleaseLen {
		biggestLen = selfPreReleaseLen
	}

	// loop for parts to find the first difference
	for i := 0; i < biggestLen; i = i + 1 {
		partSelfPre := ""
		if i < selfPreReleaseLen {
			partSelfPre = selfPreReleaseMeta[i]
		}

		partOtherPre := ""
		if i < otherPreReleaseLen {
			partOtherPre = otherPreReleaseMeta[i]
		}

		compare := comparePart(partSelfPre, partOtherPre)
		// if parts are equals, continue the loop
		if compare != 0 {
			return compare
		}
	}

	return 0
}

// Core returns a new version constructed from only the MAJOR.MINOR.PATCH
// segments of the version, without prerelease or metadata.
func (v *Version) Core() *Version {
	segments := v.Segments64()
	segmentsOnly := fmt.Sprintf("%d.%d.%d", segments[0], segments[1], segments[2])
	return Must(NewVersion(segmentsOnly))
}

// Equal tests if two versions are equal.
func (v *Version) Equal(o *Version) bool {
	if v == nil || o == nil {
		return v == o
	}

	return v.Compare(o) == 0
}

// GreaterThan tests if this version is greater than another version.
func (v *Version) GreaterThan(o *Version) bool {
	return v.Compare(o) > 0
}

// GreaterThanOrEqual tests if this version is greater than or equal to another version.
func (v *Version) GreaterThanOrEqual(o *Version) bool {
	return v.Compare(o) >= 0
}

// LessThan tests if this version is less than another version.
func (v *Version) LessThan(o *Version) bool {
	return v.Compare(o) < 0
}

// LessThanOrEqual tests if this version is less than or equal to another version.
func (v *Version) LessThanOrEqual(o *Version) bool {
	return v.Compare(o) <= 0
}

// Metadata returns any metadata that was part of the version
// string.
//
// Metadata is anything that comes after the "+" in the version.
// For example, with "1.2.3+beta", the metadata is "beta".
func (v *Version) Metadata() string {
	return v.metadata
}

// Prerelease returns any prerelease data that is part of the version,
// or blank if there is no prerelease data.
//
// Prerelease information is anything that comes after the "-" in the
// version (but before any metadata). For example, with "1.2.3-beta",
// the prerelease information is "beta".
func (v *Version) Prerelease() string {
	return v.pre
}

// Segments returns the numeric segments of the version as a slice of ints.
//
// This excludes any metadata or pre-release information. For example,
// for a version "1.2.3-beta", segments will return a slice of
// 1, 2, 3.
func (v *Version) Segments() []int {
	segmentSlice := make([]int, len(v.segments))
	for i, v := range v.segments {
		segmentSlice[i] = int(v)
	}
	return segmentSlice
}

// Segments64 returns the numeric segments of the version as a slice of int64s.
//
// This excludes any metadata or pre-release information. For example,
// for a version "1.2.3-beta", segments will return a slice of
// 1, 2, 3.
func (v *Version) Segments64() []int64 {
	result := make([]int64, len(v.segments))
	copy(result, v.segments)
	return result
}

// String returns the full version string included pre-release
// and metadata information.
//
// This value is rebuilt according to the parsed segments and other
// information. Therefore, ambiguities in the version string such as
// prefixed zeroes (1.04.0 => 1.4.0), `v` prefix (v1.0.0 => 1.0.0), and
// missing parts (1.0 => 1.0.0) will be made into a canonicalized form
// as shown in the parenthesized examples.
func (v *Version) String() string {
	var buf bytes.Buffer
	fmtParts := make([]string, len(v.segments))
	for i, s := range v.segments {
		// We can ignore err here since we've pre-parsed the values in segments
		str := strconv.FormatInt(s, 10)
		fmtParts[i] = str
	}
	fmt.Fprintf(&buf, strings.Join(fmtParts, "."))
	if v.pre != "" {
		fmt.Fprintf(&buf, "-%s", v.pre)
	}
	if v.metadata != "" {
		fmt.Fprintf(&buf, "+%s", v.metadata)
	}

	return buf.String()
}

// Original returns the original parsed version as-is, including any
// potential whitespace, `v` prefix, etc.
func (v *Version) Original() string {
	return v.original
}

// UnmarshalText implements encoding.TextUnmarshaler interface.
func (v *Version) UnmarshalText(b []byte) error {
	temp, err := NewVersion(string(b))
	if err != nil {
		return err
	}

	*v = *temp

	return nil
}

// MarshalText implements encoding.TextMarshaler interface.
func (v *Version) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// Scan implements the sql.Scanner interface.
func (v *Version) Scan(src interface{}) error {
	switch src := src.(type) {
	case string:
		return v.UnmarshalText([]byte(src))
	case nil:
		return nil
	default:
		return fmt.Errorf("cannot scan %T as Version", src)
	}
}

// Value implements the driver.Valuer interface.
func (v *Version) Value() (driver.Value, error) {
	return v.String(), nil
}
