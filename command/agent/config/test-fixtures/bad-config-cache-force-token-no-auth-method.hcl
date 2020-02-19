pid_file = "./pidfile"

cache {
  use_auto_auth_token = "force"
}

listener "tcp" {
  address     = "127.0.0.1:8300"
  tls_disable = true
}
