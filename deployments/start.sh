#!/usr/bin/env sh
CURDIR=$(cd "$(dirname "$0")"; pwd)

echo "启动... ssl证书监控";
nohup sh $CURDIR/ssl/watch.sh >> $CURDIR/ssl/watch.log
