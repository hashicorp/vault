package vault

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"net/http"
	"strings"
	"time"

	proto "github.com/golang/protobuf/proto"
	"github.com/hashicorp/errwrap"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/namespace"
	physicalstd "github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/vault/seal"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/net/http2"
)

var raftTLSStoragePath = "core/raft/tls"

type raftTLSConfig struct {
	Cert      []byte                     `json:"cluster_cert,omitempty"`
	KeyParams *certutil.ClusterKeyParams `json:"cluster_key_params,omitempty"`
}

func generateRaftTLS() (*raftTLSConfig, error) {
	key, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return nil, err
	}

	host, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}
	host = fmt.Sprintf("raft-%s", host)
	template := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: host,
		},
		DNSNames: []string{host},
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement | x509.KeyUsageCertSign,
		SerialNumber: big.NewInt(mathrand.Int63()),
		NotBefore:    time.Now().Add(-30 * time.Second),
		// 30 years of single-active uptime ought to be enough for anybody
		NotAfter:              time.Now().Add(262980 * time.Hour),
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, key.Public(), key)
	if err != nil {
		return nil, errwrap.Wrapf("unable to generate local cluster certificate: {{err}}", err)
	}

	return &raftTLSConfig{
		Cert: certBytes,
		KeyParams: &certutil.ClusterKeyParams{
			Type: certutil.PrivateKeyTypeP521,
			X:    key.PublicKey.X,
			Y:    key.PublicKey.Y,
			D:    key.D,
		},
	}, nil
}

func (c *Core) startRaftStorage(ctx context.Context) error {
	if raftStorage, ok := c.underlyingPhysical.(*raft.RaftBackend); ok {
		if raftStorage.Initialized() {
			return nil
		}

		raftTLSEntry, err := c.barrier.Get(ctx, raftTLSStoragePath)
		if err != nil {
			return err
		}
		if raftTLSEntry == nil {
			return errors.New("could not find raft TLS configuration")
		}

		raftTLS := new(raftTLSConfig)
		if err := raftTLSEntry.DecodeJSON(raftTLS); err != nil {
			return err
		}

		raftStorage.SetRestoreCallback(c.raftSnapshotRestoreCallback(true))
		if err := raftStorage.SetupCluster(ctx, &physicalstd.NetworkConfig{
			Addr:      c.clusterListenerAddrs[0],
			Cert:      raftTLS.Cert,
			KeyParams: raftTLS.KeyParams,
		}, c.clusterListener); err != nil {
			return err
		}
	}

	return nil
}

// handleSnapshotRestore is for the raft backend to hook back into core after a
// snapshot is restored so we can clear the necessary caches and handle changing
// keyrings or master keys
func (c *Core) raftSnapshotRestoreCallback(grabLock bool) func() error {
	return func() error {
		c.logger.Info("running post snapshot restore invalidations")

		if grabLock {
			// Grab statelock
			if stopped := grabLockOrStop(c.stateLock.Lock, c.stateLock.Unlock, c.standbyStopCh); stopped {
				c.logger.Error("did not apply snapshot; vault is shutting down")
				return errors.New("did not apply snapshot; vault is shutting down")
			}
			defer c.stateLock.Unlock()
		}
		ctx, ctxCancel := context.WithCancel(namespace.RootContext(nil))

		// Purge the cache so we make sure we are operating on fresh data
		c.physicalCache.Purge(ctx)

		// Reload the keyring in case it changed. If this fails it's likely
		// we've changed master keys.
		err := c.performKeyUpgrades(ctx)
		if err != nil {
			// The snapshot contained a master key or keyring we couldn't
			// recover
			switch c.seal.BarrierType() {
			case seal.Shamir:
				// If we are a shamir seal we can't do anything. Just
				// seal all nodes.

				// Seal ourselves
				c.logger.Info("failed to perform key upgrades, sealing", "error", err)
				c.sealInternalWithOptions(false, false)
				return err
			default:
				// If we are using an auto-unseal we can try to use the seal to
				// unseal again. If the auto-unseal mechanism has changed then
				// there isn't anything we can do but seal.
				c.logger.Info("failed to perform key upgrades, reloading using auto seal")
				keys, err := c.seal.GetStoredKeys(ctx)
				if err != nil {
					c.logger.Error("raft snapshot restore failed to get stored keys", "error", err)
					c.sealInternalWithOptions(false, false)
					return err
				}
				if err := c.barrier.Seal(); err != nil {
					c.logger.Error("raft snapshot restore failed to seal barrier", "error", err)
					c.sealInternalWithOptions(false, false)
					return err
				}
				if err := c.barrier.Unseal(ctx, keys[0]); err != nil {
					c.logger.Error("raft snapshot restore failed to unseal barrier", "error", err)
					c.sealInternalWithOptions(false, false)
					return err
				}
				c.logger.Info("done reloading master key using auto seal")
			}
		}

		if !c.standby {
			// Go through a preSeal, postUnseal cycle to clear the rest of
			// vault's in-memory caches
			c.logger.Info("restarting system after raft snapshot restore")
			if err := c.preSeal(); err != nil {
				c.logger.Error("raft snapshot restore failed preSeal", "error", err)
				return err
			}
			{
				// If the snapshot was taken while another node was leader we
				// need to reset the leader information to this node.
				if err := c.underlyingPhysical.Put(ctx, &physical.Entry{
					Key:   CoreLockPath,
					Value: []byte(c.leaderUUID),
				}); err != nil {
					c.logger.Error("cluster setup failed", "error", err)
					return err
				}
				// re-advertise our cluster information
				if err := c.advertiseLeader(ctx, c.leaderUUID, nil); err != nil {
					c.logger.Error("cluster setup failed", "error", err)
					return err
				}
			}
			if err := c.postUnseal(ctx, ctxCancel, standardUnsealStrategy{}); err != nil {
				c.logger.Error("raft snapshot restore failed postUnseal", "error", err)
				return err
			}
		}

		return nil
	}
}

