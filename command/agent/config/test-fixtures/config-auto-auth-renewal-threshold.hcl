# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

vault {
  address = "http://127.0.0.1:8200"
}

auto_auth {
  method {
    type      = "approle"
    lease_renewal_threshold = 0.75
    config = {
      role_id_file_path = "/tmp/role-id"
      secret_id_file_path = "/tmp/secret-id"
    }
  }

  sink {
    type = "file"
    config = {
      path = "/tmp/token"
    }
  }
}
