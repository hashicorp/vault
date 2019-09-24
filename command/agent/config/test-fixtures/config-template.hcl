pid_file = "./pidfile"

auto_auth {
  method {
    type      = "aws"
    namespace = "/my-namespace"

    config = {
      role = "foobar"
    }
  }

  sink {
    type = "file"

    config = {
      path = "/tmp/file-foo"
    }

    aad     = "foobar"
    dh_type = "curve25519"
    dh_path = "/tmp/file-foo-dhpath"
  }
}

template {
  source           = "/path/on/disk/to/template.ctmpl"
  destination      = "/path/on/disk/where/template/will/render.txt"
  create_dest_dirs = true

  #contents = "{{ keyOrDefault \"service/redis/maxconns@east-aws\" \"5\" }}"
  command              = "restart service foo"
  command_timeout      = "60s"
  error_on_missing_key = false
  perms                = 0600
  backup               = true
  left_delimiter       = "{{"
  right_delimiter      = "}}"
  function_blacklist   = []
  sandbox_path         = ""

  wait {
    min = "2s"
    max = "10s"
  }
}
