# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

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
    arch                 = "amd64"
    artifact_type        = "bundle"
    backend_license_path = abspath(var.backend_license_path != null ? var.backend_license_path : joinpath(path.root, "./support/consul.hclic"))
    backend_tag_key      = "VaultStorage"
    build_tags = {
      "oss" = ["ui"]
      "ent" = ["ui", "enterprise", "ent"]
    }
    bundle_path    = abspath(var.vault_artifact_path)
    distro         = "ubuntu"
    consul_version = "1.14.2"
    seal           = "awskms"
    tags = merge({
      "Project Name" : var.project_name
      "Project" : "Enos",
      "Environment" : "ci"
    }, var.tags)
    vault_install_dir_packages = {
      rhel   = "/bin"
      ubuntu = "/usr/bin"
    }
    vault_install_dir  = var.vault_install_dir
    vault_license_path = abspath(var.vault_license_path != null ? var.vault_license_path : joinpath(path.root, "./support/vault.hclic"))
    vault_tag_key      = "Type" // enos_vault_start expects Type as the tag key
    ui_test_filter     = var.ui_test_filter != null && try(trimspace(var.ui_test_filter), "") != "" ? var.ui_test_filter : (matrix.edition == "oss") ? "!enterprise" : null
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
    skip_step = matrix.backend == "raft" || var.backend_edition == "oss"
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
      enos = provider.enos.ubuntu
    }

    variables {
      ami_id                = step.ec2_info.ami_ids[local.arch][local.distro][var.ubuntu_distro_version]
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
      step.create_vault_cluster_backend_targets,
    ]

    providers = {
      enos = provider.enos.ubuntu
    }

    variables {
      cluster_name    = step.create_vault_cluster_backend_targets.cluster_name
      cluster_tag_key = local.backend_tag_key
      license         = (matrix.backend == "consul" && var.backend_edition == "ent") ? step.read_backend_license.license : null
      release = {
        edition = var.backend_edition
        version = local.consul_version
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
      enos = provider.enos.ubuntu
    }

    variables {
      awskms_unseal_key_arn   = step.create_vpc.kms_key_arn
      backend_cluster_name    = step.create_vault_cluster_backend_targets.cluster_name
      backend_cluster_tag_key = local.backend_tag_key
      cluster_name            = step.create_vault_cluster_targets.cluster_name
      consul_license          = (matrix.backend == "consul" && var.backend_edition == "ent") ? step.read_backend_license.license : null
      consul_release = matrix.backend == "consul" ? {
        edition = var.backend_edition
        version = local.consul_version
      } : null
      enable_file_audit_device = var.vault_enable_file_audit_device
      install_dir              = local.vault_install_dir
      license                  = matrix.edition != "oss" ? step.read_vault_license.license : null
      local_artifact_path      = local.bundle_path
      storage_backend          = matrix.backend
      target_hosts             = step.create_vault_cluster_targets.hosts
      unseal_method            = local.seal
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

  output "root_token" {
    description = "The Vault cluster root token"
    value       = step.create_vault_cluster.root_token
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

  output "unseal_keys_b64" {
    description = "The Vault cluster unseal keys"
    value       = step.create_vault_cluster.unseal_keys_b64
  }

  output "unseal_keys_hex" {
    description = "The Vault cluster unseal keys hex"
    value       = step.create_vault_cluster.unseal_keys_hex
  }
}
