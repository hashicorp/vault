pid_file = "./pidfile"

auto_auth {
	method {
		type = "aws"
		wrap_ttl = 300
		config = {
			role = "foobar"
		}
		max_backoff = "2m"
	}

	sink {
		type = "file"
		config = {
			path = "/tmp/file-foo"
		}
	}
}
