schema_version = 1

project {
  license          = "BUSL-1.1"
  copyright_year   = 2026
  copyright_holder = "IBM Corp."

  # (OPTIONAL) A list of globs that should not have copyright/license headers.
  # Supports doublestar glob patterns for more flexibility in defining which
  # files or folders should be ignored
  header_ignore = [
    "enos/.enos/**",
    "enos/.terraform/**",
    "enos/k8s/.enos/**",
    "enos/modules/k8s_deploy_vault/raft-config.hcl",
    "helper/pkcs7/**",
    "plugins/database/postgresql/scram/**",
    "tools/pipeline/internal/pkg/generate/fixtures/*",
    "ui/node_modules/**",
    "ui/pnpm-lock.yaml",
    "ui/pnpm-workspace.yaml",
  ]
}
