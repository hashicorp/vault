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
      ]
    }

    // Ignore whole directories where no enterprise only code should ever exist.
    ignore {
      base_dir = [
        joinpath(".release", "docker"),
        joinpath("enos", "modules"),
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

    // Match whole directories that are enterprise only
    match {
      base_dir = [
        "vault_ent",
        joinpath("scripts", "dev", "hsm"),
        joinpath("scripts", "testing"),
        joinpath("specs"),
        joinpath(".release", "ibm-pao"),
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
        "_ent",
        "_ent.pb.go",
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

  // The "ui" group matches the Web UI source
  group "ui" {
    match {
      base_dir = ["ui"]
    }
  }
}
