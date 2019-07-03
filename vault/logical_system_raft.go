package vault

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"strings"

	proto "github.com/golang/protobuf/proto"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/vault/seal"
)

// raftStoragePaths returns paths for use when raft is the storage mechanism.
func (b *SystemBackend) raftStoragePaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "storage/raft/bootstrap/answer",

			Fields: map[string]*framework.FieldSchema{
				"server_id": {
					Type: framework.TypeString,
				},
				"answer": {
					Type: framework.TypeString,
				},
				"cluster_addr": {
					Type: framework.TypeString,
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleRaftBootstrapAnswerWrite(),
					Summary:  "Accepts an answer from the peer to be joined to the fact cluster.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysRaftHelp["raft-bootstrap-answer"][0]),
			HelpDescription: strings.TrimSpace(sysRaftHelp["raft-bootstrap-answer"][1]),
		},
		{
			Pattern: "storage/raft/bootstrap/challenge",

			Fields: map[string]*framework.FieldSchema{
				"server_id": {
					Type: framework.TypeString,
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleRaftBootstrapChallengeWrite(),
					Summary:  "Creates a challenge for the new peer to be joined to the raft cluster.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysRaftHelp["raft-bootstrap-challenge"][0]),
			HelpDescription: strings.TrimSpace(sysRaftHelp["raft-bootstrap-challenge"][1]),
		},
		{
			Pattern: "storage/raft/remove-peer",

			Fields: map[string]*framework.FieldSchema{
				"server_id": {
					Type: framework.TypeString,
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleRaftRemovePeerUpdate(),
					Summary:  "Remove a peer from the raft cluster.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysRaftHelp["raft-remove-peer"][0]),
			HelpDescription: strings.TrimSpace(sysRaftHelp["raft-remove-peer"][1]),
		},
		{
			Pattern: "storage/raft/configuration",

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleRaftConfigurationGet(),
					Summary:  "Returns the configuration of the raft cluster.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysRaftHelp["raft-remove-peer"][0]),
			HelpDescription: strings.TrimSpace(sysRaftHelp["raft-remove-peer"][1]),
		},
		{
			Pattern: "storage/raft/snapshot",
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleStorageRaftSnapshotRead(),
					Summary:  "Retruns a snapshot of the current state of vault.",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleStorageRaftSnapshotWrite(false),
					Summary:  "Installs the provided snapshot, returning the cluster to the state defined in it.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysRaftHelp["raft-remove-peer"][0]),
			HelpDescription: strings.TrimSpace(sysRaftHelp["raft-remove-peer"][1]),
		},
		{
			Pattern: "storage/raft/snapshot-force",
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleStorageRaftSnapshotWrite(true),
					Summary:  "Installs the provided snapshot, returning the cluster to the state defined in it. This bypasses checks ensuring the current Autounseal or Shamir keys are consistent with the snapshot data.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysRaftHelp["raft-remove-peer"][0]),
			HelpDescription: strings.TrimSpace(sysRaftHelp["raft-remove-peer"][1]),
		},
	}
}

func (b *SystemBackend) handleRaftConfigurationGet() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

		raftStorage, ok := b.Core.underlyingPhysical.(*raft.RaftBackend)
		if !ok {
			return logical.ErrorResponse("raft storage is not in use"), logical.ErrInvalidRequest
		}

		config, err := raftStorage.GetConfiguration(ctx)
		if err != nil {
			return nil, err
		}

		return &logical.Response{
			Data: map[string]interface{}{
				"config": config,
			},
		}, nil
	}
}

func (b *SystemBackend) handleRaftRemovePeerUpdate() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		serverID := d.Get("server_id").(string)
		if len(serverID) == 0 {
			return logical.ErrorResponse("no server id provided"), logical.ErrInvalidRequest
		}

		raftStorage, ok := b.Core.underlyingPhysical.(*raft.RaftBackend)
		if !ok {
			return logical.ErrorResponse("raft storage is not in use"), logical.ErrInvalidRequest
		}

		if err := raftStorage.RemovePeer(ctx, serverID); err != nil {
			return nil, err
		}
		if b.Core.raftFollowerStates != nil {
			b.Core.raftFollowerStates.delete(serverID)
		}

		return nil, nil
	}
}

