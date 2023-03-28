// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package pluginhelpers contains testhelpers that don't depend on package
// vault, and thus can be used within vault (as well as elsewhere.)
package pluginhelpers

import (
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sync"

	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/mitchellh/go-testing-interface"
)

var (
	testPluginCacheLock sync.Mutex
	testPluginCache     = map[string][]byte{}
)

type TestPlugin struct {
	Name     string
	Typ      consts.PluginType
	Version  string
	FileName string
	Sha256   string
}

func GetPlugin(t testing.T, typ consts.PluginType) (string, string, string, string) {
	t.Helper()
	var pluginName string
	var pluginType string
	var pluginMain string
	var pluginVersionLocation string

	switch typ {
	case consts.PluginTypeCredential:
		pluginType = "approle"
		pluginName = "vault-plugin-auth-" + pluginType
		pluginMain = filepath.Join("builtin", "credential", pluginType, "cmd", pluginType, "main.go")
		pluginVersionLocation = fmt.Sprintf("github.com/hashicorp/vault/builtin/credential/%s.ReportedVersion", pluginType)
	case consts.PluginTypeSecrets:
		pluginType = "consul"
		pluginName = "vault-plugin-secrets-" + pluginType
		pluginMain = filepath.Join("builtin", "logical", pluginType, "cmd", pluginType, "main.go")
		pluginVersionLocation = fmt.Sprintf("github.com/hashicorp/vault/builtin/logical/%s.ReportedVersion", pluginType)
	case consts.PluginTypeDatabase:
		pluginType = "postgresql"
		pluginName = "vault-plugin-database-" + pluginType
		pluginMain = filepath.Join("plugins", "database", pluginType, fmt.Sprintf("%s-database-plugin", pluginType), "main.go")
		pluginVersionLocation = fmt.Sprintf("github.com/hashicorp/vault/plugins/database/%s.ReportedVersion", pluginType)
	default:
		t.Fatal(typ.String())
	}
	return pluginName, pluginType, pluginMain, pluginVersionLocation
}

// to mount a plugin, we need a working binary plugin, so we compile one here.
// pluginVersion is used to override the plugin's self-reported version
func CompilePlugin(t testing.T, typ consts.PluginType, pluginVersion string, pluginDir string) TestPlugin {
	t.Helper()

	pluginName, pluginType, pluginMain, pluginVersionLocation := GetPlugin(t, typ)

	testPluginCacheLock.Lock()
	defer testPluginCacheLock.Unlock()

	var pluginBytes []byte

	dir := ""
	var err error
	pluginRootDir := "builtin"
	if typ == consts.PluginTypeDatabase {
		pluginRootDir = "plugins"
	}
	for {
		dir, err = os.Getwd()
		if err != nil {
			t.Fatal(err)
		}
		// detect if we are in a subdirectory or the root directory and compensate
		if _, err := os.Stat(pluginRootDir); os.IsNotExist(err) {
			err := os.Chdir("..")
			if err != nil {
				t.Fatal(err)
			}
		} else {
			break
		}
	}

	pluginPath := path.Join(pluginDir, pluginName)
	if pluginVersion != "" {
		pluginPath += "-" + pluginVersion
	}

	key := fmt.Sprintf("%s %s %s", pluginName, pluginType, pluginVersion)
	// cache the compilation to only run once
	var ok bool
	pluginBytes, ok = testPluginCache[key]
	if !ok {
		// we need to compile
		line := []string{"build"}
		if pluginVersion != "" {
			line = append(line, "-ldflags", fmt.Sprintf("-X %s=%s", pluginVersionLocation, pluginVersion))
		}
		line = append(line, "-o", pluginPath, pluginMain)
		cmd := exec.Command("go", line...)
		cmd.Dir = dir
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(fmt.Errorf("error running go build %v output: %s", err, output))
		}
		testPluginCache[key], err = os.ReadFile(pluginPath)
		if err != nil {
			t.Fatal(err)
		}
		pluginBytes = testPluginCache[key]
	}

	// write the cached plugin if necessary
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		err = os.WriteFile(pluginPath, pluginBytes, 0o755)
	}
	if err != nil {
		t.Fatal(err)
	}

	sha := sha256.New()
	_, err = sha.Write(pluginBytes)
	if err != nil {
		t.Fatal(err)
	}
	return TestPlugin{
		Name:     pluginName,
		Typ:      typ,
		Version:  pluginVersion,
		FileName: path.Base(pluginPath),
		Sha256:   fmt.Sprintf("%x", sha.Sum(nil)),
	}
}
