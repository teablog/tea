FROM alpine:latest
RUN  echo 'http://mirrors.ustc.edu.cn/alpine/v3.5/main' > /etc/apk/repositories \
    && echo 'http://mirrors.ustc.edu.cn/alpine/v3.5/community' >>/etc/apk/repositories \
&& apk update && apk add tzdata \
&& ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
&& echo "Asia/Shanghai" > /etc/timezone
RUN apk add --no-cache --upgrade bash
