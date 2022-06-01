package command

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/helper/namespace"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/version"
	"github.com/hashicorp/vault/vault"
	vaultseal "github.com/hashicorp/vault/vault/seal"
	"github.com/mitchellh/go-testing-interface"
)

func initDevCore(c *ServerCommand, coreConfig *vault.CoreConfig, config *server.Config, core *vault.Core) error {
	if c.flagDev && !c.flagDevSkipInit {

		init, qsOptions, err := c.enableDev(core, coreConfig)
		if err != nil {
			return fmt.Errorf("Error initializing Dev mode: %s", err)
		}

		var plugins, pluginsNotLoaded []string
		if c.flagDevPluginDir != "" && c.flagDevPluginInit {

			f, err := os.Open(c.flagDevPluginDir)
			if err != nil {
				return fmt.Errorf("Error reading plugin dir: %s", err)
			}

			list, err := f.Readdirnames(0)
			f.Close()
			if err != nil {
				return fmt.Errorf("Error listing plugins: %s", err)
			}

			for _, name := range list {
				path := filepath.Join(f.Name(), name)
				if err := c.addPlugin(path, init.RootToken, core); err != nil {
					if !errwrap.Contains(err, vault.ErrPluginBadType.Error()) {
						return fmt.Errorf("Error enabling plugin %s: %s", name, err)
					}
					pluginsNotLoaded = append(pluginsNotLoaded, name)
					continue
				}
				plugins = append(plugins, name)
			}

			sort.Strings(plugins)
		}

		var qw *quiescenceSink
		var qwo sync.Once
		qw = &quiescenceSink{
			t: time.AfterFunc(100*time.Millisecond, func() {
				qwo.Do(func() {
					c.logger.DeregisterSink(qw)
					endpointURL := "http://" + config.Listeners[0].Address

					// Print the big dev mode warning!
					c.UI.Warn(wrapAtLength(
						"WARNING! dev mode is enabled! In this mode, Vault runs entirely " +
							"in-memory and starts unsealed with a single unseal key. The root " +
							"token is already authenticated to the CLI, so you can immediately " +
							"begin using Vault."))
					c.UI.Warn("")

					if c.flagDevQuickStart {
						qsBanner(c, qsOptions, endpointURL)
					}

					c.UI.Warn("You may need to set the following environment variable:")
					c.UI.Warn("")

					if runtime.GOOS == "windows" {
						c.UI.Warn("PowerShell:")
						c.UI.Warn(fmt.Sprintf("    $env:VAULT_ADDR=\"%s\"", endpointURL))
						c.UI.Warn("cmd.exe:")
						c.UI.Warn(fmt.Sprintf("    set VAULT_ADDR=%s", endpointURL))
					} else {
						c.UI.Warn(fmt.Sprintf("    $ export VAULT_ADDR='%s'", endpointURL))
					}

					// Unseal key is not returned if stored shares is supported
					if len(init.SecretShares) > 0 {
						c.UI.Warn("")
						c.UI.Warn(wrapAtLength(
							"The unseal key and root token are displayed below in case you want " +
								"to seal/unseal the Vault or re-authenticate."))
						c.UI.Warn("")
						c.UI.Warn(fmt.Sprintf("Unseal Key: %s", base64.StdEncoding.EncodeToString(init.SecretShares[0])))
					}

					if len(init.RecoveryShares) > 0 {
						c.UI.Warn("")
						c.UI.Warn(wrapAtLength(
							"The recovery key and root token are displayed below in case you want " +
								"to seal/unseal the Vault or re-authenticate."))
						c.UI.Warn("")
						c.UI.Warn(fmt.Sprintf("Recovery Key: %s", base64.StdEncoding.EncodeToString(init.RecoveryShares[0])))
					}

					c.UI.Warn(fmt.Sprintf("Root Token: %s", init.RootToken))

					if len(plugins) > 0 {
						c.UI.Warn("")
						c.UI.Warn(wrapAtLength(
							"The following dev plugins are registered in the catalog:"))
						for _, p := range plugins {
							c.UI.Warn(fmt.Sprintf("    - %s", p))
						}
					}

					if len(pluginsNotLoaded) > 0 {
						c.UI.Warn("")
						c.UI.Warn(wrapAtLength(
							"The following dev plugins FAILED to be registered in the catalog due to unknown type:"))
						for _, p := range pluginsNotLoaded {
							c.UI.Warn(fmt.Sprintf("    - %s", p))
						}
					}

					c.UI.Warn("")
					c.UI.Warn(wrapAtLength(
						"Development mode should NOT be used in production installations!"))
					c.UI.Warn("")
				})
			}),
		}
		c.logger.RegisterSink(qw)
	}
	return nil
}

