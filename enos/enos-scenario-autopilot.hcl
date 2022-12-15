scenario "autopilot" {
  matrix {
    arch             = ["amd64", "arm64"]
    artifact_source  = ["local", "crt", "artifactory"]
    artifact_type    = ["bundle", "package"]
    distro           = ["ubuntu", "rhel"]
    edition          = ["ent", "ent.fips1402", "ent.hsm", "ent.hsm.fips1402"]
    seal             = ["awskms", "shamir"]
    undo_logs_status = ["0", "1"]
  }

  terraform_cli = terraform_cli.default
  terraform     = terraform.default
  providers = [
    provider.aws.default,
    provider.enos.ubuntu,
    provider.enos.rhel
  ]

  locals {
    build_tags = {
      "ent"              = ["ui", "enterprise", "ent"]
      "ent.fips1402"     = ["ui", "enterprise", "cgo", "hsm", "fips", "fips_140_2", "ent.fips1402"]
      "ent.hsm"          = ["ui", "enterprise", "cgo", "hsm", "venthsm"]
      "ent.hsm.fips1402" = ["ui", "enterprise", "cgo", "hsm", "fips", "fips_140_2", "ent.hsm.fips1402"]
    }
    bundle_path             = matrix.artifact_source != "artifactory" ? abspath(var.vault_bundle_path) : null
    dependencies_to_install = ["jq"]
    enos_provider = {
      rhel   = provider.enos.rhel
      ubuntu = provider.enos.ubuntu
    }
    tags = merge({
      "Project Name" : var.project_name
      "Project" : "Enos",
      "Environment" : "ci"
    }, var.tags)
    vault_instance_types = {
      amd64 = "t3a.small"
      arm64 = "t4g.small"
    }

    enable_undo_logs = matrix.undo_logs_status == "1" && semverconstraint(var.vault_product_version, ">=1.13.0-0") ? true : false

    vault_instance_type = coalesce(var.vault_instance_type, local.vault_instance_types[matrix.arch])
    vault_license_path  = abspath(var.vault_license_path != null ? var.vault_license_path : joinpath(path.root, "./support/vault.hclic"))
    vault_install_dir_packages = {
      rhel   = "/bin"
      ubuntu = "/usr/bin"
    }
    vault_install_dir = matrix.artifact_type == "bundle" ? var.vault_install_dir : local.vault_install_dir_packages[matrix.distro]
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
      instance_type        = matrix.artifact_source == "artifactory" ? local.vault_instance_type : null
      revision             = var.vault_revision
    }
  }

  step "find_azs" {
    module = module.az_finder

    variables {
      instance_type = [
        local.vault_instance_type
      ]
    }
  }

  step "create_vpc" {
    module     = module.create_vpc
    depends_on = [step.find_azs]

    variables {
      ami_architectures  = [matrix.arch]
      availability_zones = step.find_azs.availability_zones
      common_tags        = local.tags
    }
  }

  step "read_license" {
    module = module.read_license

    variables {
      file_name = local.vault_license_path
    }
  }

  # This step creates a Vault cluster using a bundle downloaded from
  # releases.hashicorp.com, with the version specified in var.vault_autopilot_initial_release
  step "create_vault_cluster" {
    module = module.vault_cluster
    depends_on = [
      step.create_vpc,
      step.build_vault,
    ]
    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      ami_id                  = step.create_vpc.ami_ids[matrix.distro][matrix.arch]
      common_tags             = local.tags
      dependencies_to_install = local.dependencies_to_install
      instance_type           = local.vault_instance_type
      kms_key_arn             = step.create_vpc.kms_key_arn
      storage_backend         = "raft"
      storage_backend_addl_config = {
        autopilot_upgrade_version = var.vault_autopilot_initial_release.version
      }
      unseal_method     = matrix.seal
      vault_install_dir = local.vault_install_dir
      vault_release     = var.vault_autopilot_initial_release
      vault_license     = step.read_license.license
      vpc_id            = step.create_vpc.vpc_id
    }
  }

  step "get_local_metadata" {
    skip_step = matrix.artifact_source != "local"
    module    = module.get_local_metadata
  }

  step "get_vault_cluster_ips" {
    module     = module.vault_cluster_ips
    depends_on = [step.create_vault_cluster]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances   = step.create_vault_cluster.vault_instances
      vault_install_dir = local.vault_install_dir
      vault_root_token  = step.create_vault_cluster.vault_root_token
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
      vault_instances   = step.create_vault_cluster.vault_instances
      vault_install_dir = local.vault_install_dir
      vault_root_token  = step.create_vault_cluster.vault_root_token
    }
  }

  step "create_autopilot_upgrade_storageconfig" {
    module = module.autopilot_upgrade_storageconfig

    variables {
      vault_product_version = matrix.artifact_source == "local" ? step.get_local_metadata.version : var.vault_product_version
    }
  }

  # This step creates a new Vault cluster using a bundle or package
  # from the matrix.artifact_source, with the var.vault_product_version
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
      ami_id                      = step.create_vpc.ami_ids[matrix.distro][matrix.arch]
      common_tags                 = local.tags
      dependencies_to_install     = local.dependencies_to_install
      instance_type               = local.vault_instance_type
      kms_key_arn                 = step.create_vpc.kms_key_arn
      storage_backend             = "raft"
      storage_backend_addl_config = step.create_autopilot_upgrade_storageconfig.storage_addl_config
      unseal_method               = matrix.seal
      vault_cluster_tag           = step.create_vault_cluster.vault_cluster_tag
      vault_init                  = false
      vault_install_dir           = local.vault_install_dir
      vault_license               = step.read_license.license
      vault_local_artifact_path   = local.bundle_path
      vault_artifactory_release   = matrix.artifact_source == "artifactory" ? step.build_vault.vault_artifactory_release : null
      vault_node_prefix           = "upgrade_node"
      vault_root_token            = step.create_vault_cluster.vault_root_token
      vault_unseal_when_no_init   = matrix.seal == "shamir"
      vault_unseal_keys           = matrix.seal == "shamir" ? step.create_vault_cluster.vault_unseal_keys_hex : null
      vpc_id                      = step.create_vpc.vpc_id
      vault_environment           = { "VAULT_REPLICATION_USE_UNDO_LOGS" : local.enable_undo_logs }
    }
  }

  step "get_updated_vault_cluster_ips" {
    module = module.vault_cluster_ips
    depends_on = [
      step.create_vault_cluster,
      step.get_vault_cluster_ips,
      step.upgrade_vault_cluster_with_autopilot
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances       = step.create_vault_primary_cluster.vault_instances
      vault_install_dir     = local.vault_install_dir
      added_vault_instances = step.upgrade_vault_cluster_with_autopilot.vault_instances
      vault_root_token      = step.create_vault_primary_cluster.vault_root_token
      node_public_ip        = step.get_vault_cluster_ips.leader_public_ip
    }
  }

  step "verify_vault_unsealed" {
    module = module.vault_verify_unsealed
    depends_on = [
      step.create_vault_cluster,
      step.upgrade_vault_cluster_with_autopilot,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_install_dir = local.vault_install_dir
      vault_instances   = step.create_vault_cluster.vault_instances
    }
  }

  step "verify_read_test_data" {
    module = module.vault_verify_read_data
    depends_on = [
      step.get_updated_vault_cluster_ips,
      step.verify_write_test_data,
      step.verify_vault_unsealed
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      node_public_ips   = step.get_vault_cluster_ips.follower_public_ips
      vault_install_dir = local.vault_install_dir
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
      vault_instances   = step.create_vault_cluster.vault_instances
      vault_root_token  = step.create_vault_cluster.vault_root_token
    }
  }

  step "verify_autopilot_upgraded_vault_cluster" {
    module = module.vault_verify_autopilot
    depends_on = [
      step.upgrade_vault_cluster_with_autopilot,
      step.verify_vault_unsealed
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_autopilot_upgrade_version = matrix.artifact_source == "local" ? step.get_local_metadata.version : var.vault_product_version
      vault_autopilot_upgrade_status  = "await-server-removal"
      vault_install_dir               = local.vault_install_dir
      vault_instances                 = step.create_vault_cluster.vault_instances
      vault_root_token                = step.create_vault_cluster.vault_root_token
    }
  }

  step "verify_undo_logs_status" {
    skip_step = semverconstraint(var.vault_product_version, "<1.13.0-0")
    module    = module.vault_verify_undo_logs
    depends_on = [
      step.upgrade_vault_cluster_with_autopilot,
      step.verify_vault_unsealed
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_install_dir               = local.vault_install_dir
      vault_autopilot_upgrade_version = matrix.artifact_source == "local" ? step.get_local_metadata.version : var.vault_product_version
      vault_undo_logs_status          = matrix.undo_logs_status
      vault_instances                 = step.upgrade_vault_cluster_with_autopilot.vault_instances
      vault_root_token                = step.create_vault_cluster.vault_root_token
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

  output "upgraded_vault_cluster_instance_ids" {
    description = "The Vault cluster instance IDs"
    value       = step.upgrade_vault_cluster_with_autopilot.instance_ids
  }

  output "upgraded_vault_cluster_pub_ips" {
    description = "The Vault cluster public IPs"
    value       = step.upgrade_vault_cluster_with_autopilot.instance_public_ips
  }

  output "upgraded_vault_cluster_priv_ips" {
    description = "The Vault cluster private IPs"
    value       = step.upgrade_vault_cluster_with_autopilot.instance_private_ips
  }
}
