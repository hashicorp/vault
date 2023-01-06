//go:build !deadlock

package vault

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	mathrand "math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/mitchellh/go-testing-interface"

	"github.com/hashicorp/go-secure-stdlib/reloadutil"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"
	physInmem "github.com/hashicorp/vault/sdk/physical/inmem"
	"github.com/hashicorp/vault/vault/cluster"
)

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
	var baseAddr *net.TCPAddr
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

	var testCluster TestCluster
	testCluster.base = base

	switch {
	case opts != nil && opts.Logger != nil:
		testCluster.Logger = opts.Logger
	default:
		testCluster.Logger = NewTestLogger(t)
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
		BuiltinRegistry:    NewMockBuiltinRegistry(),
	}

	if base != nil {
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
		AddNoopAudit(coreConfig, nil)
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

	// Create cores
	testCluster.cleanupFuncs = []func(){}
	cores := []*Core{}
	coreConfigs := []*CoreConfig{}

	for i := 0; i < numCores; i++ {
		cleanup, c, localConfig, handler := testCluster.newCore(t, i, coreConfig, opts, listeners[i], testCluster.LicensePublicKey)

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
			TLSConfig:            tlsConfigs[i],
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
		if l, ok := testCluster.Logger.(*TestLogger); ok {
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
