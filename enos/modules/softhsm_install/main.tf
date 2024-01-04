# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

variable "hosts" {
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
  description = "The hosts that will have access to the softhsm. We assume they're all the same platform and architecture"
}

variable "include_tools" {
  type        = bool
  default     = false
  description = "Install opensc pkcs11-tools along with softhsm"
}

variable "retry_interval" {
  type        = string
  default     = "2"
  description = "How long to wait between retries"
}

variable "timeout" {
  type        = string
  default     = "15"
  description = "How many seconds to wait before timing out"
}

locals {
  packages = var.include_tools ? ["softhsm", "opensc"] : ["softhsm"]
}

module "install_softhsm" {
  source = "../install_packages"

  hosts    = var.hosts
  packages = local.packages
}

resource "enos_remote_exec" "find_shared_object" {
  for_each   = var.hosts
  depends_on = [module.install_softhsm]

  environment = {
    RETRY_INTERVAL  = var.retry_interval
    TIMEOUT_SECONDS = var.timeout
  }

  scripts = [abspath("${path.module}/scripts/find-shared-object.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

locals {
  object_paths = compact(distinct(values(enos_remote_exec.find_shared_object)[*].stdout))
}

output "lib" {
  value = local.object_paths[0]

  precondition {
    condition     = length(local.object_paths) == 1
    error_message = "SoftHSM targets cannot have different libsofthsm2.so shared object paths. Are they all the same Linux distro?"
  }
}
