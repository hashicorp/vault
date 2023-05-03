terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
    azurerm = {
      source = "hashicorp/azurerm"
    }
  }
}

data "azurerm_subscription" "current" {
}

resource "enos_remote_exec" "setup_auth" {
  environment = {
    VAULT_TOKEN           = var.vault_root_token
    VAULT_ADDR            = "http://127.0.0.1:8200"
    TENANT_ID             = data.azurerm_subscription.current.tenant_id
    CLIENT_ID             = var.client_id
    CLIENT_SECRET         = var.client_secret
    BOUND_SUBSCRIPTION_ID = data.azurerm_subscription.current.subscription_id
    BOUND_RESOURCE_GROUP  = data.azurerm_subscription.current.display_name
  }

  scripts = [abspath("${path.module}/scripts/setup_auth.sh")]

  transport = var.transport
}
