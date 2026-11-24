# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

pid_file = "./pidfile"

auto_auth {
  method "aws" {
    mount_path = "auth/aws"
    config = {
      role = "foobar"
    }
  }
}
