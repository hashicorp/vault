scenario "package_manager" {
  matrix {
    arch            = ["amd64", "arm64"]
    backend         = ["consul", "raft"]
    artifact_source = ["local", "crt", "artifactory"]
    consul_version  = ["1.13.2"]
    distro          = ["ubuntu", "rhel"]
    edition         = ["ent"]
    seal            = ["awskms", "shamir"]
  }

  terraform_cli = terraform_cli.default
  terraform     = terraform.default
  providers = [
    provider.aws.default,
    provider.enos.ubuntu,
    provider.enos.rhel
  ]

  locals {
    artifact_type = "package"
    build_tags = {
      "ent" = ["enterprise", "ent"]
    }
    bundle_path             = matrix.artifact_source != "artifactory" ? abspath(var.vault_bundle_path) : null
    dependencies_to_install = ["jq"]
    enos_provider = {
      rhel   = provider.enos.rhel
      ubuntu = provider.enos.ubuntu
    }
    install_artifactory_artifact = var.vault_revision != null && local.bundle_path == null
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
    vault_install_dir = {
      rhel   = "/bin"
      ubuntu = "/usr/bin"
    }
  }

  step "build_vault" {
    module = "build_${matrix.artifact_source}"

    variables {
      build_tags           = try(var.vault_local_build_tags, local.build_tags[matrix.edition])
      bundle_path          = local.bundle_path
      goarch               = matrix.arch
      goos                 = "linux"
      artifactory_host     = matrix.artifact_source == "artifactory" ? var.artifactory_host : null
      artifactory_repo     = matrix.artifact_source == "artifactory" ? var.artifactory_repo : null
      artifactory_username = matrix.artifact_source == "artifactory" ? var.artifactory_username : null
      artifactory_token    = matrix.artifact_source == "artifactory" ? var.artifactory_token : null
      arch                 = matrix.artifact_source == "artifactory" ? matrix.arch : null
      product_version      = matrix.artifact_source == "artifactory" ? var.vault_product_version : null
      artifact_type        = local.artifact_type
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

  # Follow the current structure/modules of existing enos scenarios
  step "create_vpc" {
    module = module.create_vpc

    variables {
      ami_architectures  = [matrix.arch]
      availability_zones = step.find_azs.availability_zones
      common_tags        = local.tags
    }
  }

  # Replaces module "verify_license" from original smoke test
  step "read_license" {
    skip_step = matrix.edition == "oss"
    module    = module.read_license

    variables {
      file_name = abspath(joinpath(path.root, "./support/vault.hclic"))
    }
  }

  step "create_backend_cluster" {
    module     = "backend_${matrix.backend}"
    depends_on = [step.create_vpc]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      ami_id      = step.create_vpc.ami_ids[matrix.distro][matrix.arch]
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
      vault_artifactory_release = matrix.artifact_source == "artifactory" ? step.build_vault.vault_artifactory_release : null
      common_tags               = local.tags
      consul_cluster_tag        = step.create_backend_cluster.consul_cluster_tag
      dependencies_to_install   = local.dependencies_to_install
      instance_type             = local.vault_instance_type
      kms_key_arn               = step.create_vpc.kms_key_arn
      manage_service            = false
      storage_backend           = matrix.backend
      vault_install_dir         = local.vault_install_dir[matrix.distro]
      unseal_method             = matrix.seal
      vault_local_artifact_path = local.bundle_path
      vault_license             = matrix.edition != "oss" ? step.read_license.license : null
      vpc_id                    = step.create_vpc.vpc_id
    }
  }

  # Verify version
  step "verify_vault_version" {
    module     = module.vault_verify_version
    depends_on = [step.create_vault_cluster]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_install_dir = local.vault_install_dir[matrix.distro]

      vault_instances       = step.create_vault_cluster.vault_instances
      vault_edition         = matrix.edition
      vault_product_version = var.vault_product_version
      vault_revision        = var.vault_revision
      vault_build_date      = var.vault_build_date
      vault_root_token      = step.create_vault_cluster.vault_root_token
    }
  }

  step "verify_replication" {
    module     = module.vault_verify_replication
    depends_on = [step.create_vault_cluster]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_edition     = matrix.edition
      vault_instances   = step.create_vault_cluster.vault_instances
      vault_install_dir = local.vault_install_dir[matrix.distro]
    }
  }

  step "verify_ui" {
    module     = module.vault_verify_ui
    depends_on = [step.create_vault_cluster]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_instances   = step.create_vault_cluster.vault_instances
      vault_install_dir = local.vault_install_dir[matrix.distro]
    }
  }

  step "verify_write_test_data" {
    module     = module.vault_verify_write_test_data
    depends_on = [step.create_vault_cluster]

    providers = {
      enos = local.enos_provider[matrix.distro]
    }

    variables {
      vault_root_token  = step.create_vault_cluster.vault_root_token
      vault_instances   = step.create_vault_cluster.vault_instances
      vault_install_dir = local.vault_install_dir[matrix.distro]
    }
  }

}
