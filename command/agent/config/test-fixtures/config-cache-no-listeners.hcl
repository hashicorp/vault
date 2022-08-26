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

template {
  source      = "/path/on/disk/to/template.ctmpl"
  destination = "/path/on/disk/where/template/will/render.txt"
}
