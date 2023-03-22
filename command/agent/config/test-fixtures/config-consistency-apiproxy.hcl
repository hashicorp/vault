# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

api_proxy {
	enforce_consistency = "always"
	when_inconsistent = "retry"
}

listener "tcp" {
	address = "127.0.0.1:8300"
	tls_disable = true
}
