# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

output "k6_public_ip" {
  value = aws_instance.k6.public_ip
}

output "dashboard_url" {
  value = "http://${aws_instance.metrics.public_ip}:3000"
}
