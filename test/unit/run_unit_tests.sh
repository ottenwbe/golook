#!/usr/bin/env bash

set -e
echo "" > coverage.txt

# run tests
ginkgo -coverprofile=profile.out -covermode=atomic  ./...

# collect coverprofiles from all tested packages
for d in $(find . -type d -print0 | xargs -0 echo); do
    coverprofile=${d}/${d##*/}.coverprofile
    if [ -f ${coverprofile} ]; then
        cat ${coverprofile} >> coverage.txt
        rm ${coverprofile}
    fi
done