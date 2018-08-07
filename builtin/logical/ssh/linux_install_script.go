package ssh

const (
	// This is a constant representing a script to install and uninstall public
	// key in remote hosts.
	DefaultPublicKeyInstallScript = `
#!/bin/bash
#
# This is a default script which installs or uninstalls an RSA public key to/from
# authorized_keys file in a typical linux machine.
#
# If the platform differs or if the binaries used in this script are not available
# in target machine, use the 'install_script' parameter with 'roles/' endpoint to
# register a custom script (applicable for Dynamic type only).
#
# Vault server runs this script on the target machine with the following params:
#
# $1:INSTALL_OPTION: "install" or "uninstall"
#
# $2:PUBLIC_KEY_FILE: File name containing public key to be installed. Vault server
# uses UUID as name to avoid collisions with public keys generated for other requests.
#
# $3:AUTH_KEYS_FILE: Absolute path of the authorized_keys file.
# Currently, vault uses /home/<username>/.ssh/authorized_keys as the path.
#
# [Note: This script will be run by Vault using the registered admin username.
# Notice that some commands below are run as 'sudo'. For graceful execution of
# this script there should not be any password prompts. So, disable password
# prompt for the admin username registered with Vault.

set -e

# Storing arguments into variables, to increase readability of the script.
INSTALL_OPTION=$1
PUBLIC_KEY_FILE=$2
AUTH_KEYS_FILE=$3

# Delete the public key file and the temporary file
function cleanup
{
	rm -f "$PUBLIC_KEY_FILE" temp_$PUBLIC_KEY_FILE
}

# 'cleanup' will be called if the script ends or if any command fails.
trap cleanup EXIT

# Return if the option is anything other than 'install' or 'uninstall'.
if [ "$INSTALL_OPTION" != "install" ] && [ "$INSTALL_OPTION" != "uninstall" ]; then
	exit 1
fi

# use locking to avoid parallel script execution
(
	flock --timeout 10 200
	# Create the .ssh directory and authorized_keys file if it does not exist
	SSH_DIR=$(dirname $AUTH_KEYS_FILE)
	sudo mkdir -p "$SSH_DIR"
	sudo touch "$AUTH_KEYS_FILE"
	# Remove the key from authorized_keys file if it is already present.
	# This step is common for both install and uninstall.  Note that grep's
	# return code is ignored, thus if grep fails all keys will be removed
	# rather than none and it fails secure
	sudo grep -vFf "$PUBLIC_KEY_FILE" "$AUTH_KEYS_FILE" > temp_$PUBLIC_KEY_FILE || true
	cat temp_$PUBLIC_KEY_FILE | sudo tee "$AUTH_KEYS_FILE"
	# Append the new public key to authorized_keys file
	if [ "$INSTALL_OPTION" == "install" ]; then
		cat "$PUBLIC_KEY_FILE" | sudo tee --append "$AUTH_KEYS_FILE"
	fi
) 200> ${AUTH_KEYS_FILE}.lock
`
)
