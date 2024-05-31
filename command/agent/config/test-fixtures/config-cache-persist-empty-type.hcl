# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

pid_file = "./pidfile"

cache {
    persist = {
        path = "/vault/agent-cache/"
    }
}

listener "tcp" {
    address = "127.0.0.1:8300"
    tls_disable = true
}
