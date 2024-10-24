# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

scenario "k8s" {
  description = <<-EOF
    The k8s scenario verifies Vault when running in Kubernetes mode. The build can be a container
    in a remote repository or a local container archive tarball.

    The scenario creates a new kind kubernetes cluster in Docker and creates a Vault Cluster using
    the candidate artifact and verifies behavior against the Vault cluster.
  EOF

  matrix {
    edition = ["ce", "ent", "ent.fips1402", "ent.hsm", "ent.hsm.fips1402"]
    repo    = ["docker", "ecr", "quay"]
  }

  terraform_cli = terraform_cli.default
  terraform     = terraform.k8s

  providers = [
    provider.enos.default,
    provider.helm.default,
  ]

  locals {
    // For now this works as the vault_version includes metadata. If we ever get to the point that
    // vault_version excludes metadata we'll have to include the matrix.edition here as well.
    tag_version     = replace(var.vault_version, "+ent", "-ent")
    tag_version_ubi = "${local.tag_version}-ubi"
    // When we load candidate images into our k8s cluster we verify that the archives embedded
    // repository and tag match our expectations. This is the source of truth for what we _expect_
    // various artifacts to have. The source of truth for what we use when building is defined in
    // .github/actions/containerize. If you are modifying these expectations you likely need to
    // modify the source of truth there.
    repo_metadata = {
      "ce" = {
        docker = {
          // https://hub.docker.com/r/hashicorp/vault
          repo = "hashicorp/vault"
          tag  = local.tag_version
        }
        ecr = {
          // https://gallery.ecr.aws/hashicorp/vault
          repo = "public.ecr.aws/hashicorp/vault"
          tag  = local.tag_version
        }
        quay = {
          // https://catalog.redhat.com/software/containers/hashicorp/vault/5fda55bd2937386820429e0c
          repo = "quay.io/redhat-isv-containers/5f89bb5e0b94cf64cfeb500a"
          tag  = local.tag_version_ubi
        }
      },
      "ent" = {
        docker = {
          // https://hub.docker.com/r/hashicorp/vault-enterprise
          repo = "hashicorp/vault-enterprise"
          tag  = local.tag_version
        }
        ecr = {
          // https://gallery.ecr.aws/hashicorp/vault-enterprise
          repo = "public.ecr.aws/hashicorp/vault-enterprise"
          tag  = local.tag_version
        }
        quay = {
          // https://catalog.redhat.com/software/containers/hashicorp/vault-enterprise/5fda5633ac3db90370a26443
          repo = "quay.io/redhat-isv-containers/5f89bb9242e382c85087dce2"
          tag  = local.tag_version_ubi
        }
      },
      "ent.fips1402" = {
        docker = {
          // https://hub.docker.com/r/hashicorp/vault-enterprise-fips
          repo = "hashicorp/vault-enterprise-fips"
          tag  = local.tag_version
        }
        ecr = {
          // https://gallery.ecr.aws/hashicorp/vault-enterprise-fips
          repo = "public.ecr.aws/hashicorp/vault-enterprise-fips"
          tag  = local.tag_version
        }
        quay = {
          // https://catalog.redhat.com/software/containers/hashicorp/vault-enterprise-fips/628d50e37ff70c66a88517ea
          repo = "quay.io/redhat-isv-containers/6283f645d02c6b16d9caeb8e"
          tag  = local.tag_version_ubi
        }
      },
      "ent.hsm" = {
        docker = {
          // https://hub.docker.com/r/hashicorp/vault-enterprise
          repo = "hashicorp/vault-enterprise"
          tag  = local.tag_version
        }
        ecr = {
          // https://gallery.ecr.aws/hashicorp/vault-enterprise
          repo = "public.ecr.aws/hashicorp/vault-enterprise"
          tag  = local.tag_version
        }
        quay = {
          // https://catalog.redhat.com/software/containers/hashicorp/vault-enterprise/5fda5633ac3db90370a26443
          repo = "quay.io/redhat-isv-containers/5f89bb9242e382c85087dce2"
          tag  = local.tag_version_ubi
        }
      },
      "ent.hsm.fips1402" = {
        docker = {
          // https://hub.docker.com/r/hashicorp/vault-enterprise
          repo = "hashicorp/vault-enterprise"
          tag  = local.tag_version
        }
        ecr = {
          // https://gallery.ecr.aws/hashicorp/vault-enterprise
          repo = "public.ecr.aws/hashicorp/vault-enterprise"
          tag  = local.tag_version
        }
        quay = {
          // https://catalog.redhat.com/software/containers/hashicorp/vault-enterprise/5fda5633ac3db90370a26443
          repo = "quay.io/redhat-isv-containers/5f89bb9242e382c85087dce2"
          tag  = local.tag_version_ubi
        }
      },
    }
    // The additional '-0' is required in the constraint since without it, the semver function will
    // only compare the non-pre-release parts (Major.Minor.Patch) of the version and the constraint,
    // which can lead to unexpected results.
    version_includes_build_date = semverconstraint(var.vault_version, ">=1.11.0-0")
  }

  step "read_license" {
    skip_step = matrix.edition == "ce"
    module    = module.read_license

    variables {
      file_name = abspath(joinpath(path.root, "../support/vault.hclic"))
    }
  }

  step "create_kind_cluster" {
    module = module.create_kind_cluster

    variables {
      kubeconfig_path = abspath(joinpath(path.root, "kubeconfig"))
    }
  }

  step "load_docker_image" {
    description = <<-EOF
      Load an verify the tags of a Vault container image into the kind k8s cluster. If no
      var.container_image_archive has been set it will attempt to load an image matching the
      var.vault_version from the matrix.repo.
    EOF
    module      = module.load_docker_image
    depends_on  = [step.create_kind_cluster]

    verifies = [
      quality.vault_artifact_container_alpine,
      quality.vault_artifact_container_ubi,
      quality.vault_artifact_container_tags,
    ]

    variables {
      cluster_name = step.create_kind_cluster.cluster_name
      image        = local.repo_metadata[matrix.edition][matrix.repo].repo
      tag          = local.repo_metadata[matrix.edition][matrix.repo].tag
      archive      = var.container_image_archive
    }
  }

  step "deploy_vault" {
    module = module.k8s_deploy_vault
    depends_on = [
      step.load_docker_image,
      step.create_kind_cluster,
    ]

    variables {
      image_tag         = step.load_docker_image.tag
      context_name      = step.create_kind_cluster.context_name
      image_repository  = step.load_docker_image.repository
      kubeconfig_base64 = step.create_kind_cluster.kubeconfig_base64
      vault_edition     = matrix.edition
      vault_log_level   = var.log_level
      ent_license       = matrix.edition != "ce" ? step.read_license.license : null
    }
  }

  step "verify_replication" {
    module     = module.k8s_verify_replication
    depends_on = [step.deploy_vault]

    variables {
      vault_pods        = step.deploy_vault.vault_pods
      vault_edition     = matrix.edition
      kubeconfig_base64 = step.create_kind_cluster.kubeconfig_base64
      context_name      = step.create_kind_cluster.context_name
    }
  }

  step "verify_version" {
    module     = module.k8s_verify_version
    depends_on = [step.deploy_vault]

    variables {
      vault_pods        = step.deploy_vault.vault_pods
      vault_root_token  = step.deploy_vault.vault_root_token
      vault_edition     = matrix.edition
      kubeconfig_base64 = step.create_kind_cluster.kubeconfig_base64
      context_name      = step.create_kind_cluster.context_name
      check_build_date  = local.version_includes_build_date
      vault_build_date  = var.vault_build_date
    }
  }

  step "verify_write_data" {
    module     = module.k8s_verify_write_data
    depends_on = [step.deploy_vault]

    variables {
      vault_pods        = step.deploy_vault.vault_pods
      vault_root_token  = step.deploy_vault.vault_root_token
      kubeconfig_base64 = step.create_kind_cluster.kubeconfig_base64
      context_name      = step.create_kind_cluster.context_name
    }
  }
}
