// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

scenario "dev_pr_replication" {
  description = <<-EOF
    This scenario spins up a two Vault clusters with either an external Consul cluster or
    integrated Raft for storage. The secondary cluster is configured with performance replication
    from the primary cluster. None of our test verification is included in this scenario in order
    to improve end-to-end speed. If you wish to perform such verification you'll need to use a
    non-dev scenario.

    The scenario supports finding and installing any released 'linux/amd64' or 'linux/arm64' Vault
    artifact as long as its version is >= 1.8. You can also use the 'artifact:local' variant to
    build and deploy the current branch!

    In order to execute this scenario you'll need to install the enos CLI:
      - $ brew tap hashicorp/tap && brew update && brew install hashicorp/tap/enos
    
    You'll also need access to an AWS account via Doormat, follow the guide here:
      https://eng-handbook.hashicorp.services/internal-tools/enos/getting-started/#authenticate-to-aws-with-doormat

    Follow this guide to get an SSH keypair set up in the AWS account:
      https://eng-handbook.hashicorp.services/internal-tools/enos/getting-started/#set-your-aws-key-pair-name-and-private-key

    Please note that this scenario requires several inputs variables to be set in order to function
    properly. While not all variants will require all variables, it's suggested that you look over
    the scenario outline to determine which variables affect which steps and which have inputs that
    you should set. You can use the following command to get a textual outline of the entire
    scenario:
      enos scenario outline dev_pr_replication

    You can also create an HTML version that is suitable for viewing in web browsers:
      enos scenario outline dev_pr_replication --format html > index.html
      open index.html

    To configure the required variables you have a couple of choices. You can create an
    'enos-local.vars' file in the same 'enos' directory where this scenario is defined. In it you
    declare your desired variable values. For example, you could copy the following content and
    then set the values as necessary:

    artifactory_username      = "username@hashicorp.com"
    artifactory_token         = "<ARTIFACTORY TOKEN VALUE>
    aws_region                = "us-west-2"
    aws_ssh_keypair_name      = "<YOUR REGION SPECIFIC KEYPAIR NAME>"
    aws_ssh_keypair_key_path  = "/path/to/your/private/key.pem"
    dev_build_local_ui        = false
    dev_consul_version        = "1.18.1"
    vault_license_path        = "./support/vault.hclic"
    vault_product_version     = "1.16.2"

    Alternatively, you can set them in your environment:
    export ENOS_VAR_aws_region="us-west-2"
    export ENOS_VAR_vault_license_path="./support/vault.hclic"

    After you've configured your inputs you can list and filter the available scenarios and then
    subsequently launch and destroy them.
      enos scenario list --help
      enos scenario launch --help
      enos scenario list dev_pr_replication
      enos scenario launch dev_pr_replication arch:amd64 artifact:deb distro:ubuntu edition:ent.hsm primary_backend:raft primary_seal:awskms secondary_backend:raft secondary_seal:pkcs11

    When the scenario is finished launching you refer to the scenario outputs to see information
    related to your cluster. You can use this information to SSH into nodes and/or to interact
    with vault.
      enos scenario output dev_pr_replication arch:amd64 artifact:deb distro:ubuntu edition:ent.hsm primary_backend:raft primary_seal:awskms secondary_backend:raft secondary_seal:pkcs11
      ssh -i /path/to/your/private/key.pem <PUBLIC_IP>
      vault status

    After you've finished you can tear down the cluster
      enos scenario destroy dev_pr_replication arch:amd64 artifact:deb distro:ubuntu edition:ent.hsm primary_backend:raft primary_seal:awskms secondary_backend:raft secondary_seal:pkcs11
  EOF

  // The matrix is where we define all the baseline combinations that enos can utilize to customize
  // your scenario. By default enos attempts to perform your command on the entire product of these
  // possible comginations! Most of the time you'll want to reduce that by passing in a filter.
  // Run 'enos scenario list --help' to see more about how filtering scenarios works in enos.
  matrix {
    arch              = ["amd64", "arm64"]
    artifact          = ["local", "deb", "rpm", "zip"]
    distro            = ["amzn", "leap", "rhel", "sles", "ubuntu"]
    edition           = ["ent", "ent.fips1402", "ent.hsm", "ent.hsm.fips1402"]
    primary_backend   = ["consul", "raft"]
    primary_seal      = ["awskms", "pkcs11", "shamir"]
    secondary_backend = ["consul", "raft"]
    secondary_seal    = ["awskms", "pkcs11", "shamir"]

    exclude {
      edition = ["ent.hsm", "ent.fips1402", "ent.hsm.fips1402"]
      arch    = ["arm64"]
    }

    exclude {
      artifact = ["rpm"]
      distro   = ["ubuntu"]
    }

    exclude {
      artifact = ["deb"]
      distro   = ["rhel"]
    }

    exclude {
      primary_seal = ["pkcs11"]
      edition      = ["ce", "ent", "ent.fips1402"]
    }

    exclude {
      secondary_seal = ["pkcs11"]
      edition        = ["ce", "ent", "ent.fips1402"]
    }
  }

  // Specify which Terraform configs and providers to use in this scenario. Most of the time you'll
  // never need to change this! If you wanted to test with different terraform or terraform CLI
  // settings you can define them and assign them here.
  terraform_cli = terraform_cli.default
  terraform     = terraform.default

  // Here we declare all of the providers that we might need for our scenario.
  providers = [
    provider.aws.default,
    provider.enos.ec2_user,
    provider.enos.ubuntu
  ]

  // These are variable values that are local to our scenario. They are evaluated after external
  // variables and scenario matrices but before any of our steps.
  locals {
    // The enos provider uses different ssh transport configs for different distros (as
    // specified in enos-providers.hcl), and we need to be able to access both of those here.
    enos_provider = {
      amzn   = provider.enos.ec2_user
      leap   = provider.enos.ec2_user
      rhel   = provider.enos.ec2_user
      sles   = provider.enos.ec2_user
      ubuntu = provider.enos.ubuntu
    }
    // We install vault packages from artifactory. If you wish to use one of these variants you'll
    // need to configure your artifactory credentials.
    use_artifactory = matrix.artifact == "deb" || matrix.artifact == "rpm"
    // The IP version to use for the Vault listener and associated things.
    ip_version = 4
    // Zip bundles and local builds don't come with systemd units or any associated configuration.
    // When this is true we'll let enos handle this for us.
    manage_service = matrix.artifact == "zip" || matrix.artifact == "local"
    // If you are using an ent edition, you will need a Vault license. Common convention
    // is to store it at ./support/vault.hclic, but you may change this path according
    // to your own preference.
    vault_install_dir = matrix.artifact == "zip" || matrix.artifact == "local" ? global.vault_install_dir["bundle"] : global.vault_install_dir["package"]
  }

  // Begin scenario steps. These are the steps we'll perform to get your cluster up and running.
  step "build_or_find_vault_artifact" {
    description = <<-EOF
      Depending on how we intend to get our Vault artifact, this step either builds vault from our
      current branch or finds debian or redhat packages in Artifactory. If we're using a zip bundle
      we'll get it from releases.hashicorp.com and skip this step entirely. Please note that if you
      wish to use a deb or rpm artifact you'll have to configure your artifactory credentials!

      Variables that are used in this step:

        artifactory_host:
          The artifactory host to search. It's very unlikely that you'll want to change this. The
          default value is the HashiCorp Artifactory instance.
        artifactory_repo:
          The artifactory host to search. It's very unlikely that you'll want to change this. The
          default value is where CRT will publish packages.
        artifactory_username:
          The artifactory username associated with your token. You'll need this if you wish to use
          deb or rpm artifacts! You can request access via Okta.
        artifactory_token:
          The artifactory token associated with your username. You'll need this if you wish to use
          deb or rpm artifacts! You can create a token by logging into Artifactory via Okta.
        dev_build_local_ui:
          If you are not testing any changes in the UI, set to false. This will save time by not
          building the entire UI. If you need to test the UI, set to true.
        vault_product_version:
          When using the artifact:rpm or artifact:deb variants we'll use this variable to determine
          which version of the Vault pacakge we should fetch from Artifactory.
        vault_artifact_path:
          When using the artifact:local variant we'll utilize this variable to determine where
          to create the vault.zip archive from the local branch. Default: to /tmp/vault.zip.
        vault_local_tags:
          When using the artifact:local variant we'll use this variable to inject custom build
          tags. If left unset we'll automatically use the build tags that correspond to the edition
          variant.
    EOF
    module      = matrix.artifact == "local" ? "build_local" : local.use_artifactory ? "build_artifactory_package" : "build_crt"

    variables {
      // Used for all modules
      arch            = matrix.arch
      edition         = matrix.edition
      product_version = var.vault_product_version
      // Required for the local build which will always result in using a local zip bundle
      artifact_path = matrix.artifact == "local" ? abspath(var.vault_artifact_path) : null
      build_ui      = var.dev_build_local_ui
      build_tags    = var.vault_local_build_tags != null ? var.vault_local_build_tags : global.build_tags[matrix.edition]
      goarch        = matrix.arch
      goos          = "linux"
      // Required when using a RPM or Deb package
      // Some of these variables don't have default values so we'll only set them if they are
      // required.
      artifactory_host     = local.use_artifactory ? var.artifactory_host : null
      artifactory_repo     = local.use_artifactory ? var.artifactory_repo : null
      artifactory_username = local.use_artifactory ? var.artifactory_username : null
      artifactory_token    = local.use_artifactory ? var.artifactory_token : null
      distro               = matrix.distro
    }
  }

  step "ec2_info" {
    description = "This discovers usefull metadata in Ec2 like AWS AMI IDs that we use in later modules."
    module      = module.ec2_info
  }

  step "create_vpc" {
    description = <<-EOF
      Create the VPC resources required for our scenario.

        Variables that are used in this step:
          tags:
            If you wish to add custom tags to taggable resources in AWS you can set the 'tags' variable
            and they'll be added to resources when possible.
    EOF
    module      = module.create_vpc
    depends_on  = [step.ec2_info]

    variables {
      common_tags = global.tags
    }
  }

  step "read_backend_license" {
    description = <<-EOF
      Read the contents of the backend license if we're using a Consul backend for either cluster
      and the backend_edition variable is set to "ent".

      Variables that are used in this step:
        backend_edition:
          The edition of Consul to use. If left unset it will default to CE.
        backend_license_path:
          If this variable is set we'll use it to determine the local path on disk that contains a
          Consul Enterprise license. If it is not set we'll attempt to load it from
          ./support/consul.hclic.
    EOF
    skip_step   = (var.backend_edition == "ce" || var.backend_edition == "oss") || (matrix.primary_backend == "raft" && matrix.secondary_backend == "raft")
    module      = module.read_license

    variables {
      file_name = global.backend_license_path
    }
  }

  step "read_vault_license" {
    description = <<-EOF
      Validates and reads into memory the contents of a local Vault Enterprise license if we're
      using an Enterprise edition. This step does not run when using a community edition of Vault.

      Variables that are used in this step:
        vault_license_path:
          If this variable is set we'll use it to determine the local path on disk that contains a
          Vault Enterprise license. If it is not set we'll attempt to load it from
          ./support/vault.hclic.
    EOF
    module      = module.read_license

    variables {
      file_name = global.vault_license_path
    }
  }

  step "create_primary_seal_key" {
    description = <<-EOF
      Create the necessary seal keys depending on our configured seal.

      Variables that are used in this step:
        tags:
          If you wish to add custom tags to taggable resources in AWS you can set the 'tags' variable
          and they'll be added to resources when possible.
    EOF
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
    description = <<-EOF
      Create the necessary seal keys depending on our configured seal.

      Variables that are used in this step:
        tags:
          If you wish to add custom tags to taggable resources in AWS you can set the 'tags' variable
          and they'll be added to resources when possible.
    EOF
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

  step "create_primary_cluster_targets" {
    description = <<-EOF
      Creates the necessary machine infrastructure targets for the Vault cluster. We also ensure
      that the firewall is configured to allow the necessary Vault and Consul traffic and SSH
      from the machine executing the Enos scenario.

      Variables that are used in this step:
        aws_ssh_keypair_name:
          The AWS SSH Keypair name to use for target machines.
        project_name:
          The project name is used for additional tag metadata on resources.
        tags:
          If you wish to add custom tags to taggable resources in AWS you can set the 'tags' variable
          and they'll be added to resources when possible.
        vault_instance_count:
          How many instances to provision for the Vault cluster. If left unset it will use a default
          of three.
    EOF
    module      = module.target_ec2_instances
    depends_on  = [step.create_vpc]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      ami_id          = step.ec2_info.ami_ids[matrix.arch][matrix.distro][global.distro_version[matrix.distro]]
      cluster_tag_key = global.vault_tag_key
      common_tags     = global.tags
      instance_count  = try(var.vault_instance_count, 3)
      seal_key_names  = step.create_primary_seal_key.resource_names
      vpc_id          = step.create_vpc.id
    }
  }

  step "create_primary_cluster_backend_targets" {
    description = <<-EOF
      Creates the necessary machine infrastructure targets for the backend Consul storage cluster.
      We also ensure that the firewall is configured to allow the necessary Consul traffic and SSH
      from the machine executing the Enos scenario. When using integrated storage this step is a
      no-op that does nothing.

      Variables that are used in this step:
        tags:
          If you wish to add custom tags to taggable resources in AWS you can set the 'tags' variable
          and they'll be added to resources when possible.
        project_name:
          The project name is used for additional tag metadata on resources.
        aws_ssh_keypair_name:
          The AWS SSH Keypair name to use for target machines.
    EOF
    module      = matrix.primary_backend == "consul" ? module.target_ec2_instances : module.target_ec2_shim
    depends_on  = [step.create_vpc]

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

  step "create_secondary_cluster_targets" {
    description = <<-EOF
      Creates the necessary machine infrastructure targets for the Vault cluster. We also ensure
      that the firewall is configured to allow the necessary Vault and Consul traffic and SSH
      from the machine executing the Enos scenario.
    EOF
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
    description = <<-EOF
      Creates the necessary machine infrastructure targets for the backend Consul storage cluster.
      We also ensure that the firewall is configured to allow the necessary Consul traffic and SSH
      from the machine executing the Enos scenario. When using integrated storage this step is a
      no-op that does nothing.
    EOF

    module     = matrix.secondary_backend == "consul" ? module.target_ec2_instances : module.target_ec2_shim
    depends_on = [step.create_vpc]

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
    description = <<-EOF
      Install, configure, and start the backend Consul storage cluster for the primary Vault Cluster.
      When we are using the raft storage variant this step is a no-op.

      Variables that are used in this step:
        backend_edition:
          When configured with the backend:consul variant we'll utilize this variable to determine
          the edition of Consul to use for the cluster. Note that if you set it to 'ent' you will
          also need a valid license configured for the read_backend_license step. Default: ce.
        dev_consul_version:
          When configured with the backend:consul variant we'll utilize this variable to determine
          the version of Consul to use for the cluster.
    EOF
    module      = "backend_${matrix.primary_backend}"
    depends_on = [
      step.create_primary_cluster_backend_targets
    ]

    providers = {
      enos = provider.enos.ubuntu
    }

    variables {
      cluster_name    = step.create_primary_cluster_backend_targets.cluster_name
      cluster_tag_key = global.backend_tag_key
      hosts           = step.create_primary_cluster_backend_targets.hosts
      license         = matrix.primary_backend == "consul" ? step.read_backend_license.license : null
      release = {
        edition = var.backend_edition
        version = var.dev_consul_version
      }
    }
  }

  step "create_primary_cluster" {
    description = <<-EOF
      Install, configure, start, initialize and unseal the primary Vault cluster on the specified
      target instances.

      Variables that are used in this step:
      backend_edition:
        When configured with the backend:consul variant we'll utilize this variable to determine
        which version of the consul client to install on each node for Consul storage. Note that
        if you set it to 'ent' you will also need a valid license configured for the
        read_backend_license step. If left unset we'll use an unlicensed CE version.
      dev_config_mode:
        You can set this variable to instruct enos on how to primarily configure Vault when starting
        the service. Options are 'file' and 'env' for configuration file or environment variables.
        If left unset we'll use the default value.
      dev_consul_version:
        When configured with the backend:consul variant we'll utilize this variable to determine
        which version of Consul to install. If left unset we'll utilize the default value.
      vault_artifact_path:
        When using the artifact:local variant this variable is utilized to specify where on
        the local disk the vault.zip file we've built is located. It can be left unset to use
        the default value.
      vault_enable_audit_devices:
        Whether or not to enable various audit devices after unsealing the Vault cluster. By default
        we'll configure syslog, socket, and file auditing.
      vault_product_version:
        When using the artifact:zip variant this variable is utilized to specify the version of
        Vault to download from releases.hashicorp.com.
    EOF
    module      = module.vault_cluster
    depends_on = [
      step.create_primary_backend_cluster,
      step.create_primary_cluster_targets,
      step.build_or_find_vault_artifact,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      // We set vault_artifactory_release when we want to get a .deb or .rpm package from Artifactory.
      // We set vault_release when we want to get a .zip bundle from releases.hashicorp.com
      // We only set one or the other, never both.
      artifactory_release     = local.use_artifactory ? step.build_or_find_vault_artifact.release : null
      backend_cluster_name    = step.create_primary_cluster_backend_targets.cluster_name
      backend_cluster_tag_key = global.backend_tag_key
      cluster_name            = step.create_primary_cluster_targets.cluster_name
      config_mode             = var.dev_config_mode
      consul_license          = matrix.primary_backend == "consul" ? step.read_backend_license.license : null
      consul_release = matrix.primary_backend == "consul" ? {
        edition = var.backend_edition
        version = var.dev_consul_version
      } : null
      enable_audit_devices = var.vault_enable_audit_devices
      hosts                = step.create_primary_cluster_targets.hosts
      install_dir          = local.vault_install_dir
      ip_version           = local.ip_version
      license              = step.read_vault_license.license
      local_artifact_path  = matrix.artifact == "local" ? abspath(var.vault_artifact_path) : null
      manage_service       = local.manage_service
      packages             = concat(global.packages, global.distro_packages[matrix.distro][global.distro_version[matrix.distro]])
      release              = matrix.artifact == "zip" ? { version = var.vault_product_version, edition = matrix.edition } : null
      seal_attributes      = step.create_primary_seal_key.attributes
      seal_type            = matrix.primary_seal
      storage_backend      = matrix.primary_backend
    }
  }

  step "create_secondary_backend_cluster" {
    description = <<-EOF
      Install, configure, and start the backend Consul storage cluster for the primary Vault Cluster.
      When we are using the raft storage variant this step is a no-op.

      Variables that are used in this step:
        backend_edition:
          When configured with the backend:consul variant we'll utilize this variable to determine
          the edition of Consul to use for the cluster. Note that if you set it to 'ent' you will
          also need a valid license configured for the read_backend_license step. Default: ce.
        dev_consul_version:
          When configured with the backend:consul variant we'll utilize this variable to determine
          the version of Consul to use for the cluster.
    EOF
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
      license         = matrix.secondary_backend == "consul" ? step.read_backend_license.license : null
      release = {
        edition = var.backend_edition
        version = var.dev_consul_version
      }
    }
  }

  step "create_secondary_cluster" {
    description = <<-EOF
      Install, configure, start, initialize and unseal the secondary Vault cluster on the specified
      target instances.

      Variables that are used in this step:
      backend_edition:
        When configured with the backend:consul variant we'll utilize this variable to determine
        which version of the consul client to install on each node for Consul storage. Note that
        if you set it to 'ent' you will also need a valid license configured for the
        read_backend_license step. If left unset we'll use an unlicensed CE version.
      dev_config_mode:
        You can set this variable to instruct enos on how to primarily configure Vault when starting
        the service. Options are 'file' and 'env' for configuration file or environment variables.
        If left unset we'll use the default value.
      dev_consul_version:
        When configured with the backend:consul variant we'll utilize this variable to determine
        which version of Consul to install. If left unset we'll utilize the default value.
      vault_artifact_path:
        When using the artifact:local variant this variable is utilized to specify where on
        the local disk the vault.zip file we've built is located. It can be left unset to use
        the default value.
      vault_enable_audit_devices:
        Whether or not to enable various audit devices after unsealing the Vault cluster. By default
        we'll configure syslog, socket, and file auditing.
      vault_product_version:
        When using the artifact:zip variant this variable is utilized to specify the version of
        Vault to download from releases.hashicorp.com.
    EOF
    module      = module.vault_cluster
    depends_on = [
      step.create_secondary_backend_cluster,
      step.create_secondary_cluster_targets
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      // We set vault_artifactory_release when we want to get a .deb or .rpm package from Artifactory.
      // We set vault_release when we want to get a .zip bundle from releases.hashicorp.com
      // We only set one or the other, never both.
      artifactory_release     = local.use_artifactory ? step.build_or_find_vault_artifact.release : null
      backend_cluster_name    = step.create_secondary_cluster_backend_targets.cluster_name
      backend_cluster_tag_key = global.backend_tag_key
      cluster_name            = step.create_secondary_cluster_targets.cluster_name
      config_mode             = var.dev_config_mode
      consul_license          = matrix.secondary_backend == "consul" ? step.read_backend_license.license : null
      consul_release = matrix.secondary_backend == "consul" ? {
        edition = var.backend_edition
        version = var.dev_consul_version
      } : null
      enable_audit_devices = var.vault_enable_audit_devices
      hosts                = step.create_secondary_cluster_targets.hosts
      install_dir          = local.vault_install_dir
      ip_version           = local.ip_version
      license              = step.read_vault_license.license
      local_artifact_path  = matrix.artifact == "local" ? abspath(var.vault_artifact_path) : null
      manage_service       = local.manage_service
      packages             = concat(global.packages, global.distro_packages[matrix.distro][global.distro_version[matrix.distro]])
      release              = matrix.artifact == "zip" ? { version = var.vault_product_version, edition = matrix.edition } : null
      seal_attributes      = step.create_secondary_seal_key.attributes
      seal_type            = matrix.secondary_seal
      storage_backend      = matrix.secondary_backend
    }
  }

  step "verify_that_vault_primary_cluster_is_unsealed" {
    description = <<-EOF
      Wait for the for the primary cluster to unseal and reach a healthy state.
    EOF
    module      = module.vault_wait_for_cluster_unsealed
    depends_on = [
      step.create_primary_cluster
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      hosts             = step.create_primary_cluster_targets.hosts
      vault_addr        = step.create_primary_cluster.api_addr_localhost
      vault_install_dir = local.vault_install_dir
    }
  }

  step "verify_that_vault_secondary_cluster_is_unsealed" {
    description = <<-EOF
      Wait for the for the secondary cluster to unseal and reach a healthy state.
    EOF
    module      = module.vault_wait_for_cluster_unsealed
    depends_on = [
      step.create_secondary_cluster
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      hosts             = step.create_secondary_cluster_targets.hosts
      vault_addr        = step.create_secondary_cluster.api_addr_localhost
      vault_install_dir = local.vault_install_dir
    }
  }

  step "get_primary_cluster_ips" {
    description = <<-EOF
      Determine which node is the primary and which are followers and map their private IP address
      to their public IP address. We'll use this information so that we can enable performance
      replication on the leader.
    EOF
    module      = module.vault_get_cluster_ips
    depends_on  = [step.verify_that_vault_primary_cluster_is_unsealed]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      hosts             = step.create_primary_cluster_targets.hosts
      ip_version        = local.ip_version
      vault_addr        = step.create_primary_cluster.api_addr_localhost
      vault_install_dir = local.vault_install_dir
      vault_root_token  = step.create_primary_cluster.root_token
    }
  }

  step "get_secondary_cluster_ips" {
    description = <<-EOF
      Determine which node is the primary and which are followers and map their private IP address
      to their public IP address. We'll use this information so that we can enable performance
      replication on the leader.
    EOF
    module      = module.vault_get_cluster_ips
    depends_on  = [step.verify_that_vault_secondary_cluster_is_unsealed]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      hosts             = step.create_secondary_cluster_targets.hosts
      ip_version        = local.ip_version
      vault_addr        = step.create_secondary_cluster.api_addr_localhost
      vault_install_dir = local.vault_install_dir
      vault_root_token  = step.create_secondary_cluster.root_token
    }
  }

  step "setup_userpass_for_replication_auth" {
    description = <<-EOF
      Enable the auth userpass method and create a new user.
    EOF
    module      = module.vault_verify_secrets_engines_create
    depends_on  = [step.get_primary_cluster_ips]


    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      hosts             = step.create_primary_cluster_targets.hosts
      leader_host       = step.get_primary_cluster_ips.leader_host
      vault_addr        = step.create_primary_cluster.api_addr_localhost
      vault_install_dir = local.vault_install_dir
      vault_root_token  = step.create_primary_cluster.root_token
    }
  }

  step "configure_performance_replication_primary" {
    description = <<-EOF
      Create a superuser policy write it for our new user. Activate performance replication on
      the primary.
    EOF
    module      = module.vault_setup_perf_primary
    depends_on = [
      step.get_primary_cluster_ips,
      step.get_secondary_cluster_ips,
      step.setup_userpass_for_replication_auth,
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      ip_version          = local.ip_version
      primary_leader_host = step.get_primary_cluster_ips.leader_host
      vault_addr          = step.create_primary_cluster.api_addr_localhost
      vault_install_dir   = local.vault_install_dir
      vault_root_token    = step.create_primary_cluster.root_token
    }
  }

  step "generate_secondary_token" {
    description = <<-EOF
      Create a random token and write it to sys/replication/performance/primary/secondary-token on
      the primary.
    EOF
    module      = module.generate_secondary_token
    depends_on  = [step.configure_performance_replication_primary]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      ip_version          = local.ip_version
      primary_leader_host = step.get_primary_cluster_ips.leader_host
      replication_type    = "performance"
      vault_addr          = step.create_primary_cluster.api_addr_localhost
      vault_install_dir   = local.vault_install_dir
      vault_root_token    = step.create_primary_cluster.root_token
    }
  }

  step "configure_performance_replication_secondary" {
    description = <<-EOF
      Enable performance replication on the secondary using the new shared token.
    EOF
    module      = module.vault_setup_replication_secondary
    depends_on  = [step.generate_secondary_token]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      ip_version            = local.ip_version
      secondary_leader_host = step.get_secondary_cluster_ips.leader_host
      replication_type      = "performance"
      vault_addr            = step.create_secondary_cluster.api_addr_localhost
      vault_install_dir     = local.vault_install_dir
      vault_root_token      = step.create_secondary_cluster.root_token
      wrapping_token        = step.generate_secondary_token.secondary_token
    }
  }

  step "unseal_secondary_followers" {
    description = <<-EOF
      After replication is enabled we need to unseal the followers on the secondary cluster.
      Depending on how we're configured we'll pass the unseal keys according to this guide:
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
      vault_install_dir = local.vault_install_dir
      vault_unseal_keys = matrix.primary_seal == "shamir" ? step.create_primary_cluster.unseal_keys_hex : step.create_primary_cluster.recovery_keys_hex
      vault_seal_type   = matrix.primary_seal == "shamir" ? matrix.primary_seal : matrix.secondary_seal
    }
  }

  step "verify_secondary_cluster_is_unsealed_after_enabling_replication" {
    description = <<-EOF
      Verify that the secondary cluster is unsealed after we enable PR replication.
    EOF
    module      = module.vault_wait_for_cluster_unsealed
    depends_on = [
      step.unseal_secondary_followers
    ]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      hosts             = step.create_secondary_cluster_targets.hosts
      vault_addr        = step.create_primary_cluster.api_addr_localhost
      vault_install_dir = local.vault_install_dir
    }
  }

  step "verify_performance_replication" {
    description = <<-EOF
      Check sys/replication/performance/status and ensure that all nodes are in the correct state
      after enabling performance replication.
    EOF
    module      = module.vault_verify_performance_replication
    depends_on  = [step.verify_secondary_cluster_is_unsealed_after_enabling_replication]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      ip_version            = local.ip_version
      primary_leader_host   = step.get_primary_cluster_ips.leader_host
      secondary_leader_host = step.get_secondary_cluster_ips.leader_host
      vault_addr            = step.create_primary_cluster.api_addr_localhost
      vault_install_dir     = local.vault_install_dir
    }
  }

  // When using a Consul backend, these output values will be for the Consul backend.
  // When using a Raft backend, these output values will be null.
  output "audit_device_file_path" {
    description = "The file path for the file audit device, if enabled"
    value       = step.create_primary_cluster.audit_device_file_path
  }

  output "primary_cluster_hosts" {
    description = "The Vault primary cluster target hosts"
    value       = step.create_primary_cluster_targets.hosts
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

  output "performance_secondary_token" {
    description = "The performance secondary replication token"
    value       = step.generate_secondary_token.secondary_token
  }
}
