FROM golang:alpine as builder
# RUN apt install -y git make
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && apk update && apk --no-cache add make git 
LABEL anther="cjun"
WORKDIR /tinyurl
COPY . /tinyurl
ENV GO111MODULE=on 
ENV GOPROXY=https://goproxy.cn,direct
CMD [ "go mod download" ]
RUN make build

FROM alpine:latest AS runner
# 设置alpine 时间为上海时间
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && apk update && apk --no-cache add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata
WORKDIR /app
COPY --from=builder /tinyurl/bin/tinyurl /app/tinyurl
COPY --from=builder /tinyurl/dist/ /app/dist/
COPY --from=builder /tinyurl/configs/ /app/configs/
EXPOSE 2830
CMD ["./tinyurl"]