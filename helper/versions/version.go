package versions

import (
	"fmt"
	"runtime/debug"
	"strings"
	"sync"

	semver "github.com/hashicorp/go-version"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/version"
)

var (
	buildInfoOnce         sync.Once // once is used to ensure we only parse build info once.
	buildInfo             *debug.BuildInfo
	DefaultBuiltinVersion = "v" + version.GetVersion().Version + "+builtin.vault"
)

func GetBuiltinVersion(pluginType consts.PluginType, pluginName string) string {
	buildInfoOnce.Do(func() {
		buildInfo, _ = debug.ReadBuildInfo()
	})

	// Should never happen, means the binary was built without Go modules.
	// Fall back to just the Vault version.
	if buildInfo == nil {
		return DefaultBuiltinVersion
	}

	// Vault builtin plugins are all either:
	// a) An external repo within the hashicorp org - return external repo version with +builtin
	// b) Within the Vault repo itself - return Vault version with +builtin.vault
	//
	// The repo names are predictable, but follow slightly different patterns
	// for each plugin type.
	t := pluginType.String()
	switch pluginType {
	case consts.PluginTypeDatabase:
		// Database plugin built-ins are registered as e.g. "postgresql-database-plugin"
		pluginName = strings.TrimSuffix(pluginName, "-database-plugin")
	case consts.PluginTypeSecrets:
		// Repos use "secrets", pluginType.String() is "secret".
		t = "secrets"
	}
	pluginModulePath := fmt.Sprintf("github.com/hashicorp/vault-plugin-%s-%s", t, pluginName)

	for _, dep := range buildInfo.Deps {
		if dep.Path == pluginModulePath {
			return dep.Version + "+builtin"
		}
	}

	return DefaultBuiltinVersion
}

func IsBuiltinVersion(v string) bool {
	semanticVersion, err := semver.NewSemver(v)
	if err != nil {
		return false
	}

	metadataIdentifiers := strings.Split(semanticVersion.Metadata(), ".")
	for _, identifier := range metadataIdentifiers {
		if identifier == "builtin" {
			return true
		}
	}

	return false
}
