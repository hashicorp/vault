# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

output "ids" {
  description = "The two security group IDs for grafana and prometheus"
  value = {
    "grafana"    = aws_security_group.grafana.id
    "prometheus" = aws_security_group.prometheus.id
  }
}