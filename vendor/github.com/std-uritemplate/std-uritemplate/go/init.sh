#! /bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

# This is intended to be a customizable entrypoint for each language, it has to be generic enough
(
  cd ${SCRIPT_DIR}/test && \
  rm -f go.mod stduritemplate && \
  go mod init stduritemplate && \
  go mod edit -replace github.com/std-uritemplate/std-uritemplate/go=../ && \
  go mod tidy && \
  go build
)
