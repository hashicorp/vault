# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

pid_file = "./pidfile"

listener "tcp" {
    address = "127.0.0.1:8300"
    tls_disable = true
}