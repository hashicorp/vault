package ssh

const (
	// This is a constant representing a script to install and uninstall public
	// key in remote hosts.
	DefaultPublicKeyInstallScript = `
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
# Currently, vault uses /home/<username>/.ssh/authorized_keys as the path.
#
# [Note: This is a default script and is written to provide convenience.
# If the host platform differs, or if the binaries used in this script are not
# available, write a new script that takes the above parameters and does the
# same task as this script, and register it Vault while role creation using
# 'install_script' parameter.

if [ $1 != "install" && $1 != "uninstall" ]; then
	exit 1
fi

# Remove the key from authorized_key file if it is already present.
# This step is common for both installing and uninstalling the key.
grep -vFf $2 $3 > temp_$2
cat temp_$2 | sudo tee $3

if [ $1 == "install" ]; then
# Append the new public key to authorized_keys file
cat $2 | sudo tee --append $3
fi

# Delete the auxiliary files
rm -f $2 temp_$2
`
)
