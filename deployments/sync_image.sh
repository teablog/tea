#!/usr/bin/env bash

# shellcheck disable=SC2164
CURDIR=$(cd "../../$(dirname "$0")"; pwd)

rsync -ravz "${CURDIR}/storage/images/" teablog:/data/web/images