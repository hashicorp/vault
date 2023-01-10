scenario "ui" {
  matrix {
    edition = ["oss", "ent"]
  }

  terraform_cli = terraform_cli.default
  terraform     = terraform.default
  providers = [
    provider.aws.default,
    provider.enos.ubuntu
  ]

  locals {
    arch          = "amd64"
    backend       = "raft"
    distro        = "ubuntu"
    seal          = "awskms"
    artifact_type = "bundle"
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
    ui_test_filter = var.ui_test_filter != null ? var.ui_test_filter : (matrix.edition == "oss") ? "!enterprise" : null
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

  step "create_vault_cluster" {
    module = module.vault_cluster
    depends_on = [
      step.create_vpc,
      step.build_vault,
    ]

    providers = {
      enos = provider.enos.ubuntu
    }

    variables {
      ami_id                    = step.create_vpc.ami_ids[local.distro][local.arch]
      common_tags               = local.tags
      instance_type             = local.vault_instance_type
      kms_key_arn               = step.create_vpc.kms_key_arn
      storage_backend           = local.backend
      unseal_method             = local.seal
      vault_local_artifact_path = local.bundle_path
      vault_install_dir         = local.vault_install_dir
      vault_license             = matrix.edition != "oss" ? step.read_license.license : null
      vpc_id                    = step.create_vpc.vpc_id
    }
  }

  step "test_ui" {
    module = module.vault_test_ui

    variables {
      vault_addr               = step.create_vault_cluster.instance_public_ips[0]
      vault_root_token         = step.create_vault_cluster.vault_root_token
      vault_unseal_keys        = step.create_vault_cluster.vault_recovery_keys_b64
      vault_recovery_threshold = step.create_vault_cluster.vault_recovery_threshold
      ui_test_filter           = local.ui_test_filter
    }
  }

  output "vault_cluster_instance_ids" {
    description = "The Vault cluster instance IDs"
    value       = step.create_vault_cluster.instance_ids
  }

  output "vault_cluster_pub_ips" {
    description = "The Vault cluster public IPs"
    value       = step.create_vault_cluster.instance_public_ips
  }

  output "vault_cluster_priv_ips" {
    description = "The Vault cluster private IPs"
    value       = step.create_vault_cluster.instance_private_ips
  }

  output "vault_cluster_key_id" {
    description = "The Vault cluster Key ID"
    value       = step.create_vault_cluster.key_id
  }

  output "vault_cluster_root_token" {
    description = "The Vault cluster root token"
    value       = step.create_vault_cluster.vault_root_token
  }

  output "vault_cluster_unseal_keys_b64" {
    description = "The Vault cluster unseal keys"
    value       = step.create_vault_cluster.vault_unseal_keys_b64
  }

  output "vault_cluster_unseal_keys_hex" {
    description = "The Vault cluster unseal keys hex"
    value       = step.create_vault_cluster.vault_unseal_keys_hex
  }

  output "vault_cluster_tag" {
    description = "The Vault cluster tag"
    value       = step.create_vault_cluster.vault_cluster_tag
  }

  output "ui_test_stderr" {
    description = "The stderr of the ui tests that ran"
    value       = step.test_ui.ui_test_stderr
  }

  output "ui_test_stdout" {
    description = "The stdout of the ui tests that ran"
    value       = step.test_ui.ui_test_stdout
  }

  output "ui_test_environment" {
    value       = step.test_ui.ui_test_environment
    description = "The environment variables that are required in order to run the test:enos yarn target"
  }
}
