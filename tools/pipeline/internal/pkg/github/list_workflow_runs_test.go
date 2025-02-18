// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestWorkflowRunSummaryTemplate verifies that we can correctly render a
// human readable summary from our response test fixture. We don't do strict
// value checking on rendered template. If you modify the template and/or
// response struct you'll probably need to update the test fixture.
func TestWorkflowRunSummaryTemplate(t *testing.T) {
	f, err := os.Open(filepath.Join("./testfixtures/list_workflow_runs.json"))
	require.NoError(t, err)
	bytes, err := io.ReadAll(f)
	require.NoError(t, err)
	res := &ListWorkflowRunsRes{}
	require.NoError(t, json.Unmarshal(bytes, res))
	for _, run := range res.Runs {
		require.NotNil(t, run)
		run.summary = ""
		summary, err := run.Summary()
		require.NoError(t, err)
		require.NotEmpty(t, summary)
		// t.Log(summary) // useful to see rendered output when modifying
	}
}
