// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

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
    ip_version      = global.ip_versions
    seal            = global.seals

    // Our local builder always creates bundles
    exclude {
      artifact_source = ["local"]
      artifact_type   = ["package"]
    }

    // There are no published versions of these artifacts yet. We'll update this to exclude older
    // versions after our initial publication of these editions for arm64.
    exclude {
      arch    = ["arm64"]
      edition = ["ent.fips1402", "ent.hsm", "ent.hsm.fips1402"]
    }

    // PKCS#11 can only be used on ent.hsm and ent.hsm.fips1402.
    exclude {
      seal    = ["pkcs11"]
      edition = [for e in matrix.edition : e if !strcontains(e, "hsm")]
    }

    // softhsm packages not available for leap/sles.
    exclude {
      seal   = ["pkcs11"]
      distro = ["leap", "sles"]
    }

    // Testing in IPV6 mode is currently implemented for integrated Raft storage only
    exclude {
      ip_version = ["6"]
      backend    = ["consul"]
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
      amzn   = provider.enos.ec2_user
      leap   = provider.enos.ec2_user
      rhel   = provider.enos.ec2_user
      sles   = provider.enos.ec2_user
      ubuntu = provider.enos.ubuntu
    }
    manage_service                     = matrix.artifact_type == "bundle"
    vault_install_dir                  = global.vault_install_dir[matrix.artifact_type]
    vault_autopilot_default_max_leases = semverconstraint(var.vault_upgrade_initial_version, ">=1.16.0-0") ? "300000" : ""
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
      ip_version  = matrix.ip_version
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
      instance_count  = 3
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
      ami_id          = step.ec2_info.ami_ids[matrix.arch][matrix.distro][global.distro_version[matrix.distro]]
      common_tags     = global.tags
      cluster_name    = step.create_vault_cluster_targets.cluster_name
      cluster_tag_key = global.vault_tag_key
      instance_count  = 3
      seal_key_names  = step.create_seal_key.resource_names
      vpc_id          = step.create_vpc.id
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
      quality.vault_listener_ipv4,
      quality.vault_listener_ipv6,
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
      hosts                = step.create_vault_cluster_targets.hosts
      install_dir          = local.vault_install_dir
      ip_version           = matrix.ip_version
      license              = matrix.edition != "ce" ? step.read_license.license : null
      packages             = concat(global.packages, global.distro_packages[matrix.distro][global.distro_version[matrix.distro]])
      release = {
        edition = matrix.edition
        version = var.vault_upgrade_initial_version
      }
      seal_attributes = step.create_seal_key.attributes
      seal_type       = matrix.seal
      storage_backend = "raft"
      storage_backend_addl_config = {
        autopilot_upgrade_version = var.vault_upgrade_initial_version
      }
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
      hosts             = step.create_vault_cluster_targets.hosts
      ip_version        = matrix.ip_version
      timeout           = 120 // seconds
      vault_addr        = step.create_vault_cluster.api_addr_localhost
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
      hosts             = step.create_vault_cluster.hosts
      ip_version        = matrix.ip_version
      vault_addr        = step.create_vault_cluster.api_addr_localhost
      vault_install_dir = local.vault_install_dir
      vault_root_token  = step.create_vault_cluster.root_token
    }
  }

  step "verify_secrets_engines_create" {
    description = global.description.verify_secrets_engines_create
    module      = module.vault_verify_secrets_engines_create
    depends_on = [
      step.create_vault_cluster,
      step.get_vault_cluster_ips
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_api_auth_userpass_login_write,
      quality.vault_api_auth_userpass_user_write,
      quality.vault_api_identity_entity_write,
      quality.vault_api_identity_entity_alias_write,
      quality.vault_api_identity_group_write,
      quality.vault_api_identity_oidc_config_write,
      quality.vault_api_identity_oidc_introspect_write,
      quality.vault_api_identity_oidc_key_write,
      quality.vault_api_identity_oidc_key_rotate_write,
      quality.vault_api_identity_oidc_role_write,
      quality.vault_api_identity_oidc_token_read,
      quality.vault_api_sys_auth_userpass_user_write,
      quality.vault_api_sys_policy_write,
      quality.vault_mount_auth,
      quality.vault_mount_kv,
      quality.vault_secrets_kv_write,
    ]

    variables {
      hosts             = step.create_vault_cluster.hosts
      leader_host       = step.get_vault_cluster_ips.leader_host
      vault_addr        = step.create_vault_cluster.api_addr_localhost
      vault_install_dir = local.vault_install_dir
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
      step.verify_secrets_engines_create
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      artifactory_release         = matrix.artifact_source == "artifactory" ? step.build_vault.vault_artifactory_release : null
      cluster_name                = step.create_vault_cluster_targets.cluster_name
      config_mode                 = matrix.config_mode
      enable_audit_devices        = var.vault_enable_audit_devices
      force_unseal                = matrix.seal == "shamir"
      hosts                       = step.create_vault_cluster_upgrade_targets.hosts
      initialize_cluster          = false
      install_dir                 = local.vault_install_dir
      ip_version                  = matrix.ip_version
      license                     = matrix.edition != "ce" ? step.read_license.license : null
      local_artifact_path         = local.artifact_path
      log_level                   = var.vault_log_level
      manage_service              = local.manage_service
      packages                    = concat(global.packages, global.distro_packages[matrix.distro][global.distro_version[matrix.distro]])
      root_token                  = step.create_vault_cluster.root_token
      seal_attributes             = step.create_seal_key.attributes
      seal_type                   = matrix.seal
      shamir_unseal_keys          = matrix.seal == "shamir" ? step.create_vault_cluster.unseal_keys_hex : null
      storage_backend             = "raft"
      storage_backend_addl_config = step.create_autopilot_upgrade_storageconfig.storage_addl_config
      storage_node_prefix         = "upgrade_node"
    }
  }

  step "verify_vault_unsealed" {
    description = global.description.verify_vault_unsealed
    module      = module.vault_wait_for_cluster_unsealed
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
      hosts             = step.upgrade_vault_cluster_with_autopilot.hosts
      vault_addr        = step.upgrade_vault_cluster_with_autopilot.api_addr_localhost
      vault_install_dir = local.vault_install_dir
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
      hosts             = step.upgrade_vault_cluster_with_autopilot.hosts
      ip_version        = matrix.ip_version
      vault_addr        = step.upgrade_vault_cluster_with_autopilot.api_addr_localhost
      vault_install_dir = local.vault_install_dir
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
      hosts                           = step.create_vault_cluster.hosts
      vault_addr                      = step.create_vault_cluster.api_addr_localhost
      vault_autopilot_upgrade_version = matrix.artifact_source == "local" ? step.get_local_metadata.version : var.vault_product_version
      vault_autopilot_upgrade_status  = "await-server-removal"
      vault_install_dir               = local.vault_install_dir
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
      hosts             = step.upgrade_vault_cluster_with_autopilot.hosts
      ip_version        = matrix.ip_version
      timeout           = 120 // seconds
      vault_addr        = step.create_vault_cluster.api_addr_localhost
      vault_install_dir = local.vault_install_dir
      vault_root_token  = step.create_vault_cluster.root_token
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
      hosts             = step.upgrade_vault_cluster_with_autopilot.hosts
      ip_version        = matrix.ip_version
      vault_addr        = step.create_vault_cluster.api_addr_localhost
      vault_install_dir = local.vault_install_dir
      vault_root_token  = step.create_vault_cluster.root_token
    }
  }

  step "verify_secrets_engines_read" {
    description = global.description.verify_secrets_engines_read
    module      = module.vault_verify_secrets_engines_read
    depends_on = [
      step.get_updated_vault_cluster_ips,
      step.verify_secrets_engines_create,
      step.upgrade_vault_cluster_with_autopilot,
      step.verify_raft_auto_join_voter
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_api_auth_userpass_login_write,
      quality.vault_api_identity_entity_read,
      quality.vault_api_identity_oidc_config_read,
      quality.vault_api_identity_oidc_key_read,
      quality.vault_api_identity_oidc_role_read,
      quality.vault_secrets_kv_read
    ]

    variables {
      create_state      = step.verify_secrets_engines_create.state
      hosts             = step.get_updated_vault_cluster_ips.follower_hosts
      vault_addr        = step.upgrade_vault_cluster_with_autopilot.api_addr_localhost
      vault_install_dir = local.vault_install_dir
    }
  }

  step "verify_log_secrets" {
    skip_step = !var.vault_enable_audit_devices || !var.verify_log_secrets

    description = global.description.verify_log_secrets
    module      = module.verify_log_secrets
    depends_on = [
      step.verify_secrets_engines_read,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_audit_log_secrets,
      quality.vault_journal_secrets,
      quality.vault_radar_index_create,
      quality.vault_radar_scan_file,
    ]

    variables {
      audit_log_file_path = step.create_vault_cluster.audit_device_file_path
      leader_host         = step.get_updated_vault_cluster_ips.leader_host
      vault_addr          = step.upgrade_vault_cluster_with_autopilot.api_addr_localhost
      vault_root_token    = step.create_vault_cluster.root_token
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
      hosts                   = step.create_vault_cluster.hosts
      ip_version              = matrix.ip_version
      operator_instance       = step.get_updated_vault_cluster_ips.leader_public_ip
      vault_addr              = step.upgrade_vault_cluster_with_autopilot.api_addr_localhost
      vault_cluster_addr_port = step.upgrade_vault_cluster_with_autopilot.cluster_port
      vault_install_dir       = local.vault_install_dir
      vault_root_token        = step.create_vault_cluster.root_token
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
      old_hosts = step.create_vault_cluster.hosts
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
      hosts                           = step.upgrade_vault_cluster_with_autopilot.hosts
      vault_addr                      = step.upgrade_vault_cluster_with_autopilot.api_addr_localhost
      vault_autopilot_upgrade_version = matrix.artifact_source == "local" ? step.get_local_metadata.version : var.vault_product_version
      vault_autopilot_upgrade_status  = "idle"
      vault_install_dir               = local.vault_install_dir
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
      hosts         = step.upgrade_vault_cluster_with_autopilot.hosts
      vault_addr    = step.upgrade_vault_cluster_with_autopilot.api_addr_localhost
      vault_edition = matrix.edition
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
      quality.vault_api_sys_version_history_keys,
      quality.vault_api_sys_version_history_key_info,
      quality.vault_version_build_date,
      quality.vault_version_edition,
      quality.vault_version_release,
    ]

    variables {
      hosts                 = step.upgrade_vault_cluster_with_autopilot.hosts
      vault_addr            = step.upgrade_vault_cluster_with_autopilot.api_addr_localhost
      vault_edition         = matrix.edition
      vault_install_dir     = local.vault_install_dir
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
      hosts      = step.upgrade_vault_cluster_with_autopilot.hosts
      vault_addr = step.upgrade_vault_cluster_with_autopilot.api_addr_localhost
    }
  }

  step "verify_undo_logs_enabled_on_primary" {
    skip_step   = semverconstraint(var.vault_product_version, "<1.13.0-0")
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
      expected_state    = 1 # Enabled
      hosts             = step.get_updated_vault_cluster_ips.leader_hosts
      timeout           = 180 # Seconds
      vault_addr        = step.upgrade_vault_cluster_with_autopilot.api_addr_localhost
      vault_install_dir = local.vault_install_dir
      vault_root_token  = step.create_vault_cluster.root_token
    }
  }

  step "verify_undo_logs_disabled_on_followers" {
    skip_step  = semverconstraint(var.vault_product_version, "<1.13.0-0")
    module     = module.vault_verify_undo_logs
    depends_on = [step.verify_undo_logs_enabled_on_primary]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      expected_state    = 0 # Disabled
      hosts             = step.get_updated_vault_cluster_ips.follower_hosts
      timeout           = 10 # Seconds
      vault_addr        = step.upgrade_vault_cluster_with_autopilot.api_addr_localhost
      vault_install_dir = local.vault_install_dir
      vault_root_token  = step.create_vault_cluster.root_token
    }
  }

  // Verify that upgrading from a version <1.16.0 does not introduce Default LCQ
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
      hosts                              = step.upgrade_vault_cluster_with_autopilot.hosts
      vault_addr                         = step.upgrade_vault_cluster_with_autopilot.api_addr_localhost
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
    value       = step.create_vault_cluster.hosts
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

  output "secrets_engines_state" {
    description = "The state of configured secrets engines"
    value       = step.verify_secrets_engines_create.state
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
    value       = step.upgrade_vault_cluster_with_autopilot.hosts
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
