// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestTidyUntilStable verifies that we repeat the tidy of every module until a
// pass changes none of them, so that a module tidied before a module it depends
// on is tidied again once that module is done, and that we stop as soon as a
// pass changes nothing. The tidy is faked: tidying module a copies the go.mod
// of module b into a's go.sum, as a stand-in for a depending on b, and tidying
// b rewrites its go.mod. With a listed first, the first pass leaves a with b's
// untidy go.mod and only the second pass brings it up to date.
func TestTidyUntilStable(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeFiles(t, dir, map[string]string{
		"a/go.mod": "module example.com/a\n",
		"b/go.mod": "untidy\n",
	})
	aMod, aSum := filepath.Join(dir, "a/go.mod"), filepath.Join(dir, "a/go.sum")
	bMod := filepath.Join(dir, "b/go.mod")

	var calls []string
	failed, err := tidyUntilStable(t.Context(), []string{aMod, bMod}, func(_ context.Context, p string) error {
		calls = append(calls, filepath.Base(filepath.Dir(p)))

		if p == aMod {
			b, err := os.ReadFile(bMod)
			require.NoError(t, err)

			return os.WriteFile(aSum, b, 0o644)
		}

		return os.WriteFile(bMod, []byte("tidy\n"), 0o644)
	})
	require.NoError(t, err)
	require.Empty(t, failed)

	// The second pass changed a, so a third pass confirmed nothing else changed.
	require.Equal(t, []string{"a", "b", "a", "b", "a", "b"}, calls)

	got, err := os.ReadFile(aSum)
	require.NoError(t, err)
	require.Equal(t, "tidy\n", string(got))
}

// TestTidyUntilStable_singleModule verifies that a single module is tidied
// once, because no other module can change what it needs.
func TestTidyUntilStable_singleModule(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeFiles(t, dir, map[string]string{"go.mod": "module example.com/a\n"})

	calls := 0
	failed, err := tidyUntilStable(t.Context(), []string{filepath.Join(dir, "go.mod")}, func(_ context.Context, p string) error {
		calls++

		return os.WriteFile(p, []byte("module example.com/a\n\ngo 1.21\n"), 0o644)
	})
	require.NoError(t, err)
	require.Empty(t, failed)
	require.Equal(t, 1, calls)
}

// TestTidyUntilStable_failure verifies that a module that fails to tidy is
// reported by path and left out of the remaining passes, while the other
// modules are still tidied until they settle.
func TestTidyUntilStable_failure(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeFiles(t, dir, map[string]string{
		"a/go.mod": "module example.com/a\n",
		"b/go.mod": "module example.com/b\n",
		"c/go.mod": "module example.com/c\n",
	})
	aMod, bMod, cMod := filepath.Join(dir, "a/go.mod"), filepath.Join(dir, "b/go.mod"), filepath.Join(dir, "c/go.mod")

	boom := errors.New("boom")
	calls := map[string]int{}
	failed, err := tidyUntilStable(t.Context(), []string{aMod, bMod, cMod}, func(_ context.Context, p string) error {
		calls[p]++

		switch p {
		case bMod:
			return boom
		case aMod:
			// Changes on the first pass only.
			if calls[p] == 1 {
				return os.WriteFile(p, []byte("module example.com/a\n\ngo 1.21\n"), 0o644)
			}
		}

		return nil
	})
	require.NoError(t, err)
	require.Equal(t, map[string]error{bMod: boom}, failed)
	require.Equal(t, 1, calls[bMod], "a module that failed isn't tidied again")
	require.Equal(t, 2, calls[aMod])
	require.Equal(t, 2, calls[cMod])
}

// TestTidyUntilStable_allFail verifies that nothing more is tried when every
// module fails.
func TestTidyUntilStable_allFail(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeFiles(t, dir, map[string]string{
		"a/go.mod": "module example.com/a\n",
		"b/go.mod": "module example.com/b\n",
	})
	aMod, bMod := filepath.Join(dir, "a/go.mod"), filepath.Join(dir, "b/go.mod")

	calls := 0
	failed, err := tidyUntilStable(t.Context(), []string{aMod, bMod}, func(context.Context, string) error {
		calls++

		return errors.New("boom")
	})
	require.NoError(t, err)
	require.Len(t, failed, 2)
	require.Equal(t, 2, calls)
}

// TestTidyUntilStable_neverSettles verifies that we give up with an error after
// the most passes rather than looping when every pass changes a module.
func TestTidyUntilStable_neverSettles(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeFiles(t, dir, map[string]string{
		"a/go.mod": "module example.com/a\n",
		"b/go.mod": "module example.com/b\n",
	})
	aMod, bMod := filepath.Join(dir, "a/go.mod"), filepath.Join(dir, "b/go.mod")

	calls := 0
	failed, err := tidyUntilStable(t.Context(), []string{aMod, bMod}, func(_ context.Context, p string) error {
		calls++

		return os.WriteFile(p, []byte(string(rune('a'+calls))+"\n"), 0o644)
	})
	require.ErrorContains(t, err, "still changing")
	require.Empty(t, failed)
	require.Equal(t, 2*maxTidyPasses, calls)
}

// TestTidyUntilStable_noModules verifies that there's nothing to do, and no
// error, without modules.
func TestTidyUntilStable_noModules(t *testing.T) {
	t.Parallel()

	failed, err := tidyUntilStable(t.Context(), nil, func(context.Context, string) error {
		return errors.New("must not be called")
	})
	require.NoError(t, err)
	require.Empty(t, failed)
}

// TestRunGoModTidy_missingGoMod verifies that a workspace entry without its own
// go.mod cannot redirect tidy to an enclosing module and rewrite that module.
func TestRunGoModTidy_missingGoMod(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	const parentMod = "module example.com/root\n\ngo 1.21\n\nrequire example.com/unused v1.0.0\n\nreplace example.com/unused => ./unused\n"
	writeFiles(t, dir, map[string]string{"go.mod": parentMod, "root.go": "package root\n", "tools/.keep": ""})

	missing := filepath.Join(dir, "tools/go.mod")
	require.ErrorIs(t, RunGoModTidy(t.Context(), missing), os.ErrNotExist)
	_, snapshotErr := readModuleFiles([]string{missing})
	require.ErrorIs(t, snapshotErr, os.ErrNotExist)
	got, readErr := os.ReadFile(filepath.Join(dir, "go.mod"))
	require.NoError(t, readErr)
	require.Equal(t, parentMod, string(got))
}
