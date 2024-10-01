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
  contents             = "{{ with secret \"secret/data/foo\" }}{{ .Data.data.password }}{{ end }}"
  error_on_missing_key = false
}
env_template "FOO_USER" {
  contents             = "{{ with secret \"secret/data/foo\" }}{{ .Data.data.user }}{{ end }}"
  error_on_missing_key = false
}

exec {
  command                   = ["./my-app", "arg1", "arg2"]
  restart_on_secret_changes = "always"
  restart_stop_signal       = "SIGTERM"
}

# Error: api_proxy is incompatible with env_template
api_proxy {
	use_auto_auth_token = "force"
	enforce_consistency = "always"
	when_inconsistent   = "forward"
}

# Error: listener is incompatible with env_template
listener "tcp" {
    address = "127.0.0.1:8300"
    tls_disable = true
}