func (c *ServerCommand) enableThreeNodeDevCluster(base *vault.CoreConfig, info map[string]string, infoKeys []string, devListenAddress, tempDir string) int {
	testCluster := vault.NewTestCluster(&testing.RuntimeT{}, base, &vault.TestClusterOptions{
		HandlerFunc:       vaulthttp.Handler,
		BaseListenAddress: c.flagDevListenAddr,
		Logger:            c.logger,
		TempDir:           tempDir,
	})
	defer c.cleanupGuard.Do(testCluster.Cleanup)

	info["cluster parameters path"] = testCluster.TempDir
	infoKeys = append(infoKeys, "cluster parameters path")

	for i, core := range testCluster.Cores {
		info[fmt.Sprintf("node %d api address", i)] = fmt.Sprintf("https://%s", core.Listeners[0].Address.String())
		infoKeys = append(infoKeys, fmt.Sprintf("node %d api address", i))
	}

	infoKeys = append(infoKeys, "version")
	verInfo := version.GetVersion()
	info["version"] = verInfo.FullVersionNumber(false)
	if verInfo.Revision != "" {
		info["version sha"] = strings.Trim(verInfo.Revision, "'")
		infoKeys = append(infoKeys, "version sha")
	}

	infoKeys = append(infoKeys, "cgo")
	info["cgo"] = "disabled"
	if version.CgoEnabled {
		info["cgo"] = "enabled"
	}

	infoKeys = append(infoKeys, "go version")
	info["go version"] = runtime.Version()

	fipsStatus := getFIPSInfoKey()
	if fipsStatus != "" {
		infoKeys = append(infoKeys, "fips")
		info["fips"] = fipsStatus
	}

	// Server configuration output
	padding := 24

	sort.Strings(infoKeys)
	c.UI.Output("==> Vault server configuration:\n")

	for _, k := range infoKeys {
		c.UI.Output(fmt.Sprintf(
			"%s%s: %s",
			strings.Repeat(" ", padding-len(k)),
			strings.Title(k),
			info[k]))
	}

	c.UI.Output("")

	for _, core := range testCluster.Cores {
		core.Server.Handler = vaulthttp.Handler(&vault.HandlerProperties{
			Core: core.Core,
		})
		core.SetClusterHandler(core.Server.Handler)
	}

	testCluster.Start()

	ctx := namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace)

	if base.DevToken != "" {
		req := &logical.Request{
			ID:          "dev-gen-root",
			Operation:   logical.UpdateOperation,
			ClientToken: testCluster.RootToken,
			Path:        "auth/token/create",
			Data: map[string]interface{}{
				"id":                base.DevToken,
				"policies":          []string{"root"},
				"no_parent":         true,
				"no_default_policy": true,
			},
		}
		resp, err := testCluster.Cores[0].HandleRequest(ctx, req)
		if err != nil {
			c.UI.Error(fmt.Sprintf("failed to create root token with ID %s: %s", base.DevToken, err))
			return 1
		}
		if resp == nil {
			c.UI.Error(fmt.Sprintf("nil response when creating root token with ID %s", base.DevToken))
			return 1
		}
		if resp.Auth == nil {
			c.UI.Error(fmt.Sprintf("nil auth when creating root token with ID %s", base.DevToken))
			return 1
		}

		testCluster.RootToken = resp.Auth.ClientToken

		req.ID = "dev-revoke-init-root"
		req.Path = "auth/token/revoke-self"
		req.Data = nil
		_, err = testCluster.Cores[0].HandleRequest(ctx, req)
		if err != nil {
			c.UI.Output(fmt.Sprintf("failed to revoke initial root token: %s", err))
			return 1
		}
	}

	// Set the token
	tokenHelper, err := c.TokenHelper()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error getting token helper: %s", err))
		return 1
	}
	if err := tokenHelper.Store(testCluster.RootToken); err != nil {
		c.UI.Error(fmt.Sprintf("Error storing in token helper: %s", err))
		return 1
	}

	if err := ioutil.WriteFile(filepath.Join(testCluster.TempDir, "root_token"), []byte(testCluster.RootToken), 0o600); err != nil {
		c.UI.Error(fmt.Sprintf("Error writing token to tempfile: %s", err))
		return 1
	}

	c.UI.Output(fmt.Sprintf(
		"==> Three node dev mode is enabled\n\n" +
			"The unseal key and root token are reproduced below in case you\n" +
			"want to seal/unseal the Vault or play with authentication.\n",
	))

	for i, key := range testCluster.BarrierKeys {
		c.UI.Output(fmt.Sprintf(
			"Unseal Key %d: %s",
			i+1, base64.StdEncoding.EncodeToString(key),
		))
	}

	c.UI.Output(fmt.Sprintf(
		"\nRoot Token: %s\n", testCluster.RootToken,
	))

	c.UI.Output(fmt.Sprintf(
		"\nUseful env vars:\n"+
			"VAULT_TOKEN=%s\n"+
			"VAULT_ADDR=%s\n"+
			"VAULT_CACERT=%s/ca_cert.pem\n",
		testCluster.RootToken,
		testCluster.Cores[0].Client.Address(),
		testCluster.TempDir,
	))

	// Output the header that the server has started
	c.UI.Output("==> Vault server started! Log data will stream in below:\n")

	// Inform any tests that the server is ready
	select {
	case c.startedCh <- struct{}{}:
	default:
	}

	// Release the log gate.
	c.flushLog()

	// Wait for shutdown
	shutdownTriggered := false

	for !shutdownTriggered {
		select {
		case <-c.ShutdownCh:
			c.UI.Output("==> Vault shutdown triggered")

			// Stop the listeners so that we don't process further client requests.
			c.cleanupGuard.Do(testCluster.Cleanup)

			// Finalize will wait until after Vault is sealed, which means the
			// request forwarding listeners will also be closed (and also
			// waited for).
			for _, core := range testCluster.Cores {
				if err := core.Shutdown(); err != nil {
					c.UI.Error(fmt.Sprintf("Error with core shutdown: %s", err))
				}
			}

			shutdownTriggered = true

		case <-c.SighupCh:
			c.UI.Output("==> Vault reload triggered")
			for _, core := range testCluster.Cores {
				if err := c.Reload(core.ReloadFuncsLock, core.ReloadFuncs, nil); err != nil {
					c.UI.Error(fmt.Sprintf("Error(s) were encountered during reload: %s", err))
				}
			}
		}
	}

	return 0
}

