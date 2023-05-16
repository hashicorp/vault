#!/bin/bash
set -eux

CLEANUP_PAUSE=${CLEANUP_PAUSE:-0}
echo "==> Pausing for ${CLEANUP_PAUSE} seconds..."
sleep ${CLEANUP_PAUSE}

# Unique SSH keys will be generated on first boot
echo "==> Removing SSH server keys"
rm -f /etc/ssh/*_key*

# Unique machine ID will be generated on first boot
#echo "==> Removing machine ID"
rm -f /etc/machine-id
touch /etc/machine-id

echo "==> Cleaning up leftover dhcp leases"
if [ -d "/var/lib/dhcp" ]; then
    rm -f /var/lib/dhcp/*
fi

echo "==> Cleaning up tmp"
rm -rf /tmp/*

# Cleanup apt cache
apt-get -y autoremove --purge
apt-get -y clean
apt-get -y autoclean

echo "==> Installed packages"
dpkg --get-selections | grep -v deinstall

# Remove Bash history
unset HISTFILE
rm -f /root/.bash_history
rm -rf /root/.ssh/*
rm -rf /root/*

# Clean sudoers
rm -f /etc/sudoers.d/90-cloud-init-users

# Clean up log files
echo "==> Purging log files"
find /var/log -type f -delete

# Deny SSH root login
sed -i 's/.*PermitRootLogin.*/#PermitRootLogin No/g' /etc/ssh/sshd_config

# Clear root password
passwd -d root
passwd -l root

# Remove ubuntu user
echo "==> Removing ubuntu user"
userdel -f ubuntu
rm -rf /home/ubuntu
