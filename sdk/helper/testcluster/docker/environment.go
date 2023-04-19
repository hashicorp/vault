package docker

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	mathrand "math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/logging"

	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/api"
	dockhelper "github.com/hashicorp/vault/sdk/helper/docker"
	"github.com/hashicorp/vault/sdk/helper/testcluster"
	uberAtomic "go.uber.org/atomic"
	"golang.org/x/net/http2"
)

const dockerVersion = "1.40"

var (
	_ testcluster.VaultCluster     = &DockerCluster{}
	_ testcluster.VaultClusterNode = &dockerClusterNode{}
)

// DockerCluster is used to managing the lifecycle of the test Vault cluster
type DockerCluster struct {
	ClusterName string

	RaftStorage  bool
	ClusterNodes []*dockerClusterNode

	// Certificate fields
	CACert        *x509.Certificate
	CACertBytes   []byte
	CACertPEM     []byte
	CACertPEMFile string
	CAKey         *ecdsa.PrivateKey
	CAKeyPEM      []byte
	RootCAs       *x509.CertPool

	barrierKeys  [][]byte
	recoveryKeys [][]byte
	tmpDir       string

	// rootToken is the initial root token created when the Vault cluster is
	// created.
	rootToken string
	dockerAPI *docker.Client
	ID        string
}

func (dc *DockerCluster) NamedLogger(s string) log.Logger {
	// TODO implement me
	panic("implement me")
}

func (dc *DockerCluster) ClusterID() string {
	return dc.ID
}

func (dc *DockerCluster) Nodes() []testcluster.VaultClusterNode {
	ret := make([]testcluster.VaultClusterNode, len(dc.ClusterNodes))
	for i := range dc.ClusterNodes {
		ret[i] = dc.ClusterNodes[i]
	}
	return ret
}

func (dc *DockerCluster) GetBarrierKeys() [][]byte {
	return dc.barrierKeys
}

func testKeyCopy(key []byte) []byte {
	result := make([]byte, len(key))
	copy(result, key)
	return result
}

func (dc *DockerCluster) GetRecoveryKeys() [][]byte {
	ret := make([][]byte, len(dc.recoveryKeys))
	for i, k := range dc.recoveryKeys {
		ret[i] = testKeyCopy(k)
	}
	return ret
}

func (dc *DockerCluster) GetBarrierOrRecoveryKeys() [][]byte {
	return dc.GetBarrierKeys()
}

func (dc *DockerCluster) SetBarrierKeys(keys [][]byte) {
	dc.barrierKeys = make([][]byte, len(keys))
	for i, k := range keys {
		dc.barrierKeys[i] = testKeyCopy(k)
	}
}

func (dc *DockerCluster) SetRecoveryKeys(keys [][]byte) {
	dc.recoveryKeys = make([][]byte, len(keys))
	for i, k := range keys {
		dc.recoveryKeys[i] = testKeyCopy(k)
	}
}

func (dc *DockerCluster) GetCACertPEMFile() string {
	return dc.CACertPEMFile
}

func (dc *DockerCluster) Cleanup() {
	dc.cleanup()
}

