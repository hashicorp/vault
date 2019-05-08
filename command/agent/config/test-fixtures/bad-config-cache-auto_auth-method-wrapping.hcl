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
	}
}

cache {
	use_auto_auth_token = true
}

listener "tcp" {
    address = "127.0.0.1:8300"
    tls_disable = true
}


