# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

disable_cache = true
disable_mlock = true

backend "consul" {
    foo = "bar"
    disable_clustering = "true"
}

disable_clustering = false
