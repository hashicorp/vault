# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    # We need to specify the provider source in each module until we publish it
    # to the public registry
    enos = {
      source  = "registry.terraform.io/hashicorp-forge/enos"
      version = ">= 0.3.24"
    }
  }
}

data "aws_vpc" "vpc" {
  id = var.vpc_id
}

data "aws_ami" "ami" {
  filter {
    name   = "image-id"
    values = [var.ami_id]
  }
}

data "aws_ec2_instance_type_offerings" "instance" {
  filter {
    name   = "instance-type"
    values = [local.instance_type]
  }

  location_type = "availability-zone"
}

data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "zone-name"
    values = data.aws_ec2_instance_type_offerings.instance.locations
  }
}

data "aws_subnets" "vpc" {
  filter {
    name   = "availability-zone"
    values = data.aws_availability_zones.available.names
  }

  filter {
    name   = "vpc-id"
    values = [var.vpc_id]
  }
}

data "aws_iam_policy_document" "target" {
  statement {
    resources = ["*"]

    actions = [
      "ec2:DescribeInstances",
      "secretsmanager:*"
    ]
  }

  dynamic "statement" {
    for_each = var.seal_key_names

    content {
      resources = [statement.value]

      actions = [
        "kms:DescribeKey",
        "kms:ListKeys",
        "kms:Encrypt",
        "kms:Decrypt",
        "kms:GenerateDataKey"
      ]
    }
  }
}

data "aws_iam_policy_document" "target_instance_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

data "enos_environment" "localhost" {}

locals {
  cluster_name  = coalesce(var.cluster_name, random_string.cluster_name.result)
  instance_type = local.instance_types[data.aws_ami.ami.architecture]
  instance_types = {
    "arm64"  = var.instance_types["arm64"]
    "x86_64" = var.instance_types["amd64"]
  }
  instances   = toset([for idx in range(var.instance_count) : tostring(idx)])
  name_prefix = "${var.project_name}-${local.cluster_name}-${random_string.unique_id.result}"
}

resource "random_string" "cluster_name" {
  length  = 8
  lower   = true
  upper   = false
  numeric = false
  special = false
}

resource "random_string" "unique_id" {
  length  = 4
  lower   = true
  upper   = false
  numeric = false
  special = false
}

resource "aws_iam_role" "target_instance_role" {
  name               = "${local.name_prefix}-instance-role"
  assume_role_policy = data.aws_iam_policy_document.target_instance_role.json
}

resource "aws_iam_instance_profile" "target" {
  name = "${local.name_prefix}-instance-profile"
  role = aws_iam_role.target_instance_role.name
}

resource "aws_iam_role_policy" "target" {
  name   = "${local.name_prefix}-role-policy"
  role   = aws_iam_role.target_instance_role.id
  policy = data.aws_iam_policy_document.target.json
}

resource "aws_security_group" "target" {
  name        = "${local.name_prefix}-sg"
  description = "Target instance security group"
  vpc_id      = var.vpc_id

  # External ingress
  dynamic "ingress" {
    for_each = var.ports_ingress

    content {
      from_port = ingress.value.port
      to_port   = ingress.value.port
      protocol  = ingress.value.protocol
      cidr_blocks = flatten([
        formatlist("%s/32", data.enos_environment.localhost.public_ipv4_addresses),
        join(",", data.aws_vpc.vpc.cidr_block_associations.*.cidr_block),
        formatlist("%s/32", var.ssh_allow_ips)
      ])
      ipv6_cidr_blocks = data.aws_vpc.vpc.ipv6_cidr_block != "" ? [data.aws_vpc.vpc.ipv6_cidr_block] : null
    }
  }

  # Internal traffic
  ingress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"
    self      = true
  }

  # External traffic
  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = merge(
    var.common_tags,
    {
      Name = "${local.name_prefix}-sg"
    },
  )
}

resource "aws_instance" "targets" {
  for_each = local.instances

  ami                  = var.ami_id
  iam_instance_profile = aws_iam_instance_profile.target.name
  // Some scenarios (autopilot, pr_replication) shutdown instances to simulate failure. In those
  // cases we should terminate the instance entirely rather than get stuck in stopped limbo.
  instance_initiated_shutdown_behavior = "terminate"
  instance_type                        = local.instance_type
  key_name                             = var.ssh_keypair
  subnet_id                            = data.aws_subnets.vpc.ids[tonumber(each.key) % length(data.aws_subnets.vpc.ids)]
  vpc_security_group_ids               = [aws_security_group.target.id]

  tags = merge(
    var.common_tags,
    {
      Name                     = "${local.name_prefix}-${var.cluster_tag_key}-instance-target"
      "${var.cluster_tag_key}" = local.cluster_name
    },
  )
}

module "disable_selinux" {
  depends_on = [aws_instance.targets]
  source     = "../disable_selinux"
  count      = var.disable_selinux == true ? 1 : 0

  hosts = local.hosts
}
