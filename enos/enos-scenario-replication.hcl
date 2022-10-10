scenario "replication" {
  matrix {
    arch           = ["amd64", "arm64"]
    distro         = ["ubuntu", "rhel"]
    backend        = ["raft", "consul"]
    consul_version = ["1.12.3", "1.11.7", "1.10.12"]
    edition        = ["ent"]
    seal           = ["awskms", "shamir"]
  }

  terraform_cli = terraform_cli.default
  terraform     = terraform.default
  providers = [
    provider.aws.west1,
    provider.aws.west2,
    provider.enos.ubuntu,
    provider.enos.rhel
  ]

  locals {
    artifact_path           = var.artifact_path
    dependencies_to_install = ["jq"]
    enos_provider = {
      rhel   = provider.enos.rhel
      ubuntu = provider.enos.ubuntu
    }
    install_artifactory_artifact = true
    tags = merge({
      "Project Name" : var.project_name
      "Project" : "Enos",
      "Environment" : "ci"
    }, var.tags)
    // install_dev_artifact         = local.artifact_path != null
    // install_artifactory_artifact = local.revision != null && local.artifact_path == null
    // install_release_artifact     = local.revision == null && local.artifact_path == null
    vault_instance_types = {
      amd64 = "t3a.small"
      arm64 = "t4g.small"
    }
    vault_instance_type = coalesce(var.vault_instance_type, local.vault_instance_types[matrix.arch])
  }

  step "find_azs" {
    module = module.az_finder
    providers = {
      aws = provider.aws.west1
    }
    variables {
      instance_type = [
        local.vault_instance_type
      ]
    }
  }

  step "find_secondary_azs" {
    module = module.az_finder
    providers = {
      aws = provider.aws.west2
    }
    variables {
      instance_type = [
        local.vault_instance_type
      ]
    }
  }

  step "create_vpc" {
    module     = module.create_vpc
    depends_on = [step.find_azs]
    providers = {
      aws = provider.aws.west1
    }

    variables {
      ami_architectures  = [matrix.arch]
      availability_zones = step.find_azs.availability_zones
      common_tags        = local.tags
    }
  }

  step "create_vpc_2" {
    module     = module.create_vpc
    depends_on = [step.find_secondary_azs]
    providers = {
      aws = provider.aws.west2
    }

    variables {
      ami_architectures  = [matrix.arch]
      availability_zones = step.find_secondary_azs.availability_zones
      common_tags        = local.tags
    }
  }

  step "read_license" {
    module = module.read_license

    variables {
      file_name = abspath(joinpath(path.root, "./support/vault.hclic"))
    }
  }

  step "fetch_vault_artifact" {
    module = module.build_artifactory

    variables {
      artifactory_host      = var.artifactory_host
      artifactory_repo      = var.artifactory_repo
      artifactory_username  = var.artifactory_username
      artifactory_token     = var.artifactory_token
      arch                  = matrix.arch
      vault_product_version = var.vault_product_version
      artifact_type         = "bundle"
      distro                = matrix.distro
      edition               = matrix.edition
      instance_type         = local.vault_instance_type
      revision              = var.vault_revision
    }
  }

  step "create_primary_backend_cluster" {
    module     = "backend_${matrix.backend}"
    depends_on = [step.create_vpc]

    providers = {
      aws  = provider.aws.west1
      enos = provider.enos.ubuntu
    }

    variables {
      ami_id      = step.create_vpc.ami_ids["ubuntu"][matrix.arch]
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

  step "create_vault_primary_cluster" {
    module = module.vault_cluster
    depends_on = [
      step.create_primary_backend_cluster,
      step.fetch_vault_artifact,
    ]
    providers = {
      aws  = provider.aws.west1
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      ami_id                    = step.create_vpc.ami_ids[matrix.distro][matrix.arch]
      common_tags               = local.tags
      consul_cluster_tag        = step.create_primary_backend_cluster.consul_cluster_tag
      dependencies_to_install   = local.dependencies_to_install
      instance_type             = local.vault_instance_type
      kms_key_arn               = step.create_vpc.kms_key_arn
      storage_backend           = matrix.backend
      unseal_method             = matrix.seal
      vault_artifactory_release = local.install_artifactory_artifact ? step.fetch_vault_artifact.vault_artifactory_release : null
      vault_license             = step.read_license.license
      vpc_id                    = step.create_vpc.vpc_id
    }
  }

  step "create_secondary_backend_cluster" {
    module     = "backend_${matrix.backend}"
    depends_on = [step.create_vpc_2]

    providers = {
      aws  = provider.aws.west2
      enos = provider.enos.ubuntu
    }

    variables {
      ami_id      = step.create_vpc_2.ami_ids["ubuntu"][matrix.arch]
      common_tags = local.tags
      consul_release = {
        edition = var.backend_edition
        version = matrix.consul_version
      }
      instance_type = var.backend_instance_type
      kms_key_arn   = step.create_vpc_2.kms_key_arn
      vpc_id        = step.create_vpc_2.vpc_id
    }
  }

  step "create_vault_secondary_cluster" {
    module = module.vault_cluster
    depends_on = [
      step.create_secondary_backend_cluster,
      step.fetch_vault_artifact,
    ]
    providers = {
      aws  = provider.aws.west2
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      ami_id                    = step.create_vpc_2.ami_ids[matrix.distro][matrix.arch]
      common_tags               = local.tags
      consul_cluster_tag        = step.create_secondary_backend_cluster.consul_cluster_tag
      dependencies_to_install   = local.dependencies_to_install
      instance_type             = local.vault_instance_type
      kms_key_arn               = step.create_vpc_2.kms_key_arn
      storage_backend           = matrix.backend
      unseal_method             = matrix.seal
      vault_artifactory_release = local.install_artifactory_artifact ? step.fetch_vault_artifact.vault_artifactory_release : null
      vault_license             = step.read_license.license
      vpc_id                    = step.create_vpc_2.vpc_id
    }
  }

  step "verify_vault_primary_unsealed" {
    module = module.vault_verify_unsealed
    depends_on = [
      step.create_vault_primary_cluster
    ]

    providers = {
      aws  = provider.aws.west1
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances  = step.create_vault_primary_cluster.vault_instances
      vault_root_token = step.create_vault_primary_cluster.vault_root_token
    }
  }

  step "verify_vault_secondary_unsealed" {
    module = module.vault_verify_unsealed
    depends_on = [
      step.create_vault_secondary_cluster
    ]

    providers = {
      aws  = provider.aws.west2
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances  = step.create_vault_secondary_cluster.vault_instances
      vault_root_token = step.create_vault_secondary_cluster.vault_root_token
    }
  }

  output "vault_primary_cluster_instance_ids" {
    description = "The Vault primary cluster instance IDs"
    value       = step.create_vault_primary_cluster.instance_ids
  }

  output "vault_primary_cluster_pub_ips" {
    description = "The Vault primary cluster public IPs"
    value       = step.create_vault_primary_cluster.instance_public_ips
  }

  output "vault_primary_cluster_priv_ips" {
    description = "The Vault primary cluster private IPs"
    value       = step.create_vault_primary_cluster.instance_private_ips
  }

  output "vault_primary_cluster_key_id" {
    description = "The Vault primary cluster Key ID"
    value       = step.create_vault_primary_cluster.key_id
  }

  output "vault_primary_cluster_root_token" {
    description = "The Vault primary cluster root token"
    value       = step.create_vault_primary_cluster.vault_root_token
  }

  output "vault_primary_cluster_unseal_keys_b64" {
    description = "The Vault primary cluster unseal keys"
    value       = step.create_vault_primary_cluster.vault_unseal_keys_b64
  }

  output "vault_primary_cluster_unseal_keys_hex" {
    description = "The Vault primary cluster unseal keys hex"
    value       = step.create_vault_primary_cluster.vault_unseal_keys_hex
  }

  output "vault_primary_cluster_tag" {
    description = "The Vault primary cluster tag"
    value       = step.create_vault_primary_cluster.vault_cluster_tag
  }

  output "vault_secondary_cluster_instance_ids" {
    description = "The Vault secondary cluster instance IDs"
    value       = step.create_vault_secondary_cluster.instance_ids
  }

  output "vault_secondary_cluster_pub_ips" {
    description = "The Vault secondary cluster public IPs"
    value       = step.create_vault_secondary_cluster.instance_public_ips
  }

  output "vault_secondary_cluster_priv_ips" {
    description = "The Vault secondary cluster private IPs"
    value       = step.create_vault_secondary_cluster.instance_private_ips
  }

  output "vault_secondary_cluster_tag" {
    description = "The Vault secondary cluster tag"
    value       = step.create_vault_secondary_cluster.vault_cluster_tag
  }

  output "vault_secondary_cluster_key_id" {
    description = "The Vault secondary cluster Key ID"
    value       = step.create_vault_secondary_cluster.key_id
  }

  output "vault_secondary_cluster_root_token" {
    description = "The Vault secondary cluster root token"
    value       = step.create_vault_secondary_cluster.vault_root_token
  }

  output "vault_secondary_cluster_unseal_keys_b64" {
    description = "The Vault secondary cluster unseal keys"
    value       = step.create_vault_secondary_cluster.vault_unseal_keys_b64
  }

  output "vault_secondary_cluster_unseal_keys_hex" {
    description = "The Vault secondary cluster unseal keys hex"
    value       = step.create_vault_secondary_cluster.vault_unseal_keys_hex
  }
}
