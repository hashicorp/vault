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
	host = fmt.Sprintf("fw-%s", host)
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

func (c *Core) JoinRaftCluster(ctx context.Context, leaderAddr string, retry bool) (bool, error) {
	if len(leaderAddr) == 0 {
		return false, errors.New("No leader address provided")
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
			return errwrap.Wrapf("error during api client creation: {{err}}", err)
		}

		secret, err := client.Logical().Write("sys/storage/raft/bootstrap/challenge", map[string]interface{}{
			"cluster_addr": c.clusterAddr,
		})
		if err != nil {
			return errwrap.Wrapf("error during bootstrap init call: {{err}}", err)
		}
		if secret == nil {
			return errors.New("could not retrieve bootstrap package")
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
		peerIDRaw, ok := secret.Data["peer_id"]
		if !ok {
			return errors.New("error during raft bootstrap call, no peer id given")
		}

		sealAccess := c.seal.GetAccess()
		pt, err := sealAccess.Decrypt(ctx, eBlob)
		if err != nil {
			return errwrap.Wrapf("error decrypting challenge: {{err}}", err)
		}

		answerReq := client.NewRequest("PUT", "/v1/sys/storage/raft/bootstrap/answer")
		if err := answerReq.SetJSONBody(map[string]interface{}{
			"answer":       base64.StdEncoding.EncodeToString(pt),
			"cluster_addr": c.clusterAddr,
			"peer_id":      peerIDRaw.(string),
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

		var answerResp answerResp
		if err := jsonutil.DecodeJSONFromReader(answerRespJson.Body, &answerResp); err != nil {
			return err
		}

		tlsCert, err := base64.StdEncoding.DecodeString(answerResp.TLSCertRaw)
		if err != nil {
			return errwrap.Wrapf("error decoding tls cert: {{err}}", err)
		}

		c.underlyingPhysical.(*raft.RaftBackend).Bootstrap(ctx, answerResp.Peers)

		err = c.startClusterListener(ctx)
		if err != nil {
			return errwrap.Wrapf("error starting cluster: {{err}}", err)
		}

		c.underlyingPhysical.(*raft.RaftBackend).SetupCluster(ctx, &physicalstd.NetworkConfig{
			Addr:      c.clusterListenerAddrs[0],
			Cert:      tlsCert,
			KeyParams: answerResp.TLSKey,
		}, c.clusterListener)

		return nil
	}

	switch retry {
	case true:
		go func() {
			for {
				// TODO add a way to shut this down
				if err := join(); err != nil {
					c.logger.Error("failed to join raft cluster", "error", err)
					time.Sleep(time.Second * 2)
				}
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

type answerResp struct {
	Peers      []raft.Peer                `json:"peers"`
	TLSCertRaw string                     `json:"tls_cert"`
	TLSKey     *certutil.ClusterKeyParams `json:"tls_key"`
}
