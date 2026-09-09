// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"slices"
	"strings"

	gh "github.com/google/go-github/v83/github"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/releases"
	"github.com/jedib0t/go-pretty/v6/table"
)

// CheckWorkflowRunReq holds the request parameters for checking a workflow run.
type CheckWorkflowRunReq struct {
	// Workflow
	Owner    string
	Repo     string
	Workflow string
	MaxDepth uint

	// Branches branches
	Patterns         []string
	compiledPatterns []*regexp.Regexp
	Branches         []string
	Releases         *releases.DecodeRes
	OnlyLTS          bool
	OnlyCE           bool
	IncludeMain      bool
	CEBranchPrefix   string // Prefix for CE branches in enterprise repos (e.g., "ce" for "ce/release/<version>")

	// Output opts
	HideSuccessfulJobs  bool
	WriteToGithubOutput bool
	GithubOutputKey     string
}

// CheckWorkflowRunRes represents the response from checking a workflow run.
type CheckWorkflowRunRes struct {
	Success            bool                            `json:"success"`
	Results            []*CheckWorkflowRunBranchResult `json:"results"`
	Owner              string                          `json:"owner,omitempty"`
	Repo               string                          `json:"repo,omitempty"`
	HideSuccessfulJobs bool                            `json:"-"` // Not included in JSON output
	GithubOutputKey    string                          `json:"-"`
}

type CheckWorkflowRunGithubOutput struct {
	Success       bool                            `json:"success,omitempty"`
	Results       []*CheckWorkflowRunBranchResult `json:"results,omitempty"`
	BranchesCount int                             `json:"branch_count,omitempty"`
	FailedCount   int                             `json:"failed_count,omitempty"`
}

// CheckWorkflowRunBranchResult represents the result of checking a single branch.
type CheckWorkflowRunBranchResult struct {
	Branch          string                          `json:"branch,omitempty"`
	Success         bool                            `json:"success,omitempty"`
	Depth           int                             `json:"depth,omitempty"`
	WorkflowRunID   int64                           `json:"workflow_run_id,omitempty"`
	Status          string                          `json:"status,omitempty"`
	Conclusion      string                          `json:"conclusion,omitempty"`
	MatchedJobs     []*CheckWorkflowRunJobMatch     `json:"matched_jobs,omitempty"`
	MissingPatterns []string                        `json:"missing_patterns,omitempty"`
	FailedJobs      []*CheckWorkflowRunJobMatch     `json:"failed_jobs,omitempty"`
	FailureDetails  *CheckWorkflowRunFailureDetails `json:"failure_details,omitempty"`
	Error           string                          `json:"error,omitempty"`
	WorkflowURL     string                          `json:"workflow_url,omitempty"`
	CommitSHA       string                          `json:"commit_sha,omitempty"`
}

