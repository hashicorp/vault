#!/bin/bash

event_type=$1
base_ref=$2

changed_dir=""

if [[ "$1" == "pull_request" ]]; then
  change_count=$(git diff --merge-base $2 --name-only |xargs dirname |awk -F"/" '{ print $1}' |uniq |wc -l)
  if [[ $change_count -eq 1 ]]; then
    changed_dir=$(git diff --merge-base $2 --name-only |xargs dirname |awk -F"/" '{ print $1}' |uniq)
  fi
else
  change_count=$(git show --oneline --dirstat=files,0 |grep "%" |awk '{print $2}' |awk -F"/" '{print $1}' |uniq |wc -l)
  if [[ $change_count -eq 1 ]]; then
    changed_dir=$(git show --oneline --dirstat=files,0 |grep "%" |awk '{print $2}' |awk -F"/" '{print $1}' |uniq)
  fi
fi

if [[ "$changed_dir" == "website" ]]; then
  echo "is_docs_change=true" >> "$GITHUB_OUTPUT"
else
  echo "is_docs_change=false" >> "$GITHUB_OUTPUT"
fi
