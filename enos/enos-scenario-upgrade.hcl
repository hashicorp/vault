scenario "upgrade" {
  matrix {
    arch            = ["amd64", "arm64"]
    backend         = ["consul", "raft"]
    artifact_source = ["local", "crt", "artifactory"]
    artifact_type   = ["bundle", "package"]
    consul_version  = ["1.14.2", "1.13.4", "1.12.7"]
    distro          = ["ubuntu", "rhel"]
    edition         = ["oss", "ent", "ent.fips1402", "ent.hsm", "ent.hsm.fips1402"]
    seal            = ["awskms", "shamir"]

    # Packages are not offered for the oss edition
    exclude {
      edition       = ["oss"]
      artifact_type = ["package"]
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
    build_tags = {
      "oss"              = ["ui"]
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
    vault_instance_type = coalesce(var.vault_instance_type, local.vault_instance_types[matrix.arch])
    vault_license_path  = abspath(var.vault_license_path != null ? var.vault_license_path : joinpath(path.root, "./support/vault.hclic"))
    vault_install_dir_packages = {
      rhel   = "/bin"
      ubuntu = "/usr/bin"
    }
    vault_install_dir = matrix.artifact_type == "bundle" ? var.vault_install_dir : local.vault_install_dir_packages[matrix.distro]
  }

  # This step gets/builds the upgrade artifact that we will upgrade to
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
        var.backend_instance_type,
        local.vault_instance_type,
      ]
    }
  }

  step "create_vpc" {
    module = module.create_vpc

    variables {
      ami_architectures  = distinct([matrix.arch, "amd64"])
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

  step "get_local_metadata" {
    skip_step = matrix.artifact_source != "local"
    module    = module.get_local_metadata
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
        version = matrix.consul_version
      }
      instance_type = var.backend_instance_type
      kms_key_arn   = step.create_vpc.kms_key_arn
      vpc_id        = step.create_vpc.vpc_id
    }
  }

  # This step creates a Vault cluster using a bundle downloaded from
  # releases.hashicorp.com, with the version specified in var.vault_autopilot_initial_release
  step "create_vault_cluster" {
    module = module.vault_cluster
    depends_on = [
      step.create_backend_cluster,
      step.build_vault,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      ami_id             = step.create_vpc.ami_ids[matrix.distro][matrix.arch]
      common_tags        = local.tags
      consul_cluster_tag = step.create_backend_cluster.consul_cluster_tag
      consul_release = matrix.backend == "consul" ? {
        edition = var.backend_edition
        version = matrix.consul_version
      } : null
      dependencies_to_install = local.dependencies_to_install
      instance_type           = local.vault_instance_type
      kms_key_arn             = step.create_vpc.kms_key_arn
      storage_backend         = matrix.backend
      unseal_method           = matrix.seal
      vault_install_dir       = local.vault_install_dir
      vault_release           = var.vault_upgrade_initial_release
      vault_license           = matrix.edition != "oss" ? step.read_license.license : null
      vpc_id                  = step.create_vpc.vpc_id
    }
  }

  step "get_vault_cluster_ips" {
    module     = module.vault_get_cluster_ips
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

  # This step upgrades the Vault cluster to the var.vault_product_version
  # by getting a bundle or package of that version from the matrix.artifact_source
  step "upgrade_vault" {
    module = module.vault_upgrade
    depends_on = [
      step.create_vault_cluster,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_api_addr            = "http://localhost:8200"
      vault_instances           = step.create_vault_cluster.vault_instances
      vault_local_artifact_path = local.bundle_path
      vault_artifactory_release = matrix.artifact_source == "artifactory" ? step.build_vault.vault_artifactory_release : null
      vault_install_dir         = local.vault_install_dir
      vault_unseal_keys         = matrix.seal == "shamir" ? step.create_vault_cluster.vault_unseal_keys_hex : null
      vault_seal_type           = matrix.seal
    }
  }

  step "verify_vault_version" {
    module = module.vault_verify_version
    depends_on = [
      step.create_backend_cluster,
      step.upgrade_vault,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances       = step.create_vault_cluster.vault_instances
      vault_edition         = matrix.edition
      vault_install_dir     = local.vault_install_dir
      vault_product_version = matrix.artifact_source == "local" ? step.get_local_metadata.version : var.vault_product_version
      vault_revision        = matrix.artifact_source == "local" ? step.get_local_metadata.revision : var.vault_revision
      vault_build_date      = matrix.artifact_source == "local" ? step.get_local_metadata.build_date : var.vault_build_date
      vault_root_token      = step.create_vault_cluster.vault_root_token
    }
  }

  step "get_updated_vault_cluster_ips" {
    module = module.vault_get_cluster_ips
    depends_on = [
      step.create_vault_cluster,
      step.upgrade_vault
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances   = step.create_vault_cluster.vault_instances
      vault_install_dir = local.vault_install_dir
      vault_root_token  = step.create_vault_cluster.vault_root_token
    }
  }

  step "verify_vault_unsealed" {
    module = module.vault_verify_unsealed
    depends_on = [
      step.create_vault_cluster,
      step.get_updated_vault_cluster_ips,
      step.upgrade_vault,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances   = step.create_vault_cluster.vault_instances
      vault_install_dir = local.vault_install_dir
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
      node_public_ips   = step.get_updated_vault_cluster_ips.follower_public_ips
      vault_install_dir = local.vault_install_dir
    }
  }

  step "verify_raft_auto_join_voter" {
    skip_step = matrix.backend != "raft"
    module    = module.vault_verify_raft_auto_join_voter
    depends_on = [
      step.create_backend_cluster,
      step.upgrade_vault,
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
}
