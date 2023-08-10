# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

storage "raft" {
  path = "/path/to/raft"
  node_id = "raft_node_1"
}
listener "tcp" {
  address       = "127.0.0.1:8200"
  tls_cert_file = "/path/to/cert.pem"
  tls_key_file  = "/path/to/key.key"
}
seal "awskms" {
  kms_key_id       = "alias/kms-unseal-key"
}
service_registration "consul" {
  address       = "127.0.0.1:8500"
}
