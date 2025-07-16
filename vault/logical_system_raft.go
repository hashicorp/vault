// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"
	snapshot "github.com/hashicorp/raft-snapshot"
	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/vault/seal"
	"github.com/mitchellh/mapstructure"
)

// raftStoragePaths returns paths for use when raft is the storage mechanism.
func (b *SystemBackend) raftStoragePaths() []*framework.Path {
	makeSealer := func(logger hclog.Logger, use string) func() snapshot.Sealer {
		return func() snapshot.Sealer {
			return NewSealAccessSealer(b.Core.seal.GetAccess(), logger, use)
		}
	}

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
				"upgrade_version": {
					Type: framework.TypeString,
				},
				"sdk_version": {
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
					Callback: b.handleRaftBootstrapChallengeWrite(makeSealer(nil, "bootstrap_challenge_write")),
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
					Callback: b.handleStorageRaftSnapshotRead(makeSealer(nil, "snapshot_read")),
					Summary:  "Returns a snapshot of the current state of vault.",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleStorageRaftSnapshotWrite(false, makeSealer(b.logger, "snapshot_write")),
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
					Callback: b.handleStorageRaftSnapshotWrite(true, makeSealer(b.logger, "snapshot_write")),
					Summary:  "Installs the provided snapshot, returning the cluster to the state defined in it. This bypasses checks ensuring the current Autounseal or Shamir keys are consistent with the snapshot data.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysRaftHelp["raft-snapshot-force"][0]),
			HelpDescription: strings.TrimSpace(sysRaftHelp["raft-snapshot-force"][1]),
		},
		{
			Pattern: "storage/raft/autopilot/state",
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback:                  b.verifyDROperationTokenOnSecondary(b.handleStorageRaftAutopilotState(), false),
					Summary:                   "Returns the state of the raft cluster under integrated storage as seen by autopilot.",
					ForwardPerformanceStandby: true,
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysRaftHelp["raft-autopilot-state"][0]),
			HelpDescription: strings.TrimSpace(sysRaftHelp["raft-autopilot-state"][1]),
		},
		{
			Pattern: "storage/raft/autopilot/configuration",
			Fields: map[string]*framework.FieldSchema{
				"cleanup_dead_servers": {
					Type:        framework.TypeBool,
					Description: "Controls whether to remove dead servers from the Raft peer list periodically or when a new server joins.",
				},
				"last_contact_threshold": {
					Type:        framework.TypeDurationSecond,
					Description: "Limit on the amount of time a server can go without leader contact before being considered unhealthy.",
				},
				"dead_server_last_contact_threshold": {
					Type:        framework.TypeDurationSecond,
					Description: "Limit on the amount of time a server can go without leader contact before being considered failed. This takes effect only when cleanup_dead_servers is set.",
				},
				"max_trailing_logs": {
					Type:        framework.TypeInt,
					Description: "Amount of entries in the Raft Log that a server can be behind before being considered unhealthy.",
				},
				"min_quorum": {
					Type:        framework.TypeInt,
					Description: "Minimum number of servers allowed in a cluster before autopilot can prune dead servers. This should at least be 3.",
				},
				"server_stabilization_time": {
					Type:        framework.TypeDurationSecond,
					Description: "Minimum amount of time a server must be in a stable, healthy state before it can be added to the cluster.",
				},
				"disable_upgrade_migration": {
					Type:        framework.TypeBool,
					Description: "Whether or not to perform automated version upgrades.",
				},
				"dr_operation_token": {
					Type:        framework.TypeString,
					Description: "DR operation token used to authorize this request (if a DR secondary node).",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.verifyDROperationTokenOnSecondary(b.handleStorageRaftAutopilotConfigRead(), false),
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.verifyDROperationTokenOnSecondary(b.handleStorageRaftAutopilotConfigUpdate(), false),
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysRaftHelp["raft-autopilot-configuration"][0]),
			HelpDescription: strings.TrimSpace(sysRaftHelp["raft-autopilot-configuration"][1]),
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

		b.Core.raftFollowerStates.Delete(serverID)
		_, err := raftBackend.RemovedServerCleanup(ctx, serverID)
		if err != nil {
			// log the error but don't return it - we might get an error if we can't find the node in the cache, which
			// is not an error condition in this instance.
			b.logger.Info("attempted to remove node from perf standby cache but it failed, which might be fine", "server ID", serverID, "error", err)
			return nil, nil
		}

		return nil, nil
	}
}

