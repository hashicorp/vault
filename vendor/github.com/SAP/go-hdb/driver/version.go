package driver

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	versionMajor = iota
	versionMinor
	versionRevision
	versionPatch
	versionBuildID
	versionCount
)

/*
versionNumber holds the information of a hdb semantic version.

u.vv.wwx.yy.zzzzzzzzzz

u.vv:       hdb version (major.minor)
ww:         SPS number
wwx:        revision number
yy:         patch number
zzzzzzzzzz: build id

Example: 2.00.045.00.1575639312

	hdb version:     2.00
	SPS number:      04
	revision number: 045
	patch number:    0
	build id:        1575639312
*/
type versionNumber []uint64 // assumption: all fields are numeric

func (vn versionNumber) String() string {
	s := fmt.Sprintf("%d.%s.%s.%s",
		vn[versionMajor],
		formatUint64(vn[versionMinor], 2),
		formatUint64(vn[versionRevision], 3),
		formatUint64(vn[versionPatch], 2),
	)
	if vn[versionBuildID] != 0 {
		return fmt.Sprintf("%s.%d", s, vn[versionBuildID])
	}
	return s
}

func parseVersionNumber(s string) versionNumber {
	vn := make([]uint64, versionCount)

	parts := strings.SplitN(s, ".", versionCount)
	for i, part := range parts {
		vn[i], _ = strconv.ParseUint(part, 10, 64)
	}
	return vn
}

func formatUint64(i uint64, digits int) string {
	s := strings.Repeat("0", digits) + strconv.FormatUint(i, 10)
	return s[len(s)-digits:]
}

func (vn versionNumber) isZero() bool {
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

// Compare compares the version number with a second version number vn2. The result will be
//
//	0 in case the two versions are equal,
//
// -1 in case version v has lower precedence than c2,
//
//	1 in case version v has higher precedence than c2.
func (vn versionNumber) compare(vn2 versionNumber) int {
	for i := range versionCount - 1 { // ignore buildID - might not be ordered}
		if r := compareUint64(vn[i], vn2[i]); r != 0 {
			return r
		}
	}
	return 0
}

// hdbVersionNumberOne - if HANA version 1 assume version 1.00 SPS 12.
var versionNumberOne = parseVersionNumber("1.00.120")

// HDBVersion feature flags.
const (
	hdbfNone              uint64 = 1 << iota
	hdbfServerVersion            // HANA reports server version in connect options
	hdbfConnectClientInfo        // HANA accepts ClientInfo as part of the connection process
)

var hdbFeatureAvailability = map[uint64]versionNumber{
	hdbfServerVersion:     parseVersionNumber("2.00.000"),
	hdbfConnectClientInfo: parseVersionNumber("2.00.042"),
}

// Version is representing a hdb version.
type Version struct {
	vn      versionNumber
	feature uint64
}

func (v *Version) String() string { return v.vn.String() }

// Major returns the major field of a hdbVersionNumber.
func (v *Version) Major() uint64 { return v.vn[versionMajor] }

// Minor returns the minor field of a HDBVersionNumber.
func (v *Version) Minor() uint64 { return v.vn[versionMinor] }

// SPS returns the sps field of a HDBVersionNumber.
func (v *Version) SPS() uint64 { return v.vn[versionRevision] / 10 }

// Revision returns the revision field of a HDBVersionNumber.
func (v *Version) Revision() uint64 { return v.vn[versionRevision] }

// Patch returns the patch field of a HDBVersionNumber.
func (v *Version) Patch() uint64 { return v.vn[versionPatch] }

// BuildID returns the build id field of a HDBVersionNumber.
func (v *Version) BuildID() uint64 { return v.vn[versionBuildID] }

// parseVersion parses a semantic hdb version string field.
func parseVersion(s string) *Version {
	vn := parseVersionNumber(s)
	if vn.isZero() { // hdb 1.00 does not report version
		vn = versionNumberOne
	}

	var feature uint64
	// detect features
	for f, cv := range hdbFeatureAvailability {
		if vn.compare(cv) >= 0 { // v is equal or greater than cv
			feature |= f // add feature
		}
	}
	return &Version{vn: vn, feature: feature}
}

// compare compares the version with a second version v2. The result will be
//
//	0 in case the two versions are equal,
//
// -1 in case version v has lower precedence than c2,
//
//	1 in case version v has higher precedence than c2.
func (v *Version) compare(v2 *Version) int {
	return v.vn.compare(v2.vn)
}

// hasFeature returns true if HDBVersion does support feature - false otherwise.
func (v *Version) hasFeature(feature uint64) bool { return v.feature&feature != 0 }
