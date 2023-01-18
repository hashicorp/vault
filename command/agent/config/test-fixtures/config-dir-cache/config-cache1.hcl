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

listener "unix" {
    address = "/path/to/socket"
    tls_disable = true
    socket_mode = "configmode"
    socket_user = "configuser"
    socket_group = "configgroup"
}

listener "tcp" {
    address = "127.0.0.1:8300"
    tls_disable = true
}

listener {
    type = "tcp"
    address = "127.0.0.1:3000"
    tls_disable = true
    role = "metrics_only"
}

listener "tcp" {
    role = "default"
    address = "127.0.0.1:8400"
    tls_key_file = "/path/to/cakey.pem"
    tls_cert_file = "/path/to/cacert.pem"
}