// CheckWorkflowRunJobMatch represents a job that matched a pattern.
type CheckWorkflowRunJobMatch struct {
	Name       string `json:"name"`
	Pattern    string `json:"pattern"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
}

// CheckWorkflowRunFailureDetails provides details about why a workflow run check failed.
type CheckWorkflowRunFailureDetails struct {
	Reason          string                      `json:"reason"`
	MissingPatterns []string                    `json:"missing_patterns,omitempty"`
	FailedJobs      []*CheckWorkflowRunJobMatch `json:"failed_jobs,omitempty"`
	TaintedJobs     []*CheckWorkflowRunJobMatch `json:"tainted_jobs,omitempty"`
}

// CheckWorkflowRunEvaluation holds the result of evaluating a workflow run.
type CheckWorkflowRunEvaluation struct {
	IsComplete      bool // All required patterns matched
	IsSuccessful    bool // All matched jobs succeeded
	IsTainted       bool // Has failed jobs (even non-matched)
	Depth           int  // 1-based depth at which the run was found
	MatchedJobs     []*CheckWorkflowRunJobMatch
	MissingPatterns []string
	FailedJobs      []*CheckWorkflowRunJobMatch // Matched jobs that failed
	TaintedJobs     []*CheckWorkflowRunJobMatch // Non-matched jobs that failed
}

// Run executes the check workflow run operation.
func (r *CheckWorkflowRunReq) Run(ctx context.Context, client *gh.Client) (*CheckWorkflowRunRes, error) {
	slog.Default().DebugContext(
		ctx, "starting check workflow run operation",
		slog.String("owner", r.Owner),
		slog.String("repo", r.Repo),
		slog.String("workflow", r.Workflow),
		slog.Int("patterns", len(r.Patterns)),
	)

	if err := r.validate(ctx); err != nil {
		return nil, err
	}

	branches, err := r.determineBranches(ctx)
	if err != nil {
		return nil, fmt.Errorf("determining branches: %w", err)
	}

	workflow, err := getWorkflow(ctx, client, r.Owner, r.Repo, r.Workflow)
	if err != nil {
		return nil, fmt.Errorf("finding workflow: %w", err)
	}
	workflowID := workflow.GetID()

	slog.Default().DebugContext(ctx, "found workflow",
		slog.Group("workflow",
			slog.String("name", r.Workflow),
			slog.Int64("id", workflowID),
		),
	)

	// Check each branch for matching workflow runs
	branchResults := make([]*CheckWorkflowRunBranchResult, 0, len(branches))
	overallSuccess := true

	for _, branch := range branches {
		slog.Default().DebugContext(
			ctx, "checking branch",
			slog.String("branch", branch),
		)

		run, eval, err := r.findWorkflowRunWithJobs(
			ctx,
			client,
			workflowID,
			branch,
		)
		if err != nil {
			slog.Default().DebugContext(
				ctx, "no matching workflow run found for branch",
				slog.String("branch", branch),
				slog.String("error", err.Error()),
			)
			branchResults = append(branchResults, &CheckWorkflowRunBranchResult{
				Branch:  branch,
				Success: false,
				Error:   err.Error(),
			})
			overallSuccess = false
			continue
		}

		// Check if no runs were found (nil run/eval but no error)
		if run == nil || eval == nil {
			slog.Default().DebugContext(
				ctx, "no workflow runs found for branch",
				slog.String("branch", branch),
			)
			branchResults = append(branchResults, &CheckWorkflowRunBranchResult{
				Branch:  branch,
				Success: false,
				Error:   "no workflow runs found for branch",
			})
			overallSuccess = false
			continue
		}

		// Build branch result based on evaluation
		success := eval.IsComplete && eval.IsSuccessful && !eval.IsTainted
		branchResult := r.buildBranchResult(branch, run, eval, success)
		branchResults = append(branchResults, branchResult)

		if !success {
			overallSuccess = false
		}

		slog.Default().DebugContext(
			ctx, "workflow run check completed for branch",
			slog.String("branch", branch),
			slog.Bool("success", success),
			slog.Int64("run_id", run.GetID()),
		)
	}

	slices.SortStableFunc(branchResults, func(a, b *CheckWorkflowRunBranchResult) int {
		return cmp.Compare(a.Branch, b.Branch)
	})

	if len(branchResults) == 0 {
		return nil, fmt.Errorf("no branches to check")
	}

	return &CheckWorkflowRunRes{
		Success:            overallSuccess,
		Results:            branchResults,
		Owner:              r.Owner,
		Repo:               r.Repo,
		HideSuccessfulJobs: r.HideSuccessfulJobs,
		GithubOutputKey:    r.GithubOutputKey,
	}, nil
}

// ToTable marshals the response to a text table.
func (r *CheckWorkflowRunRes) ToTable() string {
	if r == nil {
		return ""
	}

	var output strings.Builder

	output.WriteString("Branch Summary:\n")
	output.WriteString(r.buildSummaryTable().Render())

	// Show Matched Jobs section if: we have failures OR we're not hiding successful jobs
	if r.hasMatchedJobs() && (r.hasFailedJobs() || !r.HideSuccessfulJobs) {
		output.WriteString("\n\nMatched Jobs:\n")
		output.WriteString(r.buildMatchedJobsTable(r.HideSuccessfulJobs).Render())
	}

	if r.hasFailedJobs() {
		output.WriteString("\n\nFailed Jobs:\n")
		output.WriteString(r.buildFailedJobsTable().Render())
	}

	if r.hasTaintedJobs() {
		output.WriteString("\n\nTainted Jobs (non-matched jobs that failed):\n")
		output.WriteString(r.buildTaintedJobsTable().Render())
	}

	if r.hasMissingPatterns() {
		output.WriteString("\n\nMissing Patterns:\n")
		output.WriteString(r.buildMissingPatternsTable().Render())
	}

	return output.String()
}

// ToMarkdown marshals the response to markdown format using table.Writer.
func (r *CheckWorkflowRunRes) ToMarkdown() string {
	if r == nil {
		return ""
	}

	var output strings.Builder

	output.WriteString("## Branch Summary\n\n")
	output.WriteString(r.buildSummaryTable().RenderMarkdown())

	// Show Matched Jobs section if: we have failures OR we're not hiding successful jobs
	if r.hasMatchedJobs() && (r.hasFailedJobs() || !r.HideSuccessfulJobs) {
		output.WriteString("\n\n## Matched Jobs\n\n")
		output.WriteString(r.buildMatchedJobsTable(r.HideSuccessfulJobs).RenderMarkdown())
	}

	if r.hasFailedJobs() {
		output.WriteString("\n\n## Failed Jobs\n\n")
		output.WriteString(r.buildFailedJobsTable().RenderMarkdown())
	}

	if r.hasTaintedJobs() {
		output.WriteString("\n\n## Tainted Jobs\n\n")
		output.WriteString("Non-matched jobs that failed and may have prevented required jobs from running.\n\n")
		output.WriteString(r.buildTaintedJobsTable().RenderMarkdown())
	}

	if r.hasMissingPatterns() {
		output.WriteString("\n\n## Missing Patterns\n\n")
		output.WriteString(r.buildMissingPatternsTable().RenderMarkdown())
	}

	return output.String()
}

// ToGithubOutput writes a simplified result for $GITHUB_OUTPUT.
func (r *CheckWorkflowRunRes) ToGithubOutput() ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized")
	}

	failedCount := 0
	for _, result := range r.Results {
		if !result.Success {
			failedCount++
		}
	}
	// TODO: Should we remove the successful jobs from github output if
	// HideSuccessfulJobs is true?
	output := &CheckWorkflowRunGithubOutput{
		Success:       r.Success,
		Results:       r.Results,
		BranchesCount: len(r.Results),
		FailedCount:   failedCount,
	}

	b, err := json.Marshal(output)
	if err != nil {
		return nil, fmt.Errorf("marshaling check workflow run GITHUB_OUTPUT to JSON: %w", err)
	}

	return b, nil
}

// String returns a human-readable representation.
func (r *CheckWorkflowRunRes) String() string {
	if r == nil {
		return "No result"
	}

	if r.Success {
		return fmt.Sprintf("All %d branches passed workflow run checks", len(r.Results))
	}

	failedCount := 0
	for _, br := range r.Results {
		if !br.Success {
			failedCount++
		}
	}

	return fmt.Sprintf("%d of %d branches failed workflow run checks", failedCount, len(r.Results))
}

// ToJSON marshals the response to JSON.
func (r *CheckWorkflowRunRes) ToJSON() ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized")
	}

	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling check workflow run response to JSON: %w", err)
	}

	return b, nil
}

// validate checks that the request is valid.
func (r *CheckWorkflowRunReq) validate(ctx context.Context) error {
	if r == nil {
		return errors.New("uninitialized")
	}

	if r.Owner == "" {
		return errors.New("owner is required")
	}

	if r.Repo == "" {
		return errors.New("repo is required")
	}

	if r.Workflow == "" {
		return errors.New("workflow is required")
	}

	// Branches can be empty if Releases is provided
	if len(r.Branches) == 0 && r.Releases == nil {
		return errors.New("either branches or releases configuration is required")
	}

	// If Releases is provided, validate it
	if r.Releases != nil {
		if err := r.Releases.Validate(ctx); err != nil {
			return fmt.Errorf("invalid releases configuration: %w", err)
		}
	}

	if len(r.Patterns) == 0 {
		return errors.New("at least one pattern is required")
	}

	// Validate that patterns are valid regex
	var err error
	r.compiledPatterns, err = r.compilePatterns()
	if err != nil {
		return err
	}

	if r.WriteToGithubOutput {
		if r.GithubOutputKey == "" {
			return errors.New("github output key required when writing to $GITHUB_OUTPUT")
		}
	}

	return nil
}

// findWorkflowRunWithJobs finds a workflow run that matches all required patterns.
// It searches up to 50 workflow runs to find one where all patterns match successfully.
// Returns the workflow run, evaluation, and error.
func (r *CheckWorkflowRunReq) findWorkflowRunWithJobs(
	ctx context.Context,
	client *gh.Client,
	workflowID int64,
	branch string,
) (*gh.WorkflowRun, *CheckWorkflowRunEvaluation, error) {
	slog.Default().DebugContext(
		ctx, "finding matching workflow run",
		slog.String("owner", r.Owner),
		slog.String("repo", r.Repo),
		slog.Int64("workflow_id", workflowID),
		slog.String("branch", branch),
		slog.Uint64("max_run_depth", uint64(r.MaxDepth)),
	)

	// Try each branch variation
	opts := &gh.ListWorkflowRunsOptions{
		Branch:      branch,
		ListOptions: gh.ListOptions{PerPage: PerPageMax},
	}

	runs, err := getWorkflowRuns(ctx, client, r.Owner, r.Repo, workflowID, int(r.MaxDepth), opts)
	if err != nil {
		slog.Default().DebugContext(
			ctx, "error listing runs for branch",
			slog.String("branch", branch),
			slog.String("error", err.Error()),
		)

		return nil, nil, err
	}

	if len(runs) == 0 {
		slog.Default().DebugContext(
			ctx, "no runs found for branch variation",
			slog.String("branch", branch),
		)

		return nil, nil, nil
	}

	slog.Default().DebugContext(
		ctx, "checking workflow runs for pattern matches",
		slog.String("branch", branch),
		slog.Int("run_count", len(runs)),
	)

	// Check each run to find one where all patterns match
	for i, wrun := range runs {
		if wrun.Run == nil {
			slog.Default().DebugContext(
				ctx, "workflow run empty",
				slog.Int("run_index", i+1),
			)
			continue
		}
		run := wrun.Run

		slog.Default().DebugContext(
			ctx, "checking workflow run",
			slog.Int("run_index", i+1),
			slog.Int("total_runs", len(runs)),
			slog.Int64("run_id", run.GetID()),
			slog.String("status", run.GetStatus()),
			slog.String("conclusion", run.GetConclusion()),
		)

		// Skip runs that haven't completed yet
		if run.GetStatus() != "completed" {
			slog.Default().DebugContext(ctx, "skipping non-completed run",
				slog.Int64("run_id", run.GetID()),
				slog.String("status", run.GetStatus()),
			)
			continue
		}

		// Fetch jobs for this run
		jobs, err := getWorkflowJobsForRun(ctx, client, r.Owner, r.Repo, run.GetID())
		if err != nil {
			slog.Default().DebugContext(
				ctx, "error fetching jobs for run",
				slog.Int64("run_id", run.GetID()),
				slog.String("error", err.Error()),
			)
			continue
		}

		// Evaluate the workflow run
		eval := r.evaluateWorkflowRun(jobs)
		eval.Depth = i + 1 // 1-based depth

		// Check if all matched jobs were skipped
		allSkipped := len(eval.MatchedJobs) > 0
		for _, job := range eval.MatchedJobs {
			if job.Conclusion != "skipped" {
				allSkipped = false
				break
			}
		}

		// Skip runs where all matched jobs were skipped - continue searching
		if allSkipped {
			slog.Default().DebugContext(ctx, "all matched jobs skipped, continuing search",
				slog.Int64("run_id", run.GetID()),
				slog.Int("matched_jobs", len(eval.MatchedJobs)),
			)
			continue
		}

		// Case 1: Complete & Successful - return immediately
		if eval.IsComplete && eval.IsSuccessful {
			slog.Default().DebugContext(ctx, "found complete successful run",
				slog.Group("run",
					slog.String("branch", branch),
					slog.Int64("id", run.GetID()),
					slog.Int("depth", eval.Depth),
				),
				slog.Group("evaluation",
					slog.Int("matched_jobs", len(eval.MatchedJobs)),
				),
			)
			return run, eval, nil
		}

		// Case 2: Complete & Failed - return immediately (stop searching)
		if eval.IsComplete && !eval.IsSuccessful {
			slog.Default().DebugContext(ctx, "found complete failed run",
				slog.Group("run",
					slog.String("branch", branch),
					slog.Int64("id", run.GetID()),
					slog.Int("depth", eval.Depth),
				),
				slog.Group("evaluation",
					slog.Int("failed_jobs", len(eval.FailedJobs)),
				),
			)
			return run, eval, nil
		}

		// Case 3: Incomplete & Tainted - return immediately (stop searching)
		if !eval.IsComplete && eval.IsTainted {
			slog.Default().DebugContext(ctx, "found incomplete tainted run",
				slog.Group("run",
					slog.String("branch", branch),
					slog.Int64("id", run.GetID()),
					slog.Int("depth", eval.Depth),
				),
				slog.Group("evaluation",
					slog.Int("tainted_jobs", len(eval.TaintedJobs)),
					slog.Int("missing_patterns", len(eval.MissingPatterns)),
				),
			)
			return run, eval, nil
		}

		// Case 4: Incomplete & Clean - continue to next run
		slog.Default().DebugContext(ctx, "incomplete clean run, continuing search",
			slog.Group("run",
				slog.Int64("id", run.GetID()),
			),
			slog.Group("evaluation",
				slog.Int("missing_patterns", len(eval.MissingPatterns)),
			),
		)
	}

	return nil, nil, fmt.Errorf("no suitable workflow runs found in %d runs for branch: %v", len(runs), branch)
}

// determineBranches returns the list of branches to check.
// If Branches is provided, uses those. Otherwise, extracts active branches from Releases.
func (r *CheckWorkflowRunReq) determineBranches(ctx context.Context) ([]string, error) {
	// If branches are explicitly provided, use them
	if len(r.Branches) > 0 {
		return r.Branches, nil
	}

	// Extract active branches from releases configuration
	if r.Releases == nil || r.Releases.Config == nil || r.Releases.Config.ActiveVersion == nil {
		return nil, errors.New("no branches provided and no releases configuration available")
	}

	// Determine if we're working with an enterprise repository
	// Enterprise repos have "-enterprise" suffix (e.g., vault-enterprise, consul-enterprise)
	isEntRepo := releases.IsEnterpriseRepo(r.Repo)

	branches := []string{}

	// Handle --include-main flag
	if r.IncludeMain {
		if r.OnlyCE && isEntRepo && r.CEBranchPrefix != "" {
			// CE-only mode with prefix: only add CE main
			branches = append(branches, fmt.Sprintf("%s/main", r.CEBranchPrefix))
		} else if r.OnlyCE {
			// CE-only mode without prefix: only add main
			branches = append(branches, "main")
		} else if isEntRepo && r.CEBranchPrefix != "" {
			// ENT repo with CE prefix: add both main and ce/main
			branches = append(branches, "main")
			branches = append(branches, fmt.Sprintf("%s/main", r.CEBranchPrefix))
		} else {
			// Default: just add main
			branches = append(branches, "main")
		}
	}

	for version, versionInfo := range r.Releases.Config.ActiveVersion.Versions {
		// Apply filters
		if r.OnlyLTS && !versionInfo.LTS {
			continue
		}
		if r.OnlyCE && !versionInfo.CEActive {
			continue
		}

		// Determine branch name based on repository type, filters, and CE active status
		// Decision tree:
		// 1. If --only-ce flag: Always use CE branch format (with prefix in ENT repos)
		// 2. If ENT repo + CEActive + prefix configured: Use CE branch format with prefix
		// 3. If ENT repo: Use ENT branch format (release/<version>+ent)
		// 4. If CE repo: Use CE branch format (release/<version>)
		var branchName string
		if r.OnlyCE {
			// --only-ce flag: Use CE branch format
			// In ENT repos with prefix: ce/release/<version>
			// Otherwise: release/<version>
			if isEntRepo && r.CEBranchPrefix != "" {
				branchName = releases.CEReleaseBranchForVersion(version, r.CEBranchPrefix)
			} else {
				branchName = releases.CEReleaseBranchForVersion(version, "")
			}
			branches = append(branches, branchName)

			continue
		}

		if isEntRepo {
			// Enterprise repository without --only-ce
			if versionInfo.CEActive && r.CEBranchPrefix != "" {
				// CE active version with prefix: ce/release/<version>
				branchName = releases.CEReleaseBranchForVersion(version, r.CEBranchPrefix)
			} else {
				// ENT-only or no prefix: release/<version>+ent
				branchName = releases.EnterpriseReleaseBranchForVersion(version)
			}

			branches = append(branches, branchName)

			continue
		}

		// CE repository: Only include CE active branches
		if !versionInfo.CEActive {
			continue
		}

		branches = append(branches, releases.CEReleaseBranchForVersion(version, ""))
	}

	if len(branches) == 0 {
		return nil, errors.New("no active branches found matching the specified filters")
	}

	slog.Default().DebugContext(
		ctx, "resolved branches from active versions",
		slog.Int("count", len(branches)),
		slog.Bool("only_lts", r.OnlyLTS),
		slog.Bool("only_ce", r.OnlyCE),
		slog.Bool("is_ent_repo", isEntRepo),
	)

	return branches, nil
}

// evaluateWorkflowRun matches jobs against patterns and determines success/failure.
func (r *CheckWorkflowRunReq) evaluateWorkflowRun(
	jobs []*gh.WorkflowJob,
) *CheckWorkflowRunEvaluation {
	eval := &CheckWorkflowRunEvaluation{
		IsComplete:      true,
		IsSuccessful:    true,
		IsTainted:       false,
		MatchedJobs:     make([]*CheckWorkflowRunJobMatch, 0),
		MissingPatterns: make([]string, 0),
		FailedJobs:      make([]*CheckWorkflowRunJobMatch, 0),
		TaintedJobs:     make([]*CheckWorkflowRunJobMatch, 0),
	}

	// Track which patterns have been matched
	patternMatched := make([]bool, len(r.Patterns))

	isJobFailed := func(job *gh.WorkflowJob) bool {
		return job.GetStatus() == "completed" &&
			job.GetConclusion() != "success" &&
			job.GetConclusion() != "skipped"
	}

	// Match jobs against patterns
	for _, job := range jobs {
		jobName := job.GetName()
		matched := false

		for i, re := range r.compiledPatterns {
			if re.MatchString(jobName) {
				matched = true
				patternMatched[i] = true

				match := &CheckWorkflowRunJobMatch{
					Name:       jobName,
					Pattern:    r.Patterns[i],
					Status:     job.GetStatus(),
					Conclusion: job.GetConclusion(),
				}

				eval.MatchedJobs = append(eval.MatchedJobs, match)

				// Check if matched job failed
				if isJobFailed(job) {
					eval.FailedJobs = append(eval.FailedJobs, match)
					eval.IsSuccessful = false
				}
			}
		}

		// Track failed jobs that didn't match any pattern (tainted)
		if !matched && isJobFailed(job) {
			eval.IsTainted = true
			eval.TaintedJobs = append(eval.TaintedJobs, &CheckWorkflowRunJobMatch{
				Name:       jobName,
				Status:     job.GetStatus(),
				Conclusion: job.GetConclusion(),
			})
		}
	}

	// Check for missing patterns
	for i := range len(r.Patterns) {
		if !patternMatched[i] {
			eval.IsComplete = false
			eval.MissingPatterns = append(eval.MissingPatterns, r.Patterns[i])
		}
	}

	return eval
}

func (r *CheckWorkflowRunReq) compilePatterns() ([]*regexp.Regexp, error) {
	compiledPatterns := make([]*regexp.Regexp, len(r.Patterns))
	for i := range len(r.Patterns) {
		re, err := regexp.Compile(r.Patterns[i])
		if err != nil {
			return nil, fmt.Errorf("compiling pattern %q: %w", r.Patterns[i], err)
		}
		compiledPatterns[i] = re
	}

	return compiledPatterns, nil
}

// buildBranchResult builds a branch result from a workflow run evaluation.
func (r *CheckWorkflowRunReq) buildBranchResult(
	branch string,
	run *gh.WorkflowRun,
	eval *CheckWorkflowRunEvaluation,
	success bool,
) *CheckWorkflowRunBranchResult {
	result := &CheckWorkflowRunBranchResult{
		Branch:        branch,
		Success:       success,
		Depth:         eval.Depth,
		WorkflowRunID: run.GetID(),
		Status:        run.GetStatus(),
		Conclusion:    run.GetConclusion(),
		MatchedJobs:   eval.MatchedJobs,
		CommitSHA:     run.GetHeadSHA(),
		WorkflowURL:   run.GetHTMLURL(),
	}

	if success {
		return result
	}

	determineFailureReason := func(eval *CheckWorkflowRunEvaluation) string {
		if eval.IsTainted {
			return "Workflow run has failed jobs that prevented required jobs from running"
		}
		if len(eval.FailedJobs) > 0 && len(eval.MissingPatterns) > 0 {
			return "Required jobs failed or missing"
		}
		if len(eval.FailedJobs) > 0 {
			return "Required jobs failed"
		}
		return "Required jobs missing"
	}

	result.MissingPatterns = eval.MissingPatterns
	result.FailedJobs = eval.FailedJobs
	result.FailureDetails = &CheckWorkflowRunFailureDetails{
		Reason:          determineFailureReason(eval),
		MissingPatterns: eval.MissingPatterns,
		FailedJobs:      eval.FailedJobs,
		TaintedJobs:     eval.TaintedJobs,
	}

	return result
}

// buildSummaryTable creates a table with branch summary information.
func (r *CheckWorkflowRunRes) buildSummaryTable() table.Writer {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false

	t.AppendHeader(table.Row{"Branch", "Success", "Depth", "Run URL", "Status", "Conclusion", "Error"})

	for _, br := range r.Results {
		depth := ""
		if br.Depth >= 0 {
			depth = fmt.Sprintf("%d", br.Depth)
		}

		t.AppendRow(table.Row{
			br.Branch,
			br.Success,
			depth,
			br.WorkflowURL,
			br.Status,
			br.Conclusion,
			br.Error,
		})
	}

	t.SuppressEmptyColumns()
	t.SuppressTrailingSpaces()
	return t
}

// buildMatchedJobsTable creates a table with all matched jobs across all branches.
func (r *CheckWorkflowRunRes) buildMatchedJobsTable(hideSuccessful bool) table.Writer {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false

	t.AppendHeader(table.Row{"Branch", "Depth", "Pattern", "Status", "Conclusion", "Job"})

	for _, br := range r.Results {
		for _, job := range br.MatchedJobs {
			// Skip successful jobs if hideSuccessful is true
			if hideSuccessful && job.Status == "completed" && (job.Conclusion == "success" || job.Conclusion == "skipped") {
				continue
			}

			depth := ""
			if br.Depth >= 0 {
				depth = fmt.Sprintf("%d", br.Depth)
			}

			t.AppendRow(table.Row{
				br.Branch,
				depth,
				job.Pattern,
				job.Status,
				job.Conclusion,
				job.Name,
			})
		}
	}

	t.SuppressEmptyColumns()
	t.SuppressTrailingSpaces()
	return t
}

// buildFailedJobsTable creates a table with all failed jobs across all branches.
func (r *CheckWorkflowRunRes) buildFailedJobsTable() table.Writer {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false

	t.AppendHeader(table.Row{"Branch", "Depth", "Pattern", "Status", "Conclusion", "Job"})

	for _, br := range r.Results {
		if br.FailureDetails != nil {
			for _, job := range br.FailureDetails.FailedJobs {
				depth := ""
				if br.Depth >= 0 {
					depth = fmt.Sprintf("%d", br.Depth)
				}

				t.AppendRow(table.Row{
					br.Branch,
					depth,
					job.Pattern,
					job.Status,
					job.Conclusion,
					job.Name,
				})
			}
		}
	}

	t.SuppressEmptyColumns()
	t.SuppressTrailingSpaces()
	return t
}

// buildTaintedJobsTable creates a table with all tainted jobs across all branches.
func (r *CheckWorkflowRunRes) buildTaintedJobsTable() table.Writer {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false

	t.AppendHeader(table.Row{"Branch", "Depth", "Status", "Conclusion", "Job"})

	for _, br := range r.Results {
		if br.FailureDetails != nil {
			for _, job := range br.FailureDetails.TaintedJobs {
				depth := ""
				if br.Depth >= 0 {
					depth = fmt.Sprintf("%d", br.Depth)
				}

				t.AppendRow(table.Row{
					br.Branch,
					depth,
					job.Status,
					job.Conclusion,
					job.Name,
				})
			}
		}
	}

	t.SuppressEmptyColumns()
	t.SuppressTrailingSpaces()
	return t
}

// buildMissingPatternsTable creates a table with all missing patterns across all branches.
func (r *CheckWorkflowRunRes) buildMissingPatternsTable() table.Writer {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false

	t.AppendHeader(table.Row{"Branch", "Missing Pattern"})

	for _, br := range r.Results {
		if br.FailureDetails != nil {
			for _, pattern := range br.FailureDetails.MissingPatterns {
				t.AppendRow(table.Row{
					br.Branch,
					pattern,
				})
			}
		}
	}

	t.SuppressEmptyColumns()
	t.SuppressTrailingSpaces()
	return t
}

// hasMatchedJobs returns true if any branch has matched jobs.
func (r *CheckWorkflowRunRes) hasMatchedJobs() bool {
	for _, br := range r.Results {
		if len(br.MatchedJobs) > 0 {
			return true
		}
	}
	return false
}

// hasFailedJobs returns true if any branch has failed jobs.
func (r *CheckWorkflowRunRes) hasFailedJobs() bool {
	for _, br := range r.Results {
		if br.FailureDetails != nil && len(br.FailureDetails.FailedJobs) > 0 {
			return true
		}
	}
	return false
}

// hasTaintedJobs returns true if any branch has tainted jobs.
func (r *CheckWorkflowRunRes) hasTaintedJobs() bool {
	for _, br := range r.Results {
		if br.FailureDetails != nil && len(br.FailureDetails.TaintedJobs) > 0 {
			return true
		}
	}
	return false
}

// hasMissingPatterns returns true if any branch has missing patterns.
func (r *CheckWorkflowRunRes) hasMissingPatterns() bool {
	for _, br := range r.Results {
		if br.FailureDetails != nil && len(br.FailureDetails.MissingPatterns) > 0 {
			return true
		}
	}
	return false
}
