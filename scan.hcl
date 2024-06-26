# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

repository {
  go_modules = true
  osv = true
  secrets {
    all = true
  } 
  dependabot {
    required = true
    check_config = true
  }
  
  plugin "semgrep" {
    use_git_ignore = true
    exclude = ["vendor"]
    config = [
      "tools/semgrep/ci",
      "p/r2c-security-audit",
      "r/trailofbits.go.hanging-goroutine.hanging-goroutine",
      "r/trailofbits.go.racy-append-to-slice.racy-append-to-slice",
      "r/trailofbits.go.racy-write-to-map.racy-write-to-map",
    ]
    exclude_rule = ["generic.html-templates.security.unquoted-attribute-var.unquoted-attribute-var"]
  }
  
  plugin "codeql" {
    languages = ["go"]
  }
}
