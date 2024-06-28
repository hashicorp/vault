# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

provider "aws" "default" {
  region = var.aws_region
}

# This default SSH user is used in RHEL, Amazon Linux, SUSE, and Leap distros
provider "enos" "ec2_user" {
  transport = {
    ssh = {
      user             = "ec2-user"
      private_key_path = abspath(var.aws_ssh_private_key_path)
    }
  }
}

# This default SSH user is used in the Ubuntu distro
provider "enos" "ubuntu" {
  transport = {
    ssh = {
      user             = "ubuntu"
      private_key_path = abspath(var.aws_ssh_private_key_path)
    }
  }
}
