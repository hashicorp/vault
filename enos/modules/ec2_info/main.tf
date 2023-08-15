# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1
  
# Note: in order to use the openSUSE Leap AMIs, the AWS account in use must have an
# active subscription. There is no additional charge for using these instances other
# than the usage AWS usage charges, but because the images belong to an AWS Marketplace 
# owner (679593333241), you must accept their terms and conditions.
# openSUSE Leap arm64 subscription: https://aws.amazon.com/marketplace/server/procurement?productId=a516e959-df54-4035-bb1a-63599b7a6df9
# opensuse leap amd64 subscription: https://aws.amazon.com/marketplace/server/procurement?productId=5535c495-72d4-4355-b169-54ffa874f849

locals {
  architectures      = toset(["arm64", "x86_64"])
  amazon_owner_id    = "591542846629"
  canonical_owner_id = "099720109477"
  sles_owner_id      = "013907871322"
  suse_owner_id      = "679593333241"
  rhel_owner_id      = "309956199498"
  ids = {
    "arm64" = {
      "rhel" = {
        "8.8" = data.aws_ami.rhel_88["arm64"].id
        "9.1" = data.aws_ami.rhel_91["arm64"].id
      }
      "ubuntu" = {
        "18.04" = data.aws_ami.ubuntu_1804["arm64"].id
        "20.04" = data.aws_ami.ubuntu_2004["arm64"].id
        "22.04" = data.aws_ami.ubuntu_2204["arm64"].id
      }
      "amazon_linux" = {
        "amzn2" = data.aws_ami.amazon_linux_2["arm64"].id
      }
      "leap" = {
        "15.4" = data.aws_ami.leap_154["arm64"].id
        "15.5" = data.aws_ami.leap_155["arm64"].id
      }
    }
    "amd64" = {
      "rhel" = {
        "7.9" = data.aws_ami.rhel_79.id
        "8.8" = data.aws_ami.rhel_88["x86_64"].id
        "9.1" = data.aws_ami.rhel_91["x86_64"].id
      }
      "ubuntu" = {
        "18.04" = data.aws_ami.ubuntu_1804["x86_64"].id
        "20.04" = data.aws_ami.ubuntu_2004["x86_64"].id
        "22.04" = data.aws_ami.ubuntu_2204["x86_64"].id
      }
      "amazon_linux" = {
        "amzn2" = data.aws_ami.amazon_linux_2["x86_64"].id
      }
      "leap" = {
        "15.4" = data.aws_ami.leap_154["x86_64"].id
        "15.5" = data.aws_ami.leap_155["x86_64"].id
      }
      "sles" = {
        "v12_sp5_standard" = data.aws_ami.sles_12_sp5_standard.id
        "v15_sp4_standard" = data.aws_ami.sles_15_sp4_standard["x86_64"].id
        "v15_sp5_standard" = data.aws_ami.sles_15_sp5_standard["x86_64"].id
      }
    }
  }
}

data "aws_ami" "ubuntu_1804" {
  most_recent = true
  for_each    = local.architectures

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-*-18.04-*-server-*"]
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

data "aws_ami" "rhel_79" {
  most_recent = true

  # Currently latest latest point release-1
  filter {
    name   = "name"
    values = ["RHEL-7.9*HVM-20*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name   = "architecture"
    values = ["x86_64"]
  }

  owners = [local.rhel_owner_id]
}

data "aws_ami" "rhel_88" {
  most_recent = true
  for_each    = local.architectures

  # Currently latest latest point release-1
  filter {
    name   = "name"
    values = ["RHEL-8.8*HVM-20*"]
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

data "aws_ami" "rhel_91" {
  most_recent = true
  for_each    = local.architectures

  # Currently latest latest point release-1
  filter {
    name   = "name"
    values = ["RHEL-9.1*HVM-20*"]
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

data "aws_ami" "amazon_linux_2" {
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

data "aws_ami" "sles_12_sp5_standard" {
  most_recent = true

  filter {
    name   = "description"
    values = ["SUSE Linux Enterprise Server 12 SP5 (HVM*"]
  }

  filter {
    name   = "architecture"
    # arm64 only available for BYOS images for SLES 12 SP5
    values = ["x86_64"]
  }

  owners = [local.sles_owner_id]
}

data "aws_ami" "sles_15_sp4_standard" {
  most_recent = true
  for_each    = local.architectures

  filter {
    name   = "description"
    values = ["SUSE Linux Enterprise Server 15 SP4 (HVM*"]
  }

  filter {
    name   = "architecture"
    values = [each.value]
  }

  owners = [local.sles_owner_id]
}

data "aws_ami" "sles_15_sp5_standard" {
  most_recent = true
  for_each    = local.architectures

  filter {
    name   = "description"
    values = ["SUSE Linux Enterprise Server 15 SP5 (HVM*"]
  }

  filter {
    name   = "architecture"
    values = [each.value]
  }

  owners = [local.sles_owner_id]
}

data "aws_ami" "leap_154" {
  most_recent = true
  for_each    = local.architectures

  filter {
    name   = "name"
    values = ["openSUSE-Leap-15-4*"]
  }

  filter {
    name   = "architecture"
    values = [each.value]
  }

  owners = [local.suse_owner_id]
}

data "aws_ami" "leap_155" {
  most_recent = true
  for_each    = local.architectures

  filter {
    name   = "name"
    values = ["openSUSE-Leap-15-5*"]
  }

  filter {
    name   = "architecture"
    values = [each.value]
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
