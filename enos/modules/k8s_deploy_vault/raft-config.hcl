ui = true
listener "tcp" {
  address = "[::]:8200"
  cluster_address = "[::]:8201"
  tls_disable = true
}

storage "raft" {
  path = "/vault/data"
}

service_registration "kubernetes" {}
