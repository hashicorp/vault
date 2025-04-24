module "create_aws_secrets_engine" {
  count  = var.create_aws_secrets_engine ? 1 : 0
  source = "./aws"

  hosts             = var.hosts
  leader_host       = var.leader_host
  vault_addr        = var.vault_addr
  vault_root_token  = var.vault_root_token
  vault_install_dir = var.vault_install_dir
}

locals {
  aws_output = try(module.create_aws_secrets_engine[0].output, null)
}

output "aws" {
  value = local.aws_output
}