func (dc *DockerCluster) cleanup() error {
	var result *multierror.Error
	for _, node := range dc.ClusterNodes {
		if err := node.cleanup(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result.ErrorOrNil()
}

// RootToken returns the root token of the cluster, if set
func (dc *DockerCluster) RootToken() string {
	return dc.rootToken
}

func (n *dockerClusterNode) Name() string {
	return n.Cluster.ClusterName + "-" + n.NodeID
}

func (dc *DockerCluster) setupNode0(ctx context.Context) error {
	client := dc.ClusterNodes[0].Client

	var resp *api.InitResponse
	var err error
	for ctx.Err() == nil {
		resp, err = client.Sys().Init(&api.InitRequest{
			SecretShares:    3,
			SecretThreshold: 3,
		})
		if err == nil && resp != nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if err != nil {
		return err
	}
	if resp == nil {
		return fmt.Errorf("nil response to init request")
	}

	for _, k := range resp.Keys {
		raw, err := hex.DecodeString(k)
		if err != nil {
			return err
		}
		dc.barrierKeys = append(dc.barrierKeys, raw)
	}

	for _, k := range resp.RecoveryKeys {
		raw, err := hex.DecodeString(k)
		if err != nil {
			return err
		}
		dc.recoveryKeys = append(dc.recoveryKeys, raw)
	}

	dc.rootToken = resp.RootToken
	client.SetToken(dc.rootToken)
	dc.ClusterNodes[0].Client = client

	err = testcluster.UnsealNode(ctx, dc, 0)
	if err != nil {
		return err
	}

	err = ensureLeaderMatches(ctx, client, func(leader *api.LeaderResponse) error {
		if !leader.IsSelf {
			return fmt.Errorf("node %d leader=%v, expected=%v", 0, leader.IsSelf, true)
		}

		return nil
	})

	status, err := client.Sys().SealStatusWithContext(ctx)
	if err != nil {
		return err
	}
	dc.ID = status.ClusterID
	return err
}

func (dc *DockerCluster) clusterReady(ctx context.Context) error {
	for i, node := range dc.ClusterNodes {
		expectLeader := i == 0
		err := ensureLeaderMatches(ctx, node.Client, func(leader *api.LeaderResponse) error {
			if expectLeader != leader.IsSelf {
				return fmt.Errorf("node %d leader=%v, expected=%v", i, leader.IsSelf, expectLeader)
			}

			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (dc *DockerCluster) setupCA(opts *DockerClusterOptions) error {
	var err error

	var caKey *ecdsa.PrivateKey
	if opts != nil && opts.CAKey != nil {
		caKey = opts.CAKey
	} else {
		caKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return err
		}
	}
	dc.CAKey = caKey

	var caBytes []byte
	if opts != nil && len(opts.CACert) > 0 {
		caBytes = opts.CACert
	} else {
		serialNumber := mathrand.New(mathrand.NewSource(time.Now().UnixNano())).Int63()
		CACertTemplate := &x509.Certificate{
			Subject: pkix.Name{
				CommonName: "localhost",
			},
			KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
			SerialNumber:          big.NewInt(serialNumber),
			NotBefore:             time.Now().Add(-30 * time.Second),
			NotAfter:              time.Now().Add(262980 * time.Hour),
			BasicConstraintsValid: true,
			IsCA:                  true,
		}
		caBytes, err = x509.CreateCertificate(rand.Reader, CACertTemplate, CACertTemplate, caKey.Public(), caKey)
		if err != nil {
			return err
		}
	}
	CACert, err := x509.ParseCertificate(caBytes)
	if err != nil {
		return err
	}
	dc.CACert = CACert
	dc.CACertBytes = caBytes

	dc.RootCAs = x509.NewCertPool()
	dc.RootCAs.AddCert(CACert)

	CACertPEMBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	}
	dc.CACertPEM = pem.EncodeToMemory(CACertPEMBlock)

	dc.CACertPEMFile = filepath.Join(dc.tmpDir, "ca", "ca.pem")
	err = os.WriteFile(dc.CACertPEMFile, dc.CACertPEM, 0o755)
	if err != nil {
		return err
	}

	marshaledCAKey, err := x509.MarshalECPrivateKey(caKey)
	if err != nil {
		return err
	}
	CAKeyPEMBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: marshaledCAKey,
	}
	dc.CAKeyPEM = pem.EncodeToMemory(CAKeyPEMBlock)

	return nil
}

func (n *dockerClusterNode) setupCert(ip string) error {
	var err error

	n.ServerKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}

	serialNumber := mathrand.New(mathrand.NewSource(time.Now().UnixNano())).Int63()
	certTemplate := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: n.Name(),
		},
		DNSNames:    []string{"localhost", n.Name()},
		IPAddresses: []net.IP{net.IPv6loopback, net.ParseIP("127.0.0.1"), net.ParseIP(ip)},
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement,
		SerialNumber: big.NewInt(serialNumber),
		NotBefore:    time.Now().Add(-30 * time.Second),
		NotAfter:     time.Now().Add(262980 * time.Hour),
	}
	n.ServerCertBytes, err = x509.CreateCertificate(rand.Reader, certTemplate, n.Cluster.CACert, n.ServerKey.Public(), n.Cluster.CAKey)
	if err != nil {
		return err
	}
	n.ServerCert, err = x509.ParseCertificate(n.ServerCertBytes)
	if err != nil {
		return err
	}
	n.ServerCertPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: n.ServerCertBytes,
	})

	marshaledKey, err := x509.MarshalECPrivateKey(n.ServerKey)
	if err != nil {
		return err
	}
	n.ServerKeyPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: marshaledKey,
	})

	n.ServerCertPEMFile = filepath.Join(n.WorkDir, "cert.pem")
	err = os.WriteFile(n.ServerCertPEMFile, n.ServerCertPEM, 0o755)
	if err != nil {
		return err
	}

	n.ServerKeyPEMFile = filepath.Join(n.WorkDir, "key.pem")
	err = os.WriteFile(n.ServerKeyPEMFile, n.ServerKeyPEM, 0o755)
	if err != nil {
		return err
	}

	tlsCert, err := tls.X509KeyPair(n.ServerCertPEM, n.ServerKeyPEM)
	if err != nil {
		return err
	}

	certGetter := NewCertificateGetter(n.ServerCertPEMFile, n.ServerKeyPEMFile, "")
	if err := certGetter.Reload(); err != nil {
		return err
	}
	tlsConfig := &tls.Config{
		Certificates:   []tls.Certificate{tlsCert},
		RootCAs:        n.Cluster.RootCAs,
		ClientCAs:      n.Cluster.RootCAs,
		ClientAuth:     tls.RequestClientCert,
		NextProtos:     []string{"h2", "http/1.1"},
		GetCertificate: certGetter.GetCertificate,
	}

	n.tlsConfig = tlsConfig

	err = os.WriteFile(filepath.Join(n.WorkDir, "ca.pem"), n.Cluster.CACertPEM, 0o755)
	if err != nil {
		return err
	}
	return nil
}

