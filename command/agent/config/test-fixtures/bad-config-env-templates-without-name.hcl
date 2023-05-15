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

# Error: env_template must have an environment variable name specified in the header
env_template {
  contents = "{{ with secret \"secret/data/foo\" }}{{ .Data.data.password }}{{ end }}"
}

exec {
  command               = ["./my-app", "arg1", "arg2"]
  restart_on_new_secret = "always"
  restart_kill_signal   = "SIGTERM"
}
