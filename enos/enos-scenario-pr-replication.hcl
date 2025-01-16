// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

scenario "pr_replication" {
  description = <<-EOF
    The PR replication scenario configures performance replication between two Vault clusters and
    verifies behavior and failure tolerance. The build can be a local branch, any CRT built Vault
    Enterprise artifact saved to the local machine, or any CRT built Vault Enterprise artifact in
    the stable channel in Artifactory.

    The scenario deploys two Vault Enterprise clusters and establishes performance replication
    between the primary cluster and the performance replication secondary cluster. Next, it simulates
    a catastrophic failure event whereby the primary leader and a primary follower are ungracefully
    removed from the cluster while running. This forces a leader election in the primary cluster
    and requires the secondary cluster to recover replication and establish replication to the new
    primary leader. The scenario also performs standard baseline verification that is not specific
    to performance replication.

    # How to run this scenario

    For general instructions on running a scenario, refer to the Enos docs: https://eng-handbook.hashicorp.services/internal-tools/enos/running-a-scenario/
    For troubleshooting tips and common errors, see https://eng-handbook.hashicorp.services/internal-tools/enos/troubleshooting/.

    Variables required for all scenario variants:
      - aws_ssh_private_key_path (more info about AWS SSH keypairs: https://eng-handbook.hashicorp.services/internal-tools/enos/getting-started/#set-your-aws-key-pair-name-and-private-key)
      - aws_ssh_keypair_name
      - vault_build_date*
      - vault_product_version
      - vault_revision*
    
    * If you don't already know what build date and revision you should be using, see
    https://eng-handbook.hashicorp.services/internal-tools/enos/troubleshooting/#execution-error-expected-vs-got-for-vault-versioneditionrevisionbuild-date.
  
    Variables required for some scenario variants:
      - artifactory_username (if using `artifact_source:artifactory` in your filter)
      - artifactory_token (if using `artifact_source:artifactory` in your filter)
      - aws_region (if different from the default value in enos-variables.hcl)
      - consul_license_path (if using an ENT edition of Consul)
      - distro_version_<distro> (if different from the default version for your target
      distro. See supported distros and default versions in the distro_version_<distro>
      definitions in enos-variables.hcl)
      - vault_artifact_path (the path to where you have a Vault artifact already downloaded,
      if using `artifact_source:crt` in your filter)
      - vault_license_path (if using an ENT edition of Vault)
  EOF

  matrix {
    arch              = global.archs
    artifact_source   = global.artifact_sources
    artifact_type     = global.artifact_types
    config_mode       = global.config_modes
    consul_edition    = global.consul_editions
    consul_version    = global.consul_versions
    distro            = global.distros
    edition           = global.enterprise_editions
    ip_version        = global.ip_versions
    primary_backend   = global.backends
    primary_seal      = global.seals
    secondary_backend = global.backends
    secondary_seal    = global.seals

    // Our local builder always creates bundles
    exclude {
      artifact_source = ["local"]
      artifact_type   = ["package"]
    }

    // PKCS#11 can only be used on ent.hsm and ent.hsm.fips1402.
    exclude {
      primary_seal = ["pkcs11"]
      edition      = [for e in matrix.edition : e if !strcontains(e, "hsm")]
    }

    exclude {
      secondary_seal = ["pkcs11"]
      edition        = [for e in matrix.edition : e if !strcontains(e, "hsm")]
    }

    // softhsm packages not available for leap/sles.
    exclude {
      primary_seal = ["pkcs11"]
      distro       = ["leap", "sles"]
    }

    exclude {
      secondary_seal = ["pkcs11"]
      distro         = ["leap", "sles"]
    }

    // Testing in IPV6 mode is currently implemented for integrated Raft storage only
    exclude {
      ip_version      = ["6"]
      primary_backend = ["consul"]
    }

    exclude {
      ip_version        = ["6"]
      secondary_backend = ["consul"]
    }
  }

  terraform_cli = terraform_cli.default
  terraform     = terraform.default
  providers = [
    provider.aws.default,
    provider.enos.ec2_user,
    provider.enos.ubuntu
  ]

  locals {
    artifact_path = matrix.artifact_source != "artifactory" ? abspath(var.vault_artifact_path) : null
    enos_provider = {
      amzn   = provider.enos.ec2_user
      leap   = provider.enos.ec2_user
      rhel   = provider.enos.ec2_user
      sles   = provider.enos.ec2_user
      ubuntu = provider.enos.ubuntu
    }
    manage_service    = matrix.artifact_type == "bundle"
    vault_install_dir = matrix.artifact_type == "bundle" ? var.vault_install_dir : global.vault_install_dir[matrix.artifact_type]
  }

  step "build_vault" {
    description = global.description.build_vault
    module      = "build_${matrix.artifact_source}"

    variables {
      build_tags           = var.vault_local_build_tags != null ? var.vault_local_build_tags : global.build_tags[matrix.edition]
      artifact_path        = local.artifact_path
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
    description = global.description.ec2_info
    module      = module.ec2_info
  }

  step "create_vpc" {
    description = global.description.create_vpc
    module      = module.create_vpc

    variables {
      common_tags = global.tags
      ip_version  = matrix.ip_version
    }
  }

  // This step reads the contents of the backend license if we're using a Consul backend and
  // an "ent" Consul edition.
  step "read_backend_license" {
    description = global.description.read_backend_license
    skip_step   = (matrix.primary_backend == "raft" && matrix.secondary_backend == "raft") || matrix.consul_edition == "ce"
    module      = module.read_license

    variables {
      file_name = global.backend_license_path
    }
  }

  step "read_vault_license" {
    description = global.description.read_vault_license
    module      = module.read_license

    variables {
      file_name = abspath(joinpath(path.root, "./support/vault.hclic"))
    }
  }

  step "create_primary_seal_key" {
    description = global.description.create_seal_key
    module      = "seal_${matrix.primary_seal}"
    depends_on  = [step.create_vpc]

    providers = {
      enos = provider.enos.ubuntu
    }

    variables {
      cluster_id   = step.create_vpc.id
      cluster_meta = "primary"
      common_tags  = global.tags
    }
  }

  step "create_secondary_seal_key" {
    description = global.description.create_seal_key
    module      = "seal_${matrix.secondary_seal}"
    depends_on  = [step.create_vpc]

    providers = {
      enos = provider.enos.ubuntu
    }

    variables {
      cluster_id      = step.create_vpc.id
      cluster_meta    = "secondary"
      common_tags     = global.tags
      other_resources = step.create_primary_seal_key.resource_names
    }
  }

  // Create all of our instances for both primary and secondary clusters
  step "create_primary_cluster_targets" {
    description = global.description.create_vault_cluster_targets
    module      = module.target_ec2_instances
    depends_on = [
      step.create_vpc,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      ami_id          = step.ec2_info.ami_ids[matrix.arch][matrix.distro][global.distro_version[matrix.distro]]
      cluster_tag_key = global.vault_tag_key
      common_tags     = global.tags
      seal_key_names  = step.create_primary_seal_key.resource_names
      vpc_id          = step.create_vpc.id
    }
  }

  step "create_primary_cluster_backend_targets" {
    description = global.description.create_vault_cluster_targets
    module      = matrix.primary_backend == "consul" ? module.target_ec2_instances : module.target_ec2_shim
    depends_on = [
      step.create_vpc,
    ]

    providers = {
      enos = provider.enos.ubuntu
    }

    variables {
      ami_id          = step.ec2_info.ami_ids["arm64"]["ubuntu"][global.distro_version["ubuntu"]]
      cluster_tag_key = global.backend_tag_key
      common_tags     = global.tags
      seal_key_names  = step.create_primary_seal_key.resource_names
      vpc_id          = step.create_vpc.id
    }
  }

  step "create_primary_cluster_additional_targets" {
    description = global.description.create_vault_cluster_targets
    module      = module.target_ec2_instances
    depends_on = [
      step.create_vpc,
      step.create_primary_cluster_targets,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      ami_id          = step.ec2_info.ami_ids[matrix.arch][matrix.distro][global.distro_version[matrix.distro]]
      cluster_name    = step.create_primary_cluster_targets.cluster_name
      cluster_tag_key = global.vault_tag_key
      common_tags     = global.tags
      seal_key_names  = step.create_primary_seal_key.resource_names
      vpc_id          = step.create_vpc.id
    }
  }

  step "create_secondary_cluster_targets" {
    description = global.description.create_vault_cluster_targets
    module      = module.target_ec2_instances
    depends_on  = [step.create_vpc]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      ami_id          = step.ec2_info.ami_ids[matrix.arch][matrix.distro][global.distro_version[matrix.distro]]
      cluster_tag_key = global.vault_tag_key
      common_tags     = global.tags
      seal_key_names  = step.create_secondary_seal_key.resource_names
      vpc_id          = step.create_vpc.id
    }
  }

  step "create_secondary_cluster_backend_targets" {
    description = global.description.create_vault_cluster_targets
    module      = matrix.secondary_backend == "consul" ? module.target_ec2_instances : module.target_ec2_shim
    depends_on  = [step.create_vpc]

    providers = {
      enos = provider.enos.ubuntu
    }

    variables {
      ami_id          = step.ec2_info.ami_ids["arm64"]["ubuntu"][global.distro_version["ubuntu"]]
      cluster_tag_key = global.backend_tag_key
      common_tags     = global.tags
      seal_key_names  = step.create_secondary_seal_key.resource_names
      vpc_id          = step.create_vpc.id
    }
  }

  step "create_primary_backend_cluster" {
    description = global.description.create_backend_cluster
    module      = "backend_${matrix.primary_backend}"
    depends_on = [
      step.create_primary_cluster_backend_targets,
    ]

    providers = {
      enos = provider.enos.ubuntu
    }

    verifies = [
      // verified in modules
      quality.consul_autojoin_aws,
      quality.consul_config_file,
      quality.consul_ha_leader_election,
      quality.consul_service_start_server,
      // verified in enos_consul_start resource
      quality.consul_api_agent_host_read,
      quality.consul_api_health_node_read,
      quality.consul_api_operator_raft_config_read,
      quality.consul_cli_validate,
      quality.consul_health_state_passing_read_nodes_minimum,
      quality.consul_operator_raft_configuration_read_voters_minimum,
      quality.consul_service_systemd_notified,
      quality.consul_service_systemd_unit,
    ]

    variables {
      cluster_name    = step.create_primary_cluster_backend_targets.cluster_name
      cluster_tag_key = global.backend_tag_key
      hosts           = step.create_primary_cluster_backend_targets.hosts
      license         = (matrix.primary_backend == "consul" && matrix.consul_edition == "ent") ? step.read_backend_license.license : null
      release = {
        edition = matrix.consul_edition
        version = matrix.consul_version
      }
    }
  }

  step "create_primary_cluster" {
    description = global.description.create_vault_cluster
    module      = module.vault_cluster
    depends_on = [
      step.create_primary_backend_cluster,
      step.build_vault,
      step.create_primary_cluster_targets
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      // verified in modules
      quality.consul_service_start_client,
      quality.vault_artifact_bundle,
      quality.vault_artifact_deb,
      quality.vault_artifact_rpm,
      quality.vault_audit_log,
      quality.vault_audit_socket,
      quality.vault_audit_syslog,
      quality.vault_autojoin_aws,
      quality.vault_config_env_variables,
      quality.vault_config_file,
      quality.vault_config_log_level,
      quality.vault_init,
      quality.vault_license_required_ent,
      quality.vault_listener_ipv4,
      quality.vault_listener_ipv6,
      quality.vault_service_start,
      quality.vault_storage_backend_consul,
      quality.vault_storage_backend_raft,
      // verified in enos_vault_start resource
      quality.vault_api_sys_config_read,
      quality.vault_api_sys_ha_status_read,
      quality.vault_api_sys_health_read,
      quality.vault_api_sys_host_info_read,
      quality.vault_api_sys_replication_status_read,
      quality.vault_api_sys_seal_status_api_read_matches_sys_health,
      quality.vault_api_sys_storage_raft_autopilot_configuration_read,
      quality.vault_api_sys_storage_raft_autopilot_state_read,
      quality.vault_api_sys_storage_raft_configuration_read,
      quality.vault_cli_status_exit_code,
      quality.vault_service_systemd_notified,
      quality.vault_service_systemd_unit,
    ]

    variables {
      artifactory_release     = matrix.artifact_source == "artifactory" ? step.build_vault.vault_artifactory_release : null
      backend_cluster_name    = step.create_primary_cluster_backend_targets.cluster_name
      backend_cluster_tag_key = global.backend_tag_key
      config_mode             = matrix.config_mode
      consul_license          = (matrix.primary_backend == "consul" && matrix.consul_edition == "ent") ? step.read_backend_license.license : null
      cluster_name            = step.create_primary_cluster_targets.cluster_name
      consul_release = matrix.primary_backend == "consul" ? {
        edition = matrix.consul_edition
        version = matrix.consul_version
      } : null
      enable_audit_devices = var.vault_enable_audit_devices
      hosts                = step.create_primary_cluster_targets.hosts
      install_dir          = global.vault_install_dir[matrix.artifact_type]
      ip_version           = matrix.ip_version
      license              = matrix.edition != "ce" ? step.read_vault_license.license : null
      local_artifact_path  = local.artifact_path
      manage_service       = local.manage_service
      packages             = concat(global.packages, global.distro_packages[matrix.distro][global.distro_version[matrix.distro]])
      seal_attributes      = step.create_primary_seal_key.attributes
      seal_type            = matrix.primary_seal
      storage_backend      = matrix.primary_backend
    }
  }

  step "get_local_metadata" {
    skip_step = matrix.artifact_source != "local"
    module    = module.get_local_metadata
  }

  step "wait_for_primary_cluster_leader" {
    description = global.description.wait_for_cluster_to_have_leader
    module      = module.vault_wait_for_leader
    depends_on  = [step.create_primary_cluster]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_api_sys_leader_read,
      quality.vault_unseal_ha_leader_election,
    ]

    variables {
      hosts             = step.create_primary_cluster_targets.hosts
      ip_version        = matrix.ip_version
      timeout           = 120 // seconds
      vault_addr        = step.create_primary_cluster.api_addr_localhost
      vault_install_dir = local.vault_install_dir
      vault_root_token  = step.create_primary_cluster.root_token
    }
  }

  step "create_secondary_backend_cluster" {
    description = global.description.create_backend_cluster
    module      = "backend_${matrix.secondary_backend}"
    depends_on = [
      step.create_secondary_cluster_backend_targets
    ]

    providers = {
      enos = provider.enos.ubuntu
    }

    variables {
      cluster_name    = step.create_secondary_cluster_backend_targets.cluster_name
      cluster_tag_key = global.backend_tag_key
      hosts           = step.create_secondary_cluster_backend_targets.hosts
      license         = (matrix.secondary_backend == "consul" && matrix.consul_edition == "ent") ? step.read_backend_license.license : null
      release = {
        edition = matrix.consul_edition
        version = matrix.consul_version
      }
    }
  }

  step "create_secondary_cluster" {
    module = module.vault_cluster
    depends_on = [
      step.create_secondary_backend_cluster,
      step.build_vault,
      step.create_secondary_cluster_targets
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      // verified in modules
      quality.consul_autojoin_aws,
      quality.consul_config_file,
      quality.consul_ha_leader_election,
      quality.consul_service_start_client,
      // verified in enos_consul_start resource
      quality.consul_api_agent_host_read,
      quality.consul_api_health_node_read,
      quality.consul_api_operator_raft_config_read,
      quality.consul_health_state_passing_read_nodes_minimum,
      quality.consul_operator_raft_configuration_read_voters_minimum,
      quality.consul_service_systemd_notified,
      quality.consul_service_systemd_unit,
    ]

    variables {
      artifactory_release     = matrix.artifact_source == "artifactory" ? step.build_vault.vault_artifactory_release : null
      backend_cluster_name    = step.create_secondary_cluster_backend_targets.cluster_name
      backend_cluster_tag_key = global.backend_tag_key
      config_mode             = matrix.config_mode
      consul_license          = (matrix.secondary_backend == "consul" && matrix.consul_edition == "ent") ? step.read_backend_license.license : null
      cluster_name            = step.create_secondary_cluster_targets.cluster_name
      consul_release = matrix.secondary_backend == "consul" ? {
        edition = matrix.consul_edition
        version = matrix.consul_version
      } : null
      enable_audit_devices = var.vault_enable_audit_devices
      hosts                = step.create_secondary_cluster_targets.hosts
      install_dir          = global.vault_install_dir[matrix.artifact_type]
      ip_version           = matrix.ip_version
      license              = matrix.edition != "ce" ? step.read_vault_license.license : null
      local_artifact_path  = local.artifact_path
      manage_service       = local.manage_service
      packages             = concat(global.packages, global.distro_packages[matrix.distro][global.distro_version[matrix.distro]])
      seal_attributes      = step.create_secondary_seal_key.attributes
      seal_type            = matrix.secondary_seal
      storage_backend      = matrix.secondary_backend
    }
  }

  step "wait_for_secondary_cluster_leader" {
    description = global.description.wait_for_cluster_to_have_leader
    module      = module.vault_wait_for_leader
    depends_on  = [step.create_secondary_cluster]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_api_sys_leader_read,
      quality.vault_unseal_ha_leader_election,
    ]

    variables {
      hosts             = step.create_secondary_cluster_targets.hosts
      ip_version        = matrix.ip_version
      timeout           = 120 // seconds
      vault_addr        = step.create_secondary_cluster.api_addr_localhost
      vault_install_dir = local.vault_install_dir
      vault_root_token  = step.create_secondary_cluster.root_token
    }
  }

  step "verify_that_vault_primary_cluster_is_unsealed" {
    description = global.description.verify_vault_unsealed
    module      = module.vault_wait_for_cluster_unsealed
    depends_on = [
      step.create_primary_cluster,
      step.wait_for_primary_cluster_leader,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_auto_unseals_after_autopilot_upgrade,
      quality.vault_seal_awskms,
      quality.vault_seal_pkcs11,
      quality.vault_seal_shamir,
    ]

    variables {
      hosts             = step.create_primary_cluster_targets.hosts
      vault_addr        = step.create_primary_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
    }
  }

  step "verify_that_vault_secondary_cluster_is_unsealed" {
    description = global.description.verify_vault_unsealed
    module      = module.vault_wait_for_cluster_unsealed
    depends_on = [
      step.create_secondary_cluster,
      step.wait_for_secondary_cluster_leader,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_auto_unseals_after_autopilot_upgrade,
      quality.vault_seal_awskms,
      quality.vault_seal_pkcs11,
      quality.vault_seal_shamir,
    ]

    variables {
      hosts             = step.create_secondary_cluster_targets.hosts
      vault_addr        = step.create_secondary_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
    }
  }

  step "get_primary_cluster_ips" {
    description = global.description.get_vault_cluster_ip_addresses
    module      = module.vault_get_cluster_ips
    depends_on  = [step.verify_that_vault_primary_cluster_is_unsealed]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_api_sys_ha_status_read,
      quality.vault_api_sys_leader_read,
      quality.vault_cli_operator_members,
    ]

    variables {
      hosts             = step.create_primary_cluster_targets.hosts
      ip_version        = matrix.ip_version
      vault_addr        = step.create_primary_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      vault_root_token  = step.create_primary_cluster.root_token
    }
  }

  step "get_secondary_cluster_ips" {
    description = global.description.get_vault_cluster_ip_addresses
    module      = module.vault_get_cluster_ips
    depends_on  = [step.verify_that_vault_secondary_cluster_is_unsealed]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_api_sys_ha_status_read,
      quality.vault_api_sys_leader_read,
      quality.vault_cli_operator_members,
    ]

    variables {
      hosts             = step.create_secondary_cluster_targets.hosts
      ip_version        = matrix.ip_version
      vault_addr        = step.create_secondary_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      vault_root_token  = step.create_secondary_cluster.root_token
    }
  }

  step "verify_vault_version" {
    description = global.description.verify_vault_version
    module      = module.vault_verify_version
    depends_on  = [step.get_primary_cluster_ips]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_api_sys_version_history_keys,
      quality.vault_api_sys_version_history_key_info,
      quality.vault_version_build_date,
      quality.vault_version_edition,
      quality.vault_version_release,
    ]

    variables {
      hosts                 = step.create_primary_cluster_targets.hosts
      vault_addr            = step.create_primary_cluster.api_addr_localhost
      vault_edition         = matrix.edition
      vault_install_dir     = global.vault_install_dir[matrix.artifact_type]
      vault_product_version = matrix.artifact_source == "local" ? step.get_local_metadata.version : var.vault_product_version
      vault_revision        = matrix.artifact_source == "local" ? step.get_local_metadata.revision : var.vault_revision
      vault_build_date      = matrix.artifact_source == "local" ? step.get_local_metadata.build_date : var.vault_build_date
      vault_root_token      = step.create_primary_cluster.root_token
    }
  }

  step "verify_ui" {
    description = global.description.verify_ui
    module      = module.vault_verify_ui
    depends_on  = [step.get_primary_cluster_ips]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = quality.vault_ui_assets

    variables {
      vault_addr = step.create_primary_cluster.api_addr_localhost
      hosts      = step.create_primary_cluster_targets.hosts
    }
  }

  step "verify_secrets_engines_on_primary" {
    description = global.description.verify_secrets_engines_create
    module      = module.vault_verify_secrets_engines_create
    depends_on  = [step.get_primary_cluster_ips]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_api_auth_userpass_login_write,
      quality.vault_api_auth_userpass_user_write,
      quality.vault_api_identity_entity_write,
      quality.vault_api_identity_entity_alias_write,
      quality.vault_api_identity_group_write,
      quality.vault_api_identity_oidc_config_write,
      quality.vault_api_identity_oidc_introspect_write,
      quality.vault_api_identity_oidc_key_write,
      quality.vault_api_identity_oidc_key_rotate_write,
      quality.vault_api_identity_oidc_role_write,
      quality.vault_api_identity_oidc_token_read,
      quality.vault_api_sys_auth_userpass_user_write,
      quality.vault_api_sys_policy_write,
      quality.vault_mount_auth,
      quality.vault_mount_kv,
      quality.vault_secrets_kv_write,
    ]

    variables {
      hosts             = step.create_primary_cluster_targets.hosts
      leader_host       = step.get_primary_cluster_ips.leader_host
      vault_addr        = step.create_primary_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      vault_root_token  = step.create_primary_cluster.root_token
    }
  }

  step "configure_performance_replication_primary" {
    description = <<-EOF
      Create the necessary superuser auth policy necessary for performance replication, assign it
      to a our previously create test user, and enable performance replication on the primary
      cluster.
    EOF
    module      = module.vault_setup_perf_primary
    depends_on = [
      // Wait for both clusters to be up and healthy...
      step.get_primary_cluster_ips,
      step.get_secondary_cluster_ips,
      step.verify_secrets_engines_on_primary,
      // Wait base verification to complete...
      step.verify_vault_version,
      step.verify_ui,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_api_auth_userpass_login_write,
      quality.vault_api_auth_userpass_user_write,
      quality.vault_api_identity_entity_write,
      quality.vault_api_identity_entity_alias_write,
      quality.vault_api_identity_group_write,
      quality.vault_api_identity_oidc_config_write,
      quality.vault_api_identity_oidc_introspect_write,
      quality.vault_api_identity_oidc_key_write,
      quality.vault_api_identity_oidc_key_rotate_write,
      quality.vault_api_identity_oidc_role_write,
      quality.vault_api_identity_oidc_token_read,
      quality.vault_api_sys_auth_userpass_user_write,
      quality.vault_api_sys_policy_write,
      quality.vault_mount_auth,
      quality.vault_mount_kv,
      quality.vault_secrets_kv_write,
    ]

    variables {
      ip_version          = matrix.ip_version
      primary_leader_host = step.get_primary_cluster_ips.leader_host
      vault_addr          = step.create_primary_cluster.api_addr_localhost
      vault_install_dir   = global.vault_install_dir[matrix.artifact_type]
      vault_root_token    = step.create_primary_cluster.root_token
    }
  }

  step "generate_secondary_token" {
    description = <<-EOF
      Generate a random token and configure the performance replication primary secondary-token and
      configure the Vault cluster primary replication with the token. Export the wrapping token
      so that secondary clusters can utilize it.
    EOF
    module      = module.generate_secondary_token
    depends_on  = [step.configure_performance_replication_primary]

    verifies = quality.vault_api_sys_replication_performance_primary_secondary_token_write

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      ip_version          = matrix.ip_version
      primary_leader_host = step.get_primary_cluster_ips.leader_host
      replication_type    = "performance"
      vault_addr          = step.create_primary_cluster.api_addr_localhost
      vault_install_dir   = global.vault_install_dir[matrix.artifact_type]
      vault_root_token    = step.create_primary_cluster.root_token
    }
  }

  step "configure_performance_replication_secondary" {
    description = <<-EOF
      Enable performance replication on the secondary cluster with the wrapping token created by
      the primary cluster.
    EOF
    module      = module.vault_setup_replication_secondary
    depends_on  = [step.generate_secondary_token]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = quality.vault_api_sys_replication_performance_secondary_enable_write

    variables {
      ip_version            = matrix.ip_version
      secondary_leader_host = step.get_secondary_cluster_ips.leader_host
      replication_type      = "performance"
      vault_addr            = step.create_secondary_cluster.api_addr_localhost
      vault_install_dir     = global.vault_install_dir[matrix.artifact_type]
      vault_root_token      = step.create_secondary_cluster.root_token
      wrapping_token        = step.generate_secondary_token.secondary_token
    }
  }

  step "unseal_secondary_followers" {
    description = <<-EOF
      After replication is enabled the secondary cluster followers need to be unsealed.
      Secondary unseal keys are passed differently depending primary and secondary seal
      type combinations. See the guide for more information:
        https://developer.hashicorp.com/vault/docs/enterprise/replication#seals
    EOF
    module      = module.vault_unseal_replication_followers
    depends_on = [
      step.create_primary_cluster,
      step.create_secondary_cluster,
      step.get_secondary_cluster_ips,
      step.configure_performance_replication_secondary
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      hosts             = step.get_secondary_cluster_ips.follower_hosts
      vault_addr        = step.create_secondary_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      vault_seal_type   = matrix.primary_seal == "shamir" ? matrix.primary_seal : matrix.secondary_seal
      vault_unseal_keys = matrix.primary_seal == "shamir" ? step.create_primary_cluster.unseal_keys_hex : step.create_primary_cluster.recovery_keys_hex
    }
  }

  step "verify_secondary_cluster_is_unsealed_after_enabling_replication" {
    description = global.description.verify_vault_unsealed
    module      = module.vault_wait_for_cluster_unsealed
    depends_on = [
      step.unseal_secondary_followers
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_auto_unseals_after_autopilot_upgrade,
      quality.vault_seal_awskms,
      quality.vault_seal_pkcs11,
      quality.vault_seal_shamir,
    ]

    variables {
      hosts             = step.create_secondary_cluster_targets.hosts
      vault_addr        = step.create_secondary_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
    }
  }

  step "verify_performance_replication" {
    description = <<-EOF
      Verify that the performance replication status meets our expectations after enabling replication
      and ensuring that all secondary nodes are unsealed.
    EOF
    module      = module.vault_verify_performance_replication
    depends_on  = [step.verify_secondary_cluster_is_unsealed_after_enabling_replication]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_api_sys_replication_performance_read_connection_status_connected,
      quality.vault_api_sys_replication_performance_status_read,
      quality.vault_api_sys_replication_performance_status_read_cluster_address,
      quality.vault_api_sys_replication_performance_status_read_state_not_idle,
      quality.vault_api_sys_replication_performance_status_known_primary_cluster_addrs,
    ]

    variables {
      ip_version            = matrix.ip_version
      primary_leader_host   = step.get_primary_cluster_ips.leader_host
      secondary_leader_host = step.get_secondary_cluster_ips.leader_host
      vault_addr            = step.create_primary_cluster.api_addr_localhost
      vault_install_dir     = global.vault_install_dir[matrix.artifact_type]
    }
  }

  step "verify_replicated_data" {
    description = global.description.verify_secrets_engines_read
    module      = module.vault_verify_secrets_engines_read
    depends_on = [
      step.verify_performance_replication,
      step.get_secondary_cluster_ips,
      step.verify_secrets_engines_on_primary
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_api_auth_userpass_login_write,
      quality.vault_api_identity_entity_read,
      quality.vault_api_identity_oidc_config_read,
      quality.vault_api_identity_oidc_key_read,
      quality.vault_api_identity_oidc_role_read,
      quality.vault_secrets_kv_read
    ]

    variables {
      create_state      = step.verify_secrets_engines_on_primary.state
      hosts             = step.get_secondary_cluster_ips.follower_hosts
      vault_addr        = step.create_secondary_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      vault_root_token  = step.create_secondary_cluster.root_token
      verify_pki_certs  = false
    }
  }

  step "add_additional_nodes_to_primary_cluster" {
    description = <<-EOF
      Add additional nodes the Vault Cluster to prepare for our catostrophic failure simulation.
      These nodes will use a different storage storage_node_prefix
    EOF
    module      = module.vault_cluster
    depends_on = [
      step.create_vpc,
      step.create_primary_backend_cluster,
      step.create_primary_cluster,
      step.verify_replicated_data,
      step.create_primary_cluster_additional_targets
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      // unique to this invocation of the module
      quality.vault_autojoins_new_nodes_into_initialized_cluster,
      // verified in modules
      quality.vault_artifact_bundle,
      quality.vault_artifact_deb,
      quality.vault_artifact_rpm,
      quality.vault_audit_log,
      quality.vault_audit_socket,
      quality.vault_audit_syslog,
      quality.vault_autojoin_aws,
      quality.vault_config_env_variables,
      quality.vault_config_file,
      quality.vault_config_log_level,
      quality.vault_init,
      quality.vault_license_required_ent,
      quality.vault_listener_ipv4,
      quality.vault_listener_ipv6,
      quality.vault_service_start,
      quality.vault_storage_backend_consul,
      quality.vault_storage_backend_raft,
      // verified in enos_vault_start resource
      quality.vault_api_sys_config_read,
      quality.vault_api_sys_ha_status_read,
      quality.vault_api_sys_health_read,
      quality.vault_api_sys_host_info_read,
      quality.vault_api_sys_replication_status_read,
      quality.vault_api_sys_seal_status_api_read_matches_sys_health,
      quality.vault_api_sys_storage_raft_configuration_read,
      quality.vault_api_sys_storage_raft_autopilot_configuration_read,
      quality.vault_api_sys_storage_raft_autopilot_state_read,
      quality.vault_service_systemd_notified,
      quality.vault_service_systemd_unit,
      quality.vault_cli_status_exit_code,
    ]

    variables {
      artifactory_release     = matrix.artifact_source == "artifactory" ? step.build_vault.vault_artifactory_release : null
      backend_cluster_name    = step.create_primary_cluster_backend_targets.cluster_name
      backend_cluster_tag_key = global.backend_tag_key
      cluster_name            = step.create_primary_cluster_targets.cluster_name
      config_mode             = matrix.config_mode
      consul_license          = (matrix.primary_backend == "consul" && matrix.consul_edition == "ent") ? step.read_backend_license.license : null
      consul_release = matrix.primary_backend == "consul" ? {
        edition = matrix.consul_edition
        version = matrix.consul_version
      } : null
      enable_audit_devices = var.vault_enable_audit_devices
      force_unseal         = matrix.primary_seal == "shamir"
      hosts                = step.create_primary_cluster_additional_targets.hosts
      // Don't init when adding nodes into the cluster.
      initialize_cluster  = false
      install_dir         = global.vault_install_dir[matrix.artifact_type]
      ip_version          = matrix.ip_version
      license             = matrix.edition != "ce" ? step.read_vault_license.license : null
      local_artifact_path = local.artifact_path
      manage_service      = local.manage_service
      packages            = concat(global.packages, global.distro_packages[matrix.distro][global.distro_version[matrix.distro]])
      root_token          = step.create_primary_cluster.root_token
      seal_attributes     = step.create_primary_seal_key.attributes
      seal_type           = matrix.primary_seal
      shamir_unseal_keys  = matrix.primary_seal == "shamir" ? step.create_primary_cluster.unseal_keys_hex : null
      storage_backend     = matrix.primary_backend
      storage_node_prefix = "newprimary_node"
    }
  }

  step "verify_additional_primary_nodes_are_unsealed" {
    description = global.description.verify_vault_unsealed
    module      = module.vault_wait_for_cluster_unsealed
    depends_on  = [step.add_additional_nodes_to_primary_cluster]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_auto_unseals_after_autopilot_upgrade,
      quality.vault_seal_awskms,
      quality.vault_seal_pkcs11,
      quality.vault_seal_shamir,
    ]

    variables {
      hosts             = step.create_primary_cluster_additional_targets.hosts
      vault_addr        = step.add_additional_nodes_to_primary_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
    }
  }

  step "verify_raft_auto_join_voter" {
    description = global.description.verify_raft_cluster_all_nodes_are_voters
    skip_step   = matrix.primary_backend != "raft"
    module      = module.vault_verify_raft_auto_join_voter
    depends_on = [
      step.add_additional_nodes_to_primary_cluster,
      step.create_primary_cluster,
      step.verify_additional_primary_nodes_are_unsealed
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = quality.vault_raft_voters

    variables {
      hosts             = step.create_primary_cluster_additional_targets.hosts
      ip_version        = matrix.ip_version
      vault_addr        = step.create_primary_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      vault_root_token  = step.create_primary_cluster.root_token
    }
  }

  step "remove_primary_follower_1" {
    description = <<-EOF
      Simulate a catostrophic failure by forcefully removing the a follower node from the Vault
      Cluster.
    EOF
    module      = module.shutdown_node
    depends_on = [
      step.verify_additional_primary_nodes_are_unsealed,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      host = step.get_primary_cluster_ips.follower_hosts["0"]
    }
  }

  step "remove_primary_leader" {
    description = <<-EOF
      Simulate a catostrophic failure by forcefully removing the the primary leader node from the
      Vault Cluster without allowing a graceful shutdown.
    EOF
    module      = module.shutdown_node
    depends_on = [
      step.get_primary_cluster_ips,
      step.remove_primary_follower_1
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      host = step.get_primary_cluster_ips.leader_host
    }
  }

  step "get_remaining_hosts_replication_data" {
    description = <<-EOF
      An arithmetic module that we use to determine various metadata about the the leader and
      follower nodes of the primary cluster so that we can correctly enable performance replication.

      We execute this to determine information about our hosts after having forced the leader
      and a follower from the cluster.
    EOF

    module = module.replication_data
    depends_on = [
      step.get_primary_cluster_ips,
      step.remove_primary_leader,
    ]

    variables {
      added_hosts           = step.create_primary_cluster_additional_targets.hosts
      initial_hosts         = step.create_primary_cluster_targets.hosts
      removed_follower_host = step.get_primary_cluster_ips.follower_hosts["0"]
      removed_primary_host  = step.get_primary_cluster_ips.leader_host
    }
  }

  step "wait_for_leader_in_remaining_hosts" {
    description = global.description.wait_for_cluster_to_have_leader
    module      = module.vault_wait_for_leader
    depends_on = [
      step.remove_primary_leader,
      step.get_remaining_hosts_replication_data,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_api_sys_leader_read,
      quality.vault_unseal_ha_leader_election,
    ]

    variables {
      hosts             = step.get_remaining_hosts_replication_data.remaining_hosts
      ip_version        = matrix.ip_version
      timeout           = 120 // seconds
      vault_addr        = step.create_primary_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      vault_root_token  = step.create_primary_cluster.root_token
    }
  }

  step "get_updated_primary_cluster_ips" {
    description = global.description.get_vault_cluster_ip_addresses
    module      = module.vault_get_cluster_ips
    depends_on = [
      step.get_remaining_hosts_replication_data,
      step.wait_for_leader_in_remaining_hosts,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_api_sys_ha_status_read,
      quality.vault_api_sys_leader_read,
      quality.vault_cli_operator_members,
    ]

    variables {
      hosts             = step.get_remaining_hosts_replication_data.remaining_hosts
      ip_version        = matrix.ip_version
      vault_addr        = step.create_primary_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      vault_root_token  = step.create_primary_cluster.root_token
    }
  }

  step "verify_updated_performance_replication" {
    description = <<-EOF
      Verify that the performance replication status meets our expectations after the new leader
      election.
    EOF

    module = module.vault_verify_performance_replication
    depends_on = [
      step.get_remaining_hosts_replication_data,
      step.wait_for_leader_in_remaining_hosts,
      step.get_updated_primary_cluster_ips,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_api_sys_replication_performance_read_connection_status_connected,
      quality.vault_api_sys_replication_performance_status_known_primary_cluster_addrs,
      quality.vault_api_sys_replication_performance_status_read,
      quality.vault_api_sys_replication_performance_status_read_state_not_idle,
      quality.vault_api_sys_replication_performance_status_read_cluster_address,
    ]

    variables {
      ip_version            = matrix.ip_version
      primary_leader_host   = step.get_updated_primary_cluster_ips.leader_host
      secondary_leader_host = step.get_secondary_cluster_ips.leader_host
      vault_addr            = step.create_primary_cluster.api_addr_localhost
      vault_install_dir     = global.vault_install_dir[matrix.artifact_type]
    }
  }

  output "audit_device_file_path" {
    description = "The file path for the file audit device, if enabled"
    value       = step.create_primary_cluster.audit_device_file_path
  }

  output "primary_cluster_hosts" {
    description = "The Vault primary cluster target hosts"
    value       = step.create_primary_cluster_targets.hosts
  }

  output "primary_cluster_additional_hosts" {
    description = "The Vault added new node on primary cluster target hosts"
    value       = step.create_primary_cluster_additional_targets.hosts
  }

  output "primary_cluster_root_token" {
    description = "The Vault primary cluster root token"
    value       = step.create_primary_cluster.root_token
  }

  output "primary_cluster_unseal_keys_b64" {
    description = "The Vault primary cluster unseal keys"
    value       = step.create_primary_cluster.unseal_keys_b64
  }

  output "primary_cluster_unseal_keys_hex" {
    description = "The Vault primary cluster unseal keys hex"
    value       = step.create_primary_cluster.unseal_keys_hex
  }

  output "primary_cluster_recovery_key_shares" {
    description = "The Vault primary cluster recovery key shares"
    value       = step.create_primary_cluster.recovery_key_shares
  }

  output "primary_cluster_recovery_keys_b64" {
    description = "The Vault primary cluster recovery keys b64"
    value       = step.create_primary_cluster.recovery_keys_b64
  }

  output "primary_cluster_recovery_keys_hex" {
    description = "The Vault primary cluster recovery keys hex"
    value       = step.create_primary_cluster.recovery_keys_hex
  }

  output "secondary_cluster_hosts" {
    description = "The Vault secondary cluster public IPs"
    value       = step.create_secondary_cluster_targets.hosts
  }

  output "secondary_cluster_root_token" {
    description = "The Vault secondary cluster root token"
    value       = step.create_secondary_cluster.root_token
  }

  output "secrets_engines_state" {
    description = "The state of configured secrets engines"
    value       = step.verify_secrets_engines_on_primary.state
  }

  output "performance_secondary_token" {
    description = "The performance secondary replication token"
    value       = step.generate_secondary_token.secondary_token
  }

  output "remaining_hosts" {
    description = "The Vault cluster primary hosts after removing the leader and follower"
    value       = step.get_remaining_hosts_replication_data.remaining_hosts
  }

  output "initial_primary_replication_status" {
    description = "The Vault primary cluster performance replication status"
    value       = step.verify_performance_replication.primary_replication_status
  }

  output "initial_known_primary_cluster_addresses" {
    description = "The initial known Vault primary cluster addresses"
    value       = step.verify_performance_replication.known_primary_cluster_addrs
  }

  output "initial_secondary_performance_replication_status" {
    description = "The Vault secondary cluster performance replication status"
    value       = step.verify_performance_replication.secondary_replication_status
  }

  output "intial_primary_replication_data_secondaries" {
    description = "The Vault primary cluster secondaries connection status"
    value       = step.verify_performance_replication.primary_replication_data_secondaries
  }

  output "initial_secondary_replication_data_primaries" {
    description = "The Vault secondary cluster primaries connection status"
    value       = step.verify_performance_replication.secondary_replication_data_primaries
  }

  output "updated_primary_replication_status" {
    description = "The Vault updated primary cluster performance replication status"
    value       = step.verify_updated_performance_replication.primary_replication_status
  }

  output "updated_known_primary_cluster_addresses" {
    description = "The Vault secondary cluster performance replication status"
    value       = step.verify_updated_performance_replication.known_primary_cluster_addrs
  }

  output "updated_secondary_replication_status" {
    description = "The Vault updated secondary cluster performance replication status"
    value       = step.verify_updated_performance_replication.secondary_replication_status
  }

  output "updated_primary_replication_data_secondaries" {
    description = "The Vault updated primary cluster secondaries connection status"
    value       = step.verify_updated_performance_replication.primary_replication_data_secondaries
  }

  output "updated_secondary_replication_data_primaries" {
    description = "The Vault updated secondary cluster primaries connection status"
    value       = step.verify_updated_performance_replication.secondary_replication_data_primaries
  }
}
