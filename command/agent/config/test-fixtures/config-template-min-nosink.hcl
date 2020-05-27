pid_file = "./pidfile"

auto_auth {
  method {
    type      = "aws"
    namespace = "/my-namespace"

    config = {
      role = "foobar"
    }
  }
}

template {
  source      = "/path/on/disk/to/template.ctmpl"
  destination = "/path/on/disk/where/template/will/render.txt"
}
