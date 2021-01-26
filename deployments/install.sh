#!/usr/bin/env sh

CURDIR=$(cd "$(dirname "$0")"; pwd)

echo "安装 acme.sh"
sh "${CURDIR}/ssl/acme.sh" --install

echo "安装 epel-release rsync"
yum install -y epel-release rsync

echo "安装 inotify-tools";
yum install -y inotify-tools