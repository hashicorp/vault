# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

pid_file = "./pidfile"
disable_keep_alives = ["auto-auth"]

auto_auth {
  method {
    type      = "aws"
    namespace = "my-namespace/"

    config = {
      role = "foobar"
    }
  }

  sink {
    type = "file"
    config = {
      path = "/tmp/file-foo"
    }
    aad = "foobar"
    dh_type = "curve25519"
    dh_path = "/tmp/file-foo-dhpath"
  }
}

vault {
  address = "http://127.0.0.1:1111"
}
