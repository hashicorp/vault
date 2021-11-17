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
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	mathrand "math/rand"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	docker "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-multierror"
	uuid "github.com/hashicorp/go-uuid"
	stepwise "github.com/hashicorp/vault-testing-stepwise"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"golang.org/x/net/http2"
)

var _ stepwise.Environment = (*DockerCluster)(nil)

const dockerVersion = "1.40"

// DockerCluster is used to managing the lifecycle of the test Vault cluster
type DockerCluster struct {
	ID string
	// PluginName is the input from the test case
	PluginName string
	// ClusterName is a UUID name of the cluster.
	ClusterName string

	// MountOptions are a set of options for registering and mounting the plugin
	MountOptions stepwise.MountOptions

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

	// networkID tracks the network ID of the created docker network
	networkID string

	barrierKeys  [][]byte
	recoveryKeys [][]byte
	tmpDir       string

	clientAuthRequired bool
	// the mountpath of the plugin under test
	mountPath string
	// rootToken is the initial root token created when the Vault cluster is
	// created.
	rootToken string
}

// Teardown stops all the containers.
func (dc *DockerCluster) Teardown() error {
	var result *multierror.Error
	for _, node := range dc.ClusterNodes {
		if err := node.Cleanup(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// clean up networks
	if dc.networkID != "" {
		cli, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithVersion(dockerVersion))
		if err != nil {
			return multierror.Append(result, err)
		}
		if err := cli.NetworkRemove(context.Background(), dc.networkID); err != nil {
			return multierror.Append(result, err)
		}
	}

	return result.ErrorOrNil()
}

// MountPath returns the path that the plugin under test is mounted at. If a
// MountPathPrefix was given, the mount path uses the prefix with a uuid
// appended. The default is the given PluginName with a uuid suffix.
func (dc *DockerCluster) MountPath() string {
	if dc.mountPath != "" {
		return dc.mountPath
	}

	uuidStr, err := uuid.GenerateUUID()
	if err != nil {
		panic(err)
	}

	prefix := dc.PluginName
	if dc.MountOptions.MountPathPrefix != "" {
		prefix = dc.MountOptions.MountPathPrefix
	}

	dc.mountPath = fmt.Sprintf("%s_%s", prefix, uuidStr)
	if dc.MountOptions.PluginType == stepwise.PluginTypeCredential {
		dc.mountPath = path.Join("auth", dc.mountPath)
	}

	return dc.mountPath
}

// RootToken returns the root token of the cluster, if set
func (dc *DockerCluster) RootToken() string {
	return dc.rootToken
}

// Name returns the name of this environment
func (dc *DockerCluster) Name() string {
	return "docker"
}

// Client returns a clone of the configured Vault API client.
func (dc *DockerCluster) Client() (*api.Client, error) {
	if len(dc.ClusterNodes) > 0 {
		if dc.ClusterNodes[0].Client != nil {
			c, err := dc.ClusterNodes[0].Client.Clone()
			if err != nil {
				return nil, err
			}
			c.SetToken(dc.ClusterNodes[0].Client.Token())
			return c, nil
		}
	}

	return nil, errors.New("no configured client found")
}

func (n *dockerClusterNode) Name() string {
	return n.Cluster.ClusterName + "-" + n.NodeID
}

func (dc *DockerCluster) Initialize(ctx context.Context) error {
	client, err := dc.ClusterNodes[0].NewAPIClient()
	if err != nil {
		return err
	}

	var resp *api.InitResponse
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

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Unseal
	for j, node := range dc.ClusterNodes {
		// copy the index value, so we're not reusing it in deeper scopes
		i := j
		client, err := node.NewAPIClient()
		if err != nil {
			return err
		}
		node.Client = client

		if i > 0 && dc.RaftStorage {
			leader := dc.ClusterNodes[0]
			resp, err := client.Sys().RaftJoin(&api.RaftJoinRequest{
				LeaderAPIAddr:    fmt.Sprintf("https://%s:%d", dc.ClusterNodes[0].Name(), leader.Address.Port),
				LeaderCACert:     string(dc.CACertPEM),
				LeaderClientCert: string(node.ServerCertPEM),
				LeaderClientKey:  string(node.ServerKeyPEM),
			})
			if err != nil {
				return err
			}
			if resp == nil || !resp.Joined {
				return fmt.Errorf("nil or negative response from raft join request: %v", resp)
			}
		}

		var unsealed bool
		for _, key := range dc.barrierKeys {
			resp, err := client.Sys().Unseal(hex.EncodeToString(key))
			if err != nil {
				return err
			}
			unsealed = !resp.Sealed
		}
		if i == 0 && !unsealed {
			return fmt.Errorf("could not unseal node %d", i)
		}
		client.SetToken(dc.rootToken)

		err = ensureHealthMatches(ctx, node.Client, func(health *api.HealthResponse) error {
			if health.Sealed {
				return fmt.Errorf("node %d is sealed: %#v", i, health)
			}
			if health.ClusterID == "" {
				return fmt.Errorf("node %d has no cluster ID", i)
			}

			dc.ID = health.ClusterID
			return nil
		})
		if err != nil {
			return err
		}

		if i == 0 {
			err = ensureLeaderMatches(ctx, node.Client, func(leader *api.LeaderResponse) error {
				if !leader.IsSelf {
					return fmt.Errorf("node %d leader=%v, expected=%v", i, leader.IsSelf, true)
				}

				return nil
			})
			if err != nil {
				return err
			}
		}
	}

	for i, node := range dc.ClusterNodes {
		expectLeader := i == 0
		err = ensureLeaderMatches(ctx, node.Client, func(leader *api.LeaderResponse) error {
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

	certIPs := []net.IP{
		net.IPv6loopback,
		net.ParseIP("127.0.0.1"),
	}

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
			DNSNames:              []string{"localhost"},
			IPAddresses:           certIPs,
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
	err = ioutil.WriteFile(dc.CACertPEMFile, dc.CACertPEM, 0o755)
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

// Don't call this until n.Address.IP is populated
func (n *dockerClusterNode) setupCert() error {
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
		IPAddresses: []net.IP{net.IPv6loopback, net.ParseIP("127.0.0.1")},
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
	err = ioutil.WriteFile(n.ServerCertPEMFile, n.ServerCertPEM, 0o755)
	if err != nil {
		return err
	}

	n.ServerKeyPEMFile = filepath.Join(n.WorkDir, "key.pem")
	err = ioutil.WriteFile(n.ServerKeyPEMFile, n.ServerKeyPEM, 0o755)
	if err != nil {
		return err
	}

	tlsCert, err := tls.X509KeyPair(n.ServerCertPEM, n.ServerKeyPEM)
	if err != nil {
		return err
	}

	certGetter := stepwise.NewCertificateGetter(n.ServerCertPEMFile, n.ServerKeyPEMFile, "")
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

	if n.Cluster.clientAuthRequired {
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}
	n.TLSConfig = tlsConfig

	return nil
}

// NewEnvironment creats a new Stepwise Environment for executing tests
func NewEnvironment(name string, options *stepwise.MountOptions) *DockerCluster {
	if options == nil {
		return nil
	}

	clusterUUID, err := uuid.GenerateUUID()
	if err != nil {
		panic(err)
	}

	return &DockerCluster{
		PluginName:   options.PluginName,
		ClusterName:  fmt.Sprintf("test-%s-%s", name, clusterUUID),
		RaftStorage:  true,
		MountOptions: *options,
	}
}

// DockerClusterNode represents a single instance of Vault in a cluster
type dockerClusterNode struct {
	NodeID            string
	Address           *net.TCPAddr
	HostPort          string
	Client            *api.Client
	ServerCert        *x509.Certificate
	ServerCertBytes   []byte
	ServerCertPEM     []byte
	ServerCertPEMFile string
	ServerKey         *ecdsa.PrivateKey
	ServerKeyPEM      []byte
	ServerKeyPEMFile  string
	TLSConfig         *tls.Config
	WorkDir           string
	Cluster           *DockerCluster
	container         *types.ContainerJSON
	dockerAPI         *docker.Client
}

// NewAPIClient creates and configures a Vault API client to communicate with
// the running Vault Cluster for this DockerClusterNode
func (n *dockerClusterNode) NewAPIClient() (*api.Client, error) {
	transport := cleanhttp.DefaultPooledTransport()
	transport.TLSClientConfig = n.TLSConfig.Clone()
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
	config.Address = fmt.Sprintf("https://127.0.0.1:%s", n.HostPort)
	config.HttpClient = client
	config.MaxRetries = 0
	apiClient, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	apiClient.SetToken(n.Cluster.RootToken())
	return apiClient, nil
}

// Cleanup kills the container of the node
func (n *dockerClusterNode) Cleanup() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return n.dockerAPI.ContainerKill(ctx, n.container.ID, "KILL")
}

func (n *dockerClusterNode) start(cli *docker.Client, caDir, netName string, netCIDR *dockerClusterNode, pluginBinPath string) error {
	n.dockerAPI = cli

	err := n.setupCert()
	if err != nil {
		return err
	}

	vaultCfg := map[string]interface{}{
		"listener": map[string]interface{}{
			"tcp": map[string]interface{}{
				"address":       fmt.Sprintf("%s:%d", "0.0.0.0", 8200),
				"tls_cert_file": "/vault/config/cert.pem",
				"tls_key_file":  "/vault/config/key.pem",
				"telemetry": map[string]interface{}{
					"unauthenticated_metrics_access": true,
				},
			},
		},
		"telemetry": map[string]interface{}{
			"disable_hostname": true,
		},
		"storage": map[string]interface{}{
			"raft": map[string]interface{}{
				"path":    "/vault/file",
				"node_id": n.NodeID,
			},
		},
		"cluster_name":         netName,
		"log_level":            "TRACE",
		"raw_storage_endpoint": true,
		"plugin_directory":     "/vault/config",
		// disable_mlock is required for working in the Docker environment with
		// custom plugins
		"disable_mlock": true,
		// These are being provided by docker-entrypoint now, since we don't know
		// the address before the container starts.
		//"api_addr": fmt.Sprintf("https://%s:%d", n.Address.IP, n.Address.Port),
		//"cluster_addr": fmt.Sprintf("https://%s:%d", n.Address.IP, n.Address.Port+1),
	}
	cfgJSON, err := json.Marshal(vaultCfg)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(n.WorkDir, "local.json"), cfgJSON, 0o644)
	if err != nil {
		return err
	}
	// setup plugin bin copy if needed
	copyFromTo := map[string]string{
		n.WorkDir: "/vault/config",
		caDir:     "/usr/local/share/ca-certificates/",
	}
	if pluginBinPath != "" {
		base := path.Base(pluginBinPath)
		copyFromTo[pluginBinPath] = filepath.Join("/vault/config", base)
	}

	r := &Runner{
		dockerAPI: cli,
		ContainerConfig: &container.Config{
			Image: "vault",
			Entrypoint: []string{"/bin/sh", "-c", "update-ca-certificates && " +
				"exec /usr/local/bin/docker-entrypoint.sh vault server -log-level=trace -dev-plugin-dir=./vault/config -config /vault/config/local.json"},
			Env: []string{
				"VAULT_CLUSTER_INTERFACE=eth0",
				"VAULT_API_ADDR=https://127.0.0.1:8200",
				fmt.Sprintf("VAULT_REDIRECT_ADDR=https://%s:8200", n.Name()),
			},
			Labels:       nil,
			ExposedPorts: nat.PortSet{"8200/tcp": {}, "8201/tcp": {}},
		},
		ContainerName: n.Name(),
		NetName:       netName,
		CopyFromTo:    copyFromTo,
	}

	n.container, err = r.Start(context.Background())
	if err != nil {
		return err
	}

	n.Address = &net.TCPAddr{
		IP:   net.ParseIP(n.container.NetworkSettings.IPAddress),
		Port: 8200,
	}
	ports := n.container.NetworkSettings.NetworkSettingsBase.Ports[nat.Port("8200/tcp")]
	if len(ports) == 0 {
		n.Cleanup()
		return fmt.Errorf("could not find port binding for 8200/tcp")
	}
	n.HostPort = ports[0].HostPort

	return nil
}

// DockerClusterOptions has options for setting up the docker cluster
type DockerClusterOptions struct {
	KeepStandbysSealed bool
	RequireClientAuth  bool
	SkipInit           bool
	CACert             []byte
	NumCores           int
	tmpDir             string
	PluginTestBin      string
	// SetupFunc is called after the cluster is started.
	SetupFunc func(t testing.T, c *DockerCluster)
	CAKey     *ecdsa.PrivateKey
}

//
// helper methods/functions
//

// ensureHealthMatches checks health
func ensureHealthMatches(ctx context.Context, client *api.Client, ready func(response *api.HealthResponse) error) error {
	var health *api.HealthResponse
	var err error
	for ctx.Err() == nil {
		health, err = client.Sys().Health()
		switch {
		case err != nil:
		case health == nil:
			err = fmt.Errorf("nil response to health check")
		default:
			err = ready(health)
			if err == nil {
				return nil
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("error checking health: %v", err)
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

// end test helper methods

// TODO: allow number of cores/servers to be configurable
var DefaultNumCores = 1

// creates a managed docker container running Vault
func (cluster *DockerCluster) setupDockerCluster(opts *DockerClusterOptions) error {
	if opts != nil && opts.tmpDir != "" {
		if _, err := os.Stat(opts.tmpDir); os.IsNotExist(err) {
			if err := os.MkdirAll(opts.tmpDir, 0o700); err != nil {
				return err
			}
		}
		cluster.tmpDir = opts.tmpDir
	} else {
		tempDir, err := ioutil.TempDir("", "vault-test-cluster-")
		if err != nil {
			return err
		}
		cluster.tmpDir = tempDir
	}
	caDir := filepath.Join(cluster.tmpDir, "ca")
	if err := os.MkdirAll(caDir, 0o755); err != nil {
		return err
	}

	var numCores int
	if opts == nil || opts.NumCores == 0 {
		numCores = DefaultNumCores
	} else {
		numCores = opts.NumCores
	}

	if opts != nil && opts.RequireClientAuth {
		cluster.clientAuthRequired = true
	}

	for i := 0; i < numCores; i++ {
		nodeID := fmt.Sprintf("vault-%d", i)
		node := &dockerClusterNode{
			NodeID:  nodeID,
			Cluster: cluster,
			WorkDir: filepath.Join(cluster.tmpDir, nodeID),
		}
		cluster.ClusterNodes = append(cluster.ClusterNodes, node)
		if err := os.MkdirAll(node.WorkDir, 0o700); err != nil {
			return err
		}
	}

	err := cluster.setupCA(opts)
	if err != nil {
		return err
	}

	cli, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithVersion(dockerVersion))
	if err != nil {
		return err
	}

	netUUID, err := uuid.GenerateUUID()
	if err != nil {
		panic(err)
	}
	netName := fmt.Sprintf("%s-%s", "vault-test", netUUID)
	netID, err := setupNetwork(cli, netName)
	if err != nil {
		return err
	}
	cluster.networkID = netID

	for _, node := range cluster.ClusterNodes {
		pluginBinPath := ""
		if opts != nil {
			pluginBinPath = opts.PluginTestBin
		}

		err := node.start(cli, caDir, netName, node, pluginBinPath)
		if err != nil {
			return err
		}
	}

	if opts == nil || !opts.SkipInit {
		if err := cluster.Initialize(context.Background()); err != nil {
			return err
		}
	}

	return nil
}

// Docker networking functions
// setupNetwork establishes networking for the Docker container
func setupNetwork(cli *docker.Client, netName string) (string, error) {
	id, err := createNetwork(cli, netName)
	if err != nil {
		return "", fmt.Errorf("couldn't create network %s: %w", netName, err)
	}
	return id, nil
}

func createNetwork(cli *docker.Client, netName string) (string, error) {
	resp, err := cli.NetworkCreate(context.Background(), netName, types.NetworkCreate{
		CheckDuplicate: true,
		Driver:         "bridge",
		Options:        map[string]string{},
		IPAM: &network.IPAM{
			Driver:  "default",
			Options: map[string]string{},
		},
	})
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

// Setup creates any temp directories needed and compiles the binary for copying to Docker
func (dc *DockerCluster) Setup() error {
	registryName := dc.MountOptions.RegistryName
	pluginName := dc.MountOptions.PluginName

	// get the working directory of the plugin being tested.
	srcDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// tmpDir gets cleaned up when the cluster is cleaned up
	tmpDir, err := ioutil.TempDir("", "bin")
	if err != nil {
		return err
	}

	binName, binPath, sha256value, err := stepwise.CompilePlugin(registryName, pluginName, srcDir, tmpDir)
	if err != nil {
		return err
	}

	dOpts := &DockerClusterOptions{PluginTestBin: binPath}
	if err := dc.setupDockerCluster(dOpts); err != nil {
		return err
	}

	cores := dc.ClusterNodes
	client := cores[0].Client

	// use client to mount plugin
	err = client.Sys().RegisterPlugin(&api.RegisterPluginInput{
		Name:    registryName,
		Type:    consts.PluginType(dc.MountOptions.PluginType),
		Command: binName,
		SHA256:  sha256value,
	})
	if err != nil {
		return err
	}

	switch dc.MountOptions.PluginType {
	case stepwise.PluginTypeCredential:
		// the mount path includes "auth/" for credential type plugins. For enabling
		// auth mounts via the /sys endpoint, we need to remove that prefix
		authPath := strings.TrimPrefix(dc.MountPath(), "auth/")
		err = client.Sys().EnableAuthWithOptions(authPath, &api.EnableAuthOptions{
			Type: registryName,
		})
	case stepwise.PluginTypeDatabase:
	case stepwise.PluginTypeSecrets:
		err = client.Sys().Mount(dc.MountPath(), &api.MountInput{
			Type: registryName,
		})
	default:
		return fmt.Errorf("unknown plugin type: %s", dc.MountOptions.PluginType.String())
	}
	return err
}
