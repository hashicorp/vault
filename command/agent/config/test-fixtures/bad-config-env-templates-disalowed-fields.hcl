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

  # Error: source and destination are not allowed in env_template
  source      = "/path/on/disk/to/template.ctmpl"
  destination = "/path/on/disk/where/template/will/render.txt"
}

exec {
  command               = ["./my-app", "arg1", "arg2"]
  restart_on_new_secret = "always"
  restart_kill_signal   = "SIGTERM"
}
