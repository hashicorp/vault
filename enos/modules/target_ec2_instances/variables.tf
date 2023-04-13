variable "ami_id" {
  description = "The machine image identifier"
  type        = string
}

variable "awskms_unseal_key_arn" {
  type        = string
  description = "The AWSKMS key ARN if using the awskms unseal method. If specified the instances will be granted kms permissions to the key"
  default     = null
}

variable "cluster_name" {
  type        = string
  description = "A unique cluster identifier"
  default     = null
}

variable "common_tags" {
  description = "Common tags for cloud resources"
  type        = map(string)
  default     = { "Project" : "Enos" }
}

variable "instance_count" {
  description = "The number of target instances to create"
  type        = number
  default     = 3
}

variable "instance_type" {
  description = "The instance machine type"
  type        = string
  default     = "t3.small"
}

variable "project_name" {
  description = "A unique project name"
  type        = string
}

variable "spot_price_max" {
  description = "Unused shim variable to match target_ec2_spot_fleet"
  type        = string
  default     = null
}

variable "ssh_allow_ips" {
  description = "Allowlisted IP addresses for SSH access to target nodes. The IP address of the machine running Enos will automatically allowlisted"
  type        = list(string)
  default     = []
}

variable "ssh_keypair" {
  description = "SSH keypair used to connect to EC2 instances"
  type        = string
}

variable "vpc_id" {
  description = "The identifier of the VPC where the target instances will be created"
  type        = string
}
