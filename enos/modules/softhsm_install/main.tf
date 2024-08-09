# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "hosts" {
  type = map(object({
    ipv6       = string
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
  packages = var.include_tools ? {
    // These packages match the distros that are currently defined in the `ec2_info` module.
    amzn = {
      "2023" = ["softhsm", "opensc"]
    }
    rhel = {
      "8.10" = ["softhsm", "opensc"]
      "9.4"  = ["softhsm", "opensc"]
    }
    ubuntu = {
      "20.04" = ["softhsm", "opensc"]
      "22.04" = ["softhsm", "opensc"]
      "24.04" = ["softhsm2", "opensc"]
    }
    } : {
    amzn = {
      "2023" = ["softhsm"]
    }
    rhel = {
      "8.10" = ["softhsm"]
      "9.4"  = ["softhsm"]
    }
    ubuntu = {
      "20.04" = ["softhsm"]
      "22.04" = ["softhsm"]
      "24.04" = ["softhsm2"]
    }
  }
}

// Get the host information so we can ensure that we install the correct packages depending on the
// distro and distro version
resource "enos_host_info" "target" {
  transport = {
    ssh = {
      host = var.hosts["0"].public_ip
    }
  }
}

module "install_softhsm" {
  source = "../install_packages"

  hosts    = var.hosts
  packages = local.packages[enos_host_info.target.distro][enos_host_info.target.distro_version]
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
