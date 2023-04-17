# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

scenario "ui" {
  matrix {
    edition = ["oss", "ent"]
    backend = ["consul", "raft"]
  }

  terraform_cli = terraform_cli.default
  terraform     = terraform.default
  providers = [
    provider.aws.default,
    provider.enos.ubuntu
  ]

  locals {
    arch           = "amd64"
    distro         = "ubuntu"
    seal           = "awskms"
    artifact_type  = "bundle"
    consul_version = "1.14.2"
    build_tags = {
      "oss" = ["ui"]
      "ent" = ["ui", "enterprise", "ent"]
    }
    bundle_path = abspath(var.vault_bundle_path)
    tags = merge({
      "Project Name" : var.project_name
      "Project" : "Enos",
      "Environment" : "ci"
    }, var.tags)
    vault_instance_types = {
      amd64 = "t3a.small"
      arm64 = "t4g.small"
    }
    vault_instance_type = coalesce(var.vault_instance_type, local.vault_instance_types[local.arch])
    vault_license_path  = abspath(var.vault_license_path != null ? var.vault_license_path : joinpath(path.root, "./support/vault.hclic"))
    vault_install_dir_packages = {
      rhel   = "/bin"
      ubuntu = "/usr/bin"
    }
    vault_install_dir = var.vault_install_dir
    ui_test_filter    = var.ui_test_filter != null && try(trimspace(var.ui_test_filter), "") != "" ? var.ui_test_filter : (matrix.edition == "oss") ? "!enterprise" : null
  }

  step "get_local_metadata" {
    module = module.get_local_metadata
  }

  step "build_vault" {
    module = module.build_local

    variables {
      build_tags      = var.vault_local_build_tags != null ? var.vault_local_build_tags : local.build_tags[matrix.edition]
      bundle_path     = local.bundle_path
      goarch          = local.arch
      goos            = "linux"
      product_version = var.vault_product_version
      artifact_type   = local.artifact_type
      revision        = var.vault_revision
    }
  }

  step "find_azs" {
    module = module.az_finder

    variables {
      instance_type = [
        var.backend_instance_type,
        local.vault_instance_type
      ]
    }
  }

  step "create_vpc" {
    module = module.create_vpc

    variables {
      ami_architectures  = [local.arch]
      availability_zones = step.find_azs.availability_zones
      common_tags        = local.tags
    }
  }

  step "read_license" {
    skip_step = matrix.edition == "oss"
    module    = module.read_license

    variables {
      file_name = local.vault_license_path
    }
  }

  step "create_backend_cluster" {
    module     = "backend_${matrix.backend}"
    depends_on = [step.create_vpc]

    providers = {
      enos = provider.enos.ubuntu
    }

    variables {
      ami_id      = step.create_vpc.ami_ids["ubuntu"]["amd64"]
      common_tags = local.tags
      consul_release = {
        edition = var.backend_edition
        version = local.consul_version
      }
      instance_type = var.backend_instance_type
      kms_key_arn   = step.create_vpc.kms_key_arn
      vpc_id        = step.create_vpc.vpc_id
    }
  }

  step "create_vault_cluster_targets" {
    module     = module.target_ec2_spot_fleet // "target_ec2_instances" can be used for on-demand instances
    depends_on = [step.create_vpc]

    providers = {
      enos = provider.enos.ubuntu
    }

    variables {
      ami_id                = step.create_vpc.ami_ids[local.distro][local.arch]
      awskms_unseal_key_arn = step.create_vpc.kms_key_arn
      common_tags           = local.tags
      instance_type         = local.vault_instance_type // only used for on-demand instances
      vpc_id                = step.create_vpc.vpc_id
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
      enos = provider.enos.ubuntu
    }

    variables {
      awskms_unseal_key_arn = step.create_vpc.kms_key_arn
      cluster_name          = step.create_vault_cluster_targets.cluster_name
      config_env_vars = {
        VAULT_LOG_LEVEL = var.vault_log_level
      }
      consul_cluster_tag = step.create_backend_cluster.consul_cluster_tag
      consul_release = matrix.backend == "consul" ? {
        edition = var.backend_edition
        version = local.consul_version
      } : null
      install_dir         = local.vault_install_dir
      license             = matrix.edition != "oss" ? step.read_license.license : null
      local_artifact_path = local.bundle_path
      storage_backend     = matrix.backend
      target_hosts        = step.create_vault_cluster_targets.hosts
      unseal_method       = local.seal
    }
  }

  step "test_ui" {
    module = module.vault_test_ui

    variables {
      vault_addr               = step.create_vault_cluster_targets.hosts[0].public_ip
      vault_root_token         = step.create_vault_cluster.root_token
      vault_unseal_keys        = step.create_vault_cluster.recovery_keys_b64
      vault_recovery_threshold = step.create_vault_cluster.recovery_threshold
      ui_test_filter           = local.ui_test_filter
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

  output "ui_test_environment" {
    value       = step.test_ui.ui_test_environment
    description = "The environment variables that are required in order to run the test:enos yarn target"
  }

  output "ui_test_stderr" {
    description = "The stderr of the ui tests that ran"
    value       = step.test_ui.ui_test_stderr
  }

  output "ui_test_stdout" {
    description = "The stdout of the ui tests that ran"
    value       = step.test_ui.ui_test_stdout
  }
}
