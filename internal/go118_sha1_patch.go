package internal

import (
	"fmt"
	"os"
	"sync"
	_ "unsafe" // for go:linkname

	goversion "github.com/hashicorp/go-version"
	"github.com/hashicorp/vault/version"
)

const sha1PatchVersionsBefore = "1.12.0"

var patchSha1 sync.Once

//go:linkname debugAllowSHA1 crypto/x509.debugAllowSHA1
var debugAllowSHA1 bool

// PatchSha1 patches Go 1.18+ to allow certificates with signatures containing SHA-1 hashes to be allowed.
// It is safe to call this function multiple times.
// This is necessary to allow Vault 1.10 and 1.11 to work with Go 1.18+ without breaking backwards compatibility
// with these certificates. See https://go.dev/doc/go1.18#sha1 and
// https://developer.hashicorp.com/vault/docs/deprecation/faq#q-what-is-the-impact-of-removing-support-for-x-509-certificates-with-signatures-that-use-sha-1
// for more details.
// TODO: remove when Vault <=1.11 is no longer supported
func PatchSha1() {
	patchSha1.Do(func() {
		// for Go 1.19.4 and later
		godebug := os.Getenv("GODEBUG")
		if godebug != "" {
			godebug += ","
		}
		godebug += "x509sha1=1"
		os.Setenv("GODEBUG", godebug)

		// for Go 1.19.3 and earlier, patch the variable
		patchBefore, err := goversion.NewSemver(sha1PatchVersionsBefore)
		if err != nil {
			panic(err)
		}

		patch := false
		v, err := goversion.NewSemver(version.GetVersion().Version)
		if err == nil {
			patch = v.LessThan(patchBefore)
		} else {
			fmt.Fprintf(os.Stderr, "Cannot parse version %s; going to apply SHA-1 deprecation patch workaround\n", version.GetVersion().Version)
			patch = true
		}

		if patch {
			debugAllowSHA1 = true
		}
	})
}