func (b *SystemBackend) handleRaftBootstrapChallengeWrite() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		_, ok := b.Core.underlyingPhysical.(*raft.RaftBackend)
		if !ok {
			return logical.ErrorResponse("raft storage is not in use"), logical.ErrInvalidRequest
		}

		serverID := d.Get("server_id").(string)
		if len(serverID) == 0 {
			return logical.ErrorResponse("no server id provided"), logical.ErrInvalidRequest
		}

		uuid, err := uuid.GenerateRandomBytes(16)
		if err != nil {
			return nil, err
		}

		sealAccess := b.Core.seal.GetAccess()
		eBlob, err := sealAccess.Encrypt(ctx, uuid)
		if err != nil {
			return nil, err
		}
		protoBlob, err := proto.Marshal(eBlob)
		if err != nil {
			return nil, err
		}

		b.Core.pendingRaftPeers[serverID] = uuid
		sealConfig, err := b.Core.seal.BarrierConfig(ctx)
		if err != nil {
			return nil, err
		}

		return &logical.Response{
			Data: map[string]interface{}{
				"challenge":   base64.StdEncoding.EncodeToString(protoBlob),
				"seal_config": sealConfig,
			},
		}, nil
	}
}

func (b *SystemBackend) handleRaftBootstrapAnswerWrite() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		raftStorage, ok := b.Core.underlyingPhysical.(*raft.RaftBackend)
		if !ok {
			return logical.ErrorResponse("raft storage is not in use"), logical.ErrInvalidRequest
		}

		serverID := d.Get("server_id").(string)
		if len(serverID) == 0 {
			return logical.ErrorResponse("no server_id provided"), logical.ErrInvalidRequest
		}
		answerRaw := d.Get("answer").(string)
		if len(answerRaw) == 0 {
			return logical.ErrorResponse("no answer provided"), logical.ErrInvalidRequest
		}
		clusterAddr := d.Get("cluster_addr").(string)
		if len(clusterAddr) == 0 {
			return logical.ErrorResponse("no cluster_addr provided"), logical.ErrInvalidRequest
		}

		answer, err := base64.StdEncoding.DecodeString(answerRaw)
		if err != nil {
			return logical.ErrorResponse("could not base64 decode answer"), logical.ErrInvalidRequest
		}

		expectedAnswer, ok := b.Core.pendingRaftPeers[serverID]
		if !ok {
			return logical.ErrorResponse("no expected answer for the server id provided"), logical.ErrInvalidRequest
		}

		delete(b.Core.pendingRaftPeers, serverID)

		if subtle.ConstantTimeCompare(answer, expectedAnswer) == 0 {
			return logical.ErrorResponse("invalid answer given"), logical.ErrInvalidRequest
		}

		tlsKeyringEntry, err := b.Core.barrier.Get(ctx, raftTLSStoragePath)
		if err != nil {
			return nil, err
		}
		if tlsKeyringEntry == nil {
			return nil, errors.New("could not find raft TLS configuration")
		}
		var keyring raft.RaftTLSKeyring
		if err := tlsKeyringEntry.DecodeJSON(&keyring); err != nil {
			return nil, errors.New("could not decode raft TLS configuration")
		}

		if err := raftStorage.AddPeer(ctx, serverID, clusterAddr); err != nil {
			return nil, err
		}
		if b.Core.raftFollowerStates != nil {
			b.Core.raftFollowerStates.update(serverID, 0)
		}

		peers, err := raftStorage.Peers(ctx)
		if err != nil {
			return nil, err
		}

		return &logical.Response{
			Data: map[string]interface{}{
				"peers":       peers,
				"tls_keyring": &keyring,
			},
		}, nil
	}
}

func (b *SystemBackend) handleStorageRaftSnapshotRead() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		raftStorage, ok := b.Core.underlyingPhysical.(*raft.RaftBackend)
		if !ok {
			return logical.ErrorResponse("raft storage is not in use"), logical.ErrInvalidRequest
		}
		if req.ResponseWriter == nil {
			return nil, errors.New("no writer for request")
		}

		err := raftStorage.Snapshot(req.ResponseWriter, b.Core.seal.GetAccess())
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}

