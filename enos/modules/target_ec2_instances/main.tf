terraform {
  required_providers {
    # We need to specify the provider source in each module until we publish it
    # to the public registry
    enos = {
      source  = "app.terraform.io/hashicorp-qti/enos"
      version = ">= 0.3.2"
    }
  }
}

data "aws_vpc" "vpc" {
  id = var.vpc_id
}

data "aws_subnets" "vpc" {
  filter {
    name   = "vpc-id"
    values = [var.vpc_id]
  }
}

data "aws_kms_key" "kms_key" {
  key_id = var.awskms_unseal_key_arn
}

data "aws_iam_policy_document" "target" {
  statement {
    resources = ["*"]

    actions = [
      "ec2:DescribeInstances",
      "secretsmanager:*"
    ]
  }

  statement {
    resources = [var.awskms_unseal_key_arn]

    actions = [
      "kms:DescribeKey",
      "kms:ListKeys",
      "kms:Encrypt",
      "kms:Decrypt",
      "kms:GenerateDataKey"
    ]
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

resource "random_string" "cluster_name" {
  length  = 8
  lower   = true
  upper   = false
  numeric = false
  special = false
}

locals {
  instances    = toset([for idx in range(var.instance_count) : tostring(idx)])
  cluster_name = coalesce(var.cluster_name, random_string.cluster_name.result)
  name_prefix  = "${var.project_name}-${local.cluster_name}"
}

resource "aws_iam_role" "target_instance_role" {
  name               = "target_instance_role-${random_string.cluster_name.result}"
  assume_role_policy = data.aws_iam_policy_document.target_instance_role.json
}

resource "aws_iam_instance_profile" "target" {
  name = "${local.name_prefix}-target"
  role = aws_iam_role.target_instance_role.name
}

resource "aws_iam_role_policy" "target" {
  name   = "${local.name_prefix}-target"
  role   = aws_iam_role.target_instance_role.id
  policy = data.aws_iam_policy_document.target.json
}

resource "aws_security_group" "target" {
  name        = "${local.name_prefix}-target"
  description = "Target instance security group"
  vpc_id      = var.vpc_id

  # SSH traffic
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["${data.enos_environment.localhost.public_ip_address}/32", join(",", data.aws_vpc.vpc.cidr_block_associations.*.cidr_block)]
  }

  # Vault traffic
  ingress {
    from_port = 8200
    to_port   = 8201
    protocol  = "tcp"
    cidr_blocks = flatten([
      "${data.enos_environment.localhost.public_ip_address}/32",
      join(",", data.aws_vpc.vpc.cidr_block_associations.*.cidr_block),
    formatlist("%s/32", var.ssh_allow_ips)])
  }

  # Consul traffic
  ingress {
    from_port   = 8301
    to_port     = 8301
    protocol    = "tcp"
    cidr_blocks = ["${data.enos_environment.localhost.public_ip_address}/32", join(",", data.aws_vpc.vpc.cidr_block_associations.*.cidr_block)]
  }

  ingress {
    from_port   = 8301
    to_port     = 8301
    protocol    = "udp"
    cidr_blocks = ["${data.enos_environment.localhost.public_ip_address}/32", join(",", data.aws_vpc.vpc.cidr_block_associations.*.cidr_block)]
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
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(
    var.common_tags,
    {
      Name = "${local.name_prefix}-sg"
    },
  )
}

resource "aws_instance" "targets" {
  for_each               = local.instances
  ami                    = var.ami_id
  instance_type          = var.instance_type
  vpc_security_group_ids = [aws_security_group.target.id]
  subnet_id              = tolist(data.aws_subnets.vpc.ids)[each.key % length(data.aws_subnets.vpc.ids)]
  key_name               = var.ssh_keypair
  iam_instance_profile   = aws_iam_instance_profile.target.name
  tags = merge(
    var.common_tags,
    {
      Name = "${local.name_prefix}-target-instance"
      Type = local.cluster_name
    },
  )
}
