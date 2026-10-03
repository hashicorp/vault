# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

locals {
  architectures      = toset(["arm64", "x86_64"])
  amazon_owner_id    = "591542846629"
  canonical_owner_id = "099720109477"
  suse_owner_id      = "013907871322"
  redhat_owner_id    = "309956199498"
  ids = {
    // NOTE: The versions here always correspond to the output of enos_host_info.distro_version. These are used in
    // several modules so if you change the keys here also consider the "artifact/metadata", "ec2_info",
    "arm64" = {
      "amzn" = {
        "2"    = data.aws_ami.amzn_2["arm64"].id
        "2023" = data.aws_ami.amzn_2023["arm64"].id
      }
      "rhel" = {
        "8.10" = data.aws_ami.rhel_8["arm64"].id
        "9.7"  = data.aws_ami.rhel_9["arm64"].id
        "10.1" = data.aws_ami.rhel_10["arm64"].id
      }
      "sles" = {
        "15.7" = data.aws_ami.sles_15["arm64"].id
        "16.0" = data.aws_ami.sles_16["arm64"].id
      }
      "ubuntu" = {
        "22.04" = data.aws_ami.ubuntu_2204["arm64"].id
        "24.04" = data.aws_ami.ubuntu_2404["arm64"].id
      }
    }
    "amd64" = {
      "amzn" = {
        "2"    = data.aws_ami.amzn_2["x86_64"].id
        "2023" = data.aws_ami.amzn_2023["x86_64"].id
      }
      "rhel" = {
        "8.10" = data.aws_ami.rhel_8["x86_64"].id
        "9.7"  = data.aws_ami.rhel_9["x86_64"].id
        "10.1" = data.aws_ami.rhel_10["x86_64"].id
      }
      "sles" = {
        "15.7" = data.aws_ami.sles_15["x86_64"].id
        "16.0" = data.aws_ami.sles_16["x86_64"].id
      }
      "ubuntu" = {
        "22.04" = data.aws_ami.ubuntu_2204["x86_64"].id
        "24.04" = data.aws_ami.ubuntu_2404["x86_64"].id
      }
    }
  }
}

data "aws_ami" "amzn_2" {
  most_recent = true
  for_each    = local.architectures

  filter {
    name   = "name"
    values = ["amzn2-ami-ecs-hvm-2.0*"]
  }

  filter {
    name   = "architecture"
    values = [each.value]
  }

  owners = [local.amazon_owner_id]
}

data "aws_ami" "amzn_2023" {
  most_recent = true
  for_each    = local.architectures

  filter {
    name   = "name"
    values = ["al2023-ami-ecs-hvm*"]
  }

  filter {
    name   = "architecture"
    values = [each.value]
  }

  owners = [local.amazon_owner_id]
}

data "aws_ami" "rhel_8" {
  most_recent = true
  for_each    = local.architectures

  filter {
    name   = "name"
    values = ["RHEL-8.10*HVM_GA-20*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name   = "architecture"
    values = [each.value]
  }

  owners = [local.redhat_owner_id]
}

data "aws_ami" "rhel_9" {
  most_recent = true
  for_each    = local.architectures

  filter {
    name   = "name"
    values = ["RHEL-9.7*HVM_GA-20*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name   = "architecture"
    values = [each.value]
  }

  owners = [local.redhat_owner_id]
}

data "aws_ami" "rhel_10" {
  most_recent = true
  for_each    = local.architectures

  filter {
    name   = "name"
    values = ["RHEL-10.1*HVM_GA-20*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name   = "architecture"
    values = [each.value]
  }

  owners = [local.redhat_owner_id]
}

data "aws_ami" "sles_15" {
  most_recent = true
  for_each    = local.architectures

  filter {
    name   = "name"
    values = ["suse-sles-15-sp7-v*-hvm-*"]
  }

  filter {
    name   = "architecture"
    values = [each.value]
  }

  owners = [local.suse_owner_id]
}

data "aws_ami" "sles_16" {
  most_recent = true
  for_each    = local.architectures

  filter {
    name   = "name"
    values = ["suse-sles-16-0-v*-hvm-ssd-*"]
  }

  filter {
    name   = "architecture"
    values = [each.value]
  }

  owners = [local.suse_owner_id]
}

data "aws_ami" "ubuntu_2204" {
  most_recent = true
  for_each    = local.architectures

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-*-22.04-*-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name   = "architecture"
    values = [each.value]
  }

  owners = [local.canonical_owner_id]
}

data "aws_ami" "ubuntu_2404" {
  most_recent = true
  for_each    = local.architectures

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd-gp3/ubuntu-noble-24.04-*-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name   = "architecture"
    values = [each.value]
  }

  owners = [local.canonical_owner_id]
}

data "aws_region" "current" {}

data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "zone-name"
    values = ["*"]
  }
}

output "ami_ids" {
  value = local.ids
}

output "current_region" {
  value = data.aws_region.current
}

output "availability_zones" {
  value = data.aws_availability_zones.available
}