const answerSize = 16

var answerMaxEncodedSize = base64.StdEncoding.EncodedLen(answerSize)

func (b *SystemBackend) handleRaftBootstrapChallengeWrite(makeSealer func() snapshot.Sealer) framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		serverID := d.Get("server_id").(string)
		if len(serverID) == 0 {
			return logical.ErrorResponse("no server id provided"), logical.ErrInvalidRequest
		}

		var answer []byte
		b.Core.pendingRaftPeersLock.RLock()
		challenge, ok := b.Core.pendingRaftPeers.Get(serverID)
		b.Core.pendingRaftPeersLock.RUnlock()
		if !ok {
			if !b.raftChallengeLimiter.Allow() {
				return logical.RespondWithStatusCode(logical.ErrorResponse("too many raft challenges in flight"), req, http.StatusTooManyRequests)
			}

			b.Core.pendingRaftPeersLock.Lock()
			defer b.Core.pendingRaftPeersLock.Unlock()

			challenge, ok = b.Core.pendingRaftPeers.Get(serverID)
			if !ok {

				var err error
				answer, err = uuid.GenerateRandomBytes(answerSize)
				if err != nil {
					return nil, err
				}

				sealer := makeSealer()
				if sealer == nil {
					return nil, errors.New("core has no seal access to write raft bootstrap challenge")
				}
				protoBlob, err := sealer.Seal(ctx, answer)
				if err != nil {
					return nil, err
				}

				challenge = &raftBootstrapChallenge{
					serverID:  serverID,
					answer:    answer,
					challenge: protoBlob,
				}
				b.Core.pendingRaftPeers.Add(serverID, challenge)
			}
		}

		sealConfig, err := b.Core.seal.BarrierConfig(ctx)
		if err != nil {
			return nil, err
		}

		return &logical.Response{
			Data: map[string]interface{}{
				"challenge":   base64.StdEncoding.EncodeToString(challenge.challenge),
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
		if len(answerRaw) > answerMaxEncodedSize {
			return logical.ErrorResponse("answer is too long"), logical.ErrInvalidRequest
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

		b.Core.pendingRaftPeersLock.Lock()
		expectedChallenge, ok := b.Core.pendingRaftPeers.Get(serverID)
		if !ok {
			b.Core.pendingRaftPeersLock.Unlock()
			return logical.ErrorResponse("no expected answer for the server id provided"), logical.ErrInvalidRequest
		}

		b.Core.pendingRaftPeers.Remove(serverID)
		b.Core.pendingRaftPeersLock.Unlock()

		if subtle.ConstantTimeCompare(answer, expectedChallenge.answer) == 0 {
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

		var desiredSuffrage string
		switch nonVoter {
		case true:
			desiredSuffrage = "non-voter"
		default:
			desiredSuffrage = "voter"
		}

		added := b.Core.raftFollowerStates.Update(&raft.EchoRequestUpdate{
			NodeID:          serverID,
			DesiredSuffrage: desiredSuffrage,
			SDKVersion:      d.Get("sdk_version").(string),
			UpgradeVersion:  d.Get("upgrade_version").(string),
		})

		switch nonVoter {
		case true:
			err = raftBackend.AddNonVotingPeer(ctx, serverID, clusterAddr)
		default:
			err = raftBackend.AddPeer(ctx, serverID, clusterAddr)
		}
		if err != nil {
			if added {
				b.Core.raftFollowerStates.Delete(serverID)
			}
			return nil, err
		}

		peers, err := raftBackend.Peers(ctx)
		if err != nil {
			return nil, err
		}

		b.logger.Info("follower node answered the raft bootstrap challenge", "follower_server_id", serverID)

		return &logical.Response{
			Data: map[string]interface{}{
				"peers":              peers,
				"tls_keyring":        &keyring,
				"autoloaded_license": b.Core.entIsLicenseAutoloaded(),
			},
		}, nil
	}
}

func (b *SystemBackend) handleStorageRaftSnapshotRead(makeSealer func() snapshot.Sealer) framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		raftStorage, ok := b.Core.underlyingPhysical.(*raft.RaftBackend)
		if !ok {
			return logical.ErrorResponse("raft storage is not in use"), logical.ErrInvalidRequest
		}
		if req.ResponseWriter == nil {
			return nil, errors.New("no writer for request")
		}

		err := raftStorage.SnapshotHTTP(req.ResponseWriter, makeSealer())
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}

func (b *SystemBackend) handleStorageRaftAutopilotState() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		raftBackend := b.Core.getRaftBackend()
		if raftBackend == nil {
			return logical.ErrorResponse("raft storage is not in use"), logical.ErrInvalidRequest
		}

		state, err := raftBackend.GetAutopilotServerState(ctx)
		if err != nil {
			return nil, err
		}

		if state == nil {
			return nil, nil
		}

		data := make(map[string]interface{})
		err = mapstructure.Decode(state, &data)
		if err != nil {
			return nil, err
		}

		return &logical.Response{
			Data: data,
		}, nil
	}
}

