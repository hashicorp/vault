scenario "k8s" {
  matrix {
    edition = ["oss", "ent"]
  }

  terraform_cli = terraform_cli.default
  terraform     = terraform.k8s

  providers = [
    provider.enos.k8s,
    provider.helm.default,
  ]

  locals {
    image_path = abspath(var.vault_docker_image_archive)

    image_repo = var.vault_image_repository != null ? var.vault_image_repository : matrix.edition == "oss" ? "hashicorp/vault" : "hashicorp/vault-enterprise"
    image_tag = matrix.edition == "oss" ? var.vault_product_version : "${var.vault_product_version}-ent"

    version_includes_build_date = semverconstraint(var.vault_product_version, ">=1.11.0")
  }

  step "read_license" {
    skip_step = matrix.edition == "oss"
    module    = module.read_license

    variables {
      file_name = abspath(joinpath(path.root, "../support/vault.hclic"))
    }
  }

  step "create_kind_cluster" {
    module = module.create_kind_cluster

    providers = {
      enos = provider.enos.k8s
    }

    variables {
      kubeconfig_path = abspath(joinpath(path.root, "kubeconfig"))
    }
  }

  step "load_docker_image" {
    module = module.load_docker_image

    providers = {
      enos = provider.enos.k8s
    }

    variables {
      cluster_name = step.create_kind_cluster.cluster_name
      image        = local.image_repo
      tag          = local.image_tag
      archive      = var.vault_docker_image_archive
    }

    depends_on = [step.create_kind_cluster]
  }

  step "deploy_vault" {
    module = module.k8s_deploy_vault

    providers = {
      enos = provider.enos.k8s
    }

    variables {
      image_tag         = step.load_docker_image.tag
      context_name      = step.create_kind_cluster.context_name
      image_repository  = step.load_docker_image.repository
      kubeconfig_base64 = step.create_kind_cluster.kubeconfig_base64
      vault_edition     = matrix.edition
      ent_license       = matrix.edition != "oss" ? step.read_license.license : null
    }

    depends_on = [step.load_docker_image, step.create_kind_cluster]
  }

  step "verify_build_date" {
    module = module.k8s_verify_build_date

    variables {
      vault_pods        = step.deploy_vault.vault_pods
      vault_root_token  = step.deploy_vault.vault_root_token
      kubeconfig_base64 = step.create_kind_cluster.kubeconfig_base64
      context_name      = step.create_kind_cluster.context_name
    }

    depends_on = [step.deploy_vault]
  }

  step "verify_replication" {
    module = module.k8s_verify_replication

    variables {
      vault_pods        = step.deploy_vault.vault_pods
      vault_edition     = matrix.edition
      kubeconfig_base64 = step.create_kind_cluster.kubeconfig_base64
      context_name      = step.create_kind_cluster.context_name
    }

    depends_on = [step.deploy_vault]
  }

  step "verify_ui" {
    module = module.k8s_verify_ui
    skip_step = matrix.edition == "oss"

    variables {
      vault_pods        = step.deploy_vault.vault_pods
      kubeconfig_base64 = step.create_kind_cluster.kubeconfig_base64
      context_name      = step.create_kind_cluster.context_name
    }

    depends_on = [step.deploy_vault]
  }

  step "verify_version" {
    module = module.k8s_verify_version

    variables {
      vault_pods        = step.deploy_vault.vault_pods
      vault_root_token  = step.deploy_vault.vault_root_token
      vault_edition     = matrix.edition
      kubeconfig_base64 = step.create_kind_cluster.kubeconfig_base64
      context_name      = step.create_kind_cluster.context_name
      check_build_date  = local.version_includes_build_date
    }

    depends_on = [step.deploy_vault]
  }

  step "verify_write_data" {
    module = module.k8s_verify_write_data

    variables {
      vault_pods        = step.deploy_vault.vault_pods
      vault_root_token  = step.deploy_vault.vault_root_token
      kubeconfig_base64 = step.create_kind_cluster.kubeconfig_base64
      context_name      = step.create_kind_cluster.context_name
    }

    depends_on = [step.deploy_vault]
  }
}
