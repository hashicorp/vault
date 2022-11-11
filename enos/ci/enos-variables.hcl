variable "terraform_plugin_cache_dir" {
  description = "The directory to cache Terraform modules and providers"
  type        = string
  default     = null
}

variable "tfc_api_token" {
  description = "The Terraform Cloud QTI Organization API token"
  type        = string
  sensitive   = true
}

variable "aws_ssh_public_key_path" {
  description = "The path to the public key to use for the ssh key"
  type        = string
}

variable "regions" {
  description = "The regions to bootstrap"
  type    = list(string)
  default = ["us-east-1", "us-east-2", "us-west-1", "us-west-2"]
}

variable "organization" {
  description = "The organization where the workspaces should be created"
  type        = string
  default     = "hashicorp-qti"
}
