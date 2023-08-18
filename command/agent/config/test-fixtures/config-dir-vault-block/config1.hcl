# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

vault {
	address = "http://127.0.0.1:1111"
	ca_cert = "config_ca_cert"
	ca_path = "config_ca_path"
	tls_skip_verify = "true"
	client_cert = "config_client_cert"
	client_key = "config_client_key"
}
