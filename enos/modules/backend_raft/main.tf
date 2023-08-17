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

variable "ami_id" {
  default = null
}
variable "common_tags" {
  default = null
}
variable "consul_license" {
  default = null
}
variable "consul_release" {
  default = null
}
variable "environment" {
  default = null
}
variable "instance_type" {
  default = null
}
variable "kms_key_arn" {
  default = null
}
variable "project_name" {
  default = null
}
variable "ssh_aws_keypair" {
  default = null
}
variable "vpc_id" {
  default = null
}

output "consul_cluster_tag" {
  value = null
}
