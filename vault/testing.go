// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	mathrand "math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/go-cleanhttp"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/reloadutil"
	raftlib "github.com/hashicorp/raft"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/audit"
	auditFile "github.com/hashicorp/vault/builtin/audit/file"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/helper/testhelpers/pluginhelpers"
	"github.com/hashicorp/vault/internalshared/configutil"
	v5 "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/helper/testcluster"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"
	physInmem "github.com/hashicorp/vault/sdk/physical/inmem"
	backendplugin "github.com/hashicorp/vault/sdk/plugin"
	"github.com/hashicorp/vault/vault/cluster"
	"github.com/hashicorp/vault/vault/seal"
	"github.com/mitchellh/copystructure"
	"github.com/mitchellh/go-testing-interface"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/net/http2"
)

// This file contains a number of methods that are useful for unit
// tests within other packages.

const (
	testSharedPublicKey = `
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC9i+hFxZHGo6KblVme4zrAcJstR6I0PTJozW286X4WyvPnkMYDQ5mnhEYC7UWCvjoTWbPEXPX7NjhRtwQTGD67bV+lrxgfyzK1JZbUXK4PwgKJvQD+XyyWYMzDgGSQY61KUSqCxymSm/9NZkPU3ElaQ9xQuTzPpztM4ROfb8f2Yv6/ZESZsTo0MTAkp8Pcy+WkioI/uJ1H7zqs0EA4OMY4aDJRu0UtP4rTVeYNEAuRXdX+eH4aW3KMvhzpFTjMbaJHJXlEeUm2SaX5TNQyTOvghCeQILfYIL/Ca2ij8iwCmulwdV6eQGfd4VDu40PvSnmfoaE38o6HaPnX0kUcnKiT
`
	testSharedPrivateKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAvYvoRcWRxqOim5VZnuM6wHCbLUeiND0yaM1tvOl+Fsrz55DG
A0OZp4RGAu1Fgr46E1mzxFz1+zY4UbcEExg+u21fpa8YH8sytSWW1FyuD8ICib0A
/l8slmDMw4BkkGOtSlEqgscpkpv/TWZD1NxJWkPcULk8z6c7TOETn2/H9mL+v2RE
mbE6NDEwJKfD3MvlpIqCP7idR+86rNBAODjGOGgyUbtFLT+K01XmDRALkV3V/nh+
GltyjL4c6RU4zG2iRyV5RHlJtkml+UzUMkzr4IQnkCC32CC/wmtoo/IsAprpcHVe
nkBn3eFQ7uND70p5n6GhN/KOh2j519JFHJyokwIDAQABAoIBAHX7VOvBC3kCN9/x
+aPdup84OE7Z7MvpX6w+WlUhXVugnmsAAVDczhKoUc/WktLLx2huCGhsmKvyVuH+
MioUiE+vx75gm3qGx5xbtmOfALVMRLopjCnJYf6EaFA0ZeQ+NwowNW7Lu0PHmAU8
Z3JiX8IwxTz14DU82buDyewO7v+cEr97AnERe3PUcSTDoUXNaoNxjNpEJkKREY6h
4hAY676RT/GsRcQ8tqe/rnCqPHNd7JGqL+207FK4tJw7daoBjQyijWuB7K5chSal
oPInylM6b13ASXuOAOT/2uSUBWmFVCZPDCmnZxy2SdnJGbsJAMl7Ma3MUlaGvVI+
Tfh1aQkCgYEA4JlNOabTb3z42wz6mz+Nz3JRwbawD+PJXOk5JsSnV7DtPtfgkK9y
6FTQdhnozGWShAvJvc+C4QAihs9AlHXoaBY5bEU7R/8UK/pSqwzam+MmxmhVDV7G
IMQPV0FteoXTaJSikhZ88mETTegI2mik+zleBpVxvfdhE5TR+lq8Br0CgYEA2AwJ
CUD5CYUSj09PluR0HHqamWOrJkKPFPwa+5eiTTCzfBBxImYZh7nXnWuoviXC0sg2
AuvCW+uZ48ygv/D8gcz3j1JfbErKZJuV+TotK9rRtNIF5Ub7qysP7UjyI7zCssVM
kuDd9LfRXaB/qGAHNkcDA8NxmHW3gpln4CFdSY8CgYANs4xwfercHEWaJ1qKagAe
rZyrMpffAEhicJ/Z65lB0jtG4CiE6w8ZeUMWUVJQVcnwYD+4YpZbX4S7sJ0B8Ydy
AhkSr86D/92dKTIt2STk6aCN7gNyQ1vW198PtaAWH1/cO2UHgHOy3ZUt5X/Uwxl9
cex4flln+1Viumts2GgsCQKBgCJH7psgSyPekK5auFdKEr5+Gc/jB8I/Z3K9+g4X
5nH3G1PBTCJYLw7hRzw8W/8oALzvddqKzEFHphiGXK94Lqjt/A4q1OdbCrhiE68D
My21P/dAKB1UYRSs9Y8CNyHCjuZM9jSMJ8vv6vG/SOJPsnVDWVAckAbQDvlTHC9t
O98zAoGAcbW6uFDkrv0XMCpB9Su3KaNXOR0wzag+WIFQRXCcoTvxVi9iYfUReQPi
oOyBJU/HMVvBfv4g+OVFLVgSwwm6owwsouZ0+D/LasbuHqYyqYqdyPJQYzWA2Y+F
+B6f4RoPdSXj24JHPg/ioRxjaj094UXJxua2yfkcecGNEuBQHSs=
-----END RSA PRIVATE KEY-----
`
)

// TestCore returns a pure in-memory, uninitialized core for testing.
func TestCore(t testing.T) *Core {
	return TestCoreWithSeal(t, nil, false)
}

// TestCoreRaw returns a pure in-memory, uninitialized core for testing. The raw
// storage endpoints are enabled with this core.
func TestCoreRaw(t testing.T) *Core {
	return TestCoreWithSeal(t, nil, true)
}

// TestCoreNewSeal returns a pure in-memory, uninitialized core with
// the new seal configuration.
func TestCoreNewSeal(t testing.T) *Core {
	seal := NewTestSeal(t, nil)
	return TestCoreWithSeal(t, seal, false)
}

// TestCoreWithConfig returns a pure in-memory, uninitialized core with the
// specified core configurations overridden for testing.
func TestCoreWithConfig(t testing.T, conf *CoreConfig) *Core {
	return TestCoreWithSealAndUI(t, conf)
}

// TestCoreWithSeal returns a pure in-memory, uninitialized core with the
// specified seal for testing.
func TestCoreWithSeal(t testing.T, testSeal Seal, enableRaw bool) *Core {
	conf := &CoreConfig{
		Seal:            testSeal,
		EnableUI:        false,
		EnableRaw:       enableRaw,
		BuiltinRegistry: corehelpers.NewMockBuiltinRegistry(),
		AuditBackends: map[string]audit.Factory{
			"file": auditFile.Factory,
		},
	}
	return TestCoreWithSealAndUI(t, conf)
}

func TestCoreWithCustomResponseHeaderAndUI(t testing.T, CustomResponseHeaders map[string]map[string]string, enableUI bool) (*Core, [][]byte, string) {
	confRaw := &server.Config{
		SharedConfig: &configutil.SharedConfig{
			Listeners: []*configutil.Listener{
				{
					Type:                  "tcp",
					Address:               "127.0.0.1",
					CustomResponseHeaders: CustomResponseHeaders,
				},
			},
			DisableMlock: true,
		},
	}
	conf := &CoreConfig{
		RawConfig:       confRaw,
		EnableUI:        enableUI,
		EnableRaw:       true,
		BuiltinRegistry: corehelpers.NewMockBuiltinRegistry(),
	}
	core := TestCoreWithSealAndUI(t, conf)
	return testCoreUnsealed(t, core)
}

func TestCoreUI(t testing.T, enableUI bool) *Core {
	conf := &CoreConfig{
		EnableUI:        enableUI,
		EnableRaw:       true,
		BuiltinRegistry: corehelpers.NewMockBuiltinRegistry(),
	}
	return TestCoreWithSealAndUI(t, conf)
}

func TestCoreWithSealAndUI(t testing.T, opts *CoreConfig) *Core {
	c := TestCoreWithSealAndUINoCleanup(t, opts)

	t.Cleanup(func() {
		defer func() {
			if r := recover(); r != nil {
				t.Log("panic closing core during cleanup", "panic", r)
			}
		}()
		err := c.ShutdownWait()
		if err != nil {
			t.Logf("shutdown returned error: %v", err)
		}
		if tl, ok := c.Logger().(*corehelpers.TestLogger); ok {
			tl.StopLogging()
		}
	})
	return c
}

func TestCoreWithSealAndUINoCleanup(t testing.T, opts *CoreConfig) *Core {
	logger := corehelpers.NewTestLogger(t)
	physicalBackend, err := physInmem.NewInmem(nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	errInjector := physical.NewErrorInjector(physicalBackend, 0, logger)

	// Start off with base test core config
	conf := testCoreConfig(t, errInjector, logger)

	corehelpers.RegisterSubloggerAdder(logger, conf)

	// Override config values with ones that gets passed in
	conf.EnableUI = opts.EnableUI
	conf.EnableRaw = opts.EnableRaw
	conf.EnableIntrospection = opts.EnableIntrospection
	conf.Seal = opts.Seal
	conf.LicensingConfig = opts.LicensingConfig
	conf.DisableKeyEncodingChecks = opts.DisableKeyEncodingChecks
	conf.MetricsHelper = opts.MetricsHelper
	conf.MetricSink = opts.MetricSink
	conf.NumExpirationWorkers = numExpirationWorkersTest
	conf.RawConfig = opts.RawConfig
	conf.EnableResponseHeaderHostname = opts.EnableResponseHeaderHostname
	conf.DisableSSCTokens = opts.DisableSSCTokens
	conf.PluginDirectory = opts.PluginDirectory
	conf.DetectDeadlocks = opts.DetectDeadlocks
	conf.Experiments = opts.Experiments
	conf.CensusAgent = opts.CensusAgent
	conf.AdministrativeNamespacePath = opts.AdministrativeNamespacePath
	conf.AllLoggers = logger.AllLoggers

	if opts.Logger != nil {
		conf.Logger = opts.Logger
	}

	if opts.RedirectAddr != "" {
		conf.RedirectAddr = opts.RedirectAddr
	}

	for k, v := range opts.LogicalBackends {
		conf.LogicalBackends[k] = v
	}
	for k, v := range opts.CredentialBackends {
		conf.CredentialBackends[k] = v
	}

	for k, v := range opts.AuditBackends {
		conf.AuditBackends[k] = v
	}
	if opts.RollbackPeriod != time.Duration(0) {
		conf.RollbackPeriod = opts.RollbackPeriod
	}

	conf.ActivityLogConfig = opts.ActivityLogConfig
	testApplyEntBaseConfig(conf, opts)

	c, err := NewCore(conf)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Switch the SubloggerHook over to core
	corehelpers.RegisterSubloggerAdder(logger, c)
	return c
}

func testCoreConfig(t testing.T, physicalBackend physical.Backend, logger log.Logger) *CoreConfig {
	t.Helper()
	noopAudits := map[string]audit.Factory{
		"noop": corehelpers.NoopAuditFactory(nil),
	}

	noopBackends := make(map[string]logical.Factory)
	noopBackends["noop"] = func(ctx context.Context, config *logical.BackendConfig) (logical.Backend, error) {
		b := new(framework.Backend)
		b.Setup(ctx, config)
		b.BackendType = logical.TypeCredential
		return b, nil
	}
	noopBackends["http"] = func(ctx context.Context, config *logical.BackendConfig) (logical.Backend, error) {
		return new(rawHTTP), nil
	}

	credentialBackends := make(map[string]logical.Factory)
	for backendName, backendFactory := range noopBackends {
		credentialBackends[backendName] = backendFactory
	}
	for backendName, backendFactory := range testCredentialBackends {
		credentialBackends[backendName] = backendFactory
	}

	logicalBackends := make(map[string]logical.Factory)
	for backendName, backendFactory := range noopBackends {
		logicalBackends[backendName] = backendFactory
	}

	logicalBackends["kv"] = LeasedPassthroughBackendFactory
	for backendName, backendFactory := range testLogicalBackends {
		logicalBackends[backendName] = backendFactory
	}

	conf := &CoreConfig{
		Physical:           physicalBackend,
		AuditBackends:      noopAudits,
		LogicalBackends:    logicalBackends,
		CredentialBackends: credentialBackends,
		DisableMlock:       true,
		Logger:             logger,
		BuiltinRegistry:    corehelpers.NewMockBuiltinRegistry(),
	}

	return conf
}

// TestCoreInit initializes the core with a single key, and returns
// the key that must be used to unseal the core and a root token.
func TestCoreInit(t testing.T, core *Core) ([][]byte, string) {
	t.Helper()
	secretShares, _, root := TestCoreInitClusterWrapperSetup(t, core, nil)
	return secretShares, root
}

func TestCoreInitClusterWrapperSetup(t testing.T, core *Core, handler http.Handler) ([][]byte, [][]byte, string) {
	t.Helper()
	core.SetClusterHandler(handler)

	barrierConfig := &SealConfig{
		SecretShares:    3,
		SecretThreshold: 3,
	}

	switch core.seal.StoredKeysSupported() {
	case seal.StoredKeysNotSupported:
		barrierConfig.StoredShares = 0
	default:
		barrierConfig.StoredShares = 1
	}

	recoveryConfig := &SealConfig{
		SecretShares:    3,
		SecretThreshold: 3,
	}

	initParams := &InitParams{
		BarrierConfig:  barrierConfig,
		RecoveryConfig: recoveryConfig,
	}
	if core.seal.StoredKeysSupported() == seal.StoredKeysNotSupported {
		initParams.LegacyShamirSeal = true
	}
	result, err := core.Initialize(context.Background(), initParams)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	innerToken, err := core.DecodeSSCToken(result.RootToken)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	return result.SecretShares, result.RecoveryShares, innerToken
}

func TestCoreUnseal(core *Core, key []byte) (bool, error) {
	return core.Unseal(key)
}

func TestCoreSeal(core *Core) error {
	return core.sealInternal()
}

// TestCoreUnsealed returns a pure in-memory core that is already
// initialized and unsealed.
func TestCoreUnsealed(t testing.T) (*Core, [][]byte, string) {
	t.Helper()
	core := TestCore(t)
	return testCoreUnsealed(t, core)
}

func SetupMetrics(conf *CoreConfig) *metrics.InmemSink {
	inmemSink := metrics.NewInmemSink(1000000*time.Hour, 2000000*time.Hour)
	conf.MetricSink = metricsutil.NewClusterMetricSink("test-cluster", inmemSink)
	conf.MetricsHelper = metricsutil.NewMetricsHelper(inmemSink, false)
	return inmemSink
}

func TestCoreUnsealedWithMetrics(t testing.T) (*Core, [][]byte, string, *metrics.InmemSink) {
	t.Helper()
	conf := &CoreConfig{
		BuiltinRegistry: corehelpers.NewMockBuiltinRegistry(),
	}
	sink := SetupMetrics(conf)
	core, keys, root := testCoreUnsealed(t, TestCoreWithSealAndUI(t, conf))
	return core, keys, root, sink
}

// TestCoreUnsealedRaw returns a pure in-memory core that is already
// initialized, unsealed, and with raw endpoints enabled.
func TestCoreUnsealedRaw(t testing.T) (*Core, [][]byte, string) {
	t.Helper()
	core := TestCoreRaw(t)
	return testCoreUnsealed(t, core)
}

// TestCoreUnsealedWithConfig returns a pure in-memory core that is already
// initialized, unsealed, with the any provided core config values overridden.
func TestCoreUnsealedWithConfig(t testing.T, conf *CoreConfig) (*Core, [][]byte, string) {
	t.Helper()
	core := TestCoreWithConfig(t, conf)
	return testCoreUnsealed(t, core)
}

func testCoreUnsealed(t testing.T, core *Core) (*Core, [][]byte, string) {
	t.Helper()
	token, keys := TestInitUnsealCore(t, core)

	testCoreAddSecretMount(t, core, token)
	return core, keys, token
}

func TestInitUnsealCore(t testing.T, core *Core) (string, [][]byte) {
	keys, token := TestCoreInit(t, core)
	for _, key := range keys {
		if _, err := TestCoreUnseal(core, TestKeyCopy(key)); err != nil {
			t.Fatalf("unseal err: %s", err)
		}
	}
	if core.Sealed() {
		t.Fatal("should not be sealed")
	}

	return token, keys
}

func testCoreAddSecretMount(t testing.T, core *Core, token string) {
	kvReq := &logical.Request{
		Operation:   logical.UpdateOperation,
		ClientToken: token,
		Path:        "sys/mounts/secret",
		Data: map[string]interface{}{
			"type":        "kv",
			"path":        "secret/",
			"description": "key/value secret storage",
			"options": map[string]string{
				"version": "1",
			},
		},
	}
	resp, err := core.HandleRequest(namespace.RootContext(nil), kvReq)
	if err != nil {
		t.Fatal(err)
	}
	if resp.IsError() {
		t.Fatal(err)
	}
}

func TestCoreUnsealedBackend(t testing.T, backend physical.Backend) (*Core, [][]byte, string) {
	t.Helper()
	logger := corehelpers.NewTestLogger(t)
	conf := testCoreConfig(t, backend, logger)
	conf.Seal = NewTestSeal(t, nil)
	conf.NumExpirationWorkers = numExpirationWorkersTest

	core, err := NewCore(conf)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	keys, token := TestCoreInit(t, core)
	for _, key := range keys {
		if _, err := TestCoreUnseal(core, TestKeyCopy(key)); err != nil {
			t.Fatalf("unseal err: %s", err)
		}
	}

	if err := core.UnsealWithStoredKeys(context.Background()); err != nil {
		t.Fatal(err)
	}

	if core.Sealed() {
		t.Fatal("should not be sealed")
	}

	t.Cleanup(func() {
		defer func() {
			if r := recover(); r != nil {
				t.Log("panic closing core during cleanup", "panic", r)
			}
		}()
		err := core.ShutdownWait()
		if err != nil {
			t.Logf("shutdown returned error: %v", err)
		}
	})

	return core, keys, token
}

// TestKeyCopy is a silly little function to just copy the key so that
// it can be used with Unseal easily.
func TestKeyCopy(key []byte) []byte {
	result := make([]byte, len(key))
	copy(result, key)
	return result
}

func TestDynamicSystemView(c *Core, ns *namespace.Namespace) *dynamicSystemView {
	me := &MountEntry{
		Config: MountConfig{
			DefaultLeaseTTL: 24 * time.Hour,
			MaxLeaseTTL:     2 * 24 * time.Hour,
		},
		NamespaceID: namespace.RootNamespace.ID,
		namespace:   namespace.RootNamespace,
	}

	if ns != nil {
		me.NamespaceID = ns.ID
		me.namespace = ns
	}

	return &dynamicSystemView{c, me, c.perfStandby}
}

// TestAddTestPlugin registers the testFunc as part of the plugin command to the
// plugin catalog. If provided, uses tmpDir as the plugin directory.
// NB: The test func you pass in MUST be in the same package as the parent test,
// or the test func won't be compiled into the test binary being run and the output
// will be something like:
// stderr (ignored by go-plugin): "testing: warning: no tests to run"
// stdout: "PASS"
func TestAddTestPlugin(t testing.T, c *Core, name string, pluginType consts.PluginType, version string, testFunc string, env []string, tempDir string) {
	file, err := os.Open(os.Args[0])
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	dirPath := filepath.Dir(os.Args[0])
	fileName := filepath.Base(os.Args[0])

	if tempDir != "" {
		fi, err := file.Stat()
		if err != nil {
			t.Fatal(err)
		}

		// Copy over the file to the temp dir
		dst := filepath.Join(tempDir, fileName)

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

		dirPath = tempDir
	}

	// Determine plugin directory full path, evaluating potential symlink path
	fullPath, err := filepath.EvalSymlinks(dirPath)
	if err != nil {
		t.Fatal(err)
	}

	reader, err := os.Open(filepath.Join(fullPath, fileName))
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	// Find out the sha256
	hash := sha256.New()

	_, err = io.Copy(hash, reader)
	if err != nil {
		t.Fatal(err)
	}

	sum := hash.Sum(nil)

	// Set core's plugin directory and plugin catalog directory
	c.pluginDirectory = fullPath
	c.pluginCatalog.directory = fullPath

	args := []string{fmt.Sprintf("--test.run=%s", testFunc)}
	err = c.pluginCatalog.Set(context.Background(), pluginutil.SetPluginInput{
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

// TestRunTestPlugin runs the testFunc which has already been registered to the
// plugin catalog and returns a pluginClient. This can be called after calling
// TestAddTestPlugin.
func TestRunTestPlugin(t testing.T, c *Core, pluginType consts.PluginType, pluginName string) *pluginClient {
	t.Helper()
	config := TestPluginClientConfig(c, pluginType, pluginName)
	client, err := c.pluginCatalog.NewPluginClient(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	return client
}

func TestPluginClientConfig(c *Core, pluginType consts.PluginType, pluginName string) pluginutil.PluginClientConfig {
	dsv := TestDynamicSystemView(c, nil)
	switch pluginType {
	case consts.PluginTypeCredential, consts.PluginTypeSecrets:
		return pluginutil.PluginClientConfig{
			Name:            pluginName,
			PluginType:      pluginType,
			PluginSets:      backendplugin.PluginSet,
			HandshakeConfig: backendplugin.HandshakeConfig,
			Logger:          log.NewNullLogger(),
			AutoMTLS:        true,
			IsMetadataMode:  false,
			Wrapper:         dsv,
		}
	case consts.PluginTypeDatabase:
		return pluginutil.PluginClientConfig{
			Name:            pluginName,
			PluginType:      pluginType,
			PluginSets:      v5.PluginSets,
			HandshakeConfig: v5.HandshakeConfig,
			Logger:          log.NewNullLogger(),
			AutoMTLS:        true,
			IsMetadataMode:  false,
			Wrapper:         dsv,
		}
	}
	return pluginutil.PluginClientConfig{}
}

var (
	testLogicalBackends    = map[string]logical.Factory{}
	testCredentialBackends = map[string]logical.Factory{}
)

// This adds a credential backend for the test core. This needs to be
// invoked before the test core is created.
func AddTestCredentialBackend(name string, factory logical.Factory) error {
	if name == "" {
		return fmt.Errorf("missing backend name")
	}
	if factory == nil {
		return fmt.Errorf("missing backend factory function")
	}
	testCredentialBackends[name] = factory
	return nil
}

// This adds a logical backend for the test core. This needs to be
// invoked before the test core is created.
func AddTestLogicalBackend(name string, factory logical.Factory) error {
	if name == "" {
		return fmt.Errorf("missing backend name")
	}
	if factory == nil {
		return fmt.Errorf("missing backend factory function")
	}
	testLogicalBackends[name] = factory
	return nil
}

type rawHTTP struct{}

func (n *rawHTTP) HandleRequest(ctx context.Context, req *logical.Request) (*logical.Response, error) {
	return &logical.Response{
		Data: map[string]interface{}{
			logical.HTTPStatusCode:  200,
			logical.HTTPContentType: "plain/text",
			logical.HTTPRawBody:     []byte("hello world"),
		},
	}, nil
}

func (n *rawHTTP) HandleExistenceCheck(ctx context.Context, req *logical.Request) (bool, bool, error) {
	return false, false, nil
}

func (n *rawHTTP) SpecialPaths() *logical.Paths {
	return &logical.Paths{Unauthenticated: []string{"*"}}
}

func (n *rawHTTP) System() logical.SystemView {
	return logical.StaticSystemView{
		DefaultLeaseTTLVal: time.Hour * 24,
		MaxLeaseTTLVal:     time.Hour * 24 * 32,
	}
}

func (n *rawHTTP) Logger() log.Logger {
	return logging.NewVaultLogger(log.Trace)
}

func (n *rawHTTP) Cleanup(ctx context.Context) {
	// noop
}

func (n *rawHTTP) Initialize(ctx context.Context, req *logical.InitializationRequest) error {
	return nil
}

func (n *rawHTTP) InvalidateKey(context.Context, string) {
	// noop
}

func (n *rawHTTP) Setup(ctx context.Context, config *logical.BackendConfig) error {
	// noop
	return nil
}

func (n *rawHTTP) Type() logical.BackendType {
	return logical.TypeLogical
}

func GenerateRandBytes(length int) ([]byte, error) {
	if length < 0 {
		return nil, fmt.Errorf("length must be >= 0")
	}

	buf := make([]byte, length)
	if length == 0 {
		return buf, nil
	}

	n, err := rand.Read(buf)
	if err != nil {
		return nil, err
	}
	if n != length {
		return nil, fmt.Errorf("unable to read %d bytes; only read %d", length, n)
	}

	return buf, nil
}

func TestWaitActive(t testing.T, core *Core) {
	t.Helper()
	if err := TestWaitActiveWithError(core); err != nil {
		t.Fatal(err)
	}
}

func TestWaitActiveForwardingReady(t testing.T, core *Core) {
	TestWaitActive(t, core)

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if _, ok := core.getClusterListener().Handler(consts.RequestForwardingALPN); ok {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatal("timed out waiting for request forwarding handler to be registered")
}

func TestWaitActiveWithError(core *Core) error {
	start := time.Now()
	var standby bool
	var err error
	for time.Now().Sub(start) < 30*time.Second {
		standby, err = core.Standby()
		if err != nil {
			return err
		}
		if !standby {
			break
		}
	}
	if standby {
		return errors.New("should not be in standby mode")
	}
	return nil
}

type TestCluster struct {
	BarrierKeys        [][]byte
	RecoveryKeys       [][]byte
	CACert             *x509.Certificate
	CACertBytes        []byte
	CACertPEM          []byte
	CACertPEMFile      string
	CAKey              *ecdsa.PrivateKey
	CAKeyPEM           []byte
	Cores              []*TestClusterCore
	ID                 string
	Plugins            []pluginhelpers.TestPlugin
	RootToken          string
	RootCAs            *x509.CertPool
	TempDir            string
	ClientAuthRequired bool
	Logger             log.Logger
	CleanupFunc        func()
	SetupFunc          func()

	cleanupFuncs      []func()
	base              *CoreConfig
	LicensePublicKey  ed25519.PublicKey
	LicensePrivateKey ed25519.PrivateKey
	opts              *TestClusterOptions
}

func (c *TestCluster) SetRootToken(token string) {
	c.RootToken = token
}

func (c *TestCluster) Start() {
}

func (c *TestCluster) start(t testing.T) {
	t.Helper()
	for i, core := range c.Cores {
		if core.Server != nil {
			for _, ln := range core.Listeners {
				c.Logger.Info("starting listener for test core", "core", i, "port", ln.Address.Port)
				go core.Server.Serve(ln)
			}
		}
	}
	if c.SetupFunc != nil {
		c.SetupFunc()
	}

	if c.opts != nil && c.opts.SkipInit {
		// SkipInit implies that vault may not be ready to service requests, or that
		// we're restarting a cluster from an existing storage.
		return
	}

	activeCore := -1
WAITACTIVE:
	for i := 0; i < 60; i++ {
		for i, core := range c.Cores {
			if standby, _ := core.Core.Standby(); !standby {
				activeCore = i
				break WAITACTIVE
			}
		}

		time.Sleep(time.Second)
	}
	if activeCore == -1 {
		t.Fatalf("no core became active")
	}

	switch {
	case c.opts == nil:
	case c.opts.NoDefaultQuotas:
	case c.opts.HandlerFunc == nil:
	// If no HandlerFunc is provided that means that we can't actually do
	// regular vault requests.
	case reflect.TypeOf(c.opts.HandlerFunc).PkgPath() != "github.com/hashicorp/vault/http":
	case reflect.TypeOf(c.opts.HandlerFunc).Name() != "Handler":
	default:
		cli := c.Cores[activeCore].Client
		_, err := cli.Logical().Write("sys/quotas/rate-limit/rl-NewTestCluster", map[string]interface{}{
			"rate": 1000000,
		})
		if err != nil {
			t.Fatalf("error setting up global rate limit quota: %v", err)
		}
		if constants.IsEnterprise {
			_, err = cli.Logical().Write("sys/quotas/lease-count/lc-NewTestCluster", map[string]interface{}{
				"max_leases": 1000000,
			})
			if err != nil {
				t.Fatalf("error setting up global lease count quota: %v", err)
			}
		}
	}
}

// UnsealCores uses the cluster barrier keys to unseal the test cluster cores
func (c *TestCluster) UnsealCores(t testing.T) {
	t.Helper()
	if err := c.UnsealCoresWithError(false); err != nil {
		t.Fatal(err)
	}
}

func (c *TestCluster) UnsealCoresWithError(useStoredKeys bool) error {
	unseal := func(core *Core) error {
		for _, key := range c.BarrierKeys {
			if _, err := core.Unseal(TestKeyCopy(key)); err != nil {
				return err
			}
		}
		return nil
	}
	if useStoredKeys {
		unseal = func(core *Core) error {
			return core.UnsealWithStoredKeys(context.Background())
		}
	}

	// Unseal first core
	if err := unseal(c.Cores[0].Core); err != nil {
		return fmt.Errorf("unseal core %d err: %s", 0, err)
	}

	// Verify unsealed
	if c.Cores[0].Sealed() {
		return fmt.Errorf("should not be sealed")
	}

	if err := TestWaitActiveWithError(c.Cores[0].Core); err != nil {
		return err
	}

	// Unseal other cores
	for i := 1; i < len(c.Cores); i++ {
		if err := unseal(c.Cores[i].Core); err != nil {
			return fmt.Errorf("unseal core %d err: %s", i, err)
		}
	}

	// Let them come fully up to standby
	time.Sleep(2 * time.Second)

	// Ensure cluster connection info is populated.
	// Other cores should not come up as leaders.
	for i := 1; i < len(c.Cores); i++ {
		isLeader, _, _, err := c.Cores[i].Leader()
		if err != nil {
			return err
		}
		if isLeader {
			return fmt.Errorf("core[%d] should not be leader", i)
		}
	}

	return nil
}

func (c *TestCluster) UnsealCore(t testing.T, core *TestClusterCore) {
	err := c.AttemptUnsealCore(core)
	if err != nil {
		t.Fatal(err)
	}
}

func (c *TestCluster) AttemptUnsealCore(core *TestClusterCore) error {
	var keys [][]byte
	if core.seal.RecoveryKeySupported() {
		keys = c.RecoveryKeys
	} else {
		keys = c.BarrierKeys
	}
	for _, key := range keys {
		if _, err := core.Core.Unseal(TestKeyCopy(key)); err != nil {
			return fmt.Errorf("unseal err: %w", err)
		}
	}
	return nil
}

func (c *TestCluster) UnsealCoreWithStoredKeys(t testing.T, core *TestClusterCore) {
	t.Helper()
	if err := core.UnsealWithStoredKeys(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func (c *TestCluster) EnsureCoresSealed(t testing.T) {
	t.Helper()
	if err := c.ensureCoresSealed(); err != nil {
		t.Fatal(err)
	}
}

func (c *TestClusterCore) Seal(t testing.T) {
	t.Helper()
	if err := c.Core.sealInternal(); err != nil {
		t.Fatal(err)
	}
}

func (c *TestClusterCore) LogicalStorage() logical.Storage {
	return c.barrier
}

func (c *TestClusterCore) stop() error {
	c.Logger().Info("stopping vault test core")

	if c.Listeners != nil {
		for _, ln := range c.Listeners {
			ln.Close()
		}
		c.Logger().Info("listeners successfully shut down")
	}
	if c.licensingStopCh != nil {
		close(c.licensingStopCh)
		c.licensingStopCh = nil
	}

	if err := c.Shutdown(); err != nil {
		return err
	}
	timeout := time.Now().Add(60 * time.Second)
	for {
		if time.Now().After(timeout) {
			return errors.New("timeout waiting for core to seal")
		}
		if c.Sealed() {
			break
		}
		time.Sleep(250 * time.Millisecond)
	}

	c.Logger().Info("vault test core stopped")
	return nil
}

func (c *TestClusterCore) StopAutomaticRollbacks() {
	c.rollback.StopTicker()
}

func (c *TestClusterCore) GrabRollbackLock() {
	// Ensure we don't hold this lock while there are in flight rollbacks.
	c.rollback.inflightAll.Wait()
	c.rollback.inflightLock.Lock()
}

func (c *TestClusterCore) ReleaseRollbackLock() {
	c.rollback.inflightLock.Unlock()
}

func (c *TestClusterCore) TriggerRollbacks() {
	c.rollback.triggerRollbacks()
}

func (c *TestClusterCore) TLSConfig() *tls.Config {
	return c.tlsConfig.Clone()
}

func (c *TestClusterCore) ClusterListener() *cluster.Listener {
	return c.getClusterListener()
}

func (c *TestCluster) Cleanup() {
	c.Logger.Info("cleaning up vault cluster")
	if tl, ok := c.Logger.(*corehelpers.TestLogger); ok {
		tl.StopLogging()
	}

	wg := &sync.WaitGroup{}
	for _, core := range c.Cores {
		wg.Add(1)
		lc := core

		go func() {
			defer wg.Done()
			if err := lc.stop(); err != nil {
				// Note that this log won't be seen if using TestLogger, due to
				// the above call to StopLogging.
				lc.Logger().Error("error during cleanup", "error", err)
			}
		}()
	}

	wg.Wait()

	// Remove any temp dir that exists
	if c.TempDir != "" {
		os.RemoveAll(c.TempDir)
	}

	// Give time to actually shut down/clean up before the next test
	time.Sleep(time.Second)
	if c.CleanupFunc != nil {
		c.CleanupFunc()
	}
}

func (c *TestCluster) ensureCoresSealed() error {
	for _, core := range c.Cores {
		if err := core.Shutdown(); err != nil {
			return err
		}
		timeout := time.Now().Add(60 * time.Second)
		for {
			if time.Now().After(timeout) {
				return fmt.Errorf("timeout waiting for core to seal")
			}
			if core.Sealed() {
				break
			}
			time.Sleep(250 * time.Millisecond)
		}
	}
	return nil
}

func SetReplicationFailureMode(core *TestClusterCore, mode uint32) {
	atomic.StoreUint32(core.Core.replicationFailure, mode)
}

type TestListener struct {
	net.Listener
	Address *net.TCPAddr
}

type TestClusterCore struct {
	*Core
	CoreConfig           *CoreConfig
	Client               *api.Client
	Handler              http.Handler
	Address              *net.TCPAddr
	Listeners            []*TestListener
	ReloadFuncs          *map[string][]reloadutil.ReloadFunc
	ReloadFuncsLock      *sync.RWMutex
	Server               *http.Server
	ServerCert           *x509.Certificate
	ServerCertBytes      []byte
	ServerCertPEM        []byte
	ServerKey            *ecdsa.PrivateKey
	ServerKeyPEM         []byte
	tlsConfig            *tls.Config
	UnderlyingStorage    physical.Backend
	UnderlyingRawStorage physical.Backend
	UnderlyingHAStorage  physical.HABackend
	Barrier              SecurityBarrier
	NodeID               string
}

type PhysicalBackendBundle struct {
	Backend   physical.Backend
	HABackend physical.HABackend
	Cleanup   func()
}

type HandlerHandler interface {
	Handler(*HandlerProperties) http.Handler
}

type TestClusterOptions struct {
	KeepStandbysSealed       bool
	SkipInit                 bool
	HandlerFunc              HandlerHandler
	DefaultHandlerProperties HandlerProperties

	// BaseListenAddress is used to explicitly assign ports in sequence to the
	// listener of each core.  It should be a string of the form
	// "127.0.0.1:20000"
	//
	// WARNING: Using an explicitly assigned port above 30000 may clash with
	// ephemeral ports that have been assigned by the OS in other tests.  The
	// use of explicitly assigned ports below 30000 is strongly recommended.
	// In addition, you should be careful to use explicitly assigned ports that
	// do not clash with any other explicitly assigned ports in other tests.
	BaseListenAddress string

	// BaseClusterListenPort is used to explicitly assign ports in sequence to
	// the cluster listener of each core.  If BaseClusterListenPort is
	// specified, then BaseListenAddress must also be specified.  Each cluster
	// listener will use the same host as the one specified in
	// BaseListenAddress.
	//
	// WARNING: Using an explicitly assigned port above 30000 may clash with
	// ephemeral ports that have been assigned by the OS in other tests.  The
	// use of explicitly assigned ports below 30000 is strongly recommended.
	// In addition, you should be careful to use explicitly assigned ports that
	// do not clash with any other explicitly assigned ports in other tests.
	BaseClusterListenPort int

	NumCores       int
	SealFunc       func() Seal
	UnwrapSealFunc func() Seal
	Logger         log.Logger
	TempDir        string
	CACert         []byte
	CAKey          *ecdsa.PrivateKey
	// PhysicalFactory is used to create backends.
	// The int argument is the index of the core within the cluster, i.e. first
	// core in cluster will have 0, second 1, etc.
	// If the backend is shared across the cluster (i.e. is not Raft) then it
	// should return nil when coreIdx != 0.
	PhysicalFactory func(t testing.T, coreIdx int, logger log.Logger, conf map[string]interface{}) *PhysicalBackendBundle
	// FirstCoreNumber is used to assign a unique number to each core within
	// a multi-cluster setup.
	FirstCoreNumber   int
	RequireClientAuth bool
	// SetupFunc is called after the cluster is started.
	SetupFunc      func(t testing.T, c *TestCluster)
	PR1103Disabled bool

	// ClusterLayers are used to override the default cluster connection layer
	ClusterLayers cluster.NetworkLayerSet
	// InmemClusterLayers is a shorthand way of asking for ClusterLayers to be
	// built using the inmem implementation.
	InmemClusterLayers bool

	// RaftAddressProvider is used to set the raft ServerAddressProvider on
	// each core.
	//
	// If SkipInit is true, then RaftAddressProvider has no effect.
	// RaftAddressProvider should only be specified if the underlying physical
	// storage is Raft.
	RaftAddressProvider raftlib.ServerAddressProvider

	CoreMetricSinkProvider func(clusterName string) (*metricsutil.ClusterMetricSink, *metricsutil.MetricsHelper)

	PhysicalFactoryConfig map[string]interface{}
	LicensePublicKey      ed25519.PublicKey
	LicensePrivateKey     ed25519.PrivateKey

	// this stores the vault version that should be used for each core config
	VersionMap             map[int]string
	RedundancyZoneMap      map[int]string
	KVVersion              string
	EffectiveSDKVersionMap map[int]string

	NoDefaultQuotas bool

	Plugins *TestPluginConfig

	// if populated, the callback is called for every request
	RequestResponseCallback func(logical.Backend, *logical.Request, *logical.Response)

	// ABCDLoggerNames names the loggers according to our ABCD convention when generating 4 clusters
	ABCDLoggerNames bool
}

type TestPluginConfig struct {
	Typ      consts.PluginType
	Versions []string
}

var DefaultNumCores = 3

type certInfo struct {
	cert      *x509.Certificate
	certPEM   []byte
	certBytes []byte
	key       *ecdsa.PrivateKey
	keyPEM    []byte
}

// NewTestCluster creates a new test cluster based on the provided core config
// and test cluster options.
//
// N.B. Even though a single base CoreConfig is provided, NewTestCluster will instantiate a
// core config for each core it creates. If separate seal per core is desired, opts.SealFunc
// can be provided to generate a seal for each one. Otherwise, the provided base.Seal will be
// shared among cores. NewCore's default behavior is to generate a new DefaultSeal if the
// provided Seal in coreConfig (i.e. base.Seal) is nil.
//
// If opts.Logger is provided, it takes precedence and will be used as the cluster
// logger and will be the basis for each core's logger.  If no opts.Logger is
// given, one will be generated based on t.Name() for the cluster logger, and if
// no base.Logger is given will also be used as the basis for each core's logger.
func NewTestCluster(t testing.T, base *CoreConfig, opts *TestClusterOptions) *TestCluster {
	var err error

	var numCores int
	if opts == nil || opts.NumCores == 0 {
		numCores = DefaultNumCores
	} else {
		numCores = opts.NumCores
	}

	certIPs := []net.IP{
		net.IPv6loopback,
		net.ParseIP("127.0.0.1"),
	}

	baseAddr, certIPs := GenerateListenerAddr(t, opts, certIPs)
	var testCluster TestCluster
	testCluster.base = base

	switch {
	case opts != nil && opts.Logger != nil:
		testCluster.Logger = opts.Logger
	default:
		testCluster.Logger = corehelpers.NewTestLogger(t)
	}

	if opts != nil && opts.TempDir != "" {
		if _, err := os.Stat(opts.TempDir); os.IsNotExist(err) {
			if err := os.MkdirAll(opts.TempDir, 0o700); err != nil {
				t.Fatal(err)
			}
		}
		testCluster.TempDir = opts.TempDir
	} else {
		tempDir, err := ioutil.TempDir("", "vault-test-cluster-")
		if err != nil {
			t.Fatal(err)
		}
		testCluster.TempDir = tempDir
	}

	var caKey *ecdsa.PrivateKey
	if opts != nil && opts.CAKey != nil {
		caKey = opts.CAKey
	} else {
		caKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
	}
	testCluster.CAKey = caKey
	var caBytes []byte
	if opts != nil && len(opts.CACert) > 0 {
		caBytes = opts.CACert
	} else {
		caCertTemplate := &x509.Certificate{
			Subject: pkix.Name{
				CommonName: "localhost",
			},
			DNSNames:              []string{"localhost"},
			IPAddresses:           certIPs,
			KeyUsage:              x509.KeyUsage(x509.KeyUsageCertSign | x509.KeyUsageCRLSign),
			SerialNumber:          big.NewInt(mathrand.Int63()),
			NotBefore:             time.Now().Add(-30 * time.Second),
			NotAfter:              time.Now().Add(262980 * time.Hour),
			BasicConstraintsValid: true,
			IsCA:                  true,
		}
		caBytes, err = x509.CreateCertificate(rand.Reader, caCertTemplate, caCertTemplate, caKey.Public(), caKey)
		if err != nil {
			t.Fatal(err)
		}
	}
	caCert, err := x509.ParseCertificate(caBytes)
	if err != nil {
		t.Fatal(err)
	}
	testCluster.CACert = caCert
	testCluster.CACertBytes = caBytes
	testCluster.RootCAs = x509.NewCertPool()
	testCluster.RootCAs.AddCert(caCert)
	caCertPEMBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	}
	testCluster.CACertPEM = pem.EncodeToMemory(caCertPEMBlock)
	testCluster.CACertPEMFile = filepath.Join(testCluster.TempDir, "ca_cert.pem")
	err = ioutil.WriteFile(testCluster.CACertPEMFile, testCluster.CACertPEM, 0o755)
	if err != nil {
		t.Fatal(err)
	}
	marshaledCAKey, err := x509.MarshalECPrivateKey(caKey)
	if err != nil {
		t.Fatal(err)
	}
	caKeyPEMBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: marshaledCAKey,
	}
	testCluster.CAKeyPEM = pem.EncodeToMemory(caKeyPEMBlock)
	err = ioutil.WriteFile(filepath.Join(testCluster.TempDir, "ca_key.pem"), testCluster.CAKeyPEM, 0o755)
	if err != nil {
		t.Fatal(err)
	}

	var certInfoSlice []*certInfo

	//
	// Certs generation
	//
	for i := 0; i < numCores; i++ {
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		certTemplate := &x509.Certificate{
			Subject: pkix.Name{
				CommonName: "localhost",
			},
			// Include host.docker.internal for the sake of benchmark-vault running on MacOS/Windows.
			// This allows Prometheus running in docker to scrape the cluster for metrics.
			DNSNames:    []string{"localhost", "host.docker.internal"},
			IPAddresses: certIPs,
			ExtKeyUsage: []x509.ExtKeyUsage{
				x509.ExtKeyUsageServerAuth,
				x509.ExtKeyUsageClientAuth,
			},
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement,
			SerialNumber: big.NewInt(mathrand.Int63()),
			NotBefore:    time.Now().Add(-30 * time.Second),
			NotAfter:     time.Now().Add(262980 * time.Hour),
		}
		certBytes, err := x509.CreateCertificate(rand.Reader, certTemplate, caCert, key.Public(), caKey)
		if err != nil {
			t.Fatal(err)
		}
		cert, err := x509.ParseCertificate(certBytes)
		if err != nil {
			t.Fatal(err)
		}
		certPEMBlock := &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: certBytes,
		}
		certPEM := pem.EncodeToMemory(certPEMBlock)
		marshaledKey, err := x509.MarshalECPrivateKey(key)
		if err != nil {
			t.Fatal(err)
		}
		keyPEMBlock := &pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: marshaledKey,
		}
		keyPEM := pem.EncodeToMemory(keyPEMBlock)

		certInfoSlice = append(certInfoSlice, &certInfo{
			cert:      cert,
			certPEM:   certPEM,
			certBytes: certBytes,
			key:       key,
			keyPEM:    keyPEM,
		})
	}

	//
	// Listener setup
	//
	addresses := []*net.TCPAddr{}
	listeners := [][]*TestListener{}
	servers := []*http.Server{}
	handlers := []http.Handler{}
	tlsConfigs := []*tls.Config{}
	certGetters := []*reloadutil.CertificateGetter{}
	for i := 0; i < numCores; i++ {
		addr := &net.TCPAddr{
			IP:   baseAddr.IP,
			Port: 0,
		}
		if baseAddr.Port != 0 {
			addr.Port = baseAddr.Port + i
		}

		ln, err := net.ListenTCP("tcp", addr)
		if err != nil {
			t.Fatal(err)
		}
		addresses = append(addresses, addr)

		certFile := filepath.Join(testCluster.TempDir, fmt.Sprintf("node%d_port_%d_cert.pem", i+1, ln.Addr().(*net.TCPAddr).Port))
		keyFile := filepath.Join(testCluster.TempDir, fmt.Sprintf("node%d_port_%d_key.pem", i+1, ln.Addr().(*net.TCPAddr).Port))
		err = ioutil.WriteFile(certFile, certInfoSlice[i].certPEM, 0o755)
		if err != nil {
			t.Fatal(err)
		}
		err = ioutil.WriteFile(keyFile, certInfoSlice[i].keyPEM, 0o755)
		if err != nil {
			t.Fatal(err)
		}
		tlsCert, err := tls.X509KeyPair(certInfoSlice[i].certPEM, certInfoSlice[i].keyPEM)
		if err != nil {
			t.Fatal(err)
		}
		certGetter := reloadutil.NewCertificateGetter(certFile, keyFile, "")
		certGetters = append(certGetters, certGetter)
		certGetter.Reload()
		tlsConfig := &tls.Config{
			Certificates:   []tls.Certificate{tlsCert},
			RootCAs:        testCluster.RootCAs,
			ClientCAs:      testCluster.RootCAs,
			ClientAuth:     tls.RequestClientCert,
			NextProtos:     []string{"h2", "http/1.1"},
			GetCertificate: certGetter.GetCertificate,
		}
		if opts != nil && opts.RequireClientAuth {
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
			testCluster.ClientAuthRequired = true
		}
		tlsConfigs = append(tlsConfigs, tlsConfig)
		lns := []*TestListener{
			{
				Listener: tls.NewListener(ln, tlsConfig),
				Address:  ln.Addr().(*net.TCPAddr),
			},
		}
		listeners = append(listeners, lns)
		var handler http.Handler = http.NewServeMux()
		handlers = append(handlers, handler)
		server := &http.Server{
			Handler:  handler,
			ErrorLog: testCluster.Logger.StandardLogger(nil),
		}
		servers = append(servers, server)
	}

	// Create three cores with the same physical and different redirect/cluster
	// addrs.
	// N.B.: On OSX, instead of random ports, it assigns new ports to new
	// listeners sequentially. Aside from being a bad idea in a security sense,
	// it also broke tests that assumed it was OK to just use the port above
	// the redirect addr. This has now been changed to 105 ports above, but if
	// we ever do more than three nodes in a cluster it may need to be bumped.
	// Note: it's 105 so that we don't conflict with a running Consul by
	// default.
	coreConfig := &CoreConfig{
		LogicalBackends:    make(map[string]logical.Factory),
		CredentialBackends: make(map[string]logical.Factory),
		AuditBackends:      make(map[string]audit.Factory),
		RedirectAddr:       fmt.Sprintf("https://127.0.0.1:%d", listeners[0][0].Address.Port),
		ClusterAddr:        "https://127.0.0.1:0",
		DisableMlock:       true,
		EnableUI:           true,
		EnableRaw:          true,
		BuiltinRegistry:    corehelpers.NewMockBuiltinRegistry(),
	}

	corehelpers.RegisterSubloggerAdder(testCluster.Logger, coreConfig)

	if base != nil {
		coreConfig.DetectDeadlocks = TestDeadlockDetection
		coreConfig.RawConfig = base.RawConfig
		coreConfig.DisableCache = base.DisableCache
		coreConfig.EnableUI = base.EnableUI
		coreConfig.DefaultLeaseTTL = base.DefaultLeaseTTL
		coreConfig.MaxLeaseTTL = base.MaxLeaseTTL
		coreConfig.CacheSize = base.CacheSize
		coreConfig.PluginDirectory = base.PluginDirectory
		coreConfig.Seal = base.Seal
		coreConfig.UnwrapSeal = base.UnwrapSeal
		coreConfig.DevToken = base.DevToken
		coreConfig.EnableRaw = base.EnableRaw
		coreConfig.DisableSealWrap = base.DisableSealWrap
		coreConfig.DisableCache = base.DisableCache
		coreConfig.LicensingConfig = base.LicensingConfig
		coreConfig.License = base.License
		coreConfig.LicensePath = base.LicensePath
		coreConfig.DisablePerformanceStandby = base.DisablePerformanceStandby
		coreConfig.MetricsHelper = base.MetricsHelper
		coreConfig.MetricSink = base.MetricSink
		coreConfig.SecureRandomReader = base.SecureRandomReader
		coreConfig.DisableSentinelTrace = base.DisableSentinelTrace
		coreConfig.ClusterName = base.ClusterName
		coreConfig.DisableAutopilot = base.DisableAutopilot
		coreConfig.AdministrativeNamespacePath = base.AdministrativeNamespacePath
		coreConfig.ServiceRegistration = base.ServiceRegistration

		if base.BuiltinRegistry != nil {
			coreConfig.BuiltinRegistry = base.BuiltinRegistry
		}

		if !coreConfig.DisableMlock {
			base.DisableMlock = false
		}

		if base.Physical != nil {
			coreConfig.Physical = base.Physical
		}

		if base.HAPhysical != nil {
			coreConfig.HAPhysical = base.HAPhysical
		}

		// Used to set something non-working to test fallback
		switch base.ClusterAddr {
		case "empty":
			coreConfig.ClusterAddr = ""
		case "":
		default:
			coreConfig.ClusterAddr = base.ClusterAddr
		}

		if base.LogicalBackends != nil {
			for k, v := range base.LogicalBackends {
				coreConfig.LogicalBackends[k] = v
			}
		}
		if base.CredentialBackends != nil {
			for k, v := range base.CredentialBackends {
				coreConfig.CredentialBackends[k] = v
			}
		}
		if base.AuditBackends != nil {
			for k, v := range base.AuditBackends {
				coreConfig.AuditBackends[k] = v
			}
		}
		if base.Logger != nil {
			coreConfig.Logger = base.Logger
		}

		coreConfig.ClusterCipherSuites = base.ClusterCipherSuites
		coreConfig.DisableCache = base.DisableCache
		coreConfig.DevToken = base.DevToken
		coreConfig.RecoveryMode = base.RecoveryMode
		coreConfig.ActivityLogConfig = base.ActivityLogConfig
		coreConfig.EnableResponseHeaderHostname = base.EnableResponseHeaderHostname
		coreConfig.EnableResponseHeaderRaftNodeID = base.EnableResponseHeaderRaftNodeID
		coreConfig.RollbackPeriod = base.RollbackPeriod
		coreConfig.PendingRemovalMountsAllowed = base.PendingRemovalMountsAllowed
		coreConfig.ExpirationRevokeRetryBase = base.ExpirationRevokeRetryBase
		testApplyEntBaseConfig(coreConfig, base)
	}
	if coreConfig.ClusterName == "" {
		coreConfig.ClusterName = t.Name()
	}

	if coreConfig.ClusterName == "" {
		coreConfig.ClusterName = t.Name()
	}

	if coreConfig.ClusterHeartbeatInterval == 0 {
		// Set this lower so that state populates quickly to standby nodes
		coreConfig.ClusterHeartbeatInterval = 2 * time.Second
	}

	if coreConfig.RawConfig == nil {
		c := new(server.Config)
		c.SharedConfig = &configutil.SharedConfig{LogFormat: logging.UnspecifiedFormat.String()}
		coreConfig.RawConfig = c
	}

	addAuditBackend := len(coreConfig.AuditBackends) == 0
	if addAuditBackend {
		coreConfig.AuditBackends["noop"] = corehelpers.NoopAuditFactory(nil)
	}

	if coreConfig.Physical == nil && (opts == nil || opts.PhysicalFactory == nil) {
		coreConfig.Physical, err = physInmem.NewInmem(nil, testCluster.Logger)
		if err != nil {
			t.Fatal(err)
		}
	}
	if coreConfig.HAPhysical == nil && (opts == nil || opts.PhysicalFactory == nil) {
		haPhys, err := physInmem.NewInmemHA(nil, testCluster.Logger)
		if err != nil {
			t.Fatal(err)
		}
		coreConfig.HAPhysical = haPhys.(physical.HABackend)
	}

	if testCluster.LicensePublicKey == nil {
		pubKey, priKey, err := GenerateTestLicenseKeys()
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		testCluster.LicensePublicKey = pubKey
		testCluster.LicensePrivateKey = priKey
	}

	if opts != nil && opts.InmemClusterLayers {
		if opts.ClusterLayers != nil {
			t.Fatalf("cannot specify ClusterLayers when InmemClusterLayers is true")
		}
		inmemCluster, err := cluster.NewInmemLayerCluster("inmem-cluster", numCores, testCluster.Logger.Named("inmem-cluster"))
		if err != nil {
			t.Fatal(err)
		}
		opts.ClusterLayers = inmemCluster
	}

	if opts != nil && opts.Plugins != nil {
		var pluginDir string
		var cleanup func(t testing.T)

		if coreConfig.PluginDirectory == "" {
			pluginDir, cleanup = corehelpers.MakeTestPluginDir(t)
			coreConfig.PluginDirectory = pluginDir
			t.Cleanup(func() { cleanup(t) })
		}

		var plugins []pluginhelpers.TestPlugin
		for _, version := range opts.Plugins.Versions {
			plugins = append(plugins, pluginhelpers.CompilePlugin(t, opts.Plugins.Typ, version, coreConfig.PluginDirectory))
		}
		testCluster.Plugins = plugins
	}

	// Create cores
	testCluster.cleanupFuncs = []func(){}
	cores := []*Core{}
	coreConfigs := []*CoreConfig{}

	for i := 0; i < numCores; i++ {
		cleanup, c, localConfig, handler := testCluster.newCore(t, i, coreConfig, opts, listeners[i], testCluster.LicensePublicKey)

		corehelpers.RegisterSubloggerAdder(testCluster.Logger, c)

		testCluster.cleanupFuncs = append(testCluster.cleanupFuncs, cleanup)
		cores = append(cores, c)
		coreConfigs = append(coreConfigs, &localConfig)

		if handler != nil {
			handlers[i] = handler
			servers[i].Handler = handlers[i]
		}
	}

	// Clustering setup
	for i := 0; i < numCores; i++ {
		testCluster.setupClusterListener(t, i, cores[i], coreConfigs[i], opts, listeners[i], handlers[i])
	}

	// Create TestClusterCores
	var ret []*TestClusterCore
	for i := 0; i < numCores; i++ {
		tcc := &TestClusterCore{
			Core:                 cores[i],
			CoreConfig:           coreConfigs[i],
			ServerKey:            certInfoSlice[i].key,
			ServerKeyPEM:         certInfoSlice[i].keyPEM,
			ServerCert:           certInfoSlice[i].cert,
			ServerCertBytes:      certInfoSlice[i].certBytes,
			ServerCertPEM:        certInfoSlice[i].certPEM,
			Address:              addresses[i],
			Listeners:            listeners[i],
			Handler:              handlers[i],
			Server:               servers[i],
			tlsConfig:            tlsConfigs[i],
			Barrier:              cores[i].barrier,
			NodeID:               fmt.Sprintf("core-%d", i),
			UnderlyingRawStorage: coreConfigs[i].Physical,
			UnderlyingHAStorage:  coreConfigs[i].HAPhysical,
		}
		tcc.ReloadFuncs = &cores[i].reloadFuncs
		tcc.ReloadFuncsLock = &cores[i].reloadFuncsLock
		tcc.ReloadFuncsLock.Lock()
		(*tcc.ReloadFuncs)["listener|tcp"] = []reloadutil.ReloadFunc{certGetters[i].Reload}
		tcc.ReloadFuncsLock.Unlock()

		testAdjustUnderlyingStorage(tcc)

		ret = append(ret, tcc)
	}
	testCluster.Cores = ret

	// Initialize cores
	if opts == nil || !opts.SkipInit {
		testCluster.initCores(t, opts, addAuditBackend)
	}

	// Assign clients
	for i := 0; i < numCores; i++ {
		testCluster.Cores[i].Client = testCluster.getAPIClient(t, opts, listeners[i][0].Address.Port, tlsConfigs[i])
	}

	// Extra Setup
	for _, tcc := range testCluster.Cores {
		testExtraTestCoreSetup(t, testCluster.LicensePrivateKey, tcc)
	}

	// Cleanup
	testCluster.CleanupFunc = func() {
		for _, c := range testCluster.cleanupFuncs {
			c()
		}
		if l, ok := testCluster.Logger.(*corehelpers.TestLogger); ok {
			if t.Failed() {
				_ = l.File.Close()
			} else {
				_ = os.Remove(l.Path)
			}
		}
	}

	// Setup
	if opts != nil {
		if opts.SetupFunc != nil {
			testCluster.SetupFunc = func() {
				opts.SetupFunc(t, &testCluster)
			}
		}
	}

	testCluster.opts = opts
	testCluster.start(t)
	return &testCluster
}

// StopCore performs an orderly shutdown of a core.
func (cluster *TestCluster) StopCore(t testing.T, idx int) {
	t.Helper()

	if idx < 0 || idx > len(cluster.Cores) {
		t.Fatalf("invalid core index %d", idx)
	}
	tcc := cluster.Cores[idx]
	tcc.Logger().Info("stopping core", "core", idx)

	// Stop listeners and call Finalize()
	if err := tcc.stop(); err != nil {
		t.Fatal(err)
	}

	// Run cleanup
	cluster.cleanupFuncs[idx]()
}

func GenerateListenerAddr(t testing.T, opts *TestClusterOptions, certIPs []net.IP) (*net.TCPAddr, []net.IP) {
	var baseAddr *net.TCPAddr
	var err error

	if opts != nil && opts.BaseListenAddress != "" {
		baseAddr, err = net.ResolveTCPAddr("tcp", opts.BaseListenAddress)

		if err != nil {
			t.Fatal("could not parse given base IP")
		}
		certIPs = append(certIPs, baseAddr.IP)
	} else {
		baseAddr = &net.TCPAddr{
			IP:   net.ParseIP("127.0.0.1"),
			Port: 0,
		}
	}

	return baseAddr, certIPs
}

// StartCore restarts a TestClusterCore that was stopped, by replacing the
// underlying Core.
func (cluster *TestCluster) StartCore(t testing.T, idx int, opts *TestClusterOptions) {
	t.Helper()

	if idx < 0 || idx > len(cluster.Cores) {
		t.Fatalf("invalid core index %d", idx)
	}
	tcc := cluster.Cores[idx]
	tcc.Logger().Info("restarting core", "core", idx)

	// Set up listeners
	ln, err := net.ListenTCP("tcp", tcc.Address)
	if err != nil {
		t.Fatal(err)
	}
	tcc.Listeners = []*TestListener{
		{
			Listener: tls.NewListener(ln, tcc.tlsConfig),
			Address:  ln.Addr().(*net.TCPAddr),
		},
	}

	tcc.Handler = http.NewServeMux()
	tcc.Server = &http.Server{
		Handler:  tcc.Handler,
		ErrorLog: cluster.Logger.StandardLogger(nil),
	}

	// Create a new Core
	cleanup, newCore, localConfig, coreHandler := cluster.newCore(t, idx, tcc.CoreConfig, opts, tcc.Listeners, cluster.LicensePublicKey)
	if coreHandler != nil {
		tcc.Handler = coreHandler
		tcc.Server.Handler = coreHandler
	}

	cluster.cleanupFuncs[idx] = cleanup
	tcc.Core = newCore
	tcc.CoreConfig = &localConfig
	tcc.UnderlyingRawStorage = localConfig.Physical

	cluster.setupClusterListener(
		t, idx, newCore, tcc.CoreConfig,
		opts, tcc.Listeners, tcc.Handler)

	tcc.Client = cluster.getAPIClient(t, opts, tcc.Listeners[0].Address.Port, tcc.tlsConfig)

	testAdjustUnderlyingStorage(tcc)
	testExtraTestCoreSetup(t, cluster.LicensePrivateKey, tcc)

	// Start listeners
	for _, ln := range tcc.Listeners {
		tcc.Logger().Info("starting listener for core", "port", ln.Address.Port)
		go tcc.Server.Serve(ln)
	}

	tcc.Logger().Info("restarted test core", "core", idx)
}

func (testCluster *TestCluster) newCore(t testing.T, idx int, coreConfig *CoreConfig, opts *TestClusterOptions, listeners []*TestListener, pubKey ed25519.PublicKey) (func(), *Core, CoreConfig, http.Handler) {
	localConfig := *coreConfig
	cleanupFunc := func() {}
	var handler http.Handler

	var disablePR1103 bool
	if opts != nil && opts.PR1103Disabled {
		disablePR1103 = true
	}

	var firstCoreNumber int
	if opts != nil {
		firstCoreNumber = opts.FirstCoreNumber
	}

	localConfig.RedirectAddr = fmt.Sprintf("https://127.0.0.1:%d", listeners[0].Address.Port)

	// if opts.SealFunc is provided, use that to generate a seal for the config instead
	if opts != nil && opts.SealFunc != nil {
		localConfig.Seal = opts.SealFunc()
	}
	if opts != nil && opts.UnwrapSealFunc != nil {
		localConfig.UnwrapSeal = opts.UnwrapSealFunc()
	}

	if coreConfig.Logger == nil || (opts != nil && opts.Logger != nil) {
		localConfig.Logger = testCluster.Logger.Named(fmt.Sprintf("core%d", idx))
	}

	if opts != nil && opts.EffectiveSDKVersionMap != nil {
		localConfig.EffectiveSDKVersion = opts.EffectiveSDKVersionMap[idx]
	}

	if opts != nil && opts.PhysicalFactory != nil {
		pfc := opts.PhysicalFactoryConfig
		if pfc == nil {
			pfc = make(map[string]interface{})
		}
		if len(opts.VersionMap) > 0 {
			pfc["autopilot_upgrade_version"] = opts.VersionMap[idx]
		}
		if len(opts.RedundancyZoneMap) > 0 {
			pfc["autopilot_redundancy_zone"] = opts.RedundancyZoneMap[idx]
		}
		physBundle := opts.PhysicalFactory(t, idx, localConfig.Logger, pfc)
		switch {
		case physBundle == nil && coreConfig.Physical != nil:
		case physBundle == nil && coreConfig.Physical == nil:
			t.Fatal("PhysicalFactory produced no physical and none in CoreConfig")
		case physBundle != nil:
			// Storage backend setup
			if physBundle.Backend != nil {
				testCluster.Logger.Info("created physical backend", "instance", idx)
				coreConfig.Physical = physBundle.Backend
				localConfig.Physical = physBundle.Backend
			}

			// HA Backend setup
			haBackend := physBundle.HABackend
			if haBackend == nil {
				if ha, ok := physBundle.Backend.(physical.HABackend); ok {
					haBackend = ha
				}
			}
			coreConfig.HAPhysical = haBackend
			localConfig.HAPhysical = haBackend

			// Cleanup setup
			if physBundle.Cleanup != nil {
				cleanupFunc = physBundle.Cleanup
			}
		}
	}

	if opts != nil && opts.ClusterLayers != nil {
		localConfig.ClusterNetworkLayer = opts.ClusterLayers.Layers()[idx]
		localConfig.ClusterAddr = "https://" + localConfig.ClusterNetworkLayer.Listeners()[0].Addr().String()
	}

	switch {
	case localConfig.LicensingConfig != nil:
		if pubKey != nil {
			localConfig.LicensingConfig.AdditionalPublicKeys = append(localConfig.LicensingConfig.AdditionalPublicKeys, pubKey)
		}
	default:
		localConfig.LicensingConfig = testGetLicensingConfig(pubKey)
	}

	if localConfig.MetricsHelper == nil {
		inm := metrics.NewInmemSink(10*time.Second, time.Minute)
		metrics.DefaultInmemSignal(inm)
		localConfig.MetricsHelper = metricsutil.NewMetricsHelper(inm, false)
	}
	if opts != nil && opts.CoreMetricSinkProvider != nil {
		localConfig.MetricSink, localConfig.MetricsHelper = opts.CoreMetricSinkProvider(localConfig.ClusterName)
	}

	if opts != nil && opts.CoreMetricSinkProvider != nil {
		localConfig.MetricSink, localConfig.MetricsHelper = opts.CoreMetricSinkProvider(localConfig.ClusterName)
	}

	localConfig.NumExpirationWorkers = numExpirationWorkersTest

	c, err := NewCore(&localConfig)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	c.coreNumber = firstCoreNumber + idx
	c.PR1103disabled = disablePR1103
	if opts != nil && opts.HandlerFunc != nil {
		props := opts.DefaultHandlerProperties
		props.Core = c
		if props.ListenerConfig != nil && props.ListenerConfig.MaxRequestDuration == 0 {
			props.ListenerConfig.MaxRequestDuration = DefaultMaxRequestDuration
		}
		handler = opts.HandlerFunc.Handler(&props)
	}

	if opts != nil && opts.RequestResponseCallback != nil {
		c.requestResponseCallback = opts.RequestResponseCallback
	}

	// Set this in case the Seal was manually set before the core was
	// created
	if localConfig.Seal != nil {
		localConfig.Seal.SetCore(c)
	}

	return cleanupFunc, c, localConfig, handler
}

func (testCluster *TestCluster) setupClusterListener(
	t testing.T, idx int, core *Core, coreConfig *CoreConfig,
	opts *TestClusterOptions, listeners []*TestListener, handler http.Handler,
) {
	if coreConfig.ClusterAddr == "" {
		return
	}

	clusterAddrGen := func(lns []*TestListener, port int) []*net.TCPAddr {
		ret := make([]*net.TCPAddr, len(lns))
		for i, ln := range lns {
			ret[i] = &net.TCPAddr{
				IP:   ln.Address.IP,
				Port: port,
			}
		}
		return ret
	}

	baseClusterListenPort := 0
	if opts != nil && opts.BaseClusterListenPort != 0 {
		if opts.BaseListenAddress == "" {
			t.Fatal("BaseListenAddress is not specified")
		}
		baseClusterListenPort = opts.BaseClusterListenPort
	}

	port := 0
	if baseClusterListenPort != 0 {
		port = baseClusterListenPort + idx
	}
	core.Logger().Info("assigning cluster listener for test core", "core", idx, "port", port)
	core.SetClusterListenerAddrs(clusterAddrGen(listeners, port))
	core.SetClusterHandler(handler)
}

func (tc *TestCluster) InitCores(t testing.T, opts *TestClusterOptions, addAuditBackend bool) {
	tc.initCores(t, opts, addAuditBackend)
}

func (tc *TestCluster) initCores(t testing.T, opts *TestClusterOptions, addAuditBackend bool) {
	leader := tc.Cores[0]

	bKeys, rKeys, root := TestCoreInitClusterWrapperSetup(t, leader.Core, leader.Handler)
	barrierKeys, _ := copystructure.Copy(bKeys)
	tc.BarrierKeys = barrierKeys.([][]byte)
	recoveryKeys, _ := copystructure.Copy(rKeys)
	tc.RecoveryKeys = recoveryKeys.([][]byte)
	tc.RootToken = root

	// Write root token and barrier keys
	err := ioutil.WriteFile(filepath.Join(tc.TempDir, "root_token"), []byte(root), 0o755)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	for i, key := range tc.BarrierKeys {
		buf.WriteString(base64.StdEncoding.EncodeToString(key))
		if i < len(tc.BarrierKeys)-1 {
			buf.WriteRune('\n')
		}
	}
	err = ioutil.WriteFile(filepath.Join(tc.TempDir, "barrier_keys"), buf.Bytes(), 0o755)
	if err != nil {
		t.Fatal(err)
	}
	for i, key := range tc.RecoveryKeys {
		buf.WriteString(base64.StdEncoding.EncodeToString(key))
		if i < len(tc.RecoveryKeys)-1 {
			buf.WriteRune('\n')
		}
	}
	err = ioutil.WriteFile(filepath.Join(tc.TempDir, "recovery_keys"), buf.Bytes(), 0o755)
	if err != nil {
		t.Fatal(err)
	}

	// Unseal first core
	for _, key := range bKeys {
		if _, err := leader.Core.Unseal(TestKeyCopy(key)); err != nil {
			t.Fatalf("unseal err: %s", err)
		}
	}

	ctx := context.Background()

	// If stored keys is supported, the above will no no-op, so trigger auto-unseal
	// using stored keys to try to unseal
	if err := leader.Core.UnsealWithStoredKeys(ctx); err != nil {
		t.Fatal(err)
	}

	// Verify unsealed
	if leader.Core.Sealed() {
		t.Fatal("should not be sealed")
	}

	TestWaitActive(t, leader.Core)

	kvVersion := "1"
	if opts != nil {
		kvVersion = opts.KVVersion
	}

	// Existing tests rely on this; we can make a toggle to disable it
	// later if we want
	kvReq := &logical.Request{
		Operation:   logical.UpdateOperation,
		ClientToken: tc.RootToken,
		Path:        "sys/mounts/secret",
		Data: map[string]interface{}{
			"type":        "kv",
			"path":        "secret/",
			"description": "key/value secret storage",
			"options": map[string]string{
				"version": kvVersion,
			},
		},
	}
	resp, err := leader.Core.HandleRequest(namespace.RootContext(ctx), kvReq)
	if err != nil {
		t.Fatal(err)
	}
	if resp.IsError() {
		t.Fatal(err)
	}

	cfg, err := leader.Core.seal.BarrierConfig(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Unseal other cores unless otherwise specified
	numCores := len(tc.Cores)
	if (opts == nil || !opts.KeepStandbysSealed) && numCores > 1 {
		for i := 1; i < numCores; i++ {
			tc.Cores[i].Core.seal.SetCachedBarrierConfig(cfg)
			for _, key := range bKeys {
				if _, err := tc.Cores[i].Core.Unseal(TestKeyCopy(key)); err != nil {
					t.Fatalf("unseal err: %s", err)
				}
			}

			// If stored keys is supported, the above will no no-op, so trigger auto-unseal
			// using stored keys
			if err := tc.Cores[i].Core.UnsealWithStoredKeys(ctx); err != nil {
				t.Fatal(err)
			}
		}

		// Let them come fully up to standby
		time.Sleep(2 * time.Second)

		// Ensure cluster connection info is populated.
		// Other cores should not come up as leaders.
		for i := 1; i < numCores; i++ {
			isLeader, _, _, err := tc.Cores[i].Core.Leader()
			if err != nil {
				t.Fatal(err)
			}
			if isLeader {
				t.Fatalf("core[%d] should not be leader", i)
			}
		}
	}

	//
	// Set test cluster core(s) and test cluster
	//
	cluster, err := leader.Core.Cluster(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	tc.ID = cluster.ID

	if addAuditBackend {
		// Enable auditing.
		auditReq := &logical.Request{
			Operation:   logical.UpdateOperation,
			ClientToken: tc.RootToken,
			Path:        "sys/audit/noop",
			Data: map[string]interface{}{
				"type": "noop",
			},
		}
		resp, err = leader.Core.HandleRequest(namespace.RootContext(ctx), auditReq)
		if err != nil {
			t.Fatal(err)
		}

		if resp.IsError() {
			t.Fatal(err)
		}
	}
}

func (testCluster *TestCluster) getAPIClient(
	t testing.T, opts *TestClusterOptions,
	port int, tlsConfig *tls.Config,
) *api.Client {
	transport := cleanhttp.DefaultPooledTransport()
	transport.TLSClientConfig = tlsConfig.Clone()
	if err := http2.ConfigureTransport(transport); err != nil {
		t.Fatal(err)
	}
	client := &http.Client{
		Transport: transport,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			// This can of course be overridden per-test by using its own client
			return fmt.Errorf("redirects not allowed in these tests")
		},
	}
	config := api.DefaultConfig()
	if config.Error != nil {
		t.Fatal(config.Error)
	}
	config.Address = fmt.Sprintf("https://127.0.0.1:%d", port)
	config.HttpClient = client
	config.MaxRetries = 0
	apiClient, err := api.NewClient(config)
	if err != nil {
		t.Fatal(err)
	}
	if opts == nil || !opts.SkipInit {
		apiClient.SetToken(testCluster.RootToken)
	}
	return apiClient
}

func (c *TestCluster) GetBarrierOrRecoveryKeys() [][]byte {
	if c.Cores[0].SealAccess().RecoveryKeySupported() {
		return c.GetRecoveryKeys()
	} else {
		return c.GetBarrierKeys()
	}
}

func (c *TestCluster) GetCACertPEMFile() string {
	return c.CACertPEMFile
}

func (c *TestCluster) ClusterID() string {
	return c.ID
}

func (c *TestCluster) Nodes() []testcluster.VaultClusterNode {
	ret := make([]testcluster.VaultClusterNode, len(c.Cores))
	for i, core := range c.Cores {
		ret[i] = core
	}
	return ret
}

func (c *TestCluster) SetBarrierKeys(keys [][]byte) {
	c.BarrierKeys = make([][]byte, len(keys))
	for i, k := range keys {
		c.BarrierKeys[i] = TestKeyCopy(k)
	}
}

func (c *TestCluster) SetRecoveryKeys(keys [][]byte) {
	c.RecoveryKeys = make([][]byte, len(keys))
	for i, k := range keys {
		c.RecoveryKeys[i] = TestKeyCopy(k)
	}
}

func (c *TestCluster) GetBarrierKeys() [][]byte {
	ret := make([][]byte, len(c.BarrierKeys))
	for i, k := range c.BarrierKeys {
		ret[i] = TestKeyCopy(k)
	}
	return ret
}

func (c *TestCluster) GetRecoveryKeys() [][]byte {
	ret := make([][]byte, len(c.RecoveryKeys))
	for i, k := range c.RecoveryKeys {
		ret[i] = TestKeyCopy(k)
	}
	return ret
}

func (c *TestCluster) NamedLogger(name string) log.Logger {
	return c.Logger.Named(name)
}

func (c *TestCluster) GetRootToken() string {
	return c.RootToken
}

func (c *TestClusterCore) Name() string {
	return c.NodeID
}

func (c *TestClusterCore) APIClient() *api.Client {
	return c.Client
}

var (
	_ testcluster.VaultCluster     = &TestCluster{}
	_ testcluster.VaultClusterNode = &TestClusterCore{}
)
