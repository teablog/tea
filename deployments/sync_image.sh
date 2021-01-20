#!/usr/bin/env bash

# shellcheck disable=SC2164
ROOTDIR=$(cd "$(dirname "$0")"; cd ..; pwd)
rsync -ravz "${ROOTDIR}/storage/images/" d2:/data/web/public/images/