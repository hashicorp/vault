# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

vault {
  address = "http://127.0.0.1:1111"
  retry {
    num_retries = 5
  }
}

template_config {}

template {
  source      = "/path/on/disk/to/template.ctmpl"
  destination = "/path/on/disk/where/template/will/render.txt"
}