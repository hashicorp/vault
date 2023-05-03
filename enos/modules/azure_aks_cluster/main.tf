terraform {
  required_providers {
    azurerm = {
      source = "hashicorp/azurerm"
    }
    random = {
      source = "hashicorp/random"
    }
  }
}



#resource "azurerm_resource_provider_registration" "enable_workload_identity_feature" {
#  name = "Microsoft.ContainerService"
#
#  feature {
#    name       = "EnableWorkloadIdentityPreview"
#    registered = true
#  }
#}

resource "random_pet" "cluster_name" {}

resource "azurerm_resource_group" "this" {
  name     = "${random_pet.cluster_name.id}-rg"
  location = "East US"
}

resource "azurerm_kubernetes_cluster" "this" {
  name                      = random_pet.cluster_name.id
  location                  = azurerm_resource_group.this.location
  resource_group_name       = azurerm_resource_group.this.name
  dns_prefix                = "exampleaks1"
  oidc_issuer_enabled       = true
  workload_identity_enabled = true

  default_node_pool {
    name       = "default"
    vm_size    = "Standard_D2_v2"
    node_count = var.node_count
  }

  identity {
    type = "SystemAssigned"
  }

  tags = {
    Environment = "Test"
  }
  #  depends_on = [
  #    azurerm_resource_provider_registration.enable_workload_identity_feature
  #  ]
}

resource "local_file" "kubeconfig_file" {
  filename = var.kubeconfig_path
  content  = azurerm_kubernetes_cluster.this.kube_config_raw
}

output "kube_config" {
  value = azurerm_kubernetes_cluster.this.kube_config.0
}

output "kubeconfig_base64" {
  value = base64encode(azurerm_kubernetes_cluster.this.kube_config_raw)
}

output "cluster_name" {
  value = azurerm_kubernetes_cluster.this.name
}

output "oidc_issuer_url" {
  value = azurerm_kubernetes_cluster.this.oidc_issuer_url
}

output "resource_group_name" {
  value = azurerm_resource_group.this.name
}

output "resource_group_location" {
  value = azurerm_resource_group.this.location
}
