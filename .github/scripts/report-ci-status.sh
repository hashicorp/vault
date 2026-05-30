#!/bin/bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e
MAX_TESTS=10

# All of these environment variables are required or an error will be returned.
[ "${GITHUB_TOKEN:?}" ]
[ "${RUN_ID:?}" ]
[ "${REPO:?}" ]
[ "${PR_NUMBER:?}" ]
[ "${RESULT:?}" ]

# Function to format table data for a specific test type
format_table_data() {
  local data="$1"
  local test_type="$2"
  
  if [ -z "$data" ]; then
    return 0
  fi

  # Remove any rows that don't have a test name
  # Only keep the test type, test package, test name, and logs column
  # Remove the scroll emoji
  # Remove "github.com/hashicorp/vault" from the package name
  local formatted_data=$(echo "$data" | awk -F\| '{if ($4 != " - ") { print "|" $2 "|" $3 "|" $4 "|" $7 }}' | sed -r 's/ :scroll://' | sed -r 's/github.com\/hashicorp\/vault\///')
  local num_failures=$(echo "$formatted_data" | wc -l)

  # Check if the number of failures is greater than the maximum tests to display
  # If so, limit the table to MAX_TESTS number of results
  if [ "$num_failures" -gt "$MAX_TESTS" ]; then
    formatted_data=$(echo "$formatted_data" | head -n "$MAX_TESTS")
    local num_other=$(( num_failures - MAX_TESTS ))
    formatted_data="${formatted_data}

and ${num_other} other tests"
  fi

  # Add the header for the table
  printf "%s" "Failures:
| Test Type | Package | Test | Logs |
| --------- | ------- | ---- | ---- |
${formatted_data}"
}

source ./.github/scripts/gh-comment.sh

# Separate enos failures from other test failures
if [ -n "$TABLE_DATA" ]; then
  ENOS_DATA=$(echo "$TABLE_DATA" | grep "^| enos |" || true)
  OTHER_DATA=$(echo "$TABLE_DATA" | grep -v "^| enos |" || true)
else
  ENOS_DATA=""
  OTHER_DATA=""
fi

# Create comment for regular test failures (non-enos)
if [ -n "$OTHER_DATA" ]; then
  other_td="$(format_table_data "$OTHER_DATA" "other")"
  
  case "$RESULT" in
    success)
      OTHER_BODY="CI Results - Go Tests:
All required Go tests succeeded but failures were detected :warning:
${other_td}"
    ;;
    *)
      OTHER_BODY="CI Results - Go Tests: ${RESULT} :x:
${other_td}"
    ;;
  esac
  
  update_or_create_comment "$REPO" "$PR_NUMBER" "CI Results - Go Tests:" "$OTHER_BODY"
fi

# Create separate comment for enos failures
if [ -n "$ENOS_DATA" ]; then
  enos_td="$(format_table_data "$ENOS_DATA" "enos")"
  
  case "$RESULT" in
    success)
      ENOS_BODY="CI Results - Enos Tests:
All required Enos tests succeeded but failures were detected :warning:
${enos_td}"
    ;;
    *)
      ENOS_BODY="CI Results - Enos Tests: ${RESULT} :x:
${enos_td}"
    ;;
  esac
  
  update_or_create_comment "$REPO" "$PR_NUMBER" "CI Results - Enos Tests:" "$ENOS_BODY"
fi

# If no failures at all, create a single success comment
if [ -z "$OTHER_DATA" ] && [ -z "$ENOS_DATA" ]; then
  if [ "$RESULT" == "success" ]; then
    SUCCESS_BODY="CI Results:
All Go tests succeeded! :white_check_mark:"
    update_or_create_comment "$REPO" "$PR_NUMBER" "CI Results:" "$SUCCESS_BODY"
  fi
fi
