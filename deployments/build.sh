#!/usr/bin/env bash

# shellcheck disable=SC2164
CURDIR=$(cd "$(dirname "$0")"; pwd)

docker build -t registry.cn-hangzhou.aliyuncs.com/douyacun/tea:latest .
docker push registry.cn-hangzhou.aliyuncs.com/douyacun/tea:latest
ssh douyacun < "${CURDIR}"/deploy.sh

echo y|docker image prune