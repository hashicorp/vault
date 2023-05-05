pid_file = "/tmp/pidfile"

vault {
  tls_skip_verify = true
  retry {
    num_retries = 10
  }
}

template {
  destination  = "/tmp/agent/render-content.txt"
  contents     = "{{ with secret \"auth/token/lookup-self\" }}orphan={{ .Data.orphan }} display_name={{ .Data.display_name }}{{ end }}"
}

auto_auth {
  method {
    type      = "azure"
    config = {
      role = "dev-role"
      resource = "https://management.azure.com/"
      authenticate_from_environment = true
    }
  }
}
