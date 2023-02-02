FROM alpine:3.17.0
MAINTAINER 1757314827@qq.com
ENV AppID=*** Secret=***
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
&& echo "Asia/shanghai" > /etc/timezone \
&& apk add --no-cache ca-certificates tzdata
WORKDIR /feishu
EXPOSE 5001
COPY feishu .
CMD ["./feishu"]
