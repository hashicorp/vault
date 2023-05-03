variable "transport" {
  type = any
}

variable "vault_root_token" {
  type = string
}

variable "client_id" {
  description = "The Azure activity directory application client id (application id)"
  type        = string
}

variable "client_secret" {
  description = "The Azure activity directory application client secret"
  type        = string
}
