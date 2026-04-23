# Copyright IBM Corp. 2016, 2026
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

locals {
  # Database-specific configurations
  database_configs = {
    postgres = {
      image_template = "docker.io/postgres:${var.db_version}"
      env_vars = {
        POSTGRES_USER     = var.username
        POSTGRES_PASSWORD = var.password
        POSTGRES_DB       = var.database
      }
    }
    mongodb = {
      image_template = "docker.io/mongo:${var.db_version}"
      env_vars = {
        MONGO_INITDB_ROOT_USERNAME = var.username
        MONGO_INITDB_ROOT_PASSWORD = var.password
        MONGO_INITDB_DATABASE      = var.database
      }
      args = "--bind_ip_all"
    }
    mysql = {
      image_template = "docker.io/mysql:${var.db_version}"
      env_vars = {
        MYSQL_ROOT_PASSWORD = var.password
        MYSQL_USER          = var.username
        MYSQL_PASSWORD      = var.password
        MYSQL_DATABASE      = var.database
      }
    }
  }

  config       = local.database_configs[var.database_type]
  image        = local.config.image_template
  env_vars_map = local.config.env_vars
  env_vars     = join(",", [for k, v in local.env_vars_map : "${k}=${v}"])
  args         = try(local.config.args, "")
}

# Creating Database Server using generic container script
resource "enos_remote_exec" "create_database" {
  depends_on = [var.depends_on_modules]

  scripts = [abspath("${path.module}/../../modules/set_up_external_integration_target/scripts/start-container.sh")]

  environment = {
    CONTAINER_IMAGE = local.image
    CONTAINER_NAME  = "${var.database_type}-${var.instance_name}"
    CONTAINER_PORTS = var.port
    CONTAINER_ENVS  = local.env_vars
    CONTAINER_ARGS  = local.args
  }

  transport = {
    ssh = {
      host = var.host.public_ip
    }
  }
}

# Outputs
output "config" {
  description = "Database configuration details"
  value = {
    type     = var.database_type
    username = var.username
    password = var.password
    database = var.database
    version  = var.db_version
    port     = var.port
    host     = var.host
  }
}
