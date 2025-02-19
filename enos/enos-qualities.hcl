// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

quality "consul_api_agent_host_read" {
  description = "The /v1/agent/host Consul API returns host info for each node in the cluster"
}

quality "consul_api_health_node_read" {
  description = <<-EOF
    The /v1/health/node/<node> Consul API returns health info for each node in the cluster
  EOF
}

quality "consul_api_operator_raft_config_read" {
  description = "The /v1/operator/raft/configuration Consul API returns raft info for the cluster"
}

quality "consul_autojoin_aws" {
  description = "The Consul cluster auto-joins with AWS tag discovery"
}

quality "consul_cli_validate" {
  description = "The 'consul validate' command validates the Consul configuration"
}

quality "consul_config_file" {
  description = "Consul starts when configured with a configuration file"
}

quality "consul_ha_leader_election" {
  description = "The Consul cluster elects a leader node on start up"
}

quality "consul_health_state_passing_read_nodes_minimum" {
  description = <<-EOF
    The Consul cluster meets the minimum of number of healthy nodes according to the
    /v1/health/state/passing Consul API
  EOF
}

quality "consul_operator_raft_configuration_read_voters_minimum" {
  description = <<-EOF
    The Consul cluster meets the minimum number of raft voters according to the
    /v1/operator/raft/configuration Consul API
  EOF
}

quality "consul_service_start_client" {
  description = "The Consul service starts in client mode"
}

quality "consul_service_start_server" {
  description = "The Consul service starts in server mode"
}

quality "consul_service_systemd_notified" {
  description = "The Consul binary notifies systemd when the service is active"
}

quality "consul_service_systemd_unit" {
  description = "The 'consul.service' systemd unit starts the service"
}

quality "vault_agent_auto_auth_approle" {
  description = <<-EOF
    Vault running in Agent mode utilizes the approle auth method to do auto-auth via a role and
    read secrets from a file source
  EOF
}

quality "vault_agent_log_template" {
  description = global.description.verify_agent_output
}

quality "vault_api_auth_userpass_login_write" {
  description = "The v1/auth/userpass/login/<user> Vault API creates a token for a user"
}

quality "vault_api_auth_userpass_user_write" {
  description = "The v1/auth/userpass/users/<user> Vault API associates a policy with a user"
}

quality "vault_api_identity_entity_read" {
  description = <<-EOF
    The v1/identity/entity Vault API returns an identity entity, has the correct metadata, and is
    associated with the expected entity-alias, groups, and policies
  EOF
}

quality "vault_api_identity_entity_write" {
  description = "The v1/identity/entity Vault API creates an identity entity"
}

quality "vault_api_identity_entity_alias_write" {
  description = "The v1/identity/entity-alias Vault API creates an identity entity alias"
}

quality "vault_api_identity_group_write" {
  description = "The v1/identity/group/<group> Vault API creates an identity group"
}

quality "vault_api_identity_oidc_config_read" {
  description = <<-EOF
    The v1/identity/oidc/config Vault API returns the built-in identity secrets engine configuration
  EOF
}

quality "vault_api_identity_oidc_config_write" {
  description = "The v1/identity/oidc/config Vault API configures the built-in identity secrets engine"
}

quality "vault_api_identity_oidc_introspect_write" {
  description = "The v1/identity/oidc/introspect Vault API creates introspect verifies the active state of a signed OIDC token"
}

quality "vault_api_identity_oidc_key_read" {
  description = <<-EOF
    The v1/identity/oidc/key Vault API returns the OIDC signing key and verifies the key's algorithm,
    rotation_period, and verification_ttl are correct
  EOF
}

quality "vault_api_identity_oidc_key_write" {
  description = "The v1/identity/oidc/key Vault API creates an OIDC signing key"
}

quality "vault_api_identity_oidc_key_rotate_write" {
  description = "The v1/identity/oidc/key/<name>/rotate Vault API rotates an OIDC signing key and applies a new verification TTL"
}

