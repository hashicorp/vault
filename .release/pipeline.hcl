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
  // The "app" group is the "vault" Go application.
  group "app" {
    ignore {
      base_dir = [
        // There's no reason to consider the pipeline tool as part of the app
        joinpath("tools", "pipeline"),
      ]
    }

    match {
      // Match all Go source code files
      extension = [".go"]
    }

    match {
      // Match go.mod or go.sum
      base_name = [
        "go.mod",
        "go.sum",
      ]
    }
  }

  // The "autopilot" group is for application code that is specific to autopilot.
  // This group is used to run some specific autopilot tests when the code changes.
  group "autopilot" {
    match {
      // Match go files that contain "raft_autopilot" in the path name
      extension = [".go"]
      contains  = ["raft_autopilot"]
    }
  }

  // The "changelog" group is our change log
  group "changelog" {
    match {
      base_dir = ["changelog"]
    }

    match {
      contains = ["CHANGELOG"]
    }
  }

  // The "community" group is for files that exist specifically for either CE
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

  // The "docs" group is our README and static website assets.
  group "docs" {
    match {
      base_name = ["README.md"]
    }

    match {
      base_dir = ["website"]
    }
  }

  // The "enos" group matches everything in the "enos" directory.
  group "enos" {
    match {
      base_dir = ["enos"]
    }
  }

  // The "enterprise" group is arguably the most important grouping configuration
  // in this repository. Among other things, this is how we enforce that no
  // enterprise files make it into CE branches or repositories.
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

  // The "github" group matches the .github directory
  group "github" {
    match {
      base_dir = [".github"]
    }
  }

  // The "gotoolchain" group matches changes to the Go toolchain used to build
  // the application
  group "gotoolchain" {
    match {
      base_name = [
        ".go-version",
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

  // The "go_app" group tracks core Vault application code (excluding tests).
  // This is the primary group checked by CI workflows to trigger Go tests.
  // Triggers full build, test cycle, and integration tests.
  group "go_app" {
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
      contains = [
        "hcp"
      ]
    }
  }

  // The "pipeline" group matches directories where we house code and
  // configuration used in the CI/CD pipeline
  group "pipeline" {
    match {
      base_dir = [
        ".build",
        joinpath(".github", "workflows"),
        joinpath(".github", "actions"),
        joinpath(".github", "scripts"),
        ".hooks",
        ".release",
        "scripts",
        joinpath("tools", "pipeline"),
      ]
    }

    match {
      base_name = [
        "Dockerfile",
        "Makefile",
      ]
    }
  }

  // The "proto" group matches protobuf files and configuration.
  group "proto" {
    match {
      extension = [".proto"]
    }

    match {
      base_name_prefix = ["buf."]
    }
  }

  // The "ui" group matches files for the Web UI
  group "ui" {
    match {
      base_dir = [
        joinpath(".github", "actions", "build-ui"),
        joinpath(".github", "actions", "setup-pnpm"),
        joinpath(".github", "workflows", "gen-diff-spec"),
        joinpath(".github", "workflows", "test-run-enos-scenario-ui"),
        joinpath(".github", "workflows", "test-ui"),
        joinpath(".github", "workflows", "ui-client-update"),
        "ui",
      ]
    }
  }

  // The "zap_scan" group matches the Zap scanner
  group "zap_scan" {
    match {
      contains = [
        "security-scan-zap",
        "zap_scan",
      ]
    }
  }
}