func (c *Core) JoinRaftCluster(ctx context.Context, leaderAddr string, tlsConfig *tls.Config, retry bool) (bool, error) {
	if len(leaderAddr) == 0 {
		return false, errors.New("No leader address provided")
	}

	raftStorage, ok := c.underlyingPhysical.(*raft.RaftBackend)
	if !ok {
		return false, errors.New("raft storage not configured")
	}

	if raftStorage.Initialized() {
		return false, errors.New("raft is alreay initialized")
	}

	init, err := c.Initialized(ctx)
	if err != nil {
		return false, errwrap.Wrapf("failed to check if core is initialized: {{err}}", err)
	}
	if init {
		return false, errwrap.Wrapf("join can't be invoked on an initialized cluster: {{err}}", ErrAlreadyInit)
	}

	transport := cleanhttp.DefaultPooledTransport()
	if tlsConfig != nil {
		transport.TLSClientConfig = tlsConfig.Clone()
		if err := http2.ConfigureTransport(transport); err != nil {
			return false, errwrap.Wrapf("failed to configure TLS: {{err}}", err)
		}
	}
	client := &http.Client{
		Transport: transport,
	}
	config := api.DefaultConfig()
	if config.Error != nil {
		return false, errwrap.Wrapf("failed to create api client: {{err}}", config.Error)
	}
	config.Address = leaderAddr
	config.HttpClient = client
	config.MaxRetries = 0
	apiClient, err := api.NewClient(config)
	if err != nil {
		return false, errwrap.Wrapf("failed to create api client: {{err}}", err)
	}

	join := func() error {
		// Unwrap the token
		secret, err := apiClient.Logical().Write("sys/storage/raft/bootstrap/challenge", map[string]interface{}{
			"server_id": raftStorage.NodeID(),
		})
		if err != nil {
			return errwrap.Wrapf("error during bootstrap init call: {{err}}", err)
		}
		if secret == nil {
			return errors.New("could not retrieve bootstrap package")
		}

		var sealConfig SealConfig
		err = mapstructure.Decode(secret.Data["seal_config"], &sealConfig)
		if err != nil {
			return err
		}

		if sealConfig.Type != c.seal.BarrierType() {
			return fmt.Errorf("mismatching seal types between leader (%s) and follower (%s)", sealConfig.Type, c.seal.BarrierType())
		}

		challengeB64, ok := secret.Data["challenge"]
		if !ok {
			return errors.New("error during raft bootstrap call, no challenge given")
		}
		challengeRaw, err := base64.StdEncoding.DecodeString(challengeB64.(string))
		if err != nil {
			return errwrap.Wrapf("error decoding challenge: {{err}}", err)
		}

		eBlob := &physical.EncryptedBlobInfo{}
		if err := proto.Unmarshal(challengeRaw, eBlob); err != nil {
			return errwrap.Wrapf("error decoding challenge: {{err}}", err)
		}

		if c.seal.BarrierType() == seal.Shamir {
			c.raftUnseal = true
			c.raftChallenge = eBlob
			c.raftLeaderClient = apiClient
			c.raftLeaderBarrierConfig = &sealConfig
			c.seal.SetBarrierConfig(ctx, &sealConfig)
			return nil
		}

		if err := c.joinRaftSendAnswer(ctx, apiClient, eBlob, c.seal.GetAccess()); err != nil {
			return errwrap.Wrapf("failed to send answer to leader node: {{err}}", err)
		}

		return nil
	}

	switch retry {
	case true:
		go func() {
			for {
				// TODO add a way to shut this down
				err := join()
				if err == nil {
					return
				}
				c.logger.Error("failed to join raft cluster", "error", err)
				time.Sleep(time.Second * 2)
			}
		}()

		// Backgrounded so return false
		return false, nil
	default:
		if err := join(); err != nil {
			c.logger.Error("failed to join raft cluster", "error", err)
			return false, errwrap.Wrapf("failed to join raft cluster: {{err}}", err)
		}
	}

	return true, nil
}

