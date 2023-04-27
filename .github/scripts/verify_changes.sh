#!/bin/bash
#  This script validates if the git diff contains on docs changes

event_type=$1 # GH event type (pull_request)
ref_name=$2 # branch reference that triggered the workflow
head_ref=$3 # PR branch head ref
base_ref=$4 # PR branch base ref

changed_dir=""

if [[ "$event_type" == "pull_request" ]]; then
  git fetch --no-tags --prune origin $head_ref
  git fetch --no-tags --prune origin $base_ref
  head_commit="origin/$head_ref"
  base_commit="origin/$base_ref"
else
  git fetch --no-tags --prune origin $ref_name
  head_commit=$(git log origin/$ref_name --oneline | head -1 | awk '{print $1}')
  base_commit=$(git log origin/$ref_name --oneline | head -2 | awk 'NR==2 {print $1}')
fi

# git diff with ... shows the differences count between base_commit and head_commit starting at the last common commit
change_count=$(git diff $base_commit...$head_commit --name-only | awk -F"/" '{ print $1}' | uniq | wc -l)

if [[ $change_count -eq 1 ]]; then
  changed_dir=$(git diff $base_commit...$head_commit --name-only | awk -F"/" '{ print $1}' | uniq)
fi

if [[ "$changed_dir" == "website" ]]; then
  echo "is_docs_change=true" >> "$GITHUB_OUTPUT"
else
  echo "is_docs_change=false" >> "$GITHUB_OUTPUT"
fi
