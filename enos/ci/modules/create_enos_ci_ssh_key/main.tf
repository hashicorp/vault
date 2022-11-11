
terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }
}

variable "region" {
  type        = string
  description = "The region to create the ssh key pair in"
}

variable "public_key" {
  type        = string
  description = "The public key to use for the ssh key pair"
}

resource "aws_key_pair" "enos_ci_key" {
  key_name   = "enos-ci-ssh-key"
  public_key = var.public_key
}

output "key_pair_id" {
  value = aws_key_pair.enos_ci_key.key_pair_id
}

output "key_pair_arn" {
  value = aws_key_pair.enos_ci_key.arn
}
