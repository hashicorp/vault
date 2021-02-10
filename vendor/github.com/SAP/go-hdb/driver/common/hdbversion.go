// SPDX-FileCopyrightText: 2019-2020 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	hdbVersionMajor = iota
	hdbVersionMinor
	hdbVersionRevision
	hdbVersionPatch
	hdbVersionBuildID
	hdbVersionCount
)

// hdbVersionNumber holds the information of a hdb semantic version.
//
// u.vv.wwx.yy.zzzzzzzzzz
//
// u.vv:       hdb version (major.minor)
// ww:         SPS number
// wwx:        revision number
// yy:         patch number
// zzzzzzzzzz: build id
//
// Example: 2.00.045.00.1575639312
//
// hdb version:     2.00
// SPS number:      04
// revision number: 045
// patch number:    0
// build id:        1575639312
type hdbVersionNumber []uint64 // assumption: all fields are numeric

func parseHDBVersionNumber(s string) hdbVersionNumber {
	vn := make([]uint64, hdbVersionCount)

	parts := strings.SplitN(s, ".", hdbVersionCount)
	for i := 0; i < len(parts); i++ {
		vn[i], _ = strconv.ParseUint(parts[i], 10, 64)
	}
	return vn
}

func formatUint64(i uint64, digits int) string {
	s := strings.Repeat("0", digits) + strconv.FormatUint(i, 10)
	return s[len(s)-digits:]
}

func (vn hdbVersionNumber) String() string {
	s := fmt.Sprintf("%d.%s.%s.%s", vn[hdbVersionMajor], formatUint64(vn[hdbVersionMinor], 2), formatUint64(vn[hdbVersionRevision], 3), formatUint64(vn[hdbVersionPatch], 2))
	if vn[hdbVersionBuildID] != 0 {
		return fmt.Sprintf("%s.%d", s, vn[hdbVersionBuildID])
	}
	return s
}

func (vn hdbVersionNumber) isZero() bool {
	for _, n := range vn {
		if n != 0 {
			return false
		}
	}
	return true
}

func compareUint64(u1, u2 uint64) int {
	switch {
	case u1 == u2:
		return 0
	case u1 > u2:
		return 1
	default:
		return -1
	}
}

// Major returns the major field of a hdbVersionNumber.
func (vn hdbVersionNumber) Major() uint64 { return vn[hdbVersionMajor] }

// Minor returns the minor field of a HDBVersionNumber.
func (vn hdbVersionNumber) Minor() uint64 { return vn[hdbVersionMinor] }

// SPS returns the sps field of a HDBVersionNumber.
func (vn hdbVersionNumber) SPS() uint64 { return vn[hdbVersionRevision] / 10 }

// Revision returns the revision field of a HDBVersionNumber.
func (vn hdbVersionNumber) Revision() uint64 { return vn[hdbVersionRevision] }

// Patch returns the patch field of a HDBVersionNumber.
func (vn hdbVersionNumber) Patch() uint64 { return vn[hdbVersionPatch] }

// BuildID returns the build id field of a HDBVersionNumber.
func (vn hdbVersionNumber) BuildID() uint64 { return vn[hdbVersionBuildID] }

// Compare compares the version number with a second version number vn2. The result will be
//  0 in case the two versions are equal,
// -1 in case version v has lower precedence than c2,
//  1 in case version v has higher precedence than c2.
func (vn hdbVersionNumber) compare(vn2 hdbVersionNumber) int {
	for i := 0; i < (hdbVersionCount - 1); i++ { // ignore buildID - might not be ordered}
		if r := compareUint64(vn[i], vn2[i]); r != 0 {
			return r
		}
	}
	return 0
}

// hdbVersionNumberOne - if HANA version 1 assume version 1.00 SPS 12.
var hdbVersionNumberOne = parseHDBVersionNumber("1.00.120")

// HDBVersion feature flags.
const (
	HDBFNone              uint64 = 1 << iota
	HDBFServerVersion            // HANA reports server version in connect options
	HDBFConnectClientInfo        // HANA accepts ClientInfo as part of the connection process
)

var hdbFeatureAvailability = map[uint64]hdbVersionNumber{
	HDBFServerVersion:     parseHDBVersionNumber("2.00.000"),
	HDBFConnectClientInfo: parseHDBVersionNumber("2.00.042"),
}

// HDBVersion is representing a hdb version.
type HDBVersion struct {
	hdbVersionNumber
	feature uint64
}

// ParseHDBVersion parses a semantic hdb version string field.
func ParseHDBVersion(s string) *HDBVersion {
	number := parseHDBVersionNumber(s)
	if number.isZero() { // hdb 1.00 does not report version
		number = hdbVersionNumberOne
	}

	var feature uint64
	// detect features
	for f, cv := range hdbFeatureAvailability {
		if number.compare(cv) >= 0 { // v is equal or greater than cv
			feature |= f // add feature
		}
	}
	return &HDBVersion{hdbVersionNumber: number, feature: feature}
}

// Compare compares the version with a second version v2. The result will be
//  0 in case the two versions are equal,
// -1 in case version v has lower precedence than c2,
//  1 in case version v has higher precedence than c2.
func (v *HDBVersion) Compare(v2 *HDBVersion) int {
	return v.hdbVersionNumber.compare(v2.hdbVersionNumber)
}

// HasFeature returns true if HDBVersion does support feature - false otherwise.
func (v *HDBVersion) HasFeature(feature uint64) bool { return v.feature&feature != 0 }