func (b *SystemBackend) handleStorageRaftAutopilotConfigRead() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		raftBackend := b.Core.getRaftBackend()
		if raftBackend == nil {
			return logical.ErrorResponse("raft storage is not in use"), logical.ErrInvalidRequest
		}

		config := raftBackend.AutopilotConfig()
		if config == nil {
			return nil, nil
		}

		return &logical.Response{
			Data: map[string]interface{}{
				"cleanup_dead_servers":               config.CleanupDeadServers,
				"last_contact_threshold":             config.LastContactThreshold.String(),
				"dead_server_last_contact_threshold": config.DeadServerLastContactThreshold.String(),
				"max_trailing_logs":                  config.MaxTrailingLogs,
				"min_quorum":                         config.MinQuorum,
				"server_stabilization_time":          config.ServerStabilizationTime.String(),
				"disable_upgrade_migration":          config.DisableUpgradeMigration,
			},
		}, nil
	}
}

func (b *SystemBackend) handleStorageRaftAutopilotConfigUpdate() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		raftBackend := b.Core.getRaftBackend()
		if raftBackend == nil {
			return logical.ErrorResponse("raft storage is not in use"), logical.ErrInvalidRequest
		}

		// Read autopilot configuration from storage
		config, err := b.Core.loadAutopilotConfiguration(ctx)
		if err != nil {
			b.logger.Error("failed to load autopilot config from storage when setting up cluster; continuing since autopilot falls back to default config", "error", err)
		}
		if config == nil {
			config = &raft.AutopilotConfig{}
		}

		persist := false
		cleanupDeadServers, ok := d.GetOk("cleanup_dead_servers")
		if ok {
			if cleanupDeadServers.(bool) {
				config.CleanupDeadServersValue = raft.CleanupDeadServersTrue
			} else {
				config.CleanupDeadServersValue = raft.CleanupDeadServersFalse
			}
			persist = true
		}
		lastContactThreshold, ok := d.GetOk("last_contact_threshold")
		if ok {
			config.LastContactThreshold = time.Duration(lastContactThreshold.(int)) * time.Second
			persist = true
		}
		deadServerLastContactThreshold, ok := d.GetOk("dead_server_last_contact_threshold")
		if ok {
			config.DeadServerLastContactThreshold = time.Duration(deadServerLastContactThreshold.(int)) * time.Second
			persist = true
		}
		maxTrailingLogs, ok := d.GetOk("max_trailing_logs")
		if ok {
			config.MaxTrailingLogs = uint64(maxTrailingLogs.(int))
			persist = true
		}
		minQuorum, ok := d.GetOk("min_quorum")
		if ok {
			config.MinQuorum = uint(minQuorum.(int))
			persist = true
		}
		serverStabilizationTime, ok := d.GetOk("server_stabilization_time")
		if ok {
			config.ServerStabilizationTime = time.Duration(serverStabilizationTime.(int)) * time.Second
			persist = true
		}
		disableUpgradeMigration, ok := d.GetOk("disable_upgrade_migration")
		if ok {
			if !constants.IsEnterprise {
				return logical.ErrorResponse("disable_upgrade_migration is only available in Vault Enterprise"), logical.ErrInvalidRequest
			}
			config.DisableUpgradeMigration = disableUpgradeMigration.(bool)
			persist = true
		}

		effectiveConf := raftBackend.AutopilotConfig()
		effectiveConf.Merge(config)

		if effectiveConf.CleanupDeadServers && effectiveConf.MinQuorum < 3 {
			return logical.ErrorResponse(fmt.Sprintf("min_quorum must be set when cleanup_dead_servers is set and it should at least be 3; cleanup_dead_servers: %#v, min_quorum: %#v", effectiveConf.CleanupDeadServers, effectiveConf.MinQuorum)), logical.ErrInvalidRequest
		}

		if effectiveConf.CleanupDeadServers && effectiveConf.DeadServerLastContactThreshold.Seconds() < 60 {
			return logical.ErrorResponse(fmt.Sprintf("dead_server_last_contact_threshold should not be set to less than 1m; received: %v", deadServerLastContactThreshold)), logical.ErrInvalidRequest
		}

		// Persist only the user supplied fields
		if persist {
			entry, err := logical.StorageEntryJSON(raftAutopilotConfigurationStoragePath, config)
			if err != nil {
				return nil, err
			}
			if err := b.Core.barrier.Put(ctx, entry); err != nil {
				return nil, err
			}
		}

		// Set the effectiveConfig
		raftBackend.SetAutopilotConfig(effectiveConf)

		return nil, nil
	}
}

