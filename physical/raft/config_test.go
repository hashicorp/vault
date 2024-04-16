// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package raft

import (
	"bytes"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/constants"
	"github.com/stretchr/testify/require"
)

func ceOnlyWarnings(warns ...string) []string {
	if !constants.IsEnterprise {
		return warns
	}
	return nil
}

func TestRaft_ParseConfig(t *testing.T) {
	// Note some of these can be parallel tests but since we need to setEnv in
	// some we can't make them all parallel so it's don inside the loop. We assume
	// if a case doesn't set anything on the Env it's safe to run in parallel.
	tcs := []struct {
		name         string
		conf         map[string]string
		env          map[string]string
		wantMutation func(cfg *RaftBackendConfig)
		wantErr      string
		wantWarns    []string
	}{
		// RAFT WAL --------------------------------------------------------------
		{
			name: "WAL backend junk",
			conf: map[string]string{
				"raft_wal": "notabooleanlol",
			},
			wantErr: "does not parse as a boolean",
		},
		{
			name: "WAL verifier junk",
			conf: map[string]string{
				"raft_wal":                  "true",
				"raft_log_verifier_enabled": "notabooleanlol",
			},
			wantErr: "does not parse as a boolean",
		},
		{
			name: "WAL verifier interval, zero",
			conf: map[string]string{
				"raft_log_verifier_enabled":      "true",
				"raft_log_verification_interval": "0s",
			},
			wantMutation: func(cfg *RaftBackendConfig) {
				cfg.RaftLogVerifierEnabled = true
				cfg.RaftLogVerificationInterval = defaultRaftLogVerificationInterval
			},
			wantWarns: []string{"raft_log_verification_interval is less than the minimum allowed"},
		},
		{
			name: "WAL verifier interval, one",
			conf: map[string]string{
				"raft_log_verifier_enabled":      "true",
				"raft_log_verification_interval": "0s",
			},
			wantMutation: func(cfg *RaftBackendConfig) {
				cfg.RaftLogVerifierEnabled = true

				// Below min so should get default
				cfg.RaftLogVerificationInterval = defaultRaftLogVerificationInterval
			},
			wantWarns: []string{"raft_log_verification_interval is less than the minimum allowed"},
		},
		{
			name: "WAL verifier interval, nothing",
			conf: map[string]string{
				"raft_log_verifier_enabled":      "true",
				"raft_log_verification_interval": "",
			},
			wantMutation: func(cfg *RaftBackendConfig) {
				cfg.RaftLogVerifierEnabled = true
				cfg.RaftLogVerificationInterval = defaultRaftLogVerificationInterval
			},
			wantWarns: []string{"raft_log_verification_interval is less than the minimum allowed"},
		},
		{
			name: "WAL verifier interval, valid",
			conf: map[string]string{
				"raft_log_verifier_enabled":      "true",
				"raft_log_verification_interval": "75s",
			},
			wantMutation: func(cfg *RaftBackendConfig) {
				cfg.RaftLogVerifierEnabled = true
				cfg.RaftLogVerificationInterval = 75 * time.Second
			},
		},
		{
			name: "WAL verifier interval, junk",
			conf: map[string]string{
				"raft_log_verifier_enabled":      "true",
				"raft_log_verification_interval": "notaduration",
			},
			wantErr: "does not parse as a duration",
		},

		// AUTOPILOT Upgrades ----------------------------------------------------
		{
			name: "Autopilot upgrade version, junk",
			conf: map[string]string{
				"autopilot_upgrade_version": "hahano",
			},
			wantErr: "does not parse",
		},

		// AUTOPILOT Redundancy Zone ---------------------------------------------
		{
			name: "Autopilot redundancy zone, ok",
			conf: map[string]string{
				"autopilot_redundancy_zone": "us-east-1a",
			},
			wantMutation: func(cfg *RaftBackendConfig) {
				cfg.AutopilotRedundancyZone = "us-east-1a"
			},
			wantWarns: ceOnlyWarnings("configuration for a Vault Enterprise feature has been ignored: field=autopilot_redundancy_zone"),
		},

		// Non-voter config ------------------------------------------------------
		{
			name: "non-voter, no retry-join, valid false",
			conf: map[string]string{
				raftNonVoterConfigKey: "false",
			},
			wantMutation: func(cfg *RaftBackendConfig) {
				// Should be default
			},
		},
		{
			name: "non-voter, retry-join, valid false",
			conf: map[string]string{
				"retry_join":          "not-empty",
				raftNonVoterConfigKey: "false",
			},
			wantMutation: func(cfg *RaftBackendConfig) {
				cfg.RetryJoin = "not-empty"
			},
		},
		{
			name: "non-voter, no retry-join, valid true",
			conf: map[string]string{
				raftNonVoterConfigKey: "true",
			},
			wantErr: "only valid if at least one retry_join stanza is specified",
		},
		{
			name: "non-voter, retry-join, valid true",
			conf: map[string]string{
				"retry_join":          "not-empty",
				raftNonVoterConfigKey: "true",
			},
			wantMutation: func(cfg *RaftBackendConfig) {
				cfg.RetryJoin = "not-empty"
				cfg.RaftNonVoter = true
			},
			wantWarns: ceOnlyWarnings("configuration for a Vault Enterprise feature has been ignored: field=retry_join_as_non_voter"),
		},
		{
			name: "non-voter, no retry-join, invalid empty",
			conf: map[string]string{
				raftNonVoterConfigKey: "",
			},
			wantErr: "failed to parse retry_join_as_non_voter",
		},
		{
			name: "non-voter, retry-join, invalid empty",
			conf: map[string]string{
				"retry_join":          "not-empty",
				raftNonVoterConfigKey: "",
			},
			wantErr: "failed to parse retry_join_as_non_voter",
		},
		{
			name: "non-voter, no retry-join, invalid truthy",
			conf: map[string]string{
				raftNonVoterConfigKey: "no",
			},
			wantErr: "failed to parse retry_join_as_non_voter",
		},
		{
			name: "non-voter, retry-join, invalid truthy",
			conf: map[string]string{
				"retry_join":          "not-empty",
				raftNonVoterConfigKey: "no",
			},
			wantErr: "failed to parse retry_join_as_non_voter",
		},
		{
			name: "non-voter, no retry-join, invalid",
			conf: map[string]string{
				raftNonVoterConfigKey: "totallywrong",
			},
			wantErr: "failed to parse retry_join_as_non_voter",
		},
		{
			name: "non-voter, retry-join, invalid",
			conf: map[string]string{
				"retry_join":          "not-empty",
				raftNonVoterConfigKey: "totallywrong",
			},
			wantErr: "failed to parse retry_join_as_non_voter",
		},
		{
			// Note for historical reasons we treat any non-empty value as true in ENV
			// vars.
			name: "non-voter, no retry-join, valid env false",
			env: map[string]string{
				EnvVaultRaftNonVoter: "false",
			},
			wantErr: "only valid if at least one retry_join stanza is specified",
		},
		{
			name: "non-voter, retry-join, valid env false",
			env: map[string]string{
				EnvVaultRaftNonVoter: "false",
			},
			conf: map[string]string{
				"retry_join": "not-empty",
			},
			wantMutation: func(cfg *RaftBackendConfig) {
				cfg.RetryJoin = "not-empty"
				cfg.RaftNonVoter = true // Any non-empty value is true
			},
			wantWarns: ceOnlyWarnings("configuration for a Vault Enterprise feature has been ignored: field=retry_join_as_non_voter"),
		},
		{
			name: "non-voter, no retry-join, valid env true",
			env: map[string]string{
				EnvVaultRaftNonVoter: "true",
			},
			wantErr: "only valid if at least one retry_join stanza is specified",
		},
		{
			name: "non-voter, retry-join, valid env true",
			env: map[string]string{
				EnvVaultRaftNonVoter: "true",
			},
			conf: map[string]string{
				"retry_join": "not-empty",
			},
			wantMutation: func(cfg *RaftBackendConfig) {
				cfg.RetryJoin = "not-empty"
				cfg.RaftNonVoter = true
			},
			wantWarns: ceOnlyWarnings("configuration for a Vault Enterprise feature has been ignored: field=retry_join_as_non_voter"),
		},
		{
			name: "non-voter, no retry-join, valid env not-boolean",
			env: map[string]string{
				EnvVaultRaftNonVoter: "anything",
			},
			wantErr: "only valid if at least one retry_join stanza is specified",
		},
		{
			name: "non-voter, retry-join, valid env not-boolean",
			env: map[string]string{
				EnvVaultRaftNonVoter: "anything",
			},
			conf: map[string]string{
				"retry_join": "not-empty",
			},
			wantMutation: func(cfg *RaftBackendConfig) {
				cfg.RetryJoin = "not-empty"
				cfg.RaftNonVoter = true
			},
			wantWarns: ceOnlyWarnings("configuration for a Vault Enterprise feature has been ignored: field=retry_join_as_non_voter"),
		},
		{
			name: "non-voter, no retry-join, valid env empty",
			env: map[string]string{
				EnvVaultRaftNonVoter: "",
			},
			wantMutation: func(cfg *RaftBackendConfig) {
				// Default
			},
		},
		{
			name: "non-voter, retry-join, valid env empty",
			env: map[string]string{
				EnvVaultRaftNonVoter: "",
			},
			conf: map[string]string{
				"retry_join": "not-empty",
			},
			wantMutation: func(cfg *RaftBackendConfig) {
				cfg.RetryJoin = "not-empty"
			},
		},
		{
			name: "non-voter, no retry-join, both set env preferred",
			env: map[string]string{
				EnvVaultRaftNonVoter: "true",
			},
			conf: map[string]string{
				raftNonVoterConfigKey: "false",
			},
			wantErr: "only valid if at least one retry_join stanza is specified",
		},
		{
			name: "non-voter, retry-join, both set env preferred",
			env: map[string]string{
				EnvVaultRaftNonVoter: "true",
			},
			conf: map[string]string{
				"retry_join":          "not-empty",
				raftNonVoterConfigKey: "false",
			},
			wantMutation: func(cfg *RaftBackendConfig) {
				cfg.RetryJoin = "not-empty"
				cfg.RaftNonVoter = true // Env should win
			},
			wantWarns: ceOnlyWarnings("configuration for a Vault Enterprise feature has been ignored: field=retry_join_as_non_voter"),
		},

		// Entry Size Limits -----------------------------------------------------
		{
			name: "entry size, happy path",
			conf: map[string]string{
				"max_entry_size": "123456",
				"max_mount_and_namespace_table_entry_size": "654321",
			},
			wantMutation: func(cfg *RaftBackendConfig) {
				cfg.MaxEntrySize = 123456
				cfg.MaxMountAndNamespaceTableEntrySize = 654321
			},
			wantWarns: ceOnlyWarnings("configuration for a Vault Enterprise feature has been ignored: field=max_mount_and_namespace_table_entry_size"),
		},
		{
			name: "entry size, junk entry size",
			conf: map[string]string{
				"max_entry_size": "sadfsaf",
				"max_mount_and_namespace_table_entry_size": "654321",
			},
			wantErr: "failed to parse 'max_entry_size'",
		},
		{
			name: "entry size, junk mount entry size",
			conf: map[string]string{
				"max_entry_size": "123456",
				"max_mount_and_namespace_table_entry_size": "1MiB",
			},
			wantErr: "failed to parse 'max_mount_and_namespace_table_entry_size'",
		},
		{
			name: "entry size, way too small mount entry size",
			conf: map[string]string{
				"max_mount_and_namespace_table_entry_size": "1",
			},
			wantErr: "'max_mount_and_namespace_table_entry_size' must be at least 1024 bytes",
		},
		{
			name: "entry size, way too big mount entry size",
			conf: map[string]string{
				"max_mount_and_namespace_table_entry_size": "20000000",
			},
			wantErr: "'max_mount_and_namespace_table_entry_size' must be at most 10,485,760 bytes (10MiB)",
		},
	}

	// Set a nodeid and path to remove noise from all the test cases.
	baseConf := map[string]string{
		"node_id": "abc123",
		"path":    "/dummy/path",
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if len(tc.env) == 0 {
				// Only run in parallel if there are no env vars to set.
				t.Parallel()
			}

			var logs bytes.Buffer
			logger := hclog.New(&hclog.LoggerOptions{
				Level:  hclog.Warn,
				Output: &logs,
			})

			if tc.conf == nil {
				tc.conf = make(map[string]string)
			}

			for k, v := range baseConf {
				if _, ok := tc.conf[k]; !ok {
					tc.conf[k] = v
				}
			}

			// Make a default-valued config to compare against later. Note we do this
			// before setting ENV as that would could change behavior!
			wantCfg, err := parseRaftBackendConfig(baseConf, hclog.NewNullLogger())
			require.NoError(t, err)

			for k, v := range tc.env {
				t.Setenv(k, v)
			}

			cfg, err := parseRaftBackendConfig(tc.conf, logger)

			if tc.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.wantErr)
				return
			}

			tc.wantMutation(wantCfg)

			require.Equal(t, wantCfg, cfg)
			allLogs := logs.String()
			for _, warn := range tc.wantWarns {
				require.Contains(t, allLogs, warn)
			}
			if len(tc.wantWarns) == 0 {
				require.NotContains(t, allLogs, "[WARN]", "no warnings expected")
			}
		})
	}
}