func (c *ServerCommand) enableDev(core *vault.Core, coreConfig *vault.CoreConfig) (*vault.InitResult, *quickstartOptions, error) {
	ctx := namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace)

	var recoveryConfig *vault.SealConfig
	barrierConfig := &vault.SealConfig{
		SecretShares:    1,
		SecretThreshold: 1,
	}

	if core.SealAccess().RecoveryKeySupported() {
		recoveryConfig = &vault.SealConfig{
			SecretShares:    1,
			SecretThreshold: 1,
		}
	}

	if core.SealAccess().StoredKeysSupported() != vaultseal.StoredKeysNotSupported {
		barrierConfig.StoredShares = 1
	}

	// Initialize it with a basic single key
	init, err := core.Initialize(ctx, &vault.InitParams{
		BarrierConfig:  barrierConfig,
		RecoveryConfig: recoveryConfig,
	})
	if err != nil {
		return nil, nil, err
	}

	// Handle unseal with stored keys
	if core.SealAccess().StoredKeysSupported() == vaultseal.StoredKeysSupportedGeneric {
		err := core.UnsealWithStoredKeys(ctx)
		if err != nil {
			return nil, nil, err
		}
	} else {
		// Copy the key so that it can be zeroed
		key := make([]byte, len(init.SecretShares[0]))
		copy(key, init.SecretShares[0])

		// Unseal the core
		unsealed, err := core.Unseal(key)
		if err != nil {
			return nil, nil, err
		}
		if !unsealed {
			return nil, nil, fmt.Errorf("failed to unseal Vault for dev mode")
		}
	}

	isLeader, _, _, err := core.Leader()
	if err != nil && err != vault.ErrHANotEnabled {
		return nil, nil, fmt.Errorf("failed to check active status: %w", err)
	}
	if err == nil {
		leaderCount := 5
		for !isLeader {
			if leaderCount == 0 {
				buf := make([]byte, 1<<16)
				runtime.Stack(buf, true)
				return nil, nil, fmt.Errorf("failed to get active status after five seconds; call stack is\n%s", buf)
			}
			time.Sleep(1 * time.Second)
			isLeader, _, _, err = core.Leader()
			if err != nil {
				return nil, nil, fmt.Errorf("failed to check active status: %w", err)
			}
			leaderCount--
		}
	}

	// Generate a dev root token if one is provided in the flag
	if coreConfig.DevToken != "" {
		req := &logical.Request{
			ID:          "dev-gen-root",
			Operation:   logical.UpdateOperation,
			ClientToken: init.RootToken,
			Path:        "auth/token/create",
			Data: map[string]interface{}{
				"id":                coreConfig.DevToken,
				"policies":          []string{"root"},
				"no_parent":         true,
				"no_default_policy": true,
			},
		}
		resp, err := core.HandleRequest(ctx, req)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create root token with ID %q: %w", coreConfig.DevToken, err)
		}
		if resp == nil {
			return nil, nil, fmt.Errorf("nil response when creating root token with ID %q", coreConfig.DevToken)
		}
		if resp.Auth == nil {
			return nil, nil, fmt.Errorf("nil auth when creating root token with ID %q", coreConfig.DevToken)
		}

		init.RootToken = resp.Auth.ClientToken

		req.ID = "dev-revoke-init-root"
		req.Path = "auth/token/revoke-self"
		req.Data = nil
		_, err = core.HandleRequest(ctx, req)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to revoke initial root token: %w", err)
		}
	}

	// Set the token
	if !c.flagDevNoStoreToken {
		tokenHelper, err := c.TokenHelper()
		if err != nil {
			return nil, nil, err
		}
		if err := tokenHelper.Store(init.RootToken); err != nil {
			return nil, nil, err
		}
	}

	kvVer := "2"
	if c.flagDevKVV1 || c.flagDevLeasedKV {
		kvVer = "1"
	}
	req := &logical.Request{
		Operation:   logical.UpdateOperation,
		ClientToken: init.RootToken,
		Path:        "sys/mounts/secret",
		Data: map[string]interface{}{
			"type":        "kv",
			"path":        "secret/",
			"description": "key/value secret storage",
			"options": map[string]string{
				"version": kvVer,
			},
		},
	}
	resp, err := core.HandleRequest(ctx, req)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating default K/V store: %w", err)
	}
	if resp.IsError() {
		return nil, nil, fmt.Errorf("failed to create default K/V store: %w", resp.Error())
	}

	var qsOption *quickstartOptions
	if c.flagDevQuickStart {
		qsOption, err = enableQuickstart(ctx, init, core)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to enable quickstart: %w", err)
		}
	}

	return init, qsOption, nil
}

