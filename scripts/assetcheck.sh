#!/usr/bin/env bash

if [[ ! -e http/web_ui/index.html ]]
then
  printf "Compiled UI assets not found. They can be built with: make static-dist\n\n"
  exit 1
fi
