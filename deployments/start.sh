#!/usr/bin/env sh
CURDIR=$(cd "$(dirname "$0")"; pwd)

echo "> restart ssl cert watch..."
for var in $(ps -aux|grep watch.sh|grep -v "grep"|awk '{print $2}')
do
    echo "> kill -9 ${var}"
    kill -9 $var
done
nohup $CURDIR/ssl/watch.sh >> $CURDIR/ssl/watch.log 2>&1 &