// This is used in tests to override the cluster address
var UpdateClusterAddrForTests bool

func (c *Core) joinRaftSendAnswer(ctx context.Context, leaderClient *api.Client, challenge *physical.EncryptedBlobInfo, sealAccess seal.Access) error {
	if challenge == nil {
		return errors.New("raft challenge is nil")
	}

	raftStorage, ok := c.underlyingPhysical.(*raft.RaftBackend)
	if !ok {
		return errors.New("raft storage not in use")
	}

	if raftStorage.Initialized() {
		return errors.New("raft is already initialized")
	}

	plaintext, err := sealAccess.Decrypt(ctx, challenge)
	if err != nil {
		return errwrap.Wrapf("error decrypting challenge: {{err}}", err)
	}

	clusterAddr := c.clusterAddr
	if UpdateClusterAddrForTests && strings.HasSuffix(clusterAddr, ":0") {
		// We are testing and have an address provider, so just create a random
		// addr, it will be overwritten later.
		var err error
		clusterAddr, err = uuid.GenerateUUID()
		if err != nil {
			return err
		}
	}

	answerReq := leaderClient.NewRequest("PUT", "/v1/sys/storage/raft/bootstrap/answer")
	if err := answerReq.SetJSONBody(map[string]interface{}{
		"answer":       base64.StdEncoding.EncodeToString(plaintext),
		"cluster_addr": clusterAddr,
		"server_id":    raftStorage.NodeID(),
	}); err != nil {
		return err
	}

	answerRespJson, err := leaderClient.RawRequestWithContext(ctx, answerReq)
	if answerRespJson != nil {
		defer answerRespJson.Body.Close()
	}
	if err != nil {
		return err
	}

	var answerResp answerRespData
	if err := jsonutil.DecodeJSONFromReader(answerRespJson.Body, &answerResp); err != nil {
		return err
	}

	tlsCert, err := base64.StdEncoding.DecodeString(answerResp.Data.TLSCertRaw)
	if err != nil {
		return errwrap.Wrapf("error decoding tls cert: {{err}}", err)
	}

	raftStorage.Bootstrap(ctx, answerResp.Data.Peers)

	err = c.startClusterListener(ctx)
	if err != nil {
		return errwrap.Wrapf("error starting cluster: {{err}}", err)
	}

	raftStorage.SetRestoreCallback(c.raftSnapshotRestoreCallback(true))
	err = raftStorage.SetupCluster(ctx, &physicalstd.NetworkConfig{
		Addr:      c.clusterListenerAddrs[0],
		Cert:      tlsCert,
		KeyParams: answerResp.Data.TLSKey,
	}, c.clusterListener)
	if err != nil {
		return errwrap.Wrapf("failed to setup raft cluster: {{err}}", err)
	}

	return nil
}

func (c *Core) isRaftUnseal() bool {
	return c.raftUnseal
}

type answerRespData struct {
	Data answerResp `json:"data"`
}

type answerResp struct {
	Peers      []raft.Peer                `json:"peers"`
	TLSCertRaw string                     `json:"tls_cert"`
	TLSKey     *certutil.ClusterKeyParams `json:"tls_key"`
}
