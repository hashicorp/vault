
variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
  default     = null
}

variable "vault_instance_count" {
  type        = number
  description = "How many vault instances are in the cluster"
}

variable "vault_instances" {
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
  description = "The vault cluster instances that were created"
}

variable "vault_root_token" {
  type        = string
  description = "The vault root token"
  default     = null
}
