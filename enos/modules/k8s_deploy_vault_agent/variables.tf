variable "client_id" {
  description = "The Azure activity directory application client id (application id)"
  type        = string
}

variable "service_account_name" {
  description = "The name of the service account to create for the agent pod"
  type        = string
}

variable "docker_image_name" {
  description = "The name of the docker image to run for the vault agent, e.g. hashicorp/vault"
  type        = string
  default     = "hashicorp/vault"
}

variable "docker_image_tag" {
  description = "The tag of the docker image to run for the vault agent, e.g. 1.13.1"
  type        = string
  default     = "1.13.1"
}
