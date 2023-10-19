# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

auto_auth {

  method {
    type = "token_file"

    config {
      token_file_path = "/Users/avean/.vault-token"
    }
  }
}

env_template "MY_DATABASE_USER" {
  contents = "{{ with secret \"secret/db-secret\" }}{{ .Data.data.user }}{{ end }}"
}

exec {
  command = ["/path/to/my/app", "arg1", "arg2"]
}
