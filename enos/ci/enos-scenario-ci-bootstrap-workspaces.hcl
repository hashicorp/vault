module "create_ci_bootstrap_workspaces" {
  source = abspath("./modules/create_bootstrap_workspaces")

  organization  = var.organization
  regions       = ["us-east-1", "us-east-2", "us-west-1", "us-west-2"]
  tfc_api_token = var.tfc_api_token
}

scenario "bootstrap_workspaces" {

  terraform     = terraform.bootstrap
  terraform_cli = terraform_cli.default

  providers = [ provider.tfe.bootstrap ]

  step "create_ci_bootstrap_workspaces" {

    providers = {
      tfe = provider.tfe.bootstrap
    }
    module = module.create_ci_bootstrap_workspaces
  }

  output "workspace_names" {
    value = step.create_ci_bootstrap_workspaces.workspace_names
  }
}
