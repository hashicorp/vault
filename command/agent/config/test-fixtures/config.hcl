pid_file = "./pidfile"

auto_auth {
	method {
		type = "aws-iam"
		mount_path = "auth/aws"
		config = {
			role = "foobar"
		}
	}

	vault {
		address = "http://127.0.0.1:8200"
		tls_skip_verify = true
	}

	token_sink "file" {
		path = "/tmp/file-foo"
	}

	token_sink "file" {
		path = "/tmp/file-bar"
	}
}
