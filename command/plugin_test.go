package command

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/consts"
)

// testPluginCreate creates a sample plugin in a tempdir and returns the shasum
// and filepath to the plugin.
func testPluginCreate(tb testing.TB, dir, name string) (string, string) {
	tb.Helper()

	pth := dir + "/" + name
	if err := ioutil.WriteFile(pth, nil, 0o755); err != nil {
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

// testPluginCreateAndRegisterVersioned creates a versioned plugin and registers it in the catalog.
func testPluginCreateAndRegisterVersioned(tb testing.TB, client *api.Client, dir, name string, pluginType consts.PluginType) (string, string, string) {
	tb.Helper()

	pth, sha256Sum := testPluginCreate(tb, dir, name)

	if err := client.Sys().RegisterPlugin(&api.RegisterPluginInput{
		Name:    name,
		Type:    pluginType,
		Command: name,
		SHA256:  sha256Sum,
		Version: "v1.0.0",
	}); err != nil {
		tb.Fatal(err)
	}

	return pth, sha256Sum, "v1.0.0"
}
