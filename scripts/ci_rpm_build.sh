#!/bin/bash -eu
set -x

cd $(dirname ${BASH_SOURCE})/..

declare -r VERSION=$1

cat > trafficquota.spec <<EOS
Name:           trafficquota
Version:        ${VERSION}
Release:        1
Summary:        Token bucket server to control web API requests.
License:        Apache License Version 2.0
URL:            https://github.com/orgplace/trafficquota

%description
Token bucket server to control web API requests.
You can check whether the request conforms
to defined limits on rate and burstiness by
querying this service at the beginning of your API.

%install
mkdir -p %{buildroot}%{_bindir}
install -m 0755 trafficquotad %{buildroot}%{_bindir}/trafficquotad

%files
%defattr(-,root,root)
%{_bindir}/trafficquotad
EOS

docker run -it --rm -v $PWD:/trafficquota rpmbuild/centos7 /bin/bash /trafficquota/scripts/build_rpm.sh
