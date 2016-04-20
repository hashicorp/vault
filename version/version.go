package version

import (
	"bytes"
	"fmt"
)

// The git commit that was compiled. This will be filled in by the compiler.
var GitCommit string
var GitDescribe string

// The main version number that is being run at the moment.
const Version = "0.5.3"

// A pre-release marker for the version. If this is "" (empty string)
// then it means that it is a final release. Otherwise, this is a pre-release
// such as "dev" (in development), "beta", "rc1", etc.
const VersionPrerelease = "dev"

// VersionInfo
type VersionInfo struct {
	Revision          string
	Version           string
	VersionPrerelease string
}

func GetVersion() *VersionInfo {
	ver := Version
	rel := VersionPrerelease
	if GitDescribe != "" {
		ver = GitDescribe
	}
	if GitDescribe == "" && rel == "" && VersionPrerelease != "" {
		rel = "dev"
	}

	return &VersionInfo{
		Revision:          GitCommit,
		Version:           ver,
		VersionPrerelease: rel,
	}
}

func (c *VersionInfo) String() string {
	var versionString bytes.Buffer

	fmt.Fprintf(&versionString, "Vault v%s", c.Version)
	if c.VersionPrerelease != "" {
		fmt.Fprintf(&versionString, "-%s", c.VersionPrerelease)

		if c.Revision != "" {
			fmt.Fprintf(&versionString, " (%s)", c.Revision)
		}
	}

	return versionString.String()
}
