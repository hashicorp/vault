// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

scenario "ui" {
  description = <<-EOF
    The UI scenario creates a new cluster and runs the existing Ember test suite
    against a live Vault cluster instead of a binary in dev mode. It builds Vault
    from the local branch.

    # How to run this scenario

    For general instructions on running a scenario, refer to the Enos docs: https://eng-handbook.hashicorp.services/internal-tools/enos/running-a-scenario/
    For troubleshooting tips and common errors, see https://eng-handbook.hashicorp.services/internal-tools/enos/troubleshooting/.

    Variables required for all scenario variants:
      - aws_ssh_private_key_path (more info about AWS SSH keypairs: https://eng-handbook.hashicorp.services/internal-tools/enos/getting-started/#set-your-aws-key-pair-name-and-private-key)
      - aws_ssh_keypair_name
      - vault_product_version

    Variables required for some scenario variants:
      - consul_license_path (if using an ENT edition of Consul)
      - vault_license_path (if using an ENT edition of Vault)
  EOF
  matrix {
    backend        = global.backends
    consul_edition = global.consul_editions
    edition        = ["ce", "ent"]
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
      "ce"  = ["ui"]
      "ent" = ["ui", "enterprise", "ent"]
    }
    artifact_path  = abspath(var.vault_artifact_path)
    distro         = "ubuntu"
    consul_version = "1.17.0"
    ip_version     = 4
    seal           = "awskms"
    tags = merge({
      "Project Name" : var.project_name
      "Project" : "Enos",
      "Environment" : "ci"
    }, var.tags)
    vault_install_dir  = var.vault_install_dir
    vault_license_path = abspath(var.vault_license_path != null ? var.vault_license_path : joinpath(path.root, "./support/vault.hclic"))
    vault_tag_key      = "Type" // enos_vault_start expects Type as the tag key
    ui_test_filter     = var.ui_test_filter != null && try(trimspace(var.ui_test_filter), "") != "" ? var.ui_test_filter : (matrix.edition == "ce") ? "!enterprise" : null
  }

  step "build_vault" {
    description = global.description.build_vault
    module      = module.build_local

    variables {
      build_tags      = var.vault_local_build_tags != null ? var.vault_local_build_tags : local.build_tags[matrix.edition]
      artifact_path   = local.artifact_path
      goarch          = local.arch
      goos            = "linux"
      product_version = var.vault_product_version
      artifact_type   = local.artifact_type
      revision        = var.vault_revision
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
      common_tags = local.tags
      ip_version  = local.ip_version
    }
  }

  // This step reads the contents of the backend license if we're using a Consul backend and
  // the edition is "ent".
  step "read_backend_license" {
    description = global.description.read_backend_license
    skip_step   = matrix.backend == "raft" || matrix.consul_edition == "ce"
    module      = module.read_license

    variables {
      file_name = local.backend_license_path
    }
  }

  step "read_vault_license" {
    description = global.description.read_vault_license
    skip_step   = matrix.edition == "ce"
    module      = module.read_license

    variables {
      file_name = local.vault_license_path
    }
  }

  step "create_seal_key" {
    description = global.description.create_seal_key
    module      = "seal_${local.seal}"

    variables {
      cluster_id  = step.create_vpc.cluster_id
      common_tags = global.tags
    }
  }

  step "create_vault_cluster_targets" {
    description = global.description.create_vault_cluster_targets
    module      = module.target_ec2_instances
    depends_on  = [step.create_vpc]

    providers = {
      enos = provider.enos.ubuntu
    }

    variables {
      ami_id          = step.ec2_info.ami_ids["arm64"]["ubuntu"][global.distro_version["ubuntu"]]
      cluster_tag_key = local.vault_tag_key
      common_tags     = local.tags
      seal_key_names  = step.create_seal_key.resource_names
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
      cluster_tag_key = local.backend_tag_key
      common_tags     = local.tags
      seal_key_names  = step.create_seal_key.resource_names
      vpc_id          = step.create_vpc.id
    }
  }

  step "create_backend_cluster" {
    description = global.description.create_backend_cluster
    module      = "backend_${matrix.backend}"
    depends_on = [
      step.create_vault_cluster_backend_targets,
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
      cluster_tag_key = local.backend_tag_key
      hosts           = step.create_vault_cluster_backend_targets.hosts
      license         = (matrix.backend == "consul" && matrix.consul_edition == "ent") ? step.read_backend_license.license : null
      release = {
        edition = matrix.consul_edition
        version = local.consul_version
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
      enos = provider.enos.ubuntu
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
      backend_cluster_name    = step.create_vault_cluster_backend_targets.cluster_name
      backend_cluster_tag_key = local.backend_tag_key
      cluster_name            = step.create_vault_cluster_targets.cluster_name
      config_mode             = "file"
      consul_license          = (matrix.backend == "consul" && matrix.consul_edition == "ent") ? step.read_backend_license.license : null
      consul_release = matrix.backend == "consul" ? {
        edition = matrix.consul_edition
        version = local.consul_version
      } : null
      enable_audit_devices = var.vault_enable_audit_devices
      hosts                = step.create_vault_cluster_targets.hosts
      install_dir          = local.vault_install_dir
      ip_version           = local.ip_version
      license              = matrix.edition != "ce" ? step.read_vault_license.license : null
      local_artifact_path  = local.artifact_path
      packages             = concat(global.packages, global.distro_packages["ubuntu"][global.distro_version["ubuntu"]])
      seal_attributes      = step.create_seal_key.attributes
      seal_type            = local.seal
      storage_backend      = matrix.backend
    }
  }

  // Wait for our cluster to elect a leader
  step "wait_for_leader" {
    description = global.description.wait_for_cluster_to_have_leader
    module      = module.vault_wait_for_leader
    depends_on  = [step.create_vault_cluster]

    providers = {
      enos = provider.enos.ubuntu
    }

    verifies = [
      quality.vault_api_sys_leader_read,
      quality.vault_unseal_ha_leader_election,
    ]

    variables {
      timeout           = 120 // seconds
      ip_version        = local.ip_version
      hosts             = step.create_vault_cluster_targets.hosts
      vault_addr        = step.create_vault_cluster.api_addr_localhost
      vault_install_dir = local.vault_install_dir
      vault_root_token  = step.create_vault_cluster.root_token
    }
  }

  step "test_ui" {
    description = <<-EOF
      Verify that the Vault Web UI test suite can run against a live cluster with the compiled
      assets.
    EOF
    module      = module.vault_test_ui
    depends_on  = [step.wait_for_leader]

    verifies = quality.vault_ui_test

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

  output "cluster_name" {
    description = "The Vault cluster name"
    value       = step.create_vault_cluster.cluster_name
  }

  output "hosts" {
    description = "The Vault cluster target hosts"
    value       = step.create_vault_cluster.hosts
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

  output "seal_name" {
    description = "The Vault cluster seal key name"
    value       = step.create_seal_key.resource_name
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
