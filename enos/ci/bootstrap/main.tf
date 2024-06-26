# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }

  cloud {
    hostname     = "app.terraform.io"
    organization = "hashicorp-qti"
    // workspace must be exported in the environment as: TF_WORKSPACE=<vault|vault-enterprise>-ci-enos-boostrap
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


locals {
  key_name = "${var.repository}-ci-ssh-key"
}

resource "aws_key_pair" "enos_ci_key_us_east_1" {
  key_name   = local.key_name
  public_key = var.aws_ssh_public_key

  provider = aws.us_east_1
}

resource "aws_key_pair" "enos_ci_key_us_east_2" {
  key_name   = local.key_name
  public_key = var.aws_ssh_public_key

  provider = aws.us_east_2
}

resource "aws_key_pair" "enos_ci_key_us_west_1" {
  key_name   = local.key_name
  public_key = var.aws_ssh_public_key

  provider = aws.us_west_1
}

resource "aws_key_pair" "enos_ci_key_us_west_2" {
  key_name   = local.key_name
  public_key = var.aws_ssh_public_key

  provider = aws.us_west_2
}
