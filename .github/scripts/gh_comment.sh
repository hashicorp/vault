#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


set -e


function update_or_create_comment {
  local REPO="$1"
  local PR_NUMBER="$2"
  local SEARCH_KEY="$3"
  local BODY="$4"
  local COMMENT_ID_QUERY
  
  # We only want for the GH bot to place one comment to report build failures
  # and if we rerun a job, that comment needs to be updated.
  # Let's try to find if the GH bot has placed a similar comment
  printf -v COMMENT_ID_QUERY '.[] | select (.body | startswith(%s)) | .id' "${SEARCH_KEY}"
  comment_id="$(
    gh api "/repos/hashicorp/${REPO}/issues/${PR_NUMBER}/comments" \
      --header "Accept: application/vnd.github+json" \
      --header "X-GitHub-Api-Version: 2022-11-28" \
      --paginate \
      --jq "${COMMENT_ID_QUERY}"
  )"

  if [[ "$comment_id" != "" ]]; then
    # Update the comment with the new body
    gh api "/repos/hashicorp/${REPO}/issues/comments/${comment_id}" \
      --method PATCH \
      --header "Accept: application/vnd.github+json" \
      --header "X-GitHub-Api-Version: 2022-11-28" \
      --raw-field "body=${BODY}"
  else
    # Create a comment with the new body
    gh api "/repos/hashicorp/${REPO}/issues/${PR_NUMBER}/comments" \
      --method POST \
      --header "Accept: application/vnd.github+json" \
      --header "X-GitHub-Api-Version: 2022-11-28" \
      --raw-field "body=${BODY}"
  fi
}
