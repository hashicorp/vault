package testcluster

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/logging"
)

type ExecDevCluster struct {
	ID                 string
	ClusterName        string
	ClusterNodes       []*execDevClusterNode
	CACertPEMFile      string
	barrierKeys        [][]byte
	recoveryKeys       [][]byte
	tmpDir             string
	clientAuthRequired bool
	rootToken          string
	stop               func()
	stopCh             chan struct{}
	Logger             log.Logger
}

func (dc *ExecDevCluster) NamedLogger(s string) log.Logger {
	return dc.Logger.Named(s)
}

var _ VaultCluster = &ExecDevCluster{}

type ExecDevClusterOptions struct {
	ClusterOptions
	BinaryPath string
}

func NewTestExecDevCluster(t *testing.T, opts *ExecDevClusterOptions) *ExecDevCluster {
	if opts == nil {
		opts = &ExecDevClusterOptions{}
	}
	if opts.ClusterName == "" {
		opts.ClusterName = strings.ReplaceAll(t.Name(), "/", "-")
	}
	if opts.Logger == nil {
		opts.Logger = logging.NewVaultLogger(log.Trace).Named(t.Name()) // .Named("container")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)

	dc, err := NewExecDevCluster(ctx, opts)
	if err != nil {
		t.Fatal(err)
	}
	return dc
}

func NewExecDevCluster(ctx context.Context, opts *ExecDevClusterOptions) (*ExecDevCluster, error) {
	dc := &ExecDevCluster{
		ClusterName: opts.ClusterName,
		stopCh:      make(chan struct{}),
	}

	if opts == nil {
		opts = &ExecDevClusterOptions{}
	}
	if opts.NumCores == 0 {
		opts.NumCores = 3
	}
	if err := dc.setupExecDevCluster(ctx, opts); err != nil {
		dc.Cleanup()
		return nil, err
	}

	return dc, nil
}

func (dc *ExecDevCluster) setupExecDevCluster(ctx context.Context, opts *ExecDevClusterOptions) (retErr error) {
	if opts == nil {
		opts = &ExecDevClusterOptions{}
	}
	if opts.Logger == nil {
		opts.Logger = log.NewNullLogger()
	}
	dc.Logger = opts.Logger

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

	// This context is used to stop the subprocess
	execCtx, cancel := context.WithCancel(context.Background())
	dc.stop = func() {
		cancel()
		close(dc.stopCh)
	}
	defer func() {
		if retErr != nil {
			cancel()
		}
	}()

	bin := opts.BinaryPath
	if bin == "" {
		bin = "vault"
	}

	clusterJsonPath := filepath.Join(dc.tmpDir, "cluster.json")
	args := []string{"server", "-dev", "-dev-cluster-json", clusterJsonPath}
	switch {
	case opts.NumCores == 3:
		args = append(args, "-dev-three-node")
	case opts.NumCores == 1:
		args = append(args, "-dev-tls")
	default:
		return fmt.Errorf("NumCores=1 and NumCores=3 are the only supported options right now")
	}
	cmd := exec.CommandContext(execCtx, bin, args...)
	if opts.Logger != nil {
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}
		go func() {
			// TODO set bigger buffer
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				opts.Logger.Trace(scanner.Text())
			}
		}()
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return err
		}
		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				opts.Logger.Trace(scanner.Text())
			}
		}()
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	ctx, _ = context.WithTimeout(ctx, 5*time.Second)
	for ctx.Err() == nil {
		if b, err := os.ReadFile(clusterJsonPath); err == nil && len(b) > 0 {
			var clusterJson ClusterJson
			if err := jsonutil.DecodeJSON(b, &clusterJson); err != nil {
				continue
			}
			dc.CACertPEMFile = clusterJson.CACertPath
			dc.rootToken = clusterJson.RootToken
			for i, node := range clusterJson.Nodes {
				config := api.DefaultConfig()
				config.Address = node.APIAddress
				err := config.ConfigureTLS(&api.TLSConfig{
					CACert: clusterJson.CACertPath,
				})
				if err != nil {
					return err
				}
				client, err := api.NewClient(config)
				if err != nil {
					return err
				}
				client.SetToken(dc.rootToken)
				_, err = client.Sys().ListMounts()
				if err != nil {
					return err
				}

				dc.ClusterNodes = append(dc.ClusterNodes, &execDevClusterNode{
					name:   fmt.Sprintf("core-%d", i),
					Client: client,
				})
			}
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return ctx.Err()
}

type execDevClusterNode struct {
	// HostPort          string
	name   string
	Client *api.Client
}

var _ VaultClusterNode = &execDevClusterNode{}

func (e *execDevClusterNode) Name() string {
	return e.name
}

func (e *execDevClusterNode) APIClient() *api.Client {
	return e.Client
}

func (e *execDevClusterNode) TLSConfig() *tls.Config {
	return e.Client.CloneConfig().TLSConfig()
}

func (dc *ExecDevCluster) ClusterID() string {
	return dc.ID
}

func (dc *ExecDevCluster) Nodes() []VaultClusterNode {
	ret := make([]VaultClusterNode, len(dc.ClusterNodes))
	for i := range dc.ClusterNodes {
		ret[i] = dc.ClusterNodes[i]
	}
	return ret
}

func (dc *ExecDevCluster) GetBarrierKeys() [][]byte {
	return dc.barrierKeys
}

func testKeyCopy(key []byte) []byte {
	result := make([]byte, len(key))
	copy(result, key)
	return result
}

func (dc *ExecDevCluster) GetRecoveryKeys() [][]byte {
	ret := make([][]byte, len(dc.recoveryKeys))
	for i, k := range dc.recoveryKeys {
		ret[i] = testKeyCopy(k)
	}
	return ret
}

func (dc *ExecDevCluster) GetBarrierOrRecoveryKeys() [][]byte {
	return dc.GetBarrierKeys()
}

func (dc *ExecDevCluster) SetBarrierKeys(keys [][]byte) {
	dc.barrierKeys = make([][]byte, len(keys))
	for i, k := range keys {
		dc.barrierKeys[i] = testKeyCopy(k)
	}
}

func (dc *ExecDevCluster) SetRecoveryKeys(keys [][]byte) {
	dc.recoveryKeys = make([][]byte, len(keys))
	for i, k := range keys {
		dc.recoveryKeys[i] = testKeyCopy(k)
	}
}

func (dc *ExecDevCluster) GetCACertPEMFile() string {
	return dc.CACertPEMFile
}

func (dc *ExecDevCluster) Cleanup() {
	dc.stop()
}

// RootToken returns the root token of the cluster, if set
func (dc *ExecDevCluster) RootToken() string {
	return dc.rootToken
}
