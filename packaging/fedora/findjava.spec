Name:      findjava
Version:   ${version}
Release:   1%{?dist}
Summary:   A tool for finding an appropriate installed JVM to run your program

License:   Apache-2.0
URL:       https://github.com/loicrouchon/findjava
Source0:   https://github.com/loicrouchon/findjava/archive/refs/tags/v${version}.tar.gz
# Source0:   https://github.com/loicrouchon/findjava/archive/refs/heads/<BRANCH>.zip

BuildArch: x86_64 aarch64
BuildRequires: make, java-latest-openjdk-devel, golang

%description
findjava is a command-line tool for finding an appropriate installed JVM to run your program.
It provides a simple and efficient way to specify what version of JVM you want
and what kind of features (java, javac, native-image, ...) it should provide.

%global debug_package %{nil}

%prep

%setup -q -n findjava-${version}
%build
GO_LD_FLAGS='-linkmode=external' GO_TAGS="-tags linux" make test build

%install
%define distdir build
find .
mkdir -p %{buildroot}/usr/bin %{buildroot}/usr/share/%{name} %{buildroot}/usr/share/%{name}/metadata-extractor %{buildroot}/etc/%{name}
ln -s ../share/%{name}/%{name} %{buildroot}/usr/bin/%{name}
install -p -m 755 %{distdir}/dist/%{name} %{buildroot}/usr/share/%{name}/%{name}
install -p -m 644 %{distdir}/dist/metadata-extractor/JvmMetadataExtractor.class %{buildroot}/usr/share/%{name}/metadata-extractor/JvmMetadataExtractor.class
install -p -m 644 packaging/fedora/config.conf %{buildroot}/etc/%{name}/config.conf

%files
%license LICENSE
/usr/bin/%{name}
/usr/share/%{name}/%{name}
/usr/share/%{name}/metadata-extractor/JvmMetadataExtractor.class
/etc/%{name}/config.conf
