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
		wrap_ttl = 300
	}
}
