# Copyright IBM Corp. 2016, 2025
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
        "GO-2022-0635", // github.com/aws/aws-sdk-go@v1.x
      ]

      // The OSV scanner will trip on several packages that are included in the
      // the UBI images. This is due to RHEL using the same base version in the
      // package name for the life of the distro regardless of whether or not
      // that version has been patched for security. Rather than enumate ever
      // single CVE that the OSV scanner will find (several tens) we'll ignore
      // the base UBI packages.
      paths = [
        "usr/lib/sysimage/rpm/*",
        "var/lib/rpm/*",
      ]
    }
  }
}
