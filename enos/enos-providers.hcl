// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

provider "aws" "default" {
  region = var.aws_region
}

provider "docker" "default" {
}

// This default SSH user is used in RHEL, Amazon Linux, SUSE, and Leap distros
provider "enos" "ec2_user" {
  transport = {
    ssh = {
      user             = "ec2-user"
      private_key_path = abspath(var.aws_ssh_private_key_path)
    }
  }
}

// This default SSH user is used in the Ubuntu distro
provider "enos" "ubuntu" {
  transport = {
    ssh = {
      user             = "ubuntu"
      private_key_path = abspath(var.aws_ssh_private_key_path)
    }
  }
}

provider "enos" "fyre_root" {
  transport = {
    ssh = {
      user             = "root"
      private_key_path = abspath(var.fyre_private_key_path)
    }
  }
}

provider "fyre" "rtp" {
  product_group_id = var.fyre_product_group_id
  site             = "rtp"
}

provider "fyre" "svl" {
  product_group_id = var.fyre_product_group_id
  site             = "svl"
}

provider "hcp" "default" {
}

provider "local" "default" {
}

provider "time" "default" {
}
