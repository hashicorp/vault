// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

globals {
  archs                = ["amd64", "arm64"]
  artifact_sources     = ["local", "crt", "artifactory"]
  artifact_types       = ["bundle", "package"]
  backends             = ["consul", "raft"]
  backend_license_path = abspath(var.backend_license_path != null ? var.backend_license_path : joinpath(path.root, "./support/consul.hclic"))
  backend_tag_key      = "VaultStorage"
  build_tags = {
    "ce"               = ["ui"]
    "ent"              = ["ui", "enterprise", "ent"]
    "ent.fips1403"     = ["ui", "enterprise", "cgo", "hsm", "fips", "fips_140_3", "ent.fips1403"]
    "ent.hsm"          = ["ui", "enterprise", "cgo", "hsm", "venthsm"]
    "ent.hsm.fips1403" = ["ui", "enterprise", "cgo", "hsm", "fips", "fips_140_3", "ent.hsm.fips1403"]
  }
  config_modes    = ["env", "file"]
  consul_editions = ["ce", "ent"]
  consul_versions = ["1.18.2", "1.19.2", "1.20.6", "1.21.1"]
  distros_aws     = ["amzn", "rhel", "sles", "ubuntu"]
  distros_fyre    = ["centos", "rhel", "rocky", "sles", "ubuntu"]
  // Different distros may require different packages, or use different aliases for the same package
  distro_packages = {
    // NOTE: These versions must always match the output of enos_host_info.target_distro. They are
    // also used in various modules `artifact`, `ec2_info`, and `softhsm_install`. If you are adding
    // or modifying keys you probably have to update those modules.
    amzn = {
      "2"    = ["nc", "openldap-clients"]
      "2023" = ["nc", "openldap-clients"]
    }
    centos = {
      "9"  = ["nc", "openldap-clients"]
      "10" = ["nc", "openldap-clients"]
    }
    rhel = {
      "8.10" = ["nc", "openldap-clients"]
      "9.6"  = ["nc", "openldap-clients"]
      "9.7"  = ["nc", "openldap-clients"]
      "9.8"  = ["nc", "openldap-clients"]
      "10.0" = ["nc", "openldap-clients"]
      "10.1" = ["nc", "openldap-clients"]
      "10.2" = ["nc", "openldap-clients"]
    }
    rocky = {
      "8.10" = ["nc", "openldap-clients"]
      "9.7"  = ["nc", "openldap-clients"]
      "9.8"  = ["nc", "openldap-clients"]
      "10.1" = ["nc", "openldap-clients"]
      "10.2" = ["nc", "openldap-clients"]
    }
    sles = {
      // When installing Vault RPM packages on a SLES AMI, the openssl package provided
      // isn't named "openssl, which rpm doesn't know how to handle. Therefore we add the
      // "correctly" named one in our package installation before installing Vault.
      "15.7" = ["netcat-openbsd", "openssl", "openldap2-client"]
      "16.0" = ["netcat-openbsd", "openssl", "openldap2-client"]
    }
    ubuntu = {
      "22.04" = ["netcat", "ldap-utils"]
      "24.04" = ["netcat-openbsd", "ldap-utils"]
      "26.04" = ["netcat-openbsd", "ldap-utils"]
    }
  }
  distro_version = {
    amzn   = var.distro_version_amzn
    centos = var.distro_version_centos
    rhel   = var.distro_version_rhel
    rocky  = var.distro_version_rhel
    sles   = var.distro_version_sles
    ubuntu = var.distro_version_ubuntu
  }
  # NOTE: These were generated from fyre_os_info for SVL on 6/9/25. They're a
  # moving target so it's likely they'll need to be updated in the future.
  # Make sure distro_packages is updated for any changes that are made here.
  distro_versions_fyre = {
    amd64 = {
      centos = ["10", "9"]
      rhel   = ["10.1", "10.0", "9.6", "8.10"]
      rocky  = ["10.1", "9.7", "8.10"]
      sles   = ["16.0", "15.7"]
      ubuntu = ["26.04", "24.04", "22.04"]
    }
    ppc64le = {
      centos = ["10", "9"]
      rhel   = ["10.2", "10.0", "9.8", "9.6", "8.10"]
      rocky  = ["10.2", "10.1", "9.7", "9.8"]
      sles   = ["16.0", "15.7"]
      ubuntu = ["26.04", "24.04", "22.04"]
    }
    s390x = {
      centos = []
      rhel   = ["10.0", "8.10", "9.6"]
      rocky  = []
      sles   = ["16.0", "15.7"]
      ubuntu = ["24.04", "22.04"]
    }
  }
  # The default versions we each architecture/distro combo. We currently use
  # the latest as that's how things are ordered up above. When we start running
  # these scenarios we'll use the distro_versions_fyre as sample attributes to
  # distribute over all versions.
  distro_version_fyre = {
    amd64 = {
      centos = global.distro_versions_fyre["amd64"]["centos"][0]
      rhel   = global.distro_versions_fyre["amd64"]["rhel"][0]
      rocky  = global.distro_versions_fyre["amd64"]["rocky"][0]
      sles   = global.distro_versions_fyre["amd64"]["sles"][0]
      ubuntu = global.distro_versions_fyre["amd64"]["ubuntu"][0]
    }
    ppc64le = {
      centos = global.distro_versions_fyre["ppc64le"]["centos"][0]
      rhel   = global.distro_versions_fyre["ppc64le"]["rhel"][0]
      rocky  = global.distro_versions_fyre["ppc64le"]["rocky"][0]
      sles   = global.distro_versions_fyre["ppc64le"]["sles"][0]
      ubuntu = global.distro_versions_fyre["ppc64le"]["ubuntu"][0]
    }
    s390x = {
      centos = null
      rhel   = global.distro_versions_fyre["s390x"]["rhel"][0]
      rocky  = null
      sles   = global.distro_versions_fyre["s390x"]["sles"][0]
      ubuntu = global.distro_versions_fyre["s390x"]["ubuntu"][0]
    }
  }
  editions            = ["ce", "ent", "ent.fips1403", "ent.hsm", "ent.hsm.fips1403"]
  enterprise_editions = [for e in global.editions : e if e != "ce"]
  ip_versions         = ["4", "6"]
  package_manager = {
    "amzn"   = "yum"
    "centos" = "yum"
    "rhel"   = "yum"
    "rocky"  = "yum"
    "sles"   = "zypper"
    "ubuntu" = "apt"
  }
  packages = ["jq"]
  // Ports that we'll open up for ingress in the security group for all target machines.
  // Port protocol maps to the IpProtocol schema: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_IpPermission.html

  // Ports that we'll open up for ingress in the security group for all Vault target machines.
  vault_cluster_ports = {
    vault_agent : {
      description = "Vault Agent"
      port        = 8100
      protocol    = "tcp"
    },
    vault_proxy : {
      description = "Vault Proxy"
      port        = 8101
      protocol    = "tcp"
    },
    vault_listener : {
      description = "Vault Addr listener"
      port        = 8200
      protocol    = "tcp"
    },
    vault_cluster : {
      description = "Vault Cluster listener"
      port        = 8201
      protocol    = "tcp"
    }
  }
  // Ports that we'll open up for ingress in the security group for all Consul target machines.
  consul_cluster_ports = {
    consul_rpc : {
      description = "Consul internal communication"
      port        = 8300
      protocol    = "tcp"
    },
    consul_serf_lan_tcp : {
      description = "Consul Serf LAN TCP"
      port        = 8301
      protocol    = "tcp"
    },
    consul_serf_lan_udp : {
      description = "Consul Serf LAN UDP"
      port        = 8301
      protocol    = "udp"
    },
    consul_serf_wan_tcp : {
      description = "Consul Serf WAN TCP"
      port        = 8302
      protocol    = "tcp"
    },
    consul_serf_wan_udp : {
      description = "Consul Serf WAN UDP"
      port        = 8302
      protocol    = "udp"
    },
    consul_http : {
      description = "Consul HTTP API"
      port        = 8500
      protocol    = "tcp"
    },
    consul_https : {
      description = "Consul HTTPS API"
      port        = 8501
      protocol    = "tcp"
    },
    consul_grpc : {
      description = "Consul gRPC API"
      port        = 8502
      protocol    = "tcp"
    },
    consul_grpc_tls : {
      description = "Consul gRPC TLS API"
      port        = 8503
      protocol    = "tcp"
    },
    consul_dns_tcp : {
      description = "Consul TCP DNS Server"
      port        = 8600
      protocol    = "tcp"
    },
    consul_dns_udp : {
      description = "Consul UDP DNS Server"
      port        = 8600
      protocol    = "udp"
    }
  }
  ingress_ports = {
    ssh : {
      description = "SSH"
      port        = 22
      protocol    = "tcp"
    }
  }
  // Database configurations for Vault database secrets engine testing
  database_configs = {
    postgres : {
      port        = 5432
      version     = "16-alpine"
      username    = "postgres"
      password    = "secret"
      database    = "postgres"
      description = "PostgreSQL Server"
    },
    mongodb : {
      port        = 27017
      version     = "7.0"
      username    = "admin"
      password    = "secret"
      database    = "admin"
      description = "MongoDB Server"
    },
    mysql : {
      port        = 3306
      version     = "8.0"
      username    = "root"
      password    = "secret"
      database    = "mysql"
      description = "MySQL Server"
    }
  }

  // Ports that we'll open up for ingress in the security group for all external integration target machines.
  integration_host_ports = merge(
    {
      ldap : {
        description = "LDAP"
        port        = 389
        protocol    = "tcp"
      },
      ldaps : {
        description = "LDAPS"
        port        = 636
        protocol    = "tcp"
      },
      kmip : {
        description = "KMIP Server"
        port        = 5696
        protocol    = "tcp"
      }
    },
    // Add database ports dynamically from database_configs
    { for db_name, db_config in global.database_configs : db_name => {
      description = db_config.description
      port        = db_config.port
      protocol    = "tcp"
    } }
  )

  // Combine all ports into a single map
  ports = merge(
    global.vault_cluster_ports,
    global.consul_cluster_ports,
    global.ingress_ports,
    global.integration_host_ports
  )

  seals = ["awskms", "pkcs11", "shamir"]
  tags = merge({
    "Project Name" : var.project_name
    "Project" : "Enos",
    "Environment" : "ci"
  }, var.tags)
  vault_install_dir = {
    bundle  = "/opt/vault/bin"
    package = "/usr/bin"
  }
  vault_license_path  = abspath(var.vault_license_path != null ? var.vault_license_path : joinpath(path.root, "./support/vault.hclic"))
  vault_tag_key       = "vault-cluster"
  vault_disable_mlock = false
}
