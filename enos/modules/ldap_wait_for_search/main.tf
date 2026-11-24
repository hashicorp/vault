# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "hosts" {
  description = "The target machines to run the test query from"
  type = map(object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  }))
}

variable "ldap_base_dn" {
  type        = string
  description = "The LDAP base dn to search from"
}

variable "ldap_bind_dn" {
  type        = string
  description = "The LDAP bind dn"
}

variable "ldap_host" {
  type = object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  })
  description = "The LDAP host"
}

variable "ldap_password" {
  type        = string
  description = "The LDAP password"
}

variable "ldap_port" {
  type        = string
  description = "The LDAP port"
}

variable "ldap_query" {
  type        = string
  description = "The LDAP query to use when testing the connection"
  default     = null
}

variable "retry_interval" {
  type        = number
  description = "How many seconds to wait between each retry"
  default     = 2
}

variable "timeout" {
  type        = number
  description = "The max number of seconds to wait before timing out"
  default     = 60
}

# Wait for the search to succeed
resource "enos_remote_exec" "wait_for_search" {
  for_each = var.hosts

  environment = {
    LDAP_BASE_DN    = var.ldap_base_dn
    LDAP_BIND_DN    = var.ldap_bind_dn
    LDAP_HOST       = var.ldap_host.public_ip
    LDAP_PASSWORD   = var.ldap_password
    LDAP_PORT       = var.ldap_port
    LDAP_QUERY      = var.ldap_query == null ? "" : var.ldap_query
    RETRY_INTERVAL  = var.retry_interval
    TIMEOUT_SECONDS = var.timeout
  }
  scripts = [abspath("${path.module}/scripts/wait-for-search.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
