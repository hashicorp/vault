// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

changed_files {
  group "app" {
    ignore {
      base_dir = [
        joinpath("tools", "pipeline"),
      ]
    }

    match {
      extension = [".go"]
    }

    match {
      base_name = [
        "go.mod",
        "go.sum",
      ]
    }
  }

  group "autopilot" {
    match {
      extension = [".go"]
      contains  = ["raft_autopilot"]
    }
  }

  group "changelog" {
    match {
      base_dir = ["changelog"]
    }

    match {
      contains = ["CHANGELOG"]
    }
  }

  group "community" {
    match {
      extension = [".go"]
      contains = [
        "_oss.go",
        "_ce.go",
      ]
    }

    match {
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

  group "docs" {
    match {
      base_name = ["README.md"]
    }

    match {
      base_dir = ["website"]
    }
  }

  group "enos" {
    match {
      base_dir = ["enos"]
    }
  }

  group "enterprise" {
    ignore {
      file = [
        # These exist on CE branches to please Github Actions.
        joinpath(".github", "workflows", "build-artifacts-ent.yml"),
        joinpath(".github", "workflows", "backport-automation-ent.yml"),
      ]
    }

    ignore {
      base_dir = [
        joinpath(".release", "docker"),
        joinpath("enos", "modules"),
        joinpath("scripts", "docker"),
      ]
    }

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

  group "github" {
    match {
      base_dir = [".github"]
    }
  }

  group "gotoolchain" {
    match {
      base_name = [
        ".go-version",
        "go.mod",
        "go.sum",
      ]
    }
  }

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

  group "proto" {
    match {
      extension = [".proto"]
    }

    match {
      base_name_prefix = ["buf."]
    }
  }

  group "ui" {
    match {
      base_dir = ["ui"]
    }
  }
}
