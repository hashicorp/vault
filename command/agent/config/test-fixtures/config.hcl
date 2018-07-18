pid_file = "./pidfile"

auto_auth {
	method {
		type = "aws-iam"
		mount_path = "auth/aws"
		config = {
			role = "foobar"
		}
	}

	sink "file" {
		path = "/tmp/file-foo"
	}

	sink "file" {
		wrap_ttl = "5m" 
		path = "/tmp/file-bar"
	}
}
