# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "hashicorp-forge/enos"
    }
  }
}

variable "cluster_id" {
  type = string
}

variable "host" {
  type = object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  })
}

variable "leader_addr" {
  type = string
}

variable "metrics_collector_host" {
  type = object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  })
}

variable "retry_interval" {
  type        = number
  description = "How many seconds to wait between each retry"
  default     = 2
}

variable "timeout" {
  type        = number
  description = "The max number of seconds to wait before timing out. This is applied to each step so total timeout will be longer."
  default     = 120
}

variable "vault_hosts" {
  type = map(object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  }))
}

variable "vault_token" {
  type = string
}

resource "random_string" "k6" {
  length  = 8
  numeric = false
  special = false
  upper   = false
}

locals {
  k6_id = "${var.cluster_id}-${random_string.k6.result}"
}

resource "enos_remote_exec" "install_k6" {
  environment = {
    RETRY_INTERVAL  = var.retry_interval
    TIMEOUT_SECONDS = var.timeout
  }

  scripts = [abspath("${path.module}/scripts/install-k6.sh")]

  transport = {
    ssh = {
      host = var.host.public_ip
    }
  }
}

resource "enos_file" "k6_scripts" {
  depends_on = [enos_remote_exec.install_k6]

  for_each    = fileset(abspath("${path.module}/../k6-templates"), "*.tpl")
  destination = "/home/ubuntu/scripts/${replace(basename(each.value), ".tpl", "")}"
  content = templatefile("${path.module}/k6-templates/${each.value}", {
    hosts       = var.vault_hosts
    vault_token = var.vault_token
    leader_addr = var.leader_addr
  })

  transport = {
    ssh = {
      host = var.host.public_ip
    }
  }
}

resource "enos_file" "k6_exec_script" {
  depends_on = [enos_remote_exec.install_k6]

  chmod       = "755"
  destination = "/home/ubuntu/k6-run.sh"
  content = templatefile("${path.module}/scripts/k6-run.sh.tpl", {
    metrics_addr = var.metrics_collector_host.private_ip
  })

  transport = {
    ssh = {
      host = var.host.public_ip
    }
  }
}

output "host" {
  value = var.host
}

output "hosts" {
  value = { "k6" : var.host }
}

output "public_ip" {
  value = var.host.public_ip
}
