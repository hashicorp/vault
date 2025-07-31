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
  test_server_address = var.ip_version == "6" ? var.hosts[0].ipv6 : var.hosts[0].public_ip
  ldap_server = {
    domain      = "enos.com"
    org         = "hashicorp"
    admin_pw    = "password1"
    version     = var.ldap_version
    port        = var.ports.ldap.port
    secure_port = var.ports.ldaps.port
    ip_version  = var.ip_version
    host        = var.hosts[0]
  }
  kmip_client = {
    // The KMIP client configuration is used to connect to the KMIP server
    // uses Percona (MySQL) as the KMIP client.
    port = var.ports.mysql.port
    host = var.hosts[0]
  }
}

# Outputs
output "state" {
  value = {
    ldap = local.ldap_server
    kmip = local.kmip_client
  }
}

# We run install_packages before we install Vault because for some combinations of
# certain Linux distros and artifact types (e.g. SLES and RPM packages), there may
# be packages that are required to perform Vault installation (e.g. openssl).
module "install_packages" {
  source   = "../install_packages"
  hosts    = var.hosts
  packages = var.packages
}

# Creating OpenLDAP Server
resource "enos_remote_exec" "setup_openldap" {
  depends_on = [module.install_packages]

  environment = {
    LDAP_CONTAINER_VERSION = local.ldap_server.version
    LDAP_DOMAIN            = local.ldap_server.domain
    LDAP_ORG               = local.ldap_server.org
    LDAP_ADMIN_PW          = local.ldap_server.admin_pw
    LDAP_IP_ADDRESS        = local.test_server_address
    LDAP_PORT              = local.ldap_server.port
    LDAPS_PORT             = local.ldap_server.secure_port
  }

  scripts = [abspath("${path.module}/scripts/set-up-openldap.sh")]

  transport = {
    ssh = {
      host = local.ldap_server.host.public_ip
    }
  }
}

# Creating KMIP Server
resource "enos_remote_exec" "create_kmip" {
  depends_on = [module.install_packages]

  environment = {
    VAULT_ADDR = var.ip_version == "6" ? var.hosts[0].ipv6 : var.hosts[0].public_ip
    KMIP_PORT  = var.ports.kmip.port
  }

  scripts = [abspath("${path.module}/scripts/setup_kmip.sh")]

  transport = {
    ssh = {
      host = local.kmip_client.host.public_ip
    }
  }
}
