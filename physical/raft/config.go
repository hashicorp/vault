// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package raft

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/go-uuid"
	goversion "github.com/hashicorp/go-version"
	autopilot "github.com/hashicorp/raft-autopilot"
	bolt "go.etcd.io/bbolt"
)

type RaftBackendConfig struct {
	Path                        string
	NodeId                      string
	ApplyDelay                  time.Duration
	RaftWal                     bool
	RaftLogVerifierEnabled      bool
	RaftLogVerificationInterval time.Duration
	SnapshotDelay               time.Duration
	MaxEntrySize                uint64
	MaxBatchEntries             int
	MaxBatchSize                int
	AutopilotReconcileInterval  time.Duration
	AutopilotUpdateInterval     time.Duration
	RetryJoin                   string

	// Enterprise only
	RaftNonVoter                       bool
	MaxMountAndNamespaceTableEntrySize uint64
	AutopilotUpgradeVersion            string
	AutopilotRedundancyZone            string
}

func parseRaftBackendConfig(conf map[string]string, logger log.Logger) (*RaftBackendConfig, error) {
	c := &RaftBackendConfig{}

	c.Path = conf["path"]
	envPath := os.Getenv(EnvVaultRaftPath)
	if envPath != "" {
		c.Path = envPath
	}

	if c.Path == "" {
		return nil, fmt.Errorf("'path' must be set")
	}

	c.NodeId = conf["node_id"]
	envNodeId := os.Getenv(EnvVaultRaftNodeID)
	if envNodeId != "" {
		c.NodeId = envNodeId
	}

	if c.NodeId == "" {
		localIDRaw, err := os.ReadFile(filepath.Join(c.Path, "node-id"))
		if err == nil && len(localIDRaw) > 0 {
			c.NodeId = string(localIDRaw)
		}
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
	}

	if c.NodeId == "" {
		id, err := uuid.GenerateUUID()
		if err != nil {
			return nil, err
		}

		if err = os.WriteFile(filepath.Join(c.Path, "node-id"), []byte(id), 0o600); err != nil {
			return nil, err
		}

		c.NodeId = id
	}

	if delayRaw, ok := conf["apply_delay"]; ok {
		delay, err := parseutil.ParseDurationSecond(delayRaw)
		if err != nil {
			return nil, fmt.Errorf("apply_delay does not parse as a duration: %w", err)
		}

		c.ApplyDelay = delay
	}

	if walRaw, ok := conf["raft_wal"]; ok {
		useRaftWal, err := strconv.ParseBool(walRaw)
		if err != nil {
			return nil, fmt.Errorf("raft_wal does not parse as a boolean: %w", err)
		}

		c.RaftWal = useRaftWal
	}

	if rlveRaw, ok := conf["raft_log_verifier_enabled"]; ok {
		rlve, err := strconv.ParseBool(rlveRaw)
		if err != nil {
			return nil, fmt.Errorf("raft_log_verifier_enabled does not parse as a boolean: %w", err)
		}
		c.RaftLogVerifierEnabled = rlve

		c.RaftLogVerificationInterval = defaultRaftLogVerificationInterval
		if rlviRaw, ok := conf["raft_log_verification_interval"]; ok {
			rlvi, err := parseutil.ParseDurationSecond(rlviRaw)
			if err != nil {
				return nil, fmt.Errorf("raft_log_verification_interval does not parse as a duration: %w", err)
			}

			// Make sure our interval is capped to a reasonable value, so e.g. people don't use 0s or 1s
			if rlvi >= minimumRaftLogVerificationInterval {
				c.RaftLogVerificationInterval = rlvi
			} else {
				logger.Warn("raft_log_verification_interval is less than the minimum allowed, using default instead",
					"given", rlveRaw,
					"minimum", minimumRaftLogVerificationInterval,
					"default", defaultRaftLogVerificationInterval)
			}
		}
	}

	if delayRaw, ok := conf["snapshot_delay"]; ok {
		delay, err := parseutil.ParseDurationSecond(delayRaw)
		if err != nil {
			return nil, fmt.Errorf("snapshot_delay does not parse as a duration: %w", err)
		}
		c.SnapshotDelay = delay
	}

	c.MaxEntrySize = defaultMaxEntrySize
	if maxEntrySizeCfg := conf["max_entry_size"]; len(maxEntrySizeCfg) != 0 {
		i, err := strconv.Atoi(maxEntrySizeCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to parse 'max_entry_size': %w", err)
		}

		c.MaxEntrySize = uint64(i)
	}

	c.MaxMountAndNamespaceTableEntrySize = c.MaxEntrySize
	if maxMNTEntrySize := conf["max_mount_and_namespace_table_entry_size"]; len(maxMNTEntrySize) != 0 {
		i, err := strconv.Atoi(maxMNTEntrySize)
		if err != nil {
			return nil, fmt.Errorf("failed to parse 'max_mount_and_namespace_table_entry_size': %w", err)
		}
		if i < 1024 {
			return nil, fmt.Errorf("'max_mount_and_namespace_table_entry_size' must be at least 1024 bytes")
		}
		if i > 10_485_760 {
			return nil, fmt.Errorf("'max_mount_and_namespace_table_entry_size' must be at most 10,485,760 bytes (10MiB)")
		}

		c.MaxMountAndNamespaceTableEntrySize = uint64(i)
		emitEntWarning(logger, "max_mount_and_namespace_table_entry_size")
	}

	c.MaxBatchEntries, c.MaxBatchSize = batchLimitsFromEnv(logger)

	if interval := conf["autopilot_reconcile_interval"]; interval != "" {
		interval, err := parseutil.ParseDurationSecond(interval)
		if err != nil {
			return nil, fmt.Errorf("autopilot_reconcile_interval does not parse as a duration: %w", err)
		}
		c.AutopilotReconcileInterval = interval
	}

	if interval := conf["autopilot_update_interval"]; interval != "" {
		interval, err := parseutil.ParseDurationSecond(interval)
		if err != nil {
			return nil, fmt.Errorf("autopilot_update_interval does not parse as a duration: %w", err)
		}
		c.AutopilotUpdateInterval = interval
	}

	effectiveReconcileInterval := autopilot.DefaultReconcileInterval
	effectiveUpdateInterval := autopilot.DefaultUpdateInterval

	if c.AutopilotReconcileInterval != 0 {
		effectiveReconcileInterval = c.AutopilotReconcileInterval
	}
	if c.AutopilotUpdateInterval != 0 {
		effectiveUpdateInterval = c.AutopilotUpdateInterval
	}

	if effectiveReconcileInterval < effectiveUpdateInterval {
		return nil, fmt.Errorf("autopilot_reconcile_interval (%v) should be larger than autopilot_update_interval (%v)", effectiveReconcileInterval, effectiveUpdateInterval)
	}

	if uv, ok := conf["autopilot_upgrade_version"]; ok && uv != "" {
		_, err := goversion.NewVersion(uv)
		if err != nil {
			return nil, fmt.Errorf("autopilot_upgrade_version does not parse as a semantic version: %w", err)
		}

		c.AutopilotUpgradeVersion = uv
	}
	if c.AutopilotUpgradeVersion != "" {
		emitEntWarning(logger, "autopilot_upgrade_version")
	}

	// Note: historically we've never parsed retry_join here because we have to
	// wait until we have leader TLS info before we can work out the final retry
	// join parameters. That happens in JoinConfig. So right now nothing uses
	// c.RetryJoin because it's not available at that point. But I think it's less
	// surprising that if the field is present in the returned struct, that it
	// should actually be populated and makes tests of this function less confusing
	// too.
	c.RetryJoin = conf["retry_join"]

	c.RaftNonVoter = false
	if v := os.Getenv(EnvVaultRaftNonVoter); v != "" {
		// Consistent with handling of other raft boolean env vars
		// VAULT_RAFT_AUTOPILOT_DISABLE and VAULT_RAFT_FREELIST_SYNC
		c.RaftNonVoter = true
	} else if v, ok := conf[raftNonVoterConfigKey]; ok {
		nonVoter, err := strconv.ParseBool(v)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s config value %q as a boolean: %w", raftNonVoterConfigKey, v, err)
		}

		c.RaftNonVoter = nonVoter
	}

	if c.RaftNonVoter && c.RetryJoin == "" {
		return nil, fmt.Errorf("setting %s to true is only valid if at least one retry_join stanza is specified", raftNonVoterConfigKey)
	}
	if c.RaftNonVoter {
		emitEntWarning(logger, raftNonVoterConfigKey)
	}

	c.AutopilotRedundancyZone = conf["autopilot_redundancy_zone"]
	if c.AutopilotRedundancyZone != "" {
		emitEntWarning(logger, "autopilot_redundancy_zone")
	}

	return c, nil
}

