# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

pid_file = "./pidfile"

auto_auth {
  method "aws" {
    mount_path = "auth/aws"
    config = {
      role = "foobar"
    }
  }
}
