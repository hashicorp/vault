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
  source      = "/path/on/disk/to/template.ctmpl"
  destination = "/path/on/disk/where/template/will/render.txt"

  create_dest_dirs = true

  command = "restart service foo"

  error_on_missing_key = false
  perms                = 0600
}

template {
  source      = "/path/on/disk/to/template2.ctmpl"
  destination = "/path/on/disk/where/template/will/render2.txt"

  perms = 0755

  backup = true

  wait {
    min = "2s"
    max = "10s"
  }
}
