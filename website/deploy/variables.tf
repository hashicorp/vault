variable "name" {
  default     = "vault-www"
  description = "Name of the website in slug format."
}

variable "github_repo" {
  default     = "hashicorp/vault"
  description = "GitHub repository of the provider in 'org/name' format."
}

variable "github_branch" {
  default     = "stable-website"
  description = "GitHub branch which netlify will continuously deploy."
}

variable "custom_site_domain" {
  default     = "www.vaultproject.io"
  description = "The custom domain to use for the Netlify site."
}
