%global provider        github
%global provider_tld    com
%global project         hashicorp
%global repo            vault
%global provider_prefix %{provider}.%{provider_tld}/%{project}/%{repo}
%global import_path     %{provider_prefix}
%if 0%{?_commit:1}
# rpmbuild -ba --define "_commit 409b441c31f83279af0db289123eb4b0b14809a6" *.spec
%global commit          %{_commit}
%else
%global commit          15982cfa072fc6898f56d8320a460b5bafb7606b
%endif
%global shortcommit     %(c=%{commit}; echo ${c:0:7})

%if ! 0%{?gobuild:1}
%define gobuild(o:) go build -ldflags "${LDFLAGS:-} -B 0x$(head -c20 /dev/urandom|od -An -tx1|tr -d ' \\n')" -a -v -x %{?**};
%endif

%if ! 0%{?gotest:1}
%define gotest() go test -ldflags "${LDFLAGS:-}" %{?**}
%endif

Name:           vault
Version:        0.9.3
Release:        2.git%{shortcommit}%{?dist}
Summary:        Vault is a tool for securely accessing secrets
License:        MPLv2.0
URL:            https://www.vaultproject.io
Source0:        https://%{provider_prefix}/archive/%{commit}/%{repo}-%{shortcommit}.tar.gz
Source1:        %{name}.config
Source2:        %{name}.service

# e.g. el6 has ppc64 arch without gcc-go, so EA tag is required
ExclusiveArch:  %{?go_arches:%{go_arches}}%{!?go_arches:%{ix86} x86_64 %{arm}}
# If go_compiler is not set to 1, there is no virtual provide. Use golang instead.
BuildRequires:  %{?go_compiler:compiler(go-compiler)}%{!?go_compiler:golang} >= 1.9
BuildRequires:  systemd
BuildRequires:  git

Requires(pre):          shadow-utils
Requires(post):         systemd libcap
Requires(preun):        systemd
Requires(postun):       systemd

%description
Vault secures, stores, and tightly controls access to tokens, passwords,
certificates, API keys, and other secrets in modern computing. Vault handles
leasing, key revocation, key rolling, and auditing. Through a unified API, users
can access an encrypted Key/Value store and network encryption-as-a-service, or
generate AWS IAM/STS credentials, SQL/NoSQL databases, X.509 certificates, SSH
credentials, and more.

%package selinux
Source3:        	%{name}.te
Source4:        	%{name}.fc
Source5:        	%{name}.if

BuildRequires:  	selinux-policy selinux-policy-devel
Requires:       	policycoreutils, libselinux-utils

Requires(post):         policycoreutils, policycoreutils-python 
Requires(postun):       policycoreutils

Summary: SELinux policy for %{name}

%prep
%setup -q -n %{name}-%{commit}

mkdir -p selinux
cd selinux
cp %{SOURCE3} %{SOURCE4} %{SOURCE5} .
cd -

%build
mkdir -p src/%{provider}.%{provider_tld}/%{project}
ln -s ../../../ src/%{provider}.%{provider_tld}/%{project}/%{repo}

ls -lhA
ls -lhA src/%{provider}.%{provider_tld}/%{project}/%{repo}

%if ! 0%{?with_bundled}
export GOPATH=$(pwd):%{gopath}
%else
export GOPATH=$(pwd):$(pwd)/Godeps/_workspace:%{gopath}
%endif

%gobuild -o bin/%{name} %{import_path} || exit 1

cd selinux
make -f /usr/share/selinux/devel/Makefile %{name}.pp || exit
cd -

%install
mkdir -p %{buildroot}%{_bindir}/
cp -p bin/%{name} %{buildroot}%{_bindir}/

mkdir -p %{buildroot}%{_sysconfdir}/%{name}
cp -p %{SOURCE1} %{buildroot}%{_sysconfdir}/%{name}/%{name}.config

mkdir -p %{buildroot}%{_sharedstatedir}/%{name}

mkdir -p %{buildroot}%{_unitdir}
cp -p %{SOURCE2} %{buildroot}%{_unitdir}

install -d %{buildroot}%{_datadir}/selinux/packages
install -m 0600 selinux/%{name}.pp %{buildroot}%{_datadir}/selinux/packages
install -d %{buildroot}%{_datadir}/selinux/devel/include/contrib
install -m 0644 selinux/%{name}.if %{buildroot}%{_datadir}/selinux/devel/include/contrib/

%clean
rm -rf %{buildroot}
rm -rf %{_builddir}/*

%files
%{_bindir}/%{name}
%dir %{_sysconfdir}/%{name}
%config(noreplace) %{_sysconfdir}/%{name}/%{name}.config
%attr(0750,%{name},%{name}) %dir %{_sharedstatedir}/%{name}
%{_unitdir}/%{name}.service

%doc CHANGELOG.md README.md LICENSE CHANGELOG.md

%pre
getent group %{name} > /dev/null || groupadd -r %{name}
getent passwd %{name} > /dev/null || \
    useradd -r -d %{_sharedstatedir}/%{name} -g %{name} \
    -s /sbin/nologin -c "Vault secret management tool" %{name}
exit 0

%post
/sbin/setcap cap_ipc_lock=+ep %{_bindir}/%{name}
%systemd_post %{name}.service

%preun
%systemd_preun %{name}.service

%postun
%systemd_postun_with_restart %{name}.service

%description selinux
SELinux policy for %{name}

%files selinux
%{_datadir}/selinux/packages/%{name}.pp
%{_datadir}/selinux/devel/include/contrib/%{name}.if

%post selinux
semodule -n -i %{_datadir}/selinux/packages/%{name}.pp
if /usr/sbin/selinuxenabled ; then
    /usr/sbin/load_policy
    %relabel_files
fi;
semanage port -a -t %{name}_port_t -p tcp "8201,8202"
semanage port -a -t %{name}_port_t -p udp "8201,8202"
exit 0

%postun selinux
if [ $1 -eq 0 ]; then
    semanage port -d -t %{name}_port_t -p tcp "8201,8202"
    semanage port -d -t %{name}_port_t -p udp "8201,8202"

    semodule -n -r %{name}
    if /usr/sbin/selinuxenabled ; then
       /usr/sbin/load_policy
       %relabel_files
    fi;
fi;
exit 0

%changelog
* Thu Feb 01 2018 fuero <fuerob@gmail.com> - 0.9.3-2.git15982cf
- splits selinux policy to extra package

* Wed Jan 31 2018 fuero <fuerob@gmail.com> - 0.9.3-1.git15982cf
- initial package