func (b *SystemBackend) handleStorageRaftSnapshotWrite(force bool) framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		raftStorage, ok := b.Core.underlyingPhysical.(*raft.RaftBackend)
		if !ok {
			return logical.ErrorResponse("raft storage is not in use"), logical.ErrInvalidRequest
		}
		if req.RequestReader == nil {
			return nil, errors.New("no reader for request")
		}

		access := b.Core.seal.GetAccess()
		if force {
			access = nil
		}

		// We want to buffer the http request reader into a temp file here so we
		// don't have to hold the full snapshot in memory. We also want to do
		// the restore in two parts so we can restore the snapshot while the
		// stateLock is write locked.
		snapFile, cleanup, metadata, err := raftStorage.WriteSnapshotToTemp(req.RequestReader, access)
		switch {
		case err == nil:
		case strings.Contains(err.Error(), "failed to open the sealed hashes"):
			switch b.Core.seal.BarrierType() {
			case seal.Shamir:
				return logical.ErrorResponse("could not verify hash file, possibly the snapshot is using a different set of unseal keys; use the snapshot-force API to bypass this check"), logical.ErrInvalidRequest
			default:
				return logical.ErrorResponse("could not verify hash file, possibly the snapshot is using a different autoseal key; use the snapshot-force API to bypass this check"), logical.ErrInvalidRequest
			}
		case err != nil:
			b.Core.logger.Error("raft snapshot restore: failed to write snapshot", "error", err)
			return nil, err
		}

		// We want to do this in a go routine so we can upgrade the lock and
		// allow the client to disconnect.
		go func() (retErr error) {
			// Cleanup the temp file
			defer cleanup()

			// Grab statelock
			if stopped := grabLockOrStop(b.Core.stateLock.Lock, b.Core.stateLock.Unlock, b.Core.standbyStopCh); stopped {
				b.Core.logger.Error("not applying snapshot; shutting down")
				return
			}
			defer b.Core.stateLock.Unlock()

			// If we failed to restore the snapshot we should seal this node as
			// it's in an unknown state
			defer func() {
				if retErr != nil {
					if err := b.Core.sealInternalWithOptions(false, false, true); err != nil {
						b.Core.logger.Error("failed to seal node", "error", err)
					}
				}
			}()

			ctx, ctxCancel := context.WithCancel(namespace.RootContext(nil))

			// We are calling the callback function synchronously here while we
			// have the lock. So set it to nil and restore the callback when we
			// finish.
			raftStorage.SetRestoreCallback(nil)
			defer raftStorage.SetRestoreCallback(b.Core.raftSnapshotRestoreCallback(true, true))

			// Do a preSeal to clear vault's in-memory caches and shut down any
			// systems that might be holding the encryption access.
			b.Core.logger.Info("shutting down prior to restoring snapshot")
			if err := b.Core.preSeal(); err != nil {
				b.Core.logger.Error("raft snapshot restore failed preSeal", "error", err)
				return err
			}

			b.Core.logger.Info("applying snapshot")
			if err := raftStorage.RestoreSnapshot(ctx, metadata, snapFile); err != nil {
				b.Core.logger.Error("error while restoring raft snapshot", "error", err)
				return err
			}

			// Run invalidation logic synchronously here
			callback := b.Core.raftSnapshotRestoreCallback(false, false)
			if err := callback(ctx); err != nil {
				return err
			}

			{
				// If the snapshot was taken while another node was leader we
				// need to reset the leader information to this node.
				if err := b.Core.underlyingPhysical.Put(ctx, &physical.Entry{
					Key:   CoreLockPath,
					Value: []byte(b.Core.leaderUUID),
				}); err != nil {
					b.Core.logger.Error("cluster setup failed", "error", err)
					return err
				}
				// re-advertise our cluster information
				if err := b.Core.advertiseLeader(ctx, b.Core.leaderUUID, nil); err != nil {
					b.Core.logger.Error("cluster setup failed", "error", err)
					return err
				}
			}
			if err := b.Core.postUnseal(ctx, ctxCancel, standardUnsealStrategy{}); err != nil {
				b.Core.logger.Error("raft snapshot restore failed postUnseal", "error", err)
				return err
			}

			return nil

		}()

		return nil, nil
	}
}

var sysRaftHelp = map[string][2]string{
	"raft-bootstrap-challenge": {
		"Creates a challenge for the new peer to be joined to the raft cluster.",
		"",
	},
	"raft-bootstrap-answer": {
		"Accepts an answer from the peer to be joined to the fact cluster.",
		"",
	},
	"raft-remove-peer": {
		"Removes a peer from the raft cluster.",
		"",
	},
}
