# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "zone-name"
    values = ["*"]
  }
}

resource "random_string" "cluster_id" {
  length  = 8
  lower   = true
  upper   = false
  numeric = false
  special = false
}

resource "aws_vpc" "vpc" {
  // Always set the ipv4 cidr block as it's required in "dual-stack" VPCs which we create.
  cidr_block                       = var.ipv4_cidr
  enable_dns_hostnames             = true
  enable_dns_support               = true
  assign_generated_ipv6_cidr_block = var.ip_version == 6

  tags = merge(
    var.common_tags,
    {
      "Name" = var.name
    },
  )
}

resource "aws_subnet" "subnet" {
  count             = length(data.aws_availability_zones.available.names)
  vpc_id            = aws_vpc.vpc.id
  availability_zone = data.aws_availability_zones.available.names[count.index]

  // IPV4, but since we need to support ipv4 connections from the machine running enos, we're
  // always going to need ipv4 available.
  map_public_ip_on_launch = true
  cidr_block              = cidrsubnet(var.ipv4_cidr, 8, count.index)

  // IPV6, only set these when we want to run in ipv6 mode.
  assign_ipv6_address_on_creation = var.ip_version == 6
  ipv6_cidr_block                 = var.ip_version == 6 ? cidrsubnet(aws_vpc.vpc.ipv6_cidr_block, 4, count.index) : null

  tags = merge(
    var.common_tags,
    {
      "Name" = "${var.name}-subnet-${data.aws_availability_zones.available.names[count.index]}"
    },
  )
}

resource "aws_internet_gateway" "ipv4" {
  vpc_id = aws_vpc.vpc.id

  tags = merge(
    var.common_tags,
    {
      "Name" = "${var.name}-igw"
    },
  )
}

resource "aws_egress_only_internet_gateway" "ipv6" {
  count  = var.ip_version == 6 ? 1 : 0
  vpc_id = aws_vpc.vpc.id
}

resource "aws_route" "igw_ipv4" {
  route_table_id         = aws_vpc.vpc.default_route_table_id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.ipv4.id
}

resource "aws_route" "igw_ipv6" {
  count                       = var.ip_version == 6 ? 1 : 0
  route_table_id              = aws_vpc.vpc.default_route_table_id
  destination_ipv6_cidr_block = "::/0"
  egress_only_gateway_id      = aws_egress_only_internet_gateway.ipv6[0].id
}

resource "aws_security_group" "default" {
  vpc_id = aws_vpc.vpc.id

  ingress {
    description      = "allow_ingress_from_all"
    from_port        = 0
    to_port          = 0
    protocol         = "tcp"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = var.ip_version == 6 ? ["::/0"] : null
  }

  egress {
    description      = "allow_egress_from_all"
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = var.ip_version == 6 ? ["::/0"] : null
  }

  tags = merge(
    var.common_tags,
    {
      "Name" = "${var.name}-default"
    },
  )
}
