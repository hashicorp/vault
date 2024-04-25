# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

pid_file = "./pidfile"

cache {
    cache_static_secrets = false
    disable_caching_dynamic_secrets = true
}

listener "tcp" {
    address = "127.0.0.1:8300"
    tls_disable = true
}

vault {
	address = "http://127.0.0.1:1111"
	tls_skip_verify = "true"
}
