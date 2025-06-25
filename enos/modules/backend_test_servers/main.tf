# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

locals {
  distro = var.distro
  ldap_server = {
    domain     = "enos.com"
    org        = "hashicorp"
    admin_pw   = "password1"
    version    = var.ldap_version
    port       = "389"
    ip_address = var.ip_version == "6" ? var.hosts[0].ipv6_addresses[0] : var.hosts[0].public_ip
  }
}

# Outputs
output "state" {
  value = {
    ldap = local.ldap_server
  }
}

# Creating OpenLDAP Server
resource "enos_remote_exec" "setup_docker" {
  scripts = [abspath("${path.module}/scripts/set-up-docker.sh")]

  environment = {
    DISTRO = local.distro
  }

  transport = {
    ssh = {
      host = local.ldap_server.ip_address
    }
  }
}

# Creating OpenLDAP Server
resource "enos_remote_exec" "setup_openldap" {
  depends_on = [enos_remote_exec.setup_docker]

  environment = {
    LDAP_CONTAINER_VERSION = local.ldap_server.version
    LDAP_DOMAIN            = local.ldap_server.domain
    LDAP_ORG               = local.ldap_server.org
    LDAP_ADMIN_PW          = local.ldap_server.admin_pw
    LDAP_PORT              = local.ldap_server.port
  }

  scripts = [abspath("${path.module}/scripts/set-up-openldap.sh")]

  transport = {
    ssh = {
      host = local.ldap_server.ip_address
    }
  }
}