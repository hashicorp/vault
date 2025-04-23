# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

binary {
  secrets    = false
  go_modules = false
  osv        = true
  oss_index  = true
  nvd        = false
}

container {
  dependencies    = true
  alpine_security = true
  secrets         = true

  triage {
    suppress {
      vulnerabilities = [
        "CVE-2025-46394", // We can't do anything about this until a new Alpine container with busybox 1.38 is available.
        "GO-2022-0635",   // github.com/aws/aws-sdk-go@v1.x
      ]
    }
  }
}
