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
		config = {
			path = "/tmp/file-foo"
		}
		aad = "foobar"
		dh_type = "curve25519"
		dh_path = "/tmp/file-foo-dhpath"
	}

	sink "file" {
		wrap_ttl = "5m" 
		aad_env_var = "TEST_AAD_ENV"
		dh_type = "curve25519"
		dh_path = "/tmp/file-foo-dhpath2"
		config = {
			path = "/tmp/file-bar"
		}
	}
}
