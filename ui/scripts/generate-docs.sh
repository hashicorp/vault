#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

for FILE in ./lib/core/addon/components/*.js;  do
  component=`eval "echo $FILE | cut -d/ -f6"`; 
  yarn generate-docfy-md $component core
done