pid_file = "./pidfile"

auto_auth {
	method "aws" {
		mount_path = "auth/aws"
		wrap_ttl = 300
		config = {
			role = "foobar"
		}
	}
}
