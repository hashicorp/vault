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

listener "tcp" {
    address = "127.0.0.1:8400"
    tls_key_file = "/path/to/cakey.pem"
    tls_cert_file = "/path/to/cacert.pem"
}

vault {
	address = "http://127.0.0.1:1111"
	ca_cert = "config_ca_cert"
	ca_path = "config_ca_path"
	tls_skip_verify = "true"
	client_cert = "config_client_cert"
	client_key = "config_client_key"
}
