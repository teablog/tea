#!/usr/bin/env sh
CURDIR=$(cd "$(dirname "$0")"; pwd)

echo "> 启动ssl证书监控";
nohup $CURDIR/ssl/watch.sh >> $CURDIR/ssl/watch.log 2>&1 &
