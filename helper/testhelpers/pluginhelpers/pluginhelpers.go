// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

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
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/consts"
)

var (
	testPluginCacheLock sync.Mutex
	testPluginCache     = map[string][]byte{}
)

type TestPlugin struct {
	Name        string
	Typ         consts.PluginType
	Version     string
	FileName    string
	Sha256      string
	Image       string
	ImageSha256 string
}

func GetPlugin(t testing.TB, typ consts.PluginType) (string, string, string, string) {
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
func CompilePlugin(t testing.TB, typ consts.PluginType, pluginVersion string, pluginDir string) TestPlugin {
	t.Helper()

	pluginName, pluginType, pluginMain, pluginVersionLocation := GetPlugin(t, typ)

	testPluginCacheLock.Lock()
	defer testPluginCacheLock.Unlock()

	var pluginBytes []byte

	dir := ""
	pluginRootDir := "builtin"
	if typ == consts.PluginTypeDatabase {
		pluginRootDir = "plugins"
	}
	for {
		// So that we can assign to dir without overshadowing the other
		// err variables.
		var getWdErr error
		dir, getWdErr = os.Getwd()
		if getWdErr != nil {
			t.Fatal(getWdErr)
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
		cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
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
	_, statErr := os.Stat(pluginPath)
	if os.IsNotExist(statErr) {
		err := os.WriteFile(pluginPath, pluginBytes, 0o755)
		if err != nil {
			t.Fatal(err)
		}
	} else {
		if statErr != nil {
			t.Fatal(statErr)
		}
	}

	sha := sha256.New()
	_, err := sha.Write(pluginBytes)
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

func BuildPluginContainerImage(t testing.TB, plugin TestPlugin, pluginDir string) (image string, sha256 string) {
	t.Helper()
	ref := plugin.Name
	if plugin.Version != "" {
		ref += ":" + strings.TrimPrefix(plugin.Version, "v")
	} else {
		ref += ":latest"
	}
	args := []string{"build", "--tag=" + ref, "--build-arg=plugin=" + plugin.FileName, "--file=vault/testdata/Dockerfile", pluginDir}
	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(fmt.Errorf("error running docker build %v output: %s", err, output))
	}

	cmd = exec.Command("docker", "images", ref, "--format={{ .ID }}", "--no-trunc")
	id, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(fmt.Errorf("error running docker build %v output: %s", err, output))
	}

	return plugin.Name, strings.TrimSpace(strings.TrimPrefix(string(id), "sha256:"))
}
