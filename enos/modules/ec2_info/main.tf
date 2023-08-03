locals {
  architectures      = toset(["arm64", "x86_64"])
  amazon_owner_id    = "013907871322"
  canonical_owner_id = "099720109477"
  openSUSE_owner_id  = "679593333241"
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
      "openSUSE" = {
        "15.5" = data.aws_ami.openSUSE_15.5["arm64"].id
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
      "openSUSE" = {
        "15.4" = data.aws_ami.openSUSE_15.4["x86_64"].id
        "15.5" = data.aws_ami.openSUSE_15.5["x86_64"].id
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
    values = ["amzn2-ami-hvm*"]
  }

  filter {
    name   = "architecture"
    values = [each.value]
  }

  owners = [local.amazon_owner_id]
}

data "aws_ami" "openSUSE_15.5" {
  most_recent = true
  for_each    = local.architectures

  filter {
    name   = "name"
    values = ["openSUSE-Leap-15-5-*-hvm-*"]
  }

  filter {
    name   = "architecture"
    values = [each.value]
  }

  owners = [local.openSUSE_owner_id]
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
