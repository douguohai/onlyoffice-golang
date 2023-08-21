FROM golang:alpine as builder

# 修改使用使用国内代理, 否则会很慢
RUN set -ex \
&& go env -w GO111MODULE=on \
&& go env -w GOPROXY=https://goproxy.cn,direct && mkdir /app

# 在镜像中创建项目目录
ADD . /app

WORKDIR /app

# 创建项目的可执行文件web-server
RUN go build -o web-server *.go


FROM alpine:latest as prod

WORKDIR /app

COPY --from=0 /app/web-server /app
COPY --from=0 /app/views /app/views
COPY --from=0 /app/conf /app/conf
COPY --from=0 /app/static /app/static
EXPOSE 8080

ENV documentServer http://192.168.1.123:30080
ENV serverUrl http://192.168.1.123:18080
ENV wsServer ws://192.168.1.123:18080/ws
ENV dbHost 192.168.1.183
ENV dbName doc
ENV dbPort 5432
ENV dbUser postgres
ENV dbPassword woshidoudou

CMD ["/app/web-server"]