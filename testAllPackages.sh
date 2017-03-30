#!/usr/bin/env bash

set -e
echo "" > coverage.txt

ginkgo -coverprofile=profile.out -covermode=atomic  *

for d in $(find . -type d -print0 | xargs -0 echo) ;do #$(go list ./... | grep -v vendor); do
    coverprofile=${d}/${d##*/}.coverprofile
    if [ -f ${coverprofile} ]; then
        cat ${coverprofile} >> coverage.txt
        rm ${coverprofile}
    fi
done