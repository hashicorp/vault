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
        // We don't actually use github.com/jackc/pgproto3/v2@v2.3.3 anywhere
        // because we've upgraded to github.com/jackc/pgx/v5 eveywhere. This is
        // only included in the go.sum because cloud.google.com/go/cloudsqlconn/postgres/pgxv4",
        // which we don't use.
        "GO-2026-4518",
        "GHSA-jqcq-xjh3-6g23",
        // This has been fixed upstream but is unlikely to be released until
        // bbolt/v5. Regardless, it's a condition that cannot happen in the real
        // world so we don't really need to worry about it. The reporter agrees
        // that it ought to be rescinded:
        // https://github.com/golang/vulndb/issues/4923#issuecomment-4200230182
        "GO-2026-4923",
        "GHSA-6jwv-w5xf-7j27",
        // These two are false positives that have been fixed as of
        // github.com/jackc/pgx/v5 @ >= v5.0.9
        // I've (ryan) opened requests to have the vulndb updated as such:
        // https://github.com/golang/vulndb/issues/4943
        // https://github.com/golang/vulndb/issues/4944
        "GO-2026-4772",
        "GO-2026-4771",
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
        // We don't actually use github.com/jackc/pgproto3/v2@v2.3.3 anywhere
        // because we've upgraded to github.com/jackc/pgx/v5 eveywhere. This is
        // only included in the go.sum because cloud.google.com/go/cloudsqlconn/postgres/pgxv4",
        // which we don't use.
        "GO-2026-4518",
        "GHSA-jqcq-xjh3-6g23",
        // This has been fixed upstream but is unlikely to be released until
        // bbolt/v5. Regardless, it's a condition that cannot happen in the real
        // world so we don't really need to worry about it. The reporter agrees
        // that it ought to be rescinded:
        // https://github.com/golang/vulndb/issues/4923#issuecomment-4200230182
        "GO-2026-4923",
        "GHSA-6jwv-w5xf-7j27",
        // These two are false positives that have been fixed as of
        // github.com/jackc/pgx/v5 @ >= v5.0.9
        // I've (ryan) opened requests to have the vulndb updated as such:
        // https://github.com/golang/vulndb/issues/4943
        // https://github.com/golang/vulndb/issues/4944
        "GO-2026-4772",
        "GO-2026-4771",
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
