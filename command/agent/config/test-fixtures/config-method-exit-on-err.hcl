# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

pid_file = "./pidfile"

auto_auth {
	method {
		type = "aws"
		wrap_ttl = 300
        exit_on_err = true
		config = {
			role = "foobar"
		}
		max_backoff = "2m"
        min_backoff = "5s"
	}

	sink {
		type = "file"
		config = {
			path = "/tmp/file-foo"
		}
	}
}
