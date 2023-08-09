#!/bin/sh
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


set -e
set -x

## Make a temp dir
tempdir=$(mktemp -d plugin-semgrep.XXXXXX)
vaultdir=$(pwd)
## Set paths
cd $tempdir

for plugin in $(grep github.com/hashicorp/vault-plugin- $vaultdir/go.mod | cut -f 2 | cut -d ' ' -f 1 | cut -d '/' -f 3)
do
	if [ -z $SKIP_MODULE_UPDATING ]
	then
		echo "Fetching $plugin..."
		git clone https://github.com/hashicorp/$plugin
        semgrep --include '*.go' --exclude 'vendor' -a -f $vaultdir/tools/semgrep/ci/ $plugin/. > $plugin.semgrep.txt
	fi
done
