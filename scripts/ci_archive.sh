#!/bin/bash -eu
set -x

cd $(dirname ${BASH_SOURCE})/..

declare -r NAME=trafficquota
declare -r ARCHIVE_NAME=$1

mkdir "./${NAME}"
cp README.md LICENSE trafficquotad "./${ARCHIVE_NAME}"
tar zcvf "${ARCHIVE_NAME}" "./${NAME}"
