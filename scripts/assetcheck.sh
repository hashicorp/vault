#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


if [[ ! -e http/web_ui/index.html ]]
then
  printf "Compiled UI assets not found. They can be built with: make static-dist\n\n"
  exit 1
else
  if [[ `find http/web_ui/index.html -mmin +10080` ]]
  then
    printf "Compiled UI assets are more than one week old. They can be rebuilt with: make static-dist\n\n"
    exit 1
  fi
fi
