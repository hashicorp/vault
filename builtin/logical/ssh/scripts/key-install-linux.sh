#!/bin/bash
#
# This script file installs or uninstalls an RSA public key to/from authoried_keys
# file in a typical linux machine. This script should be registered with vault
# server while creating a role for key type 'dynamic'.
#
# Vault server runs this script on the target machine with the following params:
#
# $1: "install" or "uninstall"
#
# $2: File name containing public key to be installed. Vault server uses UUID
# as file name to avoid collisions with public keys generated for requests.
#
# $3: Absolute path of the authorized_keys file.
#
# [Note: Modify the script if targt machine does not have the commands used in
# this script]

if [ $1 != "install" && $1 != "uninstall" ]; then
	exit 1
fi

# If the key being installed is already present in the authorized_keys file, it is
# removed and the result is stored in a temporary file.
grep -vFf $2 $3 > temp_$2

# Contents of temporary file will be the contents of authorized_keys file.
cat temp_$2 > $3

if [ $1 == "install" ]; then
# New public key is appended to authorized_keys file
cat $2 >> $3
fi

# Auxiliary files are deleted
rm -f $2 temp_$2
