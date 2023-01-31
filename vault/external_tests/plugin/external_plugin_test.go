package plugin_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/vault"
)

func TestExternalPlugin_AuthMethod(t *testing.T) {
	pluginDir, cleanup := corehelpers.MakeTestPluginDir(t)
	t.Cleanup(func() { cleanup(t) })
	coreConfig := &vault.CoreConfig{
		BuiltinRegistry: corehelpers.NewMockBuiltinRegistry(),
		PluginDirectory: pluginDir,
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		Plugins: &vault.TestPluginConfig{
			Typ:      consts.PluginTypeCredential,
			Versions: []string{""},
		},
		HandlerFunc: vaulthttp.Handler,
	})
	plugin := cluster.Plugins[0]

	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores

	vault.TestWaitActive(t, cores[0].Core)

	client := cores[0].Client
	client.SetToken(cluster.RootToken)

	// Register
	if err := client.Sys().RegisterPlugin(&api.RegisterPluginInput{
		Name:    plugin.Name,
		Type:    api.PluginType(plugin.Typ),
		Command: plugin.Name,
		SHA256:  plugin.Sha256,
		Version: plugin.Version,
	}); err != nil {
		t.Fatal(err)
	}

	pluginPath := fmt.Sprintf("%s-%d", plugin.Name, 0)
	// Enable
	if err := client.Sys().EnableAuthWithOptions(pluginPath, &api.EnableAuthOptions{
		Type: plugin.Name,
	}); err != nil {
		t.Fatal(err)
	}

	// Configure
	_, err := client.Logical().Write("auth/"+pluginPath+"/role/role1", map[string]interface{}{
		"bind_secret_id": "true",
		"period":         "300",
	})
	if err != nil {
		t.Fatal(err)
	}

	secret, err := client.Logical().Write("auth/"+pluginPath+"/role/role1/secret-id", nil)
	if err != nil {
		t.Fatal(err)
	}
	secretID := secret.Data["secret_id"].(string)

	secret, err = client.Logical().Read("auth/" + pluginPath + "/role/role1/role-id")
	if err != nil {
		t.Fatal(err)
	}
	roleID := secret.Data["role_id"].(string)

	// Login - expect SUCCESS
	authMethod, err := approle.NewAppRoleAuth(
		roleID,
		&approle.SecretID{FromString: secretID},
		approle.WithMountPath(pluginPath),
	)
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Auth().Login(context.Background(), authMethod)
	if err != nil {
		t.Fatal(err)
	}

	// Reset root token
	client.SetToken(cluster.RootToken)

	// Reload plugin
	if _, err := client.Sys().ReloadPlugin(&api.ReloadPluginInput{
		Plugin: plugin.Name,
	}); err != nil {
		t.Fatal(err)
	}

	// Login - expect SUCCESS
	resp, err := client.Auth().Login(context.Background(), authMethod)
	if err != nil {
		t.Fatal(err)
	}

	// Renew
	resp, err = client.Auth().Token().RenewSelf(30)
	if err != nil {
		t.Fatal(err)
	}

	// Login - expect SUCCESS
	resp, err = client.Auth().Login(context.Background(), authMethod)
	if err != nil {
		t.Fatal(err)
	}

	revokeToken := resp.Auth.ClientToken
	// Revoke
	if err = client.Auth().Token().RevokeSelf(revokeToken); err != nil {
		t.Fatal(err)
	}

	// Reset root token
	client.SetToken(cluster.RootToken)

	// Lookup - expect FAILURE
	resp, err = client.Auth().Token().Lookup(revokeToken)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	// Reset root token
	client.SetToken(cluster.RootToken)

	// Deregister
	if err := client.Sys().DeregisterPlugin(&api.DeregisterPluginInput{
		Name:    plugin.Name,
		Type:    api.PluginType(plugin.Typ),
		Version: plugin.Version,
	}); err != nil {
		t.Fatal(err)
	}
}

func TestExternalPlugin_SecretsEngine(t *testing.T) {
	pluginDir, cleanup := corehelpers.MakeTestPluginDir(t)
	t.Cleanup(func() { cleanup(t) })
	coreConfig := &vault.CoreConfig{
		BuiltinRegistry: corehelpers.NewMockBuiltinRegistry(),
		PluginDirectory: pluginDir,
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		Plugins: &vault.TestPluginConfig{
			Typ:      consts.PluginTypeSecrets,
			Versions: []string{""},
		},
		HandlerFunc: vaulthttp.Handler,
	})
	plugin := cluster.Plugins[0]

	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores

	vault.TestWaitActive(t, cores[0].Core)

	client := cores[0].Client
	client.SetToken(cluster.RootToken)

	// Register
	if err := client.Sys().RegisterPlugin(&api.RegisterPluginInput{
		Name:    plugin.Name,
		Type:    api.PluginType(plugin.Typ),
		Command: plugin.Name,
		SHA256:  plugin.Sha256,
		Version: plugin.Version,
	}); err != nil {
		t.Fatal(err)
	}
	// Enable
	if err := client.Sys().Mount(plugin.Name, &api.EnableAuthOptions{
		Type: plugin.Name,
	}); err != nil {
		t.Fatal(err)
	}

	// Configure

	// Operations?

	// Reload plugin
	if _, err := client.Sys().ReloadPlugin(&api.ReloadPluginInput{
		Plugin: plugin.Name,
	}); err != nil {
		t.Fatal(err)
	}

	// Operations? - expect SUCCESS

	// Deregister
	if err := client.Sys().DeregisterPlugin(&api.DeregisterPluginInput{
		Name:    plugin.Name,
		Type:    api.PluginType(plugin.Typ),
		Version: plugin.Version,
	}); err != nil {
		t.Fatal(err)
	}
}
