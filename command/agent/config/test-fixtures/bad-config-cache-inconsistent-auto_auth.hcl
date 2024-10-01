# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

pid_file = "./pidfile"

cache {
	use_auto_auth_token = true
}

listener "tcp" {
    address = "127.0.0.1:8300"
    tls_disable = true
}


