package packngo

import "runtime/debug"

var (
	// Version of the packngo package. Version will be updated at runtime.
	Version = "(devel)"

	// UserAgent is the default HTTP User-Agent Header value that will be used by NewClient.
	// init() will update the version to match the built version of packngo.
	UserAgent = "packngo/(devel)"
)

const packagePath = "github.com/packethost/packngo"

// init finds packngo in the dependency so the package Version can be properly
// reflected in API UserAgent headers and client introspection
func init() {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	for _, d := range bi.Deps {
		if d.Path == packagePath {
			Version = d.Version
			if d.Replace != nil {
				v := d.Replace.Version
				if v != "" {
					Version = v
				}
			}
			UserAgent = "packngo/" + Version
			break
		}
	}
}
