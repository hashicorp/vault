terraform {
  required_providers {
    tfe = {
      source = "hashicorp/tfe"
    }
  }
}

variable "tfc_api_token" {
  description = "The Terraform Cloud QTI Organization API token."
  type        = string
  sensitive   = true
}

variable "regions" {
  description = "The regions to create bootstrap workspaces for"
  type        = list(string)
}

variable "organization" {
  description = "The organization where the workspaces should be created"
  type        = string
}

variable "product_line" {
  description = "The product line to boostrap enos ci for"
  type        = string
  default     = "vault"
}

locals {
  workspace_names = toset([for region in var.regions : "${var.product_line}-ci-enos-boostrap-${region}"])
}

resource "tfe_workspace" "ci_bootstrap_workspaces" {
  for_each = local.workspace_names

  name           = each.key
  organization   = var.organization
  execution_mode = "local"
}

output "workspace_names" {
  value = [for workspace in tfe_workspace.ci_bootstrap_workspaces : workspace.name]
}
