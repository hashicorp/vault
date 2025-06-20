// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

scenario "benchmark" {
  description = <<-EOF
    The benchmark scenario creates the required infrastructure to run performance tests against a Vault cluster.
    A three node Vault cluster is created, along with two additional nodes: one to run prometheus and grafana,
    and one to run k6, the load generation tool.

    If you've never used Enos before, it's worth noting that the matrix parameters act as filters. You can view a full
    list of all possible values in enos-globals.hcl. If a matrix parameter listed below is not specified when you run
    the Enos scenario, then all possible parameters are used. For example, if you don't specify `backend` then this
    scenario will run for both raft and consul, because those are the two possible values listed for `backends`.
    Specifying backend:raft will only run the scenario for raft. Specifying backend:consul will create a three node
    Consul cluster and connect it to Vault.

    If you want to run your benchmarks against a released version of Vault, you can download the Vault release tarball
    from releases.hashicorp.com, place it inside the 'support' subdirectory (sometimes you have to create this, as
    it's git-ignored), and specify artifact_source:crt. You also need to add that path to enos-local.vars.hcl as:

    vault_artifact_path = "./support/vault_1.20.0_linux_amd64.zip"

    Substitute in your own actual file name. If you wish to run your benchmarks against your local Vault branch,
    specify artifact_source:local.

    Often times, when running benchmarks, you're wanting to run one very specific scenario, instead of running
    many possible scenarios at once. One example of this from the CLI would be:

    enos scenario launch benchmark config_mode:file artifact_source:crt artifact_type:bundle seal:awskms ip_version:4 consul_version:1.20.6 edition:ent consul_edition:ent backend:consul

    This would run exactly 1 scenario, since I've specified every possible matrix parameter. This is what I used when
    running Consul/IS benchmarks.

    Also note that as of this writing, this scenario does not do automatic benchmarking, results gathering, etc. It's
    manual. That means this scenario is (as of right now) meant to be run as `enos scenario launch` (plus the scenario
    name and matrix parameters, as outlined above) _not_ `enos scenario run`. This also means that when you're done
    running your benchmarks, you need to manually destroy the infrastructure with `enos scenario destroy`.

    If you're going to use an enterprise Consul backend, you'll need to specify the path to a Consul license file under
    the `backend_license_path` variable in the enos-local.vars.hcl file.

    The benchmark module that implements much of the actual benchmark logic has subdirectories for grafana dashboards
    that will get automatically uploaded and installed, as well as k6 templates for different benchmarking scenarios.

    It's worth mentioning that by default we've configured IOPS of the underlying storage volume to 16k. This is below
    the maximum but allows us to create 6 machines and stay below the 100,000 limit in AWS accounts. If you test with
    raft storage you can increase the IOPs to 24k, but beware that consul storage will need to stay at 16k unless you
    get an exception for the account from AWS. When I ran these benchmarks, I requested a quota increase from AWS for
    more IOPs, which is how I was able to run this scenario successfully. You may need to do the same, or if you're not
    benchmarking Consul's raw performance specifically, you can adjust some of the disk parameters in the
    create_vault_cluster_backend_targets to provision less IOPs, or use io1 instead of io2, etc.

    Once the scenario has been launched, and everything has finished, grab the public IP of the metrics node and open
    it in a browser on port 3000. Log into grafana with admin/admin and choose whatever dashboard you wish to see.
    All of the ones in the grafana-dashboards subdirectory will be available. Then SSH into the public IP of the k6
    node and run scenarios via the k6-run.sh shell script, e.g. `./k6-run.sh kvv2`. The argument you pass it should
    match the basename, minus the k6- of the corresponding file in the k6-templates directory, e.g. for the above it
    would match k6-kvv2-js.tpl.

    When you're done getting the results from your grafana dashboard, destroy the infrastructure with `enos destroy`.
  EOF

  // The arch and distro is hardcoded here for expediency. We need to install prometheus, grafana, k6, and
  // the prometheus node exporter. Some of those packages were not available via normal package managers and
  // had to be installed from source. Doing that for 1 arch and 1 distro was sufficient for our needs and a
  // lot easier than doing it for all possible combinations. If you need additional arch/distro combinations,
  // feel free to check out the installation shell scripts in the benchmark module and update them as necessary.
  matrix {
    arch            = ["amd64"]
    artifact_source = global.artifact_sources
    artifact_type   = global.artifact_types
    backend         = global.backends
    config_mode     = global.config_modes
    consul_edition  = global.consul_editions
    consul_version  = global.consul_versions
    distro          = ["ubuntu"]
    edition         = global.editions
    ip_version      = global.ip_versions
    seal            = global.seals

    // Our local builder always creates bundles
    exclude {
      artifact_source = ["local"]
      artifact_type   = ["package"]
    }

    // PKCS#11 can only be used on ent.hsm and ent.hsm.fips1402.
    exclude {
      seal    = ["pkcs11"]
      edition = [for e in matrix.edition : e if !strcontains(e, "hsm")]
    }

    // softhsm packages not available for leap/sles.
    exclude {
      seal   = ["pkcs11"]
      distro = ["leap", "sles"]
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

  step "benchmark_config" {
    description = "Get our configuration for our benchmark modules"
    module      = module.benchmark_config

    variables {
      ports_ingress = values(global.ports)
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

  step "read_backend_license" {
    description = global.description.read_backend_license
    module      = module.read_license
    skip_step   = matrix.backend == "raft" || matrix.consul_edition == "ce"

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

  step "create_seal_key" {
    description = global.description.create_seal_key
    module      = "seal_${matrix.seal}"
    depends_on  = [step.create_vpc]

    providers = {
      enos = provider.enos.ubuntu
    }

    variables {
      cluster_id  = step.create_vpc.id
      common_tags = global.tags
    }
  }

  step "create_k6_target" {
    description = "Create the k6 load generator target machine"
    module      = module.target_ec2_instances
    depends_on  = [step.create_vpc]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }


    variables {
      ami_id          = step.ec2_info.ami_ids[matrix.arch][matrix.distro][global.distro_version[matrix.distro]]
      cluster_tag_key = "benchmark-k6"
      common_tags     = global.tags
      instance_count  = 1
      instance_types  = step.benchmark_config.k6_instance_types
      ports_ingress   = step.benchmark_config.required_ports
      vpc_id          = step.create_vpc.id
    }
  }

  step "create_metrics_collector_target" {
    description = "Create the benchmark metrics collector target machine"
    module      = module.target_ec2_instances
    depends_on  = [step.create_vpc]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      ami_id          = step.ec2_info.ami_ids[matrix.arch][matrix.distro][global.distro_version[matrix.distro]]
      cluster_tag_key = "benchmark-collector"
      common_tags     = global.tags
      instance_count  = 1
      instance_types  = step.benchmark_config.metrics_instance_types
      ports_ingress   = step.benchmark_config.required_ports
      vpc_id          = step.create_vpc.id
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
      ami_id           = step.ec2_info.ami_ids[matrix.arch][matrix.distro][global.distro_version[matrix.distro]]
      cluster_tag_key  = global.vault_tag_key
      common_tags      = global.tags
      ebs_optimized    = true
      instance_count   = 3
      instance_types   = step.benchmark_config.vault_node_instance_types
      ports_ingress    = step.benchmark_config.required_ports
      root_volume_type = "io2"
      root_volume_size = 24
      root_volume_iops = step.benchmark_config.storage_disk_iops
      seal_key_names   = step.create_seal_key.resource_names
      vpc_id           = step.create_vpc.id
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
      ami_id          = step.ec2_info.ami_ids[matrix.arch][matrix.distro][global.distro_version[matrix.distro]]
      cluster_tag_key = global.backend_tag_key
      common_tags     = global.tags
      ebs_optimized   = true
      instance_count  = 3
      instance_types = {
        amd64 = "i3.4xlarge"
        arm64 = "t4g.small"
      }
      ports_ingress    = step.benchmark_config.required_ports
      root_volume_type = "io2"
      root_volume_size = 24
      root_volume_iops = step.benchmark_config.storage_disk_iops
      seal_key_names   = step.create_seal_key.resource_names
      vpc_id           = step.create_vpc.id
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
      step.create_vault_cluster_targets,
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
      backend_cluster_name    = step.create_vault_cluster_backend_targets.cluster_name
      backend_cluster_tag_key = global.backend_tag_key
      cluster_name            = step.create_vault_cluster_targets.cluster_name
      config_mode             = matrix.config_mode
      consul_license          = (matrix.backend == "consul" && matrix.consul_edition == "ent") ? step.read_backend_license.license : null
      consul_release = matrix.backend == "consul" ? {
        edition = matrix.consul_edition
        version = matrix.consul_version
      } : null
      enable_audit_devices = false
      enable_telemetry     = true
      hosts                = step.create_vault_cluster_targets.hosts
      install_dir          = global.vault_install_dir[matrix.artifact_type]
      ip_version           = matrix.ip_version
      license              = matrix.edition != "ce" ? step.read_vault_license.license : null
      local_artifact_path  = local.artifact_path
      manage_service       = local.manage_service
      packages             = concat(global.packages, global.distro_packages[matrix.distro][global.distro_version[matrix.distro]])
      seal_attributes      = step.create_seal_key.attributes
      seal_type            = matrix.seal
      storage_backend      = matrix.backend
    }
  }

  step "get_local_metadata" {
    description = global.description.get_local_metadata
    skip_step   = matrix.artifact_source != "local"
    module      = module.get_local_metadata
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
      timeout           = 120 // seconds
      ip_version        = matrix.ip_version
      hosts             = step.create_vault_cluster_targets.hosts
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

  step "benchmark_setup" {
    module = module.benchmark_setup
    depends_on = [
      step.create_metrics_collector_target,
      step.create_k6_target,
      step.verify_vault_unsealed,
      step.get_vault_cluster_ips,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      consul_hosts                     = matrix.backend == "consul" ? step.create_vault_cluster_backend_targets.hosts : {}
      grafana_version                  = step.benchmark_config.grafana_version
      grafana_http_port                = step.benchmark_config.grafana_http_port
      k6_host                          = step.create_k6_target.hosts[0]
      leader_addr                      = step.get_vault_cluster_ips.leader_private_ip
      metrics_host                     = step.create_metrics_collector_target.hosts[0]
      prometheus_node_exporter_version = step.benchmark_config.prometheus_node_exporter_version
      prometheus_version               = step.benchmark_config.prometheus_version
      vault_hosts                      = step.create_vault_cluster_targets.hosts
      vault_token                      = step.create_vault_cluster.root_token
      vpc_id                           = step.create_vpc.id
    }
  }

  output "audit_device_file_path" {
    description = "The file path for the file audit device, if enabled"
    value       = step.create_vault_cluster.audit_device_file_path
  }

  output "dashboard_url" {
    description = "The URL for viewing the dashboard in grafana"
    value       = step.benchmark_setup.dashboard_url
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

  output "seal_key_attributes" {
    description = "The Vault cluster seal attributes"
    value       = step.create_seal_key.attributes
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
