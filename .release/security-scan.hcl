# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

binary {
  secrets    = false
  go_modules = true
  osv        = true
}

container {
  dependencies    = true
  alpine_security = true
  secrets         = true

  triage {
    suppress {
      vulnerabilities = [
        // We can't do anything about these two CVE's until a new Alpine container with busybox 1.38 is available.
        "CVE-2025-46394",
        "CVE-2024-58251",
        "GO-2022-0635", // github.com/aws/aws-sdk-go@v1.x
      ]
    }
  }
}
