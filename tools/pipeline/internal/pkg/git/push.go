// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package git

import (
	"context"
	"fmt"
	"strings"
)

type (
	// PushRecurseSubmodules is how to handle recursive sub-modules
	PushRecurseSubmodules = string
	// PushSigned sets gpg sign mode
	PushSigned = string
)

const (
	PushRecurseSubmodulesCheck    PushRecurseSubmodules = "check"
	PushRecurseSubmodulesOnDemand PushRecurseSubmodules = "on-demand"
	PushRecurseSubmodulesOnly     PushRecurseSubmodules = "only"
	PushRecurseSubmodulesNo       PushRecurseSubmodules = "no"

	PushSignedTrue    PushSigned = "true"
	PushSignedFalse   PushSigned = "false"
	PushSignedIfAsked PushSigned = "if-asked"
)

// PushOpts are the git push flags and arguments
// See: https://git-scm.com/docs/git-push
type PushOpts struct {
	// Options
	All                 bool                  // --all
	Atomic              bool                  // --atomic
	Branches            bool                  // --branches
	Delete              bool                  //  --delete
	DryRun              bool                  // --dry-run
	Exec                string                // --exec=<git-receive-pack>
	FollowTags          bool                  //  --follow-tags
	Force               bool                  // --force
	ForceIfIncludes     bool                  // --force-if-includes
	ForceWithLease      string                // --force-with-lease=<refname>
	Mirror              bool                  // --mirror
	NoAtomic            bool                  // --no-atomic
	NoForceIfIncludes   bool                  // --no-force-if-includes
	NoForceWithLease    bool                  // --no-force-with-lease
	NoRecurseSubmodules bool                  // --no-recurse-submodules
	NoSigned            bool                  //  --no-signed
	NoThin              bool                  // --no-thin
	NoVerify            bool                  // --no-verify
	Porcelain           bool                  // --porcelain
	Progress            bool                  // --progress
	Prune               bool                  // --prune
	PushOption          string                // --push-option
	Quiet               bool                  // --quiet
	RecurseSubmodules   PushRecurseSubmodules // --recurse-submodules=<mode>
	SetUpstream         bool                  // --set-upstream
	Signed              PushSigned            //  --signed=<mode>
	Tags                bool                  //  --tags
	Thin                bool                  // --thin
	Verbose             bool                  // --verbose
	Verify              bool                  // --verify

	// Targets
	Repository string   // <repository>
	Refspec    []string // <refspec>
}

// Push runs the git push command
func (c *Client) Push(ctx context.Context, opts *PushOpts) (*ExecResponse, error) {
	return c.Exec(ctx, "push", opts)
}

// String returns the options as a string
func (o *PushOpts) String() string {
	return strings.Join(o.Strings(), " ")
}

// Strings returns the options as a string slice
func (o *PushOpts) Strings() []string {
	if o == nil {
		return nil
	}

	opts := []string{}

	if o.All {
		opts = append(opts, "--all")
	}

	if o.Atomic {
		opts = append(opts, "--atomic")
	}

	if o.Branches {
		opts = append(opts, "--branches")
	}

	if o.Delete {
		opts = append(opts, "--delete")
	}

	if o.DryRun {
		opts = append(opts, "--dry-run")
	}

	if o.Exec != "" {
		opts = append(opts, fmt.Sprintf("--exec=%s", o.Exec))
	}

	if o.FollowTags {
		opts = append(opts, "--follow-tags")
	}

	if o.Force {
		opts = append(opts, "--force")
	}

	if o.ForceIfIncludes {
		opts = append(opts, "--force-if-includes")
	}

	if o.ForceWithLease != "" {
		opts = append(opts, fmt.Sprintf("--force-with-lease=%s", o.ForceWithLease))
	}

	if o.Mirror {
		opts = append(opts, "--mirror")
	}

	if o.NoAtomic {
		opts = append(opts, "--no-atomic")
	}

	if o.NoForceIfIncludes {
		opts = append(opts, "--no-force-if-includes")
	}

	if o.NoForceWithLease {
		opts = append(opts, "--no-force-with-lease")
	}

	if o.NoRecurseSubmodules {
		opts = append(opts, "--no-recurse-submodules")
	}

	if o.NoSigned {
		opts = append(opts, "--no-signed")
	}

	if o.NoThin {
		opts = append(opts, "--no-thin")
	}

	if o.NoVerify {
		opts = append(opts, "--no-verify")
	}

	if o.Porcelain {
		opts = append(opts, "--porcelain")
	}

	if o.Progress {
		opts = append(opts, "--progress")
	}

	if o.Prune {
		opts = append(opts, "--prune")
	}

	if o.PushOption != "" {
		opts = append(opts, fmt.Sprintf("--push-option=%s", o.PushOption))
	}

	if o.Quiet {
		opts = append(opts, "--quiet")
	}

	if o.RecurseSubmodules != "" {
		opts = append(opts, fmt.Sprintf("--recurse-submodules=%s", string(o.RecurseSubmodules)))
	}

	if o.SetUpstream {
		opts = append(opts, "--set-upstream")
	}

	if o.Signed != "" {
		opts = append(opts, fmt.Sprintf("--signed=%s", string(o.Signed)))
	}

	if o.Tags {
		opts = append(opts, "--tags")
	}

	if o.Thin {
		opts = append(opts, "--thin")
	}

	if o.Verbose {
		opts = append(opts, "--verbose")
	}

	if o.Verify {
		opts = append(opts, "--verify")
	}

	if o.Repository != "" {
		opts = append(opts, o.Repository)
	}

	if len(o.Refspec) > 0 {
		opts = append(opts, o.Refspec...)
	}

	return opts
}
