package version

var (
	// The git commit that was compiled. This will be filled in by the compiler.
	GitCommit   string
	GitDescribe string

	// Whether cgo is enabled or not; set at build time
	CgoEnabled bool

	// The actual version will be generated at release time using ldflags
	Version           = "TBD"
	VersionPrerelease = ""
	VersionMetadata   = ""
)
