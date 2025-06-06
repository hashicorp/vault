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
	FileGroupCheckerDir,
	FileGroupCheckerFileName,
	FileGroupCheckerFileGo,
	FileGroupCheckerFileProto,
}

// FileGroupCheckerDir is a file group checker that groups based on the files directory
func FileGroupCheckerDir(ctx context.Context, file *File) FileGroups {
	name := file.Name()
	groups := FileGroups{}

	for dir, groups := range map[string]FileGroups{
		".github":   groups.Add(FileGroupPipeline),
		"changelog": groups.Add(FileGroupChangelog),
		"enos":      groups.Add(FileGroupEnos),
		"tools":     groups.Add(FileGroupTools),
		"ui":        groups.Add(FileGroupWebUI),
		"website":   groups.Add(FileGroupDocs),
	} {
		if strings.HasPrefix(name, dir+string(os.PathSeparator)) {
			return groups
		}
	}

	return nil
}

// FileGroupCheckerFileName is a file group checker that groups based on the files name
func FileGroupCheckerFileName(ctx context.Context, file *File) FileGroups {
	name := file.Name()
	groups := FileGroups{}

	switch {
	case strings.HasPrefix(name, "buf."):
		return groups.Add(FileGroupProto)
	case strings.HasPrefix(name, "CHANGELOG"):
		return groups.Add(FileGroupChangelog)
	case strings.HasPrefix(name, "CODEOWNERS"):
		return groups.Add(FileGroupPipeline)
	case strings.HasSuffix(name, "go.mod") || strings.HasSuffix(name, "go.sum"):
		return groups.Add(FileGroupGoModules, FileGroupGoApp)
	}

	return nil
}

// FileGroupCheckerFileGo is a file group checker that groups based on the files extension being .go
func FileGroupCheckerFileGo(ctx context.Context, file *File) FileGroups {
	name := file.Name()

	ext := filepath.Ext(name)
	if ext != ".go" {
		return nil
	}
	groups := FileGroups{}
	groups = groups.Add(FileGroupGoApp)

	if strings.Contains(name, "raft_autopilot") {
		groups = groups.Add(FileGroupAutopilot)
	}

	if strings.HasSuffix(name, "_ent.go") {
		groups = groups.Add(FileGroupEnterprise)
	} else if strings.HasSuffix(name, "_oss.go") || strings.HasSuffix(name, "_ce.go") {
		groups = groups.Add(FileGroupCommunity)
	}

	if strings.HasPrefix(name, "tools/pipeline") {
		groups = groups.Add(FileGroupPipeline)
	}

	return groups
}

// FileGroupCheckerFileProto is a file group checker that groups based on the files extension being .proto
func FileGroupCheckerFileProto(ctx context.Context, file *File) FileGroups {
	name := file.Name()

	ext := filepath.Ext(name)
	if ext != ".proto" {
		return nil
	}

	return FileGroups{FileGroupProto}
}
