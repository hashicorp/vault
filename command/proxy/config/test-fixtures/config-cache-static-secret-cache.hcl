# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

pid_file = "./pidfile"

auto_auth {
	method {
		type = "aws"
		config = {
			role = "foobar"
		}
	}

	sink {
		type = "file"
		config = {
			path = "/tmp/file-foo"
		}
		aad = "foobar"
		dh_type = "curve25519"
		dh_path = "/tmp/file-foo-dhpath"
	}
}

cache {
    cache_static_secrets = true
    static_secret_token_capability_refresh_interval = "1h"
}

listener "tcp" {
    address = "127.0.0.1:8300"
    tls_disable = true
}

vault {
	address = "http://127.0.0.1:1111"
	tls_skip_verify = "true"
}
