#!/bin/sh
CURDIR=$(cd "$(dirname "$0")"; pwd)
# 监视的文件或目录
filename="/root/.acme.sh/douyacun.com/douyacun.com.cer"
# 监视发现有增、删、改时执行的脚本
script="${CURDIR}/restart.sh"
inotifywait -mrq --format '%e' --event create,delete,modify  $filename | while read event
do
    case $event in MODIFY|CREATE|DELETE) bash $script ;;
    esac
done