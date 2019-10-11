# You need to define 'vault_version' and 'vault_bin' on the command line when building this rpm
## rpmbuild -ba --define "vault_version 1.0.3" --define "vault_bin /tmp/vault" SPECS/vault.spec
Name: vault
Version: %{vault_version}
Release: 1
Summary: Hashicorp Vault %{version}. A secrets management solution
Vendor: Hashicorp
URL: https://releases.hashicorp.com/vault/%{version}/vault_%{version}_linux_amd64.zip
License: Mozilla Public License 2.0
BuildRoot: %{_tmppath}/%{name}-%{version}-root
BuildRequires: coreutils

%description
Hashicorp Vault %{version} standard.
Secure, store and tightly control access to tokens, passwords, certificates, encryption keys for protecting secrets and other sensitive data using a UI, CLI, or HTTP API.

%prep
# Clean buildroot
rm -rf %{buildroot}

# Create required folders
mkdir -p %{buildroot}/etc/vault
mkdir -p %{buildroot}/usr/bin

# Pull new version

cp %{vault_bin} %{buildroot}/usr/bin/vault



%post
# Do nothing

%preun
# Nothing to do here

%files
%attr(0755, root, root) /usr/bin/vault
%attr(0660, root, root) /etc/vault

%clean
# Clean buildroot
rm -rf %{buildroot}

%changelog
* Fri Oct 11 2019 Evan Chaney <evanachaney@gmail.com> - 0.1
- Initial build
