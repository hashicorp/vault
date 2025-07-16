// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package changed

import (
	"context"
	"os"
	"path/filepath"
	"strings"
)

// FileGroupCheck is a function that takes a reference to a changed file and returns groups that
// the file belongs to
type FileGroupCheck func(context.Context, *File) FileGroups

// DefaultFileGroupCheckers are the default file group checkers
var DefaultFileGroupCheckers = []FileGroupCheck{
	FileGroupCheckerApp,
	FileGroupCheckerAutopilot,
	FileGroupCheckerChangelog,
	FileGroupCheckerCommunity,
	FileGroupCheckerDocs,
	FileGroupCheckerEnos,
	FileGroupCheckerEnterprise,
	FileGroupCheckerGoToolchain,
	FileGroupCheckerPipeline,
	FileGroupCheckerProto,
	FileGroupCheckerWebUI,
}

// Group takes a context, a file, and one-to-many file group checkers and adds group metadata to
// the file.
func Group(ctx context.Context, file *File, checkers ...FileGroupCheck) {
	if file == nil || len(checkers) < 1 {
		return
	}

	for _, check := range checkers {
		file.Groups = file.Groups.Add(check(ctx, file)...)
	}
}

// GroupFiles takes a context, a slice of files, and one-to-many file group checkers and adds group
// metadata to the files.
func GroupFiles(ctx context.Context, files []*File, checkers ...FileGroupCheck) {
	for _, file := range files {
		Group(ctx, file, checkers...)
	}
}

// FileGroupCheckerApp is a file group checker that groups based on the file being part of the Vault
// Go app
func FileGroupCheckerApp(ctx context.Context, file *File) FileGroups {
	name := file.Name()
	ext := filepath.Ext(name)

	switch {
	case hasBaseDir(name, filepath.Join("tools", "pipeline")):
		return nil
	case
		ext == ".go",
		strings.HasSuffix(name, "go.mod"),
		strings.HasSuffix(name, "go.sum"):
		return FileGroups{FileGroupGoApp}
	default:
		return nil
	}
}

// FileGroupCheckerAutopilot is a file group checker that groups based on the file being part of the
// raft autopilot system
func FileGroupCheckerAutopilot(ctx context.Context, file *File) FileGroups {
	name := file.Name()
	ext := filepath.Ext(name)

	if ext == ".go" && strings.Contains(name, "raft_autopilot") {
		return FileGroups{FileGroupAutopilot}
	}

	return nil
}

// FileGroupCheckerChangelog is a file group checker that groups based on the file being part of the
// CHANGELOG
func FileGroupCheckerChangelog(ctx context.Context, file *File) FileGroups {
	name := file.Name()
	if strings.HasPrefix(name, "CHANGELOG") || hasBaseDir(name, "changelog") {
		return FileGroups{FileGroupChangelog}
	}

	return nil
}

// FileGroupCheckerCommunity is a file group checker that groups based on the file being part of the
// Vault App but a community only file.
func FileGroupCheckerCommunity(ctx context.Context, file *File) FileGroups {
	name := file.Name()

	switch filepath.Ext(name) {
	case ".go":
		if strings.HasSuffix(name, "_oss.go") || strings.HasSuffix(name, "_ce.go") {
			return FileGroups{FileGroupCommunity}
		}
	case
		".hcl",
		".md",
		".sh",
		".yaml",
		".yml":
		switch {
		case
			strings.Contains(name, "-ce"),
			strings.Contains(name, "_ce"),
			strings.Contains(name, "-oss"),
			strings.Contains(name, "_oss"):
			return FileGroups{FileGroupCommunity}
		}
	}

	return nil
}

// FileGroupCheckerDocs is a file group checker that groups based on the file being part of the
// documenation.
func FileGroupCheckerDocs(ctx context.Context, file *File) FileGroups {
	name := file.Name()

	if strings.HasPrefix(name, "README.md") || hasBaseDir(name, "website") {
		return FileGroups{FileGroupDocs}
	}

	return nil
}

// FileGroupCheckerEnos is a file group checker that groups based on the file being part of the
// enos testing framework.
func FileGroupCheckerEnos(ctx context.Context, file *File) FileGroups {
	name := file.Name()

	if strings.Contains(name, "enos") || hasBaseDir(name, "enos") {
		return FileGroups{FileGroupEnos}
	}

	return nil
}