quality "vault_api_identity_oidc_role_read" {
  description = <<-EOF
    The v1/identity/oidc/role Vault API returns the OIDC role and verifies that the roles key and
    ttl are corect.
  EOF
}

quality "vault_api_identity_oidc_role_write" {
  description = "The v1/identity/oidc/role Vault API creates an OIDC role associated with a key and clients"
}

quality "vault_api_identity_oidc_token_read" {
  description = "The v1/identity/oidc/token Vault API creates an OIDC token associated with a role"
}

quality "vault_api_sys_auth_userpass_user_write" {
  description = "The v1/sys/auth/userpass/users/<user> Vault API associates a superuser policy with a user"
}

quality "vault_api_sys_config_read" {
  description = <<-EOF
    The v1/sys/config/sanitized Vault API returns sanitized configuration which matches our given
    configuration
  EOF
}

quality "vault_api_sys_ha_status_read" {
  description = "The v1/sys/ha-status Vault API returns the HA status of the cluster"
}

quality "vault_api_sys_health_read" {
  description = <<-EOF
    The v1/sys/health Vault API returns the correct codes depending on the replication and
    'seal-status' of the cluster
  EOF
}

quality "vault_api_sys_host_info_read" {
  description = "The v1/sys/host-info Vault API returns the host info for each node in the cluster"
}

quality "vault_api_sys_leader_read" {
  description = "The v1/sys/leader Vault API returns the cluster leader info"
}

quality "vault_api_sys_metrics_vault_core_replication_write_undo_logs_enabled" {
  description = <<-EOF
    The v1/sys/metrics Vault API returns metrics and verifies that
    'Gauges[vault.core.replication.write_undo_logs]' is enabled
  EOF
}

quality "vault_api_sys_policy_write" {
  description = "The v1/sys/policy Vault API writes a policy"
}

quality "vault_api_sys_quotas_lease_count_read_max_leases_default" {
  description = <<-EOF
    The v1/sys/quotas/lease-count/default Vault API returns the lease 'count' and 'max_leases' is
    set to 300,000
  EOF
}

quality "vault_api_sys_replication_dr_primary_enable_write" {
  description = <<-EOF
    The v1/sys/replication/dr/primary/enable Vault API enables DR replication
  EOF
}

quality "vault_api_sys_replication_dr_primary_secondary_token_write" {
  description = <<-EOF
    The v1/sys/replication/dr/primary/secondary-token Vault API configures the DR replication
    secondary token
  EOF
}

quality "vault_api_sys_replication_dr_secondary_enable_write" {
  description = <<-EOF
    The v1/sys/replication/dr/secondary/enable Vault API enables DR replication
  EOF
}

quality "vault_api_sys_replication_dr_read_connection_status_connected" {
  description = <<-EOF
    The v1/sys/replication/dr/status Vault API returns status info and the
    'connection_status' is correct for the given node
  EOF
}

quality "vault_api_sys_replication_dr_status_known_primary_cluster_addrs" {
  description = <<-EOF
    The v1/sys/replication/dr/status Vault API returns the DR replication status and
    'known_primary_cluster_address' is the expected primary cluster leader
  EOF
}

quality "vault_api_sys_replication_dr_status_read" {
  description = <<-EOF
    The v1/sys/replication/dr/status Vault API returns the DR replication status
  EOF
}

quality "vault_api_sys_replication_dr_status_read_cluster_address" {
  description = <<-EOF
    The v1/sys/replication/dr/status Vault API returns the DR replication status
    and the '{primaries,secondaries}[*].cluster_address' is correct for the given node
  EOF
}

quality "vault_api_sys_replication_dr_status_read_state_not_idle" {
  description = <<-EOF
    The v1/sys/replication/dr/status Vault API returns the DR replication status
    and the state is not idle
  EOF
}

