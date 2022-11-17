terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }

  cloud {
    hostname = "app.terraform.io"
  }
}

provider "aws" {
  region = "us-east-1"
  alias  = "us_east_1"
}

provider "aws" {
  region = "us-east-2"
  alias  = "us_east_2"
}

provider "aws" {
  region = "us-west-1"
  alias  = "us_west_1"
}

provider "aws" {
  region = "us-west-2"
  alias  = "us_west_2"
}

module "create_enos_ssh_key_pair_us_east_1" {
  source = "./modules/create_enos_ci_ssh_key"

  providers = {
    aws = aws.us_east_1
  }

  public_key = var.aws_ssh_public_key
}

module "create_enos_ssh_key_pair_us_east_2" {
  source = "./modules/create_enos_ci_ssh_key"

  providers = {
    aws = aws.us_east_2
  }

  public_key = var.aws_ssh_public_key
}

module "create_enos_ssh_key_pair_us_west_1" {
  source = "./modules/create_enos_ci_ssh_key"

  providers = {
    aws = aws.us_west_1
  }

  public_key = var.aws_ssh_public_key
}

module "create_enos_ssh_key_pair_us_west_2" {
  source = "./modules/create_enos_ci_ssh_key"

  providers = {
    aws = aws.us_west_2
  }

  public_key = var.aws_ssh_public_key
}
