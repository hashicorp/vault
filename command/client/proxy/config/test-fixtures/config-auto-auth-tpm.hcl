# Copyright IBM Corp. 2016, 2026
# SPDX-License-Identifier: BUSL-1.1

pid_file = "./pidfile"

auto_auth {
  method "tpm" {
    mount_path = "auth/tpm"
    config = {
      tpm_cert = "machine"
    }
  }
}

api_proxy {
  use_auto_auth_token = "force"
}

vault {
  address   = "http://127.0.0.1:1111"
  namespace = "platform"
}

listener "tcp" {
  address     = "127.0.0.1:18100"
  tls_disable = true
}

tpm_cert "machine" {
  mount_path = "auth/tpm"
  role_name  = "web-tier"

  tpm {
    device_path = "/dev/tpmrm0"
  }

  destination {
    path  = "TMPDIR"
    umask = "077"
  }
}
