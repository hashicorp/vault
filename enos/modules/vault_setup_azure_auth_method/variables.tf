variable "transport" {
  description = "The transport configuration to use when setting up the auth method. Must include the name and config, i.e. ssh = { host = ??, user = ?? }"
  type        = any # Cannot use object (even with optional properties) since it blows up.
}

variable "vault_root_token" {
  description = "The Vault root token"
  type        = string
}

variable "client_id" {
  description = "The Azure activity directory application client id (application id)"
  type        = string
}

variable "client_secret" {
  description = "The Azure activity directory application client secret"
  type        = string
}
