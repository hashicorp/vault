# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

scenario "autopilot" {
  matrix {
    arch            = ["amd64", "arm64"]
    artifact_source = ["local", "crt", "artifactory"]
    artifact_type   = ["bundle", "package"]
    distro          = ["ubuntu", "rhel"]
    edition         = ["ent", "ent.fips1402", "ent.hsm", "ent.hsm.fips1402"]
    seal            = ["awskms", "shamir"]

    # Our local builder always creates bundles
    exclude {
      artifact_source = ["local"]
      artifact_type   = ["package"]
    }

    # HSM and FIPS 140-2 are only supported on amd64
    exclude {
      arch    = ["arm64"]
      edition = ["ent.fips1402", "ent.hsm", "ent.hsm.fips1402"]
    }
  }

  terraform_cli = terraform_cli.default
  terraform     = terraform.default
  providers = [
    provider.aws.default,
    provider.enos.ubuntu,
    provider.enos.rhel
  ]

  locals {
    artifact_path = matrix.artifact_source != "artifactory" ? abspath(var.vault_artifact_path) : null
    enos_provider = {
      rhel   = provider.enos.rhel
      ubuntu = provider.enos.ubuntu
    }
    manage_service    = matrix.artifact_type == "bundle"
    vault_install_dir = matrix.artifact_type == "bundle" ? var.vault_install_dir : global.vault_install_dir_packages[matrix.distro]
  }

  step "build_vault" {
    module = "build_${matrix.artifact_source}"

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
    module = module.ec2_info
  }

  step "create_vpc" {
    module = module.create_vpc

    variables {
      common_tags = global.tags
    }
  }

  step "read_license" {
    module = module.read_license

    variables {
      file_name = global.vault_license_path
    }
  }

  step "create_vault_cluster_targets" {
    module     = module.target_ec2_instances
    depends_on = [step.create_vpc]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      ami_id                = step.ec2_info.ami_ids[matrix.arch][matrix.distro][global.distro_version[matrix.distro]]
      awskms_unseal_key_arn = step.create_vpc.kms_key_arn
      cluster_tag_key       = global.vault_tag_key
      common_tags           = global.tags
      vpc_id                = step.create_vpc.vpc_id
    }
  }

  step "create_vault_cluster" {
    module = module.vault_cluster
    depends_on = [
      step.build_vault,
      step.create_vault_cluster_targets
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      awskms_unseal_key_arn = step.create_vpc.kms_key_arn
      cluster_name          = step.create_vault_cluster_targets.cluster_name
      install_dir           = local.vault_install_dir
      license               = matrix.edition != "oss" ? step.read_license.license : null
      packages              = global.packages
      release               = var.vault_autopilot_initial_release
      storage_backend       = "raft"
      storage_backend_addl_config = {
        autopilot_upgrade_version = var.vault_autopilot_initial_release.version
      }
      target_hosts             = step.create_vault_cluster_targets.hosts
      unseal_method            = matrix.seal
      enable_file_audit_device = var.vault_enable_file_audit_device
    }
  }

  step "get_local_metadata" {
    skip_step = matrix.artifact_source != "local"
    module    = module.get_local_metadata
  }

  step "get_vault_cluster_ips" {
    module     = module.vault_get_cluster_ips
    depends_on = [step.create_vault_cluster]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances   = step.create_vault_cluster.target_hosts
      vault_install_dir = local.vault_install_dir
      vault_root_token  = step.create_vault_cluster.root_token
    }
  }

  step "verify_write_test_data" {
    module = module.vault_verify_write_data
    depends_on = [
      step.create_vault_cluster,
      step.get_vault_cluster_ips
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      leader_public_ip  = step.get_vault_cluster_ips.leader_public_ip
      leader_private_ip = step.get_vault_cluster_ips.leader_private_ip
      vault_instances   = step.create_vault_cluster.target_hosts
      vault_install_dir = local.vault_install_dir
      vault_root_token  = step.create_vault_cluster.root_token
    }
  }

  step "create_autopilot_upgrade_storageconfig" {
    module = module.autopilot_upgrade_storageconfig

    variables {
      vault_product_version = matrix.artifact_source == "local" ? step.get_local_metadata.version : var.vault_product_version
    }
  }

  step "create_vault_cluster_upgrade_targets" {
    module     = module.target_ec2_instances
    depends_on = [step.create_vpc]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      ami_id                = step.ec2_info.ami_ids[matrix.arch][matrix.distro][global.distro_version[matrix.distro]]
      awskms_unseal_key_arn = step.create_vpc.kms_key_arn
      common_tags           = global.tags
      cluster_name          = step.create_vault_cluster_targets.cluster_name
      vpc_id                = step.create_vpc.vpc_id
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
      awskms_unseal_key_arn       = step.create_vpc.kms_key_arn
      cluster_name                = step.create_vault_cluster_targets.cluster_name
      log_level                   = var.vault_log_level
      force_unseal                = matrix.seal == "shamir"
      initialize_cluster          = false
      install_dir                 = local.vault_install_dir
      license                     = matrix.edition != "oss" ? step.read_license.license : null
      local_artifact_path         = local.artifact_path
      manage_service              = local.manage_service
      packages                    = global.packages
      root_token                  = step.create_vault_cluster.root_token
      shamir_unseal_keys          = matrix.seal == "shamir" ? step.create_vault_cluster.unseal_keys_hex : null
      storage_backend             = "raft"
      storage_backend_addl_config = step.create_autopilot_upgrade_storageconfig.storage_addl_config
      storage_node_prefix         = "upgrade_node"
      target_hosts                = step.create_vault_cluster_upgrade_targets.hosts
      unseal_method               = matrix.seal
      enable_file_audit_device    = var.vault_enable_file_audit_device
    }
  }

  step "verify_vault_unsealed" {
    module = module.vault_verify_unsealed
    depends_on = [
      step.create_vault_cluster,
      step.create_vault_cluster_upgrade_targets,
      step.upgrade_vault_cluster_with_autopilot,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_install_dir = local.vault_install_dir
      vault_instances   = step.upgrade_vault_cluster_with_autopilot.target_hosts
    }
  }

  step "verify_raft_auto_join_voter" {
    module = module.vault_verify_raft_auto_join_voter
    depends_on = [
      step.upgrade_vault_cluster_with_autopilot,
      step.verify_vault_unsealed
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_install_dir = local.vault_install_dir
      vault_instances   = step.upgrade_vault_cluster_with_autopilot.target_hosts
      vault_root_token  = step.upgrade_vault_cluster_with_autopilot.root_token
    }
  }

  step "verify_autopilot_await_server_removal_state" {
    module = module.vault_verify_autopilot
    depends_on = [
      step.create_vault_cluster_upgrade_targets,
      step.upgrade_vault_cluster_with_autopilot,
      step.verify_raft_auto_join_voter
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_autopilot_upgrade_version = matrix.artifact_source == "local" ? step.get_local_metadata.version : var.vault_product_version
      vault_autopilot_upgrade_status  = "await-server-removal"
      vault_install_dir               = local.vault_install_dir
      vault_instances                 = step.create_vault_cluster.target_hosts
      vault_root_token                = step.upgrade_vault_cluster_with_autopilot.root_token
    }
  }

  step "get_updated_vault_cluster_ips" {
    module = module.vault_get_cluster_ips
    depends_on = [
      step.create_vault_cluster,
      step.create_vault_cluster_upgrade_targets,
      step.get_vault_cluster_ips,
      step.upgrade_vault_cluster_with_autopilot
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances       = step.create_vault_cluster.target_hosts
      vault_install_dir     = local.vault_install_dir
      vault_root_token      = step.create_vault_cluster.root_token
      node_public_ip        = step.get_vault_cluster_ips.leader_public_ip
      added_vault_instances = step.upgrade_vault_cluster_with_autopilot.target_hosts
    }
  }

  step "verify_read_test_data" {
    module = module.vault_verify_read_data
    depends_on = [
      step.get_updated_vault_cluster_ips,
      step.verify_write_test_data,
      step.upgrade_vault_cluster_with_autopilot,
      step.verify_raft_auto_join_voter
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      node_public_ips      = step.get_updated_vault_cluster_ips.follower_public_ips
      vault_instance_count = 6
      vault_install_dir    = local.vault_install_dir
    }
  }

  step "raft_remove_peers" {
    module = module.vault_raft_remove_peer
    depends_on = [
      step.create_vault_cluster_upgrade_targets,
      step.get_updated_vault_cluster_ips,
      step.upgrade_vault_cluster_with_autopilot,
      step.verify_autopilot_await_server_removal_state
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      operator_instance      = step.get_updated_vault_cluster_ips.leader_public_ip
      remove_vault_instances = step.create_vault_cluster.target_hosts
      vault_install_dir      = local.vault_install_dir
      vault_instance_count   = 3
      vault_root_token       = step.create_vault_cluster.root_token
    }
  }

  step "remove_old_nodes" {
    module = module.shutdown_multiple_nodes
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
    module = module.vault_verify_autopilot
    depends_on = [
      step.create_vault_cluster_upgrade_targets,
      step.upgrade_vault_cluster_with_autopilot,
      step.verify_raft_auto_join_voter,
      step.remove_old_nodes
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_autopilot_upgrade_version = matrix.artifact_source == "local" ? step.get_local_metadata.version : var.vault_product_version
      vault_autopilot_upgrade_status  = "idle"
      vault_install_dir               = local.vault_install_dir
      vault_instances                 = step.upgrade_vault_cluster_with_autopilot.target_hosts
      vault_root_token                = step.create_vault_cluster.root_token
    }
  }

  step "verify_undo_logs_status" {
    skip_step = semverconstraint(var.vault_product_version, "<1.13.0-0")
    module    = module.vault_verify_undo_logs
    depends_on = [
      step.create_vault_cluster_upgrade_targets,
      step.remove_old_nodes,
      step.upgrade_vault_cluster_with_autopilot,
      step.verify_autopilot_idle_state
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_install_dir = local.vault_install_dir
      vault_instances   = step.upgrade_vault_cluster_with_autopilot.target_hosts
      vault_root_token  = step.create_vault_cluster.root_token
    }
  }

  output "awskms_unseal_key_arn" {
    description = "The Vault cluster KMS key arn"
    value       = step.create_vpc.kms_key_arn
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

  output "vault_audit_device_file_path" {
    description = "The file path for the file audit device, if enabled"
    value       = step.create_vault_cluster.audit_device_file_path
  }
}
