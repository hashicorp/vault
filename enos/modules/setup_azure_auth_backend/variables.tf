variable "app_name" {
  description = "The name of the application to setup the active directory for."
  type        = string
  default     = "vault"
}

variable "cluster_id" {
  description = "The unique id of the cluster, typically the 'context_name'."
}

variable "oidc_issuer_url" {
  description = "The OIDC issuer URL"
  type        = string
}

variable "service_accounts" {
  description = "The names of service accounts to create credentials for"
  type        = set(string)
}
