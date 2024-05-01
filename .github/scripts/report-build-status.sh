#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


set -e

# All of these environment variables are required or an error will be returned.
[ "${GITHUB_TOKEN:?}" ]
[ "${PR_NUMBER:?}" ]
[ "${REPO:?}" ]
[ "${RUN_ID:?}" ]

# list of build jobs
[ "${ARTIFACTS:?}" ]
[ "${TEST:?}" ]
[ "${TEST_CONTAINERS:?}" ]
[ "${UI:?}" ]

# Build jobs
jobs=("artifacts:$ARTIFACTS" "test:$TEST" "test-containers:$TEST_CONTAINERS" "ui:$UI")

# Sometimes failed jobs can have a result of "cancelled". Handle both.
failed_jobs=()
for job in "${jobs[@]}";do
  if [[ "$job" == *"failure"* || "$job" == *"cancelled"* ]]; then
    failed_jobs+=("$job")
  fi
done

# Create a comment body to set on the pull request which reports failed jobs with a url to the
# failed workflow.
if [ ${#failed_jobs[@]} -eq 0 ]; then
  new_body="Build Results:
All builds succeeded! :white_check_mark:"
else
  new_body="Build Results:
Build failed for these jobs: ${failed_jobs[*]}. Please refer to this workflow to learn more: https://github.com/hashicorp/vault/actions/runs/$RUN_ID"
fi

source ./.github/scripts/gh-comment.sh

update_or_create_comment "$REPO" "$PR_NUMBER" "Build Results:" "$new_body"
