// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package testcluster

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
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

func (dc *ExecDevCluster) SetRootToken(token string) {
	dc.rootToken = token
}

func (dc *ExecDevCluster) NamedLogger(s string) log.Logger {
	return dc.Logger.Named(s)
}

var _ VaultCluster = &ExecDevCluster{}

type ExecDevClusterOptions struct {
	ClusterOptions
	BinaryPath string
	// this is -dev-listen-address, defaults to "127.0.0.1:8200"
	BaseListenAddress string
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
	if opts.VaultLicense == "" {
		opts.VaultLicense = os.Getenv(EnvVaultLicenseCI)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
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
		tempDir, err := os.MkdirTemp("", "vault-test-cluster-")
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
	if opts.BaseListenAddress != "" {
		args = append(args, "-dev-listen-address", opts.BaseListenAddress)
	}
	cmd := exec.CommandContext(execCtx, bin, args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "VAULT_LICENSE="+opts.VaultLicense)
	cmd.Env = append(cmd.Env, "VAULT_LOG_FORMAT=json")
	cmd.Env = append(cmd.Env, "VAULT_DEV_TEMP_DIR="+dc.tmpDir)
	if opts.Logger != nil {
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}
		go func() {
			outlog := opts.Logger.Named("stdout")
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				outlog.Trace(scanner.Text())
			}
		}()
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return err
		}
		go func() {
			errlog := opts.Logger.Named("stderr")
			scanner := bufio.NewScanner(stderr)
			// The default buffer is 4k, and Vault can emit bigger log lines
			scanner.Buffer(make([]byte, 64*1024), bufio.MaxScanTokenSize)
			for scanner.Scan() {
				JSONLogNoTimestamp(errlog, scanner.Text())
			}
		}()
	}

	if err := cmd.Start(); err != nil {
		return err
	}

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
					client: client,
				})
			}
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return ctx.Err()
}

type execDevClusterNode struct {
	name   string
	client *api.Client
}

var _ VaultClusterNode = &execDevClusterNode{}

func (e *execDevClusterNode) Name() string {
	return e.name
}

func (e *execDevClusterNode) APIClient() *api.Client {
	// We clone to ensure that whenever this method is called, the caller gets
	// back a pristine client, without e.g. any namespace or token changes that
	// might pollute a shared client.  We clone the config instead of the
	// client because (1) Client.clone propagates the replicationStateStore and
	// the httpClient pointers, (2) it doesn't copy the tlsConfig at all, and
	// (3) if clone returns an error, it doesn't feel as appropriate to panic
	// below.  Who knows why clone might return an error?
	cfg := e.client.CloneConfig()
	client, err := api.NewClient(cfg)
	if err != nil {
		// It seems fine to panic here, since this should be the same input
		// we provided to NewClient when we were setup, and we didn't panic then.
		// Better not to completely ignore the error though, suppose there's a
		// bug in CloneConfig?
		panic(fmt.Sprintf("NewClient error on cloned config: %v", err))
	}
	client.SetToken(e.client.Token())
	return client
}

func (e *execDevClusterNode) TLSConfig() *tls.Config {
	return e.client.CloneConfig().TLSConfig()
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

func copyKey(key []byte) []byte {
	result := make([]byte, len(key))
	copy(result, key)
	return result
}

func (dc *ExecDevCluster) GetRecoveryKeys() [][]byte {
	ret := make([][]byte, len(dc.recoveryKeys))
	for i, k := range dc.recoveryKeys {
		ret[i] = copyKey(k)
	}
	return ret
}

func (dc *ExecDevCluster) GetBarrierOrRecoveryKeys() [][]byte {
	return dc.GetBarrierKeys()
}

func (dc *ExecDevCluster) SetBarrierKeys(keys [][]byte) {
	dc.barrierKeys = make([][]byte, len(keys))
	for i, k := range keys {
		dc.barrierKeys[i] = copyKey(k)
	}
}

func (dc *ExecDevCluster) SetRecoveryKeys(keys [][]byte) {
	dc.recoveryKeys = make([][]byte, len(keys))
	for i, k := range keys {
		dc.recoveryKeys[i] = copyKey(k)
	}
}

func (dc *ExecDevCluster) GetCACertPEMFile() string {
	return dc.CACertPEMFile
}

func (dc *ExecDevCluster) Cleanup() {
	dc.stop()
}

// GetRootToken returns the root token of the cluster, if set
func (dc *ExecDevCluster) GetRootToken() string {
	return dc.rootToken
}
