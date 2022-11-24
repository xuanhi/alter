FROM alpine:3.17.0
MAINTAINER 1757314827@qq.com
ENV AppID=cli_a3f3305dc1f9d013 Secret=Ud5TZ38Sh6u1XUBfd4xOhhqClgoFf68l
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
&& echo "Asia/shanghai" > /etc/timezone \
&& apk add --no-cache ca-certificates tzdata
WORKDIR /feishu
EXPOSE 5001
COPY feishu .
CMD ["./feishu"]
