FROM golang:latest as builder
LABEL anther="cjun"
WORKDIR /tinyurl
COPY . /tinyurl
ENV GO111MODULE=on 
ENV GOPROXY=https://goproxy.cn,direct
CMD [ "go mod download" ]
RUN make build

FROM alpine:latest AS runner
# 设置alpine 时间为上海时间
RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata
WORKDIR /app
COPY --from=builder /tinyurl/bin/tinyurl /app/tinyurl
EXPOSE 2830
CMD ["./tinyurl"]