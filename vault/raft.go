package vault

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"time"

	proto "github.com/golang/protobuf/proto"
	"github.com/hashicorp/errwrap"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
	physicalstd "github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/vault/seal"
	"github.com/mitchellh/mapstructure"
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

func (c *Core) JoinRaftCluster(ctx context.Context, leaderAddr string, retry bool) (bool, error) {
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

	join := func() error {
		// Unwrap the token
		clientConf := api.DefaultConfig()
		clientConf.Address = leaderAddr
		/*	if apiTLSConfig != nil {
			err := clientConf.ConfigureTLS(apiTLSConfig)
			if err != nil {
				return nil, errwrap.Wrapf("error configuring api client {{err}}", err)
			}
		} */
		client, err := api.NewClient(clientConf)
		if err != nil {
			return errwrap.Wrapf("failed to create api client: {{err}}", err)
		}

		secret, err := client.Logical().Write("sys/storage/raft/bootstrap/challenge", map[string]interface{}{
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
			c.raftLeaderAddr = leaderAddr
			c.raftLeaderBarrierConfig = &sealConfig
			c.seal.SetBarrierConfig(ctx, &sealConfig)
			return nil
		}

		if err := c.joinRaftSendAnswer(ctx, leaderAddr, eBlob, c.seal.GetAccess()); err != nil {
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

func (c *Core) joinRaftSendAnswer(ctx context.Context, leaderAddr string, challenge *physical.EncryptedBlobInfo, sealAccess seal.Access) error {
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

	clientConf := api.DefaultConfig()
	clientConf.Address = leaderAddr
	client, err := api.NewClient(clientConf)
	if err != nil {
		return errwrap.Wrapf("failed to create api client: {{err}}", err)
	}

	answerReq := client.NewRequest("PUT", "/v1/sys/storage/raft/bootstrap/answer")
	if err := answerReq.SetJSONBody(map[string]interface{}{
		"answer":       base64.StdEncoding.EncodeToString(plaintext),
		"cluster_addr": c.clusterListenerAddrs[0].String(),
		"server_id":    raftStorage.NodeID(),
	}); err != nil {
		return err
	}

	answerRespJson, err := client.RawRequestWithContext(ctx, answerReq)
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
