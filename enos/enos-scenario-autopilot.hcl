# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

scenario "autopilot" {
  description = <<-EOF
    The autopilot scenario verifies autopilot upgrades between previously released versions of
    Vault Enterprise against another candidate build. The build can be a local branch, any CRT built
    Vault Enterprise artifact saved to the local machine, or any CRT built Vault Enterprise artifact
    in the stable channel in Artifactory.

    The scenario creates a new Vault Cluster with a previously released version of Vault, mounts
    various engines and creates data, then perform an Autopilot upgrade with any candidate build.
    The scenario also performs standard baseline verification that is not specific to the autopilot
    upgrade.

    If you want to use the 'distro:leap' variant you must first accept SUSE's terms for the AWS
    account. To verify that your account has agreed, sign-in to your AWS through Doormat,
    and visit the following links to verify your subscription or subscribe:
      arm64 AMI: https://aws.amazon.com/marketplace/server/procurement?productId=a516e959-df54-4035-bb1a-63599b7a6df9
      amd64 AMI: https://aws.amazon.com/marketplace/server/procurement?productId=5535c495-72d4-4355-b169-54ffa874f849
  EOF

  matrix {
    arch            = global.archs
    artifact_source = global.artifact_sources
    artifact_type   = global.artifact_types
    config_mode     = global.config_modes
    distro          = global.distros
    edition         = global.enterprise_editions
    initial_version = global.upgrade_initial_versions_ent
    seal            = global.seals

    # Autopilot wasn't available before 1.11.x
    exclude {
      initial_version = [for e in matrix.initial_version : e if semverconstraint(e, "<1.11.0-0")]
    }

    # Our local builder always creates bundles
    exclude {
      artifact_source = ["local"]
      artifact_type   = ["package"]
    }

    # There are no published versions of these artifacts yet. We'll update this to exclude older
    # versions after our initial publication of these editions for arm64.
    exclude {
      arch    = ["arm64"]
      edition = ["ent.fips1402", "ent.hsm", "ent.hsm.fips1402"]
    }

    # PKCS#11 can only be used on ent.hsm and ent.hsm.fips1402.
    exclude {
      seal    = ["pkcs11"]
      edition = [for e in matrix.edition : e if !strcontains(e, "hsm")]
    }

    # arm64 AMIs are not offered for Leap
    exclude {
      distro = ["leap"]
      arch   = ["arm64"]
    }

    # softhsm packages not available for leap/sles. Enos support for softhsm on amzn2 is
    # not implemented yet.
    exclude {
      seal   = ["pkcs11"]
      distro = ["amzn2", "leap", "sles"]
    }
  }

  terraform_cli = terraform_cli.default
  terraform     = terraform.default
  providers = [
    provider.aws.default,
    provider.enos.ec2_user,
    provider.enos.ubuntu
  ]

  locals {
    artifact_path = matrix.artifact_source != "artifactory" ? abspath(var.vault_artifact_path) : null
    enos_provider = {
      amzn2  = provider.enos.ec2_user
      leap   = provider.enos.ec2_user
      rhel   = provider.enos.ec2_user
      sles   = provider.enos.ec2_user
      ubuntu = provider.enos.ubuntu
    }
    manage_service                     = matrix.artifact_type == "bundle"
    vault_install_dir                  = global.vault_install_dir[matrix.artifact_type]
    vault_autopilot_default_max_leases = semverconstraint(matrix.initial_version, ">=1.16.0-0") ? "300000" : ""
  }

  step "build_vault" {
    description = global.description.build_vault
    module      = "build_${matrix.artifact_source}"

    variables {
      build_tags           = var.vault_local_build_tags != null ? var.vault_local_build_tags : global.build_tags[matrix.edition]
      artifact_path        = local.artifact_path
      goarch               = matrix.arch
      goos                 = "linux"
      artifactory_host     = matrix.artifact_source == "artifactory" ? var.artifactory_host : null
      artifactory_repo     = matrix.artifact_source == "artifactory" ? var.artifactory_repo : null
      artifactory_username = matrix.artifact_source == "artifactory" ? var.artifactory_username : null
      artifactory_token    = matrix.artifact_source == "artifactory" ? var.artifactory_token : null
      arch                 = matrix.artifact_source == "artifactory" ? matrix.arch : null
      product_version      = var.vault_product_version
      artifact_type        = matrix.artifact_type
      distro               = matrix.artifact_source == "artifactory" ? matrix.distro : null
      edition              = matrix.artifact_source == "artifactory" ? matrix.edition : null
      revision             = var.vault_revision
    }
  }

  step "ec2_info" {
    description = global.description.ec2_info
    module      = module.ec2_info
  }

  step "create_vpc" {
    description = global.description.create_vpc
    module      = module.create_vpc

    variables {
      common_tags = global.tags
    }
  }

  step "read_license" {
    description = global.description.read_vault_license
    module      = module.read_license

    variables {
      file_name = global.vault_license_path
    }
  }

  step "create_seal_key" {
    description = global.description.create_seal_key
    module      = "seal_${matrix.seal}"
    depends_on  = [step.create_vpc]

    providers = {
      enos = provider.enos.ubuntu
    }

    variables {
      cluster_id  = step.create_vpc.id
      common_tags = global.tags
    }
  }

  step "create_vault_cluster_targets" {
    description = global.description.create_vault_cluster_targets
    module      = module.target_ec2_instances
    depends_on  = [step.create_vpc]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      ami_id          = step.ec2_info.ami_ids[matrix.arch][matrix.distro][global.distro_version[matrix.distro]]
      cluster_tag_key = global.vault_tag_key
      common_tags     = global.tags
      seal_key_names  = step.create_seal_key.resource_names
      vpc_id          = step.create_vpc.id
    }
  }

  step "create_vault_cluster_upgrade_targets" {
    description = global.description.create_vault_cluster_targets
    module      = module.target_ec2_instances
    depends_on  = [step.create_vpc]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      ami_id         = step.ec2_info.ami_ids[matrix.arch][matrix.distro][global.distro_version[matrix.distro]]
      common_tags    = global.tags
      cluster_name   = step.create_vault_cluster_targets.cluster_name
      seal_key_names = step.create_seal_key.resource_names
      vpc_id         = step.create_vpc.id
    }
  }

  step "create_vault_cluster" {
    description = <<-EOF
      ${global.description.create_vault_cluster} In this instance we'll create a Vault Cluster with
      and older version and use Autopilot to upgrade to it.
    EOF

    module = module.vault_cluster
    depends_on = [
      step.build_vault,
      step.create_vault_cluster_targets
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      // verified in modules
      quality.vault_artifact_bundle,
      quality.vault_artifact_deb,
      quality.vault_artifact_rpm,
      quality.vault_audit_log,
      quality.vault_audit_socket,
      quality.vault_audit_syslog,
      quality.vault_autojoin_aws,
      quality.vault_config_env_variables,
      quality.vault_config_file,
      quality.vault_config_log_level,
      quality.vault_init,
      quality.vault_license_required_ent,
      quality.vault_service_start,
      quality.vault_storage_backend_consul,
      quality.vault_storage_backend_raft,
      // verified in enos_vault_start resource
      quality.vault_api_sys_config_read,
      quality.vault_api_sys_ha_status_read,
      quality.vault_api_sys_health_read,
      quality.vault_api_sys_host_info_read,
      quality.vault_api_sys_seal_status_api_read_matches_sys_health,
      quality.vault_api_sys_storage_raft_autopilot_configuration_read,
      quality.vault_api_sys_storage_raft_autopilot_state_read,
      quality.vault_api_sys_storage_raft_configuration_read,
      quality.vault_api_sys_replication_status_read,
      quality.vault_cli_status_exit_code,
      quality.vault_service_systemd_notified,
      quality.vault_service_systemd_unit,
    ]

    variables {
      cluster_name         = step.create_vault_cluster_targets.cluster_name
      config_mode          = matrix.config_mode
      enable_audit_devices = var.vault_enable_audit_devices
      install_dir          = global.vault_install_dir[matrix.artifact_type]
      license              = matrix.edition != "ce" ? step.read_license.license : null
      packages             = concat(global.packages, global.distro_packages[matrix.distro])
      release = {
        edition = matrix.edition
        version = matrix.initial_version
      }
      seal_attributes = step.create_seal_key.attributes
      seal_type       = matrix.seal
      storage_backend = "raft"
      storage_backend_addl_config = {
        autopilot_upgrade_version = matrix.initial_version
      }
      target_hosts = step.create_vault_cluster_targets.hosts
    }
  }

  step "get_local_metadata" {
    description = global.description.get_local_metadata
    skip_step   = matrix.artifact_source != "local"
    module      = module.get_local_metadata
  }

  // Wait for our cluster to elect a leader
  step "wait_for_leader" {
    description = global.description.wait_for_cluster_to_have_leader
    module      = module.vault_wait_for_leader
    depends_on  = [step.create_vault_cluster]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_api_sys_leader_read,
      quality.vault_unseal_ha_leader_election,
    ]

    variables {
      timeout           = 120 # seconds
      vault_hosts       = step.create_vault_cluster_targets.hosts
      vault_install_dir = local.vault_install_dir
      vault_root_token  = step.create_vault_cluster.root_token
    }
  }

  step "get_vault_cluster_ips" {
    description = global.description.get_vault_cluster_ip_addresses
    module      = module.vault_get_cluster_ips
    depends_on = [
      step.create_vault_cluster,
      step.wait_for_leader,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_api_sys_ha_status_read,
      quality.vault_api_sys_leader_read,
      quality.vault_cli_operator_members,
    ]

    variables {
      vault_hosts       = step.create_vault_cluster.target_hosts
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      vault_root_token  = step.create_vault_cluster.root_token
    }
  }

  step "verify_write_test_data" {
    description = global.description.verify_write_test_data
    module      = module.vault_verify_write_data
    depends_on = [
      step.create_vault_cluster,
      step.get_vault_cluster_ips
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_mount_auth,
      quality.vault_mount_kv,
      quality.vault_secrets_auth_user_policy_write,
      quality.vault_secrets_kv_write,
    ]

    variables {
      leader_public_ip  = step.get_vault_cluster_ips.leader_public_ip
      leader_private_ip = step.get_vault_cluster_ips.leader_private_ip
      vault_instances   = step.create_vault_cluster.target_hosts
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      vault_root_token  = step.create_vault_cluster.root_token
    }
  }

  step "create_autopilot_upgrade_storageconfig" {
    description = <<-EOF
      An arithmetic module used to dynamically create autopilot storage configuration depending on
      whether or not we're testing a local build or a candidate build.
    EOF
    module      = module.autopilot_upgrade_storageconfig

    variables {
      vault_product_version = matrix.artifact_source == "local" ? step.get_local_metadata.version : var.vault_product_version
    }
  }

  step "upgrade_vault_cluster_with_autopilot" {
    module = module.vault_cluster
    depends_on = [
      step.build_vault,
      step.create_vault_cluster,
      step.create_autopilot_upgrade_storageconfig,
      step.verify_write_test_data
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      artifactory_release         = matrix.artifact_source == "artifactory" ? step.build_vault.vault_artifactory_release : null
      enable_audit_devices        = var.vault_enable_audit_devices
      cluster_name                = step.create_vault_cluster_targets.cluster_name
      config_mode                 = matrix.config_mode
      log_level                   = var.vault_log_level
      force_unseal                = matrix.seal == "shamir"
      initialize_cluster          = false
      install_dir                 = global.vault_install_dir[matrix.artifact_type]
      license                     = matrix.edition != "ce" ? step.read_license.license : null
      local_artifact_path         = local.artifact_path
      manage_service              = local.manage_service
      packages                    = concat(global.packages, global.distro_packages[matrix.distro])
      root_token                  = step.create_vault_cluster.root_token
      seal_attributes             = step.create_seal_key.attributes
      seal_type                   = matrix.seal
      shamir_unseal_keys          = matrix.seal == "shamir" ? step.create_vault_cluster.unseal_keys_hex : null
      storage_backend             = "raft"
      storage_backend_addl_config = step.create_autopilot_upgrade_storageconfig.storage_addl_config
      storage_node_prefix         = "upgrade_node"
      target_hosts                = step.create_vault_cluster_upgrade_targets.hosts
    }
  }

  step "verify_vault_unsealed" {
    description = global.description.verify_vault_unsealed
    module      = module.vault_verify_unsealed
    depends_on = [
      step.create_vault_cluster,
      step.create_vault_cluster_upgrade_targets,
      step.upgrade_vault_cluster_with_autopilot,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_auto_unseals_after_autopilot_upgrade,
      quality.vault_seal_awskms,
      quality.vault_seal_pkcs11,
      quality.vault_seal_shamir,
    ]

    variables {
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      vault_instances   = step.upgrade_vault_cluster_with_autopilot.target_hosts
    }
  }

  step "verify_raft_auto_join_voter" {
    description = global.description.verify_raft_cluster_all_nodes_are_voters
    module      = module.vault_verify_raft_auto_join_voter
    depends_on = [
      step.upgrade_vault_cluster_with_autopilot,
      step.verify_vault_unsealed
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = quality.vault_raft_voters

    variables {
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      vault_instances   = step.upgrade_vault_cluster_with_autopilot.target_hosts
      vault_root_token  = step.upgrade_vault_cluster_with_autopilot.root_token
    }
  }

  step "verify_autopilot_await_server_removal_state" {
    description = global.description.verify_autopilot_idle_state
    module      = module.vault_verify_autopilot
    depends_on = [
      step.create_vault_cluster_upgrade_targets,
      step.upgrade_vault_cluster_with_autopilot,
      step.verify_raft_auto_join_voter
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_api_sys_storage_raft_autopilot_upgrade_info_read_status_matches,
      quality.vault_api_sys_storage_raft_autopilot_upgrade_info_target_version_read_matches_candidate,
    ]

    variables {
      vault_autopilot_upgrade_version = matrix.artifact_source == "local" ? step.get_local_metadata.version : var.vault_product_version
      vault_autopilot_upgrade_status  = "await-server-removal"
      vault_install_dir               = global.vault_install_dir[matrix.artifact_type]
      vault_instances                 = step.create_vault_cluster.target_hosts
      vault_root_token                = step.upgrade_vault_cluster_with_autopilot.root_token
    }
  }

  step "wait_for_leader_in_upgrade_targets" {
    description = global.description.wait_for_cluster_to_have_leader
    module      = module.vault_wait_for_leader
    depends_on = [
      step.create_vault_cluster,
      step.create_vault_cluster_upgrade_targets,
      step.get_vault_cluster_ips,
      step.upgrade_vault_cluster_with_autopilot
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_api_sys_leader_read,
      quality.vault_autopilot_upgrade_leader_election,
    ]

    variables {
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      vault_root_token  = step.create_vault_cluster.root_token
      vault_hosts       = step.upgrade_vault_cluster_with_autopilot.target_hosts
    }
  }

  step "get_updated_vault_cluster_ips" {
    description = global.description.get_vault_cluster_ip_addresses
    module      = module.vault_get_cluster_ips
    depends_on = [
      step.create_vault_cluster,
      step.create_vault_cluster_upgrade_targets,
      step.get_vault_cluster_ips,
      step.upgrade_vault_cluster_with_autopilot,
      step.wait_for_leader_in_upgrade_targets,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_api_sys_ha_status_read,
      quality.vault_api_sys_leader_read,
      quality.vault_cli_operator_members,
    ]

    variables {
      vault_hosts       = step.upgrade_vault_cluster_with_autopilot.target_hosts
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      vault_root_token  = step.create_vault_cluster.root_token
    }
  }

  step "verify_read_test_data" {
    description = global.description.verify_read_test_data
    module      = module.vault_verify_read_data
    depends_on = [
      step.get_updated_vault_cluster_ips,
      step.verify_write_test_data,
      step.upgrade_vault_cluster_with_autopilot,
      step.verify_raft_auto_join_voter
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = quality.vault_secrets_kv_read

    variables {
      node_public_ips      = step.get_updated_vault_cluster_ips.follower_public_ips
      vault_instance_count = 6
      vault_install_dir    = global.vault_install_dir[matrix.artifact_type]
    }
  }

  step "raft_remove_peers" {
    description = <<-EOF
      Remove the nodes that were running the prior version of Vault from the raft cluster
    EOF
    module      = module.vault_raft_remove_peer
    depends_on = [
      step.create_vault_cluster_upgrade_targets,
      step.get_updated_vault_cluster_ips,
      step.upgrade_vault_cluster_with_autopilot,
      step.verify_autopilot_await_server_removal_state
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_api_sys_storage_raft_remove_peer_write_removes_peer,
      quality.vault_cli_operator_raft_remove_peer,
    ]

    variables {
      operator_instance      = step.get_updated_vault_cluster_ips.leader_public_ip
      remove_vault_instances = step.create_vault_cluster.target_hosts
      vault_install_dir      = global.vault_install_dir[matrix.artifact_type]
      vault_instance_count   = 3
      vault_root_token       = step.create_vault_cluster.root_token
    }
  }

  step "remove_old_nodes" {
    description = global.description.shutdown_nodes
    module      = module.shutdown_multiple_nodes
    depends_on = [
      step.create_vault_cluster,
      step.raft_remove_peers
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      old_vault_instances  = step.create_vault_cluster.target_hosts
      vault_instance_count = 3
    }
  }

  step "verify_autopilot_idle_state" {
    description = global.description.verify_autopilot_idle_state
    module      = module.vault_verify_autopilot
    depends_on = [
      step.create_vault_cluster_upgrade_targets,
      step.upgrade_vault_cluster_with_autopilot,
      step.verify_raft_auto_join_voter,
      step.remove_old_nodes
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_api_sys_storage_raft_autopilot_upgrade_info_read_status_matches,
      quality.vault_api_sys_storage_raft_autopilot_upgrade_info_target_version_read_matches_candidate,
    ]

    variables {
      vault_autopilot_upgrade_version = matrix.artifact_source == "local" ? step.get_local_metadata.version : var.vault_product_version
      vault_autopilot_upgrade_status  = "idle"
      vault_install_dir               = global.vault_install_dir[matrix.artifact_type]
      vault_instances                 = step.upgrade_vault_cluster_with_autopilot.target_hosts
      vault_root_token                = step.create_vault_cluster.root_token
    }
  }

  step "verify_replication" {
    description = global.description.verify_replication_status
    module      = module.vault_verify_replication
    depends_on = [
      step.create_vault_cluster_upgrade_targets,
      step.upgrade_vault_cluster_with_autopilot,
      step.verify_raft_auto_join_voter,
      step.remove_old_nodes
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_replication_ce_disabled,
      quality.vault_replication_ent_dr_available,
      quality.vault_replication_ent_pr_available,
    ]

    variables {
      vault_edition     = matrix.edition
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      vault_instances   = step.upgrade_vault_cluster_with_autopilot.target_hosts
    }
  }

  step "verify_vault_version" {
    description = global.description.verify_vault_version
    module      = module.vault_verify_version
    depends_on = [
      step.create_vault_cluster_upgrade_targets,
      step.upgrade_vault_cluster_with_autopilot,
      step.verify_raft_auto_join_voter,
      step.remove_old_nodes
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_version_build_date,
      quality.vault_version_edition,
      quality.vault_version_release,
    ]

    variables {
      vault_instances       = step.upgrade_vault_cluster_with_autopilot.target_hosts
      vault_edition         = matrix.edition
      vault_install_dir     = global.vault_install_dir[matrix.artifact_type]
      vault_product_version = matrix.artifact_source == "local" ? step.get_local_metadata.version : var.vault_product_version
      vault_revision        = matrix.artifact_source == "local" ? step.get_local_metadata.revision : var.vault_revision
      vault_build_date      = matrix.artifact_source == "local" ? step.get_local_metadata.build_date : var.vault_build_date
      vault_root_token      = step.create_vault_cluster.root_token
    }
  }

  step "verify_ui" {
    description = global.description.verify_ui
    module      = module.vault_verify_ui
    depends_on = [
      step.create_vault_cluster_upgrade_targets,
      step.upgrade_vault_cluster_with_autopilot,
      step.verify_raft_auto_join_voter,
      step.remove_old_nodes
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = quality.vault_ui_assets

    variables {
      vault_instances = step.upgrade_vault_cluster_with_autopilot.target_hosts
    }
  }

  step "verify_undo_logs_status" {
    skip_step = true
    # NOTE: temporarily disable undo logs checking until it is fixed. See VAULT-20259
    # skip_step = semverconstraint(var.vault_product_version, "<1.13.0-0")
    module      = module.vault_verify_undo_logs
    description = <<-EOF
      Verifies that undo logs is correctly enabled on newly upgraded target hosts. For this it will
      query the metrics system backend for the vault.core.replication.write_undo_logs gauge.
    EOF

    depends_on = [
      step.create_vault_cluster_upgrade_targets,
      step.remove_old_nodes,
      step.upgrade_vault_cluster_with_autopilot,
      step.verify_autopilot_idle_state
    ]

    verifies = quality.vault_api_sys_metrics_vault_core_replication_write_undo_logs_enabled

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      vault_instances   = step.upgrade_vault_cluster_with_autopilot.target_hosts
      vault_root_token  = step.create_vault_cluster.root_token
    }
  }

  # Verify that upgrading from a version <1.16.0 does not introduce Default LCQ
  step "verify_default_lcq" {
    description = <<-EOF
      Verify that the default max lease count is 300,000 when the upgraded nodes are running
      Vault >= 1.16.0.
    EOF
    module      = module.vault_verify_default_lcq
    depends_on = [
      step.create_vault_cluster_upgrade_targets,
      step.remove_old_nodes,
      step.upgrade_vault_cluster_with_autopilot,
      step.verify_autopilot_idle_state
    ]

    verifies = quality.vault_api_sys_quotas_lease_count_read_max_leases_default

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances                    = step.upgrade_vault_cluster_with_autopilot.target_hosts
      vault_root_token                   = step.create_vault_cluster.root_token
      vault_autopilot_default_max_leases = local.vault_autopilot_default_max_leases
    }
  }

  output "audit_device_file_path" {
    description = "The file path for the file audit device, if enabled"
    value       = step.create_vault_cluster.audit_device_file_path
  }

  output "cluster_name" {
    description = "The Vault cluster name"
    value       = step.create_vault_cluster.cluster_name
  }

  output "hosts" {
    description = "The Vault cluster target hosts"
    value       = step.create_vault_cluster.target_hosts
  }

  output "private_ips" {
    description = "The Vault cluster private IPs"
    value       = step.create_vault_cluster.private_ips
  }

  output "public_ips" {
    description = "The Vault cluster public IPs"
    value       = step.create_vault_cluster.public_ips
  }

  output "root_token" {
    description = "The Vault cluster root token"
    value       = step.create_vault_cluster.root_token
  }

  output "recovery_key_shares" {
    description = "The Vault cluster recovery key shares"
    value       = step.create_vault_cluster.recovery_key_shares
  }

  output "recovery_keys_b64" {
    description = "The Vault cluster recovery keys b64"
    value       = step.create_vault_cluster.recovery_keys_b64
  }

  output "recovery_keys_hex" {
    description = "The Vault cluster recovery keys hex"
    value       = step.create_vault_cluster.recovery_keys_hex
  }

  output "seal_attributes" {
    description = "The Vault cluster seal attributes"
    value       = step.create_seal_key.attributes
  }

  output "unseal_keys_b64" {
    description = "The Vault cluster unseal keys"
    value       = step.create_vault_cluster.unseal_keys_b64
  }

  output "unseal_keys_hex" {
    description = "The Vault cluster unseal keys hex"
    value       = step.create_vault_cluster.unseal_keys_hex
  }

  output "upgrade_hosts" {
    description = "The Vault cluster target hosts"
    value       = step.upgrade_vault_cluster_with_autopilot.target_hosts
  }

  output "upgrade_private_ips" {
    description = "The Vault cluster private IPs"
    value       = step.upgrade_vault_cluster_with_autopilot.private_ips
  }

  output "upgrade_public_ips" {
    description = "The Vault cluster public IPs"
    value       = step.upgrade_vault_cluster_with_autopilot.public_ips
  }
}
