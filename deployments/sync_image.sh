#!/usr/bin/env bash

# shellcheck disable=SC2164
ROOTDIR=$(cd "$(dirname "$0")"; cd ..; pwd)
rsync -ravz --exclude=".DS_Store" "${ROOTDIR}/storage/images/" teablog:/data/web/public/images/