package vault

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"strings"
	"time"

	proto "github.com/golang/protobuf/proto"
	wrapping "github.com/hashicorp/go-kms-wrapping"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"
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
				"non_voter": {
					Type: framework.TypeBool,
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
				"dr_operation_token": {
					Type:        framework.TypeString,
					Description: "DR operation token used to authorize this request (if a DR secondary node).",
				},
				"server_id": {
					Type: framework.TypeString,
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.verifyDROperationTokenOnSecondary(b.handleRaftRemovePeerUpdate(), false),
					Summary:  "Remove a peer from the raft cluster.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysRaftHelp["raft-remove-peer"][0]),
			HelpDescription: strings.TrimSpace(sysRaftHelp["raft-remove-peer"][1]),
		},
		{
			Pattern: "storage/raft/configuration",

			Fields: map[string]*framework.FieldSchema{
				"dr_operation_token": {
					Type:        framework.TypeString,
					Description: "DR operation token used to authorize this request (if a DR secondary node).",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleRaftConfigurationGet(),
					Summary:  "Returns the configuration of the raft cluster.",
				},
				// Reading configuration on a DR secondary cluster is an update
				// operation to allow consuming the DR operation token for
				// authenticating the request.
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.verifyDROperationToken(b.handleRaftConfigurationGet(), false),
					Summary:  "Returns the configuration of the raft cluster in a DR secondary cluster.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysRaftHelp["raft-configuration"][0]),
			HelpDescription: strings.TrimSpace(sysRaftHelp["raft-configuration"][1]),
		},
		{
			Pattern: "storage/raft/snapshot",
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleStorageRaftSnapshotRead(),
					Summary:  "Returns a snapshot of the current state of vault.",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleStorageRaftSnapshotWrite(false),
					Summary:  "Installs the provided snapshot, returning the cluster to the state defined in it.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysRaftHelp["raft-snapshot"][0]),
			HelpDescription: strings.TrimSpace(sysRaftHelp["raft-snapshot"][1]),
		},
		{
			Pattern: "storage/raft/snapshot-force",
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleStorageRaftSnapshotWrite(true),
					Summary:  "Installs the provided snapshot, returning the cluster to the state defined in it. This bypasses checks ensuring the current Autounseal or Shamir keys are consistent with the snapshot data.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysRaftHelp["raft-snapshot-force"][0]),
			HelpDescription: strings.TrimSpace(sysRaftHelp["raft-snapshot-force"][1]),
		},
		{
			Pattern: "storage/raft/autopilot/configuration",

			Fields: map[string]*framework.FieldSchema{
				"cleanup_dead_servers": {
					Type: framework.TypeBool,
				},
				"last_contact_threshold": {
					Type: framework.TypeDurationSecond,
				},
				"max_trailing_logs": {
					Type: framework.TypeInt,
				},
				"min_quorum": {
					Type: framework.TypeInt,
				},
				"server_stabilization_time": {
					Type: framework.TypeDurationSecond,
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleStorageRaftAutopilotConfigRead(),
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleStorageRaftAutopilotConfigUpdate(),
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysRaftHelp["raft-snapshot-force"][0]),
			HelpDescription: strings.TrimSpace(sysRaftHelp["raft-snapshot-force"][1]),
		},
	}
}

func (b *SystemBackend) handleRaftConfigurationGet() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		raftBackend := b.Core.getRaftBackend()
		if raftBackend == nil {
			return logical.ErrorResponse("raft storage is not in use"), logical.ErrInvalidRequest
		}

		config, err := raftBackend.GetConfiguration(ctx)
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

		raftBackend := b.Core.getRaftBackend()
		if raftBackend == nil {
			return logical.ErrorResponse("raft storage is not in use"), logical.ErrInvalidRequest
		}

		if err := raftBackend.RemovePeer(ctx, serverID); err != nil {
			return nil, err
		}
		if b.Core.raftFollowerStates != nil {
			b.Core.raftFollowerStates.Delete(serverID)
		}

		return nil, nil
	}
}

