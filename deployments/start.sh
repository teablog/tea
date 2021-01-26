#!/usr/bin/env sh
CURDIR=$(cd "$(dirname "$0")"; pwd)

echo "> 重启ssl证书架空..."
ps -aux|grep ssl/watch.sh|grep -v "grep"|awk '{print $2}'|xargs kill -9
nohup $CURDIR/ssl/watch.sh >> $CURDIR/ssl/watch.log 2>&1 &
