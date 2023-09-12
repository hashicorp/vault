# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

// The replication scenario configures performance replication between two Vault clusters and verifies
// known_primary_cluster_addrs are updated on secondary Vault cluster with the IP addresses of replaced
// nodes on primary Vault cluster
scenario "replication" {
  matrix {
    arch              = ["amd64", "arm64"]
    artifact_source   = ["local", "crt", "artifactory"]
    artifact_type     = ["bundle", "package"]
    consul_version    = ["1.14.2", "1.13.4", "1.12.7"]
    distro            = ["ubuntu", "rhel"]
    edition           = ["ent", "ent.fips1402", "ent.hsm", "ent.hsm.fips1402"]
    primary_backend   = ["raft", "consul"]
    primary_seal      = ["awskms", "shamir"]
    secondary_backend = ["raft", "consul"]
    secondary_seal    = ["awskms", "shamir"]

    # Packages are not offered for the   oss, ent.fips1402, and ent.hsm.fips1402 editions
    exclude {
      edition       = ["ent.fips1402", "ent.hsm.fips1402"]
      artifact_type = ["package"]
    }

    # Our local builder always creates bundles
    exclude {
      artifact_source = ["local"]
      artifact_type   = ["package"]
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
    # The path to the backend license file (Consul Enterprise)
    backend_license_path = abspath(var.backend_license_path != null ? var.backend_license_path : joinpath(path.root, "./support/consul.hclic"))
    backend_tag_key      = "VaultStorage"
    build_tags = {
      "ent"              = ["ui", "enterprise", "ent"]
      "ent.fips1402"     = ["ui", "enterprise", "cgo", "hsm", "fips", "fips_140_2", "ent.fips1402"]
      "ent.hsm"          = ["ui", "enterprise", "cgo", "hsm", "venthsm"]
      "ent.hsm.fips1402" = ["ui", "enterprise", "cgo", "hsm", "fips", "fips_140_2", "ent.hsm.fips1402"]
    }
    distro_version = {
      "rhel"   = var.rhel_distro_version
      "ubuntu" = var.ubuntu_distro_version
    }
    bundle_path = matrix.artifact_source != "artifactory" ? abspath(var.vault_artifact_path) : null
    enos_provider = {
      rhel   = provider.enos.rhel
      ubuntu = provider.enos.ubuntu
    }
    packages = ["jq"]
    tags = merge({
      "Project Name" : var.project_name
      "Project" : "Enos",
      "Environment" : "ci"
    }, var.tags)
    vault_license_path = abspath(var.vault_license_path != null ? var.vault_license_path : joinpath(path.root, "./support/vault.hclic"))
    vault_install_dir_packages = {
      rhel   = "/bin"
      ubuntu = "/usr/bin"
    }
    vault_install_dir = matrix.artifact_type == "bundle" ? var.vault_install_dir : local.vault_install_dir_packages[matrix.distro]
    vault_tag_key     = "Type" // enos_vault_start expects Type as the tag key
  }

  step "build_vault" {
    module = "build_${matrix.artifact_source}"

    variables {
      build_tags           = var.vault_local_build_tags != null ? var.vault_local_build_tags : local.build_tags[matrix.edition]
      bundle_path          = local.bundle_path
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
      common_tags = local.tags
    }
  }

  // This step reads the contents of the backend license if we're using a Consul backend and
  // the edition is "ent".
  step "read_backend_license" {
    skip_step = (matrix.primary_backend == "raft" && matrix.secondary_backend == "raft") || var.backend_edition == "oss"
    module    = module.read_license

    variables {
      file_name = local.backend_license_path
    }
  }

  step "read_vault_license" {
    module = module.read_license

    variables {
      file_name = abspath(joinpath(path.root, "./support/vault.hclic"))
    }
  }

  # Create all of our instances for both primary and secondary clusters
  step "create_primary_cluster_targets" {
    module = module.target_ec2_instances
    depends_on = [
      step.create_vpc,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      ami_id                = step.ec2_info.ami_ids[matrix.arch][matrix.distro][local.distro_version[matrix.distro]]
      awskms_unseal_key_arn = step.create_vpc.kms_key_arn
      cluster_tag_key       = local.vault_tag_key
      common_tags           = local.tags
      vpc_id                = step.create_vpc.vpc_id
    }
  }

  step "create_primary_cluster_backend_targets" {
    module = matrix.primary_backend == "consul" ? module.target_ec2_instances : module.target_ec2_shim
    depends_on = [
      step.create_vpc,
    ]

    providers = {
      enos = provider.enos.ubuntu
    }

    variables {
      ami_id                = step.ec2_info.ami_ids["arm64"]["ubuntu"]["22.04"]
      awskms_unseal_key_arn = step.create_vpc.kms_key_arn
      cluster_tag_key       = local.backend_tag_key
      common_tags           = local.tags
      vpc_id                = step.create_vpc.vpc_id
    }
  }

  step "create_primary_cluster_additional_targets" {
    module = module.target_ec2_instances
    depends_on = [
      step.create_vpc,
      step.create_primary_cluster_targets,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      ami_id                = step.ec2_info.ami_ids[matrix.arch][matrix.distro][local.distro_version[matrix.distro]]
      awskms_unseal_key_arn = step.create_vpc.kms_key_arn
      cluster_name          = step.create_primary_cluster_targets.cluster_name
      cluster_tag_key       = local.vault_tag_key
      common_tags           = local.tags
      vpc_id                = step.create_vpc.vpc_id
    }
  }

  step "create_secondary_cluster_targets" {
    module     = module.target_ec2_instances
    depends_on = [step.create_vpc]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      ami_id                = step.ec2_info.ami_ids[matrix.arch][matrix.distro][local.distro_version[matrix.distro]]
      awskms_unseal_key_arn = step.create_vpc.kms_key_arn
      cluster_tag_key       = local.vault_tag_key
      common_tags           = local.tags
      vpc_id                = step.create_vpc.vpc_id
    }
  }

  step "create_secondary_cluster_backend_targets" {
    module     = matrix.secondary_backend == "consul" ? module.target_ec2_instances : module.target_ec2_shim
    depends_on = [step.create_vpc]

    providers = {
      enos = provider.enos.ubuntu
    }

    variables {
      ami_id                = step.ec2_info.ami_ids["arm64"]["ubuntu"]["22.04"]
      awskms_unseal_key_arn = step.create_vpc.kms_key_arn
      cluster_tag_key       = local.backend_tag_key
      common_tags           = local.tags
      vpc_id                = step.create_vpc.vpc_id
    }
  }

  step "create_primary_backend_cluster" {
    module = "backend_${matrix.primary_backend}"
    depends_on = [
      step.create_primary_cluster_backend_targets,
    ]

    providers = {
      enos = provider.enos.ubuntu
    }

    variables {
      cluster_name    = step.create_primary_cluster_backend_targets.cluster_name
      cluster_tag_key = local.backend_tag_key
      license         = (matrix.primary_backend == "consul" && var.backend_edition == "ent") ? step.read_backend_license.license : null
      release = {
        edition = var.backend_edition
        version = matrix.consul_version
      }
      target_hosts = step.create_primary_cluster_backend_targets.hosts
    }
  }

  step "create_primary_cluster" {
    module = module.vault_cluster
    depends_on = [
      step.create_primary_backend_cluster,
      step.build_vault,
      step.create_primary_cluster_targets
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      artifactory_release     = matrix.artifact_source == "artifactory" ? step.build_vault.vault_artifactory_release : null
      awskms_unseal_key_arn   = step.create_vpc.kms_key_arn
      backend_cluster_name    = step.create_primary_cluster_backend_targets.cluster_name
      backend_cluster_tag_key = local.backend_tag_key
      consul_license          = (matrix.primary_backend == "consul" && var.backend_edition == "ent") ? step.read_backend_license.license : null
      cluster_name            = step.create_primary_cluster_targets.cluster_name
      consul_release = matrix.primary_backend == "consul" ? {
        edition = var.backend_edition
        version = matrix.consul_version
      } : null
      enable_file_audit_device = var.vault_enable_file_audit_device
      install_dir              = local.vault_install_dir
      license                  = matrix.edition != "oss" ? step.read_vault_license.license : null
      local_artifact_path      = local.bundle_path
      packages                 = local.packages
      storage_backend          = matrix.primary_backend
      target_hosts             = step.create_primary_cluster_targets.hosts
      unseal_method            = matrix.primary_seal
    }
  }

  step "create_secondary_backend_cluster" {
    module = "backend_${matrix.secondary_backend}"
    depends_on = [
      step.create_secondary_cluster_backend_targets
    ]

    providers = {
      enos = provider.enos.ubuntu
    }

    variables {
      cluster_name    = step.create_secondary_cluster_backend_targets.cluster_name
      cluster_tag_key = local.backend_tag_key
      license         = (matrix.secondary_backend == "consul" && var.backend_edition == "ent") ? step.read_backend_license.license : null
      release = {
        edition = var.backend_edition
        version = matrix.consul_version
      }
      target_hosts = step.create_secondary_cluster_backend_targets.hosts
    }
  }

  step "create_secondary_cluster" {
    module = module.vault_cluster
    depends_on = [
      step.create_secondary_backend_cluster,
      step.build_vault,
      step.create_secondary_cluster_targets
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      artifactory_release     = matrix.artifact_source == "artifactory" ? step.build_vault.vault_artifactory_release : null
      awskms_unseal_key_arn   = step.create_vpc.kms_key_arn
      backend_cluster_name    = step.create_secondary_cluster_backend_targets.cluster_name
      backend_cluster_tag_key = local.backend_tag_key
      consul_license          = (matrix.secondary_backend == "consul" && var.backend_edition == "ent") ? step.read_backend_license.license : null
      cluster_name            = step.create_secondary_cluster_targets.cluster_name
      consul_release = matrix.secondary_backend == "consul" ? {
        edition = var.backend_edition
        version = matrix.consul_version
      } : null
      enable_file_audit_device = var.vault_enable_file_audit_device
      install_dir              = local.vault_install_dir
      license                  = matrix.edition != "oss" ? step.read_vault_license.license : null
      local_artifact_path      = local.bundle_path
      packages                 = local.packages
      storage_backend          = matrix.secondary_backend
      target_hosts             = step.create_secondary_cluster_targets.hosts
      unseal_method            = matrix.secondary_seal
    }
  }

  step "verify_that_vault_primary_cluster_is_unsealed" {
    module = module.vault_verify_unsealed
    depends_on = [
      step.create_primary_cluster
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances   = step.create_primary_cluster_targets.hosts
      vault_install_dir = local.vault_install_dir
    }
  }

  step "verify_that_vault_secondary_cluster_is_unsealed" {
    module = module.vault_verify_unsealed
    depends_on = [
      step.create_secondary_cluster
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances   = step.create_secondary_cluster_targets.hosts
      vault_install_dir = local.vault_install_dir
    }
  }

  step "get_primary_cluster_ips" {
    module     = module.vault_get_cluster_ips
    depends_on = [step.verify_that_vault_primary_cluster_is_unsealed]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances   = step.create_primary_cluster_targets.hosts
      vault_install_dir = local.vault_install_dir
      vault_root_token  = step.create_primary_cluster.root_token
    }
  }

  step "get_secondary_cluster_ips" {
    module     = module.vault_get_cluster_ips
    depends_on = [step.verify_that_vault_secondary_cluster_is_unsealed]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances   = step.create_secondary_cluster_targets.hosts
      vault_install_dir = local.vault_install_dir
      vault_root_token  = step.create_secondary_cluster.root_token
    }
  }

  step "write_test_data_on_primary" {
    module     = module.vault_verify_write_data
    depends_on = [step.get_primary_cluster_ips]


    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      leader_public_ip  = step.get_primary_cluster_ips.leader_public_ip
      leader_private_ip = step.get_primary_cluster_ips.leader_private_ip
      vault_instances   = step.create_primary_cluster_targets.hosts
      vault_install_dir = local.vault_install_dir
      vault_root_token  = step.create_primary_cluster.root_token
    }
  }

  step "configure_performance_replication_primary" {
    module = module.vault_setup_perf_primary
    depends_on = [
      step.get_primary_cluster_ips,
      step.write_test_data_on_primary
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      primary_leader_public_ip  = step.get_primary_cluster_ips.leader_public_ip
      primary_leader_private_ip = step.get_primary_cluster_ips.leader_private_ip
      vault_install_dir         = local.vault_install_dir
      vault_root_token          = step.create_primary_cluster.root_token
    }
  }

  step "generate_secondary_token" {
    module     = module.generate_secondary_token
    depends_on = [step.configure_performance_replication_primary]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      primary_leader_public_ip = step.get_primary_cluster_ips.leader_public_ip
      vault_install_dir        = local.vault_install_dir
      vault_root_token         = step.create_primary_cluster.root_token
    }
  }

  step "configure_performance_replication_secondary" {
    module     = module.vault_setup_perf_secondary
    depends_on = [step.generate_secondary_token]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      secondary_leader_public_ip  = step.get_secondary_cluster_ips.leader_public_ip
      secondary_leader_private_ip = step.get_secondary_cluster_ips.leader_private_ip
      vault_install_dir           = local.vault_install_dir
      vault_root_token            = step.create_secondary_cluster.root_token
      wrapping_token              = step.generate_secondary_token.secondary_token
    }
  }

  // After replication is enabled, the secondary cluster followers need to be unsealed
  // Secondary unseal keys are passed using the guide https://developer.hashicorp.com/vault/docs/enterprise/replication#seals
  step "unseal_secondary_followers" {
    module = module.vault_unseal_nodes
    depends_on = [
      step.create_primary_cluster,
      step.create_secondary_cluster,
      step.get_secondary_cluster_ips,
      step.configure_performance_replication_secondary
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      follower_public_ips = step.get_secondary_cluster_ips.follower_public_ips
      vault_install_dir   = local.vault_install_dir
      vault_unseal_keys   = matrix.primary_seal == "shamir" ? step.create_primary_cluster.unseal_keys_hex : step.create_primary_cluster.recovery_keys_hex
      vault_seal_type     = matrix.primary_seal == "shamir" ? matrix.primary_seal : matrix.secondary_seal
    }
  }

  step "verify_secondary_cluster_is_unsealed_after_enabling_replication" {
    module = module.vault_verify_unsealed
    depends_on = [
      step.unseal_secondary_followers
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances   = step.create_secondary_cluster_targets.hosts
      vault_install_dir = local.vault_install_dir
    }
  }

  step "verify_performance_replication" {
    module     = module.vault_verify_performance_replication
    depends_on = [step.verify_secondary_cluster_is_unsealed_after_enabling_replication]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      primary_leader_public_ip    = step.get_primary_cluster_ips.leader_public_ip
      primary_leader_private_ip   = step.get_primary_cluster_ips.leader_private_ip
      secondary_leader_public_ip  = step.get_secondary_cluster_ips.leader_public_ip
      secondary_leader_private_ip = step.get_secondary_cluster_ips.leader_private_ip
      vault_install_dir           = local.vault_install_dir
    }
  }

  step "verify_replicated_data" {
    module = module.vault_verify_read_data
    depends_on = [
      step.verify_performance_replication,
      step.get_secondary_cluster_ips,
      step.write_test_data_on_primary
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      node_public_ips   = step.get_secondary_cluster_ips.follower_public_ips
      vault_install_dir = local.vault_install_dir
    }
  }

  step "add_additional_nodes_to_primary_cluster" {
    module = module.vault_cluster
    depends_on = [
      step.create_vpc,
      step.create_primary_backend_cluster,
      step.create_primary_cluster,
      step.verify_replicated_data,
      step.create_primary_cluster_additional_targets
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      artifactory_release     = matrix.artifact_source == "artifactory" ? step.build_vault.vault_artifactory_release : null
      awskms_unseal_key_arn   = step.create_vpc.kms_key_arn
      backend_cluster_name    = step.create_primary_cluster_backend_targets.cluster_name
      backend_cluster_tag_key = local.backend_tag_key
      cluster_name            = step.create_primary_cluster_targets.cluster_name
      consul_license          = (matrix.primary_backend == "consul" && var.backend_edition == "ent") ? step.read_backend_license.license : null
      consul_release = matrix.primary_backend == "consul" ? {
        edition = var.backend_edition
        version = matrix.consul_version
      } : null
      force_unseal        = matrix.primary_seal == "shamir"
      initialize_cluster  = false
      install_dir         = local.vault_install_dir
      license             = matrix.edition != "oss" ? step.read_vault_license.license : null
      local_artifact_path = local.bundle_path
      packages            = local.packages
      root_token          = step.create_primary_cluster.root_token
      shamir_unseal_keys  = matrix.primary_seal == "shamir" ? step.create_primary_cluster.unseal_keys_hex : null
      storage_backend     = matrix.primary_backend
      storage_node_prefix = "newprimary_node"
      target_hosts        = step.create_primary_cluster_additional_targets.hosts
      unseal_method       = matrix.primary_seal
    }
  }

  step "verify_addtional_primary_nodes_are_unsealed" {
    module     = module.vault_verify_unsealed
    depends_on = [step.add_additional_nodes_to_primary_cluster]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances   = step.create_primary_cluster_additional_targets.hosts
      vault_install_dir = local.vault_install_dir
    }
  }

  step "verify_raft_auto_join_voter" {
    skip_step = matrix.primary_backend != "raft"
    module    = module.vault_verify_raft_auto_join_voter
    depends_on = [
      step.add_additional_nodes_to_primary_cluster,
      step.create_primary_cluster,
      step.verify_addtional_primary_nodes_are_unsealed
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances   = step.create_primary_cluster_additional_targets.hosts
      vault_install_dir = local.vault_install_dir
      vault_root_token  = step.create_primary_cluster.root_token
    }
  }

  step "remove_primary_follower_1" {
    module = module.shutdown_node
    depends_on = [
      step.get_primary_cluster_ips,
      step.verify_addtional_primary_nodes_are_unsealed
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      node_public_ip = step.get_primary_cluster_ips.follower_public_ip_1
    }
  }

  step "remove_primary_leader" {
    module = module.shutdown_node
    depends_on = [
      step.get_primary_cluster_ips,
      step.remove_primary_follower_1
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      node_public_ip = step.get_primary_cluster_ips.leader_public_ip
    }
  }

  step "get_updated_primary_cluster_ips" {
    module = module.vault_get_cluster_ips
    depends_on = [
      step.add_additional_nodes_to_primary_cluster,
      step.remove_primary_follower_1,
      step.remove_primary_leader
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances       = step.create_primary_cluster_targets.hosts
      vault_install_dir     = local.vault_install_dir
      added_vault_instances = step.create_primary_cluster_additional_targets.hosts
      vault_root_token      = step.create_primary_cluster.root_token
      node_public_ip        = step.get_primary_cluster_ips.follower_public_ip_2
    }
  }

  step "verify_updated_performance_replication" {
    module     = module.vault_verify_performance_replication
    depends_on = [step.get_updated_primary_cluster_ips]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      primary_leader_public_ip    = step.get_updated_primary_cluster_ips.leader_public_ip
      primary_leader_private_ip   = step.get_updated_primary_cluster_ips.leader_private_ip
      secondary_leader_public_ip  = step.get_secondary_cluster_ips.leader_public_ip
      secondary_leader_private_ip = step.get_secondary_cluster_ips.leader_private_ip
      vault_install_dir           = local.vault_install_dir
    }
  }

  output "audit_device_file_path" {
    description = "The file path for the file audit device, if enabled"
    value       = step.create_primary_cluster.audit_device_file_path
  }

  output "primary_cluster_hosts" {
    description = "The Vault primary cluster target hosts"
    value       = step.create_primary_cluster_targets.hosts
  }

  output "primary_cluster_additional_hosts" {
    description = "The Vault added new node on primary cluster target hosts"
    value       = step.create_primary_cluster_additional_targets.hosts
  }

  output "primary_cluster_root_token" {
    description = "The Vault primary cluster root token"
    value       = step.create_primary_cluster.root_token
  }

  output "primary_cluster_unseal_keys_b64" {
    description = "The Vault primary cluster unseal keys"
    value       = step.create_primary_cluster.unseal_keys_b64
  }

  output "primary_cluster_unseal_keys_hex" {
    description = "The Vault primary cluster unseal keys hex"
    value       = step.create_primary_cluster.unseal_keys_hex
  }

  output "primary_cluster_recovery_key_shares" {
    description = "The Vault primary cluster recovery key shares"
    value       = step.create_primary_cluster.recovery_key_shares
  }

  output "primary_cluster_recovery_keys_b64" {
    description = "The Vault primary cluster recovery keys b64"
    value       = step.create_primary_cluster.recovery_keys_b64
  }

  output "primary_cluster_recovery_keys_hex" {
    description = "The Vault primary cluster recovery keys hex"
    value       = step.create_primary_cluster.recovery_keys_hex
  }

  output "secondary_cluster_hosts" {
    description = "The Vault secondary cluster public IPs"
    value       = step.create_secondary_cluster_targets.hosts
  }

  output "initial_primary_replication_status" {
    description = "The Vault primary cluster performance replication status"
    value       = step.verify_performance_replication.primary_replication_status
  }

  output "initial_known_primary_cluster_addresses" {
    description = "The Vault secondary cluster performance replication status"
    value       = step.verify_performance_replication.known_primary_cluster_addrs
  }

  output "initial_secondary_performance_replication_status" {
    description = "The Vault secondary cluster performance replication status"
    value       = step.verify_performance_replication.secondary_replication_status
  }

  output "intial_primary_replication_data_secondaries" {
    description = "The Vault primary cluster secondaries connection status"
    value       = step.verify_performance_replication.primary_replication_data_secondaries
  }

  output "initial_secondary_replication_data_primaries" {
    description = "The Vault  secondary cluster primaries connection status"
    value       = step.verify_performance_replication.secondary_replication_data_primaries
  }

  output "updated_primary_replication_status" {
    description = "The Vault updated primary cluster performance replication status"
    value       = step.verify_updated_performance_replication.primary_replication_status
  }

  output "updated_known_primary_cluster_addresses" {
    description = "The Vault secondary cluster performance replication status"
    value       = step.verify_updated_performance_replication.known_primary_cluster_addrs
  }

  output "updated_secondary_replication_status" {
    description = "The Vault updated secondary cluster performance replication status"
    value       = step.verify_updated_performance_replication.secondary_replication_status
  }

  output "updated_primary_replication_data_secondaries" {
    description = "The Vault updated primary cluster secondaries connection status"
    value       = step.verify_updated_performance_replication.primary_replication_data_secondaries
  }

  output "updated_secondary_replication_data_primaries" {
    description = "The Vault updated secondary cluster primaries connection status"
    value       = step.verify_updated_performance_replication.secondary_replication_data_primaries
  }
}