func (b *SystemBackend) handleRaftBootstrapChallengeWrite() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		serverID := d.Get("server_id").(string)
		if len(serverID) == 0 {
			return logical.ErrorResponse("no server id provided"), logical.ErrInvalidRequest
		}

		var answer []byte
		answerRaw, ok := b.Core.pendingRaftPeers.Load(serverID)
		if !ok {
			var err error
			answer, err = uuid.GenerateRandomBytes(16)
			if err != nil {
				return nil, err
			}
			b.Core.pendingRaftPeers.Store(serverID, answer)
		} else {
			answer = answerRaw.([]byte)
		}

		sealAccess := b.Core.seal.GetAccess()

		eBlob, err := sealAccess.Encrypt(ctx, answer, nil)
		if err != nil {
			return nil, err
		}
		protoBlob, err := proto.Marshal(eBlob)
		if err != nil {
			return nil, err
		}

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
		raftBackend := b.Core.getRaftBackend()
		if raftBackend == nil {
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

		nonVoter := d.Get("non_voter").(bool)

		answer, err := base64.StdEncoding.DecodeString(answerRaw)
		if err != nil {
			return logical.ErrorResponse("could not base64 decode answer"), logical.ErrInvalidRequest
		}

		expectedAnswerRaw, ok := b.Core.pendingRaftPeers.Load(serverID)
		if !ok {
			return logical.ErrorResponse("no expected answer for the server id provided"), logical.ErrInvalidRequest
		}

		b.Core.pendingRaftPeers.Delete(serverID)

		if subtle.ConstantTimeCompare(answer, expectedAnswerRaw.([]byte)) == 0 {
			return logical.ErrorResponse("invalid answer given"), logical.ErrInvalidRequest
		}

		tlsKeyringEntry, err := b.Core.barrier.Get(ctx, raftTLSStoragePath)
		if err != nil {
			return nil, err
		}
		if tlsKeyringEntry == nil {
			return nil, errors.New("could not find raft TLS configuration")
		}
		var keyring raft.TLSKeyring
		if err := tlsKeyringEntry.DecodeJSON(&keyring); err != nil {
			return nil, errors.New("could not decode raft TLS configuration")
		}

		switch nonVoter {
		case true:
			err = raftBackend.AddNonVotingPeer(ctx, serverID, clusterAddr)
		default:
			err = raftBackend.AddPeer(ctx, serverID, clusterAddr)
		}
		if err != nil {
			return nil, err
		}

		if b.Core.raftFollowerStates != nil {
			b.Core.raftFollowerStates.Update(serverID, 0, 0)
		}

		peers, err := raftBackend.Peers(ctx)
		if err != nil {
			return nil, err
		}

		b.logger.Info("follower node answered the raft bootstrap challenge", "follower_server_id", serverID)

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

		err := raftStorage.SnapshotHTTP(req.ResponseWriter, b.Core.seal.GetAccess())
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}

func (b *SystemBackend) handleStorageRaftAutopilotConfigRead() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		raftStorage, ok := b.Core.underlyingPhysical.(*raft.RaftBackend)
		if !ok {
			return logical.ErrorResponse("raft storage is not in use"), logical.ErrInvalidRequest
		}
		config := raftStorage.AutopilotConfig()
		return &logical.Response{
			Data: map[string]interface{}{
				"cleanup_dead_servers":      config.CleanupDeadServers,
				"last_contact_threshold":    config.LastContactThreshold.String(),
				"max_trailing_logs":         config.MaxTrailingLogs,
				"min_quorum":                config.MinQuorum,
				"server_stabilization_time": config.ServerStabilizationTime.String(),
			},
		}, nil
	}
}

func (b *SystemBackend) handleStorageRaftAutopilotConfigUpdate() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		raftStorage, ok := b.Core.underlyingPhysical.(*raft.RaftBackend)
		if !ok {
			return logical.ErrorResponse("raft storage is not in use"), logical.ErrInvalidRequest
		}

		// Get the config present in the backend
		config := raftStorage.AutopilotConfig()

		// Mutate the clone of the actual config to avoid setting values in
		// failure cases.
		configClone := config.Clone()

		persist := false
		cleanupDeadServers, ok := d.GetOk("cleanup_dead_servers")
		if ok {
			configClone.CleanupDeadServers = cleanupDeadServers.(bool)
			persist = true
		}
		lastContactThreshold, ok := d.GetOk("last_contact_threshold")
		if ok {
			configClone.LastContactThreshold = time.Duration(lastContactThreshold.(int)) * time.Second
			persist = true
		}
		maxTrailingLogs, ok := d.GetOk("max_trailing_logs")
		if ok {
			configClone.MaxTrailingLogs = uint64(maxTrailingLogs.(int))
			persist = true
		}
		minQuorum, ok := d.GetOk("min_quorum")
		if ok {
			configClone.MinQuorum = uint(minQuorum.(int))
			persist = true
		}
		serverStabilizationTime, ok := d.GetOk("server_stabilization_time")
		if ok {
			configClone.ServerStabilizationTime = time.Duration(serverStabilizationTime.(int)) * time.Second
			persist = true
		}

		if persist {
			entry, err := logical.StorageEntryJSON("core/raft/autopilot/configuration", configClone)
			if err != nil {
				return nil, err
			}
			if err := b.Core.barrier.Put(ctx, entry); err != nil {
				return nil, err
			}
		}

		raftStorage.SetAutopilotConfig(configClone)

		return nil, nil
	}
}

func (b *SystemBackend) handleStorageRaftSnapshotWrite(force bool) framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		raftStorage, ok := b.Core.underlyingPhysical.(*raft.RaftBackend)
		if !ok {
			return logical.ErrorResponse("raft storage is not in use"), logical.ErrInvalidRequest
		}
		if req.HTTPRequest == nil || req.HTTPRequest.Body == nil {
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
		snapFile, cleanup, metadata, err := raftStorage.WriteSnapshotToTemp(req.HTTPRequest.Body, access)
		switch {
		case err == nil:
		case strings.Contains(err.Error(), "failed to open the sealed hashes"):
			switch b.Core.seal.BarrierType() {
			case wrapping.Shamir:
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
			if stopped := grabLockOrStop(b.Core.stateLock.Lock, b.Core.stateLock.Unlock, b.Core.standbyStopCh.Load().(chan struct{})); stopped {
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
	"raft-configuration": {
		"Returns the raft cluster configuration.",
		`On a DR secondary cluster, instead of a GET, this must be a POST or
		PUT, and furthermore a DR operation token must be provided.`,
	},
	"raft-remove-peer": {
		"Removes a peer from the raft cluster.",
		"",
	},
	"raft-snapshot": {
		"Restores and saves snapshots from the raft cluster.",
		"",
	},
	"raft-snapshot-force": {
		"Force restore a raft cluster snapshot",
		"",
	},
}
