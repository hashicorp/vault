terraform {
  required_providers {
    azuread = {
      source = "hashicorp/azuread"
    }
    azurerm = {
      source = "hashicorp/azurerm"
    }
  }
}

locals {
  app_name           = "${var.cluster_id}-${var.app_name}"
  app_rw_owned_by_id = data.azuread_service_principal.ms_graph.app_role_ids["Application.ReadWrite.All"]
}

data "azuread_application_published_app_ids" "well_known" {}

data "azuread_client_config" "current" {}
data "azurerm_subscription" "current" {}

data "azuread_service_principal" "ms_graph" {
  application_id = data.azuread_application_published_app_ids.well_known.result.MicrosoftGraph
}

# An Azure AD application which will be linked with Vault in order to provide the identities and
# service (Azure services) permissions that Vault requires to support the Azure auth method
resource "azuread_application" "vault_app" {
  display_name = local.app_name
  owners       = [data.azuread_client_config.current.object_id]

  # The Vault server requires Application.ReadWrite.All access to the Microsoft Graph service
  # The resource access requirements can be seen here: https://developer.hashicorp.com/vault/docs/auth/azure
  required_resource_access {
    resource_app_id = data.azuread_application_published_app_ids.well_known.result.MicrosoftGraph

    resource_access {
      id   = local.app_rw_owned_by_id
      type = "Role" # Application type
    }
  }
}

# The service principal linked to the AD application created above. Any roles and permissions
# that Vault requires for the Azure auth method should be assigned to the service principal.
resource "azuread_service_principal" "vault_app" {
  application_id = azuread_application.vault_app.application_id
}

# The Vault server currently uses a shared secret in order to access the VMSS api to verify the jwt
# token used during authentication using the Azure auth method. This should not be required when fully
# supported WIF since, the pod should rather use the inject federated identity token. Full documentation
# regarding how to use WIF with AKS can be found here: https://learn.microsoft.com/en-us/azure/aks/workload-identity-overview
resource "azuread_service_principal_password" "password" {
  service_principal_id = azuread_service_principal.vault_app.id
}

# Each role service account will be linked to one federated identity
resource "azuread_application_federated_identity_credential" "vault" {
  for_each = var.service_accounts

  application_object_id = azuread_application.vault_app.object_id
  display_name          = "${var.cluster_id}-${each.value}"
  audiences             = ["api://AzureADTokenExchange"]
  issuer                = var.oidc_issuer_url
  subject               = "system:serviceaccount:default:${each.value}"
}

# The Vault server currently requires VMMS (VirtualMachineScaleSet) read permissions since, in order
# to validate the jwt used token used during login it must assert identity of the actor logging in.
# Currently Pod identity is used which associates identities with the nodes (VM) of an AKS cluster.
# This should not be required once WIF is properly supported in the Azure auth plugin.
resource "azurerm_role_definition" "vmss_read" {
  name  = "${var.cluster_id}-vmms-read"
  scope = data.azurerm_subscription.current.id

  permissions {
    actions     = ["Microsoft.Compute/virtualMachineScaleSets/read"]
    not_actions = []
  }

  assignable_scopes = [
    data.azurerm_subscription.current.id,
  ]
}

resource "azurerm_role_assignment" "vmss_read" {
  scope              = data.azurerm_subscription.current.id
  role_definition_id = azurerm_role_definition.vmss_read.role_definition_resource_id
  principal_id       = azuread_service_principal.vault_app.id
}
