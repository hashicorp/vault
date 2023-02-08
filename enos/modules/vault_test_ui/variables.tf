variable "vault_addr" {
  description = "The host address for the vault instance to test"
  type        = string
}

variable "vault_root_token" {
  description = "The vault root token"
  type        = string
}

variable "ui_test_filter" {
  type        = string
  description = "A test filter to limit the ui tests to execute. Will be appended to the ember test command as '-f=<filter>'"
  default     = null
}

variable "vault_unseal_keys" {
  description = "Base64 encoded recovery keys to use for the seal/unseal test"
  type        = list(string)
}

variable "vault_recovery_threshold" {
  description = "The number of recovery keys to require when unsealing Vault"
  type        = string
}

variable "ui_run_tests" {
  type        = bool
  description = "Whether to run the UI tests or not. If set to false a cluster will be created but no tests will be run"
  default     = true
}