quality "vault_api_sys_replication_performance_primary_enable_write" {
  description = <<-EOF
    The v1/sys/replication/performance/primary/enable Vault API enables performance replication
  EOF
}

quality "vault_api_sys_replication_performance_primary_secondary_token_write" {
  description = <<-EOF
    The v1/sys/replication/performance/primary/secondary-token Vault API configures the replication
    token
  EOF
}

quality "vault_api_sys_replication_performance_secondary_enable_write" {
  description = <<-EOF
    The v1/sys/replication/performance/secondary/enable Vault API enables performance replication
  EOF
}

quality "vault_api_sys_replication_performance_read_connection_status_connected" {
  description = <<-EOF
    The v1/sys/replication/performance/status Vault API returns status info and the
    'connection_status' is correct for the given node
  EOF
}

quality "vault_api_sys_replication_performance_status_known_primary_cluster_addrs" {
  description = <<-EOF
    The v1/sys/replication/performance/status Vault API returns the replication status and
    'known_primary_cluster_address' is the expected primary cluster leader
  EOF
}

quality "vault_api_sys_replication_performance_status_read" {
  description = <<-EOF
    The v1/sys/replication/performance/status Vault API returns the performance replication status
  EOF
}

quality "vault_api_sys_replication_performance_status_read_cluster_address" {
  description = <<-EOF
    The v1/sys/replication/performance/status Vault API returns the performance replication status
    and the '{primaries,secondaries}[*].cluster_address' is correct for the given node
  EOF
}

quality "vault_api_sys_replication_performance_status_read_state_not_idle" {
  description = <<-EOF
    The v1/sys/replication/performance/status Vault API returns the performance replication status
    and the state is not idle
  EOF
}

quality "vault_api_sys_replication_status_read" {
  description = <<-EOF
    The v1/sys/replication/status Vault API returns the performance replication status of the
    cluster
  EOF
}

quality "vault_api_sys_seal_status_api_read_matches_sys_health" {
  description = <<-EOF
    The v1/sys/seal-status Vault API and v1/sys/health Vault API agree on the health of each node
    and the cluster
  EOF
}

quality "vault_api_sys_sealwrap_rewrap_read_entries_processed_eq_entries_succeeded_post_rewrap" {
  description = global.description.verify_seal_rewrap_entries_processed_eq_entries_succeeded_post_rewrap
}

quality "vault_api_sys_sealwrap_rewrap_read_entries_processed_gt_zero_post_rewrap" {
  description = global.description.verify_seal_rewrap_entries_processed_is_gt_zero_post_rewrap
}

quality "vault_api_sys_sealwrap_rewrap_read_is_running_false_post_rewrap" {
  description = global.description.verify_seal_rewrap_is_running_false_post_rewrap
}

quality "vault_api_sys_sealwrap_rewrap_read_no_entries_fail_during_rewrap" {
  description = global.description.verify_seal_rewrap_no_entries_fail_during_rewrap
}

quality "vault_api_sys_step_down_steps_down" {
  description = <<-EOF
    The v1/sys/step-down Vault API forces the cluster leader to step down and intiates a new leader
    election
  EOF
}

quality "vault_api_sys_storage_raft_autopilot_configuration_read" {
  description = <<-EOF
    The /sys/storage/raft/autopilot/configuration Vault API returns the autopilot configuration of
    the cluster
  EOF
}

quality "vault_api_sys_storage_raft_autopilot_state_read" {
  description = <<-EOF
    The v1/sys/storage/raft/autopilot/state Vault API returns the raft autopilot state of the
    cluster
  EOF
}

quality "vault_api_sys_storage_raft_autopilot_upgrade_info_read_status_matches" {
  description = <<-EOF
    The v1/sys/storage/raft/autopilot/state Vault API returns the raft autopilot state and the
    'upgrade_info.status' matches our expected state
  EOF
}

