# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

cache {
	use_auto_auth_token = true
	persist = {
		type = "kubernetes"
		path = "/vault/agent-cache/"
		keep_after_import = true
		exit_on_err = true
		service_account_token_file = "/tmp/serviceaccount/token"
	}
}

vault {
	address = "http://127.0.0.1:1111"
	ca_cert = "config_ca_cert"
	ca_path = "config_ca_path"
	tls_skip_verify = "true"
	client_cert = "config_client_cert"
	client_key = "config_client_key"
}
