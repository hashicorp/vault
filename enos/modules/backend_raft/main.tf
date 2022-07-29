// Shim module to handle the fact that Vault doesn't actually need a backend module
terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

variable "ami_id" {}
variable "consul_release" {}
variable "instance_type" {}
variable "kms_key_arn" {}
variable "vpc_id" {}
variable "common_tags" {}

output "consul_cluster_tag" {
  value = null
}
