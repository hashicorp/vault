# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

pid_file = "./pidfile"

auto_auth {
  method {
    type      = "gcp"
    namespace = "/gcp-namespace"
    config = {
      role            = "dev-role"
      service_account = "https://[2001:DB8:AC3:FE4::1]"
    }
  }

  sink {
    type = "file"
    config = {
      path = "/tmp/file-foo"
    }
    aad     = "foobar"
    dh_type = "curve25519"
    dh_path = "/tmp/file-foo-dhpath"
  }
}

api_proxy {
  use_auto_auth_token = true
  enforce_consistency = "always"
  when_inconsistent   = "retry"
}

cache {
  persist = {
    type                       = "kubernetes"
    path                       = "/vault/agent-cache/"
    keep_after_import          = true
    exit_on_err                = true
    service_account_token_file = "/tmp/serviceaccount/token"
  }
}

listener "unix" {
  address      = "/path/to/socket"
  tls_disable  = true
  socket_mode  = "configmode"
  socket_user  = "configuser"
  socket_group = "configgroup"
}

listener "tcp" {
  address     = "2001:0db8::0001:8200"
  tls_disable = true
}

listener {
  type        = "tcp"
  address     = "[2001:0:0:1:0:0:0:1]:3000"
  tls_disable = true
  role        = "metrics_only"
}

listener "tcp" {
  role          = "default"
  address       = "2001:db8:0:1:1:1:1:1:8400"
  tls_key_file  = "/path/to/cakey.pem"
  tls_cert_file = "/path/to/cacert.pem"
}

vault {
  address         = "https://[2001:0db8::0001]:8200"
  ca_cert         = "config_ca_cert"
  ca_path         = "config_ca_path"
  tls_skip_verify = "true"
  client_cert     = "config_client_cert"
  client_key      = "config_client_key"
}
