package useragent

import (
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestUserAgent(t *testing.T) {
	projectURL = "https://vault-test.com"
	rt = "go5.0"
	versionFunc = func() string { return "1.2.3" }

	act := String()

	exp := "Vault/1.2.3 (+https://vault-test.com; go5.0)"
	if exp != act {
		t.Errorf("expected %q to be %q", act, exp)
	}
}

func TestUserAgentPlugin(t *testing.T) {
	projectURL = "https://vault-test.com"
	rt = "go5.0"
	env := &logical.PluginEnvironment{
		VaultVersion: "1.2.3",
	}
	pluginName := "azure-auth"

	act := PluginString(env, pluginName)

	exp := "Vault/1.2.3 (+https://vault-test.com; azure-auth; go5.0)"
	if exp != act {
		t.Errorf("expected %q to be %q", act, exp)
	}

	pluginName = ""
	act = PluginString(env, pluginName)

	exp = "Vault/1.2.3 (+https://vault-test.com; go5.0)"
	if exp != act {
		t.Errorf("expected %q to be %q", act, exp)
	}
}
