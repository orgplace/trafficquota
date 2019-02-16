#!/bin/bash -eu
set -x

cd $(dirname ${BASH_SOURCE})/..

declare -r ARCHIVE_DIR=trafficquota
declare -r ARCHIVE_NAME=$1

mkdir "./${ARCHIVE_DIR}"
cp README.md LICENSE trafficquotad "./${ARCHIVE_DIR}"
tar zcvf "${ARCHIVE_NAME}" "./${ARCHIVE_DIR}"
