// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package changed

import (
	"context"
)

// FileGroupCheck is a function that takes a reference to a changed file and returns groups that
// the file belongs to.
type FileGroupCheck func(context.Context, *File) FileGroups

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

// Groups takes a set of Files and returns a set of all FileGroups.
func Groups(files Files) FileGroups {
	groups := FileGroups{}
	for _, file := range files {
		for _, group := range file.Groups {
			groups = groups.Add(group)
		}
	}

	return groups
}
