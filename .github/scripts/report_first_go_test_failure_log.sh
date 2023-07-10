#!/bin/bash

set -e

# this script expects the following env vars to be set
# error if these are not set
[ ${GITHUB_TOKEN:?} ]
[ ${RUN_ID:?} ]
[ ${REPO:?} ]
[ ${PR:?} ]

failed_tests_ids=$(gh api \
  -H "Accept: application/vnd.github+json" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  /repos/hashicorp/"$REPO"/actions/runs/"$RUN_ID"/attempts/1/jobs | jq -r '[.jobs[] | select(.name | startswith("Run Go tests ")) | select(.conclusion == "failure") | .id][0]')

extracted_error=$(gh run view --job $failed_tests_ids --log | grep "\#\#\[error\]")

gh pr comment "$PR" --body "$extracted_error"
