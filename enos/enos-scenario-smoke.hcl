scenario "smoke" {
  matrix {
    arch           = ["amd64", "arm64"]
    backend        = ["consul", "raft"]
    builder        = ["local", "crt"]
    consul_version = ["1.12.3", "1.11.7", "1.10.12"]
    distro         = ["ubuntu", "rhel"]
    edition        = ["oss", "ent"]
    seal           = ["awskms", "shamir"]
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
      "oss" = ["ui"]
      "ent" = ["enterprise", "ent"]
    }
    bundle_path             = abspath(var.vault_bundle_path)
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
  }

  step "build_vault" {
    module = matrix.builder == "crt" ? module.build_crt : module.build_local

    variables {
      build_tags  = var.vault_local_build_tags != null ? var.vault_local_build_tags : local.build_tags[matrix.edition]
      bundle_path = local.bundle_path
      goarch      = matrix.arch
      goos        = "linux"
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
      ami_architectures  = [matrix.arch]
      availability_zones = step.find_azs.availability_zones
      common_tags        = local.tags
    }
  }

  step "read_license" {
    skip_step = matrix.edition == "oss"
    module    = module.read_license

    variables {
      file_name = abspath(joinpath(path.root, "./support/vault.hclic"))
    }
  }

  step "create_backend_cluster" {
    module = "backend_${matrix.backend}"
    depends_on = [
      step.create_vpc,
      step.build_vault,
    ]

    providers = {
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

  step "create_vault_cluster" {
    module = module.vault_cluster
    depends_on = [
      step.create_vpc,
      step.create_backend_cluster,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      ami_id                    = step.create_vpc.ami_ids[matrix.distro][matrix.arch]
      common_tags               = local.tags
      consul_cluster_tag        = step.create_backend_cluster.consul_cluster_tag
      dependencies_to_install   = local.dependencies_to_install
      instance_type             = local.vault_instance_type
      kms_key_arn               = step.create_vpc.kms_key_arn
      storage_backend           = matrix.backend
      unseal_method             = matrix.seal
      vault_local_artifact_path = local.bundle_path
      vault_license             = matrix.edition != "oss" ? step.read_license.license : null
      vpc_id                    = step.create_vpc.vpc_id
    }
  }

  step "verify_vault_unsealed" {
    module = module.vault_verify_unsealed
    depends_on = [
      step.create_vault_cluster,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances  = step.create_vault_cluster.vault_instances
      vault_root_token = step.create_vault_cluster.vault_root_token
    }
  }

  step "verify_raft_auto_join_voter" {
    skip_step = matrix.backend != "raft"
    module    = module.vault_verify_raft_auto_join_voter
    depends_on = [
      step.create_vault_cluster,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances  = step.create_vault_cluster.vault_instances
      vault_root_token = step.create_vault_cluster.vault_root_token
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
