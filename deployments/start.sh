#!/usr/bin/env sh
CURDIR=$(cd "$(dirname "$0")"; pwd)

echo "> 重启ssl证书监控..."
for var in $(ps -aux|grep ssl/watch.sh|grep -v "grep"|awk '{print $2}')
do
  if [ ! -n "$var" ]; then
      kill -9 $var
  fi
done
nohup $CURDIR/ssl/watch.sh >> $CURDIR/ssl/watch.log 2>&1 &