func NewTestDockerCluster(t *testing.T, opts *DockerClusterOptions) *DockerCluster {
	if opts == nil {
		opts = &DockerClusterOptions{}
	}
	if opts.ClusterName == "" {
		opts.ClusterName = strings.ReplaceAll(t.Name(), "/", "-")
	}
	if opts.Logger == nil {
		opts.Logger = logging.NewVaultLogger(log.Trace).Named(t.Name()) // .Named("container")
	}
	if opts.NetworkName == "" {
		opts.NetworkName = os.Getenv("TEST_DOCKER_NETWORK_NAME")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	t.Cleanup(cancel)

	dc, err := NewDockerCluster(ctx, opts)
	if err != nil {
		t.Fatal(err)
	}
	return dc
}

func NewDockerCluster(ctx context.Context, opts *DockerClusterOptions) (*DockerCluster, error) {
	cli, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithVersion(dockerVersion))
	if err != nil {
		return nil, err
	}

	dc := &DockerCluster{
		dockerAPI:   cli,
		RaftStorage: true,
		ClusterName: opts.ClusterName,
	}

	if err := dc.setupDockerCluster(ctx, opts); err != nil {
		dc.Cleanup()
		return nil, err
	}

	return dc, nil
}

// DockerClusterNode represents a single instance of Vault in a cluster
type dockerClusterNode struct {
	NodeID            string
	HostPort          string
	Client            *api.Client
	ServerCert        *x509.Certificate
	ServerCertBytes   []byte
	ServerCertPEM     []byte
	ServerCertPEMFile string
	ServerKey         *ecdsa.PrivateKey
	ServerKeyPEM      []byte
	ServerKeyPEMFile  string
	tlsConfig         *tls.Config
	WorkDir           string
	Cluster           *DockerCluster
	container         *types.ContainerJSON
	dockerAPI         *docker.Client
	runner            *dockhelper.Runner
	Logger            log.Logger
	RealAPIAddr       string
	cleanupContainer  func()
}

func (n *dockerClusterNode) TLSConfig() *tls.Config {
	return n.tlsConfig.Clone()
}

func (n *dockerClusterNode) APIClient() *api.Client {
	return n.Client
}

func (n *dockerClusterNode) Pause(ctx context.Context) error {
	return n.dockerAPI.ContainerPause(ctx, n.container.ID)
}

func (n *dockerClusterNode) AddNetworkDelay(ctx context.Context, delay time.Duration, targetIP string) error {
	stdout, stderr, exitCode, err := n.runner.RunCmdWithOutput(ctx, n.container.ID, []string{
		"/bin/sh",
		"-xc", strings.Join([]string{
			"echo isolating node",
			"apk add iproute2",
			"tc qdisc add dev eth0 root handle 1: prio",
			fmt.Sprintf("tc filter add dev eth0 parent 1:0 protocol ip pref 55 handle ::55 u32 match ip dst %s flowid 2:1", targetIP),
			fmt.Sprintf("tc qdisc add dev eth0 parent 1:1 handle 2: netem delay %dms", delay/time.Millisecond),
		}, " && "),
	})
	if err != nil {
		return err
	}

	n.Logger.Trace(string(stdout))
	n.Logger.Trace(string(stderr))
	if exitCode != 0 {
		return fmt.Errorf("got nonzero exit code from iptables: %d", exitCode)
	}
	return nil
}

