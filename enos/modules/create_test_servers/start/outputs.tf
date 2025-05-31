# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

output "instance_id" {
  description = "ID of the EC2 instance"
  value       = aws_instance.enos_test_server.id
}

output "instance_public_ip" {
  description = "Public IP address of the EC2 instance"
  value       = aws_instance.enos_test_server.public_ip
}

output "instance_private_ip" {
  description = "Private IP address of the EC2 instance"
  value       = aws_instance.enos_test_server.private_ip
}

output "instance_ami_id" {
  description = "AMI ID used for the EC2 instance"
  value       = data.aws_ami.amazon_linux.id
}

output "instance_type" {
  description = "Instance type used for the EC2 instance"
  value       = local.instance_type
}

output "aws_region" {
  description = "AWS region where resources are deployed"
  value       = data.aws_region.current.name
}

output "subnet_id" {
  description = "Subnet ID where the EC2 instance is launched"
  value       = aws_instance.enos_test_server.subnet_id
}

output "availability_zone" {
  description = "Availability Zone of the EC2 instance"
  value       = aws_instance.enos_test_server.availability_zone
}
