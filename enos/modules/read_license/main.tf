# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "file_name" {}

output "license" {
  value = file(var.file_name)
}
