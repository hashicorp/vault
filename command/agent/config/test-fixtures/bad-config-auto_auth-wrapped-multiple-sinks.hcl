# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

pid_file = "./pidfile"

auto_auth {
  method "aws" {
    mount_path = "auth/aws"
    wrap_ttl = 300
    config = {
      role = "foobar"
    }
  }

  sink "file" {
    config = {
      path = "/tmp/file-foo"
    }
  }

  sink "file" {
    config = {
      path = "/tmp/file-bar"
    }
  }
}
