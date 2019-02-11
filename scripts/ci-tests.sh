#!/usr/bin/env bash

cd "$(dirname $0)/.."

set -e
echo "" > coverage.txt

for d in $(go list ./...); do
    go test -race -coverprofile=profile.out -covermode=atomic "$d"
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done
