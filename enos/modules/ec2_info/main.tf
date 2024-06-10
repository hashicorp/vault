# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

# Note: in order to use the openSUSE Leap AMIs, the AWS account in use must "subscribe"
# and accept SUSE's terms of use. You can do this at the links below. If the AWS account
# you are using is already subscribed, this confirmation will be displayed on each page.
# openSUSE Leap arm64 subscription: https://aws.amazon.com/marketplace/server/procurement?productId=a516e959-df54-4035-bb1a-63599b7a6df9
# openSUSE leap amd64 subscription: https://aws.amazon.com/marketplace/server/procurement?productId=5535c495-72d4-4355-b169-54ffa874f849

locals {
  architectures      = toset(["arm64", "x86_64"])
  amzn2_owner_id     = "591542846629"
  canonical_owner_id = "099720109477"
  sles_owner_id      = "013907871322"
  suse_owner_id      = "679593333241"
  rhel_owner_id      = "309956199498"
  ids = {
    "arm64" = {
      "amzn2" = {
        "2" = data.aws_ami.amzn2["arm64"].id
      }
      "rhel" = {
        "8.9" = data.aws_ami.rhel_89["arm64"].id
        "9.3" = data.aws_ami.rhel_93["arm64"].id
      }
      "sles" = {
        "v15_sp5_standard" = data.aws_ami.sles_15_sp5_standard["arm64"].id
      }
      "ubuntu" = {
        "20.04" = data.aws_ami.ubuntu_2004["arm64"].id
        "22.04" = data.aws_ami.ubuntu_2204["arm64"].id
      }
    }
    "amd64" = {
      "amzn2" = {
        "2" = data.aws_ami.amzn2["x86_64"].id
      }
      "leap" = {
        "15.5" = data.aws_ami.leap_155.id
      }
      "rhel" = {
        "8.9" = data.aws_ami.rhel_89["x86_64"].id
        "9.3" = data.aws_ami.rhel_93["x86_64"].id
      }
      "sles" = {
        "v15_sp5_standard" = data.aws_ami.sles_15_sp5_standard["x86_64"].id
      }
      "ubuntu" = {
        "20.04" = data.aws_ami.ubuntu_2004["x86_64"].id
        "22.04" = data.aws_ami.ubuntu_2204["x86_64"].id
      }
    }
  }
}

data "aws_ami" "ubuntu_2004" {
  most_recent = true
  for_each    = local.architectures

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-*-20.04-*-server-*"]
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

data "aws_ami" "rhel_89" {
  most_recent = true
  for_each    = local.architectures

  # Currently latest latest point release-1
  filter {
    name   = "name"
    values = ["RHEL-8.9*HVM-20*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name   = "architecture"
    values = [each.value]
  }

  owners = [local.rhel_owner_id]
}

data "aws_ami" "rhel_93" {
  most_recent = true
  for_each    = local.architectures

  # Currently latest latest point release-1
  filter {
    name   = "name"
    values = ["RHEL-9.3*HVM-20*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name   = "architecture"
    values = [each.value]
  }

  owners = [local.rhel_owner_id]
}

data "aws_ami" "amzn2" {
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

  owners = [local.amzn2_owner_id]
}

data "aws_ami" "sles_15_sp5_standard" {
  most_recent = true
  for_each    = local.architectures

  filter {
    name   = "name"
    values = ["suse-sles-15-sp5-v*-hvm-*"]
  }

  filter {
    name   = "architecture"
    values = [each.value]
  }

  owners = [local.sles_owner_id]
}

data "aws_ami" "leap_155" {
  most_recent = true

  filter {
    name   = "name"
    values = ["openSUSE-Leap-15.5*"]
  }

  filter {
    name = "architecture"
    # Note: arm64 AMIs are not offered for Leap.
    values = ["x86_64"]
  }

  owners = [local.suse_owner_id]
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
