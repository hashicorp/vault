#!/bin/bash
#
# This script file adds a RSA public key to the authoried_keys file in a typical
# linux machine. This script should be registered with vault server while creating
# a role for key type 'dynamic'.
#
# Vault server runs this script on the target machine with the following params:
# $1: File containing public key to be installed. Vault server uses UUID as file
# name to avoid collisions with public keys generated for requests.
#
# $2: Absolute path of the authorized_keys file.
#
# [Note: Modify the script if targt machine does not have the commands used in
# this script]

# If the key being installed is already present in the authorized_keys file, it is
# removed and the result is stored in a temporary file.
grep -vFf $1 $2 > temp_$1

# Contents of temporary file will be the contents of authorized_keys file.
cat temp_$1 > $2

# New public key is appended to authorized_keys file
cat $1 >> $2

# Auxiliary files are deleted
rm -f $1 temp_$1
