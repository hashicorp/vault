# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "hashicorp-forge/enos"
    }
  }
}

data "aws_ec2_instance_type_offerings" "instance" {
  filter {
    name   = "instance-type"
    values = [var.k6_instance_type, var.metrics_instance_type]
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

resource "aws_instance" "k6" {
  ami                    = var.ami_id
  instance_type          = var.k6_instance_type
  key_name               = var.ssh_keypair
  subnet_id              = data.aws_subnets.vpc.ids[0]
  vpc_security_group_ids = [var.metrics_security_group_ids["prometheus"], var.target_security_group_id]

  tags = {
    Name = "${var.project_name}-k6"
  }
}

resource "aws_instance" "metrics" {
  ami                    = var.ami_id
  instance_type          = var.metrics_instance_type
  key_name               = var.ssh_keypair
  subnet_id              = data.aws_subnets.vpc.ids[0]
  vpc_security_group_ids = [var.metrics_security_group_ids["grafana"], var.target_security_group_id]

  tags = {
    Name = "${var.project_name}-metrics"
  }
}

resource "enos_remote_exec" "install_k6" {
  depends_on = [aws_instance.k6]

  environment = {
    RETRY_INTERVAL  = var.retry_interval
    TIMEOUT_SECONDS = var.timeout
  }

  scripts = [abspath("${path.module}/scripts/install-k6.sh")]

  transport = {
    ssh = {
      host = aws_instance.k6.public_ip
    }
  }
}

resource "enos_remote_exec" "add_telemetry_to_consul" {
  depends_on = [
    aws_instance.metrics,
    aws_instance.k6
  ]

  for_each = var.consul_hosts

  scripts = [abspath("${path.module}/scripts/add-consul-telemetry.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

resource "enos_remote_exec" "install_prometheus_node_exporter" {
  depends_on = [
    aws_instance.metrics,
    aws_instance.k6
  ]
  for_each = local.all_hosts

  environment = {
    PROMETHEUS_NODE_EXPORTER_VERSION = var.prometheus_node_exporter_version
    RETRY_INTERVAL                   = var.retry_interval
    TIMEOUT_SECONDS                  = var.timeout
  }

  scripts = [abspath("${path.module}/scripts/install-prometheus-node-exporter.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

resource "enos_remote_exec" "run_prometheus_node_exporter" {
  depends_on = [enos_remote_exec.install_prometheus_node_exporter]
  for_each   = local.all_hosts

  environment = {
    PROMETHEUS_NODE_EXPORTER_VERSION = var.prometheus_node_exporter_version
  }

  scripts = [abspath("${path.module}/scripts/run-prometheus-node-exporter.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

module "restart_all_consuls" {
  depends_on = [enos_remote_exec.add_telemetry_to_consul]
  source     = "../restart_consul"
  hosts      = var.consul_hosts
}

resource "enos_remote_exec" "install_prometheus" {
  depends_on = [
    aws_instance.metrics,
    aws_instance.k6
  ]

  environment = local.prometheus_environment
  scripts     = [abspath("${path.module}/scripts/install-prometheus.sh")]

  transport = {
    ssh = {
      host = aws_instance.metrics.public_ip
    }
  }
}

resource "enos_remote_exec" "install_grafana" {
  depends_on = [enos_remote_exec.install_prometheus]

  environment = {
    GRAFANA_VERSION = var.grafana_version
    RETRY_INTERVAL  = var.retry_interval
    TIMEOUT_SECONDS = var.timeout
  }

  scripts = [abspath("${path.module}/scripts/install-grafana.sh")]

  transport = {
    ssh = {
      host = aws_instance.metrics.public_ip
    }
  }
}

resource "enos_file" "grafana_dashboards" {
  depends_on  = [enos_remote_exec.install_grafana]
  for_each    = fileset(abspath("${path.module}/grafana-dashboards"), "*.json")
  source      = abspath("${path.module}/grafana-dashboards/${each.value}")
  destination = "/etc/grafana/dashboards/${basename(each.value)}"

  transport = {
    ssh = {
      host = aws_instance.metrics.public_ip
    }
  }
}

resource "enos_remote_exec" "run_prometheus" {
  depends_on = [enos_remote_exec.install_grafana]

  scripts = [abspath("${path.module}/scripts/run-prometheus.sh")]

  transport = {
    ssh = {
      host = aws_instance.metrics.public_ip
    }
  }
}

resource "enos_remote_exec" "run_grafana" {
  depends_on = [enos_file.grafana_dashboards]

  environment = {
    RETRY_INTERVAL  = var.retry_interval
    TIMEOUT_SECONDS = var.timeout
  }

  scripts = [abspath("${path.module}/scripts/run-grafana.sh")]

  transport = {
    ssh = {
      host = aws_instance.metrics.public_ip
    }
  }
}

resource "enos_file" "k6_scripts" {
  depends_on = [enos_remote_exec.install_k6]

  for_each    = fileset(abspath("${path.module}/k6-templates"), "*.tpl")
  destination = "/home/ubuntu/scripts/${replace(basename(each.value), ".tpl", "")}"
  content = templatefile("${path.module}/k6-templates/${each.value}", {
    hosts       = var.hosts
    vault_token = var.vault_token
    leader_addr = var.leader_addr
  })

  transport = {
    ssh = {
      host = aws_instance.k6.public_ip
    }
  }
}

resource "enos_file" "k6_exec_script" {
  depends_on = [enos_remote_exec.install_k6]

  chmod       = "755"
  destination = "/home/ubuntu/k6-run.sh"
  content = templatefile("${path.module}/scripts/k6-run.sh.tpl", {
    metrics_addr = aws_instance.metrics.private_ip
  })

  transport = {
    ssh = {
      host = aws_instance.k6.public_ip
    }
  }
}

