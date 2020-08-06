#!/usr/bin/env bash

sudo docker pull registry.cn-hangzhou.aliyuncs.com/douyacun/tea:latest
# shellcheck disable=SC2164
pushd /data/web/tea/deployments/tea
sudo docker-compose up --force-recreate -d
# shellcheck disable=SC2164
popd
echo y|sudo docker image prune