scenario "auth" {
  matrix {
    edition      = ["oss", "ent"]
    cluster_type = ["kind", "aks"]
    auth_method  = ["azure"]

    exclude {
      cluster_type = ["kind"]
      auth_method  = ["azure"]
    }
  }

  terraform_cli = terraform_cli.default
  terraform     = terraform.azure

  providers = local.providers[matrix.cluster_type]

  locals {
    vault_image_repo = var.vault_image_repository != null ? var.vault_image_repository : matrix.edition == "oss" ? "hashicorp/vault" : "hashicorp/vault-enterprise"
    vault_image_tag  = replace(var.vault_product_version, "+ent", "-ent")

    providers = {
      kind = [
        provider.enos.default,
        provider.helm.default,
      ]
      aks = [
        provider.enos.default,
        provider.helm.default,
        provider.azurerm.default,
        provider.azuread.default,
        provider.kubernetes.default,
      ]
    }
  }

  step "create_cluster" {
    module = "create_${matrix.cluster_type}_cluster"
  }

  step "setup_auth_backend" {
    module = "setup_${matrix.auth_method}_auth_backend"

    variables {
      oidc_issuer_url  = step.create_cluster.oidc_issuer_url
      cluster_id       = step.create_cluster.cluster_name
      service_accounts = ["vault", "vault-agent"]
    }
  }

  step "read_license" {
    skip_step = matrix.edition == "oss"
    module    = module.read_license

    variables {
      file_name = try(abspath(var.vault_license_path), "")
    }
  }

  step "extra_chart_values" {
    module    = module.k8s_azure_wif_extra_chart_values
    skip_step = matrix.auth_method != "azure"

    variables {
      azure_ad_application_id = step.setup_auth_backend.application_id
    }
  }

  step "deploy_vault" {
    module = module.k8s_deploy_vault

    variables {
      image_repository     = local.vault_image_repo
      image_tag            = local.vault_image_tag
      image_pull_policy    = "Always"
      context_name         = step.create_cluster.cluster_name
      kubeconfig_base64    = step.create_cluster.kubeconfig_base64
      vault_edition        = matrix.edition
      vault_instance_count = 3
      vault_log_level      = "debug"
      ent_license          = try(step.read_license.license, "")

      node_count = 3

      extra_helm_release_values = matrix.auth_method == "azure" ? step.extra_chart_values.chart_values : null
    }

    depends_on = [
      step.create_cluster,
      step.extra_chart_values,
      step.setup_auth_backend,
    ]
  }

  step "format_transport" {
    module = module.format_transport_block

    variables {
      type          = "kubernetes"
      configuration = step.deploy_vault.transports.0
    }
  }

  step "setup_auth_method" {
    module = "setup_${matrix.auth_method}_auth_method"

    variables {
      vault_root_token = step.deploy_vault.vault_root_token
      client_id        = step.setup_auth_backend.application_id
      client_secret    = step.setup_auth_backend.client_secret
      transport        = step.format_transport.transport_config
    }

    depends_on = [
      step.deploy_vault
    ]
  }

  step "deploy_vault_agent" {
    module = module.deploy_vault_agent

    variables {
      client_id            = step.setup_auth_backend.application_id
      service_account_name = "vault-agent"
      docker_image_name    = var.vault_agent_image_name
      docker_image_tag     = var.vault_agent_image_tag
      kubernetes_context   = step.create_cluster.cluster_name
      kubeconfig_base64    = step.create_cluster.kubeconfig_base64
    }

    depends_on = [
      step.create_cluster,
      step.setup_auth_method,
    ]
  }

  step "verify_vault_agent_output" {
    module = module.vault_verify_agent_output

    variables {
      transport                        = step.deploy_vault_agent.transport
      vault_agent_template_destination = "/tmp/agent/render-content.txt"
      vault_agent_expected_output      = "orphan=true display_name=dev-role"
    }

    depends_on = [
      step.deploy_vault_agent,
    ]
  }

  output "cluster_name" {
    value = step.create_cluster.cluster_name
  }
}
