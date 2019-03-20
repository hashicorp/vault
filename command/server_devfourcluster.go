package command

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/big"
	mathrand "math/rand"
	"net"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/testhelpers"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/version"
	testing "github.com/mitchellh/go-testing-interface"
	"github.com/pkg/errors"
)

func (c *ServerCommand) enableFourClusterDev(base *vault.CoreConfig, info map[string]string, infoKeys []string, devListenAddress, tempDir string) int {
	var err error
	ctx := namespace.RootContext(nil)
	clusters := map[string]*vault.TestCluster{}

	if base.DevToken == "" {
		base.DevToken = "root"
	}
	base.EnableRaw = true

	// Without setting something in the future we get bombarded with warnings which are quite annoying during testing
	base.DevLicenseDuration = 6 * time.Hour

	// Set a base temp dir
	if tempDir == "" {
		tempDir, err = ioutil.TempDir("", "vault-test-cluster-")
		if err != nil {
			c.UI.Error(fmt.Sprintf("failed to create top-level temp dir: %s", err))
			return 1
		}
	}

	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to generate CA key: %s", err))
		return 1
	}
	certIPs := []net.IP{
		net.IPv6loopback,
		net.ParseIP("127.0.0.1"),
	}
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
	caBytes, err := x509.CreateCertificate(rand.Reader, caCertTemplate, caCertTemplate, caKey.Public(), caKey)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to generate certificate: %s", err))
		return 1
	}

	getCluster := func(name string) error {
		factory := c.PhysicalBackends["inmem_transactional_ha"]
		backend, err := factory(nil, c.logger)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error initializing storage of type %s: %s", "inmem_transactional_ha", err))
			return errors.New("")
		}
		base.Physical = backend
		base.Seal = vault.NewDefaultSeal()

		testCluster := vault.NewTestCluster(&testing.RuntimeT{}, base, &vault.TestClusterOptions{
			HandlerFunc: vaulthttp.Handler,
			//BaseListenAddress: c.flagDevListenAddr,
			Logger:  c.logger.Named(name),
			TempDir: fmt.Sprintf("%s/%s", tempDir, name),
			CAKey:   caKey,
			CACert:  caBytes,
		})

		clusters[name] = testCluster

		for i, core := range testCluster.Cores {
			info[fmt.Sprintf("%s node %d redirect address", name, i)] = fmt.Sprintf("https://%s", core.Listeners[0].Address.String())
			infoKeys = append(infoKeys, fmt.Sprintf("%s node %d redirect address", name, i))
			core.Server.Handler = vaulthttp.Handler(&vault.HandlerProperties{
				Core: core.Core,
			})
			core.SetClusterHandler(core.Server.Handler)
		}

		testCluster.Start()

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
			return errors.New("")
		}
		if resp == nil {
			c.UI.Error(fmt.Sprintf("nil response when creating root token with ID %s", base.DevToken))
			return errors.New("")
		}
		if resp.Auth == nil {
			c.UI.Error(fmt.Sprintf("nil auth when creating root token with ID %s", base.DevToken))
			return errors.New("")
		}

		testCluster.RootToken = resp.Auth.ClientToken

		req.ID = "dev-revoke-init-root"
		req.Path = "auth/token/revoke-self"
		req.Data = nil
		resp, err = testCluster.Cores[0].HandleRequest(ctx, req)
		if err != nil {
			c.UI.Output(fmt.Sprintf("failed to revoke initial root token: %s", err))
			return errors.New("")
		}

		for _, core := range testCluster.Cores {
			core.Client.SetToken(base.DevToken)
		}

		return nil
	}

	err = getCluster("perf-pri")
	if err != nil {
		return 1
	}
	err = getCluster("perf-pri-dr")
	if err != nil {
		return 1
	}
	err = getCluster("perf-sec")
	if err != nil {
		return 1
	}
	err = getCluster("perf-sec-dr")
	if err != nil {
		return 1
	}

	clusterCleanup := func() {
		for name, cluster := range clusters {
			cluster.Cleanup()

			// Shutdown will wait until after Vault is sealed, which means the
			// request forwarding listeners will also be closed (and also
			// waited for).
			for _, core := range cluster.Cores {
				if err := core.Shutdown(); err != nil {
					c.UI.Error(fmt.Sprintf("Error with cluster %s core shutdown: %s", name, err))
				}
			}
		}
	}

	defer c.cleanupGuard.Do(clusterCleanup)

	info["cluster parameters path"] = tempDir
	infoKeys = append(infoKeys, "cluster parameters path")

	verInfo := version.GetVersion()
	info["version"] = verInfo.FullVersionNumber(false)
	infoKeys = append(infoKeys, "version")
	if verInfo.Revision != "" {
		info["version sha"] = strings.Trim(verInfo.Revision, "'")
		infoKeys = append(infoKeys, "version sha")
	}
	infoKeys = append(infoKeys, "cgo")
	info["cgo"] = "disabled"
	if version.CgoEnabled {
		info["cgo"] = "enabled"
	}

	// Server configuration output
	padding := 40
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

	// Set the token
	tokenHelper, err := c.TokenHelper()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error getting token helper: %s", err))
		return 1
	}
	if err := tokenHelper.Store(base.DevToken); err != nil {
		c.UI.Error(fmt.Sprintf("Error storing in token helper: %s", err))
		return 1
	}

	if err := ioutil.WriteFile(filepath.Join(tempDir, "root_token"), []byte(base.DevToken), 0755); err != nil {
		c.UI.Error(fmt.Sprintf("Error writing token to tempfile: %s", err))
		return 1
	}

	c.UI.Output(fmt.Sprintf(
		"\nRoot Token: %s\n", base.DevToken,
	))

	for i, key := range clusters["perf-pri"].BarrierKeys {
		c.UI.Output(fmt.Sprintf(
			"Unseal Key %d: %s",
			i+1, base64.StdEncoding.EncodeToString(key),
		))
	}

	c.UI.Output(fmt.Sprintf(
		"\nUseful env vars:\n"+
			"export VAULT_TOKEN=%s\n"+
			"export VAULT_CACERT=%s/perf-pri/ca_cert.pem\n",
		base.DevToken,
		tempDir,
	))
	c.UI.Output(fmt.Sprintf("Addresses of initial active nodes:"))
	clusterNames := []string{}
	for name := range clusters {
		clusterNames = append(clusterNames, name)
	}
	sort.Strings(clusterNames)
	for _, name := range clusterNames {
		c.UI.Output(fmt.Sprintf(
			"%s:\n"+
				"export VAULT_ADDR=%s\n",
			name,
			clusters[name].Cores[0].Client.Address(),
		))
	}

	// Output the header that the server has started
	c.UI.Output("==> Vault clusters started! Log data will stream in below:\n")

	// Inform any tests that the server is ready
	select {
	case c.startedCh <- struct{}{}:
	default:
	}

	// Release the log gate.
	c.logGate.Flush()

	testhelpers.SetupFourClusterReplication(&testing.RuntimeT{},
		clusters["perf-pri"],
		clusters["perf-sec"],
		clusters["perf-pri-dr"],
		clusters["perf-sec-dr"],
	)

	// Wait for shutdown
	shutdownTriggered := false

	for !shutdownTriggered {
		select {
		case <-c.ShutdownCh:
			c.UI.Output("==> Vault shutdown triggered")

			// Stop the listeners so that we don't process further client requests.
			c.cleanupGuard.Do(clusterCleanup)

			shutdownTriggered = true

		case <-c.SighupCh:
			c.UI.Output("==> Vault reload triggered")
			for name, cluster := range clusters {
				for _, core := range cluster.Cores {
					if err := c.Reload(core.ReloadFuncsLock, core.ReloadFuncs, nil); err != nil {
						c.UI.Error(fmt.Sprintf("Error(s) were encountered during reload of cluster %s cores: %s", name, err))
					}
				}
			}
		}
	}

	return 0
}
