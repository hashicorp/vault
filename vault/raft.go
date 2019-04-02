package vault

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/physical/raft"
)

func (c *Core) joinRaftCluster(ctx context.Context, leaderAddr string, retry bool) (bool, error) {
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
		if err := json.Unmarshal(challengeRaw, eBlob); err != nil {
			return errwrap.Wrapf("error decoding challenge: {{err}}", err)
		}

		sealAccess := c.seal.GetAccess()
		pt, err := sealAccess.Decrypt(ctx, eBlob)
		if err != nil {
			return errwrap.Wrapf("error decrypting challenge: {{err}}", err)
		}

		secret, err = client.Logical().Write("sys/storage/raft/bootstrap/answer", map[string]interface{}{
			"answer": pt,
		})
		if err != nil {
			return errwrap.Wrapf("error sending answer: {{err}}", err)
		}
		if secret == nil {
			return errors.New("no response when sending answer")
		}

		tlsCertRaw, ok := secret.Data["tls_cert"]
		if !ok {
			return errors.New("error during raft bootstrap call, no tls cert given")
		}

		tlsKeyRaw, ok := secret.Data["tls_key"]
		if !ok {
			return errors.New("error during raft bootstrap call, no tls key given")
		}
		tlsCARaw, ok := secret.Data["tls_ca_cert"]
		if !ok {
			return errors.New("error during raft bootstrap call, no tls CA cert given")
		}
		peersRaw, ok := secret.Data["peers"]
		if !ok {
			return errors.New("error during raft bootstrap call, no peers given")
		}

		c.underlyingPhysical.(*raft.RaftBackend).Bootstrap(ctx, []raft.Peer{})

		err = c.startClusterListener(ctx)
		if err != nil {
			return errwrap.Wrapf("error starting cluster: {{err}}", err)
		}

		c.underlyingPhysical.(*raft.RaftBackend).SetupCluster(ctx, &physical.NetworkConfig{
			Cert: tlsCertRaw.([]byte),
			Addr: c.clusterListener.Addr(),
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
