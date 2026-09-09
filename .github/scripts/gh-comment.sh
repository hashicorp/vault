#!/bin/bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

function update_or_create_comment {
  REPO=$1
  PR_NUMBER=$2
  SEARCH_KEY=$3
  BODY=$4

  # We only want for the GH bot to place one comment to report build failures
  # and if we rerun a job, that comment needs to be updated.
  # Let's try to find if the GH bot has placed a similar comment
  comment_id=$(gh api \
                 -H "Accept: application/vnd.github+json" \
                 -H "X-GitHub-Api-Version: 2022-11-28" \
                 --paginate \
                 /repos/hashicorp/"$REPO"/issues/"$PR_NUMBER"/comments |
                 jq -r --arg SEARCH_KEY "$SEARCH_KEY" '.[] | select (.body | startswith($SEARCH_KEY)) | .id')

  if [[ "$comment_id" != "" ]]; then
    # update the comment with the new body
    gh api \
      --method PATCH \
      -H "Accept: application/vnd.github+json" \
      -H "X-GitHub-Api-Version: 2022-11-28" \
      /repos/hashicorp/"$REPO"/issues/comments/"$comment_id" \
      -f body="$BODY"
  else
    # create a comment with the new body
    gh api \
      --method POST \
      -H "Accept: application/vnd.github+json" \
      -H "X-GitHub-Api-Version: 2022-11-28" \
      /repos/hashicorp/"$REPO"/issues/"$PR_NUMBER"/comments \
      -f body="$BODY"
  fi
}

function delete_comment_if_exists {
  REPO=$1
  PR_NUMBER=$2
  SEARCH_KEY=$3

  # Look for an existing GH bot comment matching the search key, mirroring the
  # lookup logic in update_or_create_comment. Limit to the first match so that
  # multiple matching comments can't produce a multi-line id and an invalid
  # API path.
  comment_id=$(gh api \
                 -H "Accept: application/vnd.github+json" \
                 -H "X-GitHub-Api-Version: 2022-11-28" \
                 --paginate \
                 /repos/hashicorp/"$REPO"/issues/"$PR_NUMBER"/comments |
                 jq -r --arg SEARCH_KEY "$SEARCH_KEY" 'first(.[] | select (.body | startswith($SEARCH_KEY)) | .id) // empty' |
                 head -n 1)

  # If no matching comment exists, there's nothing to do.
  if [[ "$comment_id" != "" ]]; then
    gh api \
      --method DELETE \
      -H "Accept: application/vnd.github+json" \
      -H "X-GitHub-Api-Version: 2022-11-28" \
      /repos/hashicorp/"$REPO"/issues/comments/"$comment_id"
  fi
}

