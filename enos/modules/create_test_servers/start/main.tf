# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

locals {
  // Variables
  instance_types = {
    "arm64"  = var.instance_types["arm64"]
    "x86_64" = var.instance_types["amd64"]
  }
  instance_type = local.instance_types[data.aws_ami.amazon_linux.architecture]
  name_prefix = "enos-test-server-${random_id.unique_suffix.hex}"

}

data "aws_region" "current" {}

# Generates a unique suffix for the EC2 tag
resource "random_id" "unique_suffix" {
  byte_length = 4
}

# Lookup latest Amazon Linux 2 AMI dynamically
data "aws_ami" "amazon_linux" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "${local.name_prefix}-ami"
    values = ["amzn2-ami-hvm-*-x86_64-gp2"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
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
    name   = "${local.name_prefix}-zone"
    values = data.aws_ec2_instance_type_offerings.instance.locations
  }
}

# Get subnets in the given VPC across available AZs
data "aws_subnets" "vpc_subnets" {
  filter {
    name   = "${local.name_prefix}-subnets"
    values = data.aws_availability_zones.available.names
  }

  filter {
    name   = "${local.name_prefix}-vpc-id"
    values = [var.vpc_id]
  }
}

# Create the EC2 instance
resource "aws_instance" "enos_test_server" {
  ami                                  = data.aws_ami.amazon_linux.id
  key_name                             = var.ssh_keypair
  instance_initiated_shutdown_behavior = "terminate"
  instance_type                        = local.instance_type
  subnet_id                            = element(data.aws_subnets.vpc_subnets.ids, 0)
  vpc_security_group_ids               = var.vpc_security_group_ids

  # Install Docker using user_data script
  user_data = <<-EOF
              #!/bin/bash
              yum update -y
              amazon-linux-extras install docker -y
              systemctl start docker
              systemctl enable docker
              usermod -aG docker ec2-user
              EOF

  tags = {
    Name = "${local.name_prefix}-ec2"
  }
}