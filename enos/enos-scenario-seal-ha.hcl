// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

scenario "seal_ha" {
  description = <<-EOF
    The seal_ha scenario verifies Vault Enterprises seal HA capabilities. The build can be a local
    branch, any CRT built Vault Enterprise artifact saved to the local machine, or any CRT built
    Vault Enterprise artifact in the stable channel in Artifactory.

    The scenario deploys a Vault Enterprise cluster with the candidate build and enables a single
    primary seal, mounts various engines and writes data, then establishes seal HA with a secondary
    seal, then removes the primary and verifies data integrity and seal data migration. It also
    verifies that the cluster is able to recover from a forced leader election after the initial
    seal rewrap. The scenario also performs standard baseline verification that is not specific to
    seal_ha.

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
    arch            = global.archs
    artifact_source = global.artifact_sources
    artifact_type   = global.artifact_types
    backend         = global.backends
    config_mode     = global.config_modes
    consul_edition  = global.consul_editions
    consul_version  = global.consul_versions
    distro          = global.distros
    edition         = global.enterprise_editions
    ip_version      = global.ip_versions
    // Seal HA is only supported with auto-unseal devices.
    primary_seal   = ["awskms", "pkcs11"]
    secondary_seal = ["awskms", "pkcs11"]

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
      ip_version = ["6"]
      backend    = ["consul"]
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
    manage_service = matrix.artifact_type == "bundle"
  }

  step "get_local_metadata" {
    description = global.description.get_local_metadata
    skip_step   = matrix.artifact_source != "local"
    module      = module.get_local_metadata
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
  // the edition is "ent".
  step "read_backend_license" {
    description = global.description.read_backend_license
    skip_step   = matrix.backend == "raft" || matrix.consul_edition == "ce"
    module      = module.read_license

    variables {
      file_name = global.backend_license_path
    }
  }

  step "read_vault_license" {
    description = global.description.read_vault_license
    skip_step   = matrix.edition == "ce"
    module      = module.read_license

    variables {
      file_name = global.vault_license_path
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

  step "create_vault_cluster_targets" {
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

  step "create_vault_cluster_backend_targets" {
    description = global.description.create_vault_cluster_targets
    module      = matrix.backend == "consul" ? module.target_ec2_instances : module.target_ec2_shim
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

  step "create_backend_cluster" {
    description = global.description.create_backend_cluster
    module      = "backend_${matrix.backend}"
    depends_on = [
      step.create_vault_cluster_backend_targets
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
      cluster_name    = step.create_vault_cluster_backend_targets.cluster_name
      cluster_tag_key = global.backend_tag_key
      hosts           = step.create_vault_cluster_backend_targets.hosts
      license         = (matrix.backend == "consul" && matrix.consul_edition == "ent") ? step.read_backend_license.license : null
      release = {
        edition = matrix.consul_edition
        version = matrix.consul_version
      }
    }
  }

  step "create_vault_cluster" {
    description = global.description.create_vault_cluster
    module      = module.vault_cluster
    depends_on = [
      step.create_backend_cluster,
      step.build_vault,
      step.create_vault_cluster_targets
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
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
      quality.vault_api_sys_health_read,
      quality.vault_api_sys_config_read,
      quality.vault_api_sys_ha_status_read,
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
      backend_cluster_name    = step.create_vault_cluster_backend_targets.cluster_name
      backend_cluster_tag_key = global.backend_tag_key
      cluster_name            = step.create_vault_cluster_targets.cluster_name
      config_mode             = matrix.config_mode
      consul_license          = (matrix.backend == "consul" && matrix.consul_edition == "ent") ? step.read_backend_license.license : null
      consul_release = matrix.backend == "consul" ? {
        edition = matrix.consul_edition
        version = matrix.consul_version
      } : null
      enable_audit_devices = var.vault_enable_audit_devices
      hosts                = step.create_vault_cluster_targets.hosts
      install_dir          = global.vault_install_dir[matrix.artifact_type]
      ip_version           = matrix.ip_version
      license              = matrix.edition != "ce" ? step.read_vault_license.license : null
      local_artifact_path  = local.artifact_path
      manage_service       = local.manage_service
      packages             = concat(global.packages, global.distro_packages[matrix.distro][global.distro_version[matrix.distro]])
      // Only configure our primary seal during our initial cluster setup
      seal_attributes = step.create_primary_seal_key.attributes
      seal_type       = matrix.primary_seal
      storage_backend = matrix.backend
    }
  }

  // Wait for our cluster to elect a leader
  step "wait_for_leader" {
    description = global.description.wait_for_cluster_to_have_leader
    module      = module.vault_wait_for_leader
    depends_on  = [step.create_vault_cluster]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_api_sys_leader_read,
      quality.vault_unseal_ha_leader_election,
    ]

    variables {
      hosts             = step.create_vault_cluster_targets.hosts
      ip_version        = matrix.ip_version
      timeout           = 120 // seconds
      vault_addr        = step.create_vault_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      vault_root_token  = step.create_vault_cluster.root_token
    }
  }

  step "get_vault_cluster_ips" {
    description = global.description.get_vault_cluster_ip_addresses
    module      = module.vault_get_cluster_ips
    depends_on  = [step.wait_for_leader]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_api_sys_ha_status_read,
      quality.vault_api_sys_leader_read,
      quality.vault_cli_operator_members,
    ]

    variables {
      hosts             = step.create_vault_cluster_targets.hosts
      ip_version        = matrix.ip_version
      vault_addr        = step.create_vault_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      vault_root_token  = step.create_vault_cluster.root_token
    }
  }

  step "verify_vault_unsealed" {
    description = global.description.verify_vault_unsealed
    module      = module.vault_wait_for_cluster_unsealed
    depends_on  = [step.wait_for_leader]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_seal_awskms,
      quality.vault_seal_pkcs11,
      quality.vault_seal_shamir,
    ]

    variables {
      hosts             = step.create_vault_cluster_targets.hosts
      vault_addr        = step.create_vault_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
    }
  }

  // Write some test data before we create the new seal
  step "verify_secrets_engines_create" {
    description = global.description.verify_secrets_engines_create
    module      = module.vault_verify_secrets_engines_create
    depends_on = [
      step.create_vault_cluster,
      step.get_vault_cluster_ips,
      step.verify_vault_unsealed,
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
      hosts             = step.create_vault_cluster_targets.hosts
      leader_host       = step.get_vault_cluster_ips.leader_host
      vault_addr        = step.create_vault_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      vault_root_token  = step.create_vault_cluster.root_token
    }
  }

  // Wait for the initial seal rewrap to complete before we add our HA seal.
  step "wait_for_initial_seal_rewrap" {
    description = global.description.wait_for_seal_rewrap
    module      = module.vault_wait_for_seal_rewrap
    depends_on = [
      step.verify_secrets_engines_create,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_api_sys_sealwrap_rewrap_read_entries_processed_eq_entries_succeeded_post_rewrap,
      quality.vault_api_sys_sealwrap_rewrap_read_entries_processed_gt_zero_post_rewrap,
      quality.vault_api_sys_sealwrap_rewrap_read_is_running_false_post_rewrap,
      quality.vault_api_sys_sealwrap_rewrap_read_no_entries_fail_during_rewrap,
    ]

    variables {
      hosts             = step.create_vault_cluster_targets.hosts
      vault_addr        = step.create_vault_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      vault_root_token  = step.create_vault_cluster.root_token
    }
  }

  step "stop_vault" {
    description = "${global.description.stop_vault}. We do this to write new seal configuration."
    module      = module.stop_vault
    depends_on = [
      step.create_vault_cluster,
      step.verify_secrets_engines_create,
      step.wait_for_initial_seal_rewrap,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      hosts = step.create_vault_cluster_targets.hosts
    }
  }

  // Add the secondary seal to the cluster
  step "add_ha_seal_to_cluster" {
    description = global.description.enable_multiseal
    module      = module.start_vault
    depends_on  = [step.stop_vault]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = quality.vault_multiseal_enable

    variables {
      cluster_name              = step.create_vault_cluster_targets.cluster_name
      hosts                     = step.create_vault_cluster_targets.hosts
      install_dir               = global.vault_install_dir[matrix.artifact_type]
      ip_version                = matrix.ip_version
      license                   = matrix.edition != "ce" ? step.read_vault_license.license : null
      manage_service            = local.manage_service
      seal_attributes           = step.create_primary_seal_key.attributes
      seal_attributes_secondary = step.create_secondary_seal_key.attributes
      seal_type                 = matrix.primary_seal
      seal_type_secondary       = matrix.secondary_seal
      storage_backend           = matrix.backend
    }
  }

  // Wait for our cluster to elect a leader
  step "wait_for_leader_election" {
    description = global.description.wait_for_cluster_to_have_leader
    module      = module.vault_wait_for_leader
    depends_on  = [step.add_ha_seal_to_cluster]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_api_sys_leader_read,
      quality.vault_unseal_ha_leader_election,
    ]

    variables {
      hosts             = step.create_vault_cluster_targets.hosts
      ip_version        = matrix.ip_version
      timeout           = 120 // seconds
      vault_addr        = step.create_vault_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      vault_root_token  = step.create_vault_cluster.root_token
    }
  }

  step "get_leader_ip_for_step_down" {
    description = global.description.get_vault_cluster_ip_addresses
    module      = module.vault_get_cluster_ips
    depends_on  = [step.wait_for_leader_election]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_api_sys_ha_status_read,
      quality.vault_api_sys_leader_read,
      quality.vault_cli_operator_members,
    ]

    variables {
      hosts             = step.create_vault_cluster_targets.hosts
      ip_version        = matrix.ip_version
      vault_addr        = step.create_vault_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      vault_root_token  = step.create_vault_cluster.root_token
    }
  }

  // Force a step down to trigger a new leader election
  step "vault_leader_step_down" {
    description = global.description.vault_leader_step_down
    module      = module.vault_step_down
    depends_on  = [step.get_leader_ip_for_step_down]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_api_sys_step_down_steps_down,
      quality.vault_cli_operator_step_down,
    ]

    variables {
      leader_host       = step.get_leader_ip_for_step_down.leader_host
      vault_addr        = step.create_vault_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      vault_root_token  = step.create_vault_cluster.root_token
    }
  }

  // Wait for our cluster to elect a leader
  step "wait_for_new_leader" {
    description = global.description.wait_for_cluster_to_have_leader
    module      = module.vault_wait_for_leader
    depends_on  = [step.vault_leader_step_down]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_api_sys_leader_read,
      quality.vault_unseal_ha_leader_election,
    ]

    variables {
      timeout           = 120 // seconds
      ip_version        = matrix.ip_version
      hosts             = step.create_vault_cluster_targets.hosts
      vault_addr        = step.create_vault_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      vault_root_token  = step.create_vault_cluster.root_token
    }
  }

  step "get_updated_cluster_ips" {
    description = global.description.get_vault_cluster_ip_addresses
    module      = module.vault_get_cluster_ips
    depends_on  = [step.wait_for_new_leader]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_api_sys_ha_status_read,
      quality.vault_api_sys_leader_read,
      quality.vault_cli_operator_members,
    ]

    variables {
      hosts             = step.create_vault_cluster_targets.hosts
      ip_version        = matrix.ip_version
      vault_addr        = step.create_vault_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      vault_root_token  = step.create_vault_cluster.root_token
    }
  }

  step "verify_vault_unsealed_with_new_seal" {
    description = global.description.verify_vault_unsealed
    module      = module.vault_wait_for_cluster_unsealed
    depends_on  = [step.wait_for_new_leader]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_seal_awskms,
      quality.vault_seal_pkcs11,
      quality.vault_seal_shamir,
    ]

    variables {
      hosts             = step.create_vault_cluster_targets.hosts
      vault_addr        = step.create_vault_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
    }
  }

  // Wait for the seal rewrap to complete and verify that no entries failed
  step "wait_for_seal_rewrap" {
    description = global.description.wait_for_seal_rewrap
    module      = module.vault_wait_for_seal_rewrap
    depends_on = [
      step.add_ha_seal_to_cluster,
      step.verify_vault_unsealed_with_new_seal,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_api_sys_sealwrap_rewrap_read_entries_processed_eq_entries_succeeded_post_rewrap,
      quality.vault_api_sys_sealwrap_rewrap_read_entries_processed_gt_zero_post_rewrap,
      quality.vault_api_sys_sealwrap_rewrap_read_is_running_false_post_rewrap,
      quality.vault_api_sys_sealwrap_rewrap_read_no_entries_fail_during_rewrap,
    ]

    variables {
      hosts             = step.create_vault_cluster_targets.hosts
      vault_addr        = step.create_vault_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      vault_root_token  = step.create_vault_cluster.root_token
    }
  }

  // Perform all of our standard verifications after we've enabled multiseal
  step "verify_vault_version" {
    description = global.description.verify_vault_version
    module      = module.vault_verify_version
    depends_on  = [step.wait_for_seal_rewrap]

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
      hosts                 = step.create_vault_cluster_targets.hosts
      vault_addr            = step.create_vault_cluster.api_addr_localhost
      vault_edition         = matrix.edition
      vault_install_dir     = global.vault_install_dir[matrix.artifact_type]
      vault_product_version = matrix.artifact_source == "local" ? step.get_local_metadata.version : var.vault_product_version
      vault_revision        = matrix.artifact_source == "local" ? step.get_local_metadata.revision : var.vault_revision
      vault_build_date      = matrix.artifact_source == "local" ? step.get_local_metadata.build_date : var.vault_build_date
      vault_root_token      = step.create_vault_cluster.root_token
    }
  }

  step "verify_raft_auto_join_voter" {
    description = global.description.verify_raft_cluster_all_nodes_are_voters
    skip_step   = matrix.backend != "raft"
    module      = module.vault_verify_raft_auto_join_voter
    depends_on  = [step.wait_for_seal_rewrap]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = quality.vault_raft_voters

    variables {
      hosts             = step.create_vault_cluster_targets.hosts
      ip_version        = matrix.ip_version
      vault_addr        = step.create_vault_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      vault_root_token  = step.create_vault_cluster.root_token
    }
  }

  step "verify_replication" {
    description = global.description.verify_replication_status
    module      = module.vault_verify_replication
    depends_on  = [step.wait_for_seal_rewrap]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_replication_ce_disabled,
      quality.vault_replication_ent_dr_available,
      quality.vault_replication_ent_pr_available,
    ]

    variables {
      hosts         = step.create_vault_cluster_targets.hosts
      vault_addr    = step.create_vault_cluster.api_addr_localhost
      vault_edition = matrix.edition
    }
  }

  // Make sure our data is still available
  step "verify_secrets_engines_read" {
    description = global.description.verify_secrets_engines_read
    module      = module.vault_verify_secrets_engines_read
    depends_on  = [step.wait_for_seal_rewrap]

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
      create_state      = step.verify_secrets_engines_create.state
      hosts             = step.get_updated_cluster_ips.follower_hosts
      vault_addr        = step.create_vault_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      vault_root_token  = step.create_vault_cluster.root_token
    }
  }

  step "verify_log_secrets" {
    skip_step = !var.vault_enable_audit_devices || !var.verify_log_secrets

    description = global.description.verify_log_secrets
    module      = module.verify_log_secrets
    depends_on = [
      step.verify_secrets_engines_read,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_audit_log_secrets,
      quality.vault_journal_secrets,
      quality.vault_radar_index_create,
      quality.vault_radar_scan_file,
    ]

    variables {
      audit_log_file_path = step.create_vault_cluster.audit_device_file_path
      leader_host         = step.get_updated_cluster_ips.leader_host
      vault_addr          = step.create_vault_cluster.api_addr_localhost
      vault_root_token    = step.create_vault_cluster.root_token
    }
  }

  step "verify_ui" {
    description = global.description.verify_ui
    module      = module.vault_verify_ui
    depends_on  = [step.wait_for_seal_rewrap]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = quality.vault_ui_assets

    variables {
      hosts      = step.create_vault_cluster_targets.hosts
      vault_addr = step.create_vault_cluster.api_addr_localhost
    }
  }

  step "verify_seal_type" {
    description = "${global.description.verify_seal_type} In this case we expect to have 'multiseal'."
    // Don't run this on versions less than 1.16.0-beta1 until VAULT-21053 is fixed on prior branches.
    skip_step  = semverconstraint(var.vault_product_version, "< 1.16.0-beta1")
    module     = module.verify_seal_type
    depends_on = [step.wait_for_seal_rewrap]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = quality.vault_status_seal_type

    variables {
      hosts             = step.create_vault_cluster_targets.hosts
      seal_type         = "multiseal"
      vault_addr        = step.create_vault_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
    }
  }

  // Now we'll migrate away from our initial seal to our secondary seal

  // Stop the vault service on all nodes before we restart with new seal config
  step "stop_vault_for_migration" {
    description = "${global.description.stop_vault}. We do this to remove the old primary seal."
    module      = module.stop_vault
    depends_on = [
      step.wait_for_seal_rewrap,
      step.verify_secrets_engines_read,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      hosts = step.create_vault_cluster_targets.hosts
    }
  }

  // Remove the "primary" seal from the cluster. Set our "secondary" seal to priority 1. We do this
  // by restarting vault with the correct config.
  step "remove_primary_seal" {
    description = <<-EOF
      Reconfigure the vault cluster seal configuration with only our secondary seal config which
      will force a seal migration to a single seal.
    EOF
    module      = module.start_vault
    depends_on  = [step.stop_vault_for_migration]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = quality.vault_config_multiseal_is_toggleable

    variables {
      cluster_name    = step.create_vault_cluster_targets.cluster_name
      hosts           = step.create_vault_cluster_targets.hosts
      install_dir     = global.vault_install_dir[matrix.artifact_type]
      ip_version      = matrix.ip_version
      license         = matrix.edition != "ce" ? step.read_vault_license.license : null
      manage_service  = local.manage_service
      seal_alias      = "secondary"
      seal_attributes = step.create_secondary_seal_key.attributes
      seal_type       = matrix.secondary_seal
      storage_backend = matrix.backend
    }
  }

  // Wait for our cluster to elect a leader after restarting vault with a new primary seal
  step "wait_for_leader_after_migration" {
    description = global.description.wait_for_cluster_to_have_leader
    module      = module.vault_wait_for_leader
    depends_on  = [step.remove_primary_seal]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_api_sys_leader_read,
      quality.vault_unseal_ha_leader_election,
    ]

    variables {
      timeout           = 120 // seconds
      ip_version        = matrix.ip_version
      hosts             = step.create_vault_cluster_targets.hosts
      vault_addr        = step.create_vault_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      vault_root_token  = step.create_vault_cluster.root_token
    }
  }

  // Since we've restarted our cluster we might have a new leader and followers. Get the new IPs.
  step "get_cluster_ips_after_migration" {
    description = global.description.get_vault_cluster_ip_addresses
    module      = module.vault_get_cluster_ips
    depends_on  = [step.wait_for_leader_after_migration]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    verifies = [
      quality.vault_api_sys_ha_status_read,
      quality.vault_api_sys_leader_read,
      quality.vault_cli_operator_members,
    ]

    variables {
      hosts             = step.create_vault_cluster_targets.hosts
      ip_version        = matrix.ip_version
      vault_addr        = step.create_vault_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      vault_root_token  = step.create_vault_cluster.root_token
    }
  }

  // Make sure we unsealed
  step "verify_vault_unsealed_after_migration" {
    module     = module.vault_wait_for_cluster_unsealed
    depends_on = [step.wait_for_leader_after_migration]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      hosts             = step.create_vault_cluster_targets.hosts
      vault_addr        = step.create_vault_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
    }
  }

  // Wait for the seal rewrap to complete and verify that no entries failed
  step "wait_for_seal_rewrap_after_migration" {
    module = module.vault_wait_for_seal_rewrap
    depends_on = [
      step.wait_for_leader_after_migration,
      step.verify_vault_unsealed_after_migration,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      hosts             = step.create_vault_cluster_targets.hosts
      vault_addr        = step.create_vault_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      vault_root_token  = step.create_vault_cluster.root_token
    }
  }

  // Make sure our data is still available after migration
  step "verify_secrets_engines_read_after_migration" {
    module     = module.vault_verify_secrets_engines_read
    depends_on = [step.wait_for_seal_rewrap_after_migration]

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
      create_state      = step.verify_secrets_engines_create.state
      hosts             = step.get_cluster_ips_after_migration.follower_hosts
      vault_addr        = step.create_vault_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
      vault_root_token  = step.create_vault_cluster.root_token
    }
  }

  // Make sure we have our secondary seal type after migration
  step "verify_seal_type_after_migration" {
    // Don't run this on versions less than 1.16.0-beta1 until VAULT-21053 is fixed on prior branches.
    skip_step  = semverconstraint(var.vault_product_version, "<= 1.16.0-beta1")
    module     = module.verify_seal_type
    depends_on = [step.wait_for_seal_rewrap_after_migration]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      hosts             = step.create_vault_cluster_targets.hosts
      seal_type         = matrix.secondary_seal
      vault_addr        = step.create_vault_cluster.api_addr_localhost
      vault_install_dir = global.vault_install_dir[matrix.artifact_type]
    }
  }

  output "audit_device_file_path" {
    description = "The file path for the file audit device, if enabled"
    value       = step.create_vault_cluster.audit_device_file_path
  }

  output "cluster_name" {
    description = "The Vault cluster name"
    value       = step.create_vault_cluster.cluster_name
  }

  output "hosts" {
    description = "The Vault cluster target hosts"
    value       = step.create_vault_cluster.hosts
  }

  output "initial_seal_rewrap" {
    description = "The initial seal rewrap status"
    value       = step.wait_for_initial_seal_rewrap.stdout
  }

  output "post_migration_seal_rewrap" {
    description = "The seal rewrap status after migrating the primary seal"
    value       = step.wait_for_seal_rewrap_after_migration.stdout
  }

  output "primary_seal_attributes" {
    description = "The Vault cluster primary seal attributes"
    value       = step.create_primary_seal_key.attributes
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

  output "secondary_seal_attributes" {
    description = "The Vault cluster secondary seal attributes"
    value       = step.create_secondary_seal_key.attributes
  }

  output "secrets_engines_state" {
    description = "The state of configured secrets engines"
    value       = step.verify_secrets_engines_create.state
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
