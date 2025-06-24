// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package git

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test our opts structs ability to render the correct flags in the correct order
// NOTE: Many of thests use incompatible options but that's not what we care about,
// we're simply asserting that the rendered string matches what ought to be there
// give the config.
// We have chosen not to try and very flag combinations. Instead we render it
// and execute it and rely on git to handle validation of options.
func TestOptsStringers(t *testing.T) {
	t.Parallel()

	for name, expect := range map[string]struct {
		opts     OptStringer
		expected string
	}{
		"am": {
			&AmOpts{
				AllowEmpty:                true, // Only supported for --resolved
				CommitterDateIsAuthorDate: true,
				Empty:                     EmptyCommitKeep,
				Keep:                      true,
				KeepNonPatch:              true,
				MessageID:                 true,
				NoMessageID:               true,
				NoReReReAutoupdate:        true,
				NoVerify:                  true,
				Quiet:                     true,
				ReReReAutoupdate:          true,
				Signoff:                   true,
				ThreeWayMerge:             true,
				Whitespace:                ApplyWhitespaceActionFix,
				Mbox:                      []string{"/path/to/my.patch"},
			},
			"--committer-date-is-author-date --empty=keep --keep --keep-non-patch --message-id --no-message-id --no-rerere-autoupdate --no-verify --quiet --rerere-autoupdate --signoff --3way --whitespace=fix /path/to/my.patch",
		},
		"am --continue": {
			&AmOpts{
				// Unallowed options are ignored
				Empty:      EmptyCommitKeep,
				AllowEmpty: true,
				// Sequence
				Continue: true,
			},
			"--continue",
		},
		"am --abort": {
			&AmOpts{
				// Unallowed options are ignored
				Empty:      EmptyCommitKeep,
				AllowEmpty: true,
				// Sequence
				Abort: true,
			},
			"--abort",
		},
		"am --quit": {
			&AmOpts{
				// Unallowed options are ignored
				Empty:      EmptyCommitKeep,
				AllowEmpty: true,
				// Sequence
				Quit: true,
			},
			"--quit",
		},
		"am --allow-empty --resolved": {
			&AmOpts{
				// Unallowed options are ignored
				Empty: EmptyCommitKeep,
				// Allowed options are kept
				AllowEmpty: true,
				// Sequence
				Resolved: true,
			},
			"--allow-empty --resolved",
		},
		"am --retry": {
			&AmOpts{
				// Unallowed options are ignored
				Empty:      EmptyCommitKeep,
				AllowEmpty: true,
				// Sequence
				Retry: true,
			},
			"--retry",
		},
		"apply": {
			&ApplyOpts{
				AllowEmpty:    true,
				Cached:        true,
				Check:         true,
				Index:         true,
				Ours:          true,
				Recount:       true,
				Stat:          true,
				Summary:       true,
				Theirs:        true,
				ThreeWayMerge: true,
				Union:         true,
				Whitespace:    ApplyWhitespaceActionFix,
				Patch:         []string{"/path/to/my.diff"},
			},
			"--allow-empty --cached --check --index --ours --recount --stat --summary --theirs --3way --union --whitespace=fix /path/to/my.diff",
		},
		"branch copy": {
			&BranchOpts{
				Copy:      true,
				Force:     true,
				OldBranch: "my-old-branch",
				NewBranch: "my-new-branch",
			},
			"--copy --force my-old-branch my-new-branch",
		},
		"branch delete": {
			&BranchOpts{
				Delete:     true,
				Remotes:    true,
				BranchName: "my-branch",
			},
			"--delete --remotes my-branch",
		},
		"branch move": {
			&BranchOpts{
				Move:      true,
				OldBranch: "my-old-branch",
				NewBranch: "my-new-branch",
			},
			"--move my-old-branch my-new-branch",
		},
		"branch upstream set": {
			&BranchOpts{
				SetUpstream:   true,
				SetUpstreamTo: "my-upstream",
				BranchName:    "my-branch",
			},
			"--set-upstream --set-upstream-to=my-upstream my-branch",
		},
		"branch upstream unset": {
			&BranchOpts{
				UnsetUpstream: true,
				BranchName:    "my-branch",
			},
			"--unset-upstream my-branch",
		},
		"branch track": {
			&BranchOpts{
				Track:      BranchTrackInherit,
				NoTrack:    true,
				Force:      true,
				BranchName: "my-branch",
				StartPoint: "HEAD~2",
			},
			"--force --no-track --track=inherit my-branch HEAD~2",
		},
		"branch with pattern": {
			// Everything else in branch..
			&BranchOpts{
				Abbrev:      7,
				All:         true,
				Contains:    "abcd1234",
				Format:      "%%",
				List:        true,
				Merged:      "1234abcd",
				NoColor:     true,
				NoColumn:    true,
				PointsAt:    "12ab34cd",
				Remotes:     true,
				ShowCurrent: true,
				Sort:        "key",
				Pattern:     []string{"my/dir", "another/dir"},
			},
			"--abbrev=7 --all --contains=abcd1234 --format=%% --list --merged=1234abcd --no-color --no-column --points-at=12ab34cd --remotes --show-current --sort=key my/dir another/dir",
		},
		"checkout 1/2 opts": {
			&CheckoutOpts{
				Branch:                 "source",
				NewBranchForceCheckout: "new",
				Force:                  true,
				NoTrack:                true,
				Ours:                   true,
				Quiet:                  true,
			},
			"-B new --force --no-track --ours --quiet source",
		},
		"checkout 2/2 opts": {
			&CheckoutOpts{
				Branch:     "source",
				NewBranch:  "new",
				Guess:      true,
				Orphan:     "bar",
				Progress:   true,
				Theirs:     true,
				Track:      BranchTrackDirect,
				StartPoint: "HEAD~1",
			},
			"-b new --guess --orphan bar --progress --theirs --track=direct source HEAD~1",
		},
		"checkout path spec": {
			&CheckoutOpts{
				Branch:     "main",
				StartPoint: "HEAD~1",
				PathSpec:   []string{"go.mod", "go.sum"},
			},
			"main HEAD~1 -- go.mod go.sum",
		},
		"cherry-pick 1/2 opts": {
			&CherryPickOpts{
				AllowEmpty:        true,
				AllowEmptyMessage: true,
				Empty:             EmptyCommitKeep,
				FF:                true,
				GPGSign:           true,
				Mainline:          "ABCDEFGH",
				Record:            true,
				Signoff:           true,
				Commit:            "1234ABCD",
			},
			"--allow-empty --allow-empty-message --empty=keep --ff --gpg-sign --mainline=ABCDEFGH -x --signoff 1234ABCD",
		},
		"cherry-pick: 2/2 opts": {
			&CherryPickOpts{
				GPGSignKeyID:     "4321DCBA",
				ReReReAutoupdate: true,
				Strategy:         MergeStrategyResolve,
				StrategyOptions: []MergeStrategyOption{
					MergeStrategyOptionDiffAlgorithmHistogram,
					MergeStrategyOptionIgnoreSpaceChange,
				},
				Commit: "1234ABCD",
			},
			"--gpg-sign=4321DCBA --rerere-autoupdate --strategy=resolve --strategy-option=diff-algorithm=histogram --strategy-option=ignore-space-change 1234ABCD",
		},
		"cherry-pick --continue": {
			&CherryPickOpts{
				Continue: true,
				// Options are ignored
				Commit:       "1234ABCD",
				GPGSignKeyID: "4321DCBA",
			},
			"--continue",
		},
		"cherry-pick --abort": {
			&CherryPickOpts{
				Abort: true,
				// Options are ignored
				Commit:       "1234ABCD",
				GPGSignKeyID: "4321DCBA",
			},
			"--abort",
		},
		"cherry-pick --quit": {
			&CherryPickOpts{
				Quit: true,
				// Options are ignored
				Commit:       "1234ABCD",
				GPGSignKeyID: "4321DCBA",
			},
			"--quit",
		},
		"clone 1/2 opts": {
			&CloneOpts{
				Branch:     "my-branch",
				Depth:      3,
				NoCheckout: true,
				NoTags:     true,
				Origin:     "my-fork",
				Quiet:      true,
				Directory:  "some-dir",
			},
			"--branch my-branch --depth 3 --no-checkout --no-tags --origin my-fork --quiet -- some-dir",
		},
		"clone 2/2 opts": {
			&CloneOpts{
				Branch:       "my-branch",
				Progress:     true,
				Sparse:       true,
				SingleBranch: true,
				Repository:   "my-repo",
				Directory:    "some-dir",
			},
			"--branch my-branch --progress --single-branch --sparse -- my-repo some-dir",
		},
		"commit 1/2 opts": {
			&CommitOpts{
				All:               true,
				AllowEmpty:        true,
				AllowEmptyMessage: true,
				Amend:             true,
				Author:            "example@hashicorp.com",
				Branch:            true,
				Cleanup:           CommitCleanupModeWhitespace,
				Date:              "1 day ago",
				DryRun:            true,
				File:              "path/to/message/file",
				Fixup: &CommitFixup{
					FixupLog: CommitFixupLogReword,
					Commit:   "1234ABCD",
				},
				GPGSign: true,
				Long:    true,
				NoEdit:  true,
				PathSpec: []string{
					"file/a",
					"another/b",
				},
			},
			"--all --allow-empty --allow-empty-message --amend --author=example@hashicorp.com --branch --cleanup=whitespace --date=1 day ago --dry-run --file=path/to/message/file --fixup=reword:1234ABCD --gpg-sign --long --no-edit -- file/a another/b",
		},
		"commit 2/2 opts": {
			&CommitOpts{
				GPGSignKeyID:  "4321DCBA",
				Patch:         true,
				Porcelain:     true,
				Message:       "my commit message",
				NoPostRewrite: true,
				NoVerify:      true,
				Null:          true,
				Only:          true,
				ResetAuthor:   true,
				ReuseMessage:  "1234ABCD",
				Short:         true,
				Signoff:       true,
				Status:        true,
				Verbose:       true,
				PathSpec: []string{
					"file/a",
					"another/b",
				},
			},
			"--gpg-sign=4321DCBA --patch --porcelain --message=my commit message --no-post-rewrite --no-verify --null --only --reset-author --reuse-message=1234ABCD --short --status --verbose -- file/a another/b",
		},
		"fetch": {
			&FetchOpts{
				All:         true,
				Atomic:      true,
				Depth:       5,
				Deepen:      6,
				Force:       true,
				NoTags:      true,
				Porcelain:   true,
				Progress:    true,
				Prune:       true,
				Quiet:       true,
				SetUpstream: true,
				Unshallow:   true,
				Verbose:     true,
				Repository:  "my-repo",
				Refspec:     []string{"my-branch"},
			},
			"--all --atomic --depth 5 --deepen 6 --force --no-tags --porcelain --progress --prune --quiet --set-upstream --unshallow --verbose my-repo my-branch",
		},
		"merge 1/2 opts": {
			&MergeOpts{
				Autostash:        true,
				DoCommit:         true,
				Commit:           "1234ABCD",
				File:             "/path/to/file",
				FF:               true,
				FFOnly:           true,
				IntoName:         "my-other-branch",
				Log:              2,
				Message:          "merging my branch",
				Progress:         true,
				ReReReAutoupdate: true,
				Squash:           true,
				Stat:             true,
				Strategy:         MergeStrategyORT,
				StrategyOptions: []MergeStrategyOption{
					MergeStrategyOptionDiffAlgorithmMyers,
					MergeStrategyOptionFindRenames,
				},
				Verbose: true,
			},
			"--autostash --commit --file=/path/to/file --ff --ff-only --into-name my-other-branch --log=2 --rerere-autoupdate --squash --stat --strategy=ort --strategy-option=diff-algorithm=myers --strategy-option=find-renames 1234ABCD",
		},
		"merge 2/2 opts": {
			&MergeOpts{
				NoAutostash:        true,
				NoDoCommit:         true,
				NoFF:               true,
				NoLog:              true,
				NoProgress:         true,
				NoRebase:           true,
				NoReReReAutoupdate: true,
				NoSquash:           true,
				NoStat:             true,
				NoVerify:           true,
			},
			"--no-autostash --no-commit --no-ff --no-log --no-progress --no-rebase --no-rerere-autoupdate --no-squash --no-stat --no-stat --no-verify",
		},
		"merge --continue": {
			&MergeOpts{
				Continue: true,
				// Options are ignored
				Commit:  "1234ABCD",
				Message: "merging my branch",
			},
			"--continue",
		},
		"merge --abort": {
			&MergeOpts{
				Abort: true,
				// Options are ignored
				Commit:  "1234ABCD",
				Message: "merging my branch",
			},
			"--abort",
		},
		"merge --quit": {
			&MergeOpts{
				Quit: true,
				// Options are ignored
				Commit:  "1234ABCD",
				Message: "merging my branch",
			},
			"--quit",
		},
		"pull 1/3 opts": {
			&PullOpts{
				Atomic:        true,
				Autostash:     true,
				Depth:         4,
				DoCommit:      true,
				FF:            true,
				GPGSign:       true,
				NoLog:         true,
				NoStat:        true,
				Quiet:         true,
				Prune:         true,
				SetUpstream:   true,
				Squash:        true,
				Rebase:        RebaseStrategyTrue,
				Refspec:       []string{"my-branch"},
				Repository:    "my-repo",
				UpdateShallow: true,
			}, "--atomic --autostash --commit --depth 4 --ff --gpg-sign --squash --no-log --no-stat --no-stat --prune --quiet --rebase=true --set-upstream my-repo my-branch",
		},
		"pull 2/3 opts": {
			&PullOpts{
				AllowUnrelatedHistories: true,
				Append:                  true,
				Deepen:                  3,
				FFOnly:                  true,
				GPGSignKeyID:            "4321DCBA",
				Log:                     5,
				NoRebase:                true,
				NoRecurseSubmodules:     true,
				Porcelain:               true,
				Progress:                true,
				PruneTags:               true,
				Refspec:                 []string{"my-branch"},
				Repository:              "my-repo",
				Stat:                    true,
				Strategy:                MergeStrategyOctopus,
				StrategyOptions: []MergeStrategyOption{
					MergeStrategyOptionDiffAlgorithmMinimal,
					MergeStrategyOptionIgnoreCRAtEOL,
				},
				Verbose: true,
				Verify:  true,
			}, "--deepen 3 --ff-only --gpg-sign=4321DCBA --log=5 --stat -X diff-algorithm=minimal -X ignore-cr-at-eol --no-rebase --no-recurse-submodules --porcelain --progress --prune-tags --verbose my-repo my-branch",
		},
		"pull 3/3 opts": {
			&PullOpts{
				All:               true,
				Cleanup:           CommitCleanupModeDefault,
				Force:             true,
				NoAutostash:       true,
				NoDoCommit:        true,
				NoFF:              true,
				NoSquash:          true,
				NoTags:            true,
				NoVerify:          true,
				RecurseSubmodules: RecurseSubmodulesOnDemand,
				Refspec:           []string{"my-branch"},
				Repository:        "my-repo",
				Stat:              true,
				Unshallow:         true,
			}, "--all --force --stat --no-autostash --no-commit --no-ff --no-squash --no-tags --no-verify --unshallow my-repo my-branch",
		},
		"push 1/2 opts": {
			&PushOpts{
				All:                 true,
				Branches:            true,
				DryRun:              true,
				FollowTags:          true,
				ForceIfIncludes:     true,
				Mirror:              true,
				NoForceWithLease:    true,
				NoRecurseSubmodules: true,
				NoThin:              true,
				Porcelain:           true,
				Prune:               true,
				Quiet:               true,
				Refspec:             []string{"my-branch"},
				Repository:          "my-repo",
				SetUpstream:         true,
				Tags:                true,
				Verify:              true,
			}, "--all --branches --dry-run --follow-tags --force-if-includes --mirror --no-force-with-lease --no-recurse-submodules --no-thin --porcelain --prune --quiet --set-upstream --tags --verify my-repo my-branch",
		},
		"push 2/2 opts": {
			&PushOpts{
				Atomic:            true,
				Delete:            true,
				Exec:              "/path/to/git-receive-pack",
				Force:             true,
				ForceWithLease:    "another-branch",
				NoAtomic:          true,
				NoForceIfIncludes: true,
				NoSigned:          true,
				NoVerify:          true,
				Progress:          true,
				PushOption:        "a",
				RecurseSubmodules: PushRecurseSubmodulesCheck,
				Refspec:           []string{"my-branch"},
				Repository:        "my-repo",
				Signed:            PushSignedTrue,
				Verbose:           true,
			}, "--atomic --delete --exec=/path/to/git-receive-pack --force --force-with-lease=another-branch --no-atomic --no-force-if-includes --no-signed --no-verify --progress --push-option=a --recurse-submodules=check --signed=true --verbose my-repo my-branch",
		},
		"rebase 1/3 opts": {
			&RebaseOpts{
				AllowEmptyMessage:    true,
				Autosquash:           true,
				Branch:               "my-branch",
				Empty:                EmptyCommitDrop,
				ForkPoint:            true,
				IgnoreDate:           true,
				Merge:                true,
				NoAutostash:          true,
				NoRebaseMerges:       true,
				NoReReReAutoupdate:   true,
				NoVerify:             true,
				RescheduleFailedExec: true,
				Root:                 true,
				UpdateRefs:           true,
			}, "--allow-empty-message --autosquash --empty=drop --fork-point --ignore-date --merge --no-autostash --no-rebase-merges --no-rerere-autoupdate --no-verify --reschedule-failed-exec --root --update-refs my-branch",
		},
		"rebase 2/3 opts": {
			&RebaseOpts{
				Apply:                     true,
				Branch:                    "my-branch",
				CommitterDateIsAuthorDate: true,
				Exec:                      "/path/to/git-receive-pack",
				GPGSign:                   true,
				IgnoreWhitespace:          true,
				KeepEmpty:                 true,
				NoReapplyCherryPicks:      true,
				NoStat:                    true,
				Onto:                      "new-base",
				Quiet:                     true,
				ResetAuthorDate:           true,
				Stat:                      true,
				Whitespace:                WhitespaceActionFix,
			}, "--apply --committer-date-is-author-date --exec=/path/to/git-receive-pack --gpg-sign --ignore-whitespace --keep-empty --no-reapply-cherry-picks --no-stat --onto=new-base --quiet --reset-author-date --stat --whitespace=fix my-branch",
		},
		"rebase 3/3 opts": {
			&RebaseOpts{
				Autostash:              true,
				Branch:                 "my-branch",
				Context:                3,
				ForceRebase:            true,
				GPGSignKeyID:           "4321DCBA",
				KeepBase:               "another-upstream",
				NoAutosquash:           true,
				NoKeepEmpty:            true,
				NoRescheduleFailedExec: true,
				NoUpdateRefs:           true,
				ReapplyCherryPicks:     true,
				RebaseMerges:           RebaseMergesCousins,
				ReReReAutoupdate:       true,
				Strategy:               MergeStrategySubtree,
				StrategyOptions: []MergeStrategyOption{
					MergeStrategyOptionDiffAlgorithmPatience,
					MergeStrategyOptionNoRenormalize,
				},
				Verbose: true,
			}, "--autostash -C 3 --force-rebase --gpg-sign=4321DCBA --keep-base=another-upstream --no-autosquash --no-keep-empty --no-reschedule-failed-exec --no-update-refs --reapply-cherry-picks --rebase-merges=rebase-cousins --rerere-autoupdate --strategy=subtree --strategy-option=diff-algorithm=patience --strategy-option=no-renormalize --verbose my-branch",
		},
		"rebase --continue": {
			&RebaseOpts{
				Continue: true,
				// Options are ignored
				Branch:       "my-branch",
				GPGSignKeyID: "4321DCBA",
			},
			"--continue",
		},
		"rebase --abort": {
			&RebaseOpts{
				Abort: true,
				// Options are ignored
				Branch:       "my-branch",
				GPGSignKeyID: "4321DCBA",
			},
			"--abort",
		},
		"rebase --quit": {
			&RebaseOpts{
				Quit: true,
				// Options are ignored
				Branch:       "my-branch",
				GPGSignKeyID: "4321DCBA",
			},
			"--quit",
		},
		"rebase --skip": {
			&RebaseOpts{
				Skip: true,
				// Options are ignored
				Branch:       "my-branch",
				GPGSignKeyID: "4321DCBA",
			},
			"--skip",
		},
		"rebase --show-current-patch": {
			&RebaseOpts{
				ShowCurrentPatch: true,
				// Options are ignored
				Branch:       "my-branch",
				GPGSignKeyID: "4321DCBA",
			},
			"--show-current-patch",
		},
		"remote": {
			&RemoteOpts{
				Verbose: true,
			},
			"-v",
		},
		"remote add": {
			&RemoteOpts{
				Command: RemoteCommandAdd,
				Track:   []string{"branch-a", "branch-b"},
				Master:  "main",
				Fetch:   true,
				Tags:    true,
				NoTags:  true,
				Name:    "origin",
				URL:     "git@github.com:hashicorp/vault.git",
			},
			"add -f --tags --no-tags -m main -t branch-a -t branch-b origin git@github.com:hashicorp/vault.git",
		},
		"remote rename": {
			&RemoteOpts{
				Command:    RemoteCommandRename,
				Progress:   true,
				NoProgress: true,
				Old:        "old-remote",
				New:        "new-remote",
			},
			"rename --progress --no-progress old-remote new-remote",
		},
		"remote remove": {
			&RemoteOpts{
				Command: RemoteCommandRemove,
				Name:    "remote",
			},
			"remove remote",
		},
		"remote set-head": {
			&RemoteOpts{
				Command: RemoteCommandSetHead,
				Name:    "my-remote",
				Auto:    true,
				Delete:  true,
				Branch:  "my-branch",
			},
			"set-head my-remote --auto --delete my-branch",
		},
		"remote set-branches": {
			&RemoteOpts{
				Command:  RemoteCommandSetBranches,
				Add:      true,
				Name:     "my-remote",
				Branch:   "branch-a",
				Branches: []string{"branch-b", "branch-c"},
			},
			"set-branches --add my-remote branch-a branch-b branch-c",
		},
		"remote get-url": {
			&RemoteOpts{
				Command: RemoteCommandGetURL,
				Push:    true,
				All:     true,
				Name:    "my-remote",
			},
			"get-url --push --all my-remote",
		},
		"remote set-url oldurl": {
			&RemoteOpts{
				Command: RemoteCommandSetURL,
				Push:    true,
				Name:    "my-remote",
				NewURL:  "git@github.com:hashicorp/vault-enterprise.git",
				OldURL:  "git@github.com:hashicorp/vault.git",
			},
			"set-url --push my-remote git@github.com:hashicorp/vault-enterprise.git git@github.com:hashicorp/vault.git",
		},
		"remote set-url add": {
			&RemoteOpts{
				Command: RemoteCommandSetURL,
				Add:     true,
				Push:    true,
				Name:    "my-remote",
				NewURL:  "git@github.com:hashicorp/vault-enterprise.git",
			},
			"set-url --add --push my-remote git@github.com:hashicorp/vault-enterprise.git",
		},
		"remote set-url delete": {
			&RemoteOpts{
				Command: RemoteCommandSetURL,
				Delete:  true,
				Push:    true,
				Name:    "my-remote",
				URL:     "git@github.com:hashicorp/vault-enterprise.git",
			},
			"set-url --delete --push my-remote git@github.com:hashicorp/vault-enterprise.git",
		},
		"remote show": {
			&RemoteOpts{
				Command: RemoteCommandShow,
				Verbose: true,
				NoQuery: true,
				Name:    "my-remote",
			},
			"-v show -n my-remote",
		},
		"remote prune": {
			&RemoteOpts{
				Command: RemoteCommandPrune,
				NoQuery: true,
				DryRun:  true,
				Name:    "branch-a",
				Names:   []string{"branch-b", "branch-c"},
			},
			"prune -n --dry-run branch-a branch-b branch-c",
		},
		"remote update": {
			&RemoteOpts{
				Command: RemoteCommandUpdate,
				Verbose: true,
				Prune:   true,
			},
			"-v update --prune",
		},
		"reset": {
			&ResetOpts{
				Mode:      ResetModeHard,
				NoRefresh: true,
				Patch:     true,
				Quiet:     true,
				Refresh:   true,
				Commit:    "abcd1234",
				Treeish:   "HEAD~2",
				PathSpec:  []string{"vault/something_ent.go", "vault/cli/another_ent.go"},
			},
			"--hard --no-refresh --quiet --refresh --patch abcd1234 HEAD~2 -- vault/something_ent.go vault/cli/another_ent.go",
		},
		"rm": {
			&RmOpts{
				Cached:          true,
				DryRun:          true,
				Force:           true,
				IgnoreUnmatched: true,
				Quiet:           true,
				Recursive:       true,
				Sparse:          true,
				PathSpec:        []string{"vault/something_ent.go", "vault/cli/another_ent.go"},
			},
			"--cached --dry-run --force --ignore-unmatched --quiet -r --sparse -- vault/something_ent.go vault/cli/another_ent.go",
		},
		"show": {
			&ShowOpts{
				DiffAlgorithm: DiffAlgorithmHistogram,
				DiffMerges:    DiffMergeFormatDenseCombined,
				Format:        "medium",
				NoColor:       true,
				NoPatch:       true,
				Output:        "/path/to/my.diff",
				Patch:         true,
				Raw:           true,
				Object:        "HEAD",
				PathSpec:      []string{"go.mod", "go.sum"},
			},
			"--diff-algorithm=histogram --diff-merges=dense-combined --format=medium --no-color --no-patch --output=/path/to/my.diff --patch --raw HEAD -- go.mod go.sum",
		},
		"status": {
			&StatusOpts{
				AheadBehind:      true,
				Branch:           true,
				Column:           "always",
				FindRenames:      12,
				Ignored:          IgnoredModeMatching,
				IgnoreSubmodules: IgnoreSubmodulesWhenDirty,
				Long:             true,
				NoAheadBehind:    true,
				NoColumn:         true,
				NoRenames:        true,
				Porcelain:        true,
				Renames:          true,
				Short:            true,
				ShowStash:        true,
				UntrackedFiles:   UntrackedFilesAll,
				Verbose:          true,
				PathSpec:         []string{"go.mod", "go.sum"},
			},
			"--ahead-behind --branch --column=always --find-renames=12 --ignored=matching --ignore-submodules=dirty --long --no-ahead-behind --no-column --no-renames --porcelain --renames --short --show-stash --untracked-files=all --verbose -- go.mod go.sum",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, expect.expected, expect.opts.String())
		})
	}
}
