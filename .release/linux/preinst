#!/bin/bash

set -eu

USER="vault"

if ! id -u $USER > /dev/null 2>&1; then
	useradd \
		--system \
		--user-group \
		--shell /bin/false \
		$USER
fi
