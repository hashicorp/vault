storage "raft" {
  path = "/path/to/raft/data"
  node_id = "raft_node_1"
}

api_addr = "http://127.0.0.1:8200"
cluster_addr = "https://127.0.0.1:8201"

ui_config {
  enabled = true
}