// FileGroupCheckerEnterprise is a file group checker that groups based on the file being part of
// the Vault App but an enterprise only file. Ideally enterprise only files will use common filename
// schema or directories to reduce our logic here, but some legacy files have been added here.
// NOTE: Even if we miss a file or two the sky will not fall, only our automation for CE backports
// could theoretically miss a file and require the author to extract it out themselves before they
// merge it. Since such files are created in Vault Enterprise there is little risk.
func FileGroupCheckerEnterprise(ctx context.Context, file *File) FileGroups {
	name := file.Name()

	// Base directory checks
	switch {
	case
		hasBaseDir(name, "vault_ent"),
		hasBaseDir(name, filepath.Join("scripts", "dev", "hsm")),
		hasBaseDir(name, filepath.Join("scripts", "testing")),
		hasBaseDir(name, filepath.Join("specs")):
		return FileGroups{FileGroupEnterprise}
	}

	// File name checks
	switch filepath.Base(name) {
	case
		"Dockerfile-ent",
		"Dockerfile-ent-hsm":
		return FileGroups{FileGroupEnterprise}
	}

	// File extension checks
	switch filepath.Ext(name) {
	case ".go":
		switch {
		case
			strings.HasSuffix(name, "_ent.go"),
			strings.HasSuffix(name, "_ent_test.go"),
			strings.Contains(name, "_ent") && strings.HasSuffix(name, ".pb.go"):
			return FileGroups{FileGroupEnterprise}
		}
	case ".txt":
		if hasBaseDir(name, "changelog") && strings.HasPrefix(filepath.Base(name), "_") {
			return FileGroups{FileGroupEnterprise}
		}
	case
		".proto",
		".hcl",
		".md",
		".sh",
		".yaml",
		".yml":
		switch {
		case strings.Contains(name, "-ce"): // Skip workflows that might have ce and ent in the name
		case
			strings.Contains(name, "-ent"),
			strings.Contains(name, "_ent"),
			strings.Contains(name, "hsm"),
			strings.Contains(name, "merkle-tree"):
			return FileGroups{FileGroupEnterprise}
		}
	}

	return nil
}

// FileGroupCheckerGoToolchain is a file group checker that groups based on the file modifying the
// Go toolchain or dependencies.
func FileGroupCheckerGoToolchain(ctx context.Context, file *File) FileGroups {
	name := file.Name()
	switch {
	case
		name == ".go-version",
		strings.HasSuffix(name, "go.mod"),
		strings.HasSuffix(name, "go.sum"):
		return FileGroups{FileGroupGoToolchain}
	default:
		return nil
	}
}

// FileGroupCheckerPipeline is a file group checker that groups based on the file is part of the
// build or CI pipeline.
func FileGroupCheckerPipeline(ctx context.Context, file *File) FileGroups {
	name := file.Name()

	switch {
	case
		hasBaseDir(name, ".build"),
		hasBaseDir(name, ".github"),
		hasBaseDir(name, "scripts"),
		hasBaseDir(name, filepath.Join("tools", "pipeline")),
		name == "CODEOWNERS",
		name == "Dockerfile",
		name == "Makefile":
		return FileGroups{FileGroupPipeline}
	default:
		return nil
	}
}

// FileGroupCheckerProto is a file group checker that groups based on the files extension being .proto
func FileGroupCheckerProto(ctx context.Context, file *File) FileGroups {
	name := file.Name()

	ext := filepath.Ext(name)
	if ext == ".proto" || strings.HasPrefix(name, "buf.") {
		return FileGroups{FileGroupProto}
	}

	return nil
}

// FileGroupCheckerWebUI is a file group checker that groups based on the files being part of the
// web UI
func FileGroupCheckerWebUI(ctx context.Context, file *File) FileGroups {
	name := file.Name()
	if hasBaseDir(name, "ui") {
		return FileGroups{FileGroupWebUI}
	}

	return nil
}

func hasBaseDir(name, dir string) bool {
	return strings.HasPrefix(name, dir+string(os.PathSeparator))
}
