// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"fmt"
	"log/slog"

	libgit "github.com/hashicorp/vault/tools/pipeline/internal/pkg/git"
	slogctx "github.com/veqryn/slog-context"
)

func resetAndCreateNOOPCommit(ctx context.Context, git *libgit.Client, baseRef string) error {
	ctx = slogctx.Append(ctx, slog.String("base-ref", baseRef))
	slog.Default().DebugContext(ctx, "hard resetting current checkout and creating no-op commit")

	resetRes, err := git.Reset(ctx, &libgit.ResetOpts{
		Mode:    libgit.ResetModeHard,
		Treeish: baseRef,
	})
	if err != nil {
		return fmt.Errorf("resetting back to base reference: %s: %w", resetRes.String(), err)
	}

	commitRes, err := git.Commit(ctx, &libgit.CommitOpts{
		AllowEmpty: true,
		Message:    "no-op commit",
		NoVerify:   true,
		NoEdit:     true,
	})
	if err != nil {
		return fmt.Errorf("committing no-op commit: %s: %w", commitRes.String(), err)
	}

	return nil
}
