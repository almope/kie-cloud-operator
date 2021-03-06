#!/bin/sh

source ./hack/go-mod-env.sh
CFLAGS="docker"

./hack/go-csv.sh

echo
echo Building operator metadata image:
echo
if [[ ${LOCAL} != true ]]; then
    CFLAGS="osbs"
    if [[ ${2} == "release" ]]; then
        CFLAGS+=" --release"
    fi
fi

echo ${CFLAGS}
cekit -v --descriptor image-metadata.yaml --redhat build \
    ${CFLAGS}
