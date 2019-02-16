#!/bin/bash -eu
set -x

cd $(dirname ${BASH_SOURCE})/..

declare -r REPO=orgplace/trafficquota
for RPM in *.rpm; do
    package_cloud push "$REPO/el/7" "$RPM"
    package_cloud push "$REPO/fedora/29" "$RPM"
done