// NewAPIClient creates and configures a Vault API client to communicate with
// the running Vault Cluster for this DockerClusterNode
func (n *dockerClusterNode) apiConfig() (*api.Config, error) {
	transport := cleanhttp.DefaultPooledTransport()
	transport.TLSClientConfig = n.TLSConfig()
	if err := http2.ConfigureTransport(transport); err != nil {
		return nil, err
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
		return nil, config.Error
	}
	config.Address = fmt.Sprintf("https://%s", n.HostPort)
	config.HttpClient = client
	config.MaxRetries = 0
	return config, nil
}

func (n *dockerClusterNode) NewAPIClient() (*api.Client, error) {
	config, err := n.apiConfig()
	if err != nil {
		return nil, err
	}
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	client.SetToken(n.Cluster.RootToken())
	return client, nil
}

// Cleanup kills the container of the node
func (n *dockerClusterNode) Cleanup() {
	n.cleanup()
}

func (n *dockerClusterNode) cleanup() error {
	if n.container == nil || n.container.ID == "" {
		return nil
	}
	n.cleanupContainer()
	return nil
}

func (n *dockerClusterNode) start(ctx context.Context, caDir string, opts *DockerClusterOptions) error {
	vaultCfg := map[string]interface{}{}
	vaultCfg["listener"] = map[string]interface{}{
		"tcp": map[string]interface{}{
			"address":       fmt.Sprintf("%s:%d", "0.0.0.0", 8200),
			"tls_cert_file": "/vault/config/cert.pem",
			"tls_key_file":  "/vault/config/key.pem",
			"telemetry": map[string]interface{}{
				"unauthenticated_metrics_access": true,
			},
		},
	}
	vaultCfg["telemetry"] = map[string]interface{}{
		"disable_hostname": true,
	}
	raftOpts := map[string]interface{}{
		// TODO add options from vnc
		"path":    "/vault/file",
		"node_id": n.NodeID,
	}
	vaultCfg["storage"] = map[string]interface{}{
		"raft": raftOpts,
	}
	if opts != nil && opts.VaultNodeConfig != nil && len(opts.VaultNodeConfig.StorageOptions) > 0 {
		for k, v := range opts.VaultNodeConfig.StorageOptions {
			if _, ok := raftOpts[k].(string); !ok {
				raftOpts[k] = v
			}
		}
	}
	//// disable_mlock is required for working in the Docker environment with
	//// custom plugins
	vaultCfg["disable_mlock"] = true
	vaultCfg["api_addr"] = `https://{{- GetAllInterfaces | exclude "flags" "loopback" | attr "address" -}}:8200`
	vaultCfg["cluster_addr"] = `https://{{- GetAllInterfaces | exclude "flags" "loopback" | attr "address" -}}:8201`

	systemJSON, err := json.Marshal(vaultCfg)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(n.WorkDir, "system.json"), systemJSON, 0o644)
	if err != nil {
		return err
	}

	if opts.VaultNodeConfig != nil {
		localCfg := *opts.VaultNodeConfig
		if opts.VaultNodeConfig.LicensePath != "" {
			b, err := os.ReadFile(opts.VaultNodeConfig.LicensePath)
			if err != nil || len(b) == 0 {
				return fmt.Errorf("unable to read LicensePath at %q: %w", opts.VaultNodeConfig.LicensePath, err)
			}
			localCfg.LicensePath = "/vault/config/license"
			dest := filepath.Join(n.WorkDir, "license")
			err = os.WriteFile(dest, b, 0o644)
			if err != nil {
				return fmt.Errorf("error writing license to %q: %w", dest, err)
			}

		}
		userJSON, err := json.Marshal(localCfg)
		if err != nil {
			return err
		}
		err = os.WriteFile(filepath.Join(n.WorkDir, "user.json"), userJSON, 0o644)
		if err != nil {
			return err
		}
	}

	// Create a temporary cert so vault will start up
	err = n.setupCert("127.0.0.1")
	if err != nil {
		return err
	}

	// setup plugin bin copy if needed
	copyFromTo := map[string]string{
		n.WorkDir: "/vault/config",
		caDir:     "/usr/local/share/ca-certificates/",
	}

	repo := opts.ImageRepo
	if repo == "" {
		repo = "vault"
	}
	tag := opts.ImageTag
	if tag == "" {
		tag = "latest"
	}

	var wg sync.WaitGroup
	wg.Add(1)
	var seenLogs uberAtomic.Bool
	logConsumer := func(s string) {
		if seenLogs.CAS(false, true) {
			wg.Done()
		}
		n.Logger.Trace(s)
	}
	r, err := dockhelper.NewServiceRunner(dockhelper.RunOptions{
		ImageRepo: repo,
		ImageTag:  tag,
		// We don't need to run update-ca-certificates in the container, because
		// we're providing the CA in the raft join call, and otherwise Vault
		// servers don't talk to one another on the API port.
		Cmd: []string{"server"},
		Env: []string{
			// For now we're using disable_mlock, because this is for testing
			// anyway, and because it prevents us using external plugins.
			"SKIP_SETCAP=true",
		},
		Ports:           []string{"8200/tcp", "8201/tcp"},
		ContainerName:   n.Name(),
		NetworkName:     opts.NetworkName,
		CopyFromTo:      copyFromTo,
		LogConsumer:     logConsumer,
		PreDelete:       true,
		DoNotAutoRemove: true,
		PostStart: func(containerID string, realIP string) error {
			err := n.setupCert(realIP)
			if err != nil {
				return err
			}

			// If we signal Vault before it installs its sighup handler, it'll die.
			wg.Wait()
			n.Logger.Trace("running poststart", "containerID", containerID, "IP", realIP)
			return n.runner.RefreshFiles(ctx, containerID)
		},
		Capabilities: []string{"NET_ADMIN"},
	})
	if err != nil {
		return err
	}
	n.runner = r

	svc, _, err := r.StartNewService(ctx, false, false, func(ctx context.Context, host string, port int) (dockhelper.ServiceConfig, error) {
		config, err := n.apiConfig()
		if err != nil {
			return nil, err
		}
		config.Address = fmt.Sprintf("https://%s:%d", host, port)
		client, err := api.NewClient(config)
		if err != nil {
			return nil, err
		}
		_, err = client.Sys().SealStatus()
		if err != nil {
			return nil, err
		}
		return dockhelper.NewServiceHostPort(host, port), nil
	})
	if err != nil {
		return err
	}

	n.HostPort = svc.Config.Address()
	n.container = svc.Container
	netName := opts.NetworkName
	if netName == "" {
		netName = "bridge"
	}
	n.RealAPIAddr = "https://" + svc.Container.NetworkSettings.Networks[netName].IPAddress + ":8200"
	n.cleanupContainer = svc.Cleanup

	return nil
}

