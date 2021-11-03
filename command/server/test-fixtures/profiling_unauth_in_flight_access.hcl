storage "inmem" {}
listener "tcp" {
  address = "127.0.0.1:8200"
  tls_disable = true
  profiling {
     unauthenticated_in_flight_request_access = true
  }
}
disable_mlock = true