func (b *SystemBackend) handleStorageRaftSnapshotWrite(force bool, makeSealer func() snapshot.Sealer) framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		raftStorage, ok := b.Core.underlyingPhysical.(*raft.RaftBackend)
		if !ok {
			return logical.ErrorResponse("raft storage is not in use"), logical.ErrInvalidRequest
		}
		body, ok := logical.ContextOriginalBodyValue(ctx)
		if !ok {
			return nil, errors.New("no reader for request")
		}

		var sealer snapshot.Sealer
		if !force {
			sealer = makeSealer()
		}

		// We want to buffer the http request reader into a temp file here so we
		// don't have to hold the full snapshot in memory. We also want to do
		// the restore in two parts so we can restore the snapshot while the
		// stateLock is write locked.
		snapFile, cleanup, metadata, err := raftStorage.WriteSnapshotToTemp(body, sealer)
		switch {
		case err == nil:
		case strings.Contains(err.Error(), "failed to open the sealed hashes"):
			switch b.Core.seal.BarrierSealConfigType() {
			case SealConfigTypeShamir:
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
			l := newLockGrabber(b.Core.stateLock.Lock, b.Core.stateLock.Unlock, b.Core.standbyStopCh.Load().(chan struct{}))
			go l.grab()
			if stopped := l.lockOrStop(); stopped {
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
	"raft-autopilot-state": {
		"Returns the state of the raft cluster under integrated storage as seen by autopilot.",
		"",
	},
	"raft-autopilot-configuration": {
		"Returns autopilot configuration.",
		"",
	},
}

func NewSealAccessSealer(access seal.Access, logger hclog.Logger, use string) snapshot.Sealer {
	// If we have access to the seal create a sealer object
	var s snapshot.Sealer
	if access != nil {
		s = &sealer{
			access: access,
			logger: logger,
			use:    use,
		}
	}
	return s
}

// sealer implements the snapshot.Sealer interface and is used in the snapshot
// process for encrypting/decrypting the SHASUM file in snapshot archives.
type sealer struct {
	access seal.Access
	logger hclog.Logger
	use    string
}

var _ snapshot.Sealer = (*sealer)(nil)

// Seal encrypts the data with using the seal access object.
func (s *sealer) Seal(ctx context.Context, pt []byte) ([]byte, error) {
	wrappedValue, err := SealWrapValue(ctx, s.access, true, pt, func(_ context.Context, errs map[string]error) error {
		if s.logger != nil {
			s.logger.Warn("sealed value not wrapped with all seals", "use", s.use)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return MarshalSealWrappedValue(wrappedValue)
}

// Open decrypts the data using the seal access object.
func (s *sealer) Open(ctx context.Context, ct []byte) ([]byte, error) {
	wrappedEntryValue, err := UnmarshalSealWrappedValue(ct)
	if err != nil {
		return nil, err
	}

	pt, _, err := UnsealWrapValue(ctx, s.access, "", wrappedEntryValue)
	return pt, err
}
