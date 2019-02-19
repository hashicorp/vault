pid_file = "./pidfile"

auto_auth {
	method {
		type = "aws"
		wrap_ttl = 300
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

	listener "unix" {
		address = "/path/to/socket"
		tls_disable = true
	}

	listener "tcp" {
		address = "127.0.0.1:8300"
		tls_disable = true
	}

	listener "tcp" {
		address = "127.0.0.1:8400"
		tls_key_file = "/path/to/cakey.pem"
		tls_cert_file = "/path/to/cacert.pem"
	}
}
