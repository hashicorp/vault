pid_file = "./pidfile"

auto_auth {
	method "aws" {
		mount_path = "auth/aws"
		namespace = "my-namespace"
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
		derive_key = true
		config = {
			path = "/tmp/file-bar"
		}
	}
}
