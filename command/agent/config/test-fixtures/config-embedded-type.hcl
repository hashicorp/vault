pid_file = "./pidfile"

auto_auth {
	method "aws" {
		mount_path = "auth/aws"
		wrap_ttl = 300
		config = {
			role = "foobar"
		}
	}

	sink "file" {
		do_not_publish = true
		config = {
			path = "/tmp/file-foo"
		}
		aad = "foobar"
		dh_type = "curve25519"
		dh_path = "/tmp/file-foo-dhpath"
	}

	sink "file" {
		name = "nonauto"
		wrap_ttl = "5m"
		aad_env_var = "TEST_AAD_ENV"
		dh_type = "curve25519"
		dh_path = "/tmp/file-foo-dhpath2"
		config = {
			path = "/tmp/file-bar"
		}
	}

    sink "file" {
        auto = true
        name = "auto"
        type = "file"
        wrap_ttl = "5m"
    }
}
