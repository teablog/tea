#!/usr/bin/env bash
# 禁止部署配置文件
rsync -azv ${pwd}/ douyacun:${deploy_dir}/