type quickstartOptions struct {
	identity struct {
		name     string
		entityID string
		policies []string
	}

	authAppRole struct {
		roleName string
		roleID   string
		secretID string
		policies []string
	}

	authUserpass struct {
		username string
		password string
		policies []string
	}

	secrets []qsSecret

	policyMap map[string]string
}

type qsSecret struct {
	key            string
	data           map[string]interface{}
	customMetadata map[string]interface{}
}

func enableQuickstart(ctx context.Context, init *vault.InitResult, core *vault.Core) (*quickstartOptions, error) {
	readOnlyPolicy := "path \"secret/*\" {capabilities = [\"read\"]}"
	readOnlyPolicyName := "read-only"
	crudPolicy := "path \"secret/*\" {capabilities = [\"create\", \"read\", \"update\", \"delete\"]}"
	crudPolicyName := "developer"

	policies := make(map[string]string)
	policies[crudPolicyName] = crudPolicy
	policies[readOnlyPolicyName] = readOnlyPolicy

	options := &quickstartOptions{
		policyMap: policies,

		identity: struct {
			name     string
			entityID string
			policies []string
		}{name: "demo-identity", policies: []string{crudPolicyName}},

		authAppRole: struct {
			roleName string
			roleID   string
			secretID string
			policies []string
		}{roleName: "app-read-only", policies: []string{readOnlyPolicyName}},

		authUserpass: struct {
			username string
			password string
			policies []string
		}{username: "demo-username", password: "superSecretDemo", policies: nil},

		secrets: []qsSecret{{
			key: "sample-secret",
			data: map[string]interface{}{
				"foo": "bar",
				"zip": "zap",
			},
			customMetadata: map[string]interface{}{
				"team": "demo",
			},
		}},
	}

	// add policies
	for name, policy := range options.policyMap {
		if err := devAddPolicy(ctx, init, core, name, policy); err != nil {
			return nil, err
		}
	}

	// userpass
	err := qsEnableUserpass(ctx, init, core, options)
	if err != nil {
		return nil, fmt.Errorf("failed to initalize quickstart userpass: %w", err)
	}

	// app role auth
	err = qsEnableAppRole(ctx, init, core, options)
	if err != nil {
		return nil, fmt.Errorf("failed to initalize quickstart approle: %w", err)
	}

	// identity
	err = qsEnableIdentity(ctx, init, core, options)
	if err != nil {
		return nil, fmt.Errorf("failed to initalize quickstart identity: %w", err)
	}

	// sample secret
	for _, s := range options.secrets {
		err = devAddSecret(ctx, init, core, s)
		if err != nil {
			return nil, fmt.Errorf("failed to initalize quickstart seeded secret: %w", err)
		}
	}

	return options, nil
}

