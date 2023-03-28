# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

disable_mlock = true
ui            = true

telemetry {
  statsd_address = "foo"
  filter_default = false
}