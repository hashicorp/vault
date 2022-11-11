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

locals {
  workspace_names = toset([for region in var.regions : "enos-ci-bootstrap-${region}"])
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
