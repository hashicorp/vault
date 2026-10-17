# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1
#
# CI entry point: .github/workflows/test-run-enos-scenario-containers.yml
#
# That workflow injects the following variables via ENOS_VAR_* env vars
# (no -var flags needed locally when running in CI):
#
#   ENOS_VAR_terraform_plugin_cache_dir  → var.terraform_plugin_cache_dir
#   ENOS_VAR_vault_build_date            → var.vault_build_date
#   ENOS_VAR_vault_version               → var.vault_version
#   ENOS_VAR_vault_revision              → var.vault_revision
#   ENOS_VAR_container_image_archive     → var.container_image_archive
#
# When running locally you must supply at minimum:
#   -var vault_version=<version>
#   -var container_image_archive=<path/to/image.tar>  (or omit to pull from registry)
# All other variables have sensible defaults defined in enos-variables-k8s.hcl.

scenario "k8s" {
  description = <<-EOF
    The k8s scenario verifies Vault when running in Kubernetes mode. The build can be a container
    in a remote repository or a local container archive tarball.

    The scenario creates a new kind kubernetes cluster in Docker and creates a Vault Cluster using
    the candidate artifact and verifies behavior against the Vault cluster.
  EOF

  matrix {
    edition = ["ce", "ent", "ent.fips1403", "ent.hsm", "ent.hsm.fips1403"]
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
      "ent.fips1403" = {
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
      "ent.hsm.fips1403" = {
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

  step "load_client_docker_image" {
    description = <<-EOF
      Load the Vault Client container image into the kind k8s cluster for testing as a sidecar.

      TODAY  — var.client_container_image_archive is null, so this falls back to
      var.container_image_archive and loads the same vault-enterprise image that the
      server uses. The sidecar runs "vault agent" from that image.

      FUTURE — set var.client_container_image_archive to the standalone client image
      archive (.tar). That image will be loaded here and used as the sidecar container
      in deploy_client_sidecar. Also set var.kubeclient_bin_path to the path of the
      kubeclient binary inside that image (e.g. "/bin/kubeclient") to activate the
      verify_client_smoke tests.

      Both variables can be set via ENOS_VAR_* in CI or with -var flags locally:
        ENOS_VAR_client_container_image_archive=/path/to/client.tar
        ENOS_VAR_kubeclient_bin_path=/bin/kubeclient
    EOF
    module      = module.load_docker_image
    depends_on  = [step.create_kind_cluster]

    variables {
      cluster_name = step.create_kind_cluster.cluster_name
      image        = local.repo_metadata[matrix.edition][matrix.repo].repo
      tag          = local.repo_metadata[matrix.edition][matrix.repo].tag
      // TODAY  : falls back to the vault-enterprise archive (client_container_image_archive is null)
      // FUTURE : set var.client_container_image_archive to the standalone client image archive
      archive = var.client_container_image_archive != null ? var.client_container_image_archive : var.container_image_archive
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

  step "deploy_client_sidecar" {
    description = <<-EOF
      Deploy a test application pod alongside the Vault Client container running as a sidecar.
    EOF
    module      = module.k8s_deploy_client_sidecar
    depends_on = [
      step.deploy_vault,
      step.load_client_docker_image,
      step.create_kind_cluster,
    ]

    variables {
      client_image_tag        = step.load_client_docker_image.tag
      client_image_repository = step.load_client_docker_image.repository
      context_name            = step.create_kind_cluster.context_name
      kubeconfig_base64       = step.create_kind_cluster.kubeconfig_base64
      vault_root_token        = step.deploy_vault.vault_root_token
      vault_pods              = step.deploy_vault.vault_pods
      kind_cluster_name       = step.create_kind_cluster.cluster_name
    }
  }

  step "verify_client_sidecar" {
    description = <<-EOF
      Verify that the Vault Client sidecar authenticated, populated secrets to the shared volume,
      and proxies requests over localhost to the Vault cluster.
    EOF
    module      = module.k8s_vault_verify_client_sidecar
    depends_on  = [step.deploy_client_sidecar]

    variables {
      context_name      = step.create_kind_cluster.context_name
      kubeconfig_base64 = step.create_kind_cluster.kubeconfig_base64
      app_pod_name      = step.deploy_client_sidecar.app_pod_name
      app_pod_namespace = step.deploy_client_sidecar.app_pod_namespace
    }
  }

  step "verify_client_smoke" {
    description = <<-EOF
      Run kubeclient binary smoke tests inside the vault-client-sidecar container.

      TODAY  — var.kubeclient_bin_path defaults to "/bin/vault". The smoke
      module detects the fallback and skips all binary-level tests, running
      only the RBAC setup so the infrastructure is ready. No failures.

      FUTURE — when the standalone client image ships, set kubeclient_bin_path
      to the path of the kubeclient binary inside that image (e.g. "/bin/kubeclient").
      The module will exec into the vault-client-sidecar container and run:
        1. get-pod        — happy path: expects pod JSON in stdout
        2. patch-pod      — single label add, verified via re-fetch
        3. patch-pod-multi— comma-separated -patches, both labels verified
        4. invalid-call   — -call=delete-pod must exit non-zero

      The vault-client-sidecar container already runs inside the cluster, so
      KUBERNETES_SERVICE_HOST and KUBERNETES_SERVICE_PORT are set automatically
      by Kubernetes — exactly what the kubeclient binary requires.
    EOF
    module      = module.k8s_vault_verify_client_smoke
    depends_on  = [step.verify_client_sidecar]

    variables {
      context_name      = step.create_kind_cluster.context_name
      kubeconfig_base64 = step.create_kind_cluster.kubeconfig_base64
      app_pod_name      = step.deploy_client_sidecar.app_pod_name
      app_pod_namespace = step.deploy_client_sidecar.app_pod_namespace
      kind_cluster_name = step.create_kind_cluster.cluster_name
      // Change this single line when the standalone client image is ready:
      kubeclient_bin_path = var.kubeclient_bin_path
    }
  }
}
