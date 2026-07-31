// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package consts

const (
	// ExpirationRestoreWorkerCount specifies the number of workers to use while
	// restoring leases into the expiration manager
	ExpirationRestoreWorkerCount = 64

	// NamespaceHeaderName is the header set to specify which namespace the
	// request is indented for.
	NamespaceHeaderName = "X-Vault-Namespace"

	// AuthHeaderName is the name of the header containing the token.
	AuthHeaderName = "X-Vault-Token"

	// RequestHeaderName is the name of the header used by the Agent for
	// SSRF protection.
	RequestHeaderName = "X-Vault-Request"

	// WrapTTLHeaderName is the name of the header containing a directive to
	// wrap the response
	WrapTTLHeaderName = "X-Vault-Wrap-TTL"

	// PerformanceReplicationALPN is the negotiated protocol used for
	// performance replication.
	PerformanceReplicationALPN = "replication_v1"

	// DRReplicationALPN is the negotiated protocol used for dr replication.
	DRReplicationALPN = "replication_dr_v1"

	PerfStandbyALPN = "perf_standby_v1"

	RequestForwardingALPN = "req_fw_sb-act_v1"

	RaftStorageALPN = "raft_storage_v1"

	// ReplicationResolverALPN is the negotiated protocol used for
	// resolving replicaiton addresses
	ReplicationResolverALPN = "replication_resolver_v1"

	VaultEnableFilePermissionsCheckEnv = "VAULT_ENABLE_FILE_PERMISSIONS_CHECK"

	VaultDisableUserLockout = "VAULT_DISABLE_USER_LOCKOUT"

	// VaultDisableConsumptionBilling is the environment variable that disables
	// the periodic consumption billing metrics worker. When set to a non-empty
	// value, the worker that walks the entire KV store every 10 minutes is
	// not registered during post-unseal, preventing the memory and performance
	// impact of the periodic KV store enumeration on large installations.
	VaultDisableConsumptionBilling = "VAULT_DISABLE_CONSUMPTION_BILLING"

	PerformanceReplicationPathTarget = "performance"

	DRReplicationPathTarget = "dr"

	RecoverSourcePathHeader = "X-Vault-Recover-Source-Path"
)
