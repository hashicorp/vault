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
variable "common_tags" {}
variable "consul_license" {}
variable "consul_release" {}
variable "environment" {}
variable "instance_type" {}
variable "kms_key_arn" {}
variable "project_name" {}
variable "ssh_aws_keypair" {}
variable "vpc_id" {}
variable "common_tags" {}

output "consul_cluster_tag" {
  value = null
}
