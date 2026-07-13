# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

// Configure our changed-files detection for the pipeline utility. Here we define
// how we wish to group the files in the repository. We use this group metadata
// in several places in the pipeline, notably when synchronizing files between
// enterprise and CE, when creating CE backports that exclude enterprise only
// files, and for flow control on our CI workflows.
//
// If you're modifying this file you can verify file groupings using several
// different pipeline commands. You can use `pipeline git list changed-files`
// with a number of different methods. To see a commit, use the `--commit <SHA>`
// flag. To look at a range of commits, use the `--range HEAD~7..HEAD` flag,
// using any range that is currently present. If you'd like to see all of the
// files that have ever existed in the history of a branch, you can use the
// `--branch <branch>` flag.
changed_files {
  // ----------------------------------------------------------------------------
  // Automation
  // ----------------------------------------------------------------------------

  // The "automation_ce" group tracks CE (Community Edition) automation workflows.
  group "automation_ce" {
    match {
      base_dir = [joinpath(".github", "workflows")]
      base_name = [
        "add-hashicorp-contributed-label.yml",
        "backport-automation-ent.yml",
        "copy-external-contributor-pull-request-ce.yml",
        "copy-external-contributor-pull-request-ent.yml",
        "oss.yml",
        "remove-labels.yml",
        "sync-ce-branches-from-ent.yml",
      ]
    }
  }

  // The "automation_github" graup tracks all Github actions automation.
  group "automation_github" {
    match {
      base_dir = [
        joinpath(".github", "workflows"),
        joinpath(".github", "actions"),
        joinpath(".github", "scripts"),
      ]
    }
  }

  // The "automation_pr" group tracks pull request validation and automation
  // workflows.
  group "automation_pr" {
    match {
      base_dir = [joinpath(".github", "workflows")]
      base_name = [
        "approval-gate.yml",
        "changelog-checker.yml",
        "milestone-checker.yml",
      ]
    }
  }

  // The "automation_plugins" group tracks plugin update automation.
  group "automation_plugins" {
    match {
      base_dir = [joinpath(".github", "workflows")]
      base_name = [
        "plugin-update-check.yml",
        "plugin-update.yml",
      ]
    }
  }

  // The "automation_ui" group tracks UI-related automation workflows.
  group "automation_ui" {
    match {
      base_dir = [joinpath(".github", "workflows")]
      base_name = [
        "gen-diff-spec.yml",
        "ui-client-update.yml",
      ]
    }
  }

  // ----------------------------------------------------------------------------
  // Build
  // ----------------------------------------------------------------------------

  // The "build_vault" group tracks Vault binary build workflows and configuration.
  group "build_vault" {
    match {
      base_dir = [joinpath(".github", "workflows")]
      base_name = [
        "build.yml",
        "build-artifacts-ce.yml",
        "build-artifacts-ent.yml",
      ]
    }

    match {
      base_dir = [
        ".build",
        joinpath(".github", "actions", "containerize"),
        joinpath(".github", "actions", "build-vault"),
        joinpath(".github", "actions", "set-up-go"),
      ]
    }

    match {
      base_name = [
        "ci-helper.sh",
        "Dockerfile",
        "Dockerfile-ent",
        "Dockerfile-ent-hsm",
        "Makefile",
      ]
    }
  }

  // The "build_cloud" group tracks HCP/Cloud image building and testing.
  group "build_cloud" {
    match {
      file = [
        joinpath(".github", "workflows", "build-hcp-image.yml"),
        joinpath(".github", "workflows", "test-run-enos-scenario-cloud.yml"),
        joinpath("enos", "enos-scenario-cloud-ent.hcl"),
      ]
    }

    match {
      base_dir = [
        joinpath("enos", "modules", "cloud_docker_vault_cluster"),
        joinpath("enos", "modules", "hcp"),
        joinpath("tools", "pipeline", "internal", "pkg", "hcp"),
      ]
    }

    match {
      base_dir = [
        joinpath("tools", "pipeline", "internal", "cmd"),
      ]
      contains = ["hcp"]
    }
  }

  // The "build_tfe" group tracks Terraform Enterprise S390x build workflows.
  group "build_tfe" {
    match {
      base_dir  = [joinpath(".github", "workflows")]
      base_name = ["build-tfe-s390x-artifact-ent.yml"]
    }
  }

  group "build_ui" {
    match {
      base_dir = [
        joinpath(".github", "actions", "build-ui"),
        joinpath(".github", "actions", "set-up-pnpm"),
        joinpath("ui", "package.json"),
      ]
    }

    match {
      base_name = [
        "ci-helper.sh",
        "Makefile",
      ]
    }
  }

  // ----------------------------------------------------------------------------
  // Configuration
  // ----------------------------------------------------------------------------

  // The "config" group tracks build and pipeline configuration files.
  group "config" {
    match {
      base_name = [
        "Makefile",
        "CODEOWNERS",
      ]
    }

    match {
      base_dir = [".release"]
    }
  }

  // ----------------------------------------------------------------------------
  // Documentation
  // ----------------------------------------------------------------------------

  // The "changelog" group tracks our changelogs
  group "changelog" {
    match {
      base_dir = ["changelog"]
    }

    match {
      contains = ["CHANGELOG"]
    }
  }

  // The "docs" group tracks our README and static website documentation
  group "docs" {
    match {
      base_name = ["README.md"]
    }

    match {
      base_dir = ["website"]
    }
  }

  // ----------------------------------------------------------------------------
  // Editions
  // ----------------------------------------------------------------------------

  // The "community" group tracks files that exist specifically for either CE
  // builds or to document CE behavior.
  group "community" {
    match {
      // Match any Go file that has our oss or ce suffixes
      extension = [".go"]
      contains = [
        "_oss.go",
        "_ce.go",
      ]
    }

    match {
      // Match any configuration file, document, or script that include ce or
      // oss.
      extension = [
        ".hcl",
        ".md",
        ".sh",
        ".yaml",
        ".yml",
      ]
      contains = [
        "-ce",
        "_ce",
        "-oss",
        "_oss",
      ]
    }
  }

  // The "enterprise" group tracks our files that ought to exist in only in the
  // enterprise endition. It is arguably the most important grouping
  // configuration in this repository. Among other things, this is how we
  // enforce that no enterprise files make it into CE branches or repositories
  // with out automated tooling.
  group "enterprise" {
    // Ignore some files that would otherwise match our filters but exist in
    // some fashion on enterprise branches for reasons we cannot control.
    ignore {
      file = [
        # These exist on CE branches to please Github Actions.
        joinpath(".github", "workflows", "build-artifacts-ent.yml"),
        joinpath(".github", "workflows", "backport-automation-ent.yml"),
        joinpath(".github", "workflows", "test-run-enos-scenario-cloud.yml"),
      ]
    }

    // Ignore whole directories where no enterprise only code should ever exist.
    ignore {
      base_dir = [
        "website",
        joinpath(".release", "docker"),
        joinpath("scripts", "docker"),
      ]
    }

    // Ignore files that we expect to be CE
    ignore {
      extension = [
        ".proto",
        ".hcl",
        ".md",
        ".sh",
        ".yaml",
        ".yml",
      ]
      contains = [
        "-ce",
      ]
    }

    // Ignore some enos modules related to HSM
    ignore {
      base_dir = [
        # The next matcher looks for HSM but we can ignore the softhsm modules
        joinpath("enos", "modules", "softhsm_install"),
        joinpath("enos", "modules", "softhsm_create_vault_keys"),
        joinpath("enos", "modules", "softhsm_init"),
        joinpath("enos", "modules", "softhsm_distribute_vault_keys"),
        # Some filename have ent in them
        joinpath("enos", "modules", "verify_secrets_engines"),
      ]
    }

    // Make sure our zap scanner is always ent only
    match {
      base_dir = [
        joinpath("enos", "modules", "zap_scan_ent")
      ]
    }

    // Match whole directories that are enterprise only
    match {
      base_dir = [
        "vault_ent",
        joinpath("scripts", "dev", "hsm"),
        joinpath("scripts", "testing"),
        joinpath("specs"),
        joinpath(".release", "ibm-pao"),
        // Internal developer tooling that must not sync to CE
        ".agents",
        joinpath("ui", ".agents"),
      ]
    }

    match {
      base_name = [
        "Dockerfile-ent",
        "Dockerfile-ent-hsm",
      ]
    }

    match {
      extension = [".go"]
      contains = [
        "_ent.go",
        "_ent_test.go",
      ]
    }

    match {
      extension = [".pb.go"]
      contains = [
        "_ent",
      ]
    }

    match {
      extension = [
        ".proto",
        ".hcl",
        ".md",
        ".sh",
        ".yaml",
        ".yml",
      ]
      contains = [
        "-ent",
        "_ent",
        "hsm",
        "merkle-tree",
      ]
    }
  }

  // ----------------------------------------------------------------------------
  // Go
  // ----------------------------------------------------------------------------

  // The "go_tests" group isolates Go test file changes from application code.
  // This allows workflows to skip building when only tests change.
  group "go_tests" {
    ignore {
      base_dir = [
        joinpath("tools", "pipeline"),
      ]
    }

    match {
      extension = [".go"]
      contains  = ["_test.go"]
    }
  }

  // The "go_modules" group tracks Go module dependency changes.
  // Triggers dependency security scans and full test suite.
  group "go_modules" {
    ignore {
      base_dir = [
        joinpath("tools", "pipeline"),
      ]
    }

    match {
      base_name = [
        "go.mod",
        "go.sum",
      ]
    }
  }

  // The "go_toolchain" group tracks Go compiler/toolchain version changes.
  // Triggers full rebuild with new toolchain.
  group "go_toolchain" {
    match {
      base_name = [".go-version"]
    }
  }

  // The "go_app_core" group tracks core Vault application code (excluding tests).
  // Triggers full build, test cycle, and integration tests.
  group "go_app_core" {
    ignore {
      base_dir = [
        joinpath("tools", "pipeline"),
      ]
    }

    ignore {
      extension = [".go"]
      contains  = ["_test.go"]
    }

    match {
      extension = [".go"]
    }
  }

  // ----------------------------------------------------------------------------
  // Helpers
  // ----------------------------------------------------------------------------

  // The "actions_helpers" group tracks reusable setup/helper actions.
  group "actions_helpers" {
    match {
      base_dir = [
        joinpath(".github", "actions", "changed-files"),
        joinpath(".github", "actions", "checkout"),
        joinpath(".github", "actions", "create-dynamic-config"),
        joinpath(".github", "actions", "install-tools"),
        joinpath(".github", "actions", "metadata"),
        joinpath(".github", "actions", "set-up-go"),
        joinpath(".github", "actions", "set-up-pipeline"),
      ]
    }
  }

  // The "pipeline" group tracks the pipeline utility.
  group "pipeline" {
    match {
      base_dir = [
        joinpath("tools", "pipeline"),
      ]
    }
  }

  // The "proto" tracks our protobuf files and configuration.
  group "proto" {
    match {
      extension = [".proto"]
    }

    match {
      base_name_prefix = ["buf."]
    }
  }

  // ----------------------------------------------------------------------------
  // Infrastructure
  // ----------------------------------------------------------------------------

  // The "infra" group tracks CI infrastructure management workflows.
  group "infra" {
    match {
      base_dir = [joinpath(".github", "workflows")]
      base_name = [
        "test-ci-bootstrap.yml",
        "test-ci-cleanup.yml",
      ]
    }

    match {
      base_dir = [
        joinpath("enos", "ci"),
      ]
    }
  }

  // ----------------------------------------------------------------------------
  // Quality
  // ----------------------------------------------------------------------------

  // The "quality_code" group tracks code quality and linting workflows.
  group "quality_code" {
    match {
      base_dir = [joinpath(".github", "workflows")]
      base_name = [
        "code-checker.yml",
        "copywrite.yml",
        "do-not-merge-checker.yml",
      ]
    }
  }

  // The "quality_actions" group tracks GitHub Actions linting.
  group "quality_actions" {
    match {
      base_dir = [
        joinpath(".github", "workflows"),
        joinpath(".github", "actions"),
      ]
    }

    match {
      base_dir  = [".github"]
      base_name = ["actionlint.yml"]
    }
  }


  // ----------------------------------------------------------------------------
  // Releases
  // ----------------------------------------------------------------------------

  // The "release_automation" group tracks release process automation workflows.
  group "release_automation" {
    match {
      base_dir = [joinpath(".github", "workflows")]
      base_name = [
        "cancel-stale-release-runs-ent.yml",
        "changelog-curation-ent.yml",
        "check-go-version-ent.yml",
        "cicd-pipeline-checker-ent.yml",
        "collect-artifact-shas.yml",
        "generate-changelog-ent.yml",
        "missing-release-tags-check.yml",
        "open-pr-checker-ent.yml",
        "release-branch-go-ver-check-ent.yml",
        "release-build-ent.yml",
        "release-candidate-branch-cutting-ent.yml",
        "release-procedure-ent.yml",
        "release-version-checker-ent.yml",
        "trigger-promotion.yml",
        "validate-inputs-ent.yml",
        "validate-inputs.yml",
        "workflow-status-checker.yml",
      ]
    }
  }

  // The "release_testing" group tracks release artifact testing workflows.
  group "release_testing" {
    match {
      base_dir = [joinpath(".github", "workflows")]
      base_name = [
        "enos-release-testing-ent.yml",
        "enos-release-testing-oss.yml",
      ]
    }
  }

  // ----------------------------------------------------------------------------
  // Security
  // ----------------------------------------------------------------------------

  // The "security_scanner" group tracks general security scanning workflows.
  group "security_scanner" {
    match {
      base_dir  = [joinpath(".github", "workflows")]
      base_name = ["security-scan.yml"]
    }
  }

  // The "security_zap" group tracks OWASP ZAP security scanning.
  group "security_zap" {
    match {
      contains = [
        "security-scan-zap",
        "zap_scan",
      ]
    }
  }

  // The "security_mend" group tracks Mend (formerly WhiteSource) security scanning.
  group "security_mend" {
    match {
      base_dir  = [joinpath(".github", "workflows")]
      base_name = ["mend-pr-scan.yml"]
    }
  }

  // ----------------------------------------------------------------------------
  // Testing
  // ----------------------------------------------------------------------------

  // The "test_autopilot" group tracks Autopilot code and testing workflows.
  group "test_autopilot" {
    match {
      extension = [".go"]
      contains  = ["raft_autopilot"]
    }

    match {
      base_dir  = [joinpath(".github", "actions")]
      base_name = ["run-apupgrade-tests.yml"]
    }

    match {
      base_dir = [joinpath(".github", "workflows")]
      base_name = [
        "run-apupgrade-tests-ent.yml",
        "vault-autopilot-test-ent.yml",
        "vault-replication-test-ent.yml",
      ]
    }
  }
  // The "test_ci" group tracks core CI and Go testing workflows.
  group "test_ci" {
    match {
      base_dir = [joinpath(".github", "workflows")]
      base_name = [
        "ci.yml",
        "test-go.yml",
        "nightly-tests-ent.yml",
        "test-acc-dockeronly-nightly.yml",
        "test-run-acc-tests-for-path.yml",
      ]
    }
  }

  // The "test_cloud" group tracks HCP/Cloud testing.
  group "test_cloud" {
    match {
      file = [
        joinpath(".github", "workflows", "test-run-enos-scenario-cloud.yml"),
        joinpath("enos", "enos-scenario-cloud-ent.hcl"),
      ]
    }

    match {
      base_dir = [
        joinpath("enos", "modules", "cloud_docker_vault_cluster"),
        joinpath("enos", "modules", "hcp"),
        joinpath("tools", "pipeline", "internal", "pkg", "hcp"),
      ]
    }

    match {
      base_dir = [
        joinpath("tools", "pipeline", "internal", "cmd"),
      ]
      contains = ["hcp"]
    }
  }

  // The "test_enos" group tracks Enos scenarios, modules, and testing workflows.
  group "test_enos" {
    ignore {
      base_dir = [
        joinpath("enos", "ci"),
      ]
    }

    match {
      base_dir = [
        joinpath(".github", "actions", "run-enos-scenario"),
        joinpath(".github", "actions", "run-enos-scenario"),
        joinpath("vault", "external_tests", "blackbox"),
        joinpath("sdk", "helper", "testcluster", "blackbox"),
        "enos",
      ]
    }

    match {
      base_dir = [joinpath(".github", "workflows")]
      base_name = [
        "enos-lint.yml",
        "test-enos-scenario-cloud.yml",
        "test-enos-scenario-containers.yml",
        "test-enos-scenario-matrix.yml",
        "test-enos-scenario-ui.yml",
        "test-run-enos-scenario-cloud.yml",
        "test-run-enos-scenario-containers.yml",
        "test-run-enos-scenario-matrix.yml",
        "test-run-enos-scenario.yml",
      ]
    }
  }

  // The "test_ui" group tracks UI testing workflows
  group "test_ui" {
    match {
      base_dir = [
        joinpath(".github", "workflows", "test-run-enos-scenario-ui"),
        joinpath(".github", "workflows", "test-ui"),
      ]
    }
  }

  // ----------------------------------------------------------------------------
  // Web UI
  // ----------------------------------------------------------------------------

  group "ui" {
    match {
      base_dir = [
        "ui",
      ]
    }
  }

  group "ui_deps" {
    match {
      file = [
        joinpath("ui", "package-lock.json"),
        joinpath("ui", "pnpm-lock.yaml"),
      ]
    }
  }
}