// boltOptions returns a bolt.Options struct, suitable for passing to
// bolt.Open(), pre-configured with all of our preferred defaults.
func boltOptions(path string) *bolt.Options {
	o := &bolt.Options{
		Timeout:        1 * time.Second,
		FreelistType:   bolt.FreelistMapType,
		NoFreelistSync: true,
		MmapFlags:      getMmapFlags(path),
	}

	if os.Getenv("VAULT_RAFT_FREELIST_TYPE") == "array" {
		o.FreelistType = bolt.FreelistArrayType
	}

	if os.Getenv("VAULT_RAFT_FREELIST_SYNC") != "" {
		o.NoFreelistSync = false
	}

	// By default, we want to set InitialMmapSize to 100GB, but only on 64bit platforms.
	// Otherwise, we set it to whatever the value of VAULT_RAFT_INITIAL_MMAP_SIZE
	// is, assuming it can be parsed as an int. Bolt itself sets this to 0 by default,
	// so if users are wanting to turn this off, they can also set it to 0. Setting it
	// to a negative value is the same as not setting it at all.
	if os.Getenv("VAULT_RAFT_INITIAL_MMAP_SIZE") == "" {
		o.InitialMmapSize = initialMmapSize
	} else {
		imms, err := strconv.Atoi(os.Getenv("VAULT_RAFT_INITIAL_MMAP_SIZE"))

		// If there's an error here, it means they passed something that's not convertible to
		// a number. Rather than fail startup, just ignore it.
		if err == nil && imms > 0 {
			o.InitialMmapSize = imms
		}
	}

	return o
}
