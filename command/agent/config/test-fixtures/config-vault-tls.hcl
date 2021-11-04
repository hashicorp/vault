pid_file = "./pidfile"

vault {
  address = "http://127.0.0.1:1111"
  tls_min_version = "tls13"
  retry {
    num_retries = 5
  }
}
