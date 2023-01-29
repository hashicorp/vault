pid_file = "./pidfile"

listener "tcp" {
    address = "127.0.0.1:8300"
    tls_disable = true
}