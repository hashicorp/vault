# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

variable "file_name" {}

output "license" {
  value = file(var.file_name)
}
