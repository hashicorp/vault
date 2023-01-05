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
      file_name = abspath(joinpath(path.root, "./support/vault.hclic"))
    }
  }

  step "create_primary_backend_cluster" {
    module     = "backend_${matrix.primary_backend}"
    depends_on = [step.create_vpc]

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

  step "create_vault_primary_cluster" {
    module = module.vault_cluster
    depends_on = [
      step.create_primary_backend_cluster,
      step.build_vault,
    ]
    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      ami_id             = step.create_vpc.ami_ids[matrix.distro][matrix.arch]
      common_tags        = local.tags
      consul_cluster_tag = step.create_primary_backend_cluster.consul_cluster_tag
      consul_release = matrix.primary_backend == "consul" ? {
        edition = var.backend_edition
        version = matrix.consul_version
      } : null
      dependencies_to_install   = local.dependencies_to_install
      instance_type             = local.vault_instance_type
      kms_key_arn               = step.create_vpc.kms_key_arn
      storage_backend           = matrix.primary_backend
      unseal_method             = matrix.primary_seal
      vault_local_artifact_path = local.bundle_path
      vault_install_dir         = local.vault_install_dir
      vault_artifactory_release = matrix.artifact_source == "artifactory" ? step.build_vault.vault_artifactory_release : null
      vault_environment = {
        VAULT_LOG_LEVEL = "debug"
      }
      vault_license = step.read_license.license
      vpc_id        = step.create_vpc.vpc_id
    }
  }

  step "create_secondary_backend_cluster" {
    module     = "backend_${matrix.secondary_backend}"
    depends_on = [step.create_vpc]

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

  step "create_vault_secondary_cluster" {
    module = module.vault_cluster
    depends_on = [
      step.create_secondary_backend_cluster,
      step.build_vault,
    ]
    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      ami_id             = step.create_vpc.ami_ids[matrix.distro][matrix.arch]
      common_tags        = local.tags
      consul_cluster_tag = step.create_secondary_backend_cluster.consul_cluster_tag
      consul_release = matrix.secondary_backend == "consul" ? {
        edition = var.backend_edition
        version = matrix.consul_version
      } : null
      dependencies_to_install   = local.dependencies_to_install
      instance_type             = local.vault_instance_type
      kms_key_arn               = step.create_vpc.kms_key_arn
      storage_backend           = matrix.secondary_backend
      unseal_method             = matrix.secondary_seal
      vault_local_artifact_path = local.bundle_path
      vault_install_dir         = local.vault_install_dir
      vault_artifactory_release = matrix.artifact_source == "artifactory" ? step.build_vault.vault_artifactory_release : null
      vault_environment = {
        VAULT_LOG_LEVEL = "debug"
      }
      vault_license = step.read_license.license
      vpc_id        = step.create_vpc.vpc_id
    }
  }

  step "verify_vault_primary_unsealed" {
    module = module.vault_verify_unsealed
    depends_on = [
      step.create_vault_primary_cluster
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances   = step.create_vault_primary_cluster.vault_instances
      vault_install_dir = local.vault_install_dir
    }
  }

  step "verify_vault_secondary_unsealed" {
    module = module.vault_verify_unsealed
    depends_on = [
      step.create_vault_secondary_cluster
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances   = step.create_vault_secondary_cluster.vault_instances
      vault_install_dir = local.vault_install_dir
    }
  }

  step "get_primary_cluster_ips" {
    module     = module.vault_get_cluster_ips
    depends_on = [step.verify_vault_primary_unsealed]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances   = step.create_vault_primary_cluster.vault_instances
      vault_install_dir = local.vault_install_dir
      vault_root_token  = step.create_vault_primary_cluster.vault_root_token
    }
  }

  step "get_secondary_cluster_ips" {
    module     = module.vault_get_cluster_ips
    depends_on = [step.verify_vault_secondary_unsealed]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances   = step.create_vault_secondary_cluster.vault_instances
      vault_install_dir = local.vault_install_dir
      vault_root_token  = step.create_vault_secondary_cluster.vault_root_token
    }
  }

  step "verify_vault_primary_write_data" {
    module     = module.vault_verify_write_data
    depends_on = [step.get_primary_cluster_ips]


    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      leader_public_ip  = step.get_primary_cluster_ips.leader_public_ip
      leader_private_ip = step.get_primary_cluster_ips.leader_private_ip
      vault_instances   = step.create_vault_primary_cluster.vault_instances
      vault_install_dir = local.vault_install_dir
      vault_root_token  = step.create_vault_primary_cluster.vault_root_token
    }
  }

  step "configure_performance_replication_primary" {
    module = module.vault_setup_perf_primary
    depends_on = [
      step.get_primary_cluster_ips,
      step.verify_vault_primary_write_data
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      primary_leader_public_ip  = step.get_primary_cluster_ips.leader_public_ip
      primary_leader_private_ip = step.get_primary_cluster_ips.leader_private_ip
      vault_install_dir         = local.vault_install_dir
      vault_root_token          = step.create_vault_primary_cluster.vault_root_token
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
      vault_root_token         = step.create_vault_primary_cluster.vault_root_token
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
      vault_root_token            = step.create_vault_secondary_cluster.vault_root_token
      wrapping_token              = step.generate_secondary_token.secondary_token
    }
  }

  step "unseal_secondary_followers" {
    module = module.vault_unseal_nodes
    depends_on = [
      step.create_vault_primary_cluster,
      step.create_vault_secondary_cluster,
      step.get_secondary_cluster_ips,
      step.configure_performance_replication_secondary
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      follower_public_ips = step.get_secondary_cluster_ips.follower_public_ips
      vault_install_dir   = local.vault_install_dir
      vault_unseal_keys   = matrix.primary_seal == "shamir" ? step.create_vault_primary_cluster.vault_unseal_keys_hex : step.create_vault_primary_cluster.vault_recovery_keys_hex
      vault_seal_type     = matrix.primary_seal == "shamir" ? matrix.primary_seal : matrix.secondary_seal
    }
  }

  step "verify_vault_secondary_unsealed_after_replication" {
    module = module.vault_verify_unsealed
    depends_on = [
      step.unseal_secondary_followers
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances   = step.create_vault_secondary_cluster.vault_instances
      vault_install_dir = local.vault_install_dir
    }
  }

  step "verify_performance_replication" {
    module     = module.vault_verify_performance_replication
    depends_on = [step.verify_vault_secondary_unsealed_after_replication]

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
      step.verify_vault_primary_write_data
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      node_public_ips   = step.get_secondary_cluster_ips.follower_public_ips
      vault_install_dir = local.vault_install_dir
    }
  }

  step "add_primary_cluster_nodes" {
    module = module.vault_cluster
    depends_on = [
      step.create_vpc,
      step.create_primary_backend_cluster,
      step.create_vault_primary_cluster,
      step.verify_replicated_data
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      ami_id             = step.create_vpc.ami_ids[matrix.distro][matrix.arch]
      common_tags        = local.tags
      consul_cluster_tag = step.create_primary_backend_cluster.consul_cluster_tag
      consul_release = matrix.primary_backend == "consul" ? {
        edition = var.backend_edition
        version = matrix.consul_version
      } : null
      dependencies_to_install   = local.dependencies_to_install
      instance_type             = local.vault_instance_type
      kms_key_arn               = step.create_vpc.kms_key_arn
      storage_backend           = matrix.primary_backend
      unseal_method             = matrix.primary_seal
      vault_cluster_tag         = step.create_vault_primary_cluster.vault_cluster_tag
      vault_init                = false
      vault_license             = step.read_license.license
      vault_local_artifact_path = local.bundle_path
      vault_install_dir         = local.vault_install_dir
      vault_artifactory_release = matrix.artifact_source == "artifactory" ? step.build_vault.vault_artifactory_release : null
      vault_environment = {
        VAULT_LOG_LEVEL = "debug"
      }
      vault_node_prefix         = "newprimary_node"
      vault_root_token          = step.create_vault_primary_cluster.vault_root_token
      vault_unseal_when_no_init = matrix.primary_seal == "shamir"
      vault_unseal_keys         = matrix.primary_seal == "shamir" ? step.create_vault_primary_cluster.vault_unseal_keys_hex : null
      vpc_id                    = step.create_vpc.vpc_id
    }
  }

  step "verify_add_node_unsealed" {
    module     = module.vault_verify_unsealed
    depends_on = [step.add_primary_cluster_nodes]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances   = step.add_primary_cluster_nodes.vault_instances
      vault_install_dir = local.vault_install_dir
    }
  }

  step "verify_raft_auto_join_voter" {
    skip_step = matrix.primary_backend != "raft"
    module    = module.vault_verify_raft_auto_join_voter
    depends_on = [
      step.add_primary_cluster_nodes,
      step.create_vault_primary_cluster,
      step.verify_add_node_unsealed
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances   = step.add_primary_cluster_nodes.vault_instances
      vault_install_dir = local.vault_install_dir
      vault_root_token  = step.create_vault_primary_cluster.vault_root_token
    }
  }

  step "remove_primary_follower_1" {
    module = module.shutdown_node
    depends_on = [
      step.get_primary_cluster_ips,
      step.verify_add_node_unsealed
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
      step.add_primary_cluster_nodes,
      step.remove_primary_follower_1,
      step.remove_primary_leader
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances       = step.create_vault_primary_cluster.vault_instances
      vault_install_dir     = local.vault_install_dir
      added_vault_instances = step.add_primary_cluster_nodes.vault_instances
      vault_root_token      = step.create_vault_primary_cluster.vault_root_token
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

  output "vault_primary_newnode_pub_ip" {
    description = "The Vault added new node on primary cluster public IP"
    value       = step.add_primary_cluster_nodes.instance_public_ips
  }

  output "vault_primary_newnode_priv_ip" {
    description = "The Vault added new node on primary cluster private IP"
    value       = step.add_primary_cluster_nodes.instance_private_ips
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

  output "vault_primary_performance_replication_status" {
    description = "The Vault primary cluster performance replication status"
    value       = step.verify_performance_replication.primary_replication_status
  }

  output "vault_replication_known_primary_cluster_addrs" {
    description = "The Vault secondary cluster performance replication status"
    value       = step.verify_performance_replication.known_primary_cluster_addrs
  }

  output "vault_secondary_performance_replication_status" {
    description = "The Vault secondary cluster performance replication status"
    value       = step.verify_performance_replication.secondary_replication_status
  }

  output "vault_primary_updated_performance_replication_status" {
    description = "The Vault updated primary cluster performance replication status"
    value       = step.verify_updated_performance_replication.primary_replication_status
  }

  output "verify_secondary_updated_performance_replication_status" {
    description = "The Vault updated secondary cluster performance replication status"
    value       = step.verify_updated_performance_replication.secondary_replication_status
  }
}
