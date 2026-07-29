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
        // GO-2022-0635 is valid. We will remove this when the ongoing migration
        // to github.com/aws/aws-sdk-go/v2 has been completed.
        "GO-2022-0635",
        // GO-2026-5932 appears to be a false positive as it only relates to
        // golang.org/x/crypto/openpgp, which is not in use in the project.
        // https://pkg.go.dev/vuln/GO-2026-5932
        "GO-2026-5932",
        // GO-2026-5298 appears to be a false positive as the associated
        // Github Security Advisory shows that v0.6.1 contains the fix.
        // The issue has been reported but the advisory has not been updated
        // yet.
        // https://pkg.go.dev/vuln/GO-2026-5298
        // https://github.com/google/go-attestation/security/advisories/GHSA-9r4w-jg96-92mv
        // https://github.com/golang/vulndb/issues/5795
        "GO-2026-5298",
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
        // GO-2022-0635 is valid. We will remove this when the ongoing migration
        // to github.com/aws/aws-sdk-go/v2 has been completed.
        "GO-2022-0635",
        // GO-2026-5932 appears to be a false positive as it only relates to
        // golang.org/x/crypto/openpgp, which is not in use in the project.
        // https://pkg.go.dev/vuln/GO-2026-5932
        "GO-2026-5932",
        // GO-2026-5298 appears to be a false positive as the associated
        // Github Security Advisory shows that v0.6.1 contains the fix.
        // The issue has been reported but the advisory has not been updated
        // yet.
        // https://pkg.go.dev/vuln/GO-2026-5298
        // https://github.com/google/go-attestation/security/advisories/GHSA-9r4w-jg96-92mv
        // https://github.com/golang/vulndb/issues/5795
        "GO-2026-5298",
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
