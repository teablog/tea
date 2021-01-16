#!/usr/bin/env bash

# shellcheck disable=SC2164
CURDIR=$(cd "$(dirname "$0")"; pwd)

echo '--- build docker image';
docker build -t registry.cn-hangzhou.aliyuncs.com/douyacun/tea:latest .

echo '--- push image to aliyun';
docker push registry.cn-hangzhou.aliyuncs.com/douyacun/tea:latest

echo "--- deploy tea";
ssh douyacun < "${CURDIR}"/deploy.sh

echo "--- sync image...";
sh "${CURDIR}"/sync_image.sh

echo "--- docker image prune..."
echo y|docker image prune