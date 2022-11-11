module "create_enos_ci_ssh_key" {
  source = abspath("./modules/create_enos_ci_ssh_key")

  public_key = file(var.aws_ssh_public_key_path)
}

scenario "bootstrap_ci" {
  terraform     = terraform[local.region]
  terraform_cli = terraform_cli.default

  providers = [
    provider.aws[local.region],
  ]

  matrix {
    region = ["us-east-1", "us-east-2", "us-west-1", "us-west-2"]
  }

  locals {
    region = replace(matrix.region, "-", "_")
  }

  step "create_enos_ssh_key_pair" {
    module = module.create_enos_ci_ssh_key

    providers = {
      aws = provider.aws[local.region]
    }

    variables {
      region = matrix.region
    }
  }

  output "key_pair_id" {
    value = step.create_enos_ssh_key_pair.key_pair_id
  }

  output "key_pair_arn" {
    value = step.create_enos_ssh_key_pair.key_pair_arn
  }
}
