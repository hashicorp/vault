variable "vault_image_repository" {
  description = "The repository for the docker image to load, i.e. hashicorp/vault"
  type        = string
  default     = null
}

variable "vault_product_version" {
  description = "The vault product version to test"
  type        = string
  default     = null
}

variable "vault_product_revision" {
  type        = string
  description = "The vault product revision to test"
  default     = null
}

variable "vault_docker_image_archive" {
  description = "The path to the location of the docker image archive to test"
  type        = string
  default     = null
}

variable "vault_instance_count" {
  description = "How many instances to create for the Vault cluster"
  type        = number
  default     = 3
}

variable "terraform_plugin_cache_dir" {
  description = "The directory to cache Terraform modules and providers"
  type        = string
  default     = null
}

variable "tfc_api_token" {
  description = "The Terraform Cloud QTI Organization API token."
  type        = string
}

variable "vault_build_date" {
  description = "The build date for the vault docker image"
  type        = string
  default     = ""
}
