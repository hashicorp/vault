# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "hashicorp-forge/enos"
    }
  }
}

locals {
  grafana_ports    = [3000, 9090]
  prometheus_ports = [9100, 8500]
}

data "enos_environment" "localhost" {}

data "aws_vpc" "vpc" {
  id = var.vpc_id
}

resource "aws_security_group" "grafana" {
  name        = "grafana-sg"
  description = "A security group specifically for viewing grafana dashboards"
  vpc_id      = var.vpc_id

  dynamic "ingress" {
    for_each = local.grafana_ports
    content {
      from_port = ingress.value
      to_port   = ingress.value
      protocol  = "tcp"
      cidr_blocks = flatten([
        formatlist("%s/32", data.enos_environment.localhost.public_ipv4_addresses),
        join(",", data.aws_vpc.vpc.cidr_block_associations.*.cidr_block),
      ])
      ipv6_cidr_blocks = data.aws_vpc.vpc.ipv6_cidr_block != "" ? [data.aws_vpc.vpc.ipv6_cidr_block] : null
    }
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.project_name}-grafana-sg"
  }
}

resource "aws_security_group" "prometheus" {
  name        = "prometheus-sg"
  description = "A security group specifically for scraping prometheus metrics"
  vpc_id      = var.vpc_id

  dynamic "ingress" {
    for_each = local.prometheus_ports
    content {
      from_port   = ingress.value
      to_port     = ingress.value
      protocol    = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
    }
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.project_name}-prometheus-sg"
  }
}