quality "vault_api_sys_storage_raft_autopilot_upgrade_info_target_version_read_matches_candidate" {
  description = <<-EOF
    The v1/sys/storage/raft/autopilot/state Vault API returns the raft autopilot state and the
    'upgrade_info.target_version' matches the the candidate version
  EOF
}

quality "vault_api_sys_storage_raft_configuration_read" {
  description = <<-EOF
    The v1/sys/storage/raft/configuration Vault API returns the raft configuration of the cluster
  EOF
}

quality "vault_api_sys_storage_raft_remove_peer_write_removes_peer" {
  description = <<-EOF
    The v1/sys/storage/raft/remove-peer Vault API removes the desired node from the raft sub-system
  EOF
}

quality "vault_api_sys_version_history_keys" {
  description = <<-EOF
    The v1/sys/version-history Vault API returns the cluster version history and the 'keys' data
    includes our target version
  EOF
}

quality "vault_api_sys_version_history_key_info" {
  description = <<-EOF
    The v1/sys/version-history Vault API returns the cluster version history and the
    'key_info["$expected_version]' data is present for the expected version and the 'build_date'
    matches the expected build_date.
  EOF
}

quality "vault_artifact_bundle" {
  description = "The candidate binary packaged as a zip bundle is used for testing"
}

quality "vault_artifact_deb" {
  description = "The candidate binary packaged as a deb package is used for testing"
}

quality "vault_artifact_rpm" {
  description = "The candidate binary packaged as an rpm package is used for testing"
}

quality "vault_audit_log" {
  description = "The Vault audit sub-system is enabled with the log and writes to a log"
}

quality "vault_audit_log_secrets" {
  description = "The Vault audit sub-system does not output secret values"
}

quality "vault_audit_socket" {
  description = "The Vault audit sub-system is enabled with the socket and writes to a socket"
}

quality "vault_audit_syslog" {
  description = "The Vault audit sub-system is enabled with the syslog and writes to syslog"
}

quality "vault_auto_unseals_after_autopilot_upgrade" {
  description = "Vault auto-unseals after upgrading the cluster with autopilot"
}

quality "vault_autojoins_new_nodes_into_initialized_cluster" {
  description = "Vault sucessfully auto-joins new nodes into an existing cluster"
}

quality "vault_autojoin_aws" {
  description = "Vault auto-joins nodes using AWS tag discovery"
}

quality "vault_autopilot_upgrade_leader_election" {
  description = <<-EOF
    Vault elects a new leader after upgrading the cluster with autopilot
  EOF
}

quality "vault_cli_audit_enable" {
  description = "The 'vault audit enable' command enables audit devices"
}

quality "vault_cli_auth_enable_approle" {
  description = "The 'vault auth enable approle' command enables the approle auth method"
}

quality "vault_cli_operator_members" {
  description = "The 'vault operator members' command returns the expected list of members"
}

quality "vault_cli_operator_raft_remove_peer" {
  description = "The 'vault operator remove-peer' command removes the desired node"
}

quality "vault_cli_operator_step_down" {
  description = "The 'vault operator step-down' command forces the cluster leader to step down"
}

quality "vault_cli_policy_write" {
  description = "The 'vault policy write' command writes a policy"
}

quality "vault_cli_status_exit_code" {
  description = <<-EOF
    The 'vault status' command exits with the correct code depending on expected seal status
  EOF
}

quality "vault_cluster_upgrade_in_place" {
  description = <<-EOF
    Vault starts with existing data and configuration in-place migrates the data
  EOF
}

quality "vault_config_env_variables" {
  description = "Vault starts when configured primarily with environment variables"
}

quality "vault_config_file" {
  description = "Vault starts when configured primarily with a configuration file"
}

quality "vault_config_log_level" {
  description = "The 'log_level' config stanza modifies its log level"
}

quality "vault_config_multiseal_is_toggleable" {
  description = <<-EOF
    The Vault Cluster can be configured with a single unseal method regardless of the
    'enable_multiseal' config value
  EOF
}

