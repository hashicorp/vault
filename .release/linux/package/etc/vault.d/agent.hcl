pid_file = "./pidfile"

vault {
  address = "https://vault-fqdn:8200"
  retry {
    num_retries = 5
  }
}

auto_auth {
  method "aws" {
    mount_path = "auth/aws-subaccount"
    config = {
      type = "iam"
      role = "foobar"
    }
  }

  sink "file" {
    config = {
      path = "/run/vault-agent/file-foo"
    }
  }

  sink "file" {
    wrap_ttl = "5m"
    aad_env_var = "TEST_AAD_ENV"
    dh_type = "curve25519"
    dh_path = "/run/vault-agent/file-foo-dhpath2"
    config = {
      path = "/run/vault-agent/file-bar"
    }
  }
}

cache {
  // An empty cache stanza still enables caching
}

api_proxy {
  use_auto_auth_token = true
}

listener "unix" {
  address = "/path/to/socket"
  tls_disable = true

  agent_api {
    enable_quit = true
  }
}

listener "tcp" {
  address = "127.0.0.1:8100"
  tls_disable = true
}

template {
  source = "/etc/vault/server.key.ctmpl"
  destination = "/run/vault-agent/server.key"
}

template {
  source = "/etc/vault/server.crt.ctmpl"
  destination = "/run/vault-agent/server.crt"
}

