# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

binary {
  go_stdlib  = true // Scan the Go standard library used to build the binary.
  go_modules = true // Scan the Go modules included in the binary.
  osv        = true // Use the OSV vulnerability database.
  oss_index  = true // And use OSS Index vulnerability database.

  triage {
    suppress {
      vulnerabilities = [
        "GO-2022-0635", // github.com/aws/aws-sdk-go@v1.x
      ]
    }
  }
}

container {
  dependencies = true // Scan any installed packages for vulnerabilities.
  osv          = true // Use the OSV vulnerability database.

  secrets {
    all = true
  }

  triage {
    suppress {
      vulnerabilities = [
        // We can't do anything about these two CVEs until a new Alpine container with busybox 1.38 is available.
        "CVE-2025-46394",
        "CVE-2024-58251",
        "GO-2022-0635", // github.com/aws/aws-sdk-go@v1.x
      ]
    }
  }
}