quality "vault_init" {
  description = "Vault initializes the cluster with the given seal parameters"
}

quality "vault_journal_secrets" {
  description = "The Vault systemd journal does not output secret values"
}

quality "vault_license_required_ent" {
  description = "Vault Enterprise requires a license in order to start"
}

quality "vault_listener_ipv4" {
  description = "Vault operates on ipv4 TCP listeners"
}

quality "vault_listener_ipv6" {
  description = "Vault operates on ipv6 TCP listeners"
}

quality "vault_mount_auth" {
  description = "Vault mounts the auth engine"
}

quality "vault_mount_identity" {
  description = "Vault mounts the identity engine"
}

quality "vault_mount_kv" {
  description = "Vault mounts the kv engine"
}

quality "vault_multiseal_enable" {
  description = <<-EOF
    The Vault Cluster starts with 'enable_multiseal' and multiple auto-unseal methods.
  EOF
}

quality "vault_proxy_auto_auth_approle" {
  description = <<-EOF
    Vault Proxy utilizes the approle auth method to to auto auth via a roles and secrets from file.
  EOF
}

quality "vault_proxy_cli_access" {
  description = <<-EOF
    The Vault CLI accesses tokens through the Vault proxy without a VAULT_TOKEN available
  EOF
}

quality "vault_radar_index_create" {
  description = "Vault radar is able to create an index from KVv2 mounts"
}

quality "vault_radar_scan_file" {
  description = "Vault radar is able to scan a file for secrets"
}

quality "vault_raft_voters" {
  description = global.description.verify_raft_cluster_all_nodes_are_voters
}

quality "vault_replication_ce_disabled" {
  description = "Replication is not enabled for CE editions"
}

quality "vault_replication_ent_dr_available" {
  description = "DR replication is available on Enterprise"
}

quality "vault_replication_ent_pr_available" {
  description = "PR replication is available on Enterprise"
}

quality "vault_seal_awskms" {
  description = "Vault auto-unseals with the awskms seal"
}

quality "vault_seal_shamir" {
  description = <<-EOF
    Vault manually unseals with the shamir seal when given the expected number of 'key_shares'
  EOF
}

quality "vault_seal_pkcs11" {
  description = "Vault auto-unseals with the pkcs11 seal"
}

quality "vault_secrets_kv_read" {
  description = "Vault kv secrets engine data is readable"
}

quality "vault_secrets_kv_write" {
  description = "Vault kv secrets engine data is writable"
}

quality "vault_service_restart" {
  description = "Vault restarts with existing configuration"
}

quality "vault_service_start" {
  description = "Vault starts with the configuration"
}

quality "vault_service_systemd_notified" {
  description = "The Vault binary notifies systemd when the service is active"
}

quality "vault_service_systemd_unit" {
  description = "The 'vault.service' systemd unit starts the service"
}

quality "vault_status_seal_type" {
  description = global.description.verify_seal_type
}

quality "vault_storage_backend_consul" {
  description = "Vault operates using Consul for storage"
}

quality "vault_storage_backend_raft" {
  description = "Vault operates using integrated Raft storage"
}

quality "vault_ui_assets" {
  description = global.description.verify_ui
}

quality "vault_ui_test" {
  description = <<-EOF
    The Vault Web UI test suite runs against a live Vault server with the embedded static assets
  EOF
}

quality "vault_unseal_ha_leader_election" {
  description = "Vault performs a leader election after it is unsealed"
}

quality "vault_version_build_date" {
  description = "Vault's reported build date matches our expectations"
}

quality "vault_version_edition" {
  description = "Vault's reported edition matches our expectations"
}

quality "vault_version_release" {
  description = "Vault's reported release version matches our expectations"
}

quality "vault_billing_start_date" {
  description = "Vault's billing start date has adjusted to the latest billing year"
}
