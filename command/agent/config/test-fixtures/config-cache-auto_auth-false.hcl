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
	use_auto_auth_token = "false"
}

listener "tcp" {
    address = "127.0.0.1:8300"
    tls_disable = true
}

