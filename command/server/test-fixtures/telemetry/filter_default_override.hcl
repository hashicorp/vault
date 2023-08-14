# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

disable_mlock = true
ui            = true

telemetry {
  statsd_address = "foo"
  filter_default = false
}