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
        do_not_publish = true
		type = "file"
		config = {
			path = "/tmp/file-foo"
		}
		aad = "foobar"
		dh_type = "curve25519"
		dh_path = "/tmp/file-foo-dhpath"
	}

	sink {
		name = "nonauto"
		type = "file"
		wrap_ttl = "5m" 
		aad_env_var = "TEST_AAD_ENV"
		dh_type = "curve25519"
		dh_path = "/tmp/file-foo-dhpath2"
		config = {
			path = "/tmp/file-bar"
		}
	}

	sink {
        dh_auto = true
		name = "auto"
		type = "file"
		wrap_ttl = "5m"
		config = {
			path = "/tmp/file-baz"
		}
	}
}
