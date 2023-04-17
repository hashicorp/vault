auto_auth {

  method {
    type = "token_file"

    config {
      token_file_path = "$HOME/.vault-token"
    }
  }
}

vault {
  address = "http://localhost:8200"
}

env_template "FOO_DATA_LOCK" {
  contents             = "{{ with secret \"secret/data/foo\" }}{{ .Data.data.lock }}{{ end }}"
  error_on_missing_key = false
}
env_template "FOO_DATA_PASSWORD" {
  contents             = "{{ with secret \"secret/data/foo\" }}{{ .Data.data.password }}{{ end }}"
  error_on_missing_key = false
}
env_template "FOO_DATA_USER" {
  contents             = "{{ with secret \"secret/data/foo\" }}{{ .Data.data.user }}{{ end }}"
  error_on_missing_key = false
}

exec {
  command               = "env"
  args                  = []
  restart_on_new_secret = "always"
  #  restart_kill_signal   = "SIGTERM"
}