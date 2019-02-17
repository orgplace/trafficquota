#!/bin/bash -eu
set -x

cd $(dirname ${BASH_SOURCE})/..

declare -r REPO=orgplace/trafficquota
for RPM in *.rpm; do
    # Latest RHEL
    package_cloud push "$REPO/el/7" "$RPM"
    # Latest Fedora
    package_cloud push "$REPO/fedora/29" "$RPM"
done
for DEB in *.deb; do
    # Latest Ubuntu
    package_cloud push "$REPO/ubuntu/cosmic" "$DEB"
    # Latest Ubuntu LTS
    package_cloud push "$REPO/ubuntu/bionic" "$DEB"
    # Latest Debian stable
    package_cloud push "$REPO/debian/stretch" "$DEB"
done
