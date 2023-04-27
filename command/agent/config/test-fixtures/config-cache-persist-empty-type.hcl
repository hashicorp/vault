# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

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
