package ssh

const (
	// This is a constant representing a script to install and uninstall public
	// key in remote hosts.
	DefaultPublicKeyInstallScript = `
#!/bin/bash
#
# This is a default script which installs or uninstalls an RSA public key to/from
# authoried_keys file in a typical linux machine. Use 'install_script' parameter
# with 'roles/' endpoint to register a custom script (for Dynamic type).
#
# Vault server runs this script on the target machine with the following params:
#
# $1:INSTALL_OPTION: "install" or "uninstall"
#
# $2:PUBLIC_KEY_FILE: File name containing public key to be installed. Vault server
# uses UUID as file name to avoid collisions with public keys generated for requests.
#
# $3:AUTH_KEYS_FILE: Absolute path of the authorized_keys file.
# Currently, vault uses /home/<username>/.ssh/authorized_keys as the path.
#
# [Note: If the platform differs or if the binaries used in this script are not
# available in target machine, provide a custom script.]

set -e

INSTALL_OPTION=$1
PUBLIC_KEY_FILE=$2
AUTH_KEYS_FILE=$3

# Delete the public key file and the temporary file
function cleanup
{
	rm -f "$PUBLIC_KEY_FILE" temp_$PUBLIC_KEY_FILE
}

# This ensures that cleanup is called if any command fails
trap cleanup EXIT

if [ "$INSTALL_OPTION" != "install" && "$INSTALL_OPTION" != "uninstall" ]; then
	exit 1
fi

# Remove the key from authorized_key file if it is already present.
# This step is common for both installing and uninstalling the key.
grep -vFf "$PUBLIC_KEY_FILE" "$AUTH_KEYS_FILE" > temp_$PUBLIC_KEY_FILE
cat temp_$PUBLIC_KEY_FILE | sudo tee "$AUTH_KEYS_FILE"

if [ "$INSTALL_OPTION" == "install" ]; then
# Append the new public key to authorized_keys file
cat "$PUBLIC_KEY_FILE" | sudo tee --append "$AUTH_KEYS_FILE"
fi
`
)
