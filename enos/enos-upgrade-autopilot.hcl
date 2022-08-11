// Enos scenario to test the Vault autopilot upgrade available in Vault Enterprise versions >=
variable "tags" {
  description = "Tags to add to AWS resources"
  type        = map(string)
  default     = null
}

terraform_cli "default" {
  plugin_cache_dir = var.terraform_plugin_cache_dir != null ? abspath(var.terraform_plugin_cache_dir) : null

  credentials "app.terraform.io" {
    token = var.tfc_api_token
  }
}

terraform "default" {
  required_version = ">= 1.0.0"

  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }

    aws = {
      source = "hashicorp/aws"
    }
  }
}

scenario "upgrade_autopilot" {
  matrix {
    distro  = ["ubuntu", "rhel"]
    builder = ["local", "crt"]
    arch    = ["amd64", "arm64"]
    edition = ["ent", "ent.hsm"]
  }

  terraform_cli = terraform_cli.default
  terraform     = terraform.default
  providers = [
    provider.aws.default,
    provider.enos.ubuntu,
    provider.enos.rhel
  ]

  locals {
    bundle_path = abspath(var.vault_bundle_path)
    binary_path = abspath(var.vault_local_binary_path)
    enos_provider = {
      rhel   = provider.enos.rhel
      ubuntu = provider.enos.ubuntu
    }
    enos_transport_user = {
      rhel   = "ec2-user"
      ubuntu = "ubuntu"
    }
    default_instances_types = {
      amd64 = "t3a.small"
      arm64 = "t4g.small"
    }

    supported_editions_by_arch = {
      amd64 = ["ent"]
    }
    vault_initial_version = "1.11.0"
  }

  step "build_vault" {
    module = matrix.builder == "crt" ? module.build_crt : module.build_local

    variables {
      bundle_path               = local.bundle_path
      local_vault_artifact_path = local.binary_path
    }
  }

  step "find_azs" {
    module = module.az_finder

    variables {
      instance_type = [
        var.vault_instance_type
      ]
    }
  }

  step "create_vpc" {
    module = module.create_vpc

    variables {
      ami_architectures = [matrix.arch]
    }
  }

  step "read_license" {
    module = module.read_license

    variables {
      file_name = abspath(joinpath(path.root, "./support/vault.hclic"))
    }
  }

  step "vault_with_raft" {
    module     = module.vault_cluster
    depends_on = [step.create_vpc]
    providers = {
      enos = local.enos_provider[matrix.distro]
    }
    variables {
      ami_id                  = step.create_vpc.ami_ids[matrix.distro][matrix.arch]
      common_tags             = var.tags
      dependencies_to_install = ["jq"]
      vpc_id                  = step.create_vpc.vpc_id
      kms_key_arn             = step.create_vpc.kms_key_arn
      storage_backend         = "raft"
      instance_type           = local.default_instances_types[matrix.arch]
      enos_transport_user     = local.enos_transport_user[matrix.distro]
      consul_cluster_tag      = null
      vault_release = {
        version = local.vault_initial_version
        edition = matrix.edition
      }
      vault_license = matrix.edition != "oss" ? semverconstraint(local.vault_initial_version, "> 1.8.0") ? step.read_license.license : null : null
      storage_backend_addl_config = {
        autopilot_upgrade_version = local.vault_initial_version
      }
    }
  }

  step "get_vault_version" {
    module = module.get_vault_version
  }

  step "vault_autopilot_upgrade_storageconfig" {
    module     = module.vault_autopilot_upgrade_storageconfig
    depends_on = [step.get_vault_version]

    variables {
      vault_product_version = step.get_vault_version.vault_product_version
    }
  }

  step "upgrade_with_autopilot" {
    module     = module.vault_cluster
    depends_on = [step.vault_with_raft, step.vault_autopilot_upgrade_storageconfig]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      ami_id                      = step.create_vpc.ami_ids[matrix.distro][matrix.arch]
      common_tags                 = var.tags
      dependencies_to_install     = ["jq"]
      consul_cluster_tag          = null
      enos_transport_user         = local.enos_provider[matrix.distro].config.attrs.transport.ssh.user
      instance_count              = var.vault_instance_count
      instance_type               = local.default_instances_types[matrix.arch]
      kms_key_arn                 = step.create_vpc.kms_key_arn
      storage_backend             = "raft"
      vault_install_dir           = var.vault_install_dir
      vault_local_artifact_path   = local.bundle_path
      vault_license               = matrix.edition != "oss" ? step.read_license.license : null
      vpc_id                      = step.create_vpc.vpc_id
      storage_backend_addl_config = step.vault_autopilot_upgrade_storageconfig.storage_addl_config
      vault_cluster_tag           = step.vault_with_raft.vault_cluster_tag
      vault_root_token            = step.vault_with_raft.vault_root_token
      vault_node_prefix           = "upgrade_node"
      vault_init                  = "false"
    }
  }

  step "verify_autopilot" {
    module     = module.verify_autopilot
    depends_on = [step.upgrade_with_autopilot]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_install_dir               = var.vault_install_dir
      vault_autopilot_upgrade_version = step.get_vault_version.vault_product_version
      vault_autopilot_upgrade_status  = "await-server-removal"
      vault_token                     = step.upgrade_with_autopilot.vault_root_token
    }
  }

  output "vault_with_raft_instance_ids" {
    description = "The Vault cluster instance IDs"
    value       = step.vault_with_raft.instance_ids
  }

  output "vault_with_raft_pub_ips" {
    description = "The Vault cluster public IPs"
    value       = step.vault_with_raft.instance_public_ips
  }

  output "vault_with_raft_priv_ips" {
    description = "The Vault cluster private IPs"
    value       = step.vault_with_raft.instance_private_ips
  }

  output "vault_with_raft_key_id" {
    description = "The Vault cluster Key ID"
    value       = step.vault_with_raft.key_id
  }

  output "vault_with_raft_root_token" {
    description = "The Vault cluster root token"
    value       = step.vault_with_raft.vault_root_token
  }

  output "vault_cluster_upgrade_instance_ids" {
    description = "The Vault cluster instance IDs"
    value       = step.upgrade_with_autopilot.instance_ids
  }

  output "vault_cluster_upgrade_pub_ips" {
    description = "The Vault cluster public IPs"
    value       = step.upgrade_with_autopilot.instance_public_ips
  }

  output "vault_cluster_upgrade_priv_ips" {
    description = "The Vault cluster private IPs"
    value       = step.upgrade_with_autopilot.instance_private_ips
  }

  output "vault_cluster_upgrade_key_id" {
    description = "The Vault cluster Key ID"
    value       = step.upgrade_with_autopilot.key_id
  }

  output "vault_cluster_upgrade_root_token" {
    description = "The Vault cluster root token"
    value       = step.upgrade_with_autopilot.vault_root_token
  }

  output "vault_cluster_tag" {
    description = "The Vault cluster tag"
    value       = step.upgrade_with_autopilot.vault_cluster_tag
  }
}
