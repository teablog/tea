#!/usr/bin/env bash

# shellcheck disable=SC2164
CURDIR=$(cd "$(dirname "$0")"; pwd)

echo '--- build docker image';
docker build -t registry.cn-hangzhou.aliyuncs.com/douyacun/tea:latest .

echo '--- push docker image to aliyun';
docker push registry.cn-hangzhou.aliyuncs.com/douyacun/tea:latest

echo "--- deploy tea";
ssh d2 < "${CURDIR}"/deploy.sh

#echo "--- docker image prune..."
#echo y|docker image prune