func devAddSecret(ctx context.Context, init *vault.InitResult, core *vault.Core, secret qsSecret) error {
	r := &logical.Request{
		Operation:   logical.UpdateOperation,
		ClientToken: init.RootToken,
		Path:        fmt.Sprintf("secret/data/%s", secret.key),
		Data: map[string]interface{}{
			"data": secret.data,
		},
	}
	resp, err := core.HandleRequest(ctx, r)
	if err != nil {
		return fmt.Errorf("error writing sample secret: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("failed to write sample secret: %w", resp.Error())
	}

	r = &logical.Request{
		Operation:   logical.UpdateOperation,
		ClientToken: init.RootToken,
		Path:        fmt.Sprintf("secret/data/%s", secret.key),
		Data: map[string]interface{}{
			"custom_metadata": secret.customMetadata,
		},
	}
	resp, err = core.HandleRequest(ctx, r)
	if err != nil {
		return fmt.Errorf("error writing sample secret metadata: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("failed to write sample secret metadata: %w", resp.Error())
	}

	return nil
}

func devAddPolicy(ctx context.Context, init *vault.InitResult, core *vault.Core, policyName, policy string) error {
	r := &logical.Request{
		Operation:   logical.UpdateOperation,
		ClientToken: init.RootToken,
		Path:        fmt.Sprintf("sys/policy/%s", policyName),
		Data: map[string]interface{}{
			"name":   policyName,
			"policy": policy,
		},
	}
	resp, err := core.HandleRequest(ctx, r)
	if err != nil {
		return fmt.Errorf("error creating %s policy: %w", policyName, err)
	}
	if resp.IsError() {
		return fmt.Errorf("failed to create %s policy: %w", policyName, resp.Error())
	}

	return nil
}

func qsEnableUserpass(ctx context.Context, init *vault.InitResult, core *vault.Core, options *quickstartOptions) error {
	r := &logical.Request{
		Operation:   logical.UpdateOperation,
		ClientToken: init.RootToken,
		Path:        "sys/auth/userpass",
		Data: map[string]interface{}{
			"type":        "userpass",
			"path":        "userpass/",
			"description": "sample userpass method",
		},
	}
	resp, err := core.HandleRequest(ctx, r)
	if err != nil {
		return fmt.Errorf("error creating userpass: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("failed to create userpass: %w", resp.Error())
	}
	r = &logical.Request{
		Operation:   logical.UpdateOperation,
		ClientToken: init.RootToken,
		Path:        fmt.Sprintf("auth/userpass/users/%s", options.authUserpass.username),
		Data: map[string]interface{}{
			"password": options.authUserpass.password,
		},
	}
	resp, err = core.HandleRequest(ctx, r)
	if err != nil {
		return fmt.Errorf("error creating userpass user: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("failed to create userpass user: %w", resp.Error())
	}

	return nil
}

func qsEnableIdentity(ctx context.Context, init *vault.InitResult, core *vault.Core, options *quickstartOptions) error {
	r := &logical.Request{
		Operation:   logical.CreateOperation,
		ClientToken: init.RootToken,
		Path:        "identity/entity",
		Data: map[string]interface{}{
			"name": options.authAppRole.roleName,
			"metadata": map[string]string{
				"organization": "ACME Inc.",
				"team":         "Squad-2",
			},
			"policies": options.authAppRole.policies,
		},
	}
	resp, err := core.HandleRequest(ctx, r)
	if err != nil {
		return fmt.Errorf("error creating entity: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("failed to create entity: %w", resp.Error())
	}
	id, ok := resp.Data["id"].(string)
	if !ok {
		return fmt.Errorf("failed to read entity id")
	}

	options.identity.entityID = id

	r = &logical.Request{
		Operation:   logical.ReadOperation,
		ClientToken: init.RootToken,
		Path:        "sys/auth",
	}
	resp, err = core.HandleRequest(ctx, r)
	if err != nil {
		return fmt.Errorf("error reading mount accessor: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("failed to read mount accessor: %w", resp.Error())
	}
	up, ok := resp.Data["userpass/"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("failed to read mount accessor")
	}
	accessor, ok := up["accessor"].(string)
	if !ok {
		return fmt.Errorf("failed to read mount accessor")
	}

	r = &logical.Request{
		Operation:   logical.CreateOperation,
		ClientToken: init.RootToken,
		Path:        "identity/entity-alias",
		Data: map[string]interface{}{
			"name":           options.authUserpass.username,
			"canonical_id":   id,
			"mount_accessor": accessor,
			"custom-metadata": map[string]string{
				"account": "demo",
			},
		},
	}
	resp, err = core.HandleRequest(ctx, r)
	if err != nil {
		return fmt.Errorf("error creating entity: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("failed to create entity: %w", resp.Error())
	}

	return nil
}

func qsEnableAppRole(ctx context.Context, init *vault.InitResult, core *vault.Core, options *quickstartOptions) error {
	r := &logical.Request{
		Operation:   logical.UpdateOperation,
		ClientToken: init.RootToken,
		Path:        "sys/auth/approle",
		Data: map[string]interface{}{
			"type":        "approle",
			"path":        "approle/",
			"description": "sample approle method",
		},
	}
	resp, err := core.HandleRequest(ctx, r)
	if err != nil {
		return fmt.Errorf("error creating approle: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("failed to create approle: %w", resp.Error())
	}

	r = &logical.Request{
		Operation:   logical.UpdateOperation,
		ClientToken: init.RootToken,
		Path:        fmt.Sprintf("auth/approle/role/%s", options.authAppRole.roleName),
		Data: map[string]interface{}{
			"role_name":      options.authAppRole.roleName,
			"token_policies": options.authAppRole.policies,
		},
	}
	resp, err = core.HandleRequest(ctx, r)
	if err != nil {
		return fmt.Errorf("error creating approle role: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("failed to create approle role: %w", resp.Error())
	}

	r = &logical.Request{
		Operation:   logical.ReadOperation,
		ClientToken: init.RootToken,
		Path:        fmt.Sprintf("auth/approle/role/%s/role-id", options.authAppRole.roleName),
	}
	resp, err = core.HandleRequest(ctx, r)
	if err != nil {
		return fmt.Errorf("error retrieving approle id: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("failed to retrieve approle id: %w", resp.Error())
	}

	roleID, ok := resp.Data["role_id"].(string)
	if !ok {
		return fmt.Errorf("failed to parse role id")
	}
	options.authAppRole.roleID = roleID

	r = &logical.Request{
		Operation:   logical.UpdateOperation,
		ClientToken: init.RootToken,
		Path:        fmt.Sprintf("auth/approle/role/%s/secret-id", options.authAppRole.roleID),
	}
	resp, err = core.HandleRequest(ctx, r)
	if err != nil {
		return fmt.Errorf("error retrieving approle id: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("failed to retrieve approle id: %w", resp.Error())
	}

	secretID, ok := resp.Data["secret_id"].(string)
	if !ok {
		return fmt.Errorf("failed to parse secret id")
	}

	options.authAppRole.secretID = secretID

	return nil
}

func qsBanner(c *ServerCommand, o *quickstartOptions, e string) {
	if o == nil {
		// no quickstart options seem to be set
		return
	}

	c.UI.Warn(wrapAtLength("A set of quickstart options have been enabled!"))
	c.UI.Warn("Read More: https://www.vaultproject.io/docs/concepts/dev-server")
	c.UI.Warn("")
	c.UI.Warn(wrapAtLength("An identity entity with an associated alias has been created. " +
		"To explore the restricted account, authenticate with the following information"))
	c.UI.Warn("Read More About Identity: https://www.vaultproject.io/docs/concepts/identity")
	c.UI.Warn("Read More About Userpass: https://www.vaultproject.io/docs/auth/userpass")
	c.UI.Warn("Read More About Policies: https://www.vaultproject.io/docs/concepts/policies")
	c.UI.Warn("")

	c.UI.Warn(fmt.Sprintf("               Name: %s", o.identity.name))
	c.UI.Warn(fmt.Sprintf("         Enitity ID: %s", o.identity.entityID))
	c.UI.Warn(fmt.Sprintf("           Policies: %s", o.identity.policies))
	c.UI.Warn("Associated Userpass:")
	c.UI.Warn(fmt.Sprintf("         - Username: %s", o.authUserpass.username))
	c.UI.Warn(fmt.Sprintf("         - Password: %s", o.authUserpass.password))
	c.UI.Warn("")

	c.UI.Warn(wrapAtLength("Applications can login with the following AppRole credentials:"))
	c.UI.Warn("Read More: https://www.vaultproject.io/docs/auth/approle")
	c.UI.Warn("")
	c.UI.Warn(fmt.Sprintf("Login Path: %s/v1/auth/approle/login", e))
	c.UI.Warn(fmt.Sprintf(" Role Name: %s", o.authAppRole.roleName))
	c.UI.Warn(fmt.Sprintf("   Role ID: %s", o.authAppRole.roleID))
	c.UI.Warn(fmt.Sprintf(" Secret ID: %s", o.authAppRole.secretID))
	c.UI.Warn(fmt.Sprintf("  Policies: %s", o.authAppRole.policies))
	c.UI.Warn("")

	c.UI.Warn(wrapAtLength("Polices have been created to enable the following permissions:"))
	c.UI.Warn("")
	for name, policy := range o.policyMap {
		c.UI.Warn(fmt.Sprintf("Policy Name: %s", name))
		c.UI.Warn(fmt.Sprintf("Policy Path: %s/v1/sys/policy/%s", e, name))
		c.UI.Warn(fmt.Sprintf("     Policy: %s", policy))
		c.UI.Warn("")
	}

	c.UI.Warn(wrapAtLength("A sample secret is available at:"))
	c.UI.Warn("")
	c.UI.Warn(fmt.Sprintf("  Secret Path: %s/v1/secret/data/sample-secret", e))
	c.UI.Warn(fmt.Sprintf("Metadata Path: %s/v1/secret/metadata/sample-secret", e))
}
