pid_file = "./pidfile"

auto_auth {
	method {
		type = "aws"
		namespace = "/my-namespace"
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

	sink {
		type = "file"
		wrap_ttl = "5m" 
		aad_env_var = "TEST_AAD_ENV"
		dh_type = "curve25519"
		dh_path = "/tmp/file-foo-dhpath2"
		derive_key = true
		config = {
			path = "/tmp/file-bar"
		}
	}
}
