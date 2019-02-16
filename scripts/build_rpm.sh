#!/bin/bash -eu
set -x

mkdir -p /home/builder/rpm
cd /home/builder/rpm
cp /trafficquota/trafficquotad /trafficquota/trafficquota.spec .
rpmbuild -bb trafficquota.spec
sudo mv x86_64/*.rpm /trafficquota/
