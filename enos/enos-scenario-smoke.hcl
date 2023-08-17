# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

scenario "smoke" {
  matrix {
    arch            = ["amd64", "arm64"]
    backend         = ["consul", "raft"]
    artifact_source = ["local", "crt", "artifactory"]
    artifact_type   = ["bundle", "package"]
    consul_edition  = ["oss", "ent"]
    consul_version  = ["1.14.2", "1.13.4", "1.12.7"]
    distro          = ["ubuntu", "rhel"]
    edition         = ["oss", "ent", "ent.fips1402", "ent.hsm", "ent.hsm.fips1402"]
    seal            = ["awskms", "shamir"]

    # Packages are not offered for the oss, ent.fips1402, and ent.hsm.fips1402 editions
    exclude {
      edition       = ["oss", "ent.fips1402", "ent.hsm.fips1402"]
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
    backend_license_path = abspath(var.backend_license_path != null ? var.backend_license_path : joinpath(path.root, "./support/consul.hclic"))
    backend_tag_key      = "VaultStorage"
    build_tags = {
      "oss"              = ["ui"]
      "ent"              = ["ui", "enterprise", "ent"]
      "ent.fips1402"     = ["ui", "enterprise", "cgo", "hsm", "fips", "fips_140_2", "ent.fips1402"]
      "ent.hsm"          = ["ui", "enterprise", "cgo", "hsm", "venthsm"]
      "ent.hsm.fips1402" = ["ui", "enterprise", "cgo", "hsm", "fips", "fips_140_2", "ent.hsm.fips1402"]
    }
    bundle_path = matrix.artifact_source != "artifactory" ? abspath(var.vault_artifact_path) : null
    distro_version = {
      "rhel"   = var.rhel_distro_version
      "ubuntu" = var.ubuntu_distro_version
    }
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

  step "get_local_metadata" {
    skip_step = matrix.artifact_source != "local"
    module    = module.get_local_metadata
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
    skip_step = matrix.backend == "raft" || matrix.consul_edition == "oss"
    module    = module.read_license

    variables {
      file_name = local.backend_license_path
    }
  }

  step "read_vault_license" {
    skip_step = matrix.edition == "oss"
    module    = module.read_license

    variables {
      file_name = local.vault_license_path
    }
  }

  step "create_vault_cluster_targets" {
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

  step "create_vault_cluster_backend_targets" {
    module     = matrix.backend == "consul" ? module.target_ec2_instances : module.target_ec2_shim
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

  step "create_backend_cluster" {
    module = "backend_${matrix.backend}"
    depends_on = [
      step.create_vault_cluster_backend_targets
    ]

    providers = {
      enos = provider.enos.ubuntu
    }

    variables {
      cluster_name    = step.create_vault_cluster_backend_targets.cluster_name
      cluster_tag_key = local.backend_tag_key
      license         = (matrix.backend == "consul" && matrix.consul_edition == "ent") ? step.read_backend_license.license : null
      release = {
        edition = matrix.consul_edition
        version = matrix.consul_version
      }
      target_hosts = step.create_vault_cluster_backend_targets.hosts
    }
  }

  step "create_vault_cluster" {
    module = module.vault_cluster
    depends_on = [
      step.create_backend_cluster,
      step.build_vault,
      step.create_vault_cluster_targets
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      artifactory_release     = matrix.artifact_source == "artifactory" ? step.build_vault.vault_artifactory_release : null
      awskms_unseal_key_arn   = step.create_vpc.kms_key_arn
      backend_cluster_name    = step.create_vault_cluster_backend_targets.cluster_name
      backend_cluster_tag_key = local.backend_tag_key
      cluster_name            = step.create_vault_cluster_targets.cluster_name
      consul_license          = (matrix.backend == "consul" && matrix.consul_edition == "ent") ? step.read_backend_license.license : null
      consul_release = matrix.backend == "consul" ? {
        edition = matrix.consul_edition
        version = matrix.consul_version
      } : null
      enable_file_audit_device = var.vault_enable_file_audit_device
      install_dir              = local.vault_install_dir
      license                  = matrix.edition != "oss" ? step.read_vault_license.license : null
      local_artifact_path      = local.bundle_path
      packages                 = local.packages
      storage_backend          = matrix.backend
      target_hosts             = step.create_vault_cluster_targets.hosts
      unseal_method            = matrix.seal
    }
  }

  step "get_vault_cluster_ips" {
    module     = module.vault_get_cluster_ips
    depends_on = [step.create_vault_cluster]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances   = step.create_vault_cluster_targets.hosts
      vault_install_dir = local.vault_install_dir
      vault_root_token  = step.create_vault_cluster.root_token
    }
  }

  step "verify_vault_version" {
    module     = module.vault_verify_version
    depends_on = [step.create_vault_cluster]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances       = step.create_vault_cluster_targets.hosts
      vault_edition         = matrix.edition
      vault_install_dir     = local.vault_install_dir
      vault_product_version = matrix.artifact_source == "local" ? step.get_local_metadata.version : var.vault_product_version
      vault_revision        = matrix.artifact_source == "local" ? step.get_local_metadata.revision : var.vault_revision
      vault_build_date      = matrix.artifact_source == "local" ? step.get_local_metadata.build_date : var.vault_build_date
      vault_root_token      = step.create_vault_cluster.root_token
    }
  }

  step "verify_vault_unsealed" {
    module     = module.vault_verify_unsealed
    depends_on = [step.create_vault_cluster]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_install_dir = local.vault_install_dir
      vault_instances   = step.create_vault_cluster_targets.hosts
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
      vault_instances   = step.create_vault_cluster_targets.hosts
      vault_install_dir = local.vault_install_dir
      vault_root_token  = step.create_vault_cluster.root_token
    }
  }

  step "verify_raft_auto_join_voter" {
    skip_step  = matrix.backend != "raft"
    module     = module.vault_verify_raft_auto_join_voter
    depends_on = [step.create_vault_cluster]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_install_dir = local.vault_install_dir
      vault_instances   = step.create_vault_cluster_targets.hosts
      vault_root_token  = step.create_vault_cluster.root_token
    }
  }

  step "verify_replication" {
    module     = module.vault_verify_replication
    depends_on = [step.create_vault_cluster]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_edition     = matrix.edition
      vault_install_dir = local.vault_install_dir
      vault_instances   = step.create_vault_cluster_targets.hosts
    }
  }

  step "verify_read_test_data" {
    module = module.vault_verify_read_data
    depends_on = [
      step.verify_write_test_data,
      step.verify_replication
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      node_public_ips   = step.get_vault_cluster_ips.follower_public_ips
      vault_install_dir = local.vault_install_dir
    }
  }

  step "verify_ui" {
    module     = module.vault_verify_ui
    depends_on = [step.create_vault_cluster]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances   = step.create_vault_cluster_targets.hosts
      vault_install_dir = local.vault_install_dir
    }
  }

  output "audit_device_file_path" {
    description = "The file path for the file audit device, if enabled"
    value       = step.create_vault_cluster.audit_device_file_path
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
}
