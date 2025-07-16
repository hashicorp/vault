# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

pid_file = "./pidfile"
log_file = "/var/log/vault/vault-agent.log"

vault {
  address = "https://[2001:0db8::0001]:8200"
}

auto_auth {
  method {
    type      = "aws"
    namespace = "/aws-namespace"
    config = {
      role = "foobar"
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

  sink {
    type        = "file"
    wrap_ttl    = "5m"
    aad_env_var = "TEST_AAD_ENV"
    dh_type     = "curve25519"
    dh_path     = "/tmp/file-foo-dhpath2"
    derive_key  = true
    config = {
      path = "/tmp/file-bar"
    }
  }
}

listener "unix" {
  address     = "/path/to/socket"
  tls_disable = true

  agent_api {
    enable_quit = true
  }
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
