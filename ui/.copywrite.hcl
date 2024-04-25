# (OPTIONAL) Overrides the copywrite config schema version
# Default: 1
schema_version = 1

project {
  # (OPTIONAL) SPDX-compatible license identifier
  # Leave blank if you don't wish to license the project
  # Default: "MPL-2.0"
  license = "BUSL-1.1"

  # (OPTIONAL) Represents the copyright holder used in all statements
  # Default: HashiCorp, Inc.
  # copyright_holder = ""

  # (OPTIONAL) Represents the year that the project initially began
  # Default: <the year the repo was first created>
  # copyright_year = 0

  # (OPTIONAL) A list of globs that should not have copyright or license headers .
  # Supports doublestar glob patterns for more flexibility in defining which
  # files or folders should be ignored
  # Default: []
  header_ignore = [
    "node_modules/**"
  ]

  # (OPTIONAL) Links to an upstream repo for determining repo relationships
  # This is for special cases and should not normally be set.
  # Default: ""
  # upstream = "hashicorp/<REPONAME>"
}


