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

        // https://github.com/advisories/GHSA-qqj6-54q6-cxv6
        // https://pkg.go.dev/vuln/GO-2026-6665
        // gosnowflake v1 does not currently have patch for the issue, nor is
        // clear if v1 is still actively maintained. We'll want to migrate to
        // v2 ASAP.
        "GHSA-qqj6-54q6-cxv6", "GO-2026-6665",

        // github.com/go-jose/go-jose/v3 is EOL and has several open advisories.
        // Temporarily suppress them until we've fully migrated to v4.
        "GHSA-mgcc-cww5-fp2q", "GHSA-vvm8-gf78-28g4", "GHSA-4w3g-jfj2-cw2q",
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

        // https://github.com/advisories/GHSA-qqj6-54q6-cxv6
        // https://pkg.go.dev/vuln/GO-2026-6665
        // gosnowflake v1 does not currently have patch for the issue, nor is
        // clear if v1 is still actively maintained. We'll want to migrate to
        // v2 ASAP.
        "GHSA-qqj6-54q6-cxv6", "GO-2026-6665",

        // github.com/go-jose/go-jose/v3 is EOL and has several open advisories.
        // Temporarily suppress them until we've fully migrated to v4.
        "GHSA-mgcc-cww5-fp2q", "GHSA-vvm8-gf78-28g4", "GHSA-4w3g-jfj2-cw2q",
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
