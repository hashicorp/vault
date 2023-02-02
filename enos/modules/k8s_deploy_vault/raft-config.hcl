ui_config {
  enabled = true
}

listener "tcp" {
  address = "[::]:8200"
  cluster_address = "[::]:8201"
  tls_disable = true
}

storage "raft" {
  path = "/vault/data"
  autopilot {
    cleanup_dead_servers = "true"
    last_contact_threshold = "200ms"
    last_contact_failure_threshold = "10m"
    max_trailing_logs = 250000
    min_quorum = 5
    server_stabilization_time = "10s"
  }
}

service_registration "kubernetes" {}
