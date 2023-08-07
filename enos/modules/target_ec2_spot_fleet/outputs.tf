# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

output "cluster_name" {
  value = local.cluster_name
}

output "hosts" {
  description = "The ec2 fleet target hosts"
  value = { for idx in range(var.instance_count) : idx => {
    public_ip  = data.aws_instance.targets[idx].public_ip
    private_ip = data.aws_instance.targets[idx].private_ip
  } }
}
