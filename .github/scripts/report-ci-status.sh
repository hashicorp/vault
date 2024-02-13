#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e
MAX_TESTS=10

# All of these environment variables are required or an error will be returned.
[ "${GITHUB_TOKEN:?}" ]
[ "${RUN_ID:?}" ]
[ "${REPO:?}" ]
[ "${PR_NUMBER:?}" ]
[ "${RESULT:?}" ]

table_data() {
  if [ -z "$TABLE_DATA" ]; then
    return 0
  fi

  # Remove any rows that don't have a test name
  # Only keep the test type, test package, test name, and logs column
  # Remove the scroll emoji
  # Remove "github.com/hashicorp/vault" from the package name
  TABLE_DATA=$(echo "$TABLE_DATA" | awk -F\| '{if ($4 != " - ") { print "|" $2 "|" $3 "|" $4 "|" $7 }}' | sed -r 's/ :scroll://' | sed -r 's/github.com\/hashicorp\/vault\///')
  NUM_FAILURES=$(wc -l <<< "$TABLE_DATA")

  # Check if the number of failures is greater than the maximum tests to display
  # If so, limit the table to MAX_TESTS number of results
  if [ "$NUM_FAILURES" -gt "$MAX_TESTS" ]; then
    TABLE_DATA=$(echo "$TABLE_DATA" | head -n "$MAX_TESTS")
    NUM_OTHER=( "$NUM_FAILURES" - "$MAX_TESTS" )
    TABLE_DATA="${TABLE_DATA}

and ${NUM_OTHER[*]} other tests"
  fi

  # Add the header for the table
  printf "%s" "Failures:
| Test Type | Package | Test | Logs |
| --------- | ------- | ---- | ---- |
${TABLE_DATA}"
}

td="$(table_data)"

case "$RESULT" in
  success)
    if [ -z "$td" ]; then
      BODY="CI Results:
All Go tests succeeded! :white_check_mark:"
    else
      BODY="CI Results:
All required Go tests succeeded but failures were detected :warning:
${td}"
    fi
  ;;
  *)
    BODY="CI Results: ${RESULT} :x:
${td}"
  ;;
esac

source ./.github/scripts/gh-comment.sh

update_or_create_comment "$REPO" "$PR_NUMBER" "CI Results:" "$BODY"
