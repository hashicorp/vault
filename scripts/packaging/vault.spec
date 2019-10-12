# You need to define 'vault_version' on the command line when building this rpm
## rpmbuild -ba --define "vault_version 1.0.3"  SPECS/vault.spec
# Make sure that the vault binary and systemd unit file are in the SOURCES directory of the rpmbuild tree.
Name: vault
Version: %{vault_version}
Release: 1
Summary: Hashicorp Vault %{version}. A secrets management solution
Vendor: Hashicorp
URL: https://releases.hashicorp.com/vault/%{version}/vault_%{version}_linux_amd64.zip
License: Mozilla Public License 2.0
BuildRoot: %{_tmppath}/%{name}-%{version}-root
BuildRequires: coreutils
Requires(pre): /usr/sbin/useradd, /usr/bin/getent
Requires(postun): /usr/sbin/userdel
%description
Hashicorp Vault %{version} standard.
Secure, store and tightly control access to tokens, passwords, certificates, encryption keys for protecting secrets and other sensitive data using a UI, CLI, or HTTP API.

%prep
# Clean buildroot
rm -rf %{buildroot}

# Create required folders
mkdir -p %{buildroot}/etc/vault.d
mkdir -p %{buildroot}/usr/local/bin
mkdir -p %{buildroot}/usr/lib/systemd/system

# Copy in binary and unit file
cp %{_topdir}/SOURCES/vault %{buildroot}/usr/local/bin/vault
cp %{_topdir}/SOURCES/vault.service %{buildroot}/usr/lib/systemd/system/vault.service

%pre
/usr/bin/getent group vault || /usr/sbin/groupadd -r vault
/usr/bin/getent passwd vault || /usr/sbin/useradd --system --home /etc/vault.d --shell /bin/false -g vault vault

%post
# Set up autocomplete (it's ok if these fail)
/usr/local/bin/vault -autocomplete-install 2>/dev/null||true
complete -C /usr/local/bin/vault vault 2>/dev/null||true

# Give vault the ability to use mlock systemcall
setcap cap_ipc_lock=+ep /usr/local/bin/vault

# Reload systemd
systemctl daemon-reload

%preun
case "$1" in
   0) # This is a yum remove.
      /usr/sbin/userdel vault
   ;;
   1) # This is a yum upgrade.
      # do nothing
   ;;
esac

%files
%attr(0755, root, root) /usr/local/bin/vault
%attr(0660, root, root) /etc/vault.d
%attr(0664, root, root) /usr/lib/systemd/system/vault.service
%clean
# Clean buildroot
rm -rf %{buildroot}

%changelog
* Fri Oct 11 2019 Evan Chaney <evanachaney@gmail.com> - 0.1
- Initial build
