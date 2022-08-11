
variable "enos_transport_user" {
  type        = string
  description = "The enos transport user"
  default     = null
}

variable "vault_instances" {
  type = map(object({
    public_ip   = string
    private_id  = string
    instance_id = string
  }))
  description = "The vault instances for the cluster to verify"
  default     = {}
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
  default     = null
}

variable "vault_token" {
  type        = string
  description = "The directory where the Vault binary will be installed"
  default     = null
}

variable "vault_autopilot_upgrade_version" {
  type        = string
  description = "The directory where the Vault binary will be installed"
  default     = null
}

variable "vault_autopilot_upgrade_status" {
  type        = string
  description = "The directory where the Vault binary will be installed"
  default     = null
}
