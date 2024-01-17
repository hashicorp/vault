# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

auto_auth {
  method {
    type = "token_file"
    config {
      token_file_path = "/home/username/.vault-token"
    }
  }
}

env_template "MY_PASSWORD" {
  source = "/path/on/disk/to/template.ctmpl"
}

exec {
  command = ["/path/to/my/app", "arg1", "arg2"]
}