// DockerClusterOptions has options for setting up the docker cluster
type DockerClusterOptions struct {
	testcluster.ClusterOptions
	CAKey       *ecdsa.PrivateKey
	NetworkName string
	ImageRepo   string
	ImageTag    string
	CloneCA     *DockerCluster
}

func ensureLeaderMatches(ctx context.Context, client *api.Client, ready func(response *api.LeaderResponse) error) error {
	var leader *api.LeaderResponse
	var err error
	for ctx.Err() == nil {
		leader, err = client.Sys().Leader()
		switch {
		case err != nil:
		case leader == nil:
			err = fmt.Errorf("nil response to leader check")
		default:
			err = ready(leader)
			if err == nil {
				return nil
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("error checking leader: %v", err)
}

const DefaultNumCores = 3

// creates a managed docker container running Vault
func (dc *DockerCluster) setupDockerCluster(ctx context.Context, opts *DockerClusterOptions) error {
	if opts == nil {
		opts = &DockerClusterOptions{}
	}
	if opts.Logger == nil {
		opts.Logger = log.NewNullLogger()
	}

	if opts.TmpDir != "" {
		if _, err := os.Stat(opts.TmpDir); os.IsNotExist(err) {
			if err := os.MkdirAll(opts.TmpDir, 0o700); err != nil {
				return err
			}
		}
		dc.tmpDir = opts.TmpDir
	} else {
		tempDir, err := ioutil.TempDir("", "vault-test-cluster-")
		if err != nil {
			return err
		}
		dc.tmpDir = tempDir
	}
	caDir := filepath.Join(dc.tmpDir, "ca")
	if err := os.MkdirAll(caDir, 0o755); err != nil {
		return err
	}

	var numCores int
	if opts.NumCores == 0 {
		numCores = DefaultNumCores
	} else {
		numCores = opts.NumCores
	}

	if opts.CloneCA != nil {
		dc.CACert = opts.CloneCA.CACert
		dc.CACertBytes = opts.CloneCA.CACertBytes
		dc.CACertPEM = opts.CloneCA.CACertPEM
		dc.CACertPEMFile = opts.CloneCA.CACertPEMFile
		dc.CAKey = opts.CloneCA.CAKey
		dc.CAKeyPEM = opts.CloneCA.CAKeyPEM
		dc.RootCAs = opts.CloneCA.RootCAs
	} else {
		if err := dc.setupCA(opts); err != nil {
			return err
		}
	}

	for i := 0; i < numCores; i++ {
		if err := dc.addNode(ctx, opts); err != nil {
			return err
		}
		if opts.SkipInit {
			continue
		}
		if i == 0 {
			if err := dc.setupNode0(ctx); err != nil {
				return nil
			}
		} else {
			if err := dc.joinNode(ctx, i); err != nil {
				return err
			}
		}
	}

	return nil
}

func (dc *DockerCluster) AddNode(ctx context.Context, opts *DockerClusterOptions) error {
	if err := dc.addNode(ctx, opts); err != nil {
		return err
	}

	return dc.joinNode(ctx, len(dc.ClusterNodes)-1)
}

func (dc *DockerCluster) addNode(ctx context.Context, opts *DockerClusterOptions) error {
	i := len(dc.ClusterNodes)
	nodeID := fmt.Sprintf("core-%d", i)
	node := &dockerClusterNode{
		dockerAPI: dc.dockerAPI,
		NodeID:    nodeID,
		Cluster:   dc,
		WorkDir:   filepath.Join(dc.tmpDir, nodeID),
		Logger:    opts.Logger.Named(nodeID),
	}
	dc.ClusterNodes = append(dc.ClusterNodes, node)
	if err := os.MkdirAll(node.WorkDir, 0o755); err != nil {
		return err
	}
	if err := node.start(ctx, filepath.Join(dc.tmpDir, "ca"), opts); err != nil {
		return err
	}
	client, err := node.NewAPIClient()
	if err != nil {
		return err
	}
	client.SetToken(dc.rootToken)
	node.Client = client
	return nil
}

func (dc *DockerCluster) joinNode(ctx context.Context, nodeIdx int) error {
	leaderIdx, err := testcluster.LeaderNode(ctx, dc)
	if err != nil {
		return err
	}
	leader := dc.ClusterNodes[leaderIdx]

	if nodeIdx >= len(dc.ClusterNodes) {
		return fmt.Errorf("invalid node %d", nodeIdx)
	}
	node := dc.ClusterNodes[nodeIdx]
	client := node.APIClient()

	var resp *api.RaftJoinResponse
	resp, err = client.Sys().RaftJoinWithContext(ctx, &api.RaftJoinRequest{
		// When running locally on a bridge network, the containers must use their
		// actual (private) IP to talk to one another.  Our code must instead use
		// the portmapped address since we're not on their network in that case.
		LeaderAPIAddr:    leader.RealAPIAddr,
		LeaderCACert:     string(dc.CACertPEM),
		LeaderClientCert: string(node.ServerCertPEM),
		LeaderClientKey:  string(node.ServerKeyPEM),
	})
	if resp == nil || !resp.Joined {
		return fmt.Errorf("nil or negative response from raft join request: %v", resp)
	}
	if err != nil {
		return fmt.Errorf("failed to join cluster: %w", err)
	}

	return testcluster.UnsealNode(ctx, dc, nodeIdx)
}

/* Notes on testing the non-bridge network case:
- you need the test itself to be running in a container so that it can use
  the network; create the network using
    docker network create testvault
- this means that you need to mount the docker socket in that test container,
  but on macos there's stuff that prevents that from working; to hack that,
  on the host run
    sudo ln -s "$HOME/Library/Containers/com.docker.docker/Data/docker.raw.sock" /var/run/docker.sock.raw
- run the test container like
    docker run --rm -it --network testvault \
      -v /var/run/docker.sock.raw:/var/run/docker.sock \
      -v $(pwd):/home/circleci/go/src/github.com/hashicorp/vault/ \
      -w /home/circleci/go/src/github.com/hashicorp/vault/ \
      "docker.mirror.hashicorp.services/cimg/go:1.19.2" /bin/bash
- in the container you may need to chown/chmod /var/run/docker.sock; use `docker ps`
  to test if it's working

*/
