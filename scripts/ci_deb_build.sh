#!/bin/bash -eu
set -x

cd $(dirname ${BASH_SOURCE})/..

declare -r VERSION="$1"
declare -r WORK_DIR="${HOME}/deb"

mkdir -p "${WORK_DIR}/usr/bin" "${WORK_DIR}/DEBIAN" "${WORK_DIR}/usr/lib/systemd/system"

cat > "${WORK_DIR}/DEBIAN/control" <<EOS
Package: trafficquota
Maintainer: $(git log -n1 --pretty=format:"%an <%ae>")
Architecture: amd64
Version: ${VERSION}
Description: Token bucket server to control web API requests.
 You can check whether the request conforms
 to defined limits on rate and burstiness by
 querying this service at the beginning of your API.
EOS
cp ./trafficquotad "${WORK_DIR}/usr/bin"
cat > "${WORK_DIR}/usr/lib/systemd/system/trafficquotad.service" <<EOS
[Unit]
Description=Token bucket server to control web API requests.

[Service]
ExecStart=/usr/bin/trafficquotad
Restart=always

[Install]
WantedBy=multi-user.target
EOS

fakeroot dpkg-deb --build "${WORK_DIR}" .
