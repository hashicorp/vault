# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

# Note: in order to use the openSUSE Leap AMIs, the AWS account in use must "subscribe"
# and accept SUSE's terms of use. You can do this at the links below. If the AWS account
# you are using is already subscribed, this confirmation will be displayed on each page.
# openSUSE Leap arm64 subscription: https://aws.amazon.com/marketplace/server/procurement?productId=a516e959-df54-4035-bb1a-63599b7a6df9
# openSUSE Leap amd64 subscription: https://aws.amazon.com/marketplace/server/procurement?productId=5535c495-72d4-4355-b169-54ffa874f849

locals {
  architectures      = toset(["arm64", "x86_64"])
  amazon_owner_id    = "591542846629"
  canonical_owner_id = "099720109477"
  suse_owner_id      = "013907871322"
  opensuse_owner_id  = "679593333241"
  redhat_owner_id    = "309956199498"
  ids = {
    // NOTE: If you modify these versions you'll probably also need to update the `softhsm_install`
    // module to match.
    "arm64" = {
      "amzn" = {
        "2"    = data.aws_ami.amzn_2["arm64"].id
        "2023" = data.aws_ami.amzn_2023["arm64"].id
      }
      "leap" = {
        "15.6" = data.aws_ami.leap_15["arm64"].id
      }
      "rhel" = {
        "8.10" = data.aws_ami.rhel_8["arm64"].id
        "9.4"  = data.aws_ami.rhel_9["arm64"].id
      }
      "sles" = {
        "15.6" = data.aws_ami.sles_15["arm64"].id
      }
      "ubuntu" = {
        "20.04" = data.aws_ami.ubuntu_2004["arm64"].id
        "22.04" = data.aws_ami.ubuntu_2204["arm64"].id
        "24.04" = data.aws_ami.ubuntu_2404["arm64"].id
      }
    }
    "amd64" = {
      "amzn" = {
        "2"    = data.aws_ami.amzn_2["x86_64"].id
        "2023" = data.aws_ami.amzn_2023["x86_64"].id
      }
      "leap" = {
        "15.6" = data.aws_ami.leap_15["x86_64"].id
      }
      "rhel" = {
        "8.10" = data.aws_ami.rhel_8["x86_64"].id
        "9.4"  = data.aws_ami.rhel_9["x86_64"].id
      }
      "sles" = {
        "15.6" = data.aws_ami.sles_15["x86_64"].id
      }
      "ubuntu" = {
        "20.04" = data.aws_ami.ubuntu_2004["x86_64"].id
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

data "aws_ami" "leap_15" {
  most_recent = true
  for_each    = local.architectures

  filter {
    name   = "name"
    values = ["openSUSE-Leap-15-6*"]
  }

  filter {
    name   = "architecture"
    values = [each.value]
  }

  owners = [local.opensuse_owner_id]
}

data "aws_ami" "rhel_8" {
  most_recent = true
  for_each    = local.architectures

  # Currently latest latest point release-1
  filter {
    name   = "name"
    values = ["RHEL-8.10*HVM-20*"]
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

  # Currently latest latest point release-1
  filter {
    name   = "name"
    values = ["RHEL-9.4*HVM-20*"]
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
    values = ["suse-sles-15-sp6-v*-hvm-*"]
  }

  filter {
    name   = "architecture"
    values = [each.value]
  }

  owners = [local.suse_owner_id]
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
