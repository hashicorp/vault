# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

storage "inmem" {}
listener "tcp" {
  address = "127.0.0.1:8200"
  tls_disable = true
  inflight_requests_logging {
     unauthenticated_in_flight_requests_access = true
  }
}
disable_mlock = true
