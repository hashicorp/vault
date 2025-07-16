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

template_config {
  static_secret_render_interval = "5m"
  exit_on_retry_failure         = true
}

vault {
  address = "http://localhost:8200"
}

env_template "FOO_PASSWORD" {
  contents    = "{{ with secret \"secret/data/foo\" }}{{ .Data.data.password }}{{ end }}"

  # Error: destination and create_dest_dirs are not allowed in env_template
  destination      = "/path/on/disk/where/template/will/render.txt"
  create_dest_dirs = true
}

exec {
  command                   = ["./my-app", "arg1", "arg2"]
  restart_on_secret_changes = "always"
  restart_stop_signal       = "SIGTERM"
}
