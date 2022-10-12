variable "vault_instance_count" {
  type        = number
  description = "How many vault instances are in the cluster"
}

variable "vault_pods" {
  type = list(object({
    name      = string
    namespace = string
  }))
  description = "The vault instances for the cluster to verify"
}

variable "vault_bin_path" {
  type        = string
  description = "The path to the vault binary"
  default     = "/bin/vault"
}

variable "vault_product_version" {
  type        = string
  description = "The vault product version"
}

variable "vault_product_revision" {
  type        = string
  description = "The vault product revision"
}

variable "vault_edition" {
  type        = string
  description = "The vault product edition"
}

variable "vault_root_token" {
  type        = string
  description = "The vault root token"
}

variable "kubeconfig_base64" {
  type        = string
  description = "The base64 encoded version of the Kubernetes configuration file"
}

variable "context_name" {
  type        = string
  description = "The name of the k8s context for Vault"
}

variable "check_build_date" {
  type        = bool
  description = "Whether or not to verify that the version includes the build date"
}

variable "vault_build_date" {
  type        = string
  description = "The build date of the vault docker image to check"
  default     = ""
}
