// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package plugincatalog

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
)

// TestAddTestPlugin registers the testFunc as part of the plugin command to the
// plugin catalog. The plugin catalog must be configured with a pluginDirectory.
// NB: The test func you pass in MUST be in the same package as the parent test,
// or the test func won't be compiled into the test binary being run and the output
// will be something like:
// stderr (ignored by go-plugin): "testing: warning: no tests to run"
// stdout: "PASS"
func TestAddTestPlugin(t testing.TB, pluginCatalog *PluginCatalog, name string, pluginType consts.PluginType, version string, testFunc string, env []string) {
	t.Helper()
	if pluginCatalog.directory == "" {
		t.Fatal("plugin catalog must have a plugin directory set to add plugins")
	}
	file, err := os.Open(os.Args[0])
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	fileName := filepath.Base(os.Args[0])

	fi, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}

	// Copy over the file to the temp dir
	dst := filepath.Join(pluginCatalog.directory, fileName)

	// delete the file first to avoid notary failures in macOS
	_ = os.Remove(dst) // ignore error
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fi.Mode())
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()

	if _, err = io.Copy(out, file); err != nil {
		t.Fatal(err)
	}
	err = out.Sync()
	if err != nil {
		t.Fatal(err)
	}
	// Ensure that the file is closed and written. This seems to be
	// necessary on Linux systems.
	out.Close()

	// Copied the file, now seek to the start again to calculate its sha256 hash.
	_, err = file.Seek(0, 0)
	if err != nil {
		t.Fatal(err)
	}

	hash := sha256.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		t.Fatal(err)
	}
	sum := hash.Sum(nil)

	// The flag is a regex, so use ^$ to make sure we only run a single test
	// with an exact match.
	args := []string{fmt.Sprintf("--test.run=^%s$", testFunc)}
	err = pluginCatalog.Set(context.Background(), pluginutil.SetPluginInput{
		Name:    name,
		Type:    pluginType,
		Version: version,
		Command: fileName,
		Args:    args,
		Env:     env,
		Sha256:  sum,
	})
	if err != nil {
		t.Fatal(err)
	}
}
