schema_version = 1

project {
  license          = "BUSL-1.1"
  copyright_year   = 2025
  copyright_holder = "IBM Corp."

  # (OPTIONAL) A list of globs that should not have copyright/license headers.
  # Supports doublestar glob patterns for more flexibility in defining which
  # files or folders should be ignored
  header_ignore = [
    "helper/pkcs7/**",
    "ui/node_modules/**",
    "ui/pnpm-lock.yaml",
    "ui/pnpm-workspace.yaml",
    "enos/modules/k8s_deploy_vault/raft-config.hcl",
    "plugins/database/postgresql/scram/**",
    "enos/.enos/**",
    "enos/k8s/.enos/**",
    "enos/.terraform/**",
  ]
}
