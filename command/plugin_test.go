package command

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/consts"
)

// testPluginDir creates a temporary directory suitable for holding plugins.
// This helper also resolves symlinks to make tests happy on OS X.
func testPluginDir(tb testing.TB) (string, func(tb testing.TB)) {
	tb.Helper()

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		tb.Fatal(err)
	}

	// OSX tempdir are /var, but actually symlinked to /private/var
	dir, err = filepath.EvalSymlinks(dir)
	if err != nil {
		tb.Fatal(err)
	}

	return dir, func(tb testing.TB) {
		if err := os.RemoveAll(dir); err != nil {
			tb.Fatal(err)
		}
	}
}

// testPluginCreate creates a sample plugin in a tempdir and returns the shasum
// and filepath to the plugin.
func testPluginCreate(tb testing.TB, dir, name string) (string, string) {
	tb.Helper()

	pth := dir + "/" + name
	if err := ioutil.WriteFile(pth, nil, 0755); err != nil {
		tb.Fatal(err)
	}

	f, err := os.Open(pth)
	if err != nil {
		tb.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		tb.Fatal(err)
	}
	sha256Sum := fmt.Sprintf("%x", h.Sum(nil))

	return pth, sha256Sum
}

// testPluginCreateAndRegister creates a plugin and registers it in the catalog.
func testPluginCreateAndRegister(tb testing.TB, client *api.Client, dir, name string, pluginType consts.PluginType) (string, string) {
	tb.Helper()

	pth, sha256Sum := testPluginCreate(tb, dir, name)

	if err := client.Sys().RegisterPlugin(&api.RegisterPluginInput{
		Name:    name,
		Type:    pluginType,
		Command: name,
		SHA256:  sha256Sum,
	}); err != nil {
		tb.Fatal(err)
	}

	return pth, sha256Sum
}
