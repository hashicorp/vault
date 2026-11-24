// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// ValidateGoModPath checks that p names an existing module file within dir,
// without symlinked directories or module metadata that Git cannot safely stage.
func ValidateGoModPath(ctx context.Context, dir, p string) error {
	slog.Default().DebugContext(ctx, "validating module path",
		slog.String("dir", dir), slog.String("path", p),
	)
	if !filepath.IsLocal(p) {
		return fmt.Errorf("module path %q is not inside %s", p, dir)
	}

	current := dir
	for _, part := range strings.Split(filepath.Clean(p), string(filepath.Separator)) {
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if err != nil {
			return fmt.Errorf("reading %s: %w", current, err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("module path %s contains a symlink", current)
		}
	}

	info, err := os.Stat(current)
	if err != nil {
		return fmt.Errorf("reading %s: %w", current, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("module file %s is not a regular file", current)
	}

	sum := filepath.Join(filepath.Dir(current), "go.sum")
	info, err = os.Lstat(sum)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("reading %s: %w", sum, err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("module file %s is a symlink", sum)
	}

	return nil